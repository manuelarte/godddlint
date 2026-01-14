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

	nonPointerReceivers      struct{}
	internalFieldsUnexported struct{}
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

func (r internalFieldsUnexported) Apply(d Definition) []analysis.Diagnostic {
	allDiag := make([]analysis.Diagnostic, 0)

	metadata := r.metadata()
	message := fmt.Sprintf("%s: %s", metadata.Name, metadata.Description)

	if st, ok := d.TypeSpec.Type.(*ast.StructType); ok {
		if len(d.Constructors) == 0 {
			diag := analysis.Diagnostic{
				Pos:     d.TypeSpec.Pos(),
				End:     d.TypeSpec.End(),
				Message: message,
				URL:     metadata.URL,
			}
			allDiag = append(allDiag, diag)
		}

		for _, f := range st.Fields.List {
			for _, n := range f.Names {
				if n.IsExported() {
					diag := analysis.Diagnostic{
						Pos:     n.Pos(),
						End:     n.End(),
						Message: message,
						URL:     metadata.URL,
					}
					allDiag = append(allDiag, diag)
				}
			}
		}
	}

	return allDiag
}

func (r internalFieldsUnexported) metadata() ruleMetadata {
	return ruleMetadata{
		Name:        "VOX001",
		Description: "Immutable",
	}
}
