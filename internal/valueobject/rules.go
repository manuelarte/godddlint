package valueobject

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"

	"github.com/manuelarte/godddlint/internal/model"
)

var (
	_ model.Rule = new(nonPointerReceivers)
	_ model.Rule = new(immutable)
)

type (
	nonPointerReceivers struct{}
	immutable           struct{}
)

func (r nonPointerReceivers) Apply(d model.Definition) []analysis.Diagnostic {
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

func (r nonPointerReceivers) Metadata() model.RuleMetadata {
	return model.RuleMetadata{
		Name:        "VO001",
		Description: "Non Pointer Receivers",
	}
}

func (r immutable) Apply(d model.Definition) []analysis.Diagnostic {
	allDiag := make([]analysis.Diagnostic, 0)

	metadata := r.Metadata()
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

func (r immutable) Metadata() model.RuleMetadata {
	return model.RuleMetadata{
		Name:        "VOX001",
		Description: "Immutable",
	}
}
