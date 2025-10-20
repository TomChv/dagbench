package main

import (
	"dagger/dagbench/internal/dagger"
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

func (c *CLI) Run() *CmdRun {
	return newRunCmd(c.Ctr)
}

func (c *CLI) Plot() *CmdPlot {
	return newPlotCmd(c.Ctr)
}

func (c *CLI) New() *CmdNew {
	return newNewCmd(c.Ctr)
}