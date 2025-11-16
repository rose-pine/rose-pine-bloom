package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/rose-pine/rose-pine-bloom/builder"
	"github.com/spf13/cobra"
)

var (
	variant string
	output  string
)

const (
	startMarker = "<!-- BLOOM_BUILD_START -->"
	endMarker   = "<!-- BLOOM_BUILD_END -->"
)

var initCmd = &cobra.Command{
	Use:   "init [theme-file]",
	Short: "Initialise new theme",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Initialising theme...")

		if err := ensureReadme(); err != nil {
			fmt.Fprintf(os.Stderr, "Error updating README: %v\n", err)
		} else {
			fmt.Println("Updated README.md")
		}

		licenseCreated, err := ensureLicense()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error updating LICENSE: %v\n", err)
		} else if licenseCreated {
			fmt.Println("Updated LICENSE")
		}

		if len(args) > 0 {
			themeFile := args[0]
			fmt.Printf("Creating template from %s...\n", themeFile)

			err := builder.BuildTemplate(&builder.TemplateOptions{
				Input:   themeFile,
				Output:  output,
				Variant: variant,
				Prefix:  "$",
				Format:  "hex",
				Plain:   false,
				Commas:  true,
				Spaces:  true,
			})
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating template: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("Template created in %s\n", output)
		}

		fmt.Println("Theme initialised")
	},
}

func ensureReadme() error {
	fileName, err := findAndNormalizeFile("README.md")
	if err != nil {
		return err
	}
	content, err := os.ReadFile(fileName)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	contentStr := string(content)

	section := fmt.Sprintf("%s\nThis theme was built using [bloom](https://github.com/rose-pine/rose-pine-bloom):\n\n```sh\nbloom build <template>\n```\n%s", startMarker, endMarker)

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

func ensureLicense() (bool, error) {
	fileName, err := findAndNormalizeFile("LICENSE")
	if err != nil {
		return false, err
	}

	if existingContent, err := os.ReadFile(fileName); err == nil && len(existingContent) > 0 {
		return false, nil
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
	return true, os.WriteFile(fileName, []byte(contentStr), 0644)
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

		fmt.Printf("Renamed %s to %s\n", actualName, targetName)
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

func init() {
	initCmd.Flags().StringVarP(&variant, "variant", "v", "main", "theme variant (main, moon, dawn)")
	initCmd.Flags().StringVarP(&output, "output", "o", ".", "template output directory")
	RootCmd.AddCommand(initCmd)
}
