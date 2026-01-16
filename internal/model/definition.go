package model

import (
	"go/ast"
)

type (
	Definition struct {
		TypeSpec *ast.TypeSpec
		Doc      *ast.CommentGroup
		// Constructors are functions that return either the struct or the struct and an error.
		Constructors []*ast.FuncDecl
		// Methods are the methods for this entity.
		Methods []*ast.FuncDecl
	}
)

func (d *Definition) AddConstructor(constructor *ast.FuncDecl) {
	d.Constructors = append(d.Constructors, constructor)
}

func (d *Definition) AddMethod(method *ast.FuncDecl) {
	d.Methods = append(d.Methods, method)
}
