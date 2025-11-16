package cmd

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/rose-pine/rose-pine-bloom/builder"
	"github.com/rose-pine/rose-pine-bloom/color"
	"github.com/spf13/cobra"
)

var (
	outputDir string
	prefix    string
	format    string
	plain     bool
	noCommas  bool
	noSpaces  bool
)

var buildCmd = &cobra.Command{
	Use:   "build <template>",
	Short: "Generate theme files from template",
	Args:  cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validateFormat(format)
	},
	Run: func(cmd *cobra.Command, args []string) {
		template := args[0]

		fmt.Printf("Building themes from %s...\n", template)

		err := builder.Build(&builder.Options{
			Template: template,
			Output:   outputDir,
			Prefix:   prefix,
			Format:   format,
			Plain:    plain,
			Commas:   !noCommas,
			Spaces:   !noSpaces,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error building themes: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Themes generated in %s\n", outputDir)
	},
}

func init() {
	buildCmd.Flags().StringVarP(&outputDir, "output", "o", "dist", "Output directory")
	buildCmd.Flags().StringVarP(&prefix, "prefix", "p", "$", "Variable prefix")
	buildCmd.Flags().StringVarP(&format, "format", "f", "hex", formatFlagUsage())
	buildCmd.Flags().BoolVar(&plain, "plain", false, "Remove decorators")
	buildCmd.Flags().BoolVar(&noCommas, "no-commas", false, "Remove commas")
	buildCmd.Flags().BoolVar(&noSpaces, "no-spaces", false, "Remove spaces")

	rootCmd.AddCommand(buildCmd)
}

func formatFlagUsage() string {
	table, err := color.FormatsTable()
	if err != nil {
		fmt.Printf("Error generating format table: %v", err)
		os.Exit(1)
	}
	return fmt.Sprintf("Color format:\n%s", table)
}

func validateFormat(format string) error {
	if slices.Contains(color.AllFormats, format) {
		return nil
	}
	return fmt.Errorf("invalid format '%s'. Valid formats: %s", format, strings.Join(color.AllFormats, ", "))
}
