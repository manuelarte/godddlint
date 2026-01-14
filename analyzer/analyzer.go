package analyzer

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"

	"github.com/manuelarte/go-template/internal/valueObject"
)

func New() *analysis.Analyzer {
	g := godddlint{}

	return &analysis.Analyzer{
		Name:     "godddlint",
		Doc:      "checks domain structs honor best practices",
		URL:      "https://github.com/manuelarte/godddlint",
		Run:      g.run,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
}

type godddlint struct{}

func (g godddlint) run(pass *analysis.Pass) (any, error) {
	insp, found := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !found {
		//nolint:nilnil // impossible case.
		return nil, nil
	}

	nodeFilter := []ast.Node{
		(*ast.File)(nil),
		(*ast.FuncDecl)(nil),
		(*ast.GenDecl)(nil),
	}

	var valueObjectCheckers []valueObject.Checker

	insp.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.File:
			// TODO

		case *ast.FuncDecl:
			// TODO

		case *ast.GenDecl:
			if n.Tok != token.TYPE {
				return
			}
			for _, spec := range n.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				doc := typeSpec.Doc
				if doc == nil {
					doc = n.Doc
				}
				checker, ok := valueObject.New(typeSpec, doc)
				if !ok {
					continue
				}

				valueObjectCheckers = append(valueObjectCheckers, checker)
			}
		}
	})

	//nolint:nilnil //any, error
	return nil, nil
}
