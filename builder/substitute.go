package builder

import (
	"fmt"
	"strings"

	"github.com/rose-pine/rose-pine-bloom/color"
)

func (opts *BuildOpts) varName(name string) string {
	return string(opts.Prefix) + name
}

func qualifiedForms(opts *BuildOpts, name string) string {
	forms := make([]string, 0, len(variantKeys))
	for _, key := range variantKeys {
		forms = append(forms, opts.varName(key+"-"+name))
	}
	if len(forms) == 0 {
		return ""
	}
	return strings.Join(forms[:len(forms)-1], ", ") + " or " + forms[len(forms)-1]
}

func roleVariant(qualifier string, name string, opts *BuildOpts, variant *color.VariantMeta) (*color.VariantMeta, error) {
	if qualifier == "" {
		if opts.Single {
			return nil, fmt.Errorf("variant prefix required: %s (use %s)", opts.varName(name), qualifiedForms(opts, name))
		}
		if variant == nil {
			return nil, fmt.Errorf("internal error: no variant to resolve %s against", opts.varName(name))
		}

		return variant, nil
	}

	return color.VariantsByKey[qualifier], nil
}

func substituteCaptures(content string, captures []Capture, variant *color.VariantMeta, opts *BuildOpts, accentName string) (string, error) {
	var buf strings.Builder

	buf.Grow(len(content))

	for _, capture := range captures {
		start, length := Span(capture)
		text := content[start : start+length]

		switch c := capture.(type) {
		case RoleCapture:
			target, err := roleVariant(c.qualifier, c.role, opts, variant)
			if err != nil {
				return "", err
			}

			roleName := c.role
			if c.role == "accent" || c.role == "onaccent" {
				if opts.Single {
					return "", fmt.Errorf("%s is not supported in single-file mode, use %s-<name>", opts.varName(c.qualifier+"-"+c.role), opts.varName(c.qualifier))
				}

				accentColor, ok := target.Colors[accentName]
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

			roleShades, ok := target.Shades[roleName]
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
			target, err := roleVariant(c.qualifier, c.key, opts, variant)
			if err != nil {
				return "", err
			}

			switch c.key {
			case "id":
				buf.WriteString(target.Id)
			case "name":
				buf.WriteString(target.Name)
			case "appearance":
				buf.WriteString(target.Appearance)
			case "type":
				buf.WriteString(target.Appearance)
			case "description":
				buf.WriteString(target.Description)
			case "accentname":
				if opts.Single {
					return "", fmt.Errorf("%s is not supported in single-file mode", opts.varName(c.qualifier+"-accentname"))
				}
				buf.WriteString(accentName)
			}

		case VariantCapture:
			if opts.Single {
				return "", fmt.Errorf("variant values are not supported in single-file mode, use %s instead", qualifiedForms(opts, "..."))
			}

			inner := c.arms[variant.Key]

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
