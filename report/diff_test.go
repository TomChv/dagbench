package report

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestReportDiff(t *testing.T) {
	for _, tt := range []struct {
		name     string
		report1  *Report
		report2  *Report
		expected *Report
	}{
		{
			name:     "simple",
			report1:  New("report1").AddValue("foo", 100*time.Millisecond).AddValue("bar", 300*time.Millisecond),
			report2:  New("report2").AddValue("foo", 200*time.Millisecond).AddValue("bar", 200*time.Millisecond),
			expected: New("diff-report1-on-report2").AddValue("foo", -100*time.Millisecond).AddValue("bar", 100*time.Millisecond),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.report1.Diff(tt.report2)
			require.Equal(t, tt.expected.name, got.name)
			require.Equal(t, tt.expected.values, got.values)
		})
	}
}
