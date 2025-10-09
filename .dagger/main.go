package main

import (
	"context"

	"dagger/dagbench/integration"
	"dagger/dagbench/internal/dagger"

	"golang.org/x/sync/errgroup"
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

func (d *Dagbench) Bin(
	ctx context.Context,

	//+optional
	platform dagger.Platform,
) (_ *dagger.File, err error) {
	if platform == "" {
		platform, err = dag.DefaultPlatform(ctx)
		if err != nil {
			return nil, err
		}
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

	//+optional
	platform dagger.Platform,
) (*CLI, error) {
	binary, err := d.Bin(ctx, platform)
	if err != nil {
		return nil, err
	}

	return newCLI(binary, daggerVersion)
}

func (d *Dagbench) IntegrationTest(
	ctx context.Context,
) error {
	cli, err := d.CLI(ctx, "main", "")
	if err != nil {
		return err
	}

	eg, gctx := errgroup.WithContext(ctx)

	eg.Go(func() error { return integration.TestBasic(gctx, cli) })

	return eg.Wait()
}
