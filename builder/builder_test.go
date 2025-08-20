package builder

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/rose-pine/rose-pine-bloom/color"
	"github.com/rose-pine/rose-pine-bloom/config"
)

func setupTest(t *testing.T) string {
	tmpDir, err := os.MkdirTemp("", "rose-pine-test-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		err := os.RemoveAll(tmpDir)
		if err != nil {
			t.Fatal(err)
		}
	})
	return tmpDir
}

func buildFromTemplate(t *testing.T, template string, cfg *config.BuildConfig) {
	templatePath := filepath.Join(cfg.Output, "template.json")
	if err := os.WriteFile(templatePath, []byte(template), 0644); err != nil {
		t.Fatal(err)
	}
	cfg.Template = templatePath

	if err := Build(cfg); err != nil {
		t.Fatal(err)
	}
}

func readAndParseJSON(t *testing.T, path string) map[string]any {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	var result map[string]any
	if err := json.Unmarshal(content, &result); err != nil {
		t.Fatal(err)
	}
	return result
}

func assertJSONField(t *testing.T, result map[string]any, field, want string) {
	t.Helper()
	if got := result[field]; got != want {
		t.Errorf("%s = %v, want %v", field, got, want)
	}
}

// testConfig provides standard config
var testConfig = config.BuildConfig{
	Template: "",
	Output:   "",
	Prefix:   "$",
	Format:   "hex",
	Plain:    false,
	Commas:   true,
	Spaces:   true,
}

var testBuildTemplateConfig = config.BuildTemplateConfig{
	Input:   "",
	Output:  "",
	Prefix:  "$",
	Format:  "hex",
	Variant: "moon",
	Plain:   false,
	Commas:  true,
	Spaces:  true,
}

// testColor provides a standard color
var testColor = &color.Color{
	HSL: color.HSL{H: 2, S: 55, L: 83},
	RGB: color.RGB{R: 235, G: 188, B: 186},
}

// testTemplate provides a standard template
const testTemplate = `{
    "id": "$id",
    "name": "$name",
    "appearance": "$appearance",
    "description": "$description",
    "type": "$type",
    "colors": {
        "base": "$base",
        "surface": "$surface",
        "love": "$love"
    },
    "custom": "$(main|moon|dawn)"
}`

var testVariants = []struct {
	filename   string
	id         string
	name       string
	appearance string
	baseHex    string
	custom     string
}{
	{
		filename:   "rose-pine.json",
		id:         "rose-pine",
		name:       "Rosé Pine",
		appearance: "dark",
		baseHex:    "#191724",
		custom:     "main",
	},
	{
		filename:   "rose-pine-moon.json",
		id:         "rose-pine-moon",
		name:       "Rosé Pine Moon",
		appearance: "dark",
		baseHex:    "#232136",
		custom:     "moon",
	},
	{
		filename:   "rose-pine-dawn.json",
		id:         "rose-pine-dawn",
		name:       "Rosé Pine Dawn",
		appearance: "light",
		baseHex:    "#faf4ed",
		custom:     "dawn",
	},
}

func TestColorFormatting(t *testing.T) {

	tests := []struct {
		name   string
		format color.ColorFormat
		plain  bool
		commas bool
		spaces bool
		want   string
	}{
		{"hex", color.FormatHex, false, true, true, "#ebbcba"},
		{"hex plain", color.FormatHex, true, true, true, "ebbcba"},

		{"hsl", color.FormatHSL, false, true, true, "hsl(2, 55%, 83%)"},
		{"hsl no-commas", color.FormatHSL, false, false, true, "hsl(2 55% 83%)"},
		{"hsl no-spaces", color.FormatHSL, false, true, false, "hsl(2,55%,83%)"},
		{"hsl plain", color.FormatHSL, true, true, true, "2, 55%, 83%"},
		{"hsl plain no-commas", color.FormatHSL, true, false, true, "2 55% 83%"},
		{"hsl plain no-spaces", color.FormatHSL, true, true, false, "2,55%,83%"},

		{"hsl-array", color.FormatHSLArray, false, true, true, "[2, 0.55, 0.83]"},
		{"hsl-array plain", color.FormatHSLArray, true, true, true, "2, 0.55, 0.83"},
		{"hsl-array no-commas", color.FormatHSLArray, false, false, true, "[2 0.55 0.83]"},
		{"hsl-array no-spaces", color.FormatHSLArray, false, true, false, "[2,0.55,0.83]"},

		{"hsl-css", color.FormatHSLCSS, false, true, true, "hsl(2deg 55% 83%)"},
		{"hsl-css plain", color.FormatHSLCSS, true, true, true, "2deg 55% 83%"},

		{"rgb", color.FormatRGB, false, true, true, "rgb(235, 188, 186)"},
		{"rgb no-commas", color.FormatRGB, false, false, true, "rgb(235 188 186)"},
		{"rgb no-spaces", color.FormatRGB, false, true, false, "rgb(235,188,186)"},
		{"rgb plain", color.FormatRGB, true, true, true, "235, 188, 186"},
		{"rgb plain no-commas", color.FormatRGB, true, false, true, "235 188 186"},
		{"rgb plain no-spaces", color.FormatRGB, true, true, false, "235,188,186"},

		{"rgb-array", color.FormatRGBArray, false, true, true, "[235, 188, 186]"},
		{"rgb-array plain", color.FormatRGBArray, true, true, true, "235, 188, 186"},
		{"rgb-array no-commas", color.FormatRGBArray, false, false, true, "[235 188 186]"},
		{"rgb-array no-spaces", color.FormatRGBArray, false, true, false, "[235,188,186]"},

		{"rgb-css", color.FormatRGBCSS, false, true, true, "rgb(235 188 186)"},
		{"rgb-css plain", color.FormatRGBCSS, true, true, true, "235 188 186"},

		{"ansi", color.FormatAnsi, false, true, true, "235;188;186"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := color.FormatColor(testColor, tt.format, tt.plain, tt.commas, tt.spaces)
			if got != tt.want {
				t.Errorf("formatColor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlphaFormatting(t *testing.T) {
	alpha := 0.5
	c := &color.Color{
		HSL:   testColor.HSL,
		RGB:   testColor.RGB,
		Alpha: &alpha,
	}

	tests := []struct {
		name   string
		format color.ColorFormat
		plain  bool
		want   string
	}{
		{"hex", color.FormatHex, false, "#ebbcba80"},
		{"hex plain", color.FormatHex, true, "ebbcba80"},

		{"hsl", color.FormatHSL, false, "hsla(2, 55%, 83%, 0.5)"},
		{"hsl plain", color.FormatHSL, true, "2, 55%, 83%, 0.5"},

		{"hsl-css", color.FormatHSLCSS, false, "hsl(2deg 55% 83% / 0.5)"},
		{"hsl-css plain", color.FormatHSLCSS, true, "2deg 55% 83% / 0.5"},

		{"hsl-array", color.FormatHSLArray, false, "[2, 0.55, 0.83, 0.5]"},
		{"hsl-array plain", color.FormatHSLArray, true, "2, 0.55, 0.83, 0.5"},

		{"rgb", color.FormatRGB, false, "rgba(235, 188, 186, 0.5)"},
		{"rgb plain", color.FormatRGB, true, "235, 188, 186, 0.5"},

		{"rgb-css", color.FormatRGBCSS, false, "rgb(235 188 186 / 0.5)"},
		{"rgb-css plain", color.FormatRGBCSS, true, "235 188 186 / 0.5"},

		{"rgb-array", color.FormatRGBArray, false, "[235, 188, 186, 0.5]"},
		{"rgb-array plain", color.FormatRGBArray, true, "235, 188, 186, 0.5"},

		{"ansi", color.FormatAnsi, false, "235;188;186;0.5"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := color.FormatColor(c, tt.format, tt.plain, true, true)
			if got != tt.want {
				t.Errorf("formatColor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlphaVariableFormats(t *testing.T) {
	tmpDir := setupTest(t)

	templateContent := `{
        "base": "$base",
        "base25": "$base/25",
        "base50": "$base/50"
    }`

	tests := []struct {
		format string
		want   map[string]string
	}{
		{
			format: "hex",
			want: map[string]string{
				"base":   "#191724",
				"base25": "#19172440",
				"base50": "#19172480",
			},
		},
		{
			format: "hsl",
			want: map[string]string{
				"base":   "hsl(249, 22%, 12%)",
				"base25": "hsla(249, 22%, 12%, 0.25)",
				"base50": "hsla(249, 22%, 12%, 0.5)",
			},
		},
		{
			format: "hsl-css",
			want: map[string]string{
				"base":   "hsl(249deg 22% 12%)",
				"base25": "hsl(249deg 22% 12% / 0.25)",
				"base50": "hsl(249deg 22% 12% / 0.5)",
			},
		},
		{
			format: "rgb",
			want: map[string]string{
				"base":   "rgb(25, 23, 36)",
				"base25": "rgba(25, 23, 36, 0.25)",
				"base50": "rgba(25, 23, 36, 0.5)",
			},
		},
		{
			format: "rgb-css",
			want: map[string]string{
				"base":   "rgb(25 23 36)",
				"base25": "rgb(25 23 36 / 0.25)",
				"base50": "rgb(25 23 36 / 0.5)",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.format, func(t *testing.T) {
			cfg := testConfig
			cfg.Output = tmpDir
			cfg.Format = tt.format

			buildFromTemplate(t, templateContent, &cfg)

			result := readAndParseJSON(t, filepath.Join(tmpDir, "rose-pine.json"))

			for field, want := range tt.want {
				assertJSONField(t, result, field, want)
			}
		})
	}
}

func TestVariantGeneration(t *testing.T) {
	tmpDir := setupTest(t)

	cfg := testConfig
	cfg.Output = tmpDir
	cfg.Spaces = false

	buildFromTemplate(t, testTemplate, &cfg)

	for _, v := range testVariants {
		t.Run(v.filename, func(t *testing.T) {
			result := readAndParseJSON(t, filepath.Join(tmpDir, v.filename))

			assertJSONField(t, result, "id", v.id)
			assertJSONField(t, result, "name", v.name)
			assertJSONField(t, result, "appearance", v.appearance)
			assertJSONField(t, result, "custom", v.custom)

			colors := result["colors"].(map[string]any)
			if got := colors["base"]; got != v.baseHex {
				t.Errorf("base color = %v, want %v", got, v.baseHex)
			}
		})
	}
}

func TestVariantSpecificValues(t *testing.T) {
	tmpDir := setupTest(t)

	templateContent := `{
        "accent": "$(#ebbcba|#c4a7e7|#286983)",
        "name": "$(Main|Moon|Dawn)",
        "mood": "$(Dark|Dim|Light)"
    }`

	cfg := testConfig
	cfg.Output = tmpDir

	buildFromTemplate(t, templateContent, &cfg)

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
			result := readAndParseJSON(t, filepath.Join(tmpDir, tt.variant))

			assertJSONField(t, result, "accent", tt.accent)
			assertJSONField(t, result, "name", tt.name)
			assertJSONField(t, result, "mood", tt.mood)
		})
	}
}

func TestAccents(t *testing.T) {
	tmpDir := setupTest(t)

	templateContent := `{
        "accentname": "$accentname",
        "accent": "$accent",
		"onaccent": "$onaccent"
    }`

	cfg := testConfig
	cfg.Output = tmpDir

	buildFromTemplate(t, templateContent, &cfg)

	tests := []struct {
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

	for _, v := range tests {
		t.Run(v.filename, func(t *testing.T) {
			result := readAndParseJSON(t, filepath.Join(tmpDir, v.filename))

			assertJSONField(t, result, "accentname", v.accentname)
			assertJSONField(t, result, "accent", v.accent)
			assertJSONField(t, result, "onaccent", v.onaccent)
		})
	}
}

func TestDirectories(t *testing.T) {
	tmpDir := setupTest(t)

	templateDir := filepath.Join(tmpDir, "template")
	err := os.Mkdir(templateDir, 0755)
	if err != nil {
		t.Fatal(err)
	}
	templatePath := filepath.Join(templateDir, "template.json")
	template2Path := filepath.Join(templateDir, "template2.json")
	if err := os.WriteFile(templatePath, []byte(testTemplate), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(template2Path, []byte(testTemplate), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := testConfig
	cfg.Output = tmpDir
	cfg.Template = templateDir

	if err := Build(&cfg); err != nil {
		t.Fatal(err)
	}

	testFiles := []string{"template.json", "template2.json"}
	for _, variant := range testVariants {
		for _, file := range testFiles {
			filename := filepath.Join(variant.id, file)
			t.Run(filename, func(t *testing.T) {
				result := readAndParseJSON(t, filepath.Join(tmpDir, filename))

				assertJSONField(t, result, "id", variant.id)
				assertJSONField(t, result, "name", variant.name)
				assertJSONField(t, result, "custom", variant.custom)

				colors := result["colors"].(map[string]any)
				if got := colors["base"]; got != variant.baseHex {
					t.Errorf("base color = %v, want %v", got, variant.baseHex)
				}
			})
		}
	}
}

func TestCreate(t *testing.T) {
	tmpDir := setupTest(t)

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
  "main-id": "rose-pine",
  "id": "rose-pine-moon",
  "name": "Rosé Pine Moon",
  "description": "All natural pine, faux fur and a bit of soho vibes for the classy minimalist",
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
  "main-id": "rose-pine",
  "id": "$id",
  "name": "$name",
  "description": "$description",
  "dawn-name": "Rosé Pine Dawn"
}`

	filePath := filepath.Join(tmpDir, "input.json")
	if err := os.WriteFile(filePath, []byte(fileContent), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := testBuildTemplateConfig
	cfg.Output = tmpDir
	cfg.Input = filePath

	if err := BuildTemplate(&cfg); err != nil {
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
