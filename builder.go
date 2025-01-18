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
	templateContent, err := os.ReadFile(cfg.Template)
	if err != nil {
		return fmt.Errorf("failed to read template: %w", err)
	}

	variants := []struct {
		id, name, variantType string
		colors                Variant
	}{
		{"rose-pine", "Rosé Pine", "dark", MainVariant},
		{"rose-pine-moon", "Rosé Pine Moon", "dark", MoonVariant},
		{"rose-pine-dawn", "Rosé Pine Dawn", "light", DawnVariant},
	}

	for _, v := range variants {
		if err := processVariant(cfg, templateContent, v); err != nil {
			return fmt.Errorf("failed to process %s: %w", v.id, err)
		}
	}

	return nil
}

var variantValueRegex = regexp.MustCompile(`\$\((.*?)\|(.*?)\|(.*?)\)`)

func processVariant(cfg *Config, templateContent []byte, v struct {
	id, name, variantType string
	colors                Variant
}) error {
	result := string(templateContent)

	result = strings.ReplaceAll(result, cfg.Prefix+"id", v.id)
	result = strings.ReplaceAll(result, cfg.Prefix+"name", v.name)
	result = strings.ReplaceAll(result, cfg.Prefix+"type", v.variantType)
	result = strings.ReplaceAll(result, cfg.Prefix+"description",
		"All natural pine, faux fur and a bit of soho vibes for the classy minimalist")

	for colorName, color := range v.colors.Colors {
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

		switch v.id {
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

	ext := filepath.Ext(cfg.Template)
	if ext == ".json" {
		var prettyJSON bytes.Buffer
		if err := json.Indent(&prettyJSON, []byte(result), "", "  "); err != nil {
			return fmt.Errorf("failed to format JSON: %w", err)
		}
		result = prettyJSON.String()
	}

	outputPath := filepath.Join(cfg.Output, v.id+ext)
	return os.WriteFile(outputPath, []byte(result), 0644)
}
