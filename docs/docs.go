package docs

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

const (
	startMarker = "<!-- BLOOM_BUILD_START -->"
	endMarker   = "<!-- BLOOM_BUILD_END -->"
)

func EnsureReadmeWithBuildCommand(cmd, version string) error {
	fileName, err := findAndNormalizeFile("README.md")
	if err != nil {
		return err
	}
	content, err := os.ReadFile(fileName)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	contentStr := string(content)

	versionSuffix := ""
	if version != "" {
		versionSuffix = "@" + version
	}

	section := fmt.Sprintf("%s\nThis theme was built using [rose-pine-bloom](https://github.com/rose-pine/rose-pine-bloom):\n\n```sh\n%s\n```\n\nInstall via [goblin](https://goblin.run):\n\n```sh\ncurl -sf http://goblin.run/github.com/rose-pine/rose-pine-bloom%s | sh\n```\n%s", startMarker, cmd, versionSuffix, endMarker)

	markerRe := regexp.MustCompile(
		regexp.QuoteMeta(startMarker) + `(?s).*?` + regexp.QuoteMeta(endMarker),
	)
	if markerRe.MatchString(contentStr) {
		contentStr = markerRe.ReplaceAllString(contentStr, section)
	} else {
		if len(contentStr) > 0 && !strings.HasSuffix(contentStr, "\n") {
			contentStr += "\n"
		}
		contentStr += "\n" + section + "\n"
	}

	return os.WriteFile(fileName, []byte(contentStr), 0644)
}

func EnsureLicense() error {
	fileName, err := findAndNormalizeFile("LICENSE")
	if err != nil {
		return err
	}

	year := time.Now().Year()
	contentStr := fmt.Sprintf(`MIT License

Copyright (c) %d Rosé Pine

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
`, year)
	return os.WriteFile(fileName, []byte(contentStr), 0644)
}

func findAndNormalizeFile(targetName string) (string, error) {
	entries, err := os.ReadDir(".")
	if err != nil {
		return "", err
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.EqualFold(entry.Name(), targetName) {
			continue
		}

		actualName := entry.Name()
		if actualName == targetName {
			return targetName, nil
		}

		if err := renameFile(actualName, targetName); err != nil {
			return "", err
		}

		fmt.Printf("renamed %s to %s\n", actualName, targetName)
		return targetName, nil
	}

	return targetName, nil
}

func renameFile(from, to string) error {
	if !isGitRepo() {
		return os.Rename(from, to)
	}

	cmd := exec.Command("git", "mv", from, to)
	if err := cmd.Run(); err == nil {
		return nil
	}

	if renameErr := os.Rename(from, to); renameErr != nil {
		return fmt.Errorf("failed to rename %s to %s", from, to)
	}

	return nil
}

func isGitRepo() bool {
	_, err := os.Stat(".git")
	return err == nil
}
