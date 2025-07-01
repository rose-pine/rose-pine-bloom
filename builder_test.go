package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
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
		RGB: [3]int{235, 188, 186},
		HSL: [3]float64{2, 55, 83},
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
		{"hsl plain", FormatHSL, true, true, true, "2, 55%, 83%"},
		{"hsl-function", FormatHSLFunction, false, true, true, "hsl(2, 55%, 83%)"},
		{"hsl-function plain", FormatHSLFunction, true, true, true, "hsl(2, 55%, 83%)"},
		{"hsl-css", FormatHSLCSS, false, true, true, "hsl(2deg 55% 83%)"},
		{"hsl-css plain", FormatHSLCSS, true, true, true, "2, 55%, 83%"},
		{"hsl-array", FormatHSLArray, false, true, true, "[2, 0.55, 0.83]"},
		{"hsl-array plain", FormatHSLArray, true, true, true, "2, 0.55, 0.83"},

		// RGB formats
		{"rgb", FormatRGB, false, true, true, "235, 188, 186"},
		{"rgb plain", FormatRGB, true, true, true, "235, 188, 186"},
		{"rgb no-spaces", FormatRGB, true, true, false, "235,188,186"},
		{"rgb no-commas", FormatRGB, true, false, true, "235 188 186"},
		{"rgb-function", FormatRGBFunction, false, true, true, "rgb(235, 188, 186)"},
		{"rgb-function plain", FormatRGBFunction, true, true, true, "rgb(235, 188, 186)"},
		{"rgb-css", FormatRGBCSS, false, true, true, "rgb(235 188 186)"},
		{"rgb-css plain", FormatRGBCSS, true, true, true, "235, 188, 186"},
		{"rgb-array", FormatRGBArray, false, true, true, "[235, 188, 186]"},
		{"rgb-array plain", FormatRGBArray, true, true, true, "235, 188, 186"},

		// ANSI format
		{"ansi", FormatANSI, false, true, true, "235;188;186"},
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

// TestAlphaFormatting tests color formatting with alpha values across all formats
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
		plain  bool
		want   string
	}{
		// Hex formats
		{FormatHex, false, "#ebbcba80"},
		{FormatHex, true, "ebbcba80"},

		// HSL formats
		{FormatHSL, false, "2, 55%, 83%, 0.5"},
		{FormatHSL, true, "2, 55%, 83%, 0.5"},
		{FormatHSLFunction, false, "hsla(2, 55%, 83%, 0.5)"},
		{FormatHSLFunction, true, "hsla(2, 55%, 83%, 0.5)"},
		{FormatHSLCSS, false, "hsl(2deg 55% 83% / 0.5)"},
		{FormatHSLCSS, true, "2, 55%, 83%, 0.5"},
		{FormatHSLArray, false, "[2, 0.55, 0.83, 0.5]"},
		{FormatHSLArray, true, "2, 0.55, 0.83, 0.5"},

		// RGB formats
		{FormatRGB, false, "235, 188, 186, 0.5"},
		{FormatRGB, true, "235, 188, 186, 0.5"},
		{FormatRGBFunction, false, "rgba(235, 188, 186, 0.5)"},
		{FormatRGBFunction, true, "rgba(235, 188, 186, 0.5)"},
		{FormatRGBCSS, false, "rgb(235 188 186 / 0.5)"},
		{FormatRGBCSS, true, "235, 188, 186, 0.5"},
		{FormatRGBArray, false, "[235, 188, 186, 0.5]"},
		{FormatRGBArray, true, "235, 188, 186, 0.5"},

		// ANSI format
		{FormatANSI, false, "235;188;186;0.5"},
	}

	for _, tt := range tests {
		name := string(tt.format)
		if tt.plain {
			name += " plain"
		}
		t.Run(name, func(t *testing.T) {
			got := formatColor(color, tt.format, tt.plain, true, true)
			if got != tt.want {
				t.Errorf("formatColor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBasicGeneration(t *testing.T) {
	tmpDir, cleanup := setupTest(t)
	defer cleanup()

	template := `{
		"name": "$name",
		"id": "$id",
		"type": "$type",
		"base": "$base",
		"text": "$text"
	}`

	cfg := &Config{
		Template: "",
		Output:   tmpDir,
		Format:   "hex",
		Prefix:   "$",
		Plain:    false,
		Commas:   true,
		Spaces:   true,
	}

	buildFromTemplate(t, template, cfg)

	variants := []string{"rose-pine.json", "rose-pine-moon.json", "rose-pine-dawn.json"}
	for _, variant := range variants {
		path := filepath.Join(tmpDir, variant)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("Expected file %s to exist", variant)
		}
	}

	content, err := os.ReadFile(filepath.Join(tmpDir, "rose-pine.json"))
	if err != nil {
		t.Fatal(err)
	}

	var result map[string]any
	if err := json.Unmarshal(content, &result); err != nil {
		t.Fatal(err)
	}

	expected := map[string]string{
		"name": "Rosé Pine",
		"id":   "rose-pine",
		"type": "dark",
		"base": "#191724",
		"text": "#e0def4",
	}

	for key, want := range expected {
		if got := result[key]; got != want {
			t.Errorf("%s = %v, want %v", key, got, want)
		}
	}
}

func TestAccentGeneration(t *testing.T) {
	tmpDir, cleanup := setupTest(t)
	defer cleanup()

	template := `{
		"accent": "$accent",
		"accentname": "$accentname",
		"onaccent": "$onaccent"
	}`

	cfg := &Config{
		Template: "",
		Output:   tmpDir,
		Format:   "hex",
		Prefix:   "$",
		Plain:    false,
		Commas:   true,
		Spaces:   true,
	}

	buildFromTemplate(t, template, cfg)

	accentFile := filepath.Join(tmpDir, "rose-pine", "rose-pine-foam.json")
	if _, err := os.Stat(accentFile); os.IsNotExist(err) {
		t.Errorf("Expected accent file to exist: %s", accentFile)
	}

	content, err := os.ReadFile(accentFile)
	if err != nil {
		t.Fatal(err)
	}

	var result map[string]any
	if err := json.Unmarshal(content, &result); err != nil {
		t.Fatal(err)
	}

	if result["accent"] != "#9ccfd8" {
		t.Errorf("accent = %v, want #9ccfd8", result["accent"])
	}
	if result["accentname"] != "foam" {
		t.Errorf("accentname = %v, want foam", result["accentname"])
	}
	if result["onaccent"] != "#1f1d2e" {
		t.Errorf("onaccent = %v, want #1f1d2e", result["onaccent"])
	}
}

func TestAlphaVariables(t *testing.T) {
	tmpDir, cleanup := setupTest(t)
	defer cleanup()

	template := `{
		"regular": "$base",
		"alpha50": "$base/50",
		"alpha25": "$base/25"
	}`

	cfg := &Config{
		Template: "",
		Output:   tmpDir,
		Format:   "rgb",
		Prefix:   "$",
		Plain:    true,
		Commas:   true,
		Spaces:   true,
	}

	buildFromTemplate(t, template, cfg)

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
	}

	for _, tt := range tests {
		if got := result[tt.field]; got != tt.want {
			t.Errorf("%s = %v, want %v", tt.field, got, tt.want)
		}
	}
}

func TestVariantSpecificValues(t *testing.T) {
	tmpDir, cleanup := setupTest(t)
	defer cleanup()

	template := `{
		"mood": "$(Cozy|Mystical|Fresh)"
	}`

	cfg := &Config{
		Template: "",
		Output:   tmpDir,
		Format:   "hex",
		Prefix:   "$",
		Plain:    false,
		Commas:   true,
		Spaces:   true,
	}

	buildFromTemplate(t, template, cfg)

	tests := []struct {
		file string
		want string
	}{
		{"rose-pine.json", "Cozy"},
		{"rose-pine-moon.json", "Mystical"},
		{"rose-pine-dawn.json", "Fresh"},
	}

	for _, tt := range tests {
		content, err := os.ReadFile(filepath.Join(tmpDir, tt.file))
		if err != nil {
			t.Fatal(err)
		}

		var result map[string]any
		if err := json.Unmarshal(content, &result); err != nil {
			t.Fatal(err)
		}

		if got := result["mood"]; got != tt.want {
			t.Errorf("%s mood = %v, want %v", tt.file, got, tt.want)
		}
	}
}

func TestCreateTemplate(t *testing.T) {
	tmpDir, cleanup := setupTest(t)
	defer cleanup()

	themeContent := `{
		"base": "#232136",
		"surface": "#2a273f",
		"name": "Rosé Pine Moon"
	}`

	themePath := filepath.Join(tmpDir, "theme.json")
	if err := os.WriteFile(themePath, []byte(themeContent), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := &Config{
		Template: themePath,
		Create:   "moon",
		Output:   tmpDir,
		Format:   "hex",
		Prefix:   "$",
		Plain:    false,
		Commas:   true,
		Spaces:   true,
	}

	if err := Build(cfg); err != nil {
		t.Fatal(err)
	}

	templatePath := filepath.Join(tmpDir, "template.json")
	content, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatal(err)
	}

	templateStr := string(content)
	if !strings.Contains(templateStr, "$base") {
		t.Error("Expected template to contain $base variable")
	}
	if !strings.Contains(templateStr, "$surface") {
		t.Error("Expected template to contain $surface variable")
	}
	if !strings.Contains(templateStr, "$name") {
		t.Error("Expected template to contain $name variable")
	}
}

func TestDirectoryProcessing(t *testing.T) {
	tmpDir, cleanup := setupTest(t)
	defer cleanup()

	templateDir := filepath.Join(tmpDir, "templates")
	if err := os.MkdirAll(templateDir, 0755); err != nil {
		t.Fatal(err)
	}

	templates := map[string]string{
		"theme1.json": `{"name": "$name", "base": "$base"}`,
		"theme2.json": `{"id": "$id", "text": "$text"}`,
	}

	for filename, content := range templates {
		path := filepath.Join(templateDir, filename)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	cfg := &Config{
		Template: templateDir,
		Output:   tmpDir,
		Format:   "hex",
		Prefix:   "$",
		Plain:    false,
		Commas:   true,
		Spaces:   true,
	}

	if err := Build(cfg); err != nil {
		t.Fatal(err)
	}

	expectedFiles := []string{
		"rose-pine/theme1.json",
		"rose-pine/theme2.json",
		"rose-pine-moon/theme1.json",
		"rose-pine-moon/theme2.json",
		"rose-pine-dawn/theme1.json",
		"rose-pine-dawn/theme2.json",
	}

	for _, file := range expectedFiles {
		path := filepath.Join(tmpDir, file)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("Expected file to exist: %s", file)
		}
	}
}
