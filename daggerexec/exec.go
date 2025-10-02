package daggerexec

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"quartz/dagbench.io/config"
	"quartz/dagbench.io/daggerexec/hook"
)

func execDagger(c *config.Config, args []string, hooks ...hook.Hook) error {
	// Add progress plain & verbose
	baseArgs := []string{"-vv", "--progress=plain"}

	// Add --cloud if enabled
	if c.Cloud {
		baseArgs = append(baseArgs, "--cloud")
	}

	// Add the command args
	args = append(baseArgs, args...)

	cmd := exec.Command(c.BinPath, args...)
	cmd.Dir = c.Workdir

	if err := os.MkdirAll(c.Workdir, 0755); err != nil {
		return fmt.Errorf("failed to create workdir: %w", err)
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
