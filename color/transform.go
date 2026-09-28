package color

const (
	maxLightness = 100
	maxSteps     = 10
)

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
