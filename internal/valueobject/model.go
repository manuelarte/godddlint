package valueobject

import (
	"go/ast"
	"slices"
)

type (
	Definition struct {
		TypeSpec *ast.TypeSpec
		Doc      *ast.CommentGroup
		Methods  []*ast.FuncDecl
	}

	Checker struct {
		rules []Rule
	}
)

func NewChecker() Checker {
	return Checker{
		rules: []Rule{
			NonPointerReceivers{},
		},
	}
}

func (c Checker) Check(definition Definition) {
	for _, rule := range c.rules {
		// TODO diagnosis error too
		rule.Apply(definition)
	}
}

// NewDefinition creates a value object checker if the type contains the comment //godddlint:valueObject.
func NewDefinition(spec *ast.TypeSpec, doc *ast.CommentGroup) (*Definition, bool) {
	return &Definition{
		TypeSpec: spec,
		Doc:      doc,
	}, commentContainsValueObject(doc)
}

func (d *Definition) AddMethod(method *ast.FuncDecl) {
	d.Methods = append(d.Methods, method)
}

func commentContainsValueObject(doc *ast.CommentGroup) bool {
	return doc != nil && slices.ContainsFunc(doc.List, func(c *ast.Comment) bool {
		return c.Text == "//godddlint:valueObject"
	})
}
