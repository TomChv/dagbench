package main

import "dagger/dagbench/internal/dagger"

type CLIOpts struct {
	daggerVersion string
	daggerSource  *dagger.Directory
}

type CLIOptsFunc func(c *CLIOpts)

func withDaggerVersion(version string) CLIOptsFunc {
	return func(c *CLIOpts) {
		c.daggerVersion = version
	}
}

func withDaggerSource(source *dagger.Directory) CLIOptsFunc {
	return func(c *CLIOpts) {
		c.daggerSource = source
	}
}
