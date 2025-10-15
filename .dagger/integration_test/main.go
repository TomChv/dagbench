package main

import (
	"context"

	"golang.org/x/sync/errgroup"
)

type DagbenchTest struct{}

func (d *DagbenchTest) All(
	ctx context.Context,
) error {
	eg, gctx := errgroup.WithContext(ctx)

	eg.Go(func() error { return testBasic(gctx) })
	eg.Go(func() error { return testAdvanced(gctx) })

	return eg.Wait()
}

func (d *DagbenchTest) Basic(
	ctx context.Context,
) error {
	return testBasic(ctx)
}

func (d *DagbenchTest) Advanced(
	ctx context.Context,
) error {
	return testAdvanced(ctx)
}
