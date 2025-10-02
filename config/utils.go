package config

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	DefaultBinPath = "dagger"
	DefaultOutDir  = "./dagbench-out"
)

var defaultCommands = map[string][]string{
	"functions": {},
	"develop":   {},
	"default-call": {
		"call",
		"container-echo",
		"--string-arg=foo",
	},
}

func generateDefaultWorkdir(name string) string {
	return filepath.Join(os.TempDir(), "dagbench-workdir", name)
}

func getDaggerVersion(binPath string) (string, error) {
	version, err := exec.Command(binPath, "version").Output()
	if err != nil {
		return "", fmt.Errorf("failed to execute dagger version: %w", err)
	}

	return strings.TrimSpace(string(version)), nil
}
