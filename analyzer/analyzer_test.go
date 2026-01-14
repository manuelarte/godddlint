package analyzer

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testCases := []struct {
		desc     string
		patterns string
		options  map[string]string
	}{
		{
			desc:     "valueObjects",
			patterns: "valueObjects",
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			a := New()

			for k, v := range test.options {
				err := a.Flags.Set(k, v)
				if err != nil {
					t.Fatal(err)
				}
			}

			analysistest.Run(t, analysistest.TestData(), a, test.patterns)
		})
	}
}
