package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var (
	cfg      = &Config{}
	noCommas bool
	noSpaces bool
	showHelp bool
)

func detectTemplate(args []string) (string, error) {
	switch len(args) {
	case 0:
		files, err := os.ReadDir(".")
		if err != nil {
			return "", fmt.Errorf("failed to read current directory: %w", err)
		}

		for _, file := range files {
			if file.IsDir() {
				continue
			}
			name := file.Name()
			base := name[:len(name)-len(filepath.Ext(name))]
			if base == "template" {
				return name, nil
			}
		}
		return "", fmt.Errorf("unable to find template")

	case 1:
		return args[0], nil

	default:
		return "", fmt.Errorf("multiple positional arguments detected, ensure all flags come before the template")
	}
}

func printHelp() {
	helpMessage := fmt.Sprintf(`
  🌱 Bloom - The Rosé Pine theme generator

  Usage
    $ %s [options] <template>

  Options
    -o, --output <path>     Directory for generated files (default: dist)
    -p, --prefix <string>   Color variable prefix (default: $)
    -f, --format <format>   Color output format (default: hex)
    -c, --create <variant>  Create template from existing theme (default: main)
                            Variants: main, moon, dawn

    --accents               Create themes for each accent color
    --no-commas             Remove commas from color values
    --no-spaces             Remove spaces from color values

    -h, --help              Show help

  Formats
    hex           #c4a7e7
    hex-ns        c4a7e7

    hsl           267, 57%%, 78%%
    hsl-array     [267, 57%%, 78%%]
    hsl-function  hsl(267, 57%%, 78%%)

    rgb           196, 167, 231
    rgb-ansi      196;167;231
    rgb-array     [196, 167, 231]
    rgb-function  rgb(196, 167, 231)

  Examples
    $ rose-pine-bloom template.yaml
    $ rose-pine-bloom -f hsl -o dist template.json
    $ rose-pine-bloom -c dawn my-theme.toml

`, os.Args[0])
	fmt.Fprint(os.Stderr, helpMessage)
}

func main() {
	flag.StringVar(&cfg.Output, "o", "dist", "")
	flag.StringVar(&cfg.Output, "output", "dist", "")

	flag.StringVar(&cfg.Prefix, "p", "$", "")
	flag.StringVar(&cfg.Prefix, "prefix", "$", "")

	flag.StringVar(&cfg.Format, "f", "hex", "")
	flag.StringVar(&cfg.Format, "format", "hex", "")

	flag.BoolVar(&cfg.Accents, "a", false, "")
	flag.BoolVar(&cfg.Accents, "accents", false, "")

	flag.BoolVar(&noCommas, "no-commas", false, "")
	flag.BoolVar(&noSpaces, "no-spaces", false, "")

	flag.BoolVar(&showHelp, "h", false, "")
	flag.BoolVar(&showHelp, "help", false, "")

	flag.Parse()

	if showHelp {
		printHelp()
		return
	}

	args := flag.Args()

	template, err := detectTemplate(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	cfg.Template = template
	cfg.Commas = !noCommas
	cfg.Spaces = !noSpaces

	if err := Build(cfg); err != nil {
		log.Fatal(err)
	}
}
