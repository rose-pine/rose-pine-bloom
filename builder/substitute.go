package builder

import (
	"fmt"
	"strings"

	"github.com/rose-pine/rose-pine-bloom/color"
)

func substituteCaptures(content string, captures []Capture, variant color.VariantMeta, opts *BuildOpts, accentName string) (string, error) {
	var palette color.Palette
	var buf strings.Builder

	buf.Grow(len(content))

	switch variant.Id {
	case "rose-pine":
		palette = color.MainPalette
	case "rose-pine-moon":
		palette = color.MoonPalette
	case "rose-pine-dawn":
		palette = color.DawnPalette
	default:
		return "", fmt.Errorf("unknown variant `%s`", variant.Id)
	}

	for _, capture := range captures {
		start, length := Span(capture)
		text := content[start : start+length]

		switch c := capture.(type) {
		case RoleCapture:
			var clr *color.Color

			if c.role == "accent" || c.role == "onaccent" {
				accentColor, ok := variant.Colors[accentName]
				if !ok {
					return "", fmt.Errorf("unknown accent color `%s`", accentName)
				}

				switch c.role {
				case "accent":
					clr = accentColor

				case "onaccent":
					if accentColor.On == "" {
						return "", fmt.Errorf("accent color `%s` does not support onaccent", accentColor.On)
					}
					if onAccentColor, ok := variant.Colors[accentColor.On]; ok {
						clr = onAccentColor
					} else {
						return "", fmt.Errorf("invalid role value `%s` for onaccent", accentColor.On)
					}
				}
			} else {
				clr = palette[c.role]
			}

			if clr == nil {
				return "", fmt.Errorf("no such role: `%s`", c.role)
			}
			withAlpha := clr.WithAlpha(c.alpha)
			clr = &withAlpha

			formatted := color.FormatColor(clr, opts.DefaultFormat, opts.Plain, opts.Commas, opts.Spaces)
			buf.WriteString(formatted)

		case MetaCapture:
			switch text[1:] {
			case "id":
				buf.WriteString(variant.Id)
			case "name":
				buf.WriteString(variant.Name)
			case "appearance":
				buf.WriteString(variant.Appearance)
			case "type":
				buf.WriteString(variant.Appearance)
			case "description":
				buf.WriteString(variant.Description)
			case "accentname":
				buf.WriteString(accentName)
			}

		case VariantCapture:
			var inner variantArm
			switch variant.Id {
			case "rose-pine":
				inner = c.main
			case "rose-pine-moon":
				inner = c.moon
			case "rose-pine-dawn":
				inner = c.dawn
			}
			output, err := substituteCaptures(inner.content, inner.captures, variant, opts, accentName)
			if err != nil {
				return "", err
			}

			buf.WriteString(output)

		case TextCapture:
			buf.WriteString(text)
		}
	}

	return buf.String(), nil
}
