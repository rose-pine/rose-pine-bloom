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
)

// Options contains configuration for building themes from templates.
type Options struct {
	Template string
	Output   string
	Prefix   string
	Format   string
	Plain    bool
	Commas   bool
	Spaces   bool
}

// TemplateOptions contains configuration for creating templates from existing theme files.
type TemplateOptions struct {
	Input   string
	Output  string
	Variant string
	Prefix  string
	Format  string
	Plain   bool
	Commas  bool
	Spaces  bool
}

var variantValueRegex = regexp.MustCompile(`\$\((.*?)\|(.*?)\|(.*?)\)`)

func BuildTemplate(cfg *TemplateOptions) error {
	if err := os.MkdirAll(cfg.Output, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	return createTemplates(cfg)
}

func Build(cfg *Options) error {
	if err := os.MkdirAll(cfg.Output, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	return generateThemes(cfg)
}

func generateThemes(cfg *Options) error {
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

func generateThemeFile(cfg *Options, templatePath string, content []byte, variant color.VariantMeta, accent string) error {
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

func createTemplates(cfg *TemplateOptions) error {
	files, err := getTemplateFiles(cfg.Input)
	if err != nil {
		return err
	}

	variant := getVariant(cfg.Variant)

	for _, file := range files {
		rawContent, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		content := string(rawContent)

		data := []string{}

		// Replace colors with variables
		for colorName, c := range variant.Colors.Iter() {
			colorValue := color.FormatColor(c, color.ColorFormat(cfg.Format), cfg.Plain, cfg.Commas, cfg.Spaces)
			data = append(data, colorValue, cfg.Prefix+colorName)
		}

		// Replace metadata
		data = append(data, variant.Id, cfg.Prefix+"id")
		data = append(data, variant.Name, cfg.Prefix+"name")
		data = append(data, variant.Description, cfg.Prefix+"description")

		replacer := strings.NewReplacer(data...)
		result := replacer.Replace(content)

		if content == result {
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

func processTemplate(content string, cfg *Options, variant color.VariantMeta, accent string) string {
	data := []string{}

	// Replace metadata
	data = append(data, cfg.Prefix+"id", variant.Id)
	data = append(data, cfg.Prefix+"name", variant.Name)
	data = append(data, cfg.Prefix+"type", variant.Appearance)
	data = append(data, cfg.Prefix+"appearance", variant.Appearance)
	data = append(data, cfg.Prefix+"description", variant.Description)

	// Replace accent variables
	if accent != "" {
		data = append(data, cfg.Prefix+"accentname", accent)

		if c, ok := variant.Colors.Get(accent); ok {
			accentColor := color.FormatColor(c, color.ColorFormat(cfg.Format), cfg.Plain, cfg.Commas, cfg.Spaces)
			data = append(data, cfg.Prefix+"accent", accentColor)
			if c.On != "" {
				if oc, ok := variant.Colors.Get(c.On); ok {
					accentOnColor := color.FormatColor(oc, color.ColorFormat(cfg.Format), cfg.Plain, cfg.Commas, cfg.Spaces)
					data = append(data, cfg.Prefix+"onaccent", accentOnColor)
				}
			}
		}
	}

	// Replace colors and handle alpha variants
	for colorName, c := range variant.Colors.Iter() {
		varName := cfg.Prefix + colorName

		// Handle alpha variants (e.g. $base/50)
		alphaRegex := regexp.MustCompile(regexp.QuoteMeta(varName) + `/(\d+)`)
		matches := alphaRegex.FindAllStringSubmatch(content, -1)
		seen := make(map[string]bool)
		for _, match := range matches {
			if seen[match[0]] {
				continue
			}
			seen[match[0]] = true
			alpha, _ := strconv.ParseFloat(match[1], 64)
			colorCopy := *c
			normalizedAlpha := alpha / 100
			colorCopy.Alpha = &normalizedAlpha
			data = append(data, match[0], color.FormatColor(&colorCopy, color.ColorFormat(cfg.Format), cfg.Plain, cfg.Commas, cfg.Spaces))
		}

		// Replace regular color variable
		data = append(data, varName, color.FormatColor(c, color.ColorFormat(cfg.Format), cfg.Plain, cfg.Commas, cfg.Spaces))
	}

	replacer := strings.NewReplacer(data...)
	result := replacer.Replace(content)

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

func buildOutputPath(cfg *Options, templatePath string, variant color.VariantMeta, accent string) string {
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
