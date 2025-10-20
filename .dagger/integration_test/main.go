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

	eg.Go(func() error {
		sctx, span := Tracer().Start(gctx, "basic")
		defer span.End()

		return testBasic(sctx)
	})
	eg.Go(func() error {
		sctx, span := Tracer().Start(gctx, "advanced")
		defer span.End()

		return testAdvanced(sctx)
	})

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
