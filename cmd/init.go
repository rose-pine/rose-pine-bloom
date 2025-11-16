package cmd

import (
	"fmt"
	"os"

	"github.com/rose-pine/rose-pine-bloom/builder"
	"github.com/rose-pine/rose-pine-bloom/config"
	"github.com/rose-pine/rose-pine-bloom/docs"
	"github.com/spf13/cobra"
)

var (
	variant string
	output  string
)

var initCmd = &cobra.Command{
	Use:   "init [theme-file]",
	Short: "Initialise new theme",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Creating theme files...")

		if err := docs.EnsureReadme(); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating README: %v\n", err)
		} else {
			fmt.Println("Created README.md")
		}

		if err := docs.EnsureLicense(); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating LICENSE: %v\n", err)
		} else {
			fmt.Println("Created LICENSE")
		}

		if len(args) > 0 {
			themeFile := args[0]
			fmt.Printf("Creating template from %s...\n", themeFile)

			err := builder.BuildTemplate(&config.BuildTemplateConfig{
				Input:   themeFile,
				Output:  output,
				Variant: variant,
				Prefix:  "$",
				Format:  "hex",
				Plain:   false,
				Commas:  true,
				Spaces:  true,
			})
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating template: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("Template created in %s\n", output)
		}

		fmt.Println("Theme initialised")
	},
}

func init() {
	initCmd.Flags().StringVarP(&variant, "variant", "v", "main", "Theme variant (main, moon, dawn)")
	initCmd.Flags().StringVarP(&output, "output", "o", ".", "Template output directory")
	rootCmd.AddCommand(initCmd)
}
