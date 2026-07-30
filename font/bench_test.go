package font

import (
	"testing"
)

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
