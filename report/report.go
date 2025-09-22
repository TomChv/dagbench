package report

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"text/tabwriter"
)

type Report struct {
	name string

	stderr string

	values map[string]time.Duration
}

func New(name string) *Report {
	return &Report{
		name:   name,
		values: make(map[string]time.Duration),
	}
}

func NewFromCSV(path string) (*Report, error) {
	name := filepath.Base(path)

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read report file: %w", err)
	}

	csvReader := csv.NewReader(strings.NewReader(string(content)))
	csvReader.Comma = ';'

	values, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to parse report file: %w", err)
	}

	values = values[1:]

	report := New(name)
	for _, row := range values {
		if len(row) < 2 {
			return nil, fmt.Errorf("invalid row: %v", row)
		}

		key := row[0]
		value, err := time.ParseDuration(row[1])
		if err != nil {
			return nil, fmt.Errorf("failed to parse duration: %w", err)
		}

		report = report.AddValue(key, value)
	}

	return report, nil
}

func (r *Report) String() string {
	var b strings.Builder

	fmt.Fprintf(&b, "Report %s\n", r.name)
	tw := tabwriter.NewWriter(&b, 0, 0, 2, ' ', 0)

	for k, v := range r.values {
		fmt.Fprintf(tw, "%s\t:\t %s\n", k, v)
	}

	_ = tw.Flush()

	return b.String()
}

func (r *Report) Stderr() string {
	return r.stderr
}

func (r *Report) AddValue(key string, value time.Duration) *Report {
	r.values[key] = value

	return r
}

func (r *Report) AddStderr(stderr string) *Report {
	r.stderr += stderr

	return r
}
