package sdk

import (
	"fmt"
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
	cmd := exec.Command(g.config.BinPath, "-vv", "--progress=plain", "init", "--sdk=go", "--source=.", "--name=test")
	cmd.Dir = g.config.TempDir

	name := fmt.Sprintf("init-go-%s", g.config.Version())
	if g.config.Name != "" {
		name = fmt.Sprintf("%s-init", g.config.Name)
	}

	return execwrapper.NewExecWrapper(
		name,
		cmd,
		[]string{
			"run codegen",
			"moduleSource",
			"generatedContextDirectory",
		})
}

func (g *Go) Develop() *execwrapper.ExecWrapper {
	cmd := exec.Command(g.config.BinPath, "-vv", "--progress=plain", "develop")
	cmd.Dir = g.config.TempDir

	name := fmt.Sprintf("develop-go-%s", g.config.Version())
	if g.config.Name != "" {
		name = fmt.Sprintf("%s-develop", g.config.Name)
	}

	return execwrapper.NewExecWrapper(
		name,
		cmd,
		[]string{
			"develop",
			"run codegen",
			"moduleSource",
			"generatedContextDirectory",
		})
}

func (g *Go) Functions() *execwrapper.ExecWrapper {
	cmd := exec.Command(g.config.BinPath, "-vv", "--progress=plain", "functions")
	cmd.Dir = g.config.TempDir

	name := fmt.Sprintf("functions-go-%s", g.config.Version())
	if g.config.Name != "" {
		name = fmt.Sprintf("%s-functions", g.config.Name)
	}

	return execwrapper.NewExecWrapper(
		name,
		cmd,
		[]string{
			"finding module configuration",
			"initializing module",
			"getModDef",
			"loading type definitions",
			"load module",
		})
}

func (g *Go) Call(callArgs []string) *execwrapper.ExecWrapper {
	cmdArgs := append([]string{"-vv", "--progress=plain", "call"}, callArgs...)

	cmd := exec.Command(g.config.BinPath, cmdArgs...)
	cmd.Dir = g.config.TempDir

	// Extract the function name from the first argument
	functionName := convertFunctionNameToTraceMarker(callArgs[0])

	name := fmt.Sprintf("call-%s-go-%s", functionName, g.config.Version())
	if g.config.Name != "" {
		name = fmt.Sprintf("%s-call", g.config.Name)
	}

	return execwrapper.NewExecWrapper(
		name,
		cmd,
		[]string{
			"finding module configuration",
			"initializing module",
			"getModDef",
			"loading type definitions",
			"load module",
			functionName,
		})
}
