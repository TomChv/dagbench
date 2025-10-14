package benchmark

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type ParserOpts struct {
	UseCmdDuration bool
}

func ParseFile(path string, opts ParserOpts) (*File, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read benchmark file %s: %w", path, err)
	}

	metadata := make(map[string]string)
	var results []Entry
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Ignore empty line
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "#") {
			name, value := parseMetadata(line)
			metadata[name] = value

			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 3 {
			continue
		}

		name := parseEntryName(parts[0])

		// Skip cmdDuration entries if not requested or only use it if requested
		if !opts.UseCmdDuration && strings.HasSuffix(name, "cmdDuration") {
			continue
		} else if opts.UseCmdDuration && !strings.HasSuffix(name, "cmdDuration") {
			continue
		}

		timeStr := strings.TrimSuffix(parts[2], "s/op")
		timeVal, _ := strconv.ParseFloat(timeStr, 64)

		results = append(results, Entry{
			Name:        name,
			DurationSec: timeVal,
		})
	}

	return &File{
		Path:     path,
		Metadata: metadata,
		Entries:  results,
	}, scanner.Err()
}

func parseMetadata(line string) (string, string) {
	line = strings.TrimPrefix(line, "#")
	line = strings.TrimSpace(line)
	parts := strings.Split(line, ":")

	return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
}

func parseEntryName(name string) string {
	parts := strings.Split(name, "/")

	switch len(parts) {
	case 2:
		// benchmark/operation -> operation
		return parts[1]
	case 3:
		// benchmark/operation/span -> operation/span
		return fmt.Sprintf("%s/%s", parts[1], parts[2])
	default:
		return name
	}
}
