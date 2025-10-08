package hook

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtractTimeFromTraceLine(t *testing.T) {
	t.Parallel()

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
			name:     "real case",
			line:     "ModuleSource.generatedContextDirectory DONE [9.5s]",
			expected: "9.5s",
		},
		{
			name:     "with logged line",
			line:     "17  : moduleSource DONE [0.2s]",
			expected: "200ms",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := extractTimeFromTraceLine(tt.line)

			require.NoError(t, err)
			require.Equal(t, tt.expected, got.String())
		})
	}
}
