package valueobject

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"

	"github.com/manuelarte/godddlint/internal/astutils"
	"github.com/manuelarte/godddlint/internal/model"
)

var (
	_ model.Rule = new(nonPointerReceivers)
	_ model.Rule = new(immutable)
)

type (
	// Rule that checks that value objects use non-pointer receivers.
	nonPointerReceivers struct{}
	// Rule that checks that value objects have constructor(s) and unexported fields.
	immutable struct{}
	// Rule that checks that map/slices are defensively copied.
	defensiveCopy struct{}
)

func (r nonPointerReceivers) Apply(d *model.Definition) []analysis.Diagnostic {
	allDiag := make([]analysis.Diagnostic, 0)

	for _, m := range d.Methods {
		if se, ok := m.Recv.List[0].Type.(*ast.StarExpr); ok {
			metadata := r.Metadata()
			diag := analysis.Diagnostic{
				Pos:      se.Star,
				End:      se.End(),
				Category: metadata.Name,
				Message:  fmt.Sprintf("%s: Value Object's method using a pointer receiver", metadata.Code),
				URL:      metadata.URL,
			}
			allDiag = append(allDiag, diag)
		}
	}

	return allDiag
}

func (r nonPointerReceivers) Metadata() model.RuleMetadata {
	return model.RuleMetadata{
		Code: "VO001",
		Name: "Non Pointer Receivers",
	}
}

func (r immutable) Apply(d *model.Definition) []analysis.Diagnostic {
	allDiag := make([]analysis.Diagnostic, 0)

	metadata := r.Metadata()

	if st, ok := d.TypeSpec.Type.(*ast.StructType); ok {
		if len(d.Constructors) == 0 {
			diag := analysis.Diagnostic{
				Pos:      d.TypeSpec.Pos(),
				End:      d.TypeSpec.End(),
				Category: metadata.Name,
				Message:  fmt.Sprintf("%s: Constructor for Value Object not found", metadata.Code),
				URL:      metadata.URL,
			}
			allDiag = append(allDiag, diag)
		}

		for _, f := range st.Fields.List {
			for _, n := range f.Names {
				if n.IsExported() {
					diag := analysis.Diagnostic{
						Pos:     n.Pos(),
						End:     n.End(),
						Message: fmt.Sprintf("%s: Value Object's field is exported", metadata.Code),
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
		Code: "VOX001",
		Name: "Immutable",
	}
}

//nolint:gocognit,nestif // Refactor later
func (r defensiveCopy) Apply(d *model.Definition) []analysis.Diagnostic {
	allDiag := make([]analysis.Diagnostic, 0)

	for _, constructor := range d.Constructors {
		ast.Inspect(constructor.Body, func(n ast.Node) bool {
			// Check for direct assignment: v.Field = param
			if assignStmt, ok := n.(*ast.AssignStmt); ok {
				for i, lhs := range assignStmt.Lhs {
					selectorExpr, isSelector := lhs.(*ast.SelectorExpr)
					if !isSelector {
						continue
					}

					rhs := assignStmt.Rhs[i]

					rhsIdent, isIdent := rhs.(*ast.Ident)
					if !isIdent {
						continue
					}

					if r.isMapOrSliceField(selectorExpr, d) && r.isConstructorParam(rhsIdent, constructor) {
						metadata := r.Metadata()
						diag := analysis.Diagnostic{
							Pos:      assignStmt.Pos(),
							End:      assignStmt.End(),
							Category: metadata.Name,
							Message:  fmt.Sprintf("%s: Maps/Slices Not Defensive Copied", metadata.Code),
							URL:      metadata.URL,
						}
						allDiag = append(allDiag, diag)
					}
				}
			}

			// Check for struct initialization: MyStruct{Field: param}
			if compLit, isCompLit := n.(*ast.CompositeLit); isCompLit {
				if !r.isTargetStruct(compLit, d) {
					return true
				}

				for _, elt := range compLit.Elts {
					kv, ok := elt.(*ast.KeyValueExpr)
					if !ok {
						continue
					}

					keyIdent, ok := kv.Key.(*ast.Ident)
					if !ok {
						continue
					}

					valueIdent, ok := kv.Value.(*ast.Ident)
					if !ok {
						continue
					}

					if r.isMapOrSliceFieldByName(keyIdent.Name, d) && r.isConstructorParam(valueIdent, constructor) {
						metadata := r.Metadata()
						diag := analysis.Diagnostic{
							Pos:      kv.Pos(),
							End:      kv.End(),
							Category: metadata.Name,
							Message:  fmt.Sprintf("%s: Maps/Slices Not Defensive Copied", metadata.Code),
							URL:      metadata.URL,
						}
						allDiag = append(allDiag, diag)
					}
				}
			}

			return true
		})
	}

	return allDiag
}

func (r defensiveCopy) Metadata() model.RuleMetadata {
	return model.RuleMetadata{
		Code: "VOX002",
		Name: "Maps/Slices Not Defensive Copied",
	}
}

func (r defensiveCopy) isTargetStruct(compLit *ast.CompositeLit, d *model.Definition) bool {
	if compLit.Type == nil {
		return false
	}

	ident, ok := compLit.Type.(*ast.Ident)
	if !ok {
		return false
	}

	return ident.Name == d.TypeSpec.Name.Name
}

func (r defensiveCopy) isMapOrSliceFieldByName(fieldName string, d *model.Definition) bool {
	structType, ok := d.TypeSpec.Type.(*ast.StructType)
	if !ok {
		return false
	}

	return astutils.IsMapOrSliceField(fieldName, structType)
}

func (r defensiveCopy) isMapOrSliceField(selector *ast.SelectorExpr, d *model.Definition) bool {
	return r.isMapOrSliceFieldByName(selector.Sel.Name, d)
}

func (r defensiveCopy) isConstructorParam(ident *ast.Ident, constructor *ast.FuncDecl) bool {
	for _, param := range constructor.Type.Params.List {
		for _, name := range param.Names {
			if name.Name == ident.Name {
				return true
			}
		}
	}

	return false
}
