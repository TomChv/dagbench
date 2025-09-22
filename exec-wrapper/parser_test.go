package execwrapper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtractTimeFromTraceLine(t *testing.T) {
	for _, tt := range []struct {
		name     string
		line     string
		expected string
	}{
		{
			name:     "simple",
			line:     "[1.234s] message",
			expected: "1.234s",
		},
		{
			name:     "no match",
			line:     "message",
			expected: "?",
		},
		{
			name:     "real case",
			line:     "ModuleSource.generatedContextDirectory DONE [9.5s]",
			expected: "9.5s",
		},
		{
			name:     "with logged line",
			line:     "17  : moduleSource DONE [0.2s]",
			expected: "0.2s",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := extractTimeFromTraceLine(tt.line)

			require.Equal(t, tt.expected, got)
		})
	}
}
