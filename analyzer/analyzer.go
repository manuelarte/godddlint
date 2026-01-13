package analyzer

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

func New() *analysis.Analyzer {
	g := godddlint{}

	return &analysis.Analyzer{
		Name:     "godddlint",
		Doc:      "checks domain structs honor best practices",
		URL:      "https://github.com/manuelarte/godddlint",
		Run:      g.run,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
}

type godddlint struct{}

func (g godddlint) run(pass *analysis.Pass) (any, error) {
	//nolint:nilnil // nothing to do here
	return nil, nil
}
