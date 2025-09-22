package sdk

import (
	"fmt"
	"os/exec"
	"quartz/dagbenchmark.io/config"
	execwrapper "quartz/dagbenchmark.io/exec-wrapper"
)

type SDK interface {
	Init() *execwrapper.ExecWrapper

	Develop() *execwrapper.ExecWrapper

	Functions() *execwrapper.ExecWrapper

	Call(callArgs []string) *execwrapper.ExecWrapper

	PruneCache() error
}

func NewSDK(conf *config.Config) (SDK, error) {
	switch conf.Language {
	case config.Go:
		return newGo(conf), nil
	default:
		return nil, nil
	}
}

type sdk struct {
	config *config.Config
}

func (s *sdk) PruneCache() error {
	fmt.Println("Pruning cache...")

	return exec.Command(s.config.BinPath, "core", "engine", "local-cache", "prune").Run()
}
