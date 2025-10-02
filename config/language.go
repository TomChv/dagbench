package config

import "fmt"

type Language string

const (
	Go         Language = "go"
	TypeScript Language = "typescript"
)

func (l Language) String() string {
	return string(l)
}

func StringToLanguage(s string) (Language, error) {
	switch s {
	case "go":
		return Go, nil
	case "typescript":
		return TypeScript, nil
	default:
		return "", fmt.Errorf("unknown language: %s", s)
	}
}
