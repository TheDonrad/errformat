package linters

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/golangci/plugin-module-register/register"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestErrFormatLinter(t *testing.T) {
	tests := []struct {
		name     string
		testPath string
	}{
		{
			name:     "exported errors",
			testPath: "testlintdata/exported_errors",
		},
		{
			name:     "unexported errors",
			testPath: "testlintdata/unexported_errors",
		},
		{
			name:     "mixed cases",
			testPath: "testlintdata/mixed_cases",
		},
		{
			name:     "edge cases",
			testPath: "testlintdata/edge_cases",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newPlugin, err := register.GetPlugin("errformat")
			require.NoError(t, err)

			plugin, err := newPlugin(nil)
			require.NoError(t, err)

			analyzers, err := plugin.BuildAnalyzers()
			require.NoError(t, err)

			analysistest.Run(t, testdataDir(t), analyzers[0], tt.testPath)
		})
	}
}

func testdataDir(t *testing.T) string {
	t.Helper()

	_, testFilename, _, ok := runtime.Caller(1)
	if !ok {
		require.Fail(t, "unable to get current test filename")
	}

	return filepath.Join(filepath.Dir(testFilename), "testdata")
}
