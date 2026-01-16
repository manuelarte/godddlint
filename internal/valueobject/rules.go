package valueobject

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"

	"github.com/manuelarte/godddlint/internal/astutils"
	"github.com/manuelarte/godddlint/internal/model"
	"github.com/manuelarte/godddlint/rules"
)

var (
	_ model.Rule = new(nonPointerReceivers)
	_ model.Rule = new(immutable)
)

type (
	// Rule that checks that value objects use non-pointer receivers.
	nonPointerReceivers struct {
		ruleEnableChecker model.RuleEnablerChecker
	}
	// Rule that checks that value objects have constructor(s) and unexported fields.
	immutable struct {
		ruleEnableChecker model.RuleEnablerChecker
	}
	// Rule that checks that map/slices are defensively copied.
	defensiveCopy struct{}
)

func (r nonPointerReceivers) Apply(d *model.Definition) []analysis.Diagnostic {
	if !r.ruleEnableChecker.IsEnabled(d.Doc) {
		return nil
	}

	allDiag := make([]analysis.Diagnostic, 0)

	for _, m := range d.Methods {
		if astutils.CommentHasPrefix(m.Doc, "//godddlint:disable:VO001") {
			continue
		}

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

func (r nonPointerReceivers) Metadata() rules.RuleMetadata {
	return rules.NonPointerReceivers
}

func (r immutable) Apply(d *model.Definition) []analysis.Diagnostic {
	if !r.ruleEnableChecker.IsEnabled(d.Doc) {
		return nil
	}

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
				if !r.ruleEnableChecker.IsEnabled(f.Doc) {
					return nil
				}

				if n.IsExported() {
					diag := analysis.Diagnostic{
						Pos:      n.Pos(),
						End:      n.End(),
						Category: metadata.Name,
						Message:  fmt.Sprintf("%s: Value Object's field is exported", metadata.Code),
						URL:      metadata.URL,
					}
					allDiag = append(allDiag, diag)
				}
			}
		}
	}

	return allDiag
}

func (r immutable) Metadata() rules.RuleMetadata {
	return rules.Immutable
}

//nolint:gocognit,nestif // Refactor later
func (r defensiveCopy) Apply(d *model.Definition) []analysis.Diagnostic {
	structType, isStructType := d.TypeSpec.Type.(*ast.StructType)
	if !isStructType {
		return nil
	}

	allDiag := make([]analysis.Diagnostic, 0)
	for _, constructor := range d.Constructors {
		mapOrSliceFieldNames := make(map[string]struct{})

		for _, f := range astutils.GetFieldsThatMustDefensiveCopy(structType) {
			for _, name := range f.Names {
				mapOrSliceFieldNames[name.Name] = struct{}{}
			}
		}

		ast.Inspect(constructor.Body, func(n ast.Node) bool {
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

					if _, containsKey := mapOrSliceFieldNames[keyIdent.Name]; containsKey &&
						r.isConstructorParam(valueIdent, constructor) {
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

func (r defensiveCopy) Metadata() rules.RuleMetadata {
	return rules.DefensiveCopy
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
