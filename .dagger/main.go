package main

import (
	"context"

	"dagger/dagbench/internal/dagger"
)

type Dagbench struct {
	//+private
	Source *dagger.Directory
}

func New(
	//+defaultPath="/"
	source *dagger.Directory,
) *Dagbench {
	return &Dagbench{
		Source: source,
	}
}

func (d *Dagbench) Bin(ctx context.Context) (*dagger.File, error) {
	platform, err := dag.DefaultPlatform(ctx)
	if err != nil {
		return nil, err
	}

	return dag.
		DagbenchCi(d.Source).
		Build(dagger.GoBuildOpts{
			Platform: platform,
		}).File("bin/dagbench.io"), nil
}

func (d *Dagbench) CLI(
	ctx context.Context,

	//+default="main"
	daggerVersion string,
) (*CLI, error) {
	binary, err := d.Bin(ctx)
	if err != nil {
		return nil, err
	}

	return newCLI(binary, daggerVersion)
}
