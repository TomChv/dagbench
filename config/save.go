package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Save the config to a file and returns the path to that file
func (c *Config) Save(dirpath string, format Format) (string, error) {
	path := filepath.Join(dirpath, fmt.Sprintf("%s.%s", c.Name, format))

	// Create the directory if it doesn't exist and open file
	if err := os.MkdirAll(filepath.Dir(dirpath), 0o750); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer func() { _ = file.Close() }()

	// Marshal the config to the appropriate format
	var content []byte
	switch format {
	case JSON:
		content, err = json.MarshalIndent(c, "", "  ")
		if err != nil {
			return "", fmt.Errorf("failed to marshal config to json: %w", err)
		}
	case YAML:
		content, err = yaml.Marshal(c)
		if err != nil {
			return "", fmt.Errorf("failed to marshal config to yaml: %w", err)
		}
	default:
		return "", UnsupportedFormatError(format)
	}

	// Write the marshaled config to the file
	if _, err := file.Write(content); err != nil {
		return "", fmt.Errorf("failed to write config to file %s: %w", path, err)
	}

	return path, nil
}
