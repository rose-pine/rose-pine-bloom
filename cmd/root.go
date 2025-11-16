package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var version = "dev"

var RootCmd = &cobra.Command{
	Use:   "bloom",
	Short: "Generate Rosé Pine theme files from templates",
}

func Execute() {
	RootCmd.Version = version

	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
