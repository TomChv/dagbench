package main

import "dagger/dagbench/internal/dagger"

type CmdPlot struct {
	//+private
	Ctr *dagger.Container

	//+private
	Command []string

	//+private
	Files []string
}

func newPlotCmd(ctr *dagger.Container) *CmdPlot {
	return &CmdPlot{
		Ctr: ctr,
		Command: []string{
			"plot",
		},
	}
}

func (c *CmdPlot) VersionDiff() *dagger.File {
	outputFilename := "out.svg"
	commands := append(c.Command, "version-diff", "-o", outputFilename)
	commands = append(commands, c.Files...)

	return c.Ctr.
		WithExec(commands, dagger.ContainerWithExecOpts{
			UseEntrypoint: true,
		}).
		File(outputFilename)
}

func (c *CmdPlot) MultiBar() *dagger.File {
	outputFilename := "out.svg"
	commands := append(c.Command, "multi-bar", "-o", outputFilename)
	commands = append(commands, c.Files...)

	return c.Ctr.
		WithExec(commands, dagger.ContainerWithExecOpts{
			UseEntrypoint: true,
		}).
		File(outputFilename)
}

func (c *CmdPlot) WithName(name string) *CmdPlot {
	c.Command = append(c.Command, "--name", name)

	return c
}

func (c *CmdPlot) WithFile(filename string, file *dagger.File) *CmdPlot {
	c.Files = append(c.Files, filename)
	c.Ctr = c.Ctr.WithFile(filename, file)

	return c
}
