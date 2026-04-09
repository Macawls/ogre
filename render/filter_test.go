package render

import (
	"strings"
	"testing"
)

func filterTestIDGen(prefix string) string {
	return prefix + "-1"
}

func TestFilterBlur(t *testing.T) {
	defs, attr := RenderCSSFilter("blur(5px)", filterTestIDGen)
	if !strings.Contains(defs, `stdDeviation="5"`) {
		t.Errorf("expected stdDeviation=5, got %s", defs)
	}
	if attr == "" {
		t.Error("expected filter attribute")
	}
}

func TestFilterContrast(t *testing.T) {
	defs, attr := RenderCSSFilter("contrast(1.5)", filterTestIDGen)
	if !strings.Contains(defs, `<feComponentTransfer>`) {
		t.Errorf("expected feComponentTransfer, got %s", defs)
	}
	if !strings.Contains(defs, `slope="1.5"`) {
		t.Errorf("expected slope=1.5, got %s", defs)
	}
	if !strings.Contains(defs, `intercept="-0.25"`) {
		t.Errorf("expected intercept=-0.25, got %s", defs)
	}
	if attr == "" {
		t.Error("expected filter attribute")
	}
}

func TestFilterSaturate(t *testing.T) {
	defs, _ := RenderCSSFilter("saturate(2)", filterTestIDGen)
	if !strings.Contains(defs, `type="saturate"`) {
		t.Errorf("expected type=saturate, got %s", defs)
	}
	if !strings.Contains(defs, `values="2"`) {
		t.Errorf("expected values=2, got %s", defs)
	}
}

func TestFilterSepia(t *testing.T) {
	defs, _ := RenderCSSFilter("sepia(100%)", filterTestIDGen)
	if !strings.Contains(defs, `type="matrix"`) {
		t.Errorf("expected type=matrix, got %s", defs)
	}
	if !strings.Contains(defs, `0.393`) {
		t.Errorf("expected sepia matrix value 0.393, got %s", defs)
	}
}

func TestFilterHueRotate(t *testing.T) {
	defs, _ := RenderCSSFilter("hue-rotate(90deg)", filterTestIDGen)
	if !strings.Contains(defs, `type="hueRotate"`) {
		t.Errorf("expected type=hueRotate, got %s", defs)
	}
	if !strings.Contains(defs, `values="90"`) {
		t.Errorf("expected values=90, got %s", defs)
	}
}

func TestFilterInvert(t *testing.T) {
	defs, _ := RenderCSSFilter("invert(100%)", filterTestIDGen)
	if !strings.Contains(defs, `<feComponentTransfer>`) {
		t.Errorf("expected feComponentTransfer, got %s", defs)
	}
	if !strings.Contains(defs, `tableValues="1 0"`) {
		t.Errorf("expected tableValues='1 0', got %s", defs)
	}
}

func TestFilterDropShadow(t *testing.T) {
	defs, _ := RenderCSSFilter("drop-shadow(4px 4px 10px red)", filterTestIDGen)
	if !strings.Contains(defs, `<feDropShadow`) {
		t.Errorf("expected feDropShadow, got %s", defs)
	}
	if !strings.Contains(defs, `dx="4"`) {
		t.Errorf("expected dx=4, got %s", defs)
	}
	if !strings.Contains(defs, `dy="4"`) {
		t.Errorf("expected dy=4, got %s", defs)
	}
	if !strings.Contains(defs, `stdDeviation="5"`) {
		t.Errorf("expected stdDeviation=5, got %s", defs)
	}
	if !strings.Contains(defs, `flood-color="red"`) {
		t.Errorf("expected flood-color=red, got %s", defs)
	}
}

func TestMultipleFilters(t *testing.T) {
	defs, attr := RenderCSSFilter("blur(5px) grayscale(50%)", filterTestIDGen)
	if !strings.Contains(defs, `feGaussianBlur`) {
		t.Errorf("expected feGaussianBlur, got %s", defs)
	}
	if !strings.Contains(defs, `feColorMatrix`) {
		t.Errorf("expected feColorMatrix, got %s", defs)
	}
	if attr == "" {
		t.Error("expected filter attribute")
	}
}

func TestFilterNone(t *testing.T) {
	defs, attr := RenderCSSFilter("none", filterTestIDGen)
	if defs != "" || attr != "" {
		t.Errorf("expected empty for none, got defs=%q attr=%q", defs, attr)
	}
}

func TestFilterEmpty(t *testing.T) {
	defs, attr := RenderCSSFilter("", filterTestIDGen)
	if defs != "" || attr != "" {
		t.Errorf("expected empty for empty string, got defs=%q attr=%q", defs, attr)
	}
}

func TestFilterInvertPartial(t *testing.T) {
	defs, _ := RenderCSSFilter("invert(50%)", filterTestIDGen)
	if !strings.Contains(defs, `tableValues="0.5 0.5"`) {
		t.Errorf("expected tableValues='0.5 0.5', got %s", defs)
	}
}

func TestFilterDropShadowNoColor(t *testing.T) {
	defs, _ := RenderCSSFilter("drop-shadow(2px 3px 4px)", filterTestIDGen)
	if !strings.Contains(defs, `flood-color="black"`) {
		t.Errorf("expected default flood-color=black, got %s", defs)
	}
}

func TestFilterContrastPercent(t *testing.T) {
	defs, _ := RenderCSSFilter("contrast(200%)", filterTestIDGen)
	if !strings.Contains(defs, `slope="2"`) {
		t.Errorf("expected slope=2 for 200%%, got %s", defs)
	}
}
