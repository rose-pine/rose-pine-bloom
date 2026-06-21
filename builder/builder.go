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

type Options struct {
	Template string
	Output   string
	Prefix   string
	Format   string
	Plain    bool
	Commas   bool
	Spaces   bool
}

type TemplateOptions struct {
	Input   string
	Output  string
	Variant string
	Prefix  string
	Format  string
	Plain   bool
	Commas  bool
	Spaces  bool

	DetectedFormat string
	TemplatePath   string
}

const (
	warnColor  = "\033[33m"
	resetColor = "\033[0m"
)

var variantValueRe = regexp.MustCompile(`\$\((.*?)\|(.*?)\|(.*?)\)`)

func BuildTemplate(cfg *TemplateOptions) error {
	if err := os.MkdirAll(cfg.Output, 0755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	return createTemplates(cfg)
}

func Build(cfg *Options) error {
	if err := os.MkdirAll(cfg.Output, 0755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	return generateThemes(cfg)
}

func generateThemes(cfg *Options) error {
	templates, err := templateFiles(cfg.Template)
	if err != nil {
		return err
	}

	for _, tp := range templates {
		content, err := os.ReadFile(tp)
		if err != nil {
			return err
		}

		hasAccent := strings.Contains(string(content), cfg.Prefix+"accent")

		for _, v := range color.Variants {
			if hasAccent {
				for _, accent := range color.Accents {
					if err := generateThemeFile(cfg, tp, content, v, accent); err != nil {
						return err
					}
				}
			} else {
				if err := generateThemeFile(cfg, tp, content, v, ""); err != nil {
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
		var buf bytes.Buffer
		if err := json.Indent(&buf, []byte(result), "", "  "); err != nil {
			return err
		}
		result = buf.String()
	}

	outputPath := buildOutputPath(cfg, templatePath, variant, accent)
	return writeFile(outputPath, []byte(result))
}

func detectFormatOptions(content string, variant color.VariantMeta) (color.ColorFormat, bool, bool, bool) {
	bestFmt := color.ColorFormat("hex")
	bestPlain, bestCommas, bestSpaces := false, true, true
	bestCount := 0

	for _, f := range color.AllFormats {
		fmt := color.ColorFormat(f)
		for _, plain := range []bool{false, true} {
			for _, spaces := range []bool{true, false} {
				for _, commas := range []bool{true, false} {
					count := 0
					for _, c := range variant.Colors {
						val := color.FormatColor(c, fmt, plain, commas, spaces)
						if strings.Contains(content, val) {
							count++
						}
					}
					if count > bestCount {
						bestFmt, bestPlain, bestCommas, bestSpaces = fmt, plain, commas, spaces
						bestCount = count
					}
				}
			}
		}
	}

	return bestFmt, bestPlain, bestCommas, bestSpaces
}

func createTemplates(cfg *TemplateOptions) error {
	files, err := templateFiles(cfg.Input)
	if err != nil {
		return err
	}

	variant := color.MainVariantMeta
	switch cfg.Variant {
	case "moon":
		variant = color.MoonVariantMeta
	case "dawn":
		variant = color.DawnVariantMeta
	}

	for _, file := range files {
		raw, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		content := string(raw)

		formatStr, plain, commas, spaces := cfg.Format, cfg.Plain, cfg.Commas, cfg.Spaces
		if formatStr == "" {
			var df color.ColorFormat
			df, plain, commas, spaces = detectFormatOptions(content, variant)
			formatStr = string(df)
		}

		data := []string{}

		for name, c := range variant.Colors {
			val := color.FormatColor(c, color.ColorFormat(formatStr), plain, commas, spaces)
			data = append(data, val, cfg.Prefix+name)
		}

		data = append(data, variant.Id, cfg.Prefix+"id")
		data = append(data, variant.Name, cfg.Prefix+"name")
		data = append(data, variant.Description, cfg.Prefix+"description")

		result := strings.NewReplacer(data...).Replace(content)

		if content == result {
			fmt.Printf("%sNo matches for format %q. Available formats:\n  %s%s\n", warnColor, formatStr, strings.Join(color.AllFormats, ", "), resetColor)
		}

		ext := filepath.Ext(file)
		outputFile := "template" + ext
		outputPath := filepath.Join(cfg.Output, outputFile)

		cfg.DetectedFormat = formatStr
		cfg.TemplatePath = outputPath

		if err := writeFile(outputPath, []byte(result)); err != nil {
			return err
		}
	}

	return nil
}

func processTemplate(content string, cfg *Options, variant color.VariantMeta, accent string) string {
	data := []string{}

	data = append(data, cfg.Prefix+"id", variant.Id)
	data = append(data, cfg.Prefix+"name", variant.Name)
	data = append(data, cfg.Prefix+"type", variant.Appearance)
	data = append(data, cfg.Prefix+"appearance", variant.Appearance)
	data = append(data, cfg.Prefix+"description", variant.Description)

	if accent != "" {
		data = append(data, cfg.Prefix+"accentname", accent)

		if c, ok := variant.Colors[accent]; ok {
			accentColor := color.FormatColor(c, color.ColorFormat(cfg.Format), cfg.Plain, cfg.Commas, cfg.Spaces)
			data = append(data, cfg.Prefix+"accent", accentColor)
			if c.On != "" {
				if oc, ok := variant.Colors[c.On]; ok {
					onAccent := color.FormatColor(oc, color.ColorFormat(cfg.Format), cfg.Plain, cfg.Commas, cfg.Spaces)
					data = append(data, cfg.Prefix+"onaccent", onAccent)
				}
			}
		}
	}

	for name, c := range variant.Colors {
		varName := cfg.Prefix + name

		alphaRe := regexp.MustCompile(regexp.QuoteMeta(varName) + `/(\d+)`)
		matches := alphaRe.FindAllStringSubmatch(content, -1)
		seen := make(map[string]bool)
		for _, m := range matches {
			if seen[m[0]] {
				continue
			}
			seen[m[0]] = true
			alpha, _ := strconv.ParseFloat(m[1], 64)
			tmp := *c
			a := alpha / 100
			tmp.Alpha = &a
			data = append(data, m[0], color.FormatColor(&tmp, color.ColorFormat(cfg.Format), cfg.Plain, cfg.Commas, cfg.Spaces))
		}

		data = append(data, varName, color.FormatColor(c, color.ColorFormat(cfg.Format), cfg.Plain, cfg.Commas, cfg.Spaces))
	}

	result := strings.NewReplacer(data...).Replace(content)

	result = variantValueRe.ReplaceAllStringFunc(result, func(match string) string {
		parts := variantValueRe.FindStringSubmatch(match)
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
		}
		return match
	})

	return result
}

func templateFiles(path string) ([]string, error) {
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

func writeFile(outputPath string, content []byte) error {
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return err
	}
	return os.WriteFile(outputPath, content, 0644)
}

func buildOutputPath(cfg *Options, templatePath string, variant color.VariantMeta, accent string) string {
	ext := filepath.Ext(templatePath)

	if info, err := os.Stat(cfg.Template); err == nil && info.IsDir() {
		rel, _ := strings.CutPrefix(filepath.Dir(templatePath), filepath.Clean(cfg.Template))
		if accent != "" {
			return filepath.Join(cfg.Output, accent, variant.Id, rel, filepath.Base(templatePath))
		}
		return filepath.Join(cfg.Output, variant.Id, rel, filepath.Base(templatePath))
	}

	if accent != "" {
		return filepath.Join(cfg.Output, variant.Id, variant.Id+"-"+accent+ext)
	}
	return filepath.Join(cfg.Output, variant.Id+ext)
}
