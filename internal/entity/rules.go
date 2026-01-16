package entity

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"

	"github.com/manuelarte/godddlint/internal/astutils"
	"github.com/manuelarte/godddlint/internal/model"
	"github.com/manuelarte/godddlint/rules"
)

var _ model.Rule = new(pointerReceivers)

type (
	// Rule that checks that Entities use pointer receivers.
	pointerReceivers struct {
		ruleEnableChecker model.RuleEnablerChecker
	}
	customTypesOverPrimitives struct {
		ruleEnableChecker model.RuleEnablerChecker
	}
	customDomainErrors struct {
		ruleEnableChecker model.RuleEnablerChecker
	}
	unexportedFields struct {
		ruleEnableChecker model.RuleEnablerChecker
	}
)

func (r pointerReceivers) Apply(d *model.Definition) []analysis.Diagnostic {
	if !r.ruleEnableChecker.IsEnabled(d.Doc) {
		return nil
	}

	allDiag := make([]analysis.Diagnostic, 0)

	for _, m := range d.Methods {
		if !r.ruleEnableChecker.IsEnabled(m.Doc) {
			return nil
		}

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

func (r pointerReceivers) Metadata() rules.RuleMetadata {
	return rules.PointerReceivers
}

func (r customTypesOverPrimitives) Apply(d *model.Definition) []analysis.Diagnostic {
	if !r.ruleEnableChecker.IsEnabled(d.Doc) {
		return nil
	}

	allDiag := make([]analysis.Diagnostic, 0)

	metadata := r.Metadata()

	if st, ok := d.TypeSpec.Type.(*ast.StructType); ok {
		for _, f := range st.Fields.List {
			if !r.ruleEnableChecker.IsEnabled(f.Doc) {
				continue
			}

			if astutils.IsPrimitiveType(f) {
				diag := analysis.Diagnostic{
					Pos:      f.Pos(),
					End:      f.End(),
					Category: metadata.Name,
					Message:  fmt.Sprintf("%s: Prefer custom domain types to primitives", metadata.Code),
					URL:      metadata.URL,
				}
				allDiag = append(allDiag, diag)
			}
		}
	}

	return allDiag
}

func (r customTypesOverPrimitives) Metadata() rules.RuleMetadata {
	return rules.CustomTypesOverPrimitives
}

func (r customDomainErrors) Apply(d *model.Definition) []analysis.Diagnostic {
	if !r.ruleEnableChecker.IsEnabled(d.Doc) {
		return nil
	}

	allDiag := make([]analysis.Diagnostic, 0)
	metadata := r.Metadata()

	for _, m := range d.Methods {
		if !r.ruleEnableChecker.IsEnabled(m.Doc) {
			return nil
		}

		errorPos, hasError := astutils.FuncResultError(m)
		if !hasError {
			continue
		}

		ast.Inspect(m.Body, func(n ast.Node) bool {
			//nolint:gocritic // will cover more cases
			switch n := n.(type) {
			//nolint:gocritic // will cover more cases
			case *ast.ReturnStmt:
				errorReturn := n.Results[errorPos]
				switch rt := errorReturn.(type) {
				case *ast.CallExpr:
					if astutils.IsErrorsNewFun(rt) {
						diag := analysis.Diagnostic{
							Pos:      rt.Pos(),
							End:      rt.End(),
							Category: metadata.Name,
							Message:  fmt.Sprintf("%s: Returning errors.New instead of a domain error", metadata.Code),
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

func (r customDomainErrors) Metadata() rules.RuleMetadata {
	return rules.CustomDomainErrors
}

func (r unexportedFields) Apply(d *model.Definition) []analysis.Diagnostic {
	if !r.ruleEnableChecker.IsEnabled(d.Doc) {
		return nil
	}

	allDiag := make([]analysis.Diagnostic, 0)

	if st, ok := d.TypeSpec.Type.(*ast.StructType); ok {
		metadata := r.Metadata()

		for _, f := range st.Fields.List {
			if !r.ruleEnableChecker.IsEnabled(f.Doc) {
				continue
			}

			for _, n := range f.Names {
				if n.IsExported() {
					diag := analysis.Diagnostic{
						Pos:      n.Pos(),
						End:      n.End(),
						Category: metadata.Name,
						Message:  fmt.Sprintf("%s: Entity's field is exported", metadata.Code),
						URL:      metadata.URL,
					}
					allDiag = append(allDiag, diag)
				}
			}
		}
	}

	return allDiag
}

func (r unexportedFields) Metadata() rules.RuleMetadata {
	return rules.UnexportedFields
}
