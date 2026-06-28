package derive

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rose-pine/rose-pine-bloom/builder"
	"github.com/rose-pine/rose-pine-bloom/color"
)

type DeriveOpts struct {
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

func DeriveTemplate(cfg *DeriveOpts) error {
	if err := os.MkdirAll(cfg.Output, 0755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	return createTemplates(cfg)
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

func createTemplates(opts *DeriveOpts) error {
	files, err := builder.DiscoverTemplates(opts.Input)
	if err != nil {
		return err
	}

	variant := color.MainVariantMeta
	switch opts.Variant {
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

		formatStr, plain, commas, spaces := opts.Format, opts.Plain, opts.Commas, opts.Spaces
		if formatStr == "" {
			var df color.ColorFormat
			df, plain, commas, spaces = detectFormatOptions(content, variant)
			formatStr = string(df)
		}

		data := []string{}

		for name, c := range variant.Colors {
			val := color.FormatColor(c, color.ColorFormat(formatStr), plain, commas, spaces)
			data = append(data, val, opts.Prefix+name)
		}

		data = append(data, variant.Id, opts.Prefix+"id")
		data = append(data, variant.Name, opts.Prefix+"name")
		data = append(data, variant.Description, opts.Prefix+"description")

		result := strings.NewReplacer(data...).Replace(content)

		if content == result {
			fmt.Printf("%sNo matches for format %q. Available formats:\n  %s%s\n", warnColor, formatStr, strings.Join(color.AllFormats, ", "), resetColor)
		}

		ext := filepath.Ext(file)
		outputFile := "template" + ext
		outputPath := filepath.Join(opts.Output, outputFile)

		opts.DetectedFormat = formatStr
		opts.TemplatePath = outputPath

		if err := writeFile(outputPath, []byte(result)); err != nil {
			return err
		}
	}

	return nil
}

func writeFile(outputPath string, content []byte) error {
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return err
	}
	return os.WriteFile(outputPath, content, 0644)
}
