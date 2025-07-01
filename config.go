package main

type Config struct {
	Template string
	Output   string
	Prefix   string
	Format   string
	Create   string
	Plain    bool
	Commas   bool
	Spaces   bool
}

type ColorFormat string

const (
	FormatHex         ColorFormat = "hex"
	FormatHSL         ColorFormat = "hsl"
	FormatHSLFunction ColorFormat = "hsl-function"
	FormatHSLCSS      ColorFormat = "hsl-css"
	FormatHSLArray    ColorFormat = "hsl-array"
	FormatRGB         ColorFormat = "rgb"
	FormatRGBFunction ColorFormat = "rgb-function"
	FormatRGBCSS      ColorFormat = "rgb-css"
	FormatRGBArray    ColorFormat = "rgb-array"
	FormatANSI        ColorFormat = "ansi"
)
