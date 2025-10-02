package config

import (
	"fmt"
	"os/exec"
)

type Config struct {
	Name        string              `json:"name" yaml:"name"`
	SDK         string              `json:"sdk" yaml:"sdk"`
	BinPath     string              `json:"binPath" yaml:"binPath"`
	Version     string              `json:"version" yaml:"version"`
	Workdir     string              `json:"workdir" yaml:"workdir"`
	TemplateDir string              `json:"templatedir,omitempty" yaml:"templatedir,omitempty"`
	RunInit     bool                `json:"runInit" yaml:"runInit"`
	Commands    map[string][]string `json:"commands" yaml:"commands"`
	Cloud       bool                `json:"cloud" yaml:"cloud"`
}

func New(name string, sdk string, binName string, opts ...ConfigOptFunc) (c *Config, err error) {
	c = &Config{
		Name:     name,
		SDK:      sdk,
		RunInit:  true,
		Commands: make(map[string][]string),
		Cloud:    false,
	}

	for _, opt := range opts {
		opt(c)
	}

	if c.Workdir == "" {
		c.Workdir = generateDefaultWorkdir(c.Name)
	}

	if len(c.Commands) == 0 {
		c.Commands = defaultCommands
	}

	c.BinPath, err = exec.LookPath(binName)
	if err != nil {
		return nil, fmt.Errorf("failed to find dagger binary: %w", err)
	}

	c.Version, err = getDaggerVersion(c.BinPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get dagger version: %w", err)
	}

	return c, nil
}
