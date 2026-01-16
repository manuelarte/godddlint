package model

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"

	"github.com/manuelarte/godddlint/internal/astutils"
	"github.com/manuelarte/godddlint/rules"
)

type (
	Rule interface {
		Apply(definition *Definition) []analysis.Diagnostic
		Metadata() rules.RuleMetadata
	}

	RuleEnablerChecker interface {
		IsEnabled(node ast.Node) bool
	}
)

var _ RuleEnablerChecker = new(CommentRuleEnabler)

type (
	CommentRuleEnabler struct {
		RuleCode string
	}
)

func (rec CommentRuleEnabler) IsEnabled(n ast.Node) bool {
	if doc, ok := n.(*ast.CommentGroup); ok {
		return !astutils.CommentHasPrefix(doc, fmt.Sprintf("//godddlint:disable:%s", rec.RuleCode))
	}

	return true
}
