package config

import "fmt"

type Format string

const (
	JSON Format = "json"
	YAML Format = "yaml"
)

func (f Format) String() string {
	return string(f)
}

func StringToFormat(s string) (Format, error) {
	switch s {
	case "json":
		return JSON, nil
	case "yaml":
		return YAML, nil
	default:
		return "", fmt.Errorf("unknown format: %s", s)
	}
}
