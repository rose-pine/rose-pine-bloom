package builder

import (
	"fmt"
	"strings"

	"github.com/rose-pine/rose-pine-bloom/color"
)

func substituteCaptures(content string, captures []Capture, variant color.VariantMeta, opts *BuildOpts, accentName string) (string, error) {
	var shades color.PaletteShades
	var buf strings.Builder

	buf.Grow(len(content))

	switch variant.Id {
	case "rose-pine":
		shades = color.MainPaletteShades
	case "rose-pine-moon":
		shades = color.MoonPaletteShades
	case "rose-pine-dawn":
		shades = color.DawnPaletteShades
	default:
		return "", fmt.Errorf("unknown variant `%s`", variant.Id)
	}

	for _, capture := range captures {
		start, length := Span(capture)
		text := content[start : start+length]

		switch c := capture.(type) {
		case RoleCapture:
			roleName := c.role
			if c.role == "accent" || c.role == "onaccent" {
				accentColor, ok := variant.Colors[accentName]
				if !ok {
					return "", fmt.Errorf("unknown accent color `%s`", accentName)
				}

				switch c.role {
				case "accent":
					roleName = accentName
				case "onaccent":
					if accentColor.On == "" {
						return "", fmt.Errorf("accent color `%s` does not support onaccent", accentName)
					}
					roleName = accentColor.On
				}
			}

			roleShades, ok := shades[roleName]
			if !ok {
				return "", fmt.Errorf("no such role: `%s`", roleName)
			}

			shadeIndex := roleShades.Base
			if c.shade != nil {
				idx, err := color.GetShadeIndex(*c.shade)
				if err != nil {
					return "", fmt.Errorf("invalid shade value `%d` for role `%s`: %w", *c.shade, roleName, err)
				}
				shadeIndex = idx
			}

			clr := roleShades.Colors[shadeIndex].WithAlpha(c.alpha)
			formatted := color.FormatColor(&clr, opts.DefaultFormat, opts.Plain, opts.Commas, opts.Spaces)
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
