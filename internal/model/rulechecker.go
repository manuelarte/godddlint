package model

import (
	"golang.org/x/tools/go/analysis"
)

type (
	RuleChecker struct {
		rules []Rule
	}
)

func NewRuleChecker(rules []Rule) RuleChecker {
	return RuleChecker{
		rules: rules,
	}
}

func (c RuleChecker) Check(definition *Definition) []analysis.Diagnostic {
	allDiag := make([]analysis.Diagnostic, 0, len(c.rules))

	for _, r := range c.rules {
		allDiag = append(allDiag, r.Apply(definition)...)
	}

	return allDiag
}
