package render

import (
	"strings"
	"testing"

	"github.com/macawls/ogre/font"
	"github.com/macawls/ogre/style"
)

func TestRenderTextBasic(t *testing.T) {
	lines := []font.TextLine{
		{Text: "Hello World", Width: 100, Y: 0},
	}
	cs := &style.ComputedStyle{
		FontFamily: "Arial",
		FontSize:   16,
		FontWeight: 400,
		Color:      style.Color{R: 0, G: 0, B: 0, A: 1},
	}

	result := RenderText(lines, cs, 0, 0, 200, 50)

	if !strings.Contains(result.Content, "<text") {
		t.Fatal("expected <text> element")
	}
	if !strings.Contains(result.Content, `font-family="Arial"`) {
		t.Fatal("expected font-family Arial")
	}
	if !strings.Contains(result.Content, `font-size="16"`) {
		t.Fatal("expected font-size 16")
	}
	if !strings.Contains(result.Content, "Hello World") {
		t.Fatal("expected text content")
	}
	if result.Decorations != "" {
		t.Fatal("expected no decorations")
	}
}

func TestRenderTextEmpty(t *testing.T) {
	result := RenderText(nil, &style.ComputedStyle{}, 0, 0, 200, 50)
	if result.Content != "" || result.Decorations != "" {
		t.Fatal("expected empty result for nil lines")
	}
}

func TestRenderTextAlignLeft(t *testing.T) {
	lines := []font.TextLine{
		{Text: "Left", Width: 50, Y: 0},
	}
	cs := &style.ComputedStyle{
		FontSize:  16,
		TextAlign: style.TextAlignLeft,
		Color:     style.Color{R: 0, G: 0, B: 0, A: 1},
	}

	result := RenderText(lines, cs, 10, 0, 200, 50)

	if !strings.Contains(result.Content, `x="10"`) {
		t.Fatalf("expected x=10 for left align, got: %s", result.Content)
	}
}

func TestRenderTextAlignCenter(t *testing.T) {
	lines := []font.TextLine{
		{Text: "Center", Width: 60, Y: 0},
	}
	cs := &style.ComputedStyle{
		FontSize:  16,
		TextAlign: style.TextAlignCenter,
		Color:     style.Color{R: 0, G: 0, B: 0, A: 1},
	}

	result := RenderText(lines, cs, 0, 0, 200, 50)

	if !strings.Contains(result.Content, `x="70"`) {
		t.Fatalf("expected x=70 for center align (200-60)/2=70, got: %s", result.Content)
	}
}

func TestRenderTextAlignRight(t *testing.T) {
	lines := []font.TextLine{
		{Text: "Right", Width: 50, Y: 0},
	}
	cs := &style.ComputedStyle{
		FontSize:  16,
		TextAlign: style.TextAlignRight,
		Color:     style.Color{R: 0, G: 0, B: 0, A: 1},
	}

	result := RenderText(lines, cs, 0, 0, 200, 50)

	if !strings.Contains(result.Content, `x="150"`) {
		t.Fatalf("expected x=150 for right align (200-50=150), got: %s", result.Content)
	}
}

func TestRenderTextTransformUppercase(t *testing.T) {
	lines := []font.TextLine{
		{Text: "hello", Width: 50, Y: 0},
	}
	cs := &style.ComputedStyle{
		FontSize:      16,
		TextTransform: style.TextTransformUppercase,
		Color:         style.Color{R: 0, G: 0, B: 0, A: 1},
	}

	result := RenderText(lines, cs, 0, 0, 200, 50)

	if !strings.Contains(result.Content, "HELLO") {
		t.Fatalf("expected uppercase text, got: %s", result.Content)
	}
}

func TestRenderTextTransformLowercase(t *testing.T) {
	lines := []font.TextLine{
		{Text: "HELLO", Width: 50, Y: 0},
	}
	cs := &style.ComputedStyle{
		FontSize:      16,
		TextTransform: style.TextTransformLowercase,
		Color:         style.Color{R: 0, G: 0, B: 0, A: 1},
	}

	result := RenderText(lines, cs, 0, 0, 200, 50)

	if !strings.Contains(result.Content, "hello") {
		t.Fatalf("expected lowercase text, got: %s", result.Content)
	}
}

func TestRenderTextTransformCapitalize(t *testing.T) {
	lines := []font.TextLine{
		{Text: "hello world", Width: 100, Y: 0},
	}
	cs := &style.ComputedStyle{
		FontSize:      16,
		TextTransform: style.TextTransformCapitalize,
		Color:         style.Color{R: 0, G: 0, B: 0, A: 1},
	}

	result := RenderText(lines, cs, 0, 0, 200, 50)

	if !strings.Contains(result.Content, "Hello World") {
		t.Fatalf("expected capitalized text, got: %s", result.Content)
	}
}

func TestRenderTextDecorationUnderline(t *testing.T) {
	lines := []font.TextLine{
		{Text: "Underlined", Width: 80, Y: 0},
	}
	cs := &style.ComputedStyle{
		FontSize:           16,
		TextDecorationLine: style.TextDecorationUnderline,
		Color:              style.Color{R: 255, G: 0, B: 0, A: 1},
	}

	result := RenderText(lines, cs, 0, 0, 200, 50)

	if result.Decorations == "" {
		t.Fatal("expected decoration for underline")
	}
	if !strings.Contains(result.Decorations, "<line") {
		t.Fatalf("expected <line> element in decorations, got: %s", result.Decorations)
	}
	if !strings.Contains(result.Decorations, `stroke="#ff0000"`) {
		t.Fatalf("expected stroke color from text color, got: %s", result.Decorations)
	}
}

func TestRenderTextDecorationColor(t *testing.T) {
	lines := []font.TextLine{
		{Text: "Custom", Width: 50, Y: 0},
	}
	cs := &style.ComputedStyle{
		FontSize:            16,
		TextDecorationLine:  style.TextDecorationUnderline,
		TextDecorationColor: style.Color{R: 0, G: 0, B: 255, A: 1},
		Color:               style.Color{R: 255, G: 0, B: 0, A: 1},
	}

	result := RenderText(lines, cs, 0, 0, 200, 50)

	if !strings.Contains(result.Decorations, `stroke="#0000ff"`) {
		t.Fatalf("expected custom decoration color, got: %s", result.Decorations)
	}
}

func TestRenderTextMultipleLines(t *testing.T) {
	lines := []font.TextLine{
		{Text: "Line 1", Width: 50, Y: 0},
		{Text: "Line 2", Width: 60, Y: 20},
	}
	cs := &style.ComputedStyle{
		FontSize:   16,
		LineHeight: 20,
		Color:      style.Color{R: 0, G: 0, B: 0, A: 1},
	}

	result := RenderText(lines, cs, 0, 0, 200, 50)

	count := strings.Count(result.Content, "<text")
	if count != 2 {
		t.Fatalf("expected 2 <text> elements, got %d", count)
	}
}
