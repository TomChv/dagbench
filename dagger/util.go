package dagger

import (
	"regexp"
)

func stripANSI(s string) string {
	var ansi = regexp.MustCompile(`\x1b\[[0-9;?]*[ -/]*[@-~]`) // CSI sequences

	return ansi.ReplaceAllString(s, "")
}
