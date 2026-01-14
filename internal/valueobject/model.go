package valueobject

import "go/ast"

type Checker struct{}

func New(spec *ast.TypeSpec, doc *ast.CommentGroup) (Checker, bool) {
	// TODO: check for comment //godddlint:valueObject
	return Checker{}, true
}
