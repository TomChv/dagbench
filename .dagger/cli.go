package main

import (
	"errors"

	"dagger/dagbench/internal/dagger"
)

const (
	daggerRepo = "https://github.com/dagger/dagger"
)

var (
	errDaggerSourceAndVersionConflict = errors.New("dagger source and version are both set but they are mutually exclusive")
	errMissingDaggerSource            = errors.New("dagger source is missing")
)

type CLI struct {
	//+private
	Ctr *dagger.Container
}

func newCLI(dagBenchBinary *dagger.File, opts ...CLIOptsFunc) (*CLI, error) {
	cli := &CLIOpts{}

	for _, opt := range opts {
		opt(cli)
	}

	if cli.daggerVersion != "" && cli.daggerSource != nil {
		return nil, errDaggerSourceAndVersionConflict
	}

	if cli.daggerVersion != "" && cli.daggerSource == nil {
		cli.daggerSource = fetchSourceFromGit(cli.daggerVersion)
	}

	if cli.daggerSource == nil {
		return nil, errMissingDaggerSource
	}

	daggerCtr := dag.
		DaggerDev(dagger.DaggerDevOpts{
			Source: cli.daggerSource,
		}).Dev()

	return &CLI{
		Ctr: daggerCtr.
			WithMountedFile("/bin/dagbench", dagBenchBinary).
			WithEntrypoint([]string{"/bin/dagbench"}),
	}, nil
}

func (c *CLI) Container() *dagger.Container {
	return c.Ctr
}
