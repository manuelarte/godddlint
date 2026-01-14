package valueobject

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var _ rule = new(nonPointerReceivers)

type (
	ruleMetadata struct {
		Name, Description, URL string
	}

	rule interface {
		Apply(definition Definition) []analysis.Diagnostic
		metadata() ruleMetadata
	}

	nonPointerReceivers struct{}
)

func (r nonPointerReceivers) Apply(d Definition) []analysis.Diagnostic {
	allDiag := make([]analysis.Diagnostic, 0)

	for _, m := range d.Methods {
		if se, ok := m.Recv.List[0].Type.(*ast.StarExpr); ok {
			metadata := r.metadata()
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

func (r nonPointerReceivers) metadata() ruleMetadata {
	return ruleMetadata{
		Name:        "VO001",
		Description: "Non Pointer Receivers",
	}
}
