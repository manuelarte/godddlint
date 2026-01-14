package entity

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"

	"github.com/manuelarte/godddlint/internal/model"
)

var _ model.Rule = new(pointerReceivers)

type (
	pointerReceivers struct{}
)

func (r pointerReceivers) Apply(d model.Definition) []analysis.Diagnostic {
	allDiag := make([]analysis.Diagnostic, 0)

	for _, m := range d.Methods {
		if se, ok := m.Recv.List[0].Type.(*ast.StarExpr); ok {
			metadata := r.Metadata()
			message := fmt.Sprintf("%s: %s", metadata.Name, metadata.Description)
			diag := analysis.Diagnostic{
				Pos:     se.Star,
				End:     se.End(),
				Message: message,
				URL:     metadata.URL,
			}
			allDiag = append(allDiag, diag)
		}
	}

	return allDiag
}

func (r pointerReceivers) Metadata() model.RuleMetadata {
	return model.RuleMetadata{
		Name:        "E001",
		Description: "Pointer Receivers",
	}
}
