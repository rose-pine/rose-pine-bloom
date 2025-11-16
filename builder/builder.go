package builder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/rose-pine/rose-pine-bloom/color"
	"github.com/rose-pine/rose-pine-bloom/config"
)

var variantValueRegex = regexp.MustCompile(`\$\((.*?)\|(.*?)\|(.*?)\)`)

func BuildTemplate(cfg *config.BuildTemplateConfig) error {
	if err := os.MkdirAll(cfg.Output, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	return createTemplates(cfg)
}

func Build(cfg *config.BuildConfig) error {
	if err := os.MkdirAll(cfg.Output, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	return generateThemes(cfg)
}

func generateThemes(cfg *config.BuildConfig) error {
	templates, err := getTemplateFiles(cfg.Template)
	if err != nil {
		return err
	}

	for _, templatePath := range templates {
		content, err := os.ReadFile(templatePath)
		if err != nil {
			return err
		}

		hasAccent := strings.Contains(string(content), cfg.Prefix+"accent")

		for _, variant := range color.Variants {
			if hasAccent {
				for _, accent := range color.Accents {
					if err := generateThemeFile(cfg, templatePath, content, variant, accent); err != nil {
						return err
					}
				}
			} else {
				if err := generateThemeFile(cfg, templatePath, content, variant, ""); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func generateThemeFile(cfg *config.BuildConfig, templatePath string, content []byte, variant color.VariantMeta, accent string) error {
	result := processTemplate(string(content), cfg, variant, accent)

	if filepath.Ext(templatePath) == ".json" {
		var prettyJSON bytes.Buffer
		if err := json.Indent(&prettyJSON, []byte(result), "", "  "); err != nil {
			return err
		}
		result = prettyJSON.String()
	}

	outputPath := buildOutputPath(cfg, templatePath, variant, accent)
	return writeFile(outputPath, []byte(result))
}

func createTemplates(cfg *config.BuildTemplateConfig) error {
	files, err := getTemplateFiles(cfg.Input)
	if err != nil {
		return err
	}

	variant := getVariant(cfg.Variant)

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		result := string(content)
		matchFound := false

		// Replace colors with variables
		for colorName, c := range variant.Colors.Iter() {
			colorValue := color.FormatColor(c, color.ColorFormat(cfg.Format), cfg.Plain, cfg.Commas, cfg.Spaces)
			before := result
			result = strings.ReplaceAll(result, colorValue, cfg.Prefix+colorName)
			if before != result {
				matchFound = true
			}
		}

		// Replace metadata
		result = strings.ReplaceAll(result, variant.Id, cfg.Prefix+"id")
		result = strings.ReplaceAll(result, variant.Name, cfg.Prefix+"name")
		result = strings.ReplaceAll(result, variant.Description, cfg.Prefix+"description")

		if !matchFound {
			fmt.Printf("\033[33mNo matches for specified format (%s). Available formats:\n\033[0m", cfg.Format)
			table, err := color.FormatsTable()
			if err != nil {
				return fmt.Errorf("failed to get formats table: %w", err)
			}
			fmt.Println(table)
		}

		outputFile := "template" + filepath.Ext(file)
		if ext := filepath.Ext(file); ext != "" {
			outputFile = "template" + ext
		}
		outputPath := filepath.Join(cfg.Output, outputFile)

		if err := writeFile(outputPath, []byte(result)); err != nil {
			return err
		}
	}

	return nil
}

func processTemplate(content string, cfg *config.BuildConfig, variant color.VariantMeta, accent string) string {
	result := content

	// Replace metadata
	result = strings.ReplaceAll(result, cfg.Prefix+"id", variant.Id)
	result = strings.ReplaceAll(result, cfg.Prefix+"name", variant.Name)
	result = strings.ReplaceAll(result, cfg.Prefix+"type", variant.Appearance)
	result = strings.ReplaceAll(result, cfg.Prefix+"appearance", variant.Appearance)
	result = strings.ReplaceAll(result, cfg.Prefix+"description", variant.Description)

	// Replace accent variables
	if accent != "" {
		result = strings.ReplaceAll(result, cfg.Prefix+"accentname", accent)
		result = strings.ReplaceAll(result, cfg.Prefix+"accent", cfg.Prefix+accent)

		if color, ok := variant.Colors.Get(accent); ok && color.On != "" {
			result = strings.ReplaceAll(result, cfg.Prefix+"onaccent", cfg.Prefix+color.On)
		}
	}

	// Replace colors and handle alpha variants
	for colorName, c := range variant.Colors.Iter() {
		varName := cfg.Prefix + colorName

		// Handle alpha variants (e.g. $base/50)
		alphaRegex := regexp.MustCompile(regexp.QuoteMeta(varName) + `/(\d+)`)
		result = alphaRegex.ReplaceAllStringFunc(result, func(match string) string {
			alpha, _ := strconv.ParseFloat(alphaRegex.FindStringSubmatch(match)[1], 64)
			colorCopy := *c
			normalizedAlpha := alpha / 100
			colorCopy.Alpha = &normalizedAlpha
			return color.FormatColor(&colorCopy, color.ColorFormat(cfg.Format), cfg.Plain, cfg.Commas, cfg.Spaces)
		})

		// Replace regular color variable
		result = strings.ReplaceAll(result, varName, color.FormatColor(c, color.ColorFormat(cfg.Format), cfg.Plain, cfg.Commas, cfg.Spaces))
	}

	// Process variant-specific values $(main|moon|dawn)
	result = variantValueRegex.ReplaceAllStringFunc(result, func(match string) string {
		parts := variantValueRegex.FindStringSubmatch(match)
		if len(parts) != 4 {
			return match
		}

		switch variant.Id {
		case "rose-pine":
			return parts[1]
		case "rose-pine-moon":
			return parts[2]
		case "rose-pine-dawn":
			return parts[3]
		default:
			return match
		}
	})

	return result
}

func getVariant(create string) color.VariantMeta {
	switch create {
	case "main":
		return color.MainVariantMeta
	case "moon":
		return color.MoonVariantMeta
	case "dawn":
		return color.DawnVariantMeta
	default:
		return color.MainVariantMeta
	}
}

func getTemplateFiles(templatePath string) ([]string, error) {
	var files []string

	info, err := os.Stat(templatePath)
	if err != nil {
		return nil, err
	}

	if info.IsDir() {
		err := filepath.Walk(templatePath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	} else {
		files = append(files, templatePath)
	}

	return files, nil
}

func writeFile(outputPath string, content []byte) error {
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return err
	}
	return os.WriteFile(outputPath, content, 0644)
}

func buildOutputPath(cfg *config.BuildConfig, templatePath string, variant color.VariantMeta, accent string) string {
	ext := filepath.Ext(templatePath)
	var outputFile, outputDir string

	if accent != "" {
		outputDir = filepath.Join(cfg.Output, variant.Id)
		outputFile = variant.Id + "-" + accent + ext
	} else {
		outputDir = cfg.Output
		outputFile = variant.Id + ext
	}

	// Handle directory templates
	if templateInfo, err := os.Stat(cfg.Template); err == nil && templateInfo.IsDir() {
		dir, _ := strings.CutPrefix(filepath.Dir(templatePath), filepath.Clean(cfg.Template))
		if accent != "" {
			outputDir = filepath.Join(cfg.Output, accent, variant.Id) + dir
		} else {
			outputDir = filepath.Join(cfg.Output, variant.Id) + dir
		}
		outputFile = filepath.Base(templatePath)
	}

	return filepath.Join(outputDir, outputFile)
}
