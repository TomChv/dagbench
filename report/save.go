package report

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

func (r *Report) SaveAsCSVAt(dirpath string) error {
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
		if err := w.Write([]string{
			k, v.String(),
			strconv.FormatFloat(v.Seconds(), 'f', -1, 64),
			strconv.FormatInt(v.Milliseconds(), 10),
		}); err != nil {
			return fmt.Errorf("failed to write row: %w", err)
		}
	}

	w.Flush()

	return nil
}

func (r *Report) SaveOutputAt(dirpath string) error {
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
