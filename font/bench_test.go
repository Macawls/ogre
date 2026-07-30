package font

import (
	"os"
	"testing"
)

func readTestFontBench() ([]byte, error) {
	return os.ReadFile("testdata/Estedad-VF.ttf")
}

func BenchmarkTextToPath(b *testing.B) {
	mgr := NewManager()
	if err := mgr.LoadDefaults(); err != nil {
		b.Fatalf("load defaults: %v", err)
	}
	const text = "The quick brown fox jumps over the lazy dog"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		path, adv := TextToPath(mgr, text, "sans-serif", 400, "normal", 16)
		if path == "" || adv <= 0 {
			b.Fatal("empty path")
		}
	}
}

func BenchmarkTranslatePath(b *testing.B) {
	// Representative glyph path shape: mixed M/L/Q/C commands with a Z.
	const d = "M12.34 5.6L45.678 90.1Q100 200 300.5 -12.75C1.5 2.5 3.5 4.5 5.5 6.5Z"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if out := translatePath(d, 100.25, 50.5); len(out) == 0 {
			b.Fatal("empty translation")
		}
	}
}

// BenchmarkShapeTextWarm measures Manager.ShapeText on a warm cache: the
// go-text face is built on the first iteration and reused thereafter.
func BenchmarkShapeTextWarm(b *testing.B) {
	mgr := NewManager()
	if err := mgr.LoadDefaults(); err != nil {
		b.Fatalf("load defaults: %v", err)
	}
	face := mgr.Resolve("sans-serif", 400, "normal")
	if face == nil {
		b.Fatal("no face")
	}
	const text = "The quick brown fox jumps over the lazy dog"
	// Prime the cache so the timed loop is warm-only.
	if _, err := mgr.ShapeText(face, text, 16, false); err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		run, err := mgr.ShapeText(face, text, 16, false)
		if err != nil || len(run.Glyphs) == 0 {
			b.Fatal("empty run")
		}
	}
}

// BenchmarkShapeBytesUncached measures the legacy Shaper.ShapeBytes path,
// which re-parses the font on every call. This is the pre-change baseline.
func BenchmarkTextToPathVariable(b *testing.B) {
	mgr := loadVariableFontBench(b)
	const text = "The quick brown fox jumps over the lazy dog"
	vars := []Variation{{Tag: "wght", Value: 700}}
	if p, _ := TextToPathVariations(mgr, text, "Estedad", 400, "normal", 16, vars); p == "" {
		b.Fatal("empty path (setup)")
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		path, adv := TextToPathVariations(mgr, text, "Estedad", 400, "normal", 16, vars)
		if path == "" || adv <= 0 {
			b.Fatal("empty path")
		}
	}
}

func BenchmarkTextToPathVariableNoVars(b *testing.B) {
	mgr := loadVariableFontBench(b)
	const text = "The quick brown fox jumps over the lazy dog"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		path, adv := TextToPathVariations(mgr, text, "Estedad", 400, "normal", 16, nil)
		if path == "" || adv <= 0 {
			b.Fatal("empty path")
		}
	}
}

func BenchmarkShaperFaceLookup(b *testing.B) {
	mgr := loadVariableFontBench(b)
	face := mgr.Resolve("Estedad", 400, "normal")
	vars := []Variation{{Tag: "wght", Value: 700}}
	if _, err := face.shaperFace(vars); err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, err := face.shaperFace(vars); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkVariationKey(b *testing.B) {
	vars := []Variation{
		{Tag: "wght", Value: 900},
		{Tag: "SOFT", Value: 0},
		{Tag: "WONK", Value: 1},
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if variationKey(vars) == "" {
			b.Fatal("empty key")
		}
	}
}

func loadVariableFontBench(b *testing.B) *Manager {
	b.Helper()
	data, err := readTestFontBench()
	if err != nil {
		b.Fatalf("read test font: %v", err)
	}
	mgr := NewManager()
	if err := mgr.LoadFont(FontSource{Name: "Estedad", Weight: 400, Style: "normal", Data: data}); err != nil {
		b.Fatalf("load font: %v", err)
	}
	return mgr
}

func BenchmarkShapeBytesUncached(b *testing.B) {
	mgr := NewManager()
	if err := mgr.LoadDefaults(); err != nil {
		b.Fatalf("load defaults: %v", err)
	}
	face := mgr.Resolve("sans-serif", 400, "normal")
	if face == nil {
		b.Fatal("no face")
	}
	const text = "The quick brown fox jumps over the lazy dog"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		s := NewShaper()
		run, err := s.ShapeBytes(text, face.RawData, 16, false)
		if err != nil || len(run.Glyphs) == 0 {
			b.Fatal("empty run")
		}
	}
}
