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

var AllFormats = []string{
	string(FormatHex),
	string(FormatHSL),
	string(FormatHSLCSS),
	string(FormatHSLArray),
	string(FormatRGB),
	string(FormatRGBCSS),
	string(FormatRGBArray),
	string(FormatAnsi),
}
