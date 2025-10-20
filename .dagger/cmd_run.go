package main

import (
	"strconv"
	"strings"

	"dagger/dagbench/internal/dagger"
)

type CmdRun struct {
	//+private
	Ctr *dagger.Container

	//+private
	Command []string
}

func newRunCmd(ctr *dagger.Container) *CmdRun {
	return &CmdRun{
		Ctr: ctr,
		Command: []string{
			"run",
		},
	}
}

func (c *CmdRun) Exec() *dagger.File {
	outputName := "out.txt"
	commands := append(c.Command, "-o", outputName)

	return c.Ctr.
		WithExec(commands, dagger.ContainerWithExecOpts{
			UseEntrypoint: true,
		}).
		File(outputName)
}

func (c *CmdRun) WithName(name string) *CmdRun {
	c.Command = append(c.Command, "--name", name)

	return c
}

func (c *CmdRun) WithAutoInit() *CmdRun {
	c.Command = append(c.Command, "--auto-init")

	return c
}

func (c *CmdRun) WithSdk(sdk string) *CmdRun {
	c.Command = append(c.Command, "--sdk", sdk)

	return c
}

func (c *CmdRun) WithWorkdir(
	workdir string,

	//+optional
	directory *dagger.Directory,
) *CmdRun {
	if directory == nil {
		directory = dag.Directory()
	}

	c.Command = append(c.Command, "--workdir", workdir)
	c.Ctr = c.Ctr.WithDirectory(workdir, directory)

	return c
}

func (c *CmdRun) WithIteration(iteration int) *CmdRun {
	c.Command = append(c.Command, "--iteration", strconv.Itoa(iteration))

	return c
}

func (c *CmdRun) WithCommand(command []string) *CmdRun {
	c.Command = append(c.Command, "--command", strings.Join(command, " "))

	return c
}

func (c *CmdRun) WithSpan(span string) *CmdRun {
	c.Command = append(c.Command, "--span", span)

	return c
}

func (c *CmdRun) WithTemplateDir(
	directory *dagger.Directory,

	//+default="/template"
	templatePath string,
) *CmdRun {
	c.Command = append(c.Command, "--template-dir", templatePath)
	c.Ctr = c.Ctr.WithDirectory(templatePath, directory)

	return c
}

func (c *CmdRun) WithModule(
	modulePath string,

	moduleDir *dagger.Directory,
) *CmdRun {
	c.Command = append(c.Command, "--module", modulePath)
	c.Ctr = c.Ctr.WithDirectory(modulePath, moduleDir)

	return c
}

func (c *CmdRun) WithConfigFile(
	configFile *dagger.File,

	//+default="/config.json"
	configPath string,
) *CmdRun {
	c.Command = append(c.Command, "--config", configPath)
	c.Ctr = c.Ctr.WithFile(configPath, configFile)

	return c
}
