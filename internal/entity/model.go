package entity

import (
	"go/ast"
	"slices"

	"github.com/manuelarte/godddlint/internal/model"
)

func NewChecker() model.Checker {
	return model.NewChecker([]model.Rule{
		pointerReceivers{},
		customTypesOverPrimitives{},
	})
}

// NewDefinition creates an entity definition if the type contains the comment //godddlint:entity.
func NewDefinition(spec *ast.TypeSpec, doc *ast.CommentGroup) (*model.Definition, bool) {
	return &model.Definition{
		TypeSpec: spec,
		Doc:      doc,
	}, commentContainsValueObject(doc)
}

func commentContainsValueObject(doc *ast.CommentGroup) bool {
	return doc != nil && slices.ContainsFunc(doc.List, func(c *ast.Comment) bool {
		return c.Text == "//godddlint:entity"
	})
}
