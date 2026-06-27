package derive

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rose-pine/rose-pine-bloom/color"
)

var testOpts = DeriveOpts{
	Input:   "",
	Output:  "",
	Variant: "moon",
	Prefix:  "$",
	Format:  "hex",
	Plain:   false,
	Commas:  true,
	Spaces:  true,
}

func setupTest(t *testing.T) string {
	tmpDir, err := os.MkdirTemp("", "rose-pine-test-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Fatal(err)
		}
	})
	return tmpDir
}

func BenchmarkDetectFormatOptions(b *testing.B) {
	content := `{"base": "#191724", "love": "#eb6f92"}`
	for b.Loop() {
		detectFormatOptions(content, color.MainVariantMeta)
	}
}


func TestDetectFormatOptions(t *testing.T) {
	colors := func(base, love string) string {
		return fmt.Sprintf(`{"base": "%s", "love": "%s"}`, base, love)
	}

	tests := []struct {
		name       string
		content    string
		wantFormat color.ColorFormat
		wantPlain  bool
		wantCommas bool
		wantSpaces bool
	}{
		{
			name:       "hex with hash",
			content:    colors("#191724", "#eb6f92"),
			wantFormat: color.FormatHex,
			wantPlain:  false,
			wantCommas: true,
			wantSpaces: true,
		},
		{
			name:       "hex plain",
			content:    colors("191724", "eb6f92"),
			wantFormat: color.FormatHex,
			wantPlain:  true,
			wantCommas: true,
			wantSpaces: true,
		},
		{
			name:       "rgb function",
			content:    colors("rgb(25, 23, 36)", "rgb(235, 111, 146)"),
			wantFormat: color.FormatRGB,
			wantPlain:  false,
			wantCommas: true,
			wantSpaces: true,
		},
		{
			name:       "rgb no spaces",
			content:    colors("rgb(25,23,36)", "rgb(235,111,146)"),
			wantFormat: color.FormatRGB,
			wantPlain:  false,
			wantCommas: true,
			wantSpaces: false,
		},
		{
			name:       "rgb no commas",
			content:    colors("rgb(25 23 36)", "rgb(235 111 146)"),
			wantFormat: color.FormatRGB,
			wantPlain:  false,
			wantCommas: false,
			wantSpaces: true,
		},
		{
			name:       "hsl function",
			content:    colors("hsl(249, 22%, 12%)", "hsl(343, 76%, 68%)"),
			wantFormat: color.FormatHSL,
			wantPlain:  false,
			wantCommas: true,
			wantSpaces: true,
		},
		{
			name:       "hsl no spaces",
			content:    colors("hsl(249,22%,12%)", "hsl(343,76%,68%)"),
			wantFormat: color.FormatHSL,
			wantPlain:  false,
			wantCommas: true,
			wantSpaces: false,
		},
		{
			name:       "hsl css",
			content:    colors("hsl(249deg 22% 12%)", "hsl(343deg 76% 68%)"),
			wantFormat: color.FormatHSLCSS,
			wantPlain:  false,
			wantCommas: true,
			wantSpaces: true,
		},
		{
			name:       "ansi",
			content:    colors("25;23;36", "235;111;146"),
			wantFormat: color.FormatAnsi,
			wantPlain:  false,
			wantCommas: true,
			wantSpaces: true,
		},
		{
			name:       "fallback on no match",
			content:    "no colors here at all",
			wantFormat: color.FormatHex,
			wantPlain:  false,
			wantCommas: true,
			wantSpaces: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFormat, gotPlain, gotCommas, gotSpaces := detectFormatOptions(tt.content, color.MainVariantMeta)
			if gotFormat != tt.wantFormat {
				t.Errorf("format = %q, want %q", gotFormat, tt.wantFormat)
			}
			if gotPlain != tt.wantPlain {
				t.Errorf("plain = %v, want %v", gotPlain, tt.wantPlain)
			}
			if gotCommas != tt.wantCommas {
				t.Errorf("commas = %v, want %v", gotCommas, tt.wantCommas)
			}
			if gotSpaces != tt.wantSpaces {
				t.Errorf("spaces = %v, want %v", gotSpaces, tt.wantSpaces)
			}
		})
	}
}

func TestCreateAutoDetect(t *testing.T) {
	tests := []struct {
		name    string
		variant string
		input   string
		format  string
		wantVal string
	}{
		{
			name:    "auto-detect hex",
			variant: "moon",
			input:   `"base": "#232136"`,
			format:  "",
			wantVal: "$base",
		},
		{
			name:    "auto-detect hsl",
			variant: "moon",
			input:   `"base": "hsl(246, 24%, 17%)"`,
			format:  "",
			wantVal: "$base",
		},
		{
			name:    "auto-detect rgb",
			variant: "moon",
			input:   `"base": "rgb(35, 33, 54)"`,
			format:  "",
			wantVal: "$base",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := setupTest(t)

			fileContent := "{\n  " + tt.input + "\n}"
			filePath := filepath.Join(tmpDir, "input.json")
			if err := os.WriteFile(filePath, []byte(fileContent), 0o644); err != nil {
				t.Fatal(err)
			}

			opts := testOpts
			opts.Output = tmpDir
			opts.Input = filePath
			opts.Format = tt.format
			opts.Variant = tt.variant

			if err := DeriveTemplate(&opts); err != nil {
				t.Fatal(err)
			}

			content, err := os.ReadFile(filepath.Join(tmpDir, "template.json"))
			if err != nil {
				t.Fatalf("Failed to read generated file: %v", err)
			}
			if !strings.Contains(string(content), tt.wantVal) {
				t.Errorf("generated template should contain %q, got:\n%s", tt.wantVal, string(content))
			}
		})
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
	if err := os.WriteFile(filePath, []byte(fileContent), 0o644); err != nil {
		t.Fatal(err)
	}

	opts := testOpts
	opts.Output = tmpDir
	opts.Input = filePath

	if err := DeriveTemplate(&opts); err != nil {
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
