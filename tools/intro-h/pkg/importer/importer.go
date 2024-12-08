package importer

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"slices"

	"golang.org/x/tools/go/analysis"
)

func makeFileNameImportsMap(pass *analysis.Pass) map[string]*ast.GenDecl {
	funcFileImports := make(map[string]*ast.GenDecl, 1)

	for _, f := range pass.Files {
		for _, dec := range f.Decls {
			fileName := pass.Fset.File(dec.Pos()).Name()
			switch dec := dec.(type) {
			case *ast.GenDecl:
				if dec.Tok == token.IMPORT {
					funcFileImports[fileName] = dec
				}
			}
		}
	}

	return funcFileImports
}

type ImportInfo struct {
	PkgPath string
	PkgName *string
}

func AddImports(pass *analysis.Pass, importInfos map[string][]ImportInfo) {
	fileImports := makeFileNameImportsMap(pass)
	for fileName, importInfos := range importInfos {
		importDecl, ok := fileImports[fileName]
		if !ok {
			continue
		}

		newImportSpecs := make([]ast.Spec, 0, len(importDecl.Specs))
		newImportSpecs = append(newImportSpecs, importDecl.Specs...)

		for _, importInfo := range importInfos {
			// すでにimportされていたら追加しない
			if slices.ContainsFunc(importDecl.Specs, func(spec ast.Spec) bool {
				importSpec, ok := spec.(*ast.ImportSpec)
				return ok && importSpec.Path.Value == fmt.Sprintf("%q", importInfo.PkgPath)
			}) {
				continue
			}

			var name *ast.Ident
			if importInfo.PkgName != nil {
				name = &ast.Ident{
					Name: *importInfo.PkgName,
				}
			}
			newImportSpecs = append(newImportSpecs, &ast.ImportSpec{
				Name: name,
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf("%q", importInfo.PkgPath),
				},
			})
		}

		newImportDecl := &ast.GenDecl{
			Doc:    importDecl.Doc,
			Tok:    importDecl.Tok,
			TokPos: importDecl.TokPos,
			Specs:  newImportSpecs,
		}

		buf := bytes.Buffer{}
		err := format.Node(&buf, pass.Fset, newImportDecl)
		if err != nil {
			return
		}

		pass.Report(analysis.Diagnostic{
			Pos:     newImportDecl.Pos(),
			Message: "change import",
			SuggestedFixes: []analysis.SuggestedFix{
				{
					TextEdits: []analysis.TextEdit{
						{
							Pos:     importDecl.Pos(),
							End:     importDecl.End(),
							NewText: buf.Bytes(),
						},
					},
				},
			},
		})

	}
}
