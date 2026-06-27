package builder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rose-pine/rose-pine-bloom/color"
)

type BuildOpts struct {
	DefaultFormat color.ColorFormat
	Prefix        rune
	Plain         bool
	Commas        bool
	Spaces        bool
}

func DiscoverTemplates(path string) ([]string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		return []string{path}, nil
	}

	var files []string
	err = filepath.Walk(path, func(p string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !fi.IsDir() {
			files = append(files, p)
		}
		return nil
	})
	return files, err
}

func buildOutPath(templatePath string, outputPath string, variant color.VariantMeta, accent string) string {
	ext := filepath.Ext(templatePath)

	if info, err := os.Stat(templatePath); err == nil && info.IsDir() {
		if accent != "" {
			return filepath.Join(outputPath, variant.Id+"-"+accent+ext)
		}
		return filepath.Join(outputPath, variant.Id+ext)
	}

	if accent != "" {
		return filepath.Join(outputPath, variant.Id, variant.Id+"-"+accent+ext)
	}
	return filepath.Join(outputPath, variant.Id+ext)
}

func writeFile(path string, content []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, content, 0o644)
}

func writeTemplateFile(path string, content string) error {
	result, err := postProcessTemplate(path, content)
	if err != nil {
		return fmt.Errorf("post-processing hook failed: %w", err)
	}

	return writeFile(path, ([]byte)(result))
}

func hasAccentCapture(captures []Capture) bool {
	for _, capture := range captures {
		if rc, ok := capture.(RoleCapture); ok && rc.role == "accent" {
			return true
		}
	}
	return false
}

func postProcessTemplate(templatePath string, content string) (string, error) {
	result := content

	if filepath.Ext(templatePath) == ".json" {
		var buf bytes.Buffer
		if err := json.Indent(&buf, []byte(result), "", "  "); err != nil {
			return "", fmt.Errorf("failed to format template output as JSON")
		}
		result = buf.String()
	}

	return result, nil
}

func Build(templatePath string, outPath string, opts *BuildOpts) error {
	scannerOpts := ScannerOpts{Prefix: opts.Prefix}
	templates, err := DiscoverTemplates(templatePath)
	if err != nil {
		return fmt.Errorf("failed to discover templates: %w", err)
	}

	for _, templatePath := range templates {
		bytes, err := os.ReadFile(templatePath)
		if err != nil {
			return fmt.Errorf("failed to read template at `%s`: %w", templatePath, err)
		}
		content := string(bytes)
		captures, err := Scan(content, scannerOpts)
		if err != nil {
			return fmt.Errorf("failed to scan template: %w", err)
		}

		for _, variant := range color.Variants {
			if hasAccentCapture(captures) {
				for _, accent := range color.Accents {
					result, err := substituteCaptures(content, captures, variant, opts, accent)
					if err != nil {
						return err
					}
					if err := writeTemplateFile(buildOutPath(templatePath, outPath, variant, accent), result); err != nil {
						return err
					}
				}
			} else {
				result, err := substituteCaptures(content, captures, variant, opts, "")
				if err != nil {
					return err
				}

				if err := writeTemplateFile(buildOutPath(templatePath, outPath, variant, ""), result); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
