package cmd

import (
	"os"

	"github.com/rose-pine/rose-pine-bloom/version"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bloom",
	Short: "The Rosé Pine theme generator",
	Long:  `The Rosé Pine theme generator.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.Version = version.GetCurrentVersion()

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.rose-pine-bloom.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
}
