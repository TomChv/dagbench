package main

import (
	"context"

	"dagger/dagbench-ci/internal/dagger"
)

type DagbenchCi struct{}

func (d *DagbenchCi) Build(
	ctx context.Context,

	//+ignore=["**", "!**/*.go", "!go.mod", "!go.sum", ".dagger/"]
	//+defaultPath="/"
	source *dagger.Directory,

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
		Go(source).
		Build(dagger.GoBuildOpts{
			Platform: platform,
		}).File("bin/dagbench.io"), nil
}

func (d *DagbenchCi) Lint(
	ctx context.Context,

	//+ignore=["**", "!**/*.go", "!go.mod", "!go.sum", "!.golangci.yml", ".dagger/"]
	//+defaultPath="/"
	source *dagger.Directory,
) (string, error) {
	return dag.
		Container().
		From("golangci/golangci-lint:v2.5-alpine").
		WithDirectory("/app", source).
		WithWorkdir("/app").
		WithExec([]string{"golangci-lint", "run"}).
		Stdout(ctx)
}

func (d *DagbenchCi) Fmt(
	ctx context.Context,

	//+ignore=["**", "!**/*.go", "!go.mod", "!go.sum", "!.golangci.yml", ".dagger/"]
	//+defaultPath="/"
	source *dagger.Directory,
) *dagger.Directory {
	return dag.
		Container().
		From("golangci/golangci-lint:v2.5-alpine").
		WithDirectory("/app", source).
		WithWorkdir("/app").
		WithExec([]string{"golangci-lint", "fmt"}).
		Directory("/app")
}
