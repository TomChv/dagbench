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

func newCLI(dagBenchBinary *dagger.File, daggerCtr *dagger.Container) (*CLI, error) {
	return &CLI{
		Ctr: daggerCtr.
			WithMountedFile("/bin/dagbench", dagBenchBinary).
			WithEntrypoint([]string{"/bin/dagbench"}),
	}, nil
}

func (c *CLI) Container() *dagger.Container {
	return c.Ctr
}
