package style

import (
	"testing"
)

func TestParseBoxShadowBasic(t *testing.T) {
	shadows, err := ParseBoxShadow("10px 10px 5px rgba(0,0,0,0.5)")
	if err != nil {
		t.Fatal(err)
	}
	if len(shadows) != 1 {
		t.Fatalf("got %d shadows, want 1", len(shadows))
	}
	s := shadows[0]
	if s.OffsetX != 10 || s.OffsetY != 10 || s.Blur != 5 {
		t.Errorf("offsets/blur = (%v, %v, %v), want (10, 10, 5)", s.OffsetX, s.OffsetY, s.Blur)
	}
	if s.Color.A != 0.5 {
		t.Errorf("alpha = %v, want 0.5", s.Color.A)
	}
	if s.Inset {
		t.Error("inset = true, want false")
	}
}

func TestParseBoxShadowWithSpread(t *testing.T) {
	shadows, err := ParseBoxShadow("10px 10px 5px 2px red")
	if err != nil {
		t.Fatal(err)
	}
	if len(shadows) != 1 {
		t.Fatalf("got %d shadows, want 1", len(shadows))
	}
	s := shadows[0]
	if s.Spread != 2 {
		t.Errorf("spread = %v, want 2", s.Spread)
	}
	if s.Color.R != 255 || s.Color.G != 0 || s.Color.B != 0 {
		t.Errorf("color = %v, want red", s.Color)
	}
}

func TestParseBoxShadowInset(t *testing.T) {
	shadows, err := ParseBoxShadow("inset 10px 10px 5px red")
	if err != nil {
		t.Fatal(err)
	}
	if len(shadows) != 1 {
		t.Fatalf("got %d shadows, want 1", len(shadows))
	}
	if !shadows[0].Inset {
		t.Error("inset = false, want true")
	}
	if shadows[0].OffsetX != 10 || shadows[0].OffsetY != 10 {
		t.Errorf("offsets = (%v, %v), want (10, 10)", shadows[0].OffsetX, shadows[0].OffsetY)
	}
}

func TestParseBoxShadowMultiple(t *testing.T) {
	shadows, err := ParseBoxShadow("10px 10px red, 20px 20px blue")
	if err != nil {
		t.Fatal(err)
	}
	if len(shadows) != 2 {
		t.Fatalf("got %d shadows, want 2", len(shadows))
	}
	if shadows[0].OffsetX != 10 {
		t.Errorf("shadow[0] offsetX = %v, want 10", shadows[0].OffsetX)
	}
	if shadows[1].OffsetX != 20 {
		t.Errorf("shadow[1] offsetX = %v, want 20", shadows[1].OffsetX)
	}
	if shadows[0].Color.R != 255 {
		t.Errorf("shadow[0] color = %v, want red", shadows[0].Color)
	}
	if shadows[1].Color.B != 255 {
		t.Errorf("shadow[1] color = %v, want blue", shadows[1].Color)
	}
}

func TestParseBoxShadowMultipleWithParens(t *testing.T) {
	shadows, err := ParseBoxShadow("10px 10px rgba(0,0,0,0.5), 20px 20px rgba(255,0,0,1)")
	if err != nil {
		t.Fatal(err)
	}
	if len(shadows) != 2 {
		t.Fatalf("got %d shadows, want 2", len(shadows))
	}
}

func TestParseBoxShadowNone(t *testing.T) {
	shadows, err := ParseBoxShadow("none")
	if err != nil {
		t.Fatal(err)
	}
	if shadows != nil {
		t.Errorf("got %v, want nil", shadows)
	}
}

func TestParseBoxShadowColorAtStart(t *testing.T) {
	shadows, err := ParseBoxShadow("red 10px 10px 5px")
	if err != nil {
		t.Fatal(err)
	}
	if len(shadows) != 1 {
		t.Fatalf("got %d shadows, want 1", len(shadows))
	}
	s := shadows[0]
	if s.Color.R != 255 || s.Color.G != 0 || s.Color.B != 0 {
		t.Errorf("color = %v, want red", s.Color)
	}
	if s.OffsetX != 10 || s.OffsetY != 10 || s.Blur != 5 {
		t.Errorf("values = (%v, %v, %v), want (10, 10, 5)", s.OffsetX, s.OffsetY, s.Blur)
	}
}

func TestParseBoxShadowZeroValues(t *testing.T) {
	shadows, err := ParseBoxShadow("0 0 black")
	if err != nil {
		t.Fatal(err)
	}
	if len(shadows) != 1 {
		t.Fatalf("got %d shadows, want 1", len(shadows))
	}
	if shadows[0].OffsetX != 0 || shadows[0].OffsetY != 0 {
		t.Errorf("offsets = (%v, %v), want (0, 0)", shadows[0].OffsetX, shadows[0].OffsetY)
	}
}

func TestParseTextShadowBasic(t *testing.T) {
	shadows, err := ParseTextShadow("2px 2px 4px red")
	if err != nil {
		t.Fatal(err)
	}
	if len(shadows) != 1 {
		t.Fatalf("got %d shadows, want 1", len(shadows))
	}
	s := shadows[0]
	if s.OffsetX != 2 || s.OffsetY != 2 || s.Blur != 4 {
		t.Errorf("values = (%v, %v, %v), want (2, 2, 4)", s.OffsetX, s.OffsetY, s.Blur)
	}
}

func TestParseTextShadowNoBlur(t *testing.T) {
	shadows, err := ParseTextShadow("1px 1px red")
	if err != nil {
		t.Fatal(err)
	}
	if len(shadows) != 1 {
		t.Fatalf("got %d shadows, want 1", len(shadows))
	}
	if shadows[0].Blur != 0 {
		t.Errorf("blur = %v, want 0", shadows[0].Blur)
	}
}

func TestParseTextShadowMultiple(t *testing.T) {
	shadows, err := ParseTextShadow("1px 1px red, 2px 2px blue")
	if err != nil {
		t.Fatal(err)
	}
	if len(shadows) != 2 {
		t.Fatalf("got %d shadows, want 2", len(shadows))
	}
}

func TestParseTextShadowNone(t *testing.T) {
	shadows, err := ParseTextShadow("none")
	if err != nil {
		t.Fatal(err)
	}
	if shadows != nil {
		t.Errorf("got %v, want nil", shadows)
	}
}
