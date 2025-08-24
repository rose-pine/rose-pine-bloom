package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// docCmd represents the doc command
var docCmd = &cobra.Command{
	Use:   "doc <output_directory>",
	Short: "Generates documentation",
	Long:  `The doc command generates documentation for the CLI application.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		format, _ := cmd.Flags().GetString("format")
		outputDir := args[0]

		switch format {
		case "markdown", "md":
			err := doc.GenMarkdownTree(rootCmd, outputDir)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error generating markdown documentation: %v\n", err)
				os.Exit(1)
			}
		case "man":
			header := &doc.GenManHeader{
				Title:   "CLI",
				Section: "1",
			}
			err := doc.GenManTree(rootCmd, header, outputDir)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error generating man documentation: %v\n", err)
				os.Exit(1)
			}
		default:
			fmt.Fprintf(os.Stderr, "Unsupported format: %s. Use 'markdown' or 'man'\n", format)
			os.Exit(1)
		}
	},
}

func init() {
	docCmd.Flags().StringP("format", "f", "markdown", "Documentation format (markdown, md, man)")
	rootCmd.AddCommand(docCmd)
}
