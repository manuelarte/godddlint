package analyzer

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		desc     string
		patterns string
		options  map[string]string
	}{
		{
			desc:     "valueObjects",
			patterns: "valueObjects",
		},
		{
			desc:     "valueObjects disable rules",
			patterns: "valueObjects-disable-rules",
		},
		{
			desc:     "entities",
			patterns: "entities",
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

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
