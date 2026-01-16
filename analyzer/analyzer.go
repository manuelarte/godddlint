package analyzer

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"

	"github.com/manuelarte/godddlint/internal/astutils"
	"github.com/manuelarte/godddlint/internal/entity"
	"github.com/manuelarte/godddlint/internal/model"
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

// run The run function analyze the files per package and will keep track of the
// functions that can be a constructor, the methods and all the structs that
// are annotated with valueObject or entity directives.
// Then it will link the constructor and methods found to the ddd structs.
// Once we have all the constructor and methods for a particular struct,
// we will apply the opinionated rules to them.
//
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

	valueObjectChecker := valueobject.NewRuleChecker()
	valueObjectDefinitions := make(map[string]*model.Definition)

	entitiesChecker := entity.NewRuleChecker()
	entitiesDefinitions := make(map[string]*model.Definition)

	possibleConstructorDecls := make([]*ast.FuncDecl, 0)
	methodsDecls := make([]*ast.FuncDecl, 0)

	insp.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.FuncDecl:
			if n.Recv == nil {
				if astutils.IsPotentialConstructor(n) {
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

				if voDefinition, okVo := valueobject.NewDefinition(typeSpec, doc); okVo {
					valueObjectDefinitions[typeSpec.Name.Name] = voDefinition
				}

				if eDefinition, okE := entity.NewDefinition(typeSpec, doc); okE {
					entitiesDefinitions[typeSpec.Name.Name] = eDefinition
				}
			}
		}
	})

	for _, possibleConstructor := range possibleConstructorDecls {
		structIdent, isIdent := possibleConstructor.Type.Results.List[0].Type.(*ast.Ident)
		if !isIdent {
			continue
		}

		if definition, ok := valueObjectDefinitions[structIdent.Name]; ok {
			definition.AddConstructor(possibleConstructor)
		}

		if definition, ok := entitiesDefinitions[structIdent.Name]; ok {
			definition.AddConstructor(possibleConstructor)
		}
	}

	for _, methodDecl := range methodsDecls {
		rcvName, hasRcvName := astutils.GetRcvName(methodDecl.Recv.List[0].Type)
		if !hasRcvName {
			continue
		}

		if definition, ok := valueObjectDefinitions[rcvName]; ok {
			definition.AddMethod(methodDecl)

			continue
		}

		if definition, ok := entitiesDefinitions[rcvName]; ok {
			definition.AddMethod(methodDecl)

			continue
		}
	}

	for _, voDefinition := range valueObjectDefinitions {
		diags := valueObjectChecker.Check(voDefinition)
		for _, diag := range diags {
			pass.Report(diag)
		}
	}

	for _, eDefinition := range entitiesDefinitions {
		diags := entitiesChecker.Check(eDefinition)
		for _, diag := range diags {
			pass.Report(diag)
		}
	}

	//nolint:nilnil //any, error
	return nil, nil
}
