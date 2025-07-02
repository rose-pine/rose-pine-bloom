package main

import (
	"fmt"
	"strings"
)

type Color struct {
	Hex   string     `json:"hex"`
	HSL   [3]float64 `json:"hsl"`
	RGB   [3]int     `json:"rgb"`
	Alpha *float64   `json:"alpha,omitempty"`
	On    string     `json:"on,omitempty"`
}

func formatColor(c *Color, format ColorFormat, commas bool, spaces bool) string {
	workingString := ""
	formatAlpha := func(alpha float64) string {
		s := fmt.Sprintf("%.2f", alpha)
		s = strings.TrimRight(s, "0")
		s = strings.TrimRight(s, ".")
		return s
	}

	switch format {
	case FormatHex:
		workingString = "#" + c.Hex
	case FormatHexNS:
		workingString = c.Hex
	case FormatHSL:
		hsl := fmt.Sprintf("%v, %v%%, %v%%", c.HSL[0], c.HSL[1], c.HSL[2])
		if c.Alpha != nil {
			hsl += fmt.Sprintf(", %s", formatAlpha(*c.Alpha))
		}
		workingString = hsl
	case FormatHSLFunc:
		hsl := fmt.Sprintf("%v, %v%%, %v%%", c.HSL[0], c.HSL[1], c.HSL[2])
		if c.Alpha != nil {
			hsl += fmt.Sprintf(", %s", formatAlpha(*c.Alpha))
		}
		prefix := "hsl"
		if c.Alpha != nil {
			prefix = "hsla"
		}
		workingString = fmt.Sprintf("%s(%s)", prefix, hsl)
	case FormatHSLCSS:
		// CSS format uses deg and spaces, not commas
		prefix := "hsl"
		hsl := fmt.Sprintf("%s(%vdeg %v%% %v%%", prefix, c.HSL[0], c.HSL[1], c.HSL[2])
		if c.Alpha != nil {
			hsl += fmt.Sprintf(" / %s", formatAlpha(*c.Alpha))
		}
		hsl += ")"
		workingString = hsl
	case FormatHSLArray:
		hsl := fmt.Sprintf("[%v, %v%%, %v%%", c.HSL[0], c.HSL[1], c.HSL[2])
		if c.Alpha != nil {
			hsl += fmt.Sprintf(", %s", formatAlpha(*c.Alpha))
		}
		hsl += "]"
		workingString = hsl
	case FormatRGB:
		rgb := fmt.Sprintf("%d, %d, %d", c.RGB[0], c.RGB[1], c.RGB[2])
		if c.Alpha != nil {
			rgb += fmt.Sprintf(", %s", formatAlpha(*c.Alpha))
		}
		workingString = rgb
	case FormatRGBArray:
		rgb := fmt.Sprintf("[%d, %d, %d", c.RGB[0], c.RGB[1], c.RGB[2])
		if c.Alpha != nil {
			rgb += fmt.Sprintf(", %s", formatAlpha(*c.Alpha))
		}
		rgb += "]"
		workingString = rgb
	case FormatRGBFunc:
		prefix := "rgb"
		if c.Alpha != nil {
			prefix = "rgba"
		}
		rgb := fmt.Sprintf("%s(%d, %d, %d", prefix, c.RGB[0], c.RGB[1], c.RGB[2])
		if c.Alpha != nil {
			rgb += fmt.Sprintf(", %s", formatAlpha(*c.Alpha))
		}
		rgb += ")"
		workingString = rgb
	case FormatRGBAnsi:
		if c.Alpha != nil {
			return fmt.Sprintf("%d;%d;%d;%s", c.RGB[0], c.RGB[1], c.RGB[2], formatAlpha(*c.Alpha))
		}
		return fmt.Sprintf("%d;%d;%d", c.RGB[0], c.RGB[1], c.RGB[2])
	default:
		workingString = "#" + c.Hex
	}

	if commas == false {
		workingString = strings.ReplaceAll(workingString, ",", "")
	}
	if spaces == false {
		workingString = strings.ReplaceAll(workingString, " ", "")
	}

	return workingString
}

type Variant struct {
	id          string
	name        string
	variantType string
	description string
	colors      map[string]*Color
}

const description = "All natural pine, faux fur and a bit of soho vibes for the classy minimalist"

var (
	MainVariant = Variant{
		id:          "rose-pine",
		name:        "Rosé Pine",
		variantType: "dark",
		description: description,
		colors: map[string]*Color{
			"base": {
				Hex: "191724",
				HSL: [3]float64{249, 22, 12},
				RGB: [3]int{25, 23, 36},
			},
			"surface": {
				Hex: "1f1d2e",
				HSL: [3]float64{247, 23, 15},
				RGB: [3]int{31, 29, 46},
			},
			"overlay": {
				Hex: "26233a",
				HSL: [3]float64{245, 25, 18},
				RGB: [3]int{38, 35, 58},
			},
			"muted": {
				Hex: "6e6a86",
				HSL: [3]float64{249, 12, 47},
				RGB: [3]int{110, 106, 134},
			},
			"subtle": {
				Hex: "908caa",
				HSL: [3]float64{248, 15, 61},
				RGB: [3]int{144, 140, 170},
			},
			"text": {
				Hex: "e0def4",
				HSL: [3]float64{245, 50, 91},
				RGB: [3]int{224, 222, 244},
			},
			"love": {
				Hex: "eb6f92",
				HSL: [3]float64{343, 76, 68},
				RGB: [3]int{235, 111, 146},
				On:  "text",
			},
			"gold": {
				Hex: "f6c177",
				HSL: [3]float64{35, 88, 72},
				RGB: [3]int{246, 193, 119},
				On:  "surface",
			},
			"rose": {
				Hex: "ebbcba",
				HSL: [3]float64{2, 55, 83},
				RGB: [3]int{235, 188, 186},
				On:  "surface",
			},
			"pine": {
				Hex: "31748f",
				HSL: [3]float64{197, 49, 38},
				RGB: [3]int{49, 116, 143},
				On:  "text",
			},
			"foam": {
				Hex: "9ccfd8",
				HSL: [3]float64{189, 43, 73},
				RGB: [3]int{156, 207, 216},
				On:  "surface",
			},
			"iris": {
				Hex: "c4a7e7",
				HSL: [3]float64{267, 57, 78},
				RGB: [3]int{196, 167, 231},
				On:  "surface",
			},
			"highlightLow": {
				Hex: "21202e",
				HSL: [3]float64{244, 18, 15},
				RGB: [3]int{33, 32, 46},
			},
			"highlightMed": {
				Hex: "403d52",
				HSL: [3]float64{247, 15, 28},
				RGB: [3]int{64, 61, 82},
			},
			"highlightHigh": {
				Hex: "524f67",
				HSL: [3]float64{245, 13, 36},
				RGB: [3]int{82, 79, 103},
			},
		},
	}

	MoonVariant = Variant{
		id:          "rose-pine-moon",
		name:        "Rosé Pine Moon",
		variantType: "dark",
		description: description,
		colors: map[string]*Color{
			"base": {
				Hex: "232136",
				HSL: [3]float64{246, 24, 17},
				RGB: [3]int{35, 33, 54},
			},
			"surface": {
				Hex: "2a273f",
				HSL: [3]float64{248, 24, 20},
				RGB: [3]int{42, 39, 63},
			},
			"overlay": {
				Hex: "393552",
				HSL: [3]float64{248, 21, 26},
				RGB: [3]int{57, 53, 82},
			},
			"muted": {
				Hex: "6e6a86",
				HSL: [3]float64{249, 12, 47},
				RGB: [3]int{110, 106, 134},
			},
			"subtle": {
				Hex: "908caa",
				HSL: [3]float64{248, 15, 61},
				RGB: [3]int{144, 140, 170},
			},
			"text": {
				Hex: "e0def4",
				HSL: [3]float64{245, 50, 91},
				RGB: [3]int{224, 222, 244},
			},
			"love": {
				Hex: "eb6f92",
				HSL: [3]float64{343, 76, 68},
				RGB: [3]int{235, 111, 146},
				On:  "text",
			},
			"gold": {
				Hex: "f6c177",
				HSL: [3]float64{35, 88, 72},
				RGB: [3]int{246, 193, 119},
				On:  "surface",
			},
			"rose": {
				Hex: "ea9a97",
				HSL: [3]float64{2, 66, 75},
				RGB: [3]int{234, 154, 151},
				On:  "surface",
			},
			"pine": {
				Hex: "3e8fb0",
				HSL: [3]float64{197, 48, 47},
				RGB: [3]int{62, 143, 176},
				On:  "text",
			},
			"foam": {
				Hex: "9ccfd8",
				HSL: [3]float64{189, 43, 73},
				RGB: [3]int{156, 207, 216},
				On:  "surface",
			},
			"iris": {
				Hex: "c4a7e7",
				HSL: [3]float64{267, 57, 78},
				RGB: [3]int{196, 167, 231},
				On:  "surface",
			},
			"highlightLow": {
				Hex: "2a283e",
				HSL: [3]float64{245, 22, 20},
				RGB: [3]int{42, 40, 62},
			},
			"highlightMed": {
				Hex: "44415a",
				HSL: [3]float64{247, 16, 30},
				RGB: [3]int{68, 65, 90},
			},
			"highlightHigh": {
				Hex: "56526e",
				HSL: [3]float64{249, 15, 38},
				RGB: [3]int{86, 82, 110},
			},
		},
	}

	DawnVariant = Variant{
		id:          "rose-pine-dawn",
		name:        "Rosé Pine Dawn",
		variantType: "light",
		description: description,
		colors: map[string]*Color{
			"base": {
				Hex: "faf4ed",
				HSL: [3]float64{32, 57, 95},
				RGB: [3]int{250, 244, 237},
			},
			"surface": {
				Hex: "fffaf3",
				HSL: [3]float64{35, 100, 98},
				RGB: [3]int{255, 250, 243},
			},
			"overlay": {
				Hex: "f2e9e1",
				HSL: [3]float64{25, 36, 92},
				RGB: [3]int{242, 233, 225},
			},
			"muted": {
				Hex: "9893a5",
				HSL: [3]float64{254, 9, 61},
				RGB: [3]int{152, 147, 165},
			},
			"subtle": {
				Hex: "797593",
				HSL: [3]float64{249, 13, 52},
				RGB: [3]int{121, 117, 147},
			},
			"text": {
				Hex: "575279",
				HSL: [3]float64{248, 19, 40},
				RGB: [3]int{87, 82, 121},
			},
			"love": {
				Hex: "b4637a",
				HSL: [3]float64{343, 35, 55},
				RGB: [3]int{180, 99, 122},
				On:  "surface",
			},
			"gold": {
				Hex: "ea9d34",
				HSL: [3]float64{35, 81, 56},
				RGB: [3]int{234, 157, 52},
				On:  "surface",
			},
			"rose": {
				Hex: "d7827e",
				HSL: [3]float64{2, 55, 67},
				RGB: [3]int{215, 130, 126},
				On:  "surface",
			},
			"pine": {
				Hex: "286983",
				HSL: [3]float64{197, 53, 34},
				RGB: [3]int{40, 105, 131},
				On:  "surface",
			},
			"foam": {
				Hex: "56949f",
				HSL: [3]float64{189, 30, 48},
				RGB: [3]int{86, 148, 159},
				On:  "surface",
			},
			"iris": {
				Hex: "907aa9",
				HSL: [3]float64{267, 22, 57},
				RGB: [3]int{144, 122, 169},
				On:  "surface",
			},
			"highlightLow": {
				Hex: "f4ede8",
				HSL: [3]float64{25, 35, 93},
				RGB: [3]int{244, 237, 232},
			},
			"highlightMed": {
				Hex: "dfdad9",
				HSL: [3]float64{10, 9, 86},
				RGB: [3]int{223, 218, 217},
			},
			"highlightHigh": {
				Hex: "cecacd",
				HSL: [3]float64{315, 4, 80},
				RGB: [3]int{206, 202, 205},
			},
		},
	}
)

var variants = []Variant{
	MainVariant,
	MoonVariant,
	DawnVariant,
}

var accents = []string{
	"love", "gold", "rose", "pine", "foam", "iris",
}
