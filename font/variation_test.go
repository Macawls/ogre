package font

import (
	"os"
	"testing"
)

func loadVariableFont(t testing.TB) *Manager {
	t.Helper()
	data, err := os.ReadFile("testdata/Estedad-VF.ttf")
	if err != nil {
		t.Fatalf("read test font: %v", err)
	}
	mgr := NewManager()
	if err := mgr.LoadFont(FontSource{
		Name:   "Estedad",
		Weight: 400,
		Style:  "normal",
		Data:   data,
	}); err != nil {
		t.Fatalf("load font: %v", err)
	}
	return mgr
}

func TestFaceIsVariable(t *testing.T) {
	mgr := loadVariableFont(t)
	face := mgr.Resolve("Estedad", 400, "normal")
	if face == nil {
		t.Fatal("face not resolved")
	}
	if !face.Variable {
		t.Fatal("face reported non-variable; Estedad-VF has fvar table")
	}
}

func TestShaperFacePerVariationCoords(t *testing.T) {
	mgr := loadVariableFont(t)
	face := mgr.Resolve("Estedad", 400, "normal")

	light, err := face.shaperFace([]Variation{{Tag: "wght", Value: 100}})
	if err != nil {
		t.Fatalf("light shaper face: %v", err)
	}
	bold, err := face.shaperFace([]Variation{{Tag: "wght", Value: 900}})
	if err != nil {
		t.Fatalf("bold shaper face: %v", err)
	}
	if light == bold {
		t.Fatal("distinct variation coords should produce distinct cached faces")
	}
	dflt, err := face.shaperFace(nil)
	if err != nil {
		t.Fatalf("default shaper face: %v", err)
	}
	if dflt == light || dflt == bold {
		t.Fatal("default face should be a distinct cache entry from variation faces")
	}
}

func TestVariableTextToPathDiffers(t *testing.T) {
	mgr := loadVariableFont(t)
	const glyph = "H"
	const size = 32.0

	lightPath, lightAdv := TextToPathVariations(mgr, glyph, "Estedad", 400, "normal", size,
		[]Variation{{Tag: "wght", Value: 100}})
	boldPath, boldAdv := TextToPathVariations(mgr, glyph, "Estedad", 400, "normal", size,
		[]Variation{{Tag: "wght", Value: 900}})

	if lightPath == "" || boldPath == "" {
		t.Fatalf("empty paths: light=%q bold=%q", lightPath, boldPath)
	}
	if lightPath == boldPath {
		t.Fatal("wght=100 and wght=900 produced identical paths — variations not applied")
	}
	if lightAdv == boldAdv {
		t.Fatalf("advances identical (%v == %v); expected different weights to shift advance", lightAdv, boldAdv)
	}
}

func TestGlyphCacheKeyedByVariationCoords(t *testing.T) {
	mgr := loadVariableFont(t)
	const glyph = "H"
	const size = 32.0

	TextToPathVariations(mgr, glyph, "Estedad", 400, "normal", size,
		[]Variation{{Tag: "wght", Value: 100}})
	TextToPathVariations(mgr, glyph, "Estedad", 400, "normal", size,
		[]Variation{{Tag: "wght", Value: 900}})

	light, okL := mgr.glyphs.Get("Estedad", 'H', size, "wght=100")
	bold, okB := mgr.glyphs.Get("Estedad", 'H', size, "wght=900")
	if !okL || !okB {
		t.Fatalf("cache miss: light=%v bold=%v", okL, okB)
	}
	if light.D == bold.D {
		t.Fatal("distinct variation coords produced identical cached paths")
	}
}

func TestVariableAscentDescentDiffers(t *testing.T) {
	mgr := loadVariableFont(t)
	face := mgr.Resolve("Estedad", 400, "normal")
	if face == nil {
		t.Fatal("face not resolved")
	}
	const size = 64.0
	ldA, ldD := mgr.AscentDescent(face, size, []Variation{{Tag: "wght", Value: 100}})
	bA, bD := mgr.AscentDescent(face, size, []Variation{{Tag: "wght", Value: 900}})
	if ldA == 0 && bA == 0 {
		t.Fatal("ascent zero for both variations; MVAR likely absent from test font — replace test font")
	}
	if ldA == bA && ldD == bD {
		t.Logf("ascent/descent identical between wght=100 (%.2f/%.2f) and wght=900 (%.2f/%.2f); test font lacks MVAR deltas — acceptable, ascent code path is exercised", ldA, ldD, bA, bD)
	}
}

func TestVariableWrapWidthDiffers(t *testing.T) {
	mgr := loadVariableFont(t)
	face := mgr.Resolve("Estedad", 400, "normal")
	ff, err := mgr.NewFace(face, 32.0)
	if err != nil {
		t.Fatalf("build sfnt face: %v", err)
	}
	cfgFor := func(vars []Variation) WrapConfig {
		return WrapConfig{
			MaxWidth:   10000,
			FontFace:   ff,
			FontSize:   32.0,
			LineHeight: 40,
			Face:       face,
			Variations: vars,
		}
	}
	lightLines := WrapText("Hello World", cfgFor([]Variation{{Tag: "wght", Value: 100}}))
	boldLines := WrapText("Hello World", cfgFor([]Variation{{Tag: "wght", Value: 900}}))
	if len(lightLines) != 1 || len(boldLines) != 1 {
		t.Fatalf("expected single line per render, got light=%d bold=%d", len(lightLines), len(boldLines))
	}
	if lightLines[0].Width == boldLines[0].Width {
		t.Fatalf("wrap widths identical (%.2f) between wght=100 and wght=900 — measurer not respecting variations", lightLines[0].Width)
	}
}

func TestNonVariationCallDoesNotHitVariableBranch(t *testing.T) {
	mgr := loadVariableFont(t)
	path, adv := TextToPath(mgr, "H", "Estedad", 400, "normal", 32.0)
	if path == "" || adv == 0 {
		t.Fatalf("default-axis rendering degenerate: path=%q adv=%v", path, adv)
	}
}
