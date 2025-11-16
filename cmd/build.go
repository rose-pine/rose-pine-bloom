package cmd

import (
	"fmt"
	"os"

	"github.com/rose-pine/rose-pine-bloom/builder"
	"github.com/rose-pine/rose-pine-bloom/config"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build [flags] <template>",
	Short: "Builds the theme files from templates",
	Long: `Builds the theme files from templates.

	This command processes the theme templates and generates the final theme files that can be used in applications.
	The template argument specifies the template file or directory to process.
	If the template argument is a directory, it will process all templates in that directory.`,
	Args: cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validateFormat(format)
	},
	Run: func(cmd *cobra.Command, args []string) {
		template := args[0]

		err := builder.Build(&config.BuildConfig{
			Template: template,
			Output:   outputDir,
			Prefix:   prefix,
			Format:   format,
			Plain:    plain,
			Commas:   !noCommas,
			Spaces:   !noSpaces,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error building theme: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {

	buildCmd.Flags().StringVarP(&outputDir, "output", "o", "dist", "Directory for generated files")
	buildCmd.Flags().StringVarP(&prefix, "prefix", "p", "$", "Color variable prefix")
	buildCmd.Flags().StringVarP(&format, "format", "f", "hex", formatFlagUsage())
	buildCmd.Flags().BoolVar(&plain, "plain", false, "Remove decorators from color values")
	buildCmd.Flags().BoolVar(&noCommas, "no-commas", false, "Remove commas from color values")
	buildCmd.Flags().BoolVar(&noSpaces, "no-spaces", false, "Remove spaces from color values")

	rootCmd.AddCommand(buildCmd)

}
