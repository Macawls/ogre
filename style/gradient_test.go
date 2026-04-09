package style

import (
	"math"
	"testing"
)

func TestParseLinearGradientDirectionKeywords(t *testing.T) {
	tests := []struct {
		input string
		angle float64
	}{
		{"linear-gradient(to right, red, blue)", 90},
		{"linear-gradient(to left, red, blue)", 270},
		{"linear-gradient(to top, red, blue)", 0},
		{"linear-gradient(to bottom, red, blue)", 180},
		{"linear-gradient(to top right, red, blue)", 45},
		{"linear-gradient(to bottom left, red, blue)", 225},
	}

	for _, tt := range tests {
		g, err := ParseGradient(tt.input)
		if err != nil {
			t.Errorf("ParseGradient(%q): %v", tt.input, err)
			continue
		}
		if g.Angle != tt.angle {
			t.Errorf("ParseGradient(%q): angle = %v, want %v", tt.input, g.Angle, tt.angle)
		}
		if g.Type != LinearGradient {
			t.Errorf("ParseGradient(%q): type = %v, want LinearGradient", tt.input, g.Type)
		}
		if len(g.Stops) != 2 {
			t.Errorf("ParseGradient(%q): stops = %d, want 2", tt.input, len(g.Stops))
		}
	}
}

func TestParseLinearGradientAngle(t *testing.T) {
	g, err := ParseGradient("linear-gradient(45deg, red 0%, blue 100%)")
	if err != nil {
		t.Fatal(err)
	}
	if g.Angle != 45 {
		t.Errorf("angle = %v, want 45", g.Angle)
	}
	if len(g.Stops) != 2 {
		t.Fatalf("stops = %d, want 2", len(g.Stops))
	}
	if !g.Stops[0].HasPos || g.Stops[0].Position != 0 {
		t.Errorf("stop[0] position = %v (has=%v), want 0 (has=true)", g.Stops[0].Position, g.Stops[0].HasPos)
	}
	if !g.Stops[1].HasPos || g.Stops[1].Position != 1.0 {
		t.Errorf("stop[1] position = %v (has=%v), want 1 (has=true)", g.Stops[1].Position, g.Stops[1].HasPos)
	}
}

func TestParseLinearGradientAutoStops(t *testing.T) {
	g, err := ParseGradient("linear-gradient(red, green 50%, blue)")
	if err != nil {
		t.Fatal(err)
	}
	if g.Angle != 180 {
		t.Errorf("angle = %v, want 180 (default)", g.Angle)
	}
	if len(g.Stops) != 3 {
		t.Fatalf("stops = %d, want 3", len(g.Stops))
	}
	if g.Stops[0].HasPos {
		t.Errorf("stop[0] should not have position")
	}
	if !g.Stops[1].HasPos || g.Stops[1].Position != 0.5 {
		t.Errorf("stop[1] position = %v, want 0.5", g.Stops[1].Position)
	}
	if g.Stops[2].HasPos {
		t.Errorf("stop[2] should not have position")
	}
}

func TestParseLinearGradientDefaultAngle(t *testing.T) {
	g, err := ParseGradient("linear-gradient(red, blue)")
	if err != nil {
		t.Fatal(err)
	}
	if g.Angle != 180 {
		t.Errorf("angle = %v, want 180", g.Angle)
	}
}

func TestParseRadialGradientCircle(t *testing.T) {
	g, err := ParseGradient("radial-gradient(circle at center, red, blue)")
	if err != nil {
		t.Fatal(err)
	}
	if g.Type != RadialGradient {
		t.Errorf("type = %v, want RadialGradient", g.Type)
	}
	if g.Shape != "circle" {
		t.Errorf("shape = %q, want circle", g.Shape)
	}
	if g.PositionX != 50 || g.PositionY != 50 {
		t.Errorf("position = (%v, %v), want (50, 50)", g.PositionX, g.PositionY)
	}
	if len(g.Stops) != 2 {
		t.Errorf("stops = %d, want 2", len(g.Stops))
	}
}

func TestParseRadialGradientEllipse(t *testing.T) {
	g, err := ParseGradient("radial-gradient(ellipse farthest-corner at 50% 50%, red, blue)")
	if err != nil {
		t.Fatal(err)
	}
	if g.Shape != "ellipse" {
		t.Errorf("shape = %q, want ellipse", g.Shape)
	}
	if g.Size != "farthest-corner" {
		t.Errorf("size = %q, want farthest-corner", g.Size)
	}
	if g.PositionX != 50 || g.PositionY != 50 {
		t.Errorf("position = (%v, %v), want (50, 50)", g.PositionX, g.PositionY)
	}
}

func TestParseRepeatingLinearGradient(t *testing.T) {
	g, err := ParseGradient("repeating-linear-gradient(45deg, red, blue)")
	if err != nil {
		t.Fatal(err)
	}
	if g.Type != RepeatingLinearGradient {
		t.Errorf("type = %v, want RepeatingLinearGradient", g.Type)
	}
	if !g.Repeating {
		t.Error("repeating = false, want true")
	}
	if g.Angle != 45 {
		t.Errorf("angle = %v, want 45", g.Angle)
	}
}

func TestParseRepeatingRadialGradient(t *testing.T) {
	g, err := ParseGradient("repeating-radial-gradient(circle, red, blue)")
	if err != nil {
		t.Fatal(err)
	}
	if g.Type != RepeatingRadialGradient {
		t.Errorf("type = %v, want RepeatingRadialGradient", g.Type)
	}
	if !g.Repeating {
		t.Error("repeating = false, want true")
	}
}

func TestParseGradientWithRGBA(t *testing.T) {
	g, err := ParseGradient("linear-gradient(to right, rgba(255,0,0,0.5), rgba(0,0,255,1))")
	if err != nil {
		t.Fatal(err)
	}
	if len(g.Stops) != 2 {
		t.Fatalf("stops = %d, want 2", len(g.Stops))
	}
	if g.Stops[0].Color.R != 255 || g.Stops[0].Color.A != 0.5 {
		t.Errorf("stop[0] color = %v, want rgba(255,0,0,0.5)", g.Stops[0].Color)
	}
}

func TestParseGradientRGBAWithPosition(t *testing.T) {
	g, err := ParseGradient("linear-gradient(rgba(255,0,0,0.5) 25%, blue 75%)")
	if err != nil {
		t.Fatal(err)
	}
	if len(g.Stops) != 2 {
		t.Fatalf("stops = %d, want 2", len(g.Stops))
	}
	if !g.Stops[0].HasPos || g.Stops[0].Position != 0.25 {
		t.Errorf("stop[0] pos = %v (has=%v), want 0.25", g.Stops[0].Position, g.Stops[0].HasPos)
	}
}

func TestParseGradientInvalid(t *testing.T) {
	_, err := ParseGradient("not-a-gradient(red, blue)")
	if err == nil {
		t.Error("expected error for invalid gradient")
	}
}

func TestColorStopNaN(t *testing.T) {
	g, err := ParseGradient("linear-gradient(red, blue)")
	if err != nil {
		t.Fatal(err)
	}
	if !math.IsNaN(g.Stops[0].Position) {
		t.Errorf("stop without position should have NaN, got %v", g.Stops[0].Position)
	}
}
