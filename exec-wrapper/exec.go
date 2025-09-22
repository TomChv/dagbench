package execwrapper

import (
	"bufio"
	"fmt"
	"os/exec"
	"quartz/dagbenchmark.io/report"
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

func (e *ExecWrapper) Exec() (*report.Report, error) {
	fmt.Println("Executing ", e.name)

	return e.exec()
}

func (e *ExecWrapper) exec() (*report.Report, error) {
	report := report.New(e.name)

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

		report.AddStderr(line + "\n")

		for _, marker := range e.spanMarkers {
			if strings.Contains(line, marker) && strings.Contains(line, "DONE") {
				duration, err := extractTimeFromTraceLine(line)
				if err != nil {
					return nil, fmt.Errorf("failed to extract time from trace line: %w", err)
				}

				report.AddValue(marker, duration)

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
