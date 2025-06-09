package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	cfg := &Config{}
	flag.StringVar(&cfg.Output, "o", "dist", "Directory for generated files")
	flag.StringVar(&cfg.Output, "output", "dist", "Directory for generated files")
	flag.StringVar(&cfg.Prefix, "p", "$", "Variable prefix")
	flag.StringVar(&cfg.Prefix, "prefix", "$", "Variable prefix")
	flag.StringVar(&cfg.Format, "f", "hex", "Color output format")
	flag.StringVar(&cfg.Format, "format", "hex", "Color output format")
	flag.BoolVar(&cfg.StripSpaces, "s", false, "Strip spaces in output")
	flag.BoolVar(&cfg.StripSpaces, "strip-spaces", false, "Strip spaces in output")
	flag.BoolVar(&cfg.Accents, "a", false, "Generate accent files")
	flag.BoolVar(&cfg.Accents, "accents", false, "Generate accent files")

	help := flag.Bool("help", false, "Show help")
	help = flag.Bool("h", false, "Show help")

	flag.Parse()

	if *help {
		fmt.Println("🌱 Bloom - The Rosé Pine theme generator")
		fmt.Println("\nUsage:")
		fmt.Println("  Positional Arguments: (note: they must be specified last)")
		fmt.Println("\t[File] - the file to use")
		fmt.Println()
		flag.PrintDefaults()
		os.Exit(0)
	}

	if flag.NArg() != 1 {
		fmt.Println("Required 1 positional argument \"file\", Got", flag.NArg())
		os.Exit(1)
	} else {
		cfg.File = flag.Args()[0]
	}

	if err := Build(cfg); err != nil {
		log.Fatal(err)
	}
}
