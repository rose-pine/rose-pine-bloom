package main

import (
	"strconv"
	"strings"
)

type RGB struct {
	R, G, B uint8
}

type HSL struct {
	H    uint16
	S, L uint8
}

func rgb(r uint8, g uint8, b uint8) RGB {
	return RGB{r, g, b}
}

func hsl(h uint16, s uint8, l uint8) HSL {
	return HSL{h, s, l}
}

type Color struct {
	HSL   HSL      `json:"hsl"`
	RGB   RGB      `json:"rgb"`
	Alpha *float64 `json:"alpha,omitempty"`
	On    string   `json:"on,omitempty"`
}

const hex = "0123456789abcdef"

func hexComponent(c uint8) (byte, byte) {
	return hex[c>>4], hex[c&0x0f]
}

func formatAlpha(alpha float64) string {
	return strconv.FormatFloat(alpha, 'f', -1, 64)
}

func formatUint[T ~uint8 | ~uint16](n T) string {
	return strconv.FormatUint(uint64(n), 10)
}

func formatColor(c *Color, format ColorFormat, plain bool, commas bool, spaces bool) string {
	var b strings.Builder

	writeSep := func(sep byte) {
		if sep != ',' || commas {
			b.WriteByte(sep)
		}
		if spaces {
			b.WriteByte(' ')
		}
	}

	switch format {
	case FormatHex:
		{
			if !plain {
				b.WriteByte('#')
			}
			h, l := hexComponent(c.RGB.R)
			b.WriteByte(h)
			b.WriteByte(l)
			h, l = hexComponent(c.RGB.G)
			b.WriteByte(h)
			b.WriteByte(l)
			h, l = hexComponent(c.RGB.B)
			b.WriteByte(h)
			b.WriteByte(l)
			if c.Alpha != nil {
				alpha := uint8(*c.Alpha*255 + 0.5)
				h, l = hexComponent(alpha)
				b.WriteByte(h)
				b.WriteByte(l)
			}
		}
	case FormatHSL:
		{
			if !plain {
				b.WriteString("hsl")
				if c.Alpha != nil {
					b.WriteByte('a')
				}
				b.WriteByte('(')
			}
			b.WriteString(formatUint(c.HSL.H))
			writeSep(',')
			b.WriteString(formatUint(c.HSL.S))
			b.WriteByte('%')
			writeSep(',')
			b.WriteString(formatUint(c.HSL.L))
			b.WriteByte('%')
			if c.Alpha != nil {
				writeSep(',')
				b.WriteString(formatAlpha(*c.Alpha))
			}
			if !plain {
				b.WriteByte(')')
			}
		}

	case FormatHSLCSS:
		{
			if !plain {
				b.WriteString("hsl")
				b.WriteByte('(')
			}
			b.WriteString(formatUint(c.HSL.H))
			b.WriteString("deg ")
			b.WriteString(formatUint(c.HSL.S))
			b.WriteByte('%')
			b.WriteByte(' ')
			b.WriteString(formatUint(c.HSL.L))
			b.WriteByte('%')
			if c.Alpha != nil {
				b.WriteString(" / ")
				b.WriteString(formatAlpha(*c.Alpha))
			}
			if !plain {
				b.WriteByte(')')
			}
		}
	case FormatHSLArray:
		{
			if !plain {
				b.WriteByte('[')
			}
			b.WriteString(formatUint(c.HSL.H))
			writeSep(',')
			b.WriteString(formatAlpha(float64(c.HSL.S) / 100))
			writeSep(',')
			b.WriteString(formatAlpha(float64(c.HSL.L) / 100))
			if c.Alpha != nil {
				writeSep(',')
				b.WriteString(formatAlpha(*c.Alpha))
			}
			if !plain {
				b.WriteByte(']')
			}
		}
	case FormatRGB:
		{
			if !plain {
				b.WriteString("rgb")
				if c.Alpha != nil {
					b.WriteByte('a')
				}
				b.WriteByte('(')
			}
			b.WriteString(formatUint(c.RGB.R))
			writeSep(',')
			b.WriteString(formatUint(c.RGB.G))
			writeSep(',')
			b.WriteString(formatUint(c.RGB.B))
			if c.Alpha != nil {
				writeSep(',')
				b.WriteString(formatAlpha(*c.Alpha))
			}
			if !plain {
				b.WriteByte(')')
			}
		}
	case FormatRGBCSS:
		{
			if !plain {
				b.WriteString("rgb(")
			}
			b.WriteString(formatUint(c.RGB.R))
			b.WriteByte(' ')
			b.WriteString(formatUint(c.RGB.G))
			b.WriteByte(' ')
			b.WriteString(formatUint(c.RGB.B))
			if c.Alpha != nil {
				b.WriteString(" / ")
				b.WriteString(formatAlpha(*c.Alpha))
			}
			if !plain {
				b.WriteByte(')')
			}
		}
	case FormatRGBArray:
		{
			if !plain {
				b.WriteByte('[')
			}
			b.WriteString(formatUint(c.RGB.R))
			writeSep(',')
			b.WriteString(formatUint(c.RGB.G))
			writeSep(',')
			b.WriteString(formatUint(c.RGB.B))
			if c.Alpha != nil {
				writeSep(',')
				b.WriteString(formatAlpha(*c.Alpha))
			}
			if !plain {
				b.WriteByte(']')
			}
		}
	case FormatAnsi:
		{
			b.WriteString(formatUint(c.RGB.R))
			b.WriteByte(';')
			b.WriteString(formatUint(c.RGB.G))
			b.WriteByte(';')
			b.WriteString(formatUint(c.RGB.B))
			if c.Alpha != nil {
				b.WriteByte(';')
				b.WriteString(formatAlpha(*c.Alpha))
			}
		}
	}

	return b.String()
}

type Variant struct {
	id          string
	name        string
	appearance  string
	description string
	colors      map[string]*Color
}

const description = "All natural pine, faux fur and a bit of soho vibes for the classy minimalist"

var (
	MainVariant = Variant{
		id:          "rose-pine",
		name:        "Rosé Pine",
		appearance:  "dark",
		description: description,
		colors: map[string]*Color{
			"base": {
				HSL: hsl(249, 22, 12),
				RGB: rgb(25, 23, 36),
			},
			"surface": {
				HSL: hsl(247, 23, 15),
				RGB: rgb(31, 29, 46),
			},
			"overlay": {
				HSL: hsl(245, 25, 18),
				RGB: rgb(38, 35, 58),
			},
			"muted": {
				HSL: hsl(249, 12, 47),
				RGB: rgb(110, 106, 134),
			},
			"subtle": {
				HSL: hsl(248, 15, 61),
				RGB: rgb(144, 140, 170),
			},
			"text": {
				HSL: hsl(245, 50, 91),
				RGB: rgb(224, 222, 244),
			},
			"love": {
				HSL: hsl(343, 76, 68),
				RGB: rgb(235, 111, 146),
				On:  "text",
			},
			"gold": {
				HSL: hsl(35, 88, 72),
				RGB: rgb(246, 193, 119),
				On:  "surface",
			},
			"rose": {
				HSL: hsl(2, 55, 83),
				RGB: rgb(235, 188, 186),
				On:  "surface",
			},
			"pine": {
				HSL: hsl(197, 49, 38),
				RGB: rgb(49, 116, 143),
				On:  "text",
			},
			"foam": {
				HSL: hsl(189, 43, 73),
				RGB: rgb(156, 207, 216),
				On:  "surface",
			},
			"iris": {
				HSL: hsl(267, 57, 78),
				RGB: rgb(196, 167, 231),
				On:  "surface",
			},
			"highlightLow": {
				HSL: hsl(244, 18, 15),
				RGB: rgb(33, 32, 46),
			},
			"highlightMed": {
				HSL: hsl(247, 15, 28),
				RGB: rgb(64, 61, 82),
			},
			"highlightHigh": {
				HSL: hsl(245, 13, 36),
				RGB: rgb(82, 79, 103),
			},
		},
	}

	MoonVariant = Variant{
		id:          "rose-pine-moon",
		name:        "Rosé Pine Moon",
		appearance:  "dark",
		description: description,
		colors: map[string]*Color{
			"base": {
				HSL: hsl(246, 24, 17),
				RGB: rgb(35, 33, 54),
			},
			"surface": {
				HSL: hsl(248, 24, 20),
				RGB: rgb(42, 39, 63),
			},
			"overlay": {
				HSL: hsl(248, 21, 26),
				RGB: rgb(57, 53, 82),
			},
			"muted": {
				HSL: hsl(249, 12, 47),
				RGB: rgb(110, 106, 134),
			},
			"subtle": {
				HSL: hsl(248, 15, 61),
				RGB: rgb(144, 140, 170),
			},
			"text": {
				HSL: hsl(245, 50, 91),
				RGB: rgb(224, 222, 244),
			},
			"love": {
				HSL: hsl(343, 76, 68),
				RGB: rgb(235, 111, 146),
				On:  "text",
			},
			"gold": {
				HSL: hsl(35, 88, 72),
				RGB: rgb(246, 193, 119),
				On:  "surface",
			},
			"rose": {
				HSL: hsl(2, 66, 75),
				RGB: rgb(234, 154, 151),
				On:  "surface",
			},
			"pine": {
				HSL: hsl(197, 48, 47),
				RGB: rgb(62, 143, 176),
				On:  "text",
			},
			"foam": {
				HSL: hsl(189, 43, 73),
				RGB: rgb(156, 207, 216),
				On:  "surface",
			},
			"iris": {
				HSL: hsl(267, 57, 78),
				RGB: rgb(196, 167, 231),
				On:  "surface",
			},
			"highlightLow": {
				HSL: hsl(245, 22, 20),
				RGB: rgb(42, 40, 62),
			},
			"highlightMed": {
				HSL: hsl(247, 16, 30),
				RGB: rgb(68, 65, 90),
			},
			"highlightHigh": {
				HSL: hsl(249, 15, 38),
				RGB: rgb(86, 82, 110),
			},
		},
	}

	DawnVariant = Variant{
		id:          "rose-pine-dawn",
		name:        "Rosé Pine Dawn",
		appearance:  "light",
		description: description,
		colors: map[string]*Color{
			"base": {
				HSL: hsl(32, 57, 95),
				RGB: rgb(250, 244, 237),
			},
			"surface": {
				HSL: hsl(35, 100, 98),
				RGB: rgb(255, 250, 243),
			},
			"overlay": {
				HSL: hsl(25, 36, 92),
				RGB: rgb(242, 233, 225),
			},
			"muted": {
				HSL: hsl(254, 9, 61),
				RGB: rgb(152, 147, 165),
			},
			"subtle": {
				HSL: hsl(249, 13, 52),
				RGB: rgb(121, 117, 147),
			},
			"text": {
				HSL: hsl(248, 19, 40),
				RGB: rgb(87, 82, 121),
			},
			"love": {
				HSL: hsl(343, 35, 55),
				RGB: rgb(180, 99, 122),
				On:  "surface",
			},
			"gold": {
				HSL: hsl(35, 81, 56),
				RGB: rgb(234, 157, 52),
				On:  "surface",
			},
			"rose": {
				HSL: hsl(2, 55, 67),
				RGB: rgb(215, 130, 126),
				On:  "surface",
			},
			"pine": {
				HSL: hsl(197, 53, 34),
				RGB: rgb(40, 105, 131),
				On:  "surface",
			},
			"foam": {
				HSL: hsl(189, 30, 48),
				RGB: rgb(86, 148, 159),
				On:  "surface",
			},
			"iris": {
				HSL: hsl(267, 22, 57),
				RGB: rgb(144, 122, 169),
				On:  "surface",
			},
			"highlightLow": {
				HSL: hsl(25, 35, 93),
				RGB: rgb(244, 237, 232),
			},
			"highlightMed": {
				HSL: hsl(10, 9, 86),
				RGB: rgb(223, 218, 217),
			},
			"highlightHigh": {
				HSL: hsl(315, 4, 80),
				RGB: rgb(206, 202, 205),
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
