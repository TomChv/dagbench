package execwrapper

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
)

type ExecWrapper struct {
	name        string
	cmd         *exec.Cmd
	spanMarkers []string
}

func NewExecWrapper(name string, cmd *exec.Cmd, spanMarkers []string) *ExecWrapper {
	return &ExecWrapper{
		name:        name,
		cmd:         cmd,
		spanMarkers: spanMarkers,
	}
}

func (e *ExecWrapper) Exec() (*ExecReport, error) {
	fmt.Println("Executing", e.name)

	report := &ExecReport{
		name:   e.name,
		values: make(map[string]string),
	}

	stderr, err := e.cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to get stderr pipe: %w", err)
	}

	if err := e.cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start command: %w", err)
	}

	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		line := stripANSI(scanner.Text())

		report.stderr += line + "\n"

		for _, marker := range e.spanMarkers {
			if strings.Contains(line, marker) && strings.Contains(line, "DONE") {
				report.values[marker] = extractTimeFromTraceLine(line)

				break
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("scanner error:", err)
	}

	if err := e.cmd.Wait(); err != nil {
		return nil, fmt.Errorf("failed to wait for command: %w", err)
	}

	return report, nil
}
