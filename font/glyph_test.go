package font

import (
	"strings"
	"testing"

	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
)

func loadTestFont(t *testing.T) *opentype.Font {
	t.Helper()
	f, err := opentype.Parse(goregular.TTF)
	if err != nil {
		t.Fatalf("parse font: %v", err)
	}
	return f
}

func TestGlyphToPath_A(t *testing.T) {
	f := loadTestFont(t)
	gp, err := GlyphToPath(f, 'A', 16)
	if err != nil {
		t.Fatalf("GlyphToPath: %v", err)
	}
	if gp.D == "" {
		t.Fatal("expected non-empty path data")
	}
	if !strings.HasPrefix(gp.D, "M") {
		t.Fatalf("expected path to start with M, got %q", gp.D[:10])
	}
	if gp.Advance <= 0 {
		t.Fatalf("expected positive advance, got %v", gp.Advance)
	}
}

func TestGlyphToPath_Space(t *testing.T) {
	f := loadTestFont(t)
	gp, err := GlyphToPath(f, ' ', 16)
	if err != nil {
		t.Fatalf("GlyphToPath: %v", err)
	}
	if gp.D != "" {
		t.Fatalf("expected empty path for space, got %q", gp.D)
	}
	if gp.Advance <= 0 {
		t.Fatalf("expected positive advance for space, got %v", gp.Advance)
	}
}

func TestTextToPath_Hello(t *testing.T) {
	mgr := NewManager()
	if err := mgr.LoadDefaults(); err != nil {
		t.Fatalf("load defaults: %v", err)
	}

	path, advance := TextToPath(mgr, "Hello", "sans-serif", 400, "normal", 16)
	if path == "" {
		t.Fatal("expected non-empty path for Hello")
	}
	if advance <= 0 {
		t.Fatalf("expected positive advance, got %v", advance)
	}
	if advance < 30 || advance > 60 {
		t.Fatalf("advance %.2f seems unreasonable for 'Hello' at 16px", advance)
	}
}

func TestTextToPath_MissingFamily(t *testing.T) {
	mgr := NewManager()
	path, advance := TextToPath(mgr, "test", "nonexistent", 400, "normal", 16)
	if path != "" || advance != 0 {
		t.Fatal("expected empty result for missing family")
	}
}

func TestGlyphToPath_ContainsZ(t *testing.T) {
	f := loadTestFont(t)
	gp, err := GlyphToPath(f, 'O', 16)
	if err != nil {
		t.Fatalf("GlyphToPath: %v", err)
	}
	if !strings.Contains(gp.D, "Z") {
		t.Fatal("expected path to contain Z (close) command")
	}
}
