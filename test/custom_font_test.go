package test

import (
	"strings"
	"testing"

	"golang.org/x/image/font/gofont/goitalic"

	"github.com/macawls/ogre"
)

func TestCustomFont(t *testing.T) {
	html := `<div style="display:flex;width:100%;height:100%;background-color:white;align-items:center;justify-content:center">
		<div style="display:flex;font-family:custom;font-size:48px;color:#333">Custom Font Text</div>
	</div>`

	result, err := ogre.Render(html, ogre.Options{
		Width:  800,
		Height: 400,
		Fonts: []ogre.FontSource{{
			Name:   "custom",
			Weight: 400,
			Style:  "italic",
			Data:   goitalic.TTF,
		}},
	})
	if err != nil {
		t.Fatal(err)
	}

	svg := string(result.Data)
	if len(svg) == 0 {
		t.Fatal("empty result")
	}

	t.Log(svg[:min(500, len(svg))])

	if !strings.Contains(svg, "<path") {
		t.Fatal("expected <path> elements from glyph rendering, got none — custom font glyphs not embedded")
	}

	if strings.Contains(svg, `font-family="custom"`) {
		t.Fatal("found fallback <text> element with font-family=\"custom\" — font was not resolved to glyph paths")
	}
}

func TestCustomFontFallsBackWhenNameMismatch(t *testing.T) {
	html := `<div style="display:flex;width:100%;height:100%;background-color:white">
		<div style="display:flex;font-family:unknown;font-size:32px;color:black">Fallback Test</div>
	</div>`

	result, err := ogre.Render(html, ogre.Options{
		Width:  400,
		Height: 200,
		Fonts: []ogre.FontSource{{
			Name:   "custom",
			Weight: 400,
			Style:  "normal",
			Data:   goitalic.TTF,
		}},
	})
	if err != nil {
		t.Fatal(err)
	}

	svg := string(result.Data)
	if len(svg) == 0 {
		t.Fatal("empty result")
	}

	if !strings.Contains(svg, "<path") {
		t.Fatal("expected fallback to default font with glyph paths")
	}
}

func TestCustomFontWeightResolution(t *testing.T) {
	html := `<div style="display:flex;width:100%;height:100%;background-color:white">
		<div style="display:flex;font-family:custom;font-weight:700;font-size:32px;color:black">Bold Request</div>
	</div>`

	result, err := ogre.Render(html, ogre.Options{
		Width:  400,
		Height: 200,
		Fonts: []ogre.FontSource{{
			Name:   "custom",
			Weight: 400,
			Style:  "normal",
			Data:   goitalic.TTF,
		}},
	})
	if err != nil {
		t.Fatal(err)
	}

	svg := string(result.Data)
	if len(svg) == 0 {
		t.Fatal("empty result")
	}

	if !strings.Contains(svg, "<path") {
		t.Fatal("expected glyph paths — weight fallback within custom font family should still produce paths")
	}
}
