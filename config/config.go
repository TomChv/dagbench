package config

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type Command struct {
	SpanNames []string `json:"spanNames,omitempty" yaml:"spanNames,omitempty"`
	Args      []string `json:"args" yaml:"args"`
}

type Init struct {
	Name        string `json:"name,omitempty" yaml:"name,omitempty"`
	SDK         string `json:"sdk" yaml:"sdk"`
	TemplateDir string `json:"templatedir,omitempty" yaml:"templatedir,omitempty"`
}

func (c *Init) Verify() error {
	if c.SDK == "" {
		return fmt.Errorf("init.sdk is required")
	}

	return nil
}

type Config struct {
	// Name of the benchmark
	Name      string `json:"name" yaml:"name"`
	Iteration int    `json:"iteration" yaml:"iteration"`

	// Dagger binary path
	BinPath string `json:"binPath" yaml:"binPath"`
	Version string `json:"version" yaml:"version"`

	// Module to use for the benchmark (cannot be set if RunInit is true)
	Module string `json:"module,omitempty" yaml:"module,omitempty"`

	// Working directory for the benchmark (cannot be set if module is set)
	Workdir string `json:"workdir" yaml:"workdir"`

	// Automatically init the module
	Init *Init `json:"init,omitempty" yaml:"init,omitempty"`

	Commands []*Command `json:"commands" yaml:"commands"`

	// Use --cloud in commands
	Cloud bool `json:"cloud" yaml:"cloud"`

	// Enable debug mode
	debug bool
}

func New(
	name string,
	daggerBin string,
	iteration int,
	opts ...ConfigOptFunc,
) *Config {
	c := &Config{
		Name:      name,
		Iteration: iteration,
		BinPath:   daggerBin,
		Cloud:     false,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// Verify if the configuration is valid.
func (c *Config) Verify() error {
	if c.Iteration < 1 {
		return fmt.Errorf("iteration must be greater than 0")
	}

	// Verify that the binary exists and save that abs path
	absBinPath, err := exec.LookPath(c.BinPath)
	if err != nil {
		return fmt.Errorf("failed to find dagger binary %s: %w", c.BinPath, err)
	}
	c.BinPath = absBinPath

	// Save the binary version
	c.Version, err = getDaggerVersion(c.BinPath)
	if err != nil {
		return fmt.Errorf("failed to get dagger version: %w", err)
	}

	if c.Init != nil {
		if err := c.Init.Verify(); err != nil {
			return fmt.Errorf("init is set but a field is invalid: %w", err)
		}
	}

	// If a module is set, we cannot init automatically
	if c.Module != "" && c.DoAutoInit() {
		return fmt.Errorf("cannot set both module and auto-init")
	}

	// If auto-init is enabled but no workdir provided,
	if c.Workdir == "" && c.DoAutoInit() {
		c.Workdir = generateDefaultWorkdir(c.Name)

		if err := os.MkdirAll(c.Workdir, 0755); err != nil {
			return fmt.Errorf("failed to create default workdir: %w", err)
		}
	}

	// Enable debug logs
	if c.debug {
		_ = os.MkdirAll(c.DebugDir(), 0755)
	}

	return nil
}

func (c *Config) DoAutoInit() bool {
	return c.Init != nil
}

// If a workdir is set but no auto-init & no module, that means
// the workdir should be cleaned up after each iteration because
// we are likely running an init command inside it.
// If not, then either --module or --auto-init should be set.
func (c *Config) CleanUpAfterIteration() bool {
	return c.Workdir != "" && !c.DoAutoInit() && c.Module == ""
}

func (c *Config) IsDebug() bool {
	return c.debug
}

func (c *Config) EnableDebug() {
	c.debug = true
}

func (c *Config) DebugDir() string {
	return filepath.Join("/tmp", "dagbench-report", c.Name)
}

func (c *Config) Metadatas() map[string]string {
	return map[string]string{
		"name":       c.Name,
		"version":    c.Version,
		"date":       time.Now().Format("2006-01-02"),
		"cloud":      fmt.Sprintf("%t", c.Cloud),
		"iterations": fmt.Sprintf("%d", c.Iteration),
	}
}
