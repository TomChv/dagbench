package execwrapper

import (
	"regexp"
)

func stripANSI(s string) string {
	var ansi = regexp.MustCompile(`\x1b\[[0-9;?]*[ -/]*[@-~]`) // CSI sequences

	return ansi.ReplaceAllString(s, "")
}

func extractTimeFromTraceLine(line string) string {
	line = stripANSI(line)

	re := regexp.MustCompile(`\[(.*?)\]`)
	match := re.FindStringSubmatch(line)

	if len(match) > 1 {
		return match[1]
	}

	return "?"
}
