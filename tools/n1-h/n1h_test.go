package n1h_test

import (
	"testing"

	"github.com/gostaticanalysis/testutil"
	n1h "github.com/reactive-futakotamagawa/isu-isu-h/tools/n1-h"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	testdata := testutil.WithModules(t, analysistest.TestData(), nil)
	analysistest.Run(t, testdata, n1h.Analyzer, "a")
}
