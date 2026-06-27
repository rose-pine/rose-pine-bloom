package cmd

import (
	"fmt"
	"os"
	"slices"

	"github.com/rose-pine/rose-pine-bloom/builder"
	"github.com/rose-pine/rose-pine-bloom/color"
	"github.com/spf13/cobra"
)

var (
	outDir   string
	prefix   string
	format   string
	plain    bool
	noCommas bool
	noSpaces bool
)

var buildCmd = &cobra.Command{
	Use:   "build <template>",
	Short: "Generate theme files from template",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		template := args[0]

		if !slices.Contains(color.AllFormats, format) {
			fmt.Fprintf(os.Stderr, "invalid format %q\n", format)
			os.Exit(1)
		}

		if len(prefix) != 1 {
			fmt.Fprintf(os.Stderr, "invalid prefix, must be exactly one character long\n")
			os.Exit(1)
		}

		fmt.Printf("Building themes from %s...\n", template)

		opts := builder.BuildOpts{
			Prefix:        rune(prefix[0]),
			DefaultFormat: color.ColorFormat(format),
			Plain:         plain,
			Commas:        !noCommas,
			Spaces:        !noSpaces,
		}
		err := builder.Build(template, outDir, &opts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error building themes: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Themes generated in %s\n", outDir)

		cmdLine := "bloom build " + template
		cmdLine += " --output " + outDir
		cmdLine += " --prefix " + string(prefix)
		cmdLine += " --format " + format
		if plain {
			cmdLine += " --plain"
		}
		if noCommas {
			cmdLine += " --no-commas"
		}
		if noSpaces {
			cmdLine += " --no-spaces"
		}

		if err := updateReadme(readmeSection(cmdLine)); err != nil {
			fmt.Fprintf(os.Stderr, "Error updating README: %v\n", err)
		} else {
			fmt.Println("Updated README.md")
		}
	},
}

func init() {
	buildCmd.Flags().StringVarP(&outDir, "output", "o", "dist", "output directory")
	buildCmd.Flags().StringVarP(&prefix, "prefix", "p", "$", "variable prefix")
	buildCmd.Flags().StringVarP(&format, "format", "f", "hex", "hex, hsl, hsl-css, hsl-array, rgb, rgb-css, rgb-array, ansi")
	buildCmd.Flags().BoolVar(&plain, "plain", false, "strip wrappers (#, rgb(), hsl(), brackets) from output")
	buildCmd.Flags().BoolVar(&noCommas, "no-commas", false, "remove commas")
	buildCmd.Flags().BoolVar(&noSpaces, "no-spaces", false, "remove spaces")

	rootCmd.AddCommand(buildCmd)
}
