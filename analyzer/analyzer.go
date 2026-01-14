package analyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
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
		(*ast.TypeSpec)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		switch n.(type) {
		case *ast.File:
			// TODO

		case *ast.FuncDecl:
			// TODO

		case *ast.TypeSpec:
			// TODO
		}
	})

	//nolint:nilnil //any, error
	return nil, nil
}
