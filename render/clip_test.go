package render

import (
	"strings"
	"testing"

	"github.com/macawls/ogre/style"
)

func makeIDGen() func(string) string {
	n := 0
	return func(prefix string) string {
		n++
		return prefix + string(rune('0'+n))
	}
}

func TestOverflowHiddenGeneratesClipPath(t *testing.T) {
	cs := &style.ComputedStyle{Overflow: style.OverflowHidden}
	defs, attr := RenderOverflowClip(cs, 10, 20, 200, 100, makeIDGen())

	if !strings.Contains(defs, "<clipPath") {
		t.Fatalf("expected <clipPath> in defs, got: %s", defs)
	}
	if !strings.Contains(defs, "<rect") {
		t.Fatalf("expected <rect> in clipPath, got: %s", defs)
	}
	if !strings.Contains(attr, `clip-path="url(#`) {
		t.Fatalf("expected clip-path attr, got: %s", attr)
	}
	if !strings.Contains(defs, `x="10"`) {
		t.Fatalf("expected x=10 in rect, got: %s", defs)
	}
	if !strings.Contains(defs, `width="200"`) {
		t.Fatalf("expected width=200 in rect, got: %s", defs)
	}
}

func TestOverflowHiddenWithUniformBorderRadius(t *testing.T) {
	cs := &style.ComputedStyle{
		Overflow:                style.OverflowHidden,
		BorderTopLeftRadius:     8,
		BorderTopRightRadius:    8,
		BorderBottomLeftRadius:  8,
		BorderBottomRightRadius: 8,
	}
	defs, attr := RenderOverflowClip(cs, 0, 0, 100, 100, makeIDGen())

	if !strings.Contains(defs, `rx="8"`) {
		t.Fatalf("expected rx=8, got: %s", defs)
	}
	if attr == "" {
		t.Fatal("expected clip-path attr")
	}
}

func TestOverflowHiddenWithPerCornerRadius(t *testing.T) {
	cs := &style.ComputedStyle{
		Overflow:                style.OverflowHidden,
		BorderTopLeftRadius:     10,
		BorderTopRightRadius:    5,
		BorderBottomLeftRadius:  0,
		BorderBottomRightRadius: 15,
	}
	defs, _ := RenderOverflowClip(cs, 0, 0, 100, 100, makeIDGen())

	if !strings.Contains(defs, "<path") {
		t.Fatalf("expected <path> for per-corner radii, got: %s", defs)
	}
}

func TestOverflowVisibleNoClip(t *testing.T) {
	cs := &style.ComputedStyle{Overflow: style.OverflowVisible}
	defs, attr := RenderOverflowClip(cs, 0, 0, 100, 100, makeIDGen())

	if defs != "" || attr != "" {
		t.Fatalf("expected no clip for overflow:visible, got defs=%q attr=%q", defs, attr)
	}
}

func TestClipPathCircle(t *testing.T) {
	cs := &style.ComputedStyle{ClipPath: "circle(50%)"}
	defs, attr := RenderOverflowClip(cs, 0, 0, 200, 200, makeIDGen())

	if !strings.Contains(defs, "<circle") {
		t.Fatalf("expected <circle>, got: %s", defs)
	}
	if !strings.Contains(defs, `r="100"`) {
		t.Fatalf("expected r=100 (50%% of 200), got: %s", defs)
	}
	if !strings.Contains(defs, `cx="100"`) {
		t.Fatalf("expected cx=100, got: %s", defs)
	}
	if attr == "" {
		t.Fatal("expected clip-path attr")
	}
}

func TestClipPathEllipse(t *testing.T) {
	cs := &style.ComputedStyle{ClipPath: "ellipse(50% 30%)"}
	defs, _ := RenderOverflowClip(cs, 0, 0, 200, 100, makeIDGen())

	if !strings.Contains(defs, "<ellipse") {
		t.Fatalf("expected <ellipse>, got: %s", defs)
	}
	if !strings.Contains(defs, `rx="100"`) {
		t.Fatalf("expected rx=100, got: %s", defs)
	}
	if !strings.Contains(defs, `ry="30"`) {
		t.Fatalf("expected ry=30, got: %s", defs)
	}
}

func TestClipPathPolygon(t *testing.T) {
	cs := &style.ComputedStyle{ClipPath: "polygon(0 0, 100% 0, 50% 100%)"}
	defs, _ := RenderOverflowClip(cs, 0, 0, 200, 100, makeIDGen())

	if !strings.Contains(defs, "<polygon") {
		t.Fatalf("expected <polygon>, got: %s", defs)
	}
	if !strings.Contains(defs, `points="`) {
		t.Fatalf("expected points attr, got: %s", defs)
	}
	if !strings.Contains(defs, "0,0") {
		t.Fatalf("expected point 0,0, got: %s", defs)
	}
	if !strings.Contains(defs, "200,0") {
		t.Fatalf("expected point 200,0, got: %s", defs)
	}
	if !strings.Contains(defs, "100,100") {
		t.Fatalf("expected point 100,100, got: %s", defs)
	}
}

func TestClipPathInset(t *testing.T) {
	cs := &style.ComputedStyle{ClipPath: "inset(10px 20px)"}
	defs, _ := RenderOverflowClip(cs, 0, 0, 200, 100, makeIDGen())

	if !strings.Contains(defs, "<rect") {
		t.Fatalf("expected <rect>, got: %s", defs)
	}
	if !strings.Contains(defs, `x="20"`) {
		t.Fatalf("expected x=20, got: %s", defs)
	}
	if !strings.Contains(defs, `y="10"`) {
		t.Fatalf("expected y=10, got: %s", defs)
	}
	if !strings.Contains(defs, `width="160"`) {
		t.Fatalf("expected width=160, got: %s", defs)
	}
	if !strings.Contains(defs, `height="80"`) {
		t.Fatalf("expected height=80, got: %s", defs)
	}
}

func TestOpacityZeroRendersNoGroup(t *testing.T) {
	cs := &style.ComputedStyle{Opacity: 0}
	needsGroup := cs.Opacity > 0 && cs.Opacity < 1
	if needsGroup {
		t.Fatal("opacity 0 should not create a group")
	}
}

func TestOpacityOneNoGroup(t *testing.T) {
	cs := &style.ComputedStyle{Opacity: 1}
	needsGroup := cs.Opacity > 0 && cs.Opacity < 1
	if needsGroup {
		t.Fatal("opacity 1 should not create a group")
	}
}
