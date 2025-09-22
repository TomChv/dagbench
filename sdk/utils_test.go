package sdk

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func ConvertFunctionNameToTraceMarker(t *testing.T) {
	for _, tt := range []struct {
		name     string
		function string
		expected string
	}{
		{
			name:     "simple",
			function: "function",
			expected: "function",
		},
		{
			name:     "container-echo",
			function: "container-echo",
			expected: "containerEcho",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := convertFunctionNameToTraceMarker(tt.function)

			require.Equal(t, tt.expected, got)
		})
	}
}
