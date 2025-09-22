package execwrapper

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type ExecReport struct {
	name string

	stderr string

	values map[string]string
}

func (r *ExecReport) String() string {
	var res strings.Builder

	fmt.Fprintf(&res, "Report %s results\n", r.name)

	for k, v := range r.values {
		fmt.Fprintf(&res, "%s\t: %s\n", k, v)
	}

	return res.String()
}

func (r *ExecReport) Stderr() string {
	return r.stderr
}

func (r *ExecReport) SaveAsCSVAt(dirpath string) error {
	fileName := fmt.Sprintf("%s-results.csv", r.name)
	absPath, err := filepath.Abs(dirpath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	path := filepath.Join(absPath, fileName)
	fmt.Printf("Saving report at %s\n", path)

	if err := os.MkdirAll(absPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	w := csv.NewWriter(file)
	w.Comma = ';'

	if err := w.Write([]string{"span", "time", "time (s)", "time (ms)"}); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	for k, v := range r.values {
		parsedDuration, err := time.ParseDuration(v)
		if err != nil {
			return fmt.Errorf("failed to parse duration: %w", err)
		}

		if err := w.Write([]string{
			k, v,
			strconv.FormatFloat(parsedDuration.Seconds(), 'f', -1, 64),
			strconv.FormatInt(parsedDuration.Milliseconds(), 10),
		}); err != nil {
			return fmt.Errorf("failed to write row: %w", err)
		}
	}

	w.Flush()

	return nil
}

func (r *ExecReport) SaveOutputAt(dirpath string) error {
	fileName := fmt.Sprintf("%s-output.txt", r.name)
	absPath, err := filepath.Abs(dirpath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	path := filepath.Join(absPath, fileName)
	fmt.Printf("Saving output at %s\n", path)

	if err := os.MkdirAll(absPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	if _, err := file.Write([]byte(r.stderr)); err != nil {
		return fmt.Errorf("failed to write config to file: %w", err)
	}

	return nil
}
