package entity

import (
	"go/ast"

	"github.com/manuelarte/godddlint/internal/astutils"
	"github.com/manuelarte/godddlint/internal/model"
)

func NewRuleChecker() model.RuleChecker {
	return model.NewRuleChecker([]model.Rule{
		pointerReceivers{},
		customTypesOverPrimitives{},
		customDomainErrors{},
		unexportedFields{},
	})
}

// NewDefinition creates an entity definition if the type contains the comment //godddlint:entity.
func NewDefinition(spec *ast.TypeSpec, doc *ast.CommentGroup) (*model.Definition, bool) {
	return &model.Definition{
		TypeSpec: spec,
		Doc:      doc,
	}, astutils.CommentHasPrefix(doc, "//godddlint:entity")
}
