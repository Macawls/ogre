package font

import (
	"strings"
	"testing"
)

func BenchmarkWrapText(b *testing.B) {
	mgr := NewManager()
	if err := mgr.LoadDefaults(); err != nil {
		b.Fatalf("load defaults: %v", err)
	}
	face := mgr.Resolve("sans-serif", 400, "normal")
	ff, err := mgr.NewFace(face, 16)
	if err != nil {
		b.Fatal(err)
	}
	// ~500 word paragraph.
	text := strings.Repeat("The quick brown fox jumps over the lazy dog. ", 55)
	cfg := WrapConfig{
		MaxWidth:   800,
		FontFace:   ff,
		FontSize:   16,
		LineHeight: 20,
		WhiteSpace: 0,
		WordBreak:  0,
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		lines := WrapText(text, cfg)
		if len(lines) == 0 {
			b.Fatal("no lines")
		}
	}
}
