package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func Build(cfg *Config) error {
	if err := os.MkdirAll(cfg.Output, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	if cfg.Create != "" {
		if err := generateTemplates(cfg); err != nil {
			return fmt.Errorf("failed to create template: %w", err)
		}
	} else {
		if err := generateVariants(cfg); err != nil {
			return fmt.Errorf("failed to generate variants: %w", err)
		}
	}

	return nil
}

func generateTemplates(cfg *Config) error {

	var files []string
	files, err := getFiles(cfg.Template)
	if err != nil {
		return err
	}

	for _, file := range files {
		fileContent, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}
		if err := createTemplate(cfg, file, fileContent); err != nil {
			return fmt.Errorf("failed to create template from file %s: %w", file, err)
		}
	}

	return nil
}

var variants = []struct {
	id, name, variantType string
	colors                Variant
}{
	{"rose-pine", "Rosé Pine", "dark", MainVariant},
	{"rose-pine-moon", "Rosé Pine Moon", "dark", MoonVariant},
	{"rose-pine-dawn", "Rosé Pine Dawn", "light", DawnVariant},
}

var accents = []string{
	"love", "gold", "rose", "pine", "foam", "iris",
}

var onAccentMapping = map[string]map[string]string{
	"rose-pine": {
		"love": "text",
		"gold": "surface",
		"rose": "surface",
		"pine": "text",
		"foam": "surface",
		"iris": "surface",
	},
	"rose-pine-moon": {
		"love": "text",
		"gold": "surface",
		"rose": "surface",
		"pine": "text",
		"foam": "surface",
		"iris": "surface",
	},
	"rose-pine-dawn": {
		"love": "surface",
		"gold": "surface",
		"rose": "surface",
		"pine": "surface",
		"foam": "surface",
		"iris": "surface",
	},
}

func createTemplate(cfg *Config, file string, fileContent []byte) error {
	result := string(fileContent)

	matchFound := false

	variant := variants[0]
	switch cfg.Create {
	case "main":
		variant = variants[0]
	case "moon":
		variant = variants[1]
	case "dawn":
		variant = variants[2]
	}
	snapshot := result

	for colorName, color := range variant.colors.Colors {
		result = strings.ReplaceAll(result, formatColor(color, ColorFormat(cfg.Format), cfg.Commas, cfg.Spaces), cfg.Prefix+colorName)
		if snapshot != result {
			matchFound = true
		}
	}
	result = strings.ReplaceAll(result, variant.id, cfg.Prefix+"id")
	result = strings.ReplaceAll(result, variant.name, cfg.Prefix+"name")
	// result = strings.ReplaceAll(result, variant.variantType, cfg.Prefix+"type") // likely not a good idea
	result = strings.ReplaceAll(result, "All natural pine, faux fur and a bit of soho vibes for the classy minimalist",
		cfg.Prefix+"description")

	if !matchFound {
		var Yellow = "\033[33m"
		var Reset = "\033[0m"
		fmt.Printf(Yellow+"No matches for specified format (%s). Available formats:\n"+Reset, ColorFormat(cfg.Format))
		PrintFormatsTable()
	}

	var outputFile, outputDir string
	outputDir = cfg.Output
	outputFile = filepath.Base(file)

	parts := strings.SplitN(outputFile, ".", 2)
	if len(parts) == 2 {
		outputFile = "template." + parts[1]
	} else {
		outputFile = "template" + filepath.Ext(outputFile)
	}

	outputPath := filepath.Join(outputDir, outputFile)

	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	return os.WriteFile(outputPath, []byte(result), 0644)
}

func generateVariants(cfg *Config) error {

	var templates []string
	templates, err := getFiles(cfg.Template)
	if err != nil {
		return err
	}

	for _, template := range templates {
		templateContent, err := os.ReadFile(template)
		if err != nil {
			return fmt.Errorf("failed to read template: %w", err)
		}

		hasAccent := false
		if strings.Contains(string(templateContent), cfg.Prefix+"accent") {
			hasAccent = true
		}

		for _, variant := range variants {
			if hasAccent {
				for _, accent := range accents {
					if err := processVariant(cfg, template, templateContent, accent, variant); err != nil {
						return fmt.Errorf("failed to process %s: %w", variant.id, err)
					}
				}
			} else {
				if err := processVariant(cfg, template, templateContent, "accent", variant); err != nil {
					return fmt.Errorf("failed to process %s: %w", variant.id, err)
				}
			}
		}
	}

	return nil
}

var variantValueRegex = regexp.MustCompile(`\$\((.*?)\|(.*?)\|(.*?)\)`)

func processVariant(cfg *Config, template string, templateContent []byte, accent string, variant struct {
	id, name, variantType string
	colors                Variant
}) error {
	result := string(templateContent)

	result = strings.ReplaceAll(result, cfg.Prefix+"id", variant.id)
	result = strings.ReplaceAll(result, cfg.Prefix+"name", variant.name)
	result = strings.ReplaceAll(result, cfg.Prefix+"type", variant.variantType)
	result = strings.ReplaceAll(result, cfg.Prefix+"description",
		"All natural pine, faux fur and a bit of soho vibes for the classy minimalist")

	result = strings.ReplaceAll(result, cfg.Prefix+"accentname", accent)
	result = strings.ReplaceAll(result, cfg.Prefix+"accent", cfg.Prefix+accent)

	onAccentColor := ""
	if byVariant, ok := onAccentMapping[variant.id]; ok {
		if mapped, ok := byVariant[accent]; ok {
			onAccentColor = mapped
		}
	}
	if onAccentColor != "" {
		result = strings.ReplaceAll(result, cfg.Prefix+"onaccent", cfg.Prefix+onAccentColor)
	}

	for colorName, color := range variant.colors.Colors {
		varName := cfg.Prefix + colorName

		alphaRegex := regexp.MustCompile(regexp.QuoteMeta(varName) + `/(\d+)`)
		result = alphaRegex.ReplaceAllStringFunc(result, func(match string) string {
			alpha, _ := strconv.ParseFloat(alphaRegex.FindStringSubmatch(match)[1], 64)
			colorCopy := *color
			normalizedAlpha := alpha / 100
			colorCopy.Alpha = &normalizedAlpha
			return formatColor(&colorCopy, ColorFormat(cfg.Format), cfg.Commas, cfg.Spaces)
		})

		result = strings.ReplaceAll(result, varName, formatColor(color, ColorFormat(cfg.Format), cfg.Commas, cfg.Spaces))
	}

	result = variantValueRegex.ReplaceAllStringFunc(result, func(match string) string {
		parts := variantValueRegex.FindStringSubmatch(match)
		if len(parts) != 4 {
			return match
		}

		switch variant.id {
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

	ext := filepath.Ext(template) // tmp
	if ext == ".json" {
		var prettyJSON bytes.Buffer
		if err := json.Indent(&prettyJSON, []byte(result), "", "  "); err != nil {
			return fmt.Errorf("failed to format JSON: %w", err)
		}
		result = prettyJSON.String()
	}
	var outputFile, outputDir string
	if accent != "accent" {
		outputDir = cfg.Output + "/" + variant.id
		outputFile = variant.id + "-" + accent + ext
	} else {
		outputDir = cfg.Output
		outputFile = variant.id + ext
	}

	templateFileInfo, _ := os.Stat(cfg.Template)

	if templateFileInfo.IsDir() {
		dir, _ := strings.CutPrefix(filepath.Dir(template), filepath.Clean(cfg.Template))
		if accent != "accent" {
			outputDir += "/" + accent
		} else {
			outputDir += "/" + variant.id
		}
		outputDir += dir
		outputFile = filepath.Base(template)
	}

	outputPath := filepath.Join(outputDir, outputFile)

	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	return os.WriteFile(outputPath, []byte(result), 0644)
}

func getFiles(basefile string) ([]string, error) {
	var files []string

	templateFileInfo, err := os.Stat(basefile)
	if err != nil {
		return nil, fmt.Errorf("failed to open template: %w", err)
	}

	if templateFileInfo.IsDir() {
		filepath.Walk(basefile, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				files = append(files, path)
			}
			return nil
		})

	} else {
		files = append(files, basefile)
	}
	return files, nil
}
