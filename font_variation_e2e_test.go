package ogre

import (
	"os"
	"strings"
	"testing"
)

func TestFontVariationSettingsEndToEnd(t *testing.T) {
	data, err := os.ReadFile("font/testdata/Estedad-VF.ttf")
	if err != nil {
		t.Fatalf("load bundled variable font: %v", err)
	}

	renderOne := func(varSetting string) string {
		r := NewRenderer()
		if err := r.LoadFont(FontSource{Name: "VarTest", Weight: 400, Style: "normal", Data: data}); err != nil {
			t.Fatalf("load font: %v", err)
		}
		html := `<div style="width:400px;height:120px;background:white;display:flex;align-items:center;justify-content:center">` +
			`<span style="font-family:VarTest;font-size:64px;font-variation-settings:` + varSetting + `">Ogre</span>` +
			`</div>`
		result, err := r.Render(html, Options{Width: 400, Height: 120, Format: FormatSVG})
		if err != nil {
			t.Fatalf("render: %v", err)
		}
		return string(result.Data)
	}

	light := renderOne(`'wght' 100`)
	bold := renderOne(`'wght' 900`)

	if light == bold {
		t.Fatal("wght=100 and wght=900 produced identical SVG — variations not reaching the render pipeline")
	}

	countPaths := func(svg string) int { return strings.Count(svg, `<path`) }
	if countPaths(light) == 0 || countPaths(bold) == 0 {
		t.Fatalf("no <path> elements in output: light=%d bold=%d", countPaths(light), countPaths(bold))
	}
}
