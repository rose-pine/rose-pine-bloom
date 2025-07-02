package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func setupTest(t *testing.T) (string, func()) {
	tmpDir, err := os.MkdirTemp("", "rose-pine-test-*")
	if err != nil {
		t.Fatal(err)
	}
	return tmpDir, func() { os.RemoveAll(tmpDir) }
}

func buildFromTemplate(t *testing.T, template string, cfg *Config) {
	templatePath := filepath.Join(cfg.Output, "template.json")
	if err := os.WriteFile(templatePath, []byte(template), 0644); err != nil {
		t.Fatal(err)
	}
	cfg.Template = templatePath

	if err := Build(cfg); err != nil {
		t.Fatal(err)
	}
}

func TestColorFormatting(t *testing.T) {
	color := &Color{
		Hex: "ebbcba",
		HSL: [3]float64{2, 55, 83},
		RGB: [3]int{235, 188, 186},
	}

	tests := []struct {
		name   string
		format ColorFormat
		plain  bool
		commas bool
		spaces bool
		want   string
	}{
		// Hex formats
		{"hex", FormatHex, false, true, true, "#ebbcba"},
		{"hex plain", FormatHex, true, true, true, "ebbcba"},

		// HSL formats
		{"hsl", FormatHSL, false, true, true, "2, 55%, 83%"},
		{"hsl no-commas", FormatHSL, false, false, true, "2 55% 83%"},
		{"hsl no-spaces", FormatHSL, false, true, false, "2,55%,83%"},
		{"hsl-array", FormatHSLArray, false, true, true, "[2, 55%, 83%]"},
		{"hsl-array plain", FormatHSLArray, true, true, true, "2, 55%, 83%"},
		{"hsl-array no-commas", FormatHSLArray, false, false, true, "[2 55% 83%]"},
		{"hsl-array no-spaces", FormatHSLArray, false, true, false, "[2,55%,83%]"},
		{"hsl-css", FormatHSLCSS, false, true, true, "hsl(2deg 55% 83%)"},
		{"hsl-css plain", FormatHSLCSS, true, true, true, "2deg 55% 83%"},
		{"hsl-function", FormatHSLFunc, false, true, true, "hsl(2, 55%, 83%)"},
		{"hsl-function plain", FormatHSLFunc, true, true, true, "2, 55%, 83%"},
		{"hsl-function no-commas", FormatHSLFunc, false, false, true, "hsl(2 55% 83%)"},
		{"hsl-function no-spaces", FormatHSLFunc, false, true, false, "hsl(2,55%,83%)"},

		// RGB formats
		{"rgb", FormatRGB, false, true, true, "235, 188, 186"},
		{"rgb no-spaces", FormatRGB, false, true, false, "235,188,186"},
		{"rgb no-commas", FormatRGB, false, false, true, "235 188 186"},
		{"rgb-ansi", FormatRGBAnsi, false, true, true, "235;188;186"},
		{"rgb-array", FormatRGBArray, false, true, true, "[235, 188, 186]"},
		{"rgb-array plain", FormatRGBArray, true, true, true, "235, 188, 186"},
		{"rgb-array no-spaces", FormatRGBArray, false, true, false, "[235,188,186]"},
		{"rgb-css", FormatRGBCSS, false, true, true, "rgb(235 188 186)"},
		{"rgb-css plain", FormatRGBCSS, true, true, true, "235 188 186"},
		{"rgb-function", FormatRGBFunc, false, true, true, "rgb(235, 188, 186)"},
		{"rgb-function plain", FormatRGBFunc, true, true, true, "235, 188, 186"},
		{"rgb-function no-commas", FormatRGBFunc, false, false, true, "rgb(235 188 186)"},
		{"rgb-function no-spaces", FormatRGBFunc, false, true, false, "rgb(235,188,186)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatColor(color, tt.format, tt.plain, tt.commas, tt.spaces)
			if got != tt.want {
				t.Errorf("formatColor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlphaFormatting(t *testing.T) {
	alpha := 0.5
	color := &Color{
		Hex:   "ebbcba",
		HSL:   [3]float64{2, 55, 83},
		RGB:   [3]int{235, 188, 186},
		Alpha: &alpha,
	}

	tests := []struct {
		format ColorFormat
		want   string
	}{
		{FormatHex, "#ebbcba80"},

		{FormatHSL, "2, 55%, 83%, 0.5"},
		{FormatHSLArray, "[2, 55%, 83%, 0.5]"},
		{FormatHSLCSS, "hsl(2deg 55% 83% / 0.5)"},
		{FormatHSLFunc, "hsla(2, 55%, 83%, 0.5)"},

		{FormatRGB, "235, 188, 186, 0.5"},
		{FormatRGBAnsi, "235;188;186;0.5"},
		{FormatRGBArray, "[235, 188, 186, 0.5]"},
		{FormatRGBCSS, "rgb(235 188 186 / 0.5)"},
		{FormatRGBFunc, "rgba(235, 188, 186, 0.5)"},
	}

	for _, tt := range tests {
		name := string(tt.format)
		t.Run(name, func(t *testing.T) {
			got := formatColor(color, tt.format, false, true, true)
			if got != tt.want {
				t.Errorf("formatColor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlphaVariables(t *testing.T) {
	tmpDir, cleanup := setupTest(t)
	defer cleanup()

	templateContent := `{
        "regular": "$base",
        "alpha50": "$base/50",
        "alpha25": "$base/25",
        "alphaMuted50": "$muted/50"
    }`

	cfg := &Config{
		Output: tmpDir,
		Format: "rgb",
		Prefix: "$",
		Commas: true,
		Spaces: true,
	}

	buildFromTemplate(t, templateContent, cfg)

	content, err := os.ReadFile(filepath.Join(tmpDir, "rose-pine.json"))
	if err != nil {
		t.Fatal(err)
	}

	var result map[string]any
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
	tmpDir, cleanup := setupTest(t)
	defer cleanup()

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

	cfg := &Config{
		Output: tmpDir,
		Format: "hex",
		Prefix: "$",
		Spaces: false,
	}

	buildFromTemplate(t, templateContent, cfg)

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

			var result map[string]any
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

			colors := result["colors"].(map[string]any)
			if got := colors["base"]; got != v.baseHex {
				t.Errorf("base color = %v, want %v", got, v.baseHex)
			}
		})
	}
}

func TestVariantSpecificValues(t *testing.T) {
	tmpDir, cleanup := setupTest(t)
	defer cleanup()

	templateContent := `{
        "accent": "$(#ebbcba|#c4a7e7|#286983)",
        "name": "$(Main|Moon|Dawn)",
        "mood": "$(Dark|Dim|Light)"
    }`

	cfg := &Config{
		Output: tmpDir,
		Prefix: "$",
	}

	buildFromTemplate(t, templateContent, cfg)

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

			var result map[string]any
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
	tmpDir, cleanup := setupTest(t)
	defer cleanup()

	templateContent := `{
        "accentname": "$accentname",
        "accent": "$accent",
		"onaccent": "$onaccent"
    }`

	cfg := &Config{
		Output: tmpDir,
		Format: "hex",
		Prefix: "$",
		Spaces: true,
	}

	buildFromTemplate(t, templateContent, cfg)

	variants := []struct {
		filename   string
		accentname string
		accent     string
		onaccent   string
	}{
		{filename: "rose-pine/rose-pine-love.json", accentname: "love", accent: "#eb6f92", onaccent: "#e0def4"},
		{filename: "rose-pine/rose-pine-gold.json", accentname: "gold", accent: "#f6c177", onaccent: "#1f1d2e"},
		{filename: "rose-pine/rose-pine-rose.json", accentname: "rose", accent: "#ebbcba", onaccent: "#1f1d2e"},
		{filename: "rose-pine/rose-pine-pine.json", accentname: "pine", accent: "#31748f", onaccent: "#e0def4"},
		{filename: "rose-pine/rose-pine-foam.json", accentname: "foam", accent: "#9ccfd8", onaccent: "#1f1d2e"},
		{filename: "rose-pine/rose-pine-iris.json", accentname: "iris", accent: "#c4a7e7", onaccent: "#1f1d2e"},

		{filename: "rose-pine-moon/rose-pine-moon-love.json", accentname: "love", accent: "#eb6f92", onaccent: "#e0def4"},
		{filename: "rose-pine-moon/rose-pine-moon-gold.json", accentname: "gold", accent: "#f6c177", onaccent: "#2a273f"},
		{filename: "rose-pine-moon/rose-pine-moon-rose.json", accentname: "rose", accent: "#ea9a97", onaccent: "#2a273f"},
		{filename: "rose-pine-moon/rose-pine-moon-pine.json", accentname: "pine", accent: "#3e8fb0", onaccent: "#e0def4"},
		{filename: "rose-pine-moon/rose-pine-moon-foam.json", accentname: "foam", accent: "#9ccfd8", onaccent: "#2a273f"},
		{filename: "rose-pine-moon/rose-pine-moon-iris.json", accentname: "iris", accent: "#c4a7e7", onaccent: "#2a273f"},

		{filename: "rose-pine-dawn/rose-pine-dawn-love.json", accentname: "love", accent: "#b4637a", onaccent: "#fffaf3"},
		{filename: "rose-pine-dawn/rose-pine-dawn-gold.json", accentname: "gold", accent: "#ea9d34", onaccent: "#fffaf3"},
		{filename: "rose-pine-dawn/rose-pine-dawn-rose.json", accentname: "rose", accent: "#d7827e", onaccent: "#fffaf3"},
		{filename: "rose-pine-dawn/rose-pine-dawn-pine.json", accentname: "pine", accent: "#286983", onaccent: "#fffaf3"},
		{filename: "rose-pine-dawn/rose-pine-dawn-foam.json", accentname: "foam", accent: "#56949f", onaccent: "#fffaf3"},
		{filename: "rose-pine-dawn/rose-pine-dawn-iris.json", accentname: "iris", accent: "#907aa9", onaccent: "#fffaf3"},
	}

	for _, v := range variants {
		t.Run(v.filename, func(t *testing.T) {
			content, err := os.ReadFile(filepath.Join(tmpDir, v.filename))
			if err != nil {
				t.Fatalf("Failed to read generated file: %v", err)
			}

			var result map[string]any
			if err := json.Unmarshal(content, &result); err != nil {
				t.Fatalf("Failed to parse JSON: %v", err)
			}

			tests := []struct {
				field string
				want  string
			}{
				{"accentname", v.accentname},
				{"accent", v.accent},
				{"onaccent", v.onaccent},
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
	tmpDir, cleanup := setupTest(t)
	defer cleanup()

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

	templateDir := filepath.Join(tmpDir, "template")
	os.Mkdir(templateDir, 0755)
	templatePath := filepath.Join(templateDir, "template.json")
	template2Path := filepath.Join(templateDir, "template2.json")
	if err := os.WriteFile(templatePath, []byte(templateContent), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(template2Path, []byte(templateContent), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := &Config{
		Template: templateDir,
		Output:   tmpDir,
		Format:   "hex",
		Prefix:   "$",
		Spaces:   true,
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

			var result map[string]any
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

			colors := result["colors"].(map[string]any)
			if got := colors["base"]; got != v.baseHex {
				t.Errorf("base color = %v, want %v", got, v.baseHex)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	tmpDir, cleanup := setupTest(t)
	defer cleanup()

	fileContent := `{
  "base": "#232136",
  "surface": "#2a273f",
  "overlay": "#393552",
  "muted": "#6e6a86",
  "subtle": "#908caa",
  "text": "#e0def4",
  "love": "#eb6f92",
  "gold": "#f6c177",
  "rose": "#ea9a97",
  "pine": "#3e8fb0",
  "foam": "#9ccfd8",
  "iris": "#c4a7e7",
  "id": "rose-pine-moon",
  "name": "Rosé Pine Moon",
  "description": "All natural pine, faux fur and a bit of soho vibes for the classy minimalist",
  "regular-id": "rose-pine",
  "dawn-name": "Rosé Pine Dawn"
}`
	expected := `{
  "base": "$base",
  "surface": "$surface",
  "overlay": "$overlay",
  "muted": "$muted",
  "subtle": "$subtle",
  "text": "$text",
  "love": "$love",
  "gold": "$gold",
  "rose": "$rose",
  "pine": "$pine",
  "foam": "$foam",
  "iris": "$iris",
  "id": "$id",
  "name": "$name",
  "description": "$description",
  "regular-id": "rose-pine",
  "dawn-name": "Rosé Pine Dawn"
}`

	filePath := filepath.Join(tmpDir, "template.json")
	if err := os.WriteFile(filePath, []byte(fileContent), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := &Config{
		Template: filePath,
		Create:   "moon",
		Output:   tmpDir,
		Format:   "hex",
		Prefix:   "$",
		Spaces:   true,
	}

	if err := Build(cfg); err != nil {
		t.Fatal(err)
	}

	t.Run("template.json", func(t *testing.T) {
		content, err := os.ReadFile(filepath.Join(tmpDir, "template.json"))
		if err != nil {
			t.Fatalf("Failed to read generated file: %v", err)
		}
		if string(content) != expected {
			t.Errorf("want %s\n\n got %s", expected, string(content))
		}
	})
}
