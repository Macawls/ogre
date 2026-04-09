package render

import (
	"strings"
	"testing"
)

func TestRenderTransformEmpty(t *testing.T) {
	result := RenderTransform("", "", 0, 0, 100, 50)
	if result != "" {
		t.Fatalf("expected empty string for no transform, got: %s", result)
	}
}

func TestRenderTransformNone(t *testing.T) {
	result := RenderTransform("none", "", 0, 0, 100, 50)
	if result != "" {
		t.Fatalf("expected empty for 'none', got: %s", result)
	}
}

func TestRenderTransformTranslateX(t *testing.T) {
	result := RenderTransform("translateX(10px)", "", 0, 0, 100, 50)
	if !strings.Contains(result, "translate(10, 0)") {
		t.Fatalf("expected translate(10, 0), got: %s", result)
	}
}

func TestRenderTransformTranslateY(t *testing.T) {
	result := RenderTransform("translateY(20px)", "", 0, 0, 100, 50)
	if !strings.Contains(result, "translate(0, 20)") {
		t.Fatalf("expected translate(0, 20), got: %s", result)
	}
}

func TestRenderTransformTranslate(t *testing.T) {
	result := RenderTransform("translate(10px, 20px)", "", 0, 0, 100, 50)
	if !strings.Contains(result, "translate(10, 20)") {
		t.Fatalf("expected translate(10, 20), got: %s", result)
	}
}

func TestRenderTransformScale(t *testing.T) {
	result := RenderTransform("scale(2)", "", 0, 0, 100, 50)
	if !strings.Contains(result, "scale(2)") {
		t.Fatalf("expected scale(2), got: %s", result)
	}
}

func TestRenderTransformScaleXY(t *testing.T) {
	result := RenderTransform("scaleX(2)", "", 0, 0, 100, 50)
	if !strings.Contains(result, "scale(2, 1)") {
		t.Fatalf("expected scale(2, 1), got: %s", result)
	}

	result = RenderTransform("scaleY(0.5)", "", 0, 0, 100, 50)
	if !strings.Contains(result, "scale(1, 0.5)") {
		t.Fatalf("expected scale(1, 0.5), got: %s", result)
	}
}

func TestRenderTransformRotate(t *testing.T) {
	result := RenderTransform("rotate(45deg)", "", 0, 0, 100, 50)
	if !strings.Contains(result, "rotate(45)") {
		t.Fatalf("expected rotate(45), got: %s", result)
	}
}

func TestRenderTransformSkew(t *testing.T) {
	result := RenderTransform("skewX(10deg)", "", 0, 0, 100, 50)
	if !strings.Contains(result, "skewX(10)") {
		t.Fatalf("expected skewX(10), got: %s", result)
	}

	result = RenderTransform("skewY(15deg)", "", 0, 0, 100, 50)
	if !strings.Contains(result, "skewY(15)") {
		t.Fatalf("expected skewY(15), got: %s", result)
	}
}

func TestRenderTransformOriginDefault(t *testing.T) {
	result := RenderTransform("rotate(45deg)", "", 10, 20, 100, 50)
	if !strings.Contains(result, "translate(60, 45)") {
		t.Fatalf("expected default origin at center (60, 45), got: %s", result)
	}
	if !strings.Contains(result, "translate(-60, -45)") {
		t.Fatalf("expected inverse translate, got: %s", result)
	}
}

func TestRenderTransformOriginCustom(t *testing.T) {
	result := RenderTransform("rotate(90deg)", "0% 0%", 10, 20, 100, 50)
	if !strings.Contains(result, "translate(10, 20)") {
		t.Fatalf("expected origin at top-left (10, 20), got: %s", result)
	}
}

func TestRenderTransformMultiple(t *testing.T) {
	result := RenderTransform("translateX(10px) rotate(45deg)", "", 0, 0, 100, 50)
	if !strings.Contains(result, "translate(10, 0)") {
		t.Fatalf("expected translateX, got: %s", result)
	}
	if !strings.Contains(result, "rotate(45)") {
		t.Fatalf("expected rotate, got: %s", result)
	}
}

func TestRenderTransformOriginKeywords(t *testing.T) {
	result := RenderTransform("scale(2)", "left top", 0, 0, 100, 50)
	if !strings.Contains(result, "translate(0, 0)") {
		t.Fatalf("expected origin at (0, 0) for left top, got: %s", result)
	}
}

func TestParseAngleRad(t *testing.T) {
	deg := parseAngleArg("3.14159rad")
	if deg < 179 || deg > 181 {
		t.Fatalf("expected ~180 degrees, got: %.4g", deg)
	}
}

func TestParseAngleTurn(t *testing.T) {
	deg := parseAngleArg("0.5turn")
	if deg != 180 {
		t.Fatalf("expected 180 degrees, got: %.4g", deg)
	}
}
