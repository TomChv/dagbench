package report

import (
	"maps"
	"os"
	"time"
)

type Report struct {
	name   string
	output string
	values map[string]time.Duration
}

func New(name string) *Report {
	return &Report{
		name:   name,
		values: make(map[string]time.Duration),
	}
}

func (r *Report) WithValues(values map[string]time.Duration) *Report {
	maps.Copy(r.values, values)

	return r
}

func (r *Report) WithValue(key string, value time.Duration) *Report {
	r.values[key] = value

	return r
}

func (r *Report) WithOutput(output string) *Report {
	r.output = output

	return r
}

func (r *Report) HasOutput() bool {
	return r.output != ""
}

func (r *Report) SaveOutput(path string) error {
	return os.WriteFile(path, []byte(r.output), 0644)
}
