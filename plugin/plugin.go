// Package plugin add supports to add this analyzer as a golangci-lint
// plugin.
package plugin

import (
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"

	"github.com/manuelarte/godddlint/analyzer"
)

//nolint:gochecknoinits // init needed for plugin
func init() {
	register.Plugin("godddlint", New)
}

func New(_ any) (register.LinterPlugin, error) {
	return &godddlintPlugin{}, nil
}

var _ register.LinterPlugin = new(godddlintPlugin)

type godddlintPlugin struct{}

func (u godddlintPlugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		analyzer.New(),
	}, nil
}

func (u godddlintPlugin) GetLoadMode() string {
	return register.LoadModeSyntax
}
