package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestColorFormatting(t *testing.T) {
	color := &Color{
		Hex: "ebbcba",
		RGB: [3]int{235, 188, 186},
		HSL: [3]float64{2, 55, 83},
	}

	tests := []struct {
		name        string
		format      ColorFormat
		stripSpaces bool
		want        string
	}{
		{"hex", FormatHex, false, "#ebbcba"},
		{"hex-ns", FormatHexNS, false, "ebbcba"},
		{"rgb", FormatRGB, false, "235, 188, 186"},
		{"rgb stripped", FormatRGB, true, "235,188,186"},
		{"rgb-ns", FormatRGBNS, false, "235 188 186"},
		{"rgb-ns stripped", FormatRGBNS, true, "235188186"},
		{"rgb-ansi", FormatRGBAnsi, false, "235;188;186"},
		{"rgb-array", FormatRGBArray, false, "[235, 188, 186]"},
		{"rgb-array stripped", FormatRGBArray, true, "[235,188,186]"},
		{"rgb-function", FormatRGBFunc, false, "rgb(235, 188, 186)"},
		{"rgb-function stripped", FormatRGBFunc, true, "rgb(235,188,186)"},
		{"hsl", FormatHSL, false, "2, 55%, 83%"},
		{"hsl stripped", FormatHSL, true, "2,55%,83%"},
		{"hsl-ns", FormatHSLNS, false, "2 55% 83%"},
		{"hsl-ns stripped", FormatHSLNS, true, "255%83%"},
		{"hsl-array", FormatHSLArray, false, "[2, 55%, 83%]"},
		{"hsl-array stripped", FormatHSLArray, true, "[2,55%,83%]"},
		{"hsl-function", FormatHSLFunc, false, "hsl(2, 55%, 83%)"},
		{"hsl-function stripped", FormatHSLFunc, true, "hsl(2,55%,83%)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatColor(color, tt.format, tt.stripSpaces)
			if got != tt.want {
				t.Errorf("formatColor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestColorFormattingWithAlpha(t *testing.T) {
	alpha := 0.5
	color := &Color{
		Hex:   "ebbcba",
		RGB:   [3]int{235, 188, 186},
		HSL:   [3]float64{2, 55, 83},
		Alpha: &alpha,
	}

	tests := []struct {
		name        string
		format      ColorFormat
		stripSpaces bool
		want        string
	}{
		{"rgb with alpha", FormatRGB, false, "235, 188, 186, 0.5"},
		{"rgb-ns with alpha", FormatRGBNS, false, "235 188 186 0.5"},
		{"rgb-ansi with alpha", FormatRGBAnsi, false, "235;188;186;0.5"},
		{"rgb-array with alpha", FormatRGBArray, false, "[235, 188, 186, 0.5]"},
		{"rgb-function with alpha", FormatRGBFunc, false, "rgba(235, 188, 186, 0.5)"},
		{"hsl with alpha", FormatHSL, false, "2, 55%, 83%, 0.5"},
		{"hsl-ns with alpha", FormatHSLNS, false, "2 55% 83% 0.5"},
		{"hsl-array with alpha", FormatHSLArray, false, "[2, 55%, 83%, 0.5]"},
		{"hsl-function with alpha", FormatHSLFunc, false, "hsla(2, 55%, 83%, 0.5)"},

		{"rgb-array with alpha stripped", FormatRGBArray, true, "[235,188,186,0.5]"},
		{"hsl-array with alpha stripped", FormatHSLArray, true, "[2,55%,83%,0.5]"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatColor(color, tt.format, tt.stripSpaces)
			if got != tt.want {
				t.Errorf("formatColor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlphaVariables(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "rose-pine-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	templateContent := `{
        "regular": "$base",
        "alpha50": "$base/50",
        "alpha25": "$base/25",
        "alphaMuted50": "$muted/50"
    }`

	templatePath := filepath.Join(tmpDir, "template.json")
	if err := os.WriteFile(templatePath, []byte(templateContent), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := &Config{
		Template:    templatePath,
		Output:      tmpDir,
		Format:      "rgb",
		Prefix:      "$",
		StripSpaces: false,
	}

	if err := Build(cfg); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(filepath.Join(tmpDir, "rose-pine.json"))
	if err != nil {
		t.Fatal(err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(content, &result); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		field string
		want  string
	}{
		{"regular", "25, 23, 36"},
		{"alpha50", "25, 23, 36, 0.5"},
		{"alpha25", "25, 23, 36, 0.25"},
		{"alphaMuted50", "110, 106, 134, 0.5"},
	}

	for _, tt := range tests {
		t.Run(tt.field, func(t *testing.T) {
			if got := result[tt.field]; got != tt.want {
				t.Errorf("%s = %v, want %v", tt.field, got, tt.want)
			}
		})
	}
}

func TestVariantGeneration(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "rose-pine-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	templateContent := `{
        "id": "$id",
        "name": "$name",
        "description": "$description",
        "type": "$type",
        "colors": {
            "base": "$base",
            "surface": "$surface",
            "love": "$love"
        },
        "custom": "$(main|moon|dawn)"
    }`

	templatePath := filepath.Join(tmpDir, "template.json")
	if err := os.WriteFile(templatePath, []byte(templateContent), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := &Config{
		Template:    templatePath,
		Output:      tmpDir,
		Format:      "hex",
		Prefix:      "$",
		StripSpaces: false,
	}

	if err := Build(cfg); err != nil {
		t.Fatal(err)
	}

	variants := []struct {
		filename string
		id       string
		name     string
		varType  string
		baseHex  string
		custom   string
	}{
		{
			filename: "rose-pine.json",
			id:       "rose-pine",
			name:     "Rosé Pine",
			varType:  "dark",
			baseHex:  "#191724",
			custom:   "main",
		},
		{
			filename: "rose-pine-moon.json",
			id:       "rose-pine-moon",
			name:     "Rosé Pine Moon",
			varType:  "dark",
			baseHex:  "#232136",
			custom:   "moon",
		},
		{
			filename: "rose-pine-dawn.json",
			id:       "rose-pine-dawn",
			name:     "Rosé Pine Dawn",
			varType:  "light",
			baseHex:  "#faf4ed",
			custom:   "dawn",
		},
	}

	for _, v := range variants {
		t.Run(v.filename, func(t *testing.T) {
			content, err := os.ReadFile(filepath.Join(tmpDir, v.filename))
			if err != nil {
				t.Fatalf("Failed to read generated file: %v", err)
			}

			var result map[string]interface{}
			if err := json.Unmarshal(content, &result); err != nil {
				t.Fatalf("Failed to parse JSON: %v", err)
			}

			tests := []struct {
				field string
				want  string
			}{
				{"id", v.id},
				{"name", v.name},
				{"type", v.varType},
				{"custom", v.custom},
			}

			for _, tt := range tests {
				if got := result[tt.field]; got != tt.want {
					t.Errorf("%s = %v, want %v", tt.field, got, tt.want)
				}
			}

			colors := result["colors"].(map[string]interface{})
			if got := colors["base"]; got != v.baseHex {
				t.Errorf("base color = %v, want %v", got, v.baseHex)
			}
		})
	}
}

func TestVariantSpecificValues(t *testing.T) {
	templateContent := `{
        "accent": "$(#ebbcba|#c4a7e7|#286983)",
        "name": "$(Main|Moon|Dawn)",
        "mood": "$(Dark|Dim|Light)"
    }`

	tmpDir, err := os.MkdirTemp("", "rose-pine-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	templatePath := filepath.Join(tmpDir, "template.json")
	if err := os.WriteFile(templatePath, []byte(templateContent), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := &Config{
		Template: templatePath,
		Output:   tmpDir,
		Prefix:   "$",
	}

	if err := Build(cfg); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		variant string
		accent  string
		name    string
		mood    string
	}{
		{"rose-pine.json", "#ebbcba", "Main", "Dark"},
		{"rose-pine-moon.json", "#c4a7e7", "Moon", "Dim"},
		{"rose-pine-dawn.json", "#286983", "Dawn", "Light"},
	}

	for _, tt := range tests {
		t.Run(tt.variant, func(t *testing.T) {
			content, err := os.ReadFile(filepath.Join(tmpDir, tt.variant))
			if err != nil {
				t.Fatal(err)
			}

			var result map[string]interface{}
			if err := json.Unmarshal(content, &result); err != nil {
				t.Fatal(err)
			}

			if got := result["accent"]; got != tt.accent {
				t.Errorf("accent = %v, want %v", got, tt.accent)
			}
			if got := result["name"]; got != tt.name {
				t.Errorf("name = %v, want %v", got, tt.name)
			}
			if got := result["mood"]; got != tt.mood {
				t.Errorf("mood = %v, want %v", got, tt.mood)
			}
		})
	}
}

func TestAccents(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "rose-pine-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	templateContent := `{
        "accent": "$accent"
    }`

	templatePath := filepath.Join(tmpDir, "template.json")
	if err := os.WriteFile(templatePath, []byte(templateContent), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := &Config{
		Template:    templatePath,
		Output:      tmpDir,
		Format:      "hex",
		Prefix:      "$",
		StripSpaces: false,
		Accents:     true,
	}

	if err := Build(cfg); err != nil {
		t.Fatal(err)
	}

	variants := []struct {
		filename string
		accent   string
	}{
		{filename: "rose-pine/rose-pine-foam.json", accent: "#9ccfd8"},
		{filename: "rose-pine/rose-pine-gold.json", accent: "#f6c177"},
		{filename: "rose-pine/rose-pine-iris.json", accent: "#c4a7e7"},
		{filename: "rose-pine/rose-pine-love.json", accent: "#eb6f92"},
		{filename: "rose-pine/rose-pine-pine.json", accent: "#31748f"},
		{filename: "rose-pine/rose-pine-rose.json", accent: "#ebbcba"},

		{filename: "rose-pine-dawn/rose-pine-dawn-foam.json", accent: "#56949f"},
		{filename: "rose-pine-dawn/rose-pine-dawn-gold.json", accent: "#ea9d34"},
		{filename: "rose-pine-dawn/rose-pine-dawn-iris.json", accent: "#907aa9"},
		{filename: "rose-pine-dawn/rose-pine-dawn-love.json", accent: "#b4637a"},
		{filename: "rose-pine-dawn/rose-pine-dawn-pine.json", accent: "#286983"},
		{filename: "rose-pine-dawn/rose-pine-dawn-rose.json", accent: "#d7827e"},

		{filename: "rose-pine-moon/rose-pine-moon-foam.json", accent: "#9ccfd8"},
		{filename: "rose-pine-moon/rose-pine-moon-gold.json", accent: "#f6c177"},
		{filename: "rose-pine-moon/rose-pine-moon-iris.json", accent: "#c4a7e7"},
		{filename: "rose-pine-moon/rose-pine-moon-love.json", accent: "#eb6f92"},
		{filename: "rose-pine-moon/rose-pine-moon-pine.json", accent: "#3e8fb0"},
		{filename: "rose-pine-moon/rose-pine-moon-rose.json", accent: "#ea9a97"},
	}

	for _, v := range variants {
		t.Run(v.filename, func(t *testing.T) {
			content, err := os.ReadFile(filepath.Join(tmpDir, v.filename))
			if err != nil {
				t.Fatalf("Failed to read generated file: %v", err)
			}

			var result map[string]interface{}
			if err := json.Unmarshal(content, &result); err != nil {
				t.Fatalf("Failed to parse JSON: %v", err)
			}

			tests := []struct {
				field string
				want  string
			}{
				{"accent", v.accent},
			}

			for _, tt := range tests {
				if got := result[tt.field]; got != tt.want {
					t.Errorf("%s = %v, want %v", tt.field, got, tt.want)
				}
			}
		})
	}
}

func TestAccentNames(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "rose-pine-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	templateContent := `{
        "accentname": "$accentname"
    }`

	templatePath := filepath.Join(tmpDir, "template.json")
	if err := os.WriteFile(templatePath, []byte(templateContent), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := &Config{
		Template:    templatePath,
		Output:      tmpDir,
		Format:      "hex",
		Prefix:      "$",
		StripSpaces: false,
		Accents:     true,
	}

	if err := Build(cfg); err != nil {
		t.Fatal(err)
	}

	variants := []struct {
		filename   string
		accentname string
	}{
		{filename: "rose-pine/rose-pine-foam.json", accentname: "foam"},
		{filename: "rose-pine/rose-pine-gold.json", accentname: "gold"},
		{filename: "rose-pine/rose-pine-iris.json", accentname: "iris"},
		{filename: "rose-pine/rose-pine-love.json", accentname: "love"},
		{filename: "rose-pine/rose-pine-pine.json", accentname: "pine"},
		{filename: "rose-pine/rose-pine-rose.json", accentname: "rose"},

		{filename: "rose-pine-dawn/rose-pine-dawn-foam.json", accentname: "foam"},
		{filename: "rose-pine-dawn/rose-pine-dawn-gold.json", accentname: "gold"},
		{filename: "rose-pine-dawn/rose-pine-dawn-iris.json", accentname: "iris"},
		{filename: "rose-pine-dawn/rose-pine-dawn-love.json", accentname: "love"},
		{filename: "rose-pine-dawn/rose-pine-dawn-pine.json", accentname: "pine"},
		{filename: "rose-pine-dawn/rose-pine-dawn-rose.json", accentname: "rose"},

		{filename: "rose-pine-moon/rose-pine-moon-foam.json", accentname: "foam"},
		{filename: "rose-pine-moon/rose-pine-moon-gold.json", accentname: "gold"},
		{filename: "rose-pine-moon/rose-pine-moon-iris.json", accentname: "iris"},
		{filename: "rose-pine-moon/rose-pine-moon-love.json", accentname: "love"},
		{filename: "rose-pine-moon/rose-pine-moon-pine.json", accentname: "pine"},
		{filename: "rose-pine-moon/rose-pine-moon-rose.json", accentname: "rose"},
	}

	for _, v := range variants {
		t.Run(v.filename, func(t *testing.T) {
			content, err := os.ReadFile(filepath.Join(tmpDir, v.filename))
			if err != nil {
				t.Fatalf("Failed to read generated file: %v", err)
			}

			var result map[string]interface{}
			if err := json.Unmarshal(content, &result); err != nil {
				t.Fatalf("Failed to parse JSON: %v", err)
			}

			tests := []struct {
				field string
				want  string
			}{
				{"accentname", v.accentname},
			}

			for _, tt := range tests {
				if got := result[tt.field]; got != tt.want {
					t.Errorf("%s = %v, want %v", tt.field, got, tt.want)
				}
			}
		})
	}
}

func TestDirectories(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "rose-pine-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	templateContent := `{
        "id": "$id",
        "name": "$name",
        "description": "$description",
        "type": "$type",
        "colors": {
            "base": "$base",
            "surface": "$surface",
            "love": "$love"
        },
        "custom": "$(main|moon|dawn)"
    }`

	os.Mkdir(filepath.Join(tmpDir, "template"), 0755)
	templatePath := filepath.Join(tmpDir, "template/template.json")
	template2Path := filepath.Join(tmpDir, "template/template2.json")
	if err := os.WriteFile(templatePath, []byte(templateContent), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(template2Path, []byte(templateContent), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := &Config{
		Template:    filepath.Join(tmpDir, "template"),
		Output:      tmpDir,
		Format:      "hex",
		Prefix:      "$",
		StripSpaces: false,
	}

	if err := Build(cfg); err != nil {
		t.Fatal(err)
	}

	variants := []struct {
		filename string
		id       string
		name     string
		varType  string
		baseHex  string
		custom   string
	}{
		{
			filename: "rose-pine/template.json",
			id:       "rose-pine",
			name:     "Rosé Pine",
			varType:  "dark",
			baseHex:  "#191724",
			custom:   "main",
		},
		{
			filename: "rose-pine/template2.json",
			id:       "rose-pine",
			name:     "Rosé Pine",
			varType:  "dark",
			baseHex:  "#191724",
			custom:   "main",
		},
		{
			filename: "rose-pine-moon/template.json",
			id:       "rose-pine-moon",
			name:     "Rosé Pine Moon",
			varType:  "dark",
			baseHex:  "#232136",
			custom:   "moon",
		},
		{
			filename: "rose-pine-moon/template2.json",
			id:       "rose-pine-moon",
			name:     "Rosé Pine Moon",
			varType:  "dark",
			baseHex:  "#232136",
			custom:   "moon",
		},
		{
			filename: "rose-pine-dawn/template.json",
			id:       "rose-pine-dawn",
			name:     "Rosé Pine Dawn",
			varType:  "light",
			baseHex:  "#faf4ed",
			custom:   "dawn",
		},
		{
			filename: "rose-pine-dawn/template2.json",
			id:       "rose-pine-dawn",
			name:     "Rosé Pine Dawn",
			varType:  "light",
			baseHex:  "#faf4ed",
			custom:   "dawn",
		},
	}

	for _, v := range variants {
		t.Run(v.filename, func(t *testing.T) {
			content, err := os.ReadFile(filepath.Join(tmpDir, v.filename))
			if err != nil {
				t.Fatalf("Failed to read generated file: %v", err)
			}

			var result map[string]interface{}
			if err := json.Unmarshal(content, &result); err != nil {
				t.Fatalf("Failed to parse JSON: %v", err)
			}

			tests := []struct {
				field string
				want  string
			}{
				{"id", v.id},
				{"name", v.name},
				{"type", v.varType},
				{"custom", v.custom},
			}

			for _, tt := range tests {
				if got := result[tt.field]; got != tt.want {
					t.Errorf("%s = %v, want %v", tt.field, got, tt.want)
				}
			}

			colors := result["colors"].(map[string]interface{})
			if got := colors["base"]; got != v.baseHex {
				t.Errorf("base color = %v, want %v", got, v.baseHex)
			}
		})
	}
}
