package sdk

import (
	"os/exec"
	"quartz/dagbenchmark.io/config"

	execwrapper "quartz/dagbenchmark.io/exec-wrapper"
)

type Go struct {
	sdk
}

func newGo(config *config.Config) *Go {
	return &Go{
		sdk: sdk{config: config},
	}
}

func (g *Go) Init() *execwrapper.ExecWrapper {
	cmd := exec.Command(g.config.BinPath, "-vv", "init", "--sdk=go", "--source=.", "--name=test")
	cmd.Dir = g.config.TempDir

	return execwrapper.NewExecWrapper(cmd, []string{
		"run codegen",
		"moduleSource",
		"generatedContextDirectory",
	})
}
