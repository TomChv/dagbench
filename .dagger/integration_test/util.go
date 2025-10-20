package main

import (
	"fmt"
	"testing"

	"dagger/dagbench-test/internal/dagger"

	"github.com/dagger/testctx"
	"github.com/dagger/testctx/oteltest"
	"github.com/google/uuid"
)

func Middlewares() []testctx.Middleware[*testing.T] {
	return []testctx.Middleware[*testing.T]{
		testctx.WithParallel(),             // run tests in parallel
		oteltest.WithTracing[*testing.T](), // trace each test and subtest
		//oteltest.WithLogging[*testing.T](), // direct t.Log etc. to span logs
	}
}

func getTestCLI(testName string) *dagger.DagbenchCli {
	return dag.
		Dagbench().
		Cli(dagger.DagbenchCliOpts{
			DaggerTag:         "v0.19.2",
			DaggerCacheVolume: fmt.Sprintf("dagbench-test-%s-%s", testName, uuid.NewString()),
		})
}
