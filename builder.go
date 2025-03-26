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

	if err := generateVariants(cfg); err != nil {
		return fmt.Errorf("failed to generate variants: %w", err)
	}

	return nil
}

func generateVariants(cfg *Config) error {

	variants := []struct {
		id, name, variantType string
		colors                Variant
	}{
		{"rose-pine", "Rosé Pine", "dark", MainVariant},
		{"rose-pine-moon", "Rosé Pine Moon", "dark", MoonVariant},
		{"rose-pine-dawn", "Rosé Pine Dawn", "light", DawnVariant},
	}
	accents := []string{
		"love", "gold", "rose", "pine", "foam", "iris",
	}

	templateFileInfo, err := os.Stat(cfg.Template)
	if err != nil {
		return fmt.Errorf("failed to open template: %w", err)
	}

	var templates []string

	if templateFileInfo.IsDir() {
		filepath.Walk(cfg.Template, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				templates = append(templates, path)
			}
			return nil
		})

	} else {
		templates = append(templates, cfg.Template)
	}

	for _, template := range templates {
		templateFileInfo, _ := os.Stat(template)
		if templateFileInfo.IsDir() {
			continue
		} // skip if dir (redudant)

		templateContent, err := os.ReadFile(template)
		if err != nil {
			return fmt.Errorf("failed to read template: %w", err)
		}
		for _, variant := range variants {
			if cfg.Accents {
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

	for colorName, color := range variant.colors.Colors {
		varName := cfg.Prefix + colorName

		alphaRegex := regexp.MustCompile(regexp.QuoteMeta(varName) + `/(\d+)`)
		result = alphaRegex.ReplaceAllStringFunc(result, func(match string) string {
			alpha, _ := strconv.ParseFloat(alphaRegex.FindStringSubmatch(match)[1], 64)
			colorCopy := *color
			normalizedAlpha := alpha / 100
			colorCopy.Alpha = &normalizedAlpha
			return formatColor(&colorCopy, ColorFormat(cfg.Format), cfg.StripSpaces)
		})

		result = strings.ReplaceAll(result, varName, formatColor(color, ColorFormat(cfg.Format), cfg.StripSpaces))
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
	if cfg.Accents {
		outputDir = cfg.Output + "/" + variant.id
		outputFile = variant.id + "-" + accent + ext
	} else {
		outputDir = cfg.Output
		outputFile = variant.id + ext
	}

	templateFileInfo, _ := os.Stat(cfg.Template)

	if templateFileInfo.IsDir() {
		dir, _ := strings.CutPrefix(filepath.Dir(template), filepath.Clean(cfg.Template))
		if cfg.Accents {
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
