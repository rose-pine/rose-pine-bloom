package color

import "math"

const (
	maxLightness = 100
	maxSteps     = 10
)

func Blend(base *Color, color *Color, alpha float64) Color {
	return ColorFromRGB(RGB{
		R: blendComponent(base.RGB.R, color.RGB.R, alpha),
		G: blendComponent(base.RGB.G, color.RGB.G, alpha),
		B: blendComponent(base.RGB.B, color.RGB.B, alpha),
	})
}

func blendComponent(base uint8, color uint8, alpha float64) uint8 {
	return uint8(math.Round(float64(base)*(1-alpha) + float64(color)*alpha))
}

func Lighten(color *Color, steps uint8) Color {
	stepSize := (maxLightness - color.HSL.L) / maxSteps
	l := min(color.HSL.L+steps*stepSize, maxLightness)

	return ColorFromHSL(HSL{
		H: color.HSL.H,
		S: color.HSL.S,
		L: l,
	})
}

func Darken(color *Color, steps uint8) Color {
	stepSize := color.HSL.L / maxSteps
	sub := steps * stepSize
	l := color.HSL.L
	if sub < l {
		l -= sub
	} else {
		l = 0
	}

	return ColorFromHSL(HSL{
		H: color.HSL.H,
		S: color.HSL.S,
		L: l,
	})
}
