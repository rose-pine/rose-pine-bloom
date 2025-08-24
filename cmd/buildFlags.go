package cmd

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/rose-pine/rose-pine-bloom/color"
)

var (
	outputDir string
	prefix    string
	format    string
	plain     bool
	noCommas  bool
	noSpaces  bool
)

func formatFlagUsage() string {
	table, err := color.FormatsTable()
	if err != nil {
		fmt.Printf("Error generating format table: %v", err)
		os.Exit(1)
	}
	return fmt.Sprintf("Color output format:\n%s", table)
}

func validateFormat(format string) error {
	if slices.Contains(color.AllFormats, format) {
		return nil
	}
	return fmt.Errorf("invalid format '%s'. Valid formats are: %s", format, strings.Join(color.AllFormats, ", "))
}
