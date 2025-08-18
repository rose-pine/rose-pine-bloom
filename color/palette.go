package color

import (
	"iter"
)

type Variant int

const (
	Main Variant = iota
	Moon
	Dawn
)

type Palette struct {
	Base          Color
	Surface       Color
	Overlay       Color
	Muted         Color
	Subtle        Color
	Text          Color
	Love          Color
	Gold          Color
	Rose          Color
	Pine          Color
	Foam          Color
	Iris          Color
	HighlightLow  Color
	HighlightMed  Color
	HighlightHigh Color
}

type VariantMeta struct {
	Id          string
	Name        string
	Appearance  string
	Description string
	Colors      *Palette
}

var Accents = []string{
	"love", "gold", "rose", "pine", "foam", "iris",
}

var (
	MainPalette = Palette{
		Base: Color{
			HSL: Hsl(249, 22, 12),
			RGB: Rgb(25, 23, 36),
		},
		Surface: Color{
			HSL: Hsl(247, 23, 15),
			RGB: Rgb(31, 29, 46),
		},
		Overlay: Color{
			HSL: Hsl(245, 25, 18),
			RGB: Rgb(38, 35, 58),
		},
		Muted: Color{
			HSL: Hsl(249, 12, 47),
			RGB: Rgb(110, 106, 134),
		},
		Subtle: Color{
			HSL: Hsl(248, 15, 61),
			RGB: Rgb(144, 140, 170),
		},
		Text: Color{
			HSL: Hsl(245, 50, 91),
			RGB: Rgb(224, 222, 244),
		},
		Love: Color{
			HSL: Hsl(343, 76, 68),
			RGB: Rgb(235, 111, 146),
			On:  "text",
		},
		Gold: Color{
			HSL: Hsl(35, 88, 72),
			RGB: Rgb(246, 193, 119),
			On:  "surface",
		},
		Rose: Color{
			HSL: Hsl(2, 55, 83),
			RGB: Rgb(235, 188, 186),
			On:  "surface",
		},
		Pine: Color{
			HSL: Hsl(197, 49, 38),
			RGB: Rgb(49, 116, 143),
			On:  "text",
		},
		Foam: Color{
			HSL: Hsl(189, 43, 73),
			RGB: Rgb(156, 207, 216),
			On:  "surface",
		},
		Iris: Color{
			HSL: Hsl(267, 57, 78),
			RGB: Rgb(196, 167, 231),
			On:  "surface",
		},
		HighlightLow: Color{
			HSL: Hsl(244, 18, 15),
			RGB: Rgb(33, 32, 46),
		},
		HighlightMed: Color{
			HSL: Hsl(247, 15, 28),
			RGB: Rgb(64, 61, 82),
		},
		HighlightHigh: Color{
			HSL: Hsl(245, 13, 36),
			RGB: Rgb(82, 79, 103),
		},
	}

	MoonPalette = Palette{
		Base: Color{
			HSL: Hsl(246, 24, 17),
			RGB: Rgb(35, 33, 54),
		},
		Surface: Color{
			HSL: Hsl(248, 24, 20),
			RGB: Rgb(42, 39, 63),
		},
		Overlay: Color{
			HSL: Hsl(248, 21, 26),
			RGB: Rgb(57, 53, 82),
		},
		Muted: Color{
			HSL: Hsl(249, 12, 47),
			RGB: Rgb(110, 106, 134),
		},
		Subtle: Color{
			HSL: Hsl(248, 15, 61),
			RGB: Rgb(144, 140, 170),
		},
		Text: Color{
			HSL: Hsl(245, 50, 91),
			RGB: Rgb(224, 222, 244),
		},
		Love: Color{
			HSL: Hsl(343, 76, 68),
			RGB: Rgb(235, 111, 146),
			On:  "text",
		},
		Gold: Color{
			HSL: Hsl(35, 88, 72),
			RGB: Rgb(246, 193, 119),
			On:  "surface",
		},
		Rose: Color{
			HSL: Hsl(2, 66, 75),
			RGB: Rgb(234, 154, 151),
			On:  "surface",
		},
		Pine: Color{
			HSL: Hsl(197, 48, 47),
			RGB: Rgb(62, 143, 176),
			On:  "text",
		},
		Foam: Color{
			HSL: Hsl(189, 43, 73),
			RGB: Rgb(156, 207, 216),
			On:  "surface",
		},
		Iris: Color{
			HSL: Hsl(267, 57, 78),
			RGB: Rgb(196, 167, 231),
			On:  "surface",
		},
		HighlightLow: Color{
			HSL: Hsl(245, 22, 20),
			RGB: Rgb(42, 40, 62),
		},
		HighlightMed: Color{
			HSL: Hsl(247, 16, 30),
			RGB: Rgb(68, 65, 90),
		},
		HighlightHigh: Color{
			HSL: Hsl(249, 15, 38),
			RGB: Rgb(86, 82, 110),
		},
	}

	DawnPalette = Palette{
		Base: Color{
			HSL: Hsl(32, 57, 95),
			RGB: Rgb(250, 244, 237),
		},
		Surface: Color{
			HSL: Hsl(35, 100, 98),
			RGB: Rgb(255, 250, 243),
		},
		Overlay: Color{
			HSL: Hsl(25, 36, 92),
			RGB: Rgb(242, 233, 225),
		},
		Muted: Color{
			HSL: Hsl(254, 9, 61),
			RGB: Rgb(152, 147, 165),
		},
		Subtle: Color{
			HSL: Hsl(249, 13, 52),
			RGB: Rgb(121, 117, 147),
		},
		Text: Color{
			HSL: Hsl(248, 19, 40),
			RGB: Rgb(87, 82, 121),
		},
		Love: Color{
			HSL: Hsl(343, 35, 55),
			RGB: Rgb(180, 99, 122),
			On:  "surface",
		},
		Gold: Color{
			HSL: Hsl(35, 81, 56),
			RGB: Rgb(234, 157, 52),
			On:  "surface",
		},
		Rose: Color{
			HSL: Hsl(2, 55, 67),
			RGB: Rgb(215, 130, 126),
			On:  "surface",
		},
		Pine: Color{
			HSL: Hsl(197, 53, 34),
			RGB: Rgb(40, 105, 131),
			On:  "surface",
		},
		Foam: Color{
			HSL: Hsl(189, 30, 48),
			RGB: Rgb(86, 148, 159),
			On:  "surface",
		},
		Iris: Color{
			HSL: Hsl(267, 22, 57),
			RGB: Rgb(144, 122, 169),
			On:  "surface",
		},
		HighlightLow: Color{
			HSL: Hsl(25, 35, 93),
			RGB: Rgb(244, 237, 232),
		},
		HighlightMed: Color{
			HSL: Hsl(10, 9, 86),
			RGB: Rgb(223, 218, 217),
		},
		HighlightHigh: Color{
			HSL: Hsl(315, 4, 80),
			RGB: Rgb(206, 202, 205),
		},
	}
)

const description = "All natural pine, faux fur and a bit of soho vibes for the classy minimalist"

var (
	MainVariantMeta = VariantMeta{
		Id:          "rose-pine",
		Name:        "Rosé Pine",
		Appearance:  "dark",
		Description: description,
		Colors:      &MainPalette,
	}

	MoonVariantMeta = VariantMeta{
		Id:          "rose-pine-moon",
		Name:        "Rosé Pine Moon",
		Appearance:  "dark",
		Description: description,
		Colors:      &MoonPalette,
	}

	DawnVariantMeta = VariantMeta{
		Id:          "rose-pine-dawn",
		Name:        "Rosé Pine Dawn",
		Appearance:  "light",
		Description: description,
		Colors:      &DawnPalette,
	}
)

var Variants = []VariantMeta{
	MainVariantMeta,
	MoonVariantMeta,
	DawnVariantMeta,
}

func (p *Palette) Get(role string) (*Color, bool) {
	switch role {
	case "base":
		return &p.Base, true
	case "surface":
		return &p.Surface, true
	case "overlay":
		return &p.Overlay, true
	case "muted":
		return &p.Muted, true
	case "subtle":
		return &p.Subtle, true
	case "text":
		return &p.Text, true
	case "love":
		return &p.Love, true
	case "gold":
		return &p.Gold, true
	case "rose":
		return &p.Rose, true
	case "pine":
		return &p.Pine, true
	case "foam":
		return &p.Foam, true
	case "iris":
		return &p.Iris, true
	case "highlightLow":
		return &p.HighlightLow, true
	case "highlightMed":
		return &p.HighlightMed, true
	case "highlightHigh":
		return &p.HighlightHigh, true
	default:
		return nil, false
	}
}

func (p *Palette) Iter() iter.Seq2[string, *Color] {
	return func(yield func(string, *Color) bool) {
		if !yield("base", &p.Base) {
			return
		}
		if !yield("surface", &p.Surface) {
			return
		}
		if !yield("overlay", &p.Overlay) {
			return
		}
		if !yield("muted", &p.Muted) {
			return
		}
		if !yield("subtle", &p.Subtle) {
			return
		}
		if !yield("text", &p.Text) {
			return
		}
		if !yield("love", &p.Love) {
			return
		}
		if !yield("gold", &p.Gold) {
			return
		}
		if !yield("rose", &p.Rose) {
			return
		}
		if !yield("pine", &p.Pine) {
			return
		}
		if !yield("foam", &p.Foam) {
			return
		}
		if !yield("iris", &p.Iris) {
			return
		}
		if !yield("highlightLow", &p.HighlightLow) {
			return
		}
		if !yield("highlightMed", &p.HighlightMed) {
			return
		}
		if !yield("highlightHigh", &p.HighlightHigh) {
			return
		}
	}
}
