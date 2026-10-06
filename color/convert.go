package color

import "math"

func hslToRgb(hsl *HSL) RGB {
	h := float64(hsl.H)
	s := float64(hsl.S) / 100.0
	l := float64(hsl.L) / 100.0

	c := (1.0 - math.Abs(2.0*l-1.0)) * s
	x := c * (1.0 - math.Abs(math.Mod(h/60.0, 2.0)-1.0))
	m := l - c/2.0

	var r, g, b float64
	switch {
	case h >= 0 && h < 60:
		r, g, b = c, x, 0
	case h >= 60 && h < 120:
		r, g, b = x, c, 0
	case h >= 120 && h < 180:
		r, g, b = 0, c, x
	case h >= 180 && h < 240:
		r, g, b = 0, x, c
	case h >= 240 && h < 300:
		r, g, b = x, 0, c
	case h >= 300 && h < 360:
		r, g, b = c, 0, x
	default:
		r, g, b = 0, 0, 0
	}

	return RGB{
		R: uint8(math.Round((r + m) * 255.0)),
		G: uint8(math.Round((g + m) * 255.0)),
		B: uint8(math.Round((b + m) * 255.0)),
	}
}

func rgbToHsl(rgb *RGB) HSL {
	r := float64(rgb.R) / 255.0
	g := float64(rgb.G) / 255.0
	b := float64(rgb.B) / 255.0

	max := math.Max(r, math.Max(g, b))
	min := math.Min(r, math.Min(g, b))
	l := (max + min) / 2

	var h, s float64
	if d := max - min; d != 0 {
		if l > 0.5 {
			s = d / (2 - max - min)
		} else {
			s = d / (max + min)
		}

		switch {
		case r == max:
			h = (g - b) / d
			if g < b {
				h += 6
			}
		case g == max:
			h = (b-r)/d + 2
		default:
			h = (r-g)/d + 4
		}
		h *= 60
	}

	return HSL{
		H: uint16(math.Round(math.Mod(h, 360))),
		S: uint8(math.Round(s * 100)),
		L: uint8(math.Round(l * 100)),
	}
}

func ColorFromHSL(hsl HSL) Color {
	return Color{
		HSL: hsl,
		RGB: hslToRgb(&hsl),
	}
}

func ColorFromRGB(rgb RGB) Color {
	return Color{
		HSL: rgbToHsl(&rgb),
		RGB: rgb,
	}
}
