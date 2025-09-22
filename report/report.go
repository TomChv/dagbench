package report

import (
	"fmt"
	"strings"
	"time"
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

func (r *Report) String() string {
	var res strings.Builder

	fmt.Fprintf(&res, "Report %s results\n", r.name)

	for k, v := range r.values {
		fmt.Fprintf(&res, "%s\t: %s\n", k, v)
	}

	return res.String()
}

func (r *Report) Stderr() string {
	return r.stderr
}

func (r *Report) AddValue(key string, value time.Duration) {
	r.values[key] = value
}

func (r *Report) AddStderr(stderr string) {
	r.stderr += stderr
}
