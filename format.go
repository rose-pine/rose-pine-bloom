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
