package render

import (
	"fmt"
	"strings"
	"testing"

	"github.com/macawls/ogre/style"
)

func testIDGen(prefix string) string {
	return prefix + "1"
}

func TestSolidColor(t *testing.T) {
	cs := &style.ComputedStyle{
		BackgroundColor: style.Color{R: 255, G: 0, B: 0, A: 1},
	}
	r := RenderBackground(cs, 0, 0, 100, 100, testIDGen)
	if r.Defs != "" {
		t.Errorf("expected empty defs, got %q", r.Defs)
	}
	if r.Fill != "#ff0000" {
		t.Errorf("expected #ff0000, got %q", r.Fill)
	}
}

func TestNoBackground(t *testing.T) {
	cs := &style.ComputedStyle{
		BackgroundColor: style.Color{R: 0, G: 0, B: 0, A: 0},
	}
	r := RenderBackground(cs, 0, 0, 100, 100, testIDGen)
	if r.Fill != "none" {
		t.Errorf("expected none, got %q", r.Fill)
	}
}

func TestLinearGradientWithAngle(t *testing.T) {
	cs := &style.ComputedStyle{
		BackgroundImage: "linear-gradient(90deg, red, blue)",
	}
	r := RenderBackground(cs, 0, 0, 100, 100, testIDGen)
	if !strings.Contains(r.Fill, "url(#lg1)") {
		t.Errorf("expected url(#lg1), got %q", r.Fill)
	}
	if !strings.Contains(r.Defs, "<linearGradient") {
		t.Errorf("expected linearGradient in defs, got %q", r.Defs)
	}
	if !strings.Contains(r.Defs, `x1="0%"`) {
		t.Errorf("expected x1=0%% for 90deg, got %q", r.Defs)
	}
	if !strings.Contains(r.Defs, `x2="100%"`) {
		t.Errorf("expected x2=100%% for 90deg, got %q", r.Defs)
	}
	if !strings.Contains(r.Defs, "stop-color") {
		t.Errorf("expected stop-color in defs, got %q", r.Defs)
	}
}

func TestLinearGradientWithDirection(t *testing.T) {
	cs := &style.ComputedStyle{
		BackgroundImage: "linear-gradient(to right, #ff0000, #0000ff)",
	}
	r := RenderBackground(cs, 0, 0, 100, 100, testIDGen)
	if !strings.Contains(r.Defs, "<linearGradient") {
		t.Errorf("expected linearGradient in defs, got %q", r.Defs)
	}
	if !strings.Contains(r.Defs, `x1="0%"`) {
		t.Errorf("expected x1=0%% for to right, got %q", r.Defs)
	}
}

func TestLinearGradientToTop(t *testing.T) {
	cs := &style.ComputedStyle{
		BackgroundImage: "linear-gradient(to top, red, blue)",
	}
	r := RenderBackground(cs, 0, 0, 100, 100, testIDGen)
	if !strings.Contains(r.Defs, `y1="100%"`) {
		t.Errorf("expected y1=100%% for to top, got %q", r.Defs)
	}
	if !strings.Contains(r.Defs, `y2="0%"`) {
		t.Errorf("expected y2=0%% for to top, got %q", r.Defs)
	}
}

func TestRadialGradient(t *testing.T) {
	cs := &style.ComputedStyle{
		BackgroundImage: "radial-gradient(circle, red, blue)",
	}
	r := RenderBackground(cs, 0, 0, 100, 100, testIDGen)
	if !strings.Contains(r.Fill, "url(#rg1)") {
		t.Errorf("expected url(#rg1), got %q", r.Fill)
	}
	if !strings.Contains(r.Defs, "<radialGradient") {
		t.Errorf("expected radialGradient in defs, got %q", r.Defs)
	}
	if !strings.Contains(r.Defs, `cx="50%"`) {
		t.Errorf("expected cx=50%%, got %q", r.Defs)
	}
}

func TestColorStopsWithPositions(t *testing.T) {
	cs := &style.ComputedStyle{
		BackgroundImage: "linear-gradient(90deg, red 0%, green 50%, blue 100%)",
	}
	r := RenderBackground(cs, 0, 0, 100, 100, testIDGen)
	if !strings.Contains(r.Defs, `offset="0%"`) {
		t.Errorf("expected offset=0%%, got %q", r.Defs)
	}
	if !strings.Contains(r.Defs, `offset="50%"`) {
		t.Errorf("expected offset=50%%, got %q", r.Defs)
	}
	if !strings.Contains(r.Defs, `offset="100%"`) {
		t.Errorf("expected offset=100%%, got %q", r.Defs)
	}
}

func TestColorStopsWithoutPositions(t *testing.T) {
	cs := &style.ComputedStyle{
		BackgroundImage: "linear-gradient(90deg, red, green, blue)",
	}
	r := RenderBackground(cs, 0, 0, 100, 100, testIDGen)
	if !strings.Contains(r.Defs, `offset="0%"`) {
		t.Errorf("expected offset=0%% for first stop, got %q", r.Defs)
	}
	if !strings.Contains(r.Defs, `offset="50%"`) {
		t.Errorf("expected offset=50%% for middle stop, got %q", r.Defs)
	}
	if !strings.Contains(r.Defs, `offset="100%"`) {
		t.Errorf("expected offset=100%% for last stop, got %q", r.Defs)
	}
}

func TestColorStopsMixedPositions(t *testing.T) {
	cs := &style.ComputedStyle{
		BackgroundImage: "linear-gradient(90deg, red 0%, green, yellow, blue 100%)",
	}
	r := RenderBackground(cs, 0, 0, 100, 100, testIDGen)
	if !strings.Contains(r.Defs, `offset="0%"`) {
		t.Errorf("expected offset=0%%, got %q", r.Defs)
	}
	expected33 := fmt.Sprintf(`offset="%.6g%%"`, 100.0/3)
	expected66 := fmt.Sprintf(`offset="%.6g%%"`, 200.0/3)
	if !strings.Contains(r.Defs, expected33) {
		t.Errorf("expected offset ~33.33%% for interpolated stop, got %q", r.Defs)
	}
	if !strings.Contains(r.Defs, expected66) {
		t.Errorf("expected offset ~66.67%% for interpolated stop, got %q", r.Defs)
	}
}

func TestURLBackgroundGeneratesPattern(t *testing.T) {
	cs := &style.ComputedStyle{
		BackgroundImage: `url('https://example.com/image.png')`,
	}
	r := RenderBackground(cs, 10, 20, 200, 100, testIDGen)
	if !strings.Contains(r.Defs, "<pattern") {
		t.Errorf("expected <pattern in defs, got %q", r.Defs)
	}
	if !strings.Contains(r.Defs, `<image href="https://example.com/image.png"`) {
		t.Errorf("expected image href in defs, got %q", r.Defs)
	}
	if !strings.Contains(r.Fill, "url(#bg1)") {
		t.Errorf("expected url(#bg1), got %q", r.Fill)
	}
}

func TestURLBackgroundSizeCover(t *testing.T) {
	cs := &style.ComputedStyle{
		BackgroundImage: `url("bg.jpg")`,
		BackgroundSize:  "cover",
	}
	r := RenderBackground(cs, 0, 0, 300, 200, testIDGen)
	if !strings.Contains(r.Defs, `preserveAspectRatio="xMidYMid slice"`) {
		t.Errorf("expected xMidYMid slice for cover, got %q", r.Defs)
	}
}

func TestURLBackgroundSizeContain(t *testing.T) {
	cs := &style.ComputedStyle{
		BackgroundImage: `url(bg.jpg)`,
		BackgroundSize:  "contain",
	}
	r := RenderBackground(cs, 0, 0, 300, 200, testIDGen)
	if !strings.Contains(r.Defs, `preserveAspectRatio="xMidYMid meet"`) {
		t.Errorf("expected xMidYMid meet for contain, got %q", r.Defs)
	}
}

func TestURLBackgroundRepeatNoRepeat(t *testing.T) {
	cs := &style.ComputedStyle{
		BackgroundImage:  `url(bg.jpg)`,
		BackgroundRepeat: "no-repeat",
	}
	r := RenderBackground(cs, 0, 0, 300, 200, testIDGen)
	if !strings.Contains(r.Defs, `width="300"`) {
		t.Errorf("expected pattern width=300 for no-repeat, got %q", r.Defs)
	}
	if !strings.Contains(r.Defs, `height="200"`) {
		t.Errorf("expected pattern height=200 for no-repeat, got %q", r.Defs)
	}
}

func TestURLBackgroundPositionCenter(t *testing.T) {
	cs := &style.ComputedStyle{
		BackgroundImage:    `url(bg.jpg)`,
		BackgroundPosition: "center",
		BackgroundRepeat:   "no-repeat",
	}
	r := RenderBackground(cs, 0, 0, 400, 200, testIDGen)
	if !strings.Contains(r.Defs, `x="0"`) {
		t.Errorf("expected x=0 for centered no-repeat (pattern=element size), got %q", r.Defs)
	}
}

func TestURLBackgroundExplicitSize(t *testing.T) {
	cs := &style.ComputedStyle{
		BackgroundImage: `url(bg.jpg)`,
		BackgroundSize:  "50px 80px",
	}
	r := RenderBackground(cs, 0, 0, 200, 200, testIDGen)
	if !strings.Contains(r.Defs, `width="50"`) {
		t.Errorf("expected image width=50, got %q", r.Defs)
	}
	if !strings.Contains(r.Defs, `height="80"`) {
		t.Errorf("expected image height=80, got %q", r.Defs)
	}
}

func TestURLBackgroundPercentSize(t *testing.T) {
	cs := &style.ComputedStyle{
		BackgroundImage: `url(bg.jpg)`,
		BackgroundSize:  "50% 25%",
	}
	r := RenderBackground(cs, 0, 0, 200, 400, testIDGen)
	if !strings.Contains(r.Defs, `width="100"`) {
		t.Errorf("expected image width=100 (50%% of 200), got %q", r.Defs)
	}
	if !strings.Contains(r.Defs, `height="100"`) {
		t.Errorf("expected image height=100 (25%% of 400), got %q", r.Defs)
	}
}

func TestExtractURL(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{`url('https://example.com/img.png')`, "https://example.com/img.png"},
		{`url("https://example.com/img.png")`, "https://example.com/img.png"},
		{`url(https://example.com/img.png)`, "https://example.com/img.png"},
	}
	for _, tt := range tests {
		got := extractURL(tt.input)
		if got != tt.want {
			t.Errorf("extractURL(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestDistributeStops(t *testing.T) {
	stops := []style.ColorStop{
		{Color: style.Color{R: 255, A: 1}, HasPos: false},
		{Color: style.Color{G: 255, A: 1}, HasPos: false},
		{Color: style.Color{B: 255, A: 1}, HasPos: false},
	}
	distributeStops(stops)
	if stops[0].Position != 0 {
		t.Errorf("first stop should be 0, got %f", stops[0].Position)
	}
	if stops[1].Position != 0.5 {
		t.Errorf("middle stop should be 0.5, got %f", stops[1].Position)
	}
	if stops[2].Position != 1 {
		t.Errorf("last stop should be 1, got %f", stops[2].Position)
	}
}

func TestDistributeStopsPartial(t *testing.T) {
	stops := []style.ColorStop{
		{Color: style.Color{R: 255, A: 1}, Position: 0, HasPos: true},
		{Color: style.Color{G: 128, A: 1}, HasPos: false},
		{Color: style.Color{G: 255, A: 1}, HasPos: false},
		{Color: style.Color{B: 255, A: 1}, Position: 1, HasPos: true},
	}
	distributeStops(stops)
	if stops[1].Position != 1.0/3 {
		t.Errorf("expected 1/3, got %f", stops[1].Position)
	}
	if stops[2].Position != 2.0/3 {
		t.Errorf("expected 2/3, got %f", stops[2].Position)
	}
}
