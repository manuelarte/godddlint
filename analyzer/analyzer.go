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
		(*ast.FuncDecl)(nil),
		(*ast.GenDecl)(nil),
	}

	valueObjectChecker := valueobject.NewChecker()
	valueObjectDefinitions := make(map[string]*valueobject.Definition)

	possibleConstructorDecls := make([]*ast.FuncDecl, 0)
	methodsDecls := make([]*ast.FuncDecl, 0)

	insp.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.FuncDecl:
			if n.Recv == nil {
				if potentialConstructor(n) {
					possibleConstructorDecls = append(possibleConstructorDecls, n)

					return
				}
			}

			if n.Recv != nil && len(n.Recv.List) == 1 {
				methodsDecls = append(methodsDecls, n)

				return
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

	for _, possibleConstructor := range possibleConstructorDecls {
		structIdent, isIdent := possibleConstructor.Type.Results.List[0].Type.(*ast.Ident)
		if !isIdent {
			continue
		}

		definition, ok := valueObjectDefinitions[structIdent.Name]
		if !ok {
			continue
		}

		definition.AddConstructor(possibleConstructor)
	}

	for _, methodDecl := range methodsDecls {
		rcvName, hasRcvName := getRcvName(methodDecl.Recv.List[0].Type)
		if !hasRcvName {
			continue
		}

		definition, ok := valueObjectDefinitions[rcvName]
		if !ok {
			continue
		}

		definition.AddMethod(methodDecl)
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
