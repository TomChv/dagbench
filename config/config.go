package config

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type Config struct {
	Name      string   `json:"name"`
	Language  Language `json:"language"`
	BinPath   string   `json:"binPath"`
	TempDir   string   `json:"tempDir"`
	ReportDir string   `json:"reportDir"`

	version string
}

type ConfigOpts func(c *Config)

func WithBinPath(binPath string) ConfigOpts {
	return func(c *Config) {
		c.BinPath = binPath
	}
}

func WithName(name string) ConfigOpts {
	return func(c *Config) {
		c.Name = name
	}
}

func WithReportDir(reportDir string) ConfigOpts {
	return func(c *Config) {
		c.ReportDir = reportDir
	}
}

func NewConfig(language string, opts ...ConfigOpts) (*Config, error) {
	parsedLanguage, err := stringToLanguage(language)
	if err != nil {
		return nil, fmt.Errorf("failed to convert language: %w", err)
	}

	c := &Config{
		Language: parsedLanguage,
	}

	for _, opt := range opts {
		opt(c)
	}

	if c.ReportDir == "" {
		wd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed to get current working directory: %w", err)
		}

		c.ReportDir = filepath.Join(wd, "bench-reports", c.Name)
	}

	if c.BinPath == "" {
		daggerBinPath, err := exec.LookPath("dagger")
		if err != nil {
			return nil, fmt.Errorf("failed to find dagger binary: %w", err)
		}

		c.BinPath = daggerBinPath
	}

	version, err := exec.Command(c.BinPath, "version").Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get dagger version: %w", err)
	}
	c.version = strings.TrimSpace(string(version))

	c.TempDir = filepath.Join(os.TempDir(), "dagger-benchmark", uuid.NewString())
	if err := os.MkdirAll(c.TempDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", err)
	}

	return c, nil
}

func NewConfigFromFile(path string) (*Config, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	fmt.Println("Loading config from", absPath)
	content, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var c Config
	if err := json.Unmarshal(content, &c); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	version, err := exec.Command(c.BinPath, "version").Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get dagger version: %w", err)
	}
	c.version = strings.TrimSpace(string(version))

	return &c, nil
}

func (c *Config) Save(dirpath string) error {
	fileName := fmt.Sprintf("%s-%s.json", c.Language, c.Version())
	if c.Name != "" {
		fileName = fmt.Sprintf("%s.json", c.Name)
	}

	absPath, err := filepath.Abs(dirpath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}
	path := filepath.Join(absPath, fileName)

	fmt.Println("Saving config at", path)

	if err := os.MkdirAll(absPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	content, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if _, err := file.Write(content); err != nil {
		return fmt.Errorf("failed to write config to file: %w", err)
	}

	return nil
}

func (c *Config) String() string {
	return fmt.Sprintf(
		"Name\t\t: %s\nDagger version\t: %s\nDagger path\t: %s\nLanguage\t: %s\nTemp dir\t: %s\nReport dir\t: %s",
		c.Name, c.version, c.BinPath, c.Language, c.TempDir, c.ReportDir,
	)
}

func (c *Config) Version() string {
	return strings.Split(c.version, " ")[1]
}
