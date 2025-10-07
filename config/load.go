package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func NewFromFile(path string) (*Config, error) {
	// Read the file content
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Get the format from the file extension
	ext := filepath.Ext(path)
	format := Format(strings.TrimPrefix(ext, "."))

	// Unmarshal the content to the Config struct
	var c Config
	switch format {
	case JSON:
		if err := json.Unmarshal(content, &c); err != nil {
			return nil, fmt.Errorf("failed to unmarshal config: %w", err)
		}
	case YAML:
		if err := yaml.Unmarshal(content, &c); err != nil {
			return nil, fmt.Errorf("failed to unmarshal config: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported file extension: %s", ext)
	}

	return &c, nil
}
