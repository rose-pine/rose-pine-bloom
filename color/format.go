package color

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
)

type RGB struct {
	R, G, B uint8
}

type HSL struct {
	H    uint16
	S, L uint8
}

type format struct {
	Name    string
	Example string
}

var formats = [...]format{
	{Name: "hex", Example: "#ebbcba"},
	{Name: "hex --plain", Example: "ebbcba"},

	{Name: "hsl", Example: "hsl(2, 55%, 83%)"},
	{Name: "hsl --plain", Example: "2, 55%, 83%"},
	{Name: "hsl-css", Example: "hsl(2deg 55% 83%)"},
	{Name: "hsl-css --plain", Example: "2deg 55% 83%"},
	{Name: "hsl-array", Example: "[2, 0.55, 0.83]"},
	{Name: "hsl-array --plain", Example: "2, 0.55, 0.83"},

	{Name: "rgb", Example: "rgb(235, 188, 186)"},
	{Name: "rgb --plain", Example: "235, 188, 186"},
	{Name: "rgb-css", Example: "rgb(235 188 186)"},
	{Name: "rgb-css --plain", Example: "235 188 186"},
	{Name: "rgb-array", Example: "[235, 188, 186]"},
	{Name: "rgb-array --plain", Example: "235, 188, 186"},

	{Name: "ansi", Example: "235;188;186"},
}

func FormatsTable() string {
	var sb strings.Builder
	w := tabwriter.NewWriter(&sb, 1, 1, 1, ' ', 0)
	for _, f := range formats {
		fmt.Fprintf(w, "    %-23s %s\n", f.Name, f.Example)
	}
	w.Flush()
	return sb.String()
}

func PrintFormatsTable() {
	fmt.Fprint(os.Stdout, FormatsTable())
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

func FormatColor(c *Color, format ColorFormat, plain bool, commas bool, spaces bool) string {
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
