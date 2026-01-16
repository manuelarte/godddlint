package valueobject

import (
	"go/ast"

	"github.com/manuelarte/godddlint/internal/astutils"
	"github.com/manuelarte/godddlint/internal/model"
	"github.com/manuelarte/godddlint/rules"
)

func NewRuleChecker() model.RuleChecker {
	return model.NewRuleChecker([]model.Rule{
		nonPointerReceivers{
			model.CommentRuleEnabler{RuleCode: rules.NonPointerReceivers.Code},
		},
		immutable{
			model.CommentRuleEnabler{RuleCode: rules.Immutable.Code},
		},
		defensiveCopy{},
	})
}

// NewDefinition creates a value object definition if the type contains the comment //godddlint:valueObject.
func NewDefinition(spec *ast.TypeSpec, doc *ast.CommentGroup) (*model.Definition, bool) {
	return &model.Definition{
		TypeSpec: spec,
		Doc:      doc,
	}, astutils.CommentHasPrefix(doc, "//godddlint:valueObject")
}
