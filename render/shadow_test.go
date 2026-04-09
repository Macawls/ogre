package render

import (
	"fmt"
	"strings"
	"testing"

	"github.com/macawls/ogre/style"
)

func shadowIDGen() func(string) string {
	n := 0
	return func(prefix string) string {
		n++
		return fmt.Sprintf("%s%d", prefix, n)
	}
}

func TestRenderBoxShadowEmpty(t *testing.T) {
	result := RenderBoxShadow(nil, 0, 0, 100, 50, 0, shadowIDGen())
	if result != "" {
		t.Fatalf("expected empty string for nil shadows, got: %s", result)
	}
}

func TestRenderBoxShadowSimple(t *testing.T) {
	shadows := []style.Shadow{
		{
			OffsetX: 5,
			OffsetY: 5,
			Blur:    10,
			Spread:  0,
			Color:   style.Color{R: 0, G: 0, B: 0, A: 0.5},
		},
	}

	result := RenderBoxShadow(shadows, 10, 20, 100, 50, 0, shadowIDGen())

	if !strings.Contains(result, "<defs>") {
		t.Fatal("expected <defs> element")
	}
	if !strings.Contains(result, "<filter") {
		t.Fatal("expected <filter> element")
	}
	if !strings.Contains(result, "feGaussianBlur") {
		t.Fatal("expected feGaussianBlur")
	}
	if !strings.Contains(result, `stdDeviation="5"`) {
		t.Fatalf("expected stdDeviation=5, got: %s", result)
	}
	if !strings.Contains(result, `dx="5"`) {
		t.Fatalf("expected dx=5, got: %s", result)
	}
	if !strings.Contains(result, `dy="5"`) {
		t.Fatalf("expected dy=5, got: %s", result)
	}
	if !strings.Contains(result, `flood-color="#000000"`) {
		t.Fatalf("expected flood-color, got: %s", result)
	}
	if !strings.Contains(result, `flood-opacity="0.5"`) {
		t.Fatalf("expected flood-opacity, got: %s", result)
	}
	if !strings.Contains(result, "<rect") {
		t.Fatal("expected shadow rect")
	}
}

func TestRenderBoxShadowWithSpread(t *testing.T) {
	shadows := []style.Shadow{
		{
			OffsetX: 0,
			OffsetY: 2,
			Blur:    4,
			Spread:  3,
			Color:   style.Color{R: 255, G: 0, B: 0, A: 1},
		},
	}

	result := RenderBoxShadow(shadows, 10, 10, 100, 50, 0, shadowIDGen())

	if !strings.Contains(result, `x="7"`) {
		t.Fatalf("expected spread-adjusted x=7, got: %s", result)
	}
	if !strings.Contains(result, `width="106"`) {
		t.Fatalf("expected spread-adjusted width=106, got: %s", result)
	}
}

func TestRenderBoxShadowWithRadius(t *testing.T) {
	shadows := []style.Shadow{
		{
			OffsetX: 2,
			OffsetY: 2,
			Blur:    6,
			Color:   style.Color{R: 0, G: 0, B: 0, A: 1},
		},
	}

	result := RenderBoxShadow(shadows, 0, 0, 100, 100, 8, shadowIDGen())

	if !strings.Contains(result, `rx="8"`) {
		t.Fatalf("expected border radius on shadow rect, got: %s", result)
	}
}

func TestRenderBoxShadowInsetSkipped(t *testing.T) {
	shadows := []style.Shadow{
		{
			OffsetX: 5,
			OffsetY: 5,
			Blur:    10,
			Color:   style.Color{R: 0, G: 0, B: 0, A: 1},
			Inset:   true,
		},
	}

	result := RenderBoxShadow(shadows, 0, 0, 100, 50, 0, shadowIDGen())

	if strings.Contains(result, "<filter") {
		t.Fatal("inset shadows should be skipped by RenderBoxShadow")
	}
}

func TestRenderInsetBoxShadowBasic(t *testing.T) {
	shadows := []style.Shadow{
		{
			OffsetX: 0,
			OffsetY: 0,
			Blur:    10,
			Color:   style.Color{R: 0, G: 0, B: 0, A: 0.5},
			Inset:   true,
		},
	}

	result := RenderInsetBoxShadow(shadows, 10, 20, 100, 50, 0, shadowIDGen())

	if !strings.Contains(result, "<clipPath") {
		t.Fatal("expected <clipPath> for inset shadow")
	}
	if !strings.Contains(result, "<filter") {
		t.Fatal("expected <filter> for inset shadow")
	}
	if !strings.Contains(result, "feGaussianBlur") {
		t.Fatal("expected feGaussianBlur")
	}
	rectCount := strings.Count(result, "<rect")
	if rectCount < 5 {
		t.Fatalf("expected at least 5 rects (1 clip + 4 sides), got %d", rectCount)
	}
}

func TestRenderInsetBoxShadowSkipsOutset(t *testing.T) {
	shadows := []style.Shadow{
		{
			OffsetX: 5,
			OffsetY: 5,
			Blur:    10,
			Color:   style.Color{R: 0, G: 0, B: 0, A: 1},
			Inset:   false,
		},
	}

	result := RenderInsetBoxShadow(shadows, 0, 0, 100, 50, 0, shadowIDGen())

	if strings.Contains(result, "<filter") {
		t.Fatal("outset shadows should be skipped by RenderInsetBoxShadow")
	}
}

func TestRenderInsetBoxShadowWithRadius(t *testing.T) {
	shadows := []style.Shadow{
		{
			OffsetX: 0,
			OffsetY: 0,
			Blur:    8,
			Color:   style.Color{R: 0, G: 0, B: 0, A: 1},
			Inset:   true,
		},
	}

	result := RenderInsetBoxShadow(shadows, 0, 0, 100, 100, 10, shadowIDGen())

	if !strings.Contains(result, `rx="10"`) {
		t.Fatalf("expected border radius on clip rect, got: %s", result)
	}
}

func TestRenderInsetBoxShadowEmpty(t *testing.T) {
	result := RenderInsetBoxShadow(nil, 0, 0, 100, 50, 0, shadowIDGen())
	if result != "" {
		t.Fatalf("expected empty string for nil shadows, got: %s", result)
	}
}

func TestRenderBoxShadowMultiple(t *testing.T) {
	shadows := []style.Shadow{
		{OffsetX: 2, OffsetY: 2, Blur: 4, Color: style.Color{R: 0, G: 0, B: 0, A: 1}},
		{OffsetX: 4, OffsetY: 4, Blur: 8, Color: style.Color{R: 255, G: 0, B: 0, A: 1}},
	}

	result := RenderBoxShadow(shadows, 0, 0, 100, 50, 0, shadowIDGen())

	filterCount := strings.Count(result, "<filter")
	if filterCount != 2 {
		t.Fatalf("expected 2 filters, got %d: %s", filterCount, result)
	}
	rectCount := strings.Count(result, "<rect")
	if rectCount != 2 {
		t.Fatalf("expected 2 shadow rects, got %d", rectCount)
	}
}
