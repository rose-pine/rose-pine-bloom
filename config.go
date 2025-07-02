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
	FormatHex ColorFormat = "hex"

	FormatHSL      ColorFormat = "hsl"
	FormatHSLArray ColorFormat = "hsl-array"
	FormatHSLCSS   ColorFormat = "hsl-css"
	FormatHSLFunc  ColorFormat = "hsl-function"

	FormatRGB      ColorFormat = "rgb"
	FormatRGBArray ColorFormat = "rgb-array"
	FormatRGBCSS   ColorFormat = "rgb-css"
	FormatRGBFunc  ColorFormat = "rgb-function"

	FormatAnsi ColorFormat = "ansi"
)
