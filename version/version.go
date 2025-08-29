package version

import (
	"os/exec"
	"strings"
)

// TODO this needs to be replaced with a proper versioning system
func GetCurrentVersion() string {
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	version := strings.TrimSpace(string(output))
	return version
}
