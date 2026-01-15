package entity

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"

	"github.com/manuelarte/godddlint/internal/astutils"
	"github.com/manuelarte/godddlint/internal/model"
)

var _ model.Rule = new(pointerReceivers)

type (
	// Rule that checks that Entities use pointer receivers.
	pointerReceivers          struct{}
	customTypesOverPrimitives struct{}
	customDomainErrors        struct{}
)

func (r pointerReceivers) Apply(d *model.Definition) []analysis.Diagnostic {
	allDiag := make([]analysis.Diagnostic, 0)

	for _, m := range d.Methods {
		if ident, ok := m.Recv.List[0].Type.(*ast.Ident); ok {
			metadata := r.Metadata()
			message := fmt.Sprintf("%s: Entity's method not using pointer receiver", metadata.Code)
			diag := analysis.Diagnostic{
				Pos:      ident.Pos(),
				End:      ident.End(),
				Category: metadata.Name,
				Message:  message,
				URL:      metadata.URL,
			}
			allDiag = append(allDiag, diag)
		}
	}

	return allDiag
}

func (r pointerReceivers) Metadata() model.RuleMetadata {
	return model.RuleMetadata{
		Code: "E001",
		Name: "Pointer Receivers",
	}
}

func (r customTypesOverPrimitives) Apply(d *model.Definition) []analysis.Diagnostic {
	allDiag := make([]analysis.Diagnostic, 0)

	metadata := r.Metadata()

	if st, ok := d.TypeSpec.Type.(*ast.StructType); ok {
		for _, f := range st.Fields.List {
			if astutils.IsPrimitiveType(f) {
				diag := analysis.Diagnostic{
					Pos:     f.Pos(),
					End:     f.End(),
					Message: fmt.Sprintf("%s: Prefer custom domain types to primitives", metadata.Code),
					URL:     metadata.URL,
				}
				allDiag = append(allDiag, diag)
			}
		}
	}

	return allDiag
}

func (r customTypesOverPrimitives) Metadata() model.RuleMetadata {
	return model.RuleMetadata{
		Code: "E003",
		Name: "Custom Types Over Primitives",
	}
}

func (r customDomainErrors) Apply(d *model.Definition) []analysis.Diagnostic {
	// TODO:
	return nil
}

func (r customDomainErrors) Metadata() model.RuleMetadata {
	return model.RuleMetadata{
		Code: "E004",
		Name: "Custom Domain Errors",
	}
}
