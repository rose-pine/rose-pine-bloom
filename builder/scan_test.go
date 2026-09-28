package builder

import "testing"

func TestScanRoleShade(t *testing.T) {
	opts := ScannerOpts{Prefix: '$'}

	tests := []struct {
		input     string
		wantRole  string
		wantShade int
	}{
		{"$base-50", "base", 50},
		{"$love-500", "love", 500},
		{"$pine-900", "pine", 900},
	}

	for _, tt := range tests {
		captures, err := Scan(tt.input, opts)
		if err != nil {
			t.Errorf("Scan(%q) error: %v", tt.input, err)
			continue
		}
		if len(captures) != 1 {
			t.Errorf("Scan(%q) got %d captures, want 1", tt.input, len(captures))
			continue
		}
		rc, ok := captures[0].(RoleCapture)
		if !ok {
			t.Errorf("Scan(%q) capture is not RoleCapture", tt.input)
			continue
		}
		if rc.role != tt.wantRole {
			t.Errorf("Scan(%q) role = %q, want %q", tt.input, rc.role, tt.wantRole)
		}
		if rc.shade == nil {
			t.Errorf("Scan(%q) shade is nil, want %d", tt.input, tt.wantShade)
		} else if *rc.shade != tt.wantShade {
			t.Errorf("Scan(%q) shade = %d, want %d", tt.input, *rc.shade, tt.wantShade)
		}
	}
}

func TestScanRoleShadeRejectsInvalidValues(t *testing.T) {
	opts := ScannerOpts{Prefix: '$'}

	for _, input := range []string{"$base-42", "$base-000", "$base-950"} {
		if _, err := Scan(input, opts); err == nil {
			t.Errorf("Scan(%q) expected an error, got nil", input)
		}
	}
}

func TestScanRoleShadeWithAlpha(t *testing.T) {
	opts := ScannerOpts{Prefix: '$'}
	captures, err := Scan("$rose-500/80", opts)
	if err != nil {
		t.Fatalf("Scan() error: %v", err)
	}

	rc, ok := captures[0].(RoleCapture)
	if !ok {
		t.Fatalf("capture is not RoleCapture")
	}
	if rc.shade == nil || *rc.shade != 500 {
		t.Fatalf("shade = %v, want 500", rc.shade)
	}
	if rc.alpha == nil || *rc.alpha != 0.8 {
		t.Fatalf("alpha = %v, want 0.8", rc.alpha)
	}
}

func TestScanRoleAlphaRejectsOutOfRangeValues(t *testing.T) {
	opts := ScannerOpts{Prefix: '$'}
	if _, err := Scan("$rose/150", opts); err == nil {
		t.Fatal("Scan() expected an error for alpha > 100")
	}
}

func TestScanVariantCaptureRejectsMalformedInput(t *testing.T) {
	opts := ScannerOpts{Prefix: '$'}

	for _, input := range []string{
		"$(main|moon)",
		"$(main|moon|dawn|extra)",
		"$(main|moon|dawn",
	} {
		if _, err := Scan(input, opts); err == nil {
			t.Errorf("Scan(%q) expected an error, got nil", input)
		}
	}
}
