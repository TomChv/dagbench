package main

import (
	"context"
	"fmt"

	"dagger/dagbench-test/internal/dagger"

	"golang.org/x/sync/errgroup"
)

type DagbenchTest struct{}

func (d *DagbenchTest) All(
	ctx context.Context,
) error {
	ctr, err := getDagBenchContainer(ctx)
	if err != nil {
		return err
	}

	eg, gctx := errgroup.WithContext(ctx)

	eg.Go(func() error { return testBasic(gctx, ctr) })

	return eg.Wait()
}

func (d *DagbenchTest) Basic(
	ctx context.Context,
) error {
	ctr, err := getDagBenchContainer(ctx)
	if err != nil {
		return err
	}

	return testBasic(ctx, ctr)
}

func getDagBenchContainer(ctx context.Context) (*dagger.Container, error) {
	cli := dag.Dagbench().Cli()

	ctr, err := cli.Container().Sync(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get dagbench container: %w", err)
	}

	return ctr, nil
}
