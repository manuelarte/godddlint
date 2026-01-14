package valueObject

import "go/ast"

type Checker struct {
}

func New(*ast.TypeSpec) (Checker, bool) {
	// TODO: check for comment //godddlint:valueObject
	return Checker{}, true
}
