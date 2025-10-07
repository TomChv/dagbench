package hook

type Hook interface {
	Hook(line string) error
}
