package main

import (
	"dagger/dagbench/internal/dagger"
	"fmt"
)

type CmdNew struct {
	//+private
	Ctr *dagger.Container

	//+private
	Command []string

	//+private
	Name string
}

func newNewCmd(ctr *dagger.Container) *CmdNew {
	return &CmdNew{
		Ctr: ctr,
		Command: []string{
			"new",
		},
	}
}

func (c *CmdNew) Exec() *dagger.File {
	outputName := fmt.Sprintf("%s.json", c.Name)

	return c.Ctr.
		WithExec(c.Command, dagger.ContainerWithExecOpts{
			UseEntrypoint: true,
		}).
		File(outputName)
}

func (c *CmdNew) WithName(name string) *CmdNew {
	c.Command = append(c.Command, "--name", name)
	c.Name = name

	return c
}

func (c *CmdNew) WithSDK(sdk string) *CmdNew {
	c.Command = append(c.Command, "--sdk", sdk)

	return c
}

func (c *CmdNew) WithRecipe(recipe string) *CmdNew {
	c.Command = append(c.Command, "--recipe", recipe)

	return c
}

func (c *CmdNew) WithAutoInit() *CmdNew {
	c.Command = append(c.Command, "--auto-init")

	return c
}
