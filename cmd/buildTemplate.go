package cmd

import (
	"fmt"
	"os"

	"github.com/rose-pine/rose-pine-bloom/builder"
	"github.com/rose-pine/rose-pine-bloom/config"
	"github.com/spf13/cobra"
)

// buildTemplateCmd represents the buildTemplate command
var buildTemplateCmd = &cobra.Command{
	Use:   "build-template [flags] <input>",
	Short: "Builds the template file from a given theme file",
	Long:  `Builds the template file from a given theme file.`,
	Args:  cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validateFormat(format)
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := builder.BuildTemplate(&config.BuildTemplateConfig{
			Input:  args[0],
			Output: outputDir,
			Prefix: prefix,
			Format: format,
			Plain:  plain,
			Commas: !noCommas,
			Spaces: !noSpaces,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error building theme: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {

	buildTemplateCmd.Flags().StringVarP(&outputDir, "output", "o", "dist", "Directory for generated files")
	buildTemplateCmd.Flags().StringVarP(&prefix, "prefix", "p", "$", "Color variable prefix")
	buildTemplateCmd.Flags().StringVarP(&format, "format", "f", "hex", formatFlagUsage())
	buildTemplateCmd.Flags().BoolVar(&plain, "plain", false, "Remove decorators from color values")
	buildTemplateCmd.Flags().BoolVar(&noCommas, "no-commas", false, "Remove commas from color values")
	buildTemplateCmd.Flags().BoolVar(&noSpaces, "no-spaces", false, "Remove spaces from color values")
	rootCmd.AddCommand(buildTemplateCmd)
}
