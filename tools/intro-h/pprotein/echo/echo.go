package echo

import (
	"flag"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"slices"
	"strings"

	"github.com/reactive-futakotamagawa/isu-isu-h/tools/intro-h/pkg/importer"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = "introtein is ..."

// Analyzer is ...
var Analyzer = &analysis.Analyzer{
	Name: "introtein_initialize",
	Doc:  doc,
	Run:  runInitialize,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
	Flags: *flag.NewFlagSet("introtein_initialize", flag.ExitOnError),
	// Flags: ,
}

var webhookURL string

func init() {
	Analyzer.Flags.StringVar(&webhookURL, "webhook", "https://example.com/api/group/collect", "webhook url")
}

func runInitialize(pass *analysis.Pass) (any, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.CallExpr:
			se, ok := n.Fun.(*ast.SelectorExpr)
			if !ok {
				return
			}

			obj := pass.TypesInfo.ObjectOf(se.Sel)
			if obj.Name() != "POST" || !strings.HasSuffix(obj.Pkg().Path(), "github.com/labstack/echo/v4") {
				return
			}

			recv := obj.(*types.Func).Type().(*types.Signature).Recv().Type().(*types.Pointer).Elem().(*types.Named).Obj()
			if !slices.Contains([]string{"Echo", "Group"}, recv.Name()) {
				return
			}

			path, ok := n.Args[0].(*ast.BasicLit)
			if !ok || path.Kind != token.STRING {
				return
			}

			if !strings.Contains(path.Value, "initialize") {
				return
			}

			handlerExpr := n.Args[1]
			var handlerBody *ast.BlockStmt
			switch handlerExpr := handlerExpr.(type) {
			case *ast.FuncLit:
				handlerBody = handlerExpr.Body
			case *ast.Ident:
				handlerBody = handlerExpr.Obj.Decl.(*ast.FuncDecl).Body
			default:
				return
			}

			code := fmt.Sprintf(`go func() {
	if _, err := http.Get("%s"); err != nil {
		log.Printf("failed to communicate with pprotein: %%v", err)
		}
}()`, webhookURL)

			pass.Report(analysis.Diagnostic{
				Pos:     handlerBody.Lbrace,
				Message: "initialize",
				SuggestedFixes: []analysis.SuggestedFix{
					{
						Message: "Send webhook",
						TextEdits: []analysis.TextEdit{
							{
								Pos:     handlerBody.Lbrace + 1,
								NewText: []byte(code),
							},
						},
					},
				},
			})

			importer.AddImports(pass, map[string][]importer.ImportInfo{
				pass.Fset.File(n.Pos()).Name(): {
					{
						PkgPath: "net/http",
					},
					{
						PkgPath: "log",
					},
				},
			})
		}
	})

	return nil, nil
}
