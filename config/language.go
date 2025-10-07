package config

type Language string

const (
	Go         Language = "go"
	TypeScript Language = "typescript"
)

func (l Language) String() string {
	return string(l)
}
