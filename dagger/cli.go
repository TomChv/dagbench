package dagger

import (
	"bufio"
	"fmt"
	"os/exec"

	"quartz/dagbench.io/config"
	"quartz/dagbench.io/hook"
)

type Command struct {
	bin     string
	args    []string
	workdir string
}

func NewCLI(conf *config.Config) *Command {
	args := []string{"-vv", "--progress=plain"}

	if conf.Cloud {
		args = append(args, "--cloud")
	}

	if conf.Module != "" {
		args = append(args, "-m", conf.Module)
	}

	return &Command{
		bin:     conf.BinPath,
		args:    args,
		workdir: conf.Workdir,
	}
}

func (c *Command) Exec(args []string, hooks ...hook.Hook) error {
	cmd := exec.Command(c.bin, append(c.args, args...)...)

	if c.workdir != "" {
		cmd.Dir = c.workdir
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to get stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}

	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		line := stripANSI(scanner.Text())

		for _, hook := range hooks {
			if err := hook.Hook(line); err != nil {
				return fmt.Errorf("failed to execute benchmark hook: %w", err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("scanner error:", err)
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("failed to wait for command: %w", err)
	}

	return nil
}

func (c *Command) PruneCache() error {
	cmd := exec.Command(c.bin, "core", "engine", "local-cache", "prune")

	return cmd.Run()
}

func (c *Command) Init(moduleName string, sdk string, hooks ...hook.Hook) error {
	return c.Exec(
		[]string{"init", "--sdk", sdk, "--source=.", "--name", moduleName},
		hooks...,
	)
}
