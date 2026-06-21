package cmd

import (
	"os"
	"strings"
	"testing"
)

func TestEnsureReadme(t *testing.T) {
	tests := []struct {
		name         string
		existing     string
		templatePath string
		prefix       string
		wantLines    []string
		notWantLines []string
	}{
		{
			name:         "no existing readme",
			existing:     "",
			templatePath: "template.json",
			prefix:       "$",
			wantLines:    []string{"bloom build template.json --prefix $"},
		},
		{
			name:         "existing readme without marker",
			existing:     "# My Theme\n\nSome description\n",
			templatePath: "template.json",
			prefix:       "$",
			wantLines:    []string{"bloom build template.json --prefix $"},
		},
		{
			name: "existing readme with marker",
			existing: "# My Theme\n\n" +
				"<!-- BLOOM_BUILD_START -->\n" +
				"This theme was built using [bloom](https://github.com/rose-pine/rose-pine-bloom):\n\n" +
				"```sh\nbloom build old-template.json\n```\n" +
				"<!-- BLOOM_BUILD_END -->\n",
			templatePath: "templates/new-template.json",
			prefix:       "$",
			wantLines:    []string{"templates/new-template.json", "--prefix $"},
			notWantLines: []string{"old-template.json"},
		},
		{
			name:         "custom prefix",
			existing:     "",
			templatePath: "template.json",
			prefix:       "#",
			wantLines:    []string{"bloom build template.json --prefix #"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			origDir, err := os.Getwd()
			if err != nil {
				t.Fatal(err)
			}
			if err := os.Chdir(tmpDir); err != nil {
				t.Fatal(err)
			}
			defer func() {
				if err := os.Chdir(origDir); err != nil {
					t.Fatal(err)
				}
			}()

			if tt.existing != "" {
				if err := os.WriteFile("README.md", []byte(tt.existing), 0644); err != nil {
					t.Fatal(err)
				}
			}

			if err := ensureReadme(tt.templatePath, tt.prefix); err != nil {
				t.Fatal(err)
			}

			content, err := os.ReadFile("README.md")
			if err != nil {
				t.Fatal(err)
			}
			got := string(content)

			for _, want := range tt.wantLines {
				if !strings.Contains(got, want) {
					t.Errorf("README should contain %q\n%s", want, got)
				}
			}
			for _, notWant := range tt.notWantLines {
				if strings.Contains(got, notWant) {
					t.Errorf("README should NOT contain %q\n%s", notWant, got)
				}
			}
		})
	}
}
