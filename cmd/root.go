package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var version = "dev"

var rootCmd = &cobra.Command{
	Use:   "bloom",
	Short: "Generate Rosé Pine theme files from templates",
}

func Execute() {
	rootCmd.Version = version

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
