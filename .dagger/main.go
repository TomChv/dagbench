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

	// The dagger tag to create the engine
	//+default="latest"
	//+optional
	daggerTag string,

	// The dagger cache volume to use
	//+optional
	daggerCacheVolume string,

	// The dagger engine container to use with the dagbench CLI.
	//+optional
	daggerCtr *dagger.Container,

	// The platform to use to build the dagbench binary.
	//+optional
	platform dagger.Platform,
) (*CLI, error) {
	binary := d.Bin(ctx, platform)

	if daggerCtr != nil {
		return newCLI(binary, daggerCtr)
	}

	if daggerCacheVolume == "" {
		daggerCacheVolume = buildUniqueCacheVolumeName(daggerTag)
	}

	daggerCtr, err := daggerContainerFromTag(ctx, daggerCacheVolume, daggerTag)
	if err != nil {
		return nil, err
	}

	return newCLI(binary, daggerCtr)
}
