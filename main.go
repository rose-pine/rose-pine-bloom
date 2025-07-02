package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"
)

type format struct {
	Name    string
	Example string
}

var formats = [...]format{
	{Name: "hex", Example: "#ebbcba"},
	{Name: "hex-ns", Example: "ebbcba"},

	{Name: "hsl", Example: "2, 55%, 83%"},
	{Name: "hsl-array", Example: "[2, 55%, 83%]"},
	{Name: "hsl-css", Example: "hsl(2deg 55% 83%)"},
	{Name: "hsl-function", Example: "hsl(2, 55%, 83%)"},

	{Name: "rgb", Example: "235, 188, 186"},
	{Name: "rgb-array", Example: "[235, 188, 186]"},
	{Name: "rgb-css", Example: "rgb(235 188 186)"},
	{Name: "rgb-function", Example: "rgb(235, 188, 186)"},

	{Name: "ansi", Example: "235;188;186"},
}

func formatsTable() string {
	var sb strings.Builder
	w := tabwriter.NewWriter(&sb, 1, 1, 1, ' ', 0)
	for _, f := range formats {
		fmt.Fprintf(w, "    %-13s %s\n", f.Name, f.Example)
	}
	w.Flush()
	return sb.String()
}

func printFormatsTable() {
	fmt.Fprint(os.Stdout, formatsTable())
}

func helpText() string {
	return fmt.Sprintf(`
  🌱 Bloom - The Rosé Pine theme generator

  Usage
    $ %s [options] <template>

  Options
    -o, --output <path>     Directory for generated files (default: dist)
    -p, --prefix <string>   Color variable prefix (default: $)
    -f, --format <format>   Color output format (default: hex)
    -c, --create <variant>  Create template from existing theme (default: main)
                            Variants: main, moon, dawn

    --plain                 Remove decorators from color values
    --no-commas             Remove commas from color values
    --no-spaces             Remove spaces from color values

    -h, --help              Show help

  Formats
%s
  Examples
    $ %[1]s template.yaml
    $ %[1]s --format hsl --output ./themes template.json
    $ %[1]s --create dawn my-theme.toml

`, os.Args[0], formatsTable())
}

func printHelp() {
	fmt.Fprint(os.Stderr, helpText())
}

func main() {
	cfg := &Config{}

	flag.StringVar(&cfg.Output, "o", "dist", "")
	flag.StringVar(&cfg.Output, "output", "dist", "")

	flag.StringVar(&cfg.Create, "c", "", "")
	flag.StringVar(&cfg.Create, "create", "", "")

	flag.StringVar(&cfg.Prefix, "p", "$", "")
	flag.StringVar(&cfg.Prefix, "prefix", "$", "")

	flag.StringVar(&cfg.Format, "f", "hex", "")
	flag.StringVar(&cfg.Format, "format", "hex", "")

	flag.BoolVar(&cfg.Plain, "plain", false, "")

	noCommas := flag.Bool("no-commas", false, "")
	noSpaces := flag.Bool("no-spaces", false, "")

	flag.Usage = printHelp
	flag.Parse()

	args := flag.Args()

	template, err := findTemplate(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		printHelp()
		os.Exit(1)
	}

	outputPassed := wasFlagPassed("o") || wasFlagPassed("output")
	createPassed := wasFlagPassed("c") || wasFlagPassed("create")

	if !outputPassed && createPassed {
		cfg.Output = "."
	}

	cfg.Template = template
	cfg.Commas = !*noCommas
	cfg.Spaces = !*noSpaces

	// Backward compatibility: hex-ns is equivalent to hex --plain
	if cfg.Format == "hex-ns" {
		cfg.Format = "hex"
		cfg.Plain = true
	}

	// Backward compatibility: rgb-ansi is equivalent to ansi
	if cfg.Format == "rgb-ansi" {
		cfg.Format = "ansi"
	}

	if err := Build(cfg); err != nil {
		log.Fatal(err)
	}
}

func wasFlagPassed(name string) bool {
	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		if arg == "-"+name || arg == "--"+name {
			return true
		}
		if strings.Contains(arg, "=") && (strings.HasPrefix(arg, "-"+name+"=") || strings.HasPrefix(arg, "--"+name+"=")) {
			return true
		}
	}
	return false
}

func findTemplate(args []string) (string, error) {
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
		return "", fmt.Errorf("unable to find template file")
	case 1:
		return args[0], nil
	default:
		return "", fmt.Errorf("multiple positional arguments detected, ensure all flags come before the template")
	}
}
