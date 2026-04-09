package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/macawls/ogre"
)

func BenchmarkRenderSVG(b *testing.B) {
	fixtures, _ := filepath.Glob("fixtures/*.html")
	for _, f := range fixtures {
		name := filepath.Base(f)
		html, _ := os.ReadFile(f)
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ogre.Render(string(html), ogre.Options{Width: 1200, Height: 630})
			}
		})
	}
}

func BenchmarkRenderPNG(b *testing.B) {
	fixtures, _ := filepath.Glob("fixtures/*.html")
	for _, f := range fixtures {
		name := filepath.Base(f)
		html, _ := os.ReadFile(f)
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ogre.Render(string(html), ogre.Options{Width: 1200, Height: 630, Format: ogre.FormatPNG})
			}
		})
	}
}
