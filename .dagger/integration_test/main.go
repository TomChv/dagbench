package main

import (
	"dagger/dagbench-test/dagtest"
	"fmt"
)

type DagbenchTest struct{}

func (d *DagbenchTest) All() error {
	// TODO: allow more granularity
	testRunner := dagtest.New().
		WithTest("Basic", TestBasic).
		WithTest("Advanced", TestAdvanced)

	if code := testRunner.Run(); code != 0 {
		return fmt.Errorf("dagbench tests failed with code %d", code)
	}

	return nil
}
