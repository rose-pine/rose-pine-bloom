package main

import (
	"fmt"
	"os"
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
	{Name: "hsl-function", Example: "hsl(2, 55%, 83%)"},
	{Name: "rgb", Example: "235, 188, 186"},
	{Name: "rgb-ansi", Example: "235;188;186"},
	{Name: "rgb-array", Example: "[235, 188, 186]"},
	{Name: "rgb-function", Example: "rgb(235, 188, 186)"},
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

func PrintFormatsTable() {
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

    --accents               Create themes for each accent color
    --no-commas             Remove commas from color values
    --no-spaces             Remove spaces from color values

    -h, --help              Show help

  Formats
%s
  Examples
    $ %[1]s template.yaml
    $ %[1]s -f hsl -o dist template.json
    $ %[1]s -c dawn my-theme.toml

`, os.Args[0], formatsTable())
}

func PrintHelp() {
	fmt.Fprint(os.Stderr, helpText())
}
