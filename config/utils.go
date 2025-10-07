package config

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

func generateDefaultWorkdir(name string) string {
	return filepath.Join(os.TempDir(), "dagbench-workdir", name, uuid.NewString())
}

func getDaggerVersion(binPath string) (string, error) {
	version, err := exec.Command(binPath, "version").Output()
	if err != nil {
		return "", fmt.Errorf("failed to execute dagger version: %w", err)
	}

	return strings.TrimSpace(string(version)), nil
}
