package main

type Config struct {
	Template    string
	Output      string
	Prefix      string
	Format      string
	StripSpaces bool
}

type ColorFormat string

const (
	FormatHex      ColorFormat = "hex"
	FormatHexNS    ColorFormat = "hex-ns"
	FormatRGB      ColorFormat = "rgb"
	FormatRGBNS    ColorFormat = "rgb-ns"
	FormatRGBAnsi  ColorFormat = "rgb-ansi"
	FormatRGBArray ColorFormat = "rgb-array"
	FormatRGBFunc  ColorFormat = "rgb-function"
	FormatHSL      ColorFormat = "hsl"
	FormatHSLNS    ColorFormat = "hsl-ns"
	FormatHSLArray ColorFormat = "hsl-array"
	FormatHSLFunc  ColorFormat = "hsl-function"
)
