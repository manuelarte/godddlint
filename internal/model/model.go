package model

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"

	"github.com/manuelarte/godddlint/rules"
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

	Checker struct {
		rules []Rule
	}

	Rule interface {
		Apply(definition *Definition) []analysis.Diagnostic
		Metadata() rules.RuleMetadata
	}
)

func NewChecker(rules []Rule) Checker {
	return Checker{
		rules: rules,
	}
}

func (c Checker) Check(definition *Definition) []analysis.Diagnostic {
	allDiag := make([]analysis.Diagnostic, 0)
	for _, r := range c.rules {
		allDiag = append(allDiag, r.Apply(definition)...)
	}

	return allDiag
}

func (d *Definition) AddConstructor(constructor *ast.FuncDecl) {
	d.Constructors = append(d.Constructors, constructor)
}

func (d *Definition) AddMethod(method *ast.FuncDecl) {
	d.Methods = append(d.Methods, method)
}
