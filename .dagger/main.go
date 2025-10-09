package main

import (
	"context"

	"dagger/dagbench/internal/dagger"
)

type Dagbench struct{}

func (d *Dagbench) Bin(
	ctx context.Context,

	//+optional
	platform dagger.Platform,
) *dagger.File {
	opts := dagger.DagbenchCiBuildOpts{}

	if platform != "" {
		opts.Platform = platform
	}

	return dag.DagbenchCi().Build(opts)
}

func (d *Dagbench) CLI(
	ctx context.Context,

	//+default="main"
	daggerVersion string,

	//+optional
	platform dagger.Platform,
) *CLI {
	binary := d.Bin(ctx, platform)

	return newCLI(binary, daggerVersion)
}
