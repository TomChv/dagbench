package execwrapper

import (
	"fmt"
	"regexp"
	"time"
)

func stripANSI(s string) string {
	var ansi = regexp.MustCompile(`\x1b\[[0-9;?]*[ -/]*[@-~]`) // CSI sequences

	return ansi.ReplaceAllString(s, "")
}

// extractTimeFromTraceLine extracts the time from a trace line
// We assume the trace is with `--progress=plain`
func extractTimeFromTraceLine(line string) (time.Duration, error) {
	re := regexp.MustCompile(`\[(.*?)\]`)
	match := re.FindStringSubmatch(line)

	if len(match) > 1 {
		return time.ParseDuration(match[1])
	}

	return time.Duration(0), fmt.Errorf("failed to extract time from trace line: %s", line)
}
