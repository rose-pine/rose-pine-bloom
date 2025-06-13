package main

import (
	"fmt"
	"strings"
)

type Color struct {
	Hex   string     `json:"hex"`
	RGB   [3]int     `json:"rgb"`
	HSL   [3]float64 `json:"hsl"`
	Alpha *float64   `json:"alpha,omitempty"`
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
	case FormatHSLArray:
		hsl := fmt.Sprintf("[%v, %v%%, %v%%", c.HSL[0], c.HSL[1], c.HSL[2])
		if c.Alpha != nil {
			hsl += fmt.Sprintf(", %s", formatAlpha(*c.Alpha))
		}
		hsl += "]"
		workingString = hsl
	case FormatHSLFunc:
		prefix := "hsl"
		if c.Alpha != nil {
			prefix = "hsla"
		}
		hsl := fmt.Sprintf("%s(%v, %v%%, %v%%", prefix, c.HSL[0], c.HSL[1], c.HSL[2])
		if c.Alpha != nil {
			hsl += fmt.Sprintf(", %s", formatAlpha(*c.Alpha))
		}
		hsl += ")"
		workingString = hsl
	case FormatRGB:
		rgb := fmt.Sprintf("%d, %d, %d", c.RGB[0], c.RGB[1], c.RGB[2])
		if c.Alpha != nil {
			rgb += fmt.Sprintf(", %s", formatAlpha(*c.Alpha))
		}
		workingString = rgb
	case FormatRGBAnsi:
		if c.Alpha != nil {
			return fmt.Sprintf("%d;%d;%d;%s", c.RGB[0], c.RGB[1], c.RGB[2], formatAlpha(*c.Alpha))
		}
		return fmt.Sprintf("%d;%d;%d", c.RGB[0], c.RGB[1], c.RGB[2])
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
