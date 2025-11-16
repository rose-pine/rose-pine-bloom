//go:build ignore

package main

import (
	"log"

	"github.com/rose-pine/rose-pine-bloom/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	err := doc.GenMarkdownTree(cmd.RootCmd, "./docs")
	if err != nil {
		log.Fatal(err)
	}
}
