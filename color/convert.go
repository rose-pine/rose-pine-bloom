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

func ColorFromHSL(hsl HSL) Color {
	return Color{
		HSL: hsl,
		RGB: hslToRgb(&hsl),
	}
}
