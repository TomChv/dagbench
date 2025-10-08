package config

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

type Command struct {
	// SpanNames is the list of spans to record.
	SpanNames []string `json:"spanNames,omitempty" yaml:"spanNames,omitempty"`

	// Args is the list of arguments to pass to the dagger binary.
	Args []string `json:"args" yaml:"args"`
}

type Init struct {
	Name        string `json:"name,omitempty"        yaml:"name,omitempty"`
	SDK         string `json:"sdk"                   yaml:"sdk"`
	TemplateDir string `json:"templatedir,omitempty" yaml:"templatedir,omitempty"`
}

func (c *Init) Verify() error {
	if c.SDK == "" {
		return ErrMissingSDKInInit
	}

	return nil
}

// Config is a type representing the configuration of a benchmark.
type Config struct {
	// Name of the benchmark
	Name string `json:"name" yaml:"name"`

	// Iteration is the number of time the commands will be run.
	Iteration int `json:"iteration" yaml:"iteration"`

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
	opts ...OptFunc,
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
//
//nolint:cyclop
func (c *Config) Verify(ctx context.Context) error {
	if c.Iteration < 1 {
		return ErrNonPositiveIteration
	}

	// Verify that the binary exists and save that abs path
	absBinPath, err := exec.LookPath(c.BinPath)
	if err != nil {
		return fmt.Errorf("failed to find dagger binary %s: %w", c.BinPath, err)
	}
	c.BinPath = absBinPath

	// Save the binary version
	c.Version, err = getDaggerVersion(ctx, c.BinPath)
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
		return ErrInitAndModuleOverlap
	}

	// If auto-init is enabled but no workdir provided,
	if c.Workdir == "" && c.DoAutoInit() {
		c.Workdir = generateDefaultWorkdir(c.Name)

		if err := os.MkdirAll(c.Workdir, 0o750); err != nil {
			return fmt.Errorf("failed to create default workdir: %w", err)
		}
	}

	// Enable debug logs
	if c.debug {
		_ = os.MkdirAll(c.DebugDir(), 0o750)
	}

	return nil
}

// DoAutoInit returns true if the config has an init command.
func (c *Config) DoAutoInit() bool {
	return c.Init != nil
}

// CleanUpAfterIteration returns true if a workdir is set but no auto-init
// & no module.
//
// That means the workdir should be cleaned up after each iteration because
// we are likely running an init command inside it.
// If not, then either --module or --auto-init should be set.
func (c *Config) CleanUpAfterIteration() bool {
	return c.Workdir != "" && !c.DoAutoInit() && c.Module == ""
}

// IsDebug returns true if debug is enabled.
func (c *Config) IsDebug() bool {
	return c.debug
}

// EnableDebug enables debug mode.
func (c *Config) EnableDebug() {
	c.debug = true
}

// DebugDir returns the path to the debug directory.
func (c *Config) DebugDir() string {
	return filepath.Join("/tmp", "dagbench-report", c.Name)
}

// Metadatas returns the metadatas of the config.
func (c *Config) Metadatas() map[string]string {
	return map[string]string{
		"name":       c.Name,
		"version":    c.Version,
		"date":       time.Now().Format("2006-01-02"),
		"cloud":      strconv.FormatBool(c.Cloud),
		"iterations": strconv.Itoa(c.Iteration),
	}
}
