package main

import (
	"fmt"

	"dagger/dagbench-test/internal/dagger"

	"github.com/google/uuid"
)

func getTestCLI(testName string) *dagger.DagbenchCli {
	return dag.
		Dagbench().
		Cli(dagger.DagbenchCliOpts{
			DaggerTag:         "v0.19.2",
			DaggerCacheVolume: fmt.Sprintf("dagbench-test-%s-%s", testName, uuid.NewString()),
		})
}
