package main

import (
	"dagger/dagbench-test/internal/dagger"
	"fmt"

	"github.com/google/uuid"
)

func getTestCtr(testName string) *dagger.Container {
	return dag.
		Dagbench().
		Cli(dagger.DagbenchCliOpts{
			DaggerTag: "v0.19.2",
			DaggerCacheVolume: fmt.Sprintf("dagbench-test-%s-%s", testName, uuid.NewString()),
		}).Container()
}
