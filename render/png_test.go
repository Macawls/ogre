package render

import (
	"bytes"
	"image/color"
	"image/png"
	"testing"

	"github.com/macawls/ogre/layout"
	"github.com/macawls/ogre/parse"
	"github.com/macawls/ogre/style"
)

func buildTestTree(cs *style.ComputedStyle, w, h float64) (*layout.LayoutTree, map[*parse.Node]*style.ComputedStyle) {
	pn := &parse.Node{Type: parse.ElementNode, Tag: "div"}
	styles := map[*parse.Node]*style.ComputedStyle{pn: cs}

	ln := layout.NewNode(layout.Style{
		Width:  layout.Pt(w),
		Height: layout.Pt(h),
	})
	ln.Layout = layout.Layout{X: 0, Y: 0, Width: w, Height: h}

	tree := &layout.LayoutTree{
		Root:    ln,
		NodeMap: map[*parse.Node]*layout.Node{pn: ln},
	}
	return tree, styles
}

func TestRenderPNG_ValidBytes(t *testing.T) {
	tree, styles := buildTestTree(&style.ComputedStyle{
		BackgroundColor: style.Color{R: 100, G: 150, B: 200, A: 1},
	}, 200, 100)

	data, err := RenderPNG(tree, styles, nil, 200, 100)
	if err != nil {
		t.Fatalf("RenderPNG returned error: %v", err)
	}

	if len(data) < 8 {
		t.Fatal("output too short to contain PNG header")
	}

	magic := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	if !bytes.HasPrefix(data, magic) {
		t.Errorf("output does not start with PNG magic bytes, got %x", data[:8])
	}
}

func TestRenderPNG_Dimensions(t *testing.T) {
	tree, styles := buildTestTree(&style.ComputedStyle{}, 320, 240)

	data, err := RenderPNG(tree, styles, nil, 320, 240)
	if err != nil {
		t.Fatalf("RenderPNG returned error: %v", err)
	}

	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("failed to decode PNG: %v", err)
	}

	bounds := img.Bounds()
	if bounds.Dx() != 320 || bounds.Dy() != 240 {
		t.Errorf("expected 320x240, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestRenderPNG_BackgroundColor(t *testing.T) {
	bg := style.Color{R: 255, G: 0, B: 0, A: 1}
	tree, styles := buildTestTree(&style.ComputedStyle{
		BackgroundColor: bg,
	}, 100, 100)

	data, err := RenderPNG(tree, styles, nil, 100, 100)
	if err != nil {
		t.Fatalf("RenderPNG returned error: %v", err)
	}

	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("failed to decode PNG: %v", err)
	}

	cx, cy := 50, 50
	r, g, b, a := img.At(cx, cy).RGBA()
	wantR, wantG, wantB, wantA := color.RGBA{255, 0, 0, 255}.RGBA()
	if r != wantR || g != wantG || b != wantB || a != wantA {
		t.Errorf("pixel at (%d,%d) = (%d,%d,%d,%d), want (%d,%d,%d,%d)",
			cx, cy, r>>8, g>>8, b>>8, a>>8, wantR>>8, wantG>>8, wantB>>8, wantA>>8)
	}
}

func TestRenderPNG_LinearGradient(t *testing.T) {
	cs := &style.ComputedStyle{
		BackgroundImage: "linear-gradient(90deg, red, blue)",
	}
	tree, styles := buildTestTree(cs, 100, 100)

	data, err := RenderPNG(tree, styles, nil, 100, 100)
	if err != nil {
		t.Fatalf("RenderPNG returned error: %v", err)
	}

	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("failed to decode PNG: %v", err)
	}

	hasNonBlack := false
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			if r > 0 || g > 0 || b > 0 {
				hasNonBlack = true
				break
			}
		}
		if hasNonBlack {
			break
		}
	}
	if !hasNonBlack {
		t.Error("linear gradient output is solid black")
	}

	lr, _, lb, _ := img.At(5, 50).RGBA()
	rr, _, rb, _ := img.At(95, 50).RGBA()
	if lr>>8 < 200 {
		t.Errorf("left edge should be red-ish, got R=%d", lr>>8)
	}
	if rb>>8 < 200 {
		t.Errorf("right edge should be blue-ish, got B=%d", rb>>8)
	}
	if lb>>8 > 80 {
		t.Errorf("left edge should have little blue, got B=%d", lb>>8)
	}
	if rr>>8 > 80 {
		t.Errorf("right edge should have little red, got R=%d", rr>>8)
	}
}

func buildTestTreeWithChild(parentCS, childCS *style.ComputedStyle, pw, ph, cw, ch, cx, cy float64) (*layout.LayoutTree, map[*parse.Node]*style.ComputedStyle) {
	parentPN := &parse.Node{Type: parse.ElementNode, Tag: "div"}
	childPN := &parse.Node{Type: parse.ElementNode, Tag: "div"}
	parentPN.Children = []*parse.Node{childPN}

	styles := map[*parse.Node]*style.ComputedStyle{
		parentPN: parentCS,
		childPN:  childCS,
	}

	parentLN := layout.NewNode(layout.Style{
		Width:  layout.Pt(pw),
		Height: layout.Pt(ph),
	})
	parentLN.Layout = layout.Layout{X: 0, Y: 0, Width: pw, Height: ph}

	childLN := layout.NewNode(layout.Style{
		Width:  layout.Pt(cw),
		Height: layout.Pt(ch),
	})
	childLN.Layout = layout.Layout{X: cx, Y: cy, Width: cw, Height: ch}
	parentLN.Children = []*layout.Node{childLN}

	tree := &layout.LayoutTree{
		Root: parentLN,
		NodeMap: map[*parse.Node]*layout.Node{
			parentPN: parentLN,
			childPN:  childLN,
		},
	}
	return tree, styles
}

func TestRenderPNG_OverflowHiddenClipsChildren(t *testing.T) {
	parentCS := &style.ComputedStyle{
		Overflow: style.OverflowHidden,
	}
	childCS := &style.ComputedStyle{
		BackgroundColor: style.Color{R: 255, G: 0, B: 0, A: 1},
	}

	tree, styles := buildTestTreeWithChild(parentCS, childCS, 100, 100, 200, 200, 0, 0)

	data, err := RenderPNG(tree, styles, nil, 200, 200)
	if err != nil {
		t.Fatalf("RenderPNG returned error: %v", err)
	}

	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("failed to decode PNG: %v", err)
	}

	r50, _, _, _ := img.At(50, 50).RGBA()
	if r50>>8 < 200 {
		t.Errorf("pixel inside parent (50,50) should be red, got R=%d", r50>>8)
	}

	r150, g150, b150, _ := img.At(150, 150).RGBA()
	if r150>>8 > 10 || g150>>8 < 200 || b150>>8 < 200 {
		isWhite := r150>>8 > 240 && g150>>8 > 240 && b150>>8 > 240
		if !isWhite {
			t.Errorf("pixel outside parent (150,150) should be white (clipped), got R=%d G=%d B=%d", r150>>8, g150>>8, b150>>8)
		}
	}
}

func TestRenderPNG_RoundedCorners(t *testing.T) {
	cs := &style.ComputedStyle{
		BackgroundColor:         style.Color{R: 255, G: 0, B: 0, A: 1},
		BorderTopLeftRadius:     20,
		BorderTopRightRadius:    20,
		BorderBottomLeftRadius:  20,
		BorderBottomRightRadius: 20,
	}
	tree, styles := buildTestTree(cs, 100, 100)

	data, err := RenderPNG(tree, styles, nil, 100, 100)
	if err != nil {
		t.Fatalf("RenderPNG returned error: %v", err)
	}

	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("failed to decode PNG: %v", err)
	}

	cr, _, _, _ := img.At(50, 50).RGBA()
	if cr>>8 < 200 {
		t.Errorf("center pixel should be red, got R=%d", cr>>8)
	}

	cornR, cornG, cornB, _ := img.At(0, 0).RGBA()
	isWhite := cornR>>8 > 240 && cornG>>8 > 240 && cornB>>8 > 240
	if !isWhite {
		t.Errorf("corner pixel (0,0) should be white (rounded off), got R=%d G=%d B=%d", cornR>>8, cornG>>8, cornB>>8)
	}
}

func TestRenderPNG_BoxShadow(t *testing.T) {
	cs := &style.ComputedStyle{
		BackgroundColor: style.Color{R: 255, G: 255, B: 255, A: 1},
		BoxShadow:       "5px 5px 10px rgba(0,0,0,0.5)",
	}
	tree, styles := buildTestTree(cs, 60, 60)

	data, err := RenderPNG(tree, styles, nil, 100, 100)
	if err != nil {
		t.Fatalf("RenderPNG returned error: %v", err)
	}

	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("failed to decode PNG: %v", err)
	}

	r, g, b, _ := img.At(70, 70).RGBA()
	if r>>8 > 250 && g>>8 > 250 && b>>8 > 250 {
		t.Error("shadow region (70,70) should not be pure white")
	}

	cr, cg, cb, _ := img.At(30, 30).RGBA()
	if cr>>8 < 240 || cg>>8 < 240 || cb>>8 < 240 {
		t.Errorf("center (30,30) should be white background, got R=%d G=%d B=%d", cr>>8, cg>>8, cb>>8)
	}
}

func TestRenderPNG_BoxShadowNoBlur(t *testing.T) {
	cs := &style.ComputedStyle{
		BackgroundColor: style.Color{R: 0, G: 0, B: 255, A: 1},
		BoxShadow:       "10px 10px 0 rgba(255,0,0,1)",
	}
	tree, styles := buildTestTree(cs, 50, 50)

	data, err := RenderPNG(tree, styles, nil, 100, 100)
	if err != nil {
		t.Fatalf("RenderPNG returned error: %v", err)
	}

	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("failed to decode PNG: %v", err)
	}

	r, _, _, _ := img.At(55, 55).RGBA()
	if r>>8 < 200 {
		t.Errorf("shadow at (55,55) should be red-ish, got R=%d", r>>8)
	}
}

func TestRenderPNG_BoxShadowSpread(t *testing.T) {
	cs := &style.ComputedStyle{
		BackgroundColor: style.Color{R: 255, G: 255, B: 255, A: 1},
		BoxShadow:       "0 0 0 5px rgba(255,0,0,1)",
	}
	tree, styles := buildTestTree(cs, 40, 40)

	data, err := RenderPNG(tree, styles, nil, 100, 100)
	if err != nil {
		t.Fatalf("RenderPNG returned error: %v", err)
	}

	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("failed to decode PNG: %v", err)
	}

	rr, _, _, _ := img.At(20, 42).RGBA()
	if rr>>8 < 200 {
		t.Errorf("spread region should be red, got R=%d at (20,42)", rr>>8)
	}
}

func TestRenderPNG_InsetBoxShadow(t *testing.T) {
	cs := &style.ComputedStyle{
		BackgroundColor: style.Color{R: 255, G: 255, B: 255, A: 1},
		BoxShadow:       "inset 0 0 10px rgba(0,0,0,0.8)",
	}
	tree, styles := buildTestTree(cs, 80, 80)

	data, err := RenderPNG(tree, styles, nil, 100, 100)
	if err != nil {
		t.Fatalf("RenderPNG returned error: %v", err)
	}

	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("failed to decode PNG: %v", err)
	}

	er, eg, eb, _ := img.At(1, 40).RGBA()
	edgeBrightness := (er>>8 + eg>>8 + eb>>8) / 3

	cr, cg, cb, _ := img.At(40, 40).RGBA()
	centerBrightness := (cr>>8 + cg>>8 + cb>>8) / 3

	if edgeBrightness >= centerBrightness {
		t.Errorf("edge should be darker than center for inset shadow: edge=%d center=%d", edgeBrightness, centerBrightness)
	}
}

func TestRenderPNG_MultipleShadows(t *testing.T) {
	cs := &style.ComputedStyle{
		BackgroundColor: style.Color{R: 255, G: 255, B: 255, A: 1},
		BoxShadow:       "5px 5px 5px rgba(0,0,0,0.5), inset 0 0 5px rgba(0,0,0,0.5)",
	}
	tree, styles := buildTestTree(cs, 60, 60)

	data, err := RenderPNG(tree, styles, nil, 100, 100)
	if err != nil {
		t.Fatalf("RenderPNG returned error: %v", err)
	}

	if len(data) < 8 {
		t.Fatal("output too short")
	}
}

func TestRenderPNG_RadialGradient(t *testing.T) {
	cs := &style.ComputedStyle{
		BackgroundImage: "radial-gradient(circle, red, blue)",
	}
	tree, styles := buildTestTree(cs, 100, 100)

	data, err := RenderPNG(tree, styles, nil, 100, 100)
	if err != nil {
		t.Fatalf("RenderPNG returned error: %v", err)
	}

	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("failed to decode PNG: %v", err)
	}

	hasNonBlack := false
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			if r > 0 || g > 0 || b > 0 {
				hasNonBlack = true
				break
			}
		}
		if hasNonBlack {
			break
		}
	}
	if !hasNonBlack {
		t.Error("radial gradient output is solid black")
	}

	cr, _, cb, _ := img.At(50, 50).RGBA()
	er, _, eb, _ := img.At(0, 0).RGBA()
	if cr>>8 < 200 {
		t.Errorf("center should be red-ish, got R=%d", cr>>8)
	}
	if cb>>8 > 50 {
		t.Errorf("center should have little blue, got B=%d", cb>>8)
	}
	if eb>>8 < 150 {
		t.Errorf("corner should be blue-ish, got B=%d", eb>>8)
	}
	if er>>8 > 100 {
		t.Errorf("corner should have less red, got R=%d", er>>8)
	}
}
