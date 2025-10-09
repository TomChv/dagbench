package main

import (
	"context"

	"dagger/dagbench/internal/dagger"
)

type Dagbench struct{}

// Get the dagbench binary file.
func (d *Dagbench) Bin(
	ctx context.Context,

	// The platform to use to build the dagbench binary.
	//+optional
	platform dagger.Platform,
) *dagger.File {
	opts := dagger.DagbenchCiBuildOpts{}

	if platform != "" {
		opts.Platform = platform
	}

	return dag.DagbenchCi().Build(opts)
}

// Get a CLI object to run dagbench inside a container.
func (d *Dagbench) CLI(
	ctx context.Context,

	// The dagger engine to use with dagbench to run the benchmarks
	//+default="main"
	daggerVersion string,

	// The dagger repository to use to build the dagger container
	// This should only be used to test a local dagger version.
	// This will override the daggerVersion if set.
	//+optional
	daggerSource *dagger.Directory,

	// The platform to use to build the dagbench binary.
	//+optional
	platform dagger.Platform,
) (*CLI, error) {
	binary := d.Bin(ctx, platform)

	opts := []CLIOptsFunc{}

	if daggerSource != nil {
		opts = append(opts, withDaggerSource(daggerSource))
	}

	if daggerVersion != "" && daggerSource == nil {
		opts = append(opts, withDaggerVersion(daggerVersion))
	}

	return newCLI(binary, opts...)
}
