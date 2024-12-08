package integrate

import (
	"go/ast"

	"github.com/reactive-futakotamagawa/isu-isu-h/tools/intro-h/pkg/importer"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = "introtein is ..."

// IntegrateMain is ...
var IntegrateMain = &analysis.Analyzer{
	Name: "integrate_main",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func run(pass *analysis.Pass) (any, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
		(*ast.FuncLit)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.FuncDecl:
			if n.Name.Name != "main" {
				return
			}
			addIntegrateToMain(pass, n)
		}
	})

	return nil, nil
}

func addIntegrateToMain(pass *analysis.Pass, n *ast.FuncDecl) {
	f := pass.Fset.File(n.Pos())

	pass.Report(analysis.Diagnostic{
		Pos:     n.Pos(),
		Message: "main",
		SuggestedFixes: []analysis.SuggestedFix{
			{
				Message: "Add a comment",
				TextEdits: []analysis.TextEdit{
					{
						Pos:     n.Body.Lbrace + 2,
						End:     n.Body.Lbrace + 2,
						NewText: []byte("go standalone.Integrate(\":8888\")\n"),
					},
				},
			},
		},
	})

	importer.AddImports(pass, map[string][]importer.ImportInfo{
		f.Name(): {
			{
				PkgPath: "github.com/kaz/pprotein/integration/standalone",
			},
		},
	})
}
