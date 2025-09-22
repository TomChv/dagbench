package report

import "fmt"

// Diff computes the difference between two reports
// Only values are merged, stderr is ignored
// If a value is not in the other report, it's inserted as
// it is.
func (r *Report) Diff(other *Report) *Report {
	diffReport := New(fmt.Sprintf("diff-%s-on-%s", r.name, other.name))

	for k, v := range r.values {
		// If the value isn't in the other report, it's a new value
		otherValue, ok := other.values[k]
		if !ok {
			diffReport = diffReport.AddValue(k, v)
		}

		diffReport = diffReport.AddValue(k, v-otherValue)
	}

	return diffReport
}
