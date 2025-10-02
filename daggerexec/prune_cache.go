package daggerexec

import (
	"quartz/dagbench.io/config"
)

func pruneCache(c *config.Config) error {
	return execDagger(
		c,
		[]string{"core", "engine", "local-cache", "prune"},
	)
}
