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

