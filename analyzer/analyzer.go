package analyzer

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"

	"github.com/manuelarte/godddlint/internal/valueobject"
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

//nolint:gocognit // Refactor later
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

	valueObjectChecker := valueobject.NewChecker()
	valueObjectDefinitions := make(map[string]*valueobject.Definition)

	var funcDecls []*ast.FuncDecl

	insp.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.File:
			// TODO

		case *ast.FuncDecl:
			if n.Recv != nil && len(n.Recv.List) == 1 {
				funcDecls = append(funcDecls, n)
			}

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

				checker, ok := valueobject.NewDefinition(typeSpec, doc)
				if !ok {
					continue
				}

				valueObjectDefinitions[typeSpec.Name.Name] = checker
			}
		}
	})

	for _, funcDecl := range funcDecls {
		rcvName, hasRcvName := getRcvName(funcDecl.Recv.List[0].Type)
		if !hasRcvName {
			continue
		}

		definition, ok := valueObjectDefinitions[rcvName]
		if !ok {
			continue
		}

		definition.AddMethod(funcDecl)
	}

	for _, voDefinition := range valueObjectDefinitions {
		diags := valueObjectChecker.Check(*voDefinition)
		for _, diag := range diags {
			pass.Report(diag)
		}
	}

	//nolint:nilnil //any, error
	return nil, nil
}

func getRcvName(expr ast.Expr) (string, bool) {
	switch expr := expr.(type) {
	case *ast.Ident:
		return expr.Name, true
	case *ast.StarExpr:
		return getRcvName(expr.X)
	default:
		return "", false
	}
}
