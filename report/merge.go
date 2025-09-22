package report

import (
	"fmt"
	"time"
)

// Merge merges multiple reports into a single report
// Only values are merged, stderr is ignored
// TODO: add tests
func Merge(reports ...*Report) (*Report, error) {
	if len(reports) == 0 {
		return nil, fmt.Errorf("no reports provided")
	}

	name := reports[0].name
	mergedReportName := fmt.Sprintf("avg-%s-on-%d-runs", name, len(reports))

	mergedReport := New(mergedReportName)
	mergedValuesList := make(map[string][]time.Duration)

	// TODO: Do it in parallel
	for _, report := range reports {
		for k, v := range report.values {
			mergedValuesList[k] = append(mergedValuesList[k], v)
		}
	}

	// TODO: Do it in parallel
	for k, v := range mergedValuesList {
		mergedReport.AddValue(k, computeAverage(v))
	}

	return mergedReport, nil
}

func computeAverage(values []time.Duration) time.Duration {
	n := int64(len(values))

	var sum int64
	for _, duration := range values {
		sum += duration.Nanoseconds()
	}

	return time.Duration((sum + n/2) / n)
}
