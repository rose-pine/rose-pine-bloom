package main

type Config struct {
	Template string
	Output   string
	Prefix   string
	Format   string
	Create   string
	Accents  bool
	Commas   bool
	Spaces   bool
}

type ColorFormat string

const (
	FormatHex      ColorFormat = "hex"
	FormatHexNS    ColorFormat = "hex-ns"
	FormatHSL      ColorFormat = "hsl"
	FormatHSLArray ColorFormat = "hsl-array"
	FormatHSLCSS   ColorFormat = "hsl-css"
	FormatHSLFunc  ColorFormat = "hsl-function"
	FormatRGB      ColorFormat = "rgb"
	FormatRGBAnsi  ColorFormat = "rgb-ansi"
	FormatRGBArray ColorFormat = "rgb-array"
	FormatRGBFunc  ColorFormat = "rgb-function"
)
