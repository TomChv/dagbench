package dagtest

import (
	"testing"

	"github.com/dagger/testctx/oteltest"
)

type Test func(t *testing.T)

type DagTest struct {
	tests map[string]Test
}

func New() *DagTest {
	return &DagTest{
		tests: make(map[string]Test),
	}
}

func (d *DagTest) WithTest(name string, test Test) *DagTest{
	d.tests[name] = test

	return d
}

func (d *DagTest) Run() int {
	suites := []testing.InternalTest{}

	for name, test := range d.tests {
		suites = append(suites, testing.InternalTest{
			Name: name,
			F:    test,
		})
	}

	testDep := matchStringOnly(func(pat, str string) (bool, error) {
		return true, nil
	})

	m := testing.MainStart(testDep, suites, nil, nil, nil)

	return oteltest.Main(m)
}
