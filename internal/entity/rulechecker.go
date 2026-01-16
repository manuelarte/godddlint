package entity

import (
	"go/ast"

	"github.com/manuelarte/godddlint/internal/astutils"
	"github.com/manuelarte/godddlint/internal/model"
	"github.com/manuelarte/godddlint/rules"
)

func NewRuleChecker() model.RuleChecker {
	return model.NewRuleChecker([]model.Rule{
		pointerReceivers{
			model.CommentRuleEnabler{RuleCode: rules.PointerReceivers.Code},
		},
		customTypesOverPrimitives{
			model.CommentRuleEnabler{RuleCode: rules.CustomTypesOverPrimitives.Code},
		},
		customDomainErrors{
			model.CommentRuleEnabler{RuleCode: rules.CustomDomainErrors.Code},
		},
		unexportedFields{
			model.CommentRuleEnabler{RuleCode: rules.UnexportedFields.Code},
		},
	})
}

// NewDefinition creates an entity definition if the type contains the comment //godddlint:entity.
func NewDefinition(spec *ast.TypeSpec, doc *ast.CommentGroup) (*model.Definition, bool) {
	return &model.Definition{
		TypeSpec: spec,
		Doc:      doc,
	}, astutils.CommentHasPrefix(doc, "//godddlint:entity")
}
