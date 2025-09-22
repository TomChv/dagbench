package report

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestReportMerge(t *testing.T) {
	for _, tt := range []struct {
		name     string
		reports  []*Report
		expected *Report
	}{
		{
			name: "simple",
			reports: []*Report{
				New("report1").AddValue("foo", 100*time.Millisecond).AddValue("bar", 200*time.Millisecond),
				New("report2").AddValue("foo", 200*time.Millisecond).AddValue("bar", 300*time.Millisecond),
			},
			expected: New("avg-report1-on-2-runs").AddValue("foo", 150*time.Millisecond).AddValue("bar", 250*time.Millisecond),
		},
		{
			name: "10 reports",
			reports: []*Report{
				New("report1").AddValue("foo", 100*time.Millisecond).AddValue("bar", 100*time.Millisecond),
				New("report2").AddValue("foo", 100*time.Millisecond).AddValue("bar", 100*time.Millisecond),
				New("report3").AddValue("foo", 100*time.Millisecond).AddValue("bar", 100*time.Millisecond),
				New("report4").AddValue("foo", 100*time.Millisecond).AddValue("bar", 100*time.Millisecond),
				New("report5").AddValue("foo", 100*time.Millisecond).AddValue("bar", 100*time.Millisecond),
				New("report6").AddValue("foo", 100*time.Millisecond).AddValue("bar", 100*time.Millisecond),
				New("report7").AddValue("foo", 100*time.Millisecond).AddValue("bar", 100*time.Millisecond),
				New("report8").AddValue("foo", 100*time.Millisecond).AddValue("bar", 100*time.Millisecond),
				New("report9").AddValue("foo", 100*time.Millisecond).AddValue("bar", 100*time.Millisecond),
				New("report10").AddValue("foo", 100*time.Millisecond).AddValue("bar", 100*time.Millisecond),
			},
			expected: New("avg-report1-on-10-runs").AddValue("foo", 100*time.Millisecond).AddValue("bar", 100*time.Millisecond),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Merge(tt.reports...)
			require.NoError(t, err)

			require.Equal(t, tt.expected.name, got.name)
			require.Equal(t, tt.expected.values, got.values)
		})
	}
}
