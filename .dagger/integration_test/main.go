package main

import (
	"context"

	"golang.org/x/sync/errgroup"
)

type DagbenchTest struct{}

func (d *DagbenchTest) All(
	ctx context.Context,
) error {
	cli := dag.Dagbench().Cli()

	eg, gctx := errgroup.WithContext(ctx)

	eg.Go(func() error { return testBasic(gctx, cli) })

	return eg.Wait()
}

func (d *DagbenchTest) Basic(
	ctx context.Context,
) error {
	cli := dag.Dagbench().Cli()

	return testBasic(ctx, cli)
}
