package valueobject

import (
	"go/ast"
	"slices"

	"golang.org/x/tools/go/analysis"
)

type (
	Definition struct {
		TypeSpec *ast.TypeSpec
		Doc      *ast.CommentGroup
		Methods  []*ast.FuncDecl
	}

	Checker struct {
		rules []rule
	}
)

func NewChecker() Checker {
	return Checker{
		rules: []rule{
			nonPointerReceivers{},
		},
	}
}

func (c Checker) Check(definition Definition) []analysis.Diagnostic {
	allDiag := make([]analysis.Diagnostic, 0)
	for _, r := range c.rules {
		allDiag = append(allDiag, r.Apply(definition)...)
	}

	return allDiag
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
