package ogre

import (
	"strings"
	"testing"
)

func TestRenderBasicDiv(t *testing.T) {
	result, err := Render(`<div style="width:200px;height:100px;background-color:red"></div>`, Options{Width: 400, Height: 300})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("result is nil")
	}
	if result.ContentType != "image/svg+xml" {
		t.Errorf("content type = %q, want %q", result.ContentType, "image/svg+xml")
	}
	svg := string(result.Data)
	if !strings.Contains(svg, "<svg") {
		t.Error("output does not contain <svg")
	}
	if !strings.Contains(svg, "<rect") {
		t.Error("output does not contain <rect")
	}
	if !strings.Contains(svg, "#ff0000") && !strings.Contains(svg, "red") {
		t.Error("output does not contain red fill")
	}
	if result.Width != 400 || result.Height != 300 {
		t.Errorf("dimensions = %dx%d, want 400x300", result.Width, result.Height)
	}
}

func TestRenderTextContent(t *testing.T) {
	result, err := Render(`<div style="font-size:24px;color:blue">Hello World</div>`, Options{Width: 400, Height: 300})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("result is nil")
	}
	svg := string(result.Data)
	if !strings.Contains(svg, "<path") {
		t.Error("output does not contain <path element for embedded text")
	}
	if !strings.Contains(svg, `fill="`) {
		t.Error("output does not contain fill attribute")
	}
}

func TestRenderDefaults(t *testing.T) {
	result, err := Render(`<div style="width:100px;height:50px;background-color:green"></div>`, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Width != 1200 {
		t.Errorf("default width = %d, want 1200", result.Width)
	}
	if result.Height != 630 {
		t.Errorf("default height = %d, want 630", result.Height)
	}
}

func TestRenderPNG(t *testing.T) {
	result, err := Render(`<div style="width:100px;height:100px;background-color:red"></div>`, Options{Format: FormatPNG})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ContentType != "image/png" {
		t.Errorf("content type = %q, want image/png", result.ContentType)
	}
	if len(result.Data) == 0 {
		t.Fatal("expected non-empty PNG data")
	}
}
