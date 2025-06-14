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
	FormatRGB      ColorFormat = "rgb"
	FormatRGBAnsi  ColorFormat = "rgb-ansi"
	FormatRGBArray ColorFormat = "rgb-array"
	FormatRGBFunc  ColorFormat = "rgb-function"
	FormatHSL      ColorFormat = "hsl"
	FormatHSLArray ColorFormat = "hsl-array"
	FormatHSLFunc  ColorFormat = "hsl-function"
)
