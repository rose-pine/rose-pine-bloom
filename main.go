package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rose-pine/rose-pine-bloom/builder"
	"github.com/rose-pine/rose-pine-bloom/color"
	"github.com/rose-pine/rose-pine-bloom/config"
	"github.com/rose-pine/rose-pine-bloom/docs"
)

func helpText() string {
	return fmt.Sprintf(`
  🌱 Bloom %s - The Rosé Pine theme generator

  Usage
    $ bloom [options] <template>

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
    -v, --version           Show bloom version

  Formats
%s
  Examples
    $ bloom template.yaml
    $ bloom --format hsl --output ./themes template.json
    $ bloom --create dawn my-theme.toml

`, getCurrentVersion(), color.FormatsTable())
}

func printHelp() {
	fmt.Fprint(os.Stderr, helpText())
}

func main() {
	cfg := &config.Config{}

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

	showVersion := false
	flag.BoolVar(&showVersion, "v", false, "")
	flag.BoolVar(&showVersion, "version", false, "")

	flag.Usage = printHelp
	flag.Parse()

	if showVersion {
		fmt.Println("bloom ", getCurrentVersion())
		os.Exit(0)
	}

	args := flag.Args()

	template, err := findTemplate(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
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

	if err := builder.Build(cfg); err != nil {
		log.Fatal(err)
	}

	buildCmd := "bloom " + strings.Join(os.Args[1:], " ")
	if err := docs.EnsureReadmeWithBuildCommand(buildCmd, getCurrentVersion()); err != nil {
		fmt.Println("unable to update README:", err)
	}

	if err := docs.EnsureLicense(); err != nil {
		fmt.Println("unable to update LICENSE:", err)
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

func getCurrentVersion() string {
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	version := strings.TrimSpace(string(output))
	return version
}
