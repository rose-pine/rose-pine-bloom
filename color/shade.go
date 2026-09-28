package color

import (
	"fmt"
	"math"
)

type Shades struct {
	Base   int
	Colors [10]Color
}

func InferBaseLightness(c *Color) int {
	return min(int(math.Round(float64(c.HSL.L)/10)), 9)
}

func GetShadeIndex(value int) (int, error) {
	if value == 50 {
		return 0, nil
	}
	if value >= 100 && value <= 900 && value%100 == 0 {
		return value / 100, nil
	}

	return 0, fmt.Errorf("unknown shade: `%d`", value)
}

func genShades(c *Color) Shades {
	base := int(InferBaseLightness(c))
	var colors [10]Color

	for i := range 10 {
		if i < base {
			colors[i] = Darken(c, uint8(base-i))
		} else if i > base {
			colors[i] = Lighten(c, uint8(i-base))
		} else {
			colors[i] = *c
		}
	}

	return Shades{
		base,
		colors,
	}
}

type PaletteShades map[string]Shades

func genPaletteShades(p Palette) PaletteShades {
	out := make(PaletteShades, len(p))
	for k, v := range p {
		out[k] = genShades(v)
	}
	return out
}

var (
	MainPaletteShades = genPaletteShades(MainPalette)
	MoonPaletteShades = genPaletteShades(MoonPalette)
	DawnPaletteShades = genPaletteShades(DawnPalette)
)
