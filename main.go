package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var noSpaces bool

func main() {
	cfg := &Config{}
	flag.StringVar(&cfg.Template, "t", "template.json", "Path to template file")
	flag.StringVar(&cfg.Template, "template", "template.json", "Path to template file")
	flag.StringVar(&cfg.Output, "o", "dist", "Directory for generated files")
	flag.StringVar(&cfg.Output, "output", "dist", "Directory for generated files")
	flag.StringVar(&cfg.Prefix, "p", "$", "Variable prefix")
	flag.StringVar(&cfg.Prefix, "prefix", "$", "Variable prefix")
	flag.StringVar(&cfg.Format, "f", "hex", "Color output format")
	flag.StringVar(&cfg.Format, "format", "hex", "Color output format")
	flag.BoolVar(&noSpaces, "no-spaces", false, "Remove spaces from color values")
	flag.BoolVar(&cfg.Accents, "a", false, "Generate accent files")
	flag.BoolVar(&cfg.Accents, "accents", false, "Generate accent files")

	help := flag.Bool("help", false, "Show help")
	flag.Bool("h", false, "Show help")

	flag.Parse()

	cfg.Spaces = !noSpaces

	if *help {
		fmt.Println("🌱 Bloom - The Rosé Pine theme generator")
		fmt.Println("\nUsage:")
		flag.PrintDefaults()
		os.Exit(0)
	}

	if err := Build(cfg); err != nil {
		log.Fatal(err)
	}
}
