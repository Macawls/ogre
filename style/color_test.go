package style

import (
	"math"
	"testing"
)

func colorsEqual(a, b Color) bool {
	return a.R == b.R && a.G == b.G && a.B == b.B && math.Abs(a.A-b.A) < 0.002
}

func TestParseNamedColors(t *testing.T) {
	tests := []struct {
		input string
		want  Color
	}{
		{"red", Color{255, 0, 0, 1}},
		{"Blue", Color{0, 0, 255, 1}},
		{"GREEN", Color{0, 128, 0, 1}},
		{"aliceblue", Color{240, 248, 255, 1}},
		{"yellowgreen", Color{154, 205, 50, 1}},
		{"rebeccapurple", Color{102, 51, 153, 1}},
		{"darkslategray", Color{47, 79, 79, 1}},
		{"darkslategrey", Color{47, 79, 79, 1}},
	}
	for _, tt := range tests {
		got, err := ParseColor(tt.input)
		if err != nil {
			t.Errorf("ParseColor(%q) error: %v", tt.input, err)
			continue
		}
		if !colorsEqual(got, tt.want) {
			t.Errorf("ParseColor(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestParseTransparent(t *testing.T) {
	c, err := ParseColor("transparent")
	if err != nil {
		t.Fatal(err)
	}
	if c.R != 0 || c.G != 0 || c.B != 0 || c.A != 0 {
		t.Errorf("transparent = %v, want rgba(0,0,0,0)", c)
	}
	if !c.IsTransparent() {
		t.Error("transparent should be transparent")
	}
}

func TestParseCurrentColor(t *testing.T) {
	c, err := ParseColor("currentColor")
	if err != nil {
		t.Fatal(err)
	}
	if c != CurrentColor {
		t.Errorf("currentColor = %v, want sentinel %v", c, CurrentColor)
	}

	c2, err := ParseColor("CURRENTCOLOR")
	if err != nil {
		t.Fatal(err)
	}
	if c2 != CurrentColor {
		t.Errorf("CURRENTCOLOR = %v, want sentinel", c2)
	}
}

func TestParseHex(t *testing.T) {
	tests := []struct {
		input string
		want  Color
	}{
		{"#fff", Color{255, 255, 255, 1}},
		{"#000", Color{0, 0, 0, 1}},
		{"#f00", Color{255, 0, 0, 1}},
		{"#ff0000", Color{255, 0, 0, 1}},
		{"#00ff00", Color{0, 255, 0, 1}},
		{"#0000ff", Color{0, 0, 255, 1}},
		{"#FF0000", Color{255, 0, 0, 1}},
		{"#f00f", Color{255, 0, 0, 1}},
		{"#f008", Color{255, 0, 0, 0.533}},
		{"#ff000080", Color{255, 0, 0, 0.502}},
		{"#ff0000ff", Color{255, 0, 0, 1}},
		{"#ff000000", Color{255, 0, 0, 0}},
	}
	for _, tt := range tests {
		got, err := ParseColor(tt.input)
		if err != nil {
			t.Errorf("ParseColor(%q) error: %v", tt.input, err)
			continue
		}
		if !colorsEqual(got, tt.want) {
			t.Errorf("ParseColor(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestParseRGB(t *testing.T) {
	tests := []struct {
		input string
		want  Color
	}{
		{"rgb(255, 0, 0)", Color{255, 0, 0, 1}},
		{"rgb(0,128,0)", Color{0, 128, 0, 1}},
		{"rgb(100%, 0%, 0%)", Color{255, 0, 0, 1}},
		{"rgb(50%, 50%, 50%)", Color{128, 128, 128, 1}},
		{"rgba(255, 0, 0, 0.5)", Color{255, 0, 0, 0.5}},
		{"rgba(255, 0, 0, 1)", Color{255, 0, 0, 1}},
		{"rgba(255, 0, 0, 0)", Color{255, 0, 0, 0}},
		{"rgb(255 0 0)", Color{255, 0, 0, 1}},
		{"rgb(255 0 0 / 0.5)", Color{255, 0, 0, 0.5}},
		{"rgb(255 0 0 / 50%)", Color{255, 0, 0, 0.5}},
		{"rgba(100%, 0%, 0%, 0.5)", Color{255, 0, 0, 0.5}},
	}
	for _, tt := range tests {
		got, err := ParseColor(tt.input)
		if err != nil {
			t.Errorf("ParseColor(%q) error: %v", tt.input, err)
			continue
		}
		if !colorsEqual(got, tt.want) {
			t.Errorf("ParseColor(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestParseHSL(t *testing.T) {
	tests := []struct {
		input string
		want  Color
	}{
		{"hsl(0, 100%, 50%)", Color{255, 0, 0, 1}},
		{"hsl(120, 100%, 50%)", Color{0, 255, 0, 1}},
		{"hsl(240, 100%, 50%)", Color{0, 0, 255, 1}},
		{"hsl(0, 0%, 0%)", Color{0, 0, 0, 1}},
		{"hsl(0, 0%, 100%)", Color{255, 255, 255, 1}},
		{"hsl(0, 0%, 50%)", Color{128, 128, 128, 1}},
		{"hsla(120, 100%, 50%, 0.5)", Color{0, 255, 0, 0.5}},
		{"hsl(120 100% 50%)", Color{0, 255, 0, 1}},
		{"hsl(120 100% 50% / 0.5)", Color{0, 255, 0, 0.5}},
		{"hsl(120 100% 50% / 50%)", Color{0, 255, 0, 0.5}},
		{"hsl(60, 100%, 50%)", Color{255, 255, 0, 1}},
		{"hsl(180, 100%, 50%)", Color{0, 255, 255, 1}},
		{"hsl(300, 100%, 50%)", Color{255, 0, 255, 1}},
		{"hsl(-120, 100%, 50%)", Color{0, 0, 255, 1}},
		{"hsl(480, 100%, 50%)", Color{0, 255, 0, 1}},
	}
	for _, tt := range tests {
		got, err := ParseColor(tt.input)
		if err != nil {
			t.Errorf("ParseColor(%q) error: %v", tt.input, err)
			continue
		}
		if !colorsEqual(got, tt.want) {
			t.Errorf("ParseColor(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestParseOKLCH(t *testing.T) {
	// Chromatic anchors verified against Tailwind v4's palette (converted from
	// their published OKLCH definitions). Tolerance of ±3 per channel absorbs
	// gamma-encoding rounding variance across implementations.
	tests := []struct {
		input   string
		wantR   uint8
		wantG   uint8
		wantB   uint8
		wantA   float64
		wantErr bool
	}{
		{input: "oklch(100% 0 0)", wantR: 255, wantG: 255, wantB: 255, wantA: 1},
		{input: "oklch(0% 0 0)", wantR: 0, wantG: 0, wantB: 0, wantA: 1},
		{input: "oklch(55.6% 0 none)", wantR: 115, wantG: 115, wantB: 115, wantA: 1}, // v4 neutral-500
		{input: "oklch(62.3% 0.214 259.815)", wantR: 43, wantG: 127, wantB: 255, wantA: 1}, // v4 blue-500
		{input: "oklch(63.7% 0.237 25.331)", wantR: 251, wantG: 44, wantB: 54, wantA: 1},   // v4 red-500
		{input: "oklch(72.3% 0.219 149.579)", wantR: 0, wantG: 201, wantB: 80, wantA: 1},   // v4 green-500
		{input: "oklch(55.4% 0.046 257.417 / 0.5)", wantR: 98, wantG: 116, wantB: 142, wantA: 0.5}, // v4 slate-500
		{input: "oklch(0.554 0.046 257.417)", wantR: 98, wantG: 116, wantB: 142, wantA: 1},
	}
	for _, tt := range tests {
		got, err := ParseColor(tt.input)
		if tt.wantErr {
			if err == nil {
				t.Errorf("ParseColor(%q) expected error, got %v", tt.input, got)
			}
			continue
		}
		if err != nil {
			t.Errorf("ParseColor(%q) error: %v", tt.input, err)
			continue
		}
		if !nearU8(got.R, tt.wantR, 3) || !nearU8(got.G, tt.wantG, 3) || !nearU8(got.B, tt.wantB, 3) {
			t.Errorf("ParseColor(%q) = rgb(%d,%d,%d), want ~rgb(%d,%d,%d)",
				tt.input, got.R, got.G, got.B, tt.wantR, tt.wantG, tt.wantB)
		}
		if math.Abs(got.A-tt.wantA) > 0.01 {
			t.Errorf("ParseColor(%q) alpha = %g, want %g", tt.input, got.A, tt.wantA)
		}
	}
}

func nearU8(a, b, tol uint8) bool {
	d := int(a) - int(b)
	if d < 0 {
		d = -d
	}
	return d <= int(tol)
}

func TestColorString(t *testing.T) {
	c := Color{255, 0, 0, 1}
	if s := c.String(); s != "rgba(255,0,0,1)" {
		t.Errorf("String() = %q, want %q", s, "rgba(255,0,0,1)")
	}

	c2 := Color{0, 128, 255, 0.5}
	if s := c2.String(); s != "rgba(0,128,255,0.5)" {
		t.Errorf("String() = %q, want %q", s, "rgba(0,128,255,0.5)")
	}
}

func TestColorHex(t *testing.T) {
	tests := []struct {
		color Color
		want  string
	}{
		{Color{255, 0, 0, 1}, "#ff0000"},
		{Color{0, 128, 255, 1}, "#0080ff"},
		{Color{255, 0, 0, 0.5}, "#ff000080"},
		{Color{0, 0, 0, 0}, "#00000000"},
	}
	for _, tt := range tests {
		if got := tt.color.Hex(); got != tt.want {
			t.Errorf("%v.Hex() = %q, want %q", tt.color, got, tt.want)
		}
	}
}

func TestIsTransparent(t *testing.T) {
	if !(Color{0, 0, 0, 0}).IsTransparent() {
		t.Error("alpha 0 should be transparent")
	}
	if (Color{0, 0, 0, 1}).IsTransparent() {
		t.Error("alpha 1 should not be transparent")
	}
	if (Color{0, 0, 0, 0.5}).IsTransparent() {
		t.Error("alpha 0.5 should not be transparent")
	}
}

func TestParseWhitespace(t *testing.T) {
	c, err := ParseColor("  red  ")
	if err != nil {
		t.Fatal(err)
	}
	if c.R != 255 || c.G != 0 || c.B != 0 {
		t.Errorf("trimmed color = %v, want red", c)
	}
}

func TestParseInvalid(t *testing.T) {
	invalids := []string{
		"",
		"notacolor",
		"#gg0000",
		"#12345",
		"rgb()",
		"rgb(1,2)",
		"hsl()",
	}
	for _, s := range invalids {
		_, err := ParseColor(s)
		if err == nil {
			t.Errorf("ParseColor(%q) should have returned an error", s)
		}
	}
}

func TestNamedColorCount(t *testing.T) {
	if len(namedColors) < 148 {
		t.Errorf("expected at least 148 named colors, got %d", len(namedColors))
	}
}

func TestClampValues(t *testing.T) {
	c, err := ParseColor("rgb(300, -10, 128)")
	if err != nil {
		t.Fatal(err)
	}
	if c.R != 255 {
		t.Errorf("R should be clamped to 255, got %d", c.R)
	}
	if c.G != 0 {
		t.Errorf("G should be clamped to 0, got %d", c.G)
	}
}
