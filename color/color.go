package color

type ColorFormat string

const (
	FormatHex ColorFormat = "hex"

	FormatHSL      ColorFormat = "hsl"
	FormatHSLCSS   ColorFormat = "hsl-css"
	FormatHSLArray ColorFormat = "hsl-array"

	FormatRGB      ColorFormat = "rgb"
	FormatRGBCSS   ColorFormat = "rgb-css"
	FormatRGBArray ColorFormat = "rgb-array"

	FormatAnsi ColorFormat = "ansi"
)
