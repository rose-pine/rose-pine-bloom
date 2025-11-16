package cmd

import (
	"fmt"

	"github.com/rose-pine/rose-pine-bloom/docs"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Project initialization",
	Long:  `Creates files to start a new Rosé Pine theme project.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := docs.EnsureReadme(); err != nil {
			fmt.Println("unable to update README:", err)
		}
		if err := docs.EnsureLicense(); err != nil {
			fmt.Println("unable to update LICENSE:", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
