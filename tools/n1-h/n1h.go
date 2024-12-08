package n1h

import (
	"go/ast"
	"go/token"
	"go/types"
	"slices"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/callgraph"
	"golang.org/x/tools/go/callgraph/static"
	"golang.org/x/tools/go/ssa"
)

const (
	doc = "n1-h is a tool to find N+1 problem with github.com/jmoiron/sqlx"
)

var (
	sqlxFuncs = []string{
		"Get", "GetContext",
		"Select", "SelectContext",
		"Query", "QueryContext",
		"QueryRow", "QueryRowContext",
		"Queryx", "QueryxContext",
		"QueryRowx", "QueryRowxContext",
		"Exec", "ExecContext",
		"NamedExec", "NamedExecContext",
		"NamedQuery", "NamedQueryContext",
	}
)

// Analyzer is ...
var Analyzer = &analysis.Analyzer{
	Name: "n1h",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		buildssa.Analyzer,
		inspect.Analyzer,
	},
}

func run(pass *analysis.Pass) (any, error) {
	s := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA)
	callGraph := static.CallGraph(s.Pkg.Prog)

	dbRetrieveCalls := make([]*ssa.Call, 0)
	dbRetrieveCallFuncs := make(map[*ssa.Function]struct{}, 0)

	for _, f := range s.SrcFuncs {
		for _, b := range f.Blocks {
			for _, instr := range b.Instrs {
				if call, ok := instr.(*ssa.Call); ok {
					if !checkSqlxCall(call) {
						continue
					}
					dbRetrieveCalls = append(dbRetrieveCalls, call)
					dbRetrieveCallFuncs[f] = struct{}{}
				}
			}
		}
	}

	dbRetrieveCallEdges := make(map[*callgraph.Edge]struct{}, len(dbRetrieveCallFuncs))
	walker := &callWalker{checkedFuncs: make(map[*ssa.Function]struct{}, len(dbRetrieveCallFuncs))}
	for f := range dbRetrieveCallFuncs {
		for _, edge := range walker.allAncestorCalls(f, callGraph) {
			dbRetrieveCallEdges[edge] = struct{}{}
		}
	}

	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	inspect.Preorder([]ast.Node{
		(*ast.ForStmt)(nil),
		(*ast.RangeStmt)(nil),
	}, func(n ast.Node) {
		var forStart, forEnd token.Pos
		switch n := n.(type) {
		case *ast.ForStmt:
			forStart, forEnd = n.Pos(), n.End()
		case *ast.RangeStmt:
			forStart, forEnd = n.Pos(), n.End()
		}

		for _, call := range dbRetrieveCalls {
			callPos := call.Pos()
			if forStart <= callPos && callPos <= forEnd {
				pass.Reportf(call.Pos(), "maybe N+1. `%v` in for loop", call.Call.Value)
			}
		}
		for edge := range dbRetrieveCallEdges {
			if edge.Pos() >= forStart && edge.Pos() <= forEnd {
				pass.Reportf(edge.Pos(), "maybe N+1. `%s` has more than 1 db call(s)", edge.Callee.Func.Name())
			}
		}
	})
	return nil, nil
}

func checkSqlxCall(call *ssa.Call) bool {
	if len(call.Call.Args) == 0 {
		return false
	}
	ptr, ok := call.Call.Args[0].Type().(*types.Pointer)
	if !ok {
		return false
	}
	named, ok := ptr.Elem().(*types.Named)
	if !ok {
		return false
	}

	dbField, _, _ := types.LookupFieldOrMethod(named.Obj().Type(), false, nil, "DB")
	txField, _, _ := types.LookupFieldOrMethod(named.Obj().Type(), false, nil, "Tx")
	isDB := named.Obj().Pkg().Path() == "database/sql" && named.Obj().Name() == "DB"
	if dbField == nil && txField == nil && !isDB {
		return false
	}

	if !slices.Contains(sqlxFuncs, call.Call.Value.Name()) {
		return false
	}
	return true
}

type callWalker struct {
	checkedFuncs map[*ssa.Function]struct{}
}

func (cw *callWalker) allAncestorCalls(f *ssa.Function, callGraph *callgraph.Graph) []*callgraph.Edge {
	// メモ化再帰
	if _, ok := cw.checkedFuncs[f]; ok {
		return nil
	}
	cw.checkedFuncs[f] = struct{}{}

	parents := make([]*callgraph.Edge, 0)
	node, ok := callGraph.Nodes[f]
	if !ok {
		return parents
	}
	for _, edge := range node.In {
		parents = append(parents, edge)
		parents = append(parents, cw.allAncestorCalls(edge.Caller.Func, callGraph)...)
	}
	return parents
}
