package style

import (
	"testing"
)

func TestSplitValues(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{"10px", []string{"10px"}},
		{"10px 20px", []string{"10px", "20px"}},
		{"rgb(255, 0, 0)", []string{"rgb(255, 0, 0)"}},
		{"1px solid rgb(255, 0, 0)", []string{"1px", "solid", "rgb(255, 0, 0)"}},
		{"calc(100% - 20px) auto", []string{"calc(100% - 20px)", "auto"}},
		{"'hello world' test", []string{"'hello world'", "test"}},
		{`"hello world" test`, []string{`"hello world"`, "test"}},
		{"", nil},
		{"  10px  20px  ", []string{"10px", "20px"}},
	}
	for _, tt := range tests {
		got := splitValues(tt.input)
		if len(got) != len(tt.want) {
			t.Errorf("splitValues(%q) = %v, want %v", tt.input, got, tt.want)
			continue
		}
		for i := range got {
			if got[i] != tt.want[i] {
				t.Errorf("splitValues(%q)[%d] = %q, want %q", tt.input, i, got[i], tt.want[i])
			}
		}
	}
}

func TestExpandMargin(t *testing.T) {
	tests := []struct {
		name  string
		input map[string]string
		check map[string]string
	}{
		{
			name:  "one value",
			input: map[string]string{"margin": "10px"},
			check: map[string]string{
				"margin-top": "10px", "margin-right": "10px",
				"margin-bottom": "10px", "margin-left": "10px",
			},
		},
		{
			name:  "two values",
			input: map[string]string{"margin": "10px 20px"},
			check: map[string]string{
				"margin-top": "10px", "margin-right": "20px",
				"margin-bottom": "10px", "margin-left": "20px",
			},
		},
		{
			name:  "three values",
			input: map[string]string{"margin": "10px 20px 30px"},
			check: map[string]string{
				"margin-top": "10px", "margin-right": "20px",
				"margin-bottom": "30px", "margin-left": "20px",
			},
		},
		{
			name:  "four values",
			input: map[string]string{"margin": "10px 20px 30px 40px"},
			check: map[string]string{
				"margin-top": "10px", "margin-right": "20px",
				"margin-bottom": "30px", "margin-left": "40px",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExpandShorthands(tt.input)
			for k, want := range tt.check {
				if got[k] != want {
					t.Errorf("%s: got[%q] = %q, want %q", tt.name, k, got[k], want)
				}
			}
			if got["margin"] != tt.input["margin"] {
				t.Errorf("original margin property not preserved")
			}
		})
	}
}

func TestExpandPadding(t *testing.T) {
	got := ExpandShorthands(map[string]string{"padding": "5px 10px"})
	if got["padding-top"] != "5px" || got["padding-left"] != "10px" {
		t.Errorf("padding expansion failed: %v", got)
	}
}

func TestExpandBorder(t *testing.T) {
	got := ExpandShorthands(map[string]string{"border": "1px solid red"})
	sides := []string{"top", "right", "bottom", "left"}
	for _, side := range sides {
		if got["border-"+side+"-width"] != "1px" {
			t.Errorf("border-%s-width = %q, want %q", side, got["border-"+side+"-width"], "1px")
		}
		if got["border-"+side+"-style"] != "solid" {
			t.Errorf("border-%s-style = %q, want %q", side, got["border-"+side+"-style"], "solid")
		}
		if got["border-"+side+"-color"] != "red" {
			t.Errorf("border-%s-color = %q, want %q", side, got["border-"+side+"-color"], "red")
		}
	}
}

func TestExpandBorderSide(t *testing.T) {
	got := ExpandShorthands(map[string]string{"border-top": "2px dashed blue"})
	if got["border-top-width"] != "2px" {
		t.Errorf("border-top-width = %q, want %q", got["border-top-width"], "2px")
	}
	if got["border-top-style"] != "dashed" {
		t.Errorf("border-top-style = %q, want %q", got["border-top-style"], "dashed")
	}
	if got["border-top-color"] != "blue" {
		t.Errorf("border-top-color = %q, want %q", got["border-top-color"], "blue")
	}
	if _, ok := got["border-bottom-width"]; ok {
		t.Errorf("border-top should not expand to bottom")
	}
}

func TestExpandBorderRadius(t *testing.T) {
	tests := []struct {
		name  string
		input string
		check map[string]string
	}{
		{
			name:  "one value",
			input: "10px",
			check: map[string]string{
				"border-top-left-radius":     "10px",
				"border-top-right-radius":    "10px",
				"border-bottom-right-radius": "10px",
				"border-bottom-left-radius":  "10px",
			},
		},
		{
			name:  "two values",
			input: "10px 20px",
			check: map[string]string{
				"border-top-left-radius":     "10px",
				"border-top-right-radius":    "20px",
				"border-bottom-right-radius": "10px",
				"border-bottom-left-radius":  "20px",
			},
		},
		{
			name:  "four values",
			input: "10px 20px 30px 40px",
			check: map[string]string{
				"border-top-left-radius":     "10px",
				"border-top-right-radius":    "20px",
				"border-bottom-right-radius": "30px",
				"border-bottom-left-radius":  "40px",
			},
		},
		{
			name:  "slash syntax",
			input: "10px / 20px",
			check: map[string]string{
				"border-top-left-radius":     "10px 20px",
				"border-top-right-radius":    "10px 20px",
				"border-bottom-right-radius": "10px 20px",
				"border-bottom-left-radius":  "10px 20px",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExpandShorthands(map[string]string{"border-radius": tt.input})
			for k, want := range tt.check {
				if got[k] != want {
					t.Errorf("%s: got[%q] = %q, want %q", tt.name, k, got[k], want)
				}
			}
		})
	}
}

func TestExpandFlex(t *testing.T) {
	tests := []struct {
		name  string
		input string
		grow  string
		shrk  string
		basis string
	}{
		{"single number", "1", "1", "1", "0%"},
		{"three values", "1 1 100px", "1", "1", "100px"},
		{"none", "none", "0", "0", "auto"},
		{"auto", "auto", "1", "1", "auto"},
		{"two with basis", "2 100px", "2", "1", "100px"},
		{"two numbers", "2 3", "2", "3", "0%"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExpandShorthands(map[string]string{"flex": tt.input})
			if got["flex-grow"] != tt.grow {
				t.Errorf("flex-grow = %q, want %q", got["flex-grow"], tt.grow)
			}
			if got["flex-shrink"] != tt.shrk {
				t.Errorf("flex-shrink = %q, want %q", got["flex-shrink"], tt.shrk)
			}
			if got["flex-basis"] != tt.basis {
				t.Errorf("flex-basis = %q, want %q", got["flex-basis"], tt.basis)
			}
		})
	}
}

func TestExpandGap(t *testing.T) {
	got := ExpandShorthands(map[string]string{"gap": "10px"})
	if got["row-gap"] != "10px" || got["column-gap"] != "10px" {
		t.Errorf("gap single: %v", got)
	}

	got = ExpandShorthands(map[string]string{"gap": "10px 20px"})
	if got["row-gap"] != "10px" || got["column-gap"] != "20px" {
		t.Errorf("gap two: %v", got)
	}
}

func TestExpandBackground(t *testing.T) {
	got := ExpandShorthands(map[string]string{"background": "red"})
	if got["background-color"] != "red" {
		t.Errorf("background color: %v", got)
	}

	got = ExpandShorthands(map[string]string{"background": "linear-gradient(to right, red, blue)"})
	if got["background-image"] != "linear-gradient(to right, red, blue)" {
		t.Errorf("background gradient: %v", got)
	}

	got = ExpandShorthands(map[string]string{"background": "url(image.png)"})
	if got["background-image"] != "url(image.png)" {
		t.Errorf("background url: %v", got)
	}
}

func TestExpandFont(t *testing.T) {
	got := ExpandShorthands(map[string]string{"font": "italic bold 16px/1.5 Arial, sans-serif"})
	if got["font-style"] != "italic" {
		t.Errorf("font-style = %q", got["font-style"])
	}
	if got["font-weight"] != "bold" {
		t.Errorf("font-weight = %q", got["font-weight"])
	}
	if got["font-size"] != "16px" {
		t.Errorf("font-size = %q", got["font-size"])
	}
	if got["line-height"] != "1.5" {
		t.Errorf("line-height = %q", got["line-height"])
	}
	if got["font-family"] != "Arial, sans-serif" {
		t.Errorf("font-family = %q", got["font-family"])
	}
}

func TestExpandFontSimple(t *testing.T) {
	got := ExpandShorthands(map[string]string{"font": "16px Arial"})
	if got["font-size"] != "16px" {
		t.Errorf("font-size = %q", got["font-size"])
	}
	if got["font-family"] != "Arial" {
		t.Errorf("font-family = %q", got["font-family"])
	}
}

func TestExpandTextDecoration(t *testing.T) {
	got := ExpandShorthands(map[string]string{"text-decoration": "underline red wavy"})
	if got["text-decoration-line"] != "underline" {
		t.Errorf("text-decoration-line = %q", got["text-decoration-line"])
	}
	if got["text-decoration-color"] != "red" {
		t.Errorf("text-decoration-color = %q", got["text-decoration-color"])
	}
	if got["text-decoration-style"] != "wavy" {
		t.Errorf("text-decoration-style = %q", got["text-decoration-style"])
	}
}

func TestExpandOverflow(t *testing.T) {
	got := ExpandShorthands(map[string]string{"overflow": "hidden"})
	if got["overflow-x"] != "hidden" || got["overflow-y"] != "hidden" {
		t.Errorf("overflow single: %v", got)
	}

	got = ExpandShorthands(map[string]string{"overflow": "hidden visible"})
	if got["overflow-x"] != "hidden" || got["overflow-y"] != "visible" {
		t.Errorf("overflow two: %v", got)
	}
}

func TestExpandBorderWidth(t *testing.T) {
	got := ExpandShorthands(map[string]string{"border-width": "1px 2px 3px 4px"})
	if got["border-top-width"] != "1px" {
		t.Errorf("border-top-width = %q", got["border-top-width"])
	}
	if got["border-right-width"] != "2px" {
		t.Errorf("border-right-width = %q", got["border-right-width"])
	}
	if got["border-bottom-width"] != "3px" {
		t.Errorf("border-bottom-width = %q", got["border-bottom-width"])
	}
	if got["border-left-width"] != "4px" {
		t.Errorf("border-left-width = %q", got["border-left-width"])
	}
}

func TestExpandBorderStyle(t *testing.T) {
	got := ExpandShorthands(map[string]string{"border-style": "solid dashed"})
	if got["border-top-style"] != "solid" || got["border-right-style"] != "dashed" {
		t.Errorf("border-style: %v", got)
	}
}

func TestExpandBorderColor(t *testing.T) {
	got := ExpandShorthands(map[string]string{"border-color": "red blue green yellow"})
	if got["border-top-color"] != "red" || got["border-right-color"] != "blue" ||
		got["border-bottom-color"] != "green" || got["border-left-color"] != "yellow" {
		t.Errorf("border-color: %v", got)
	}
}

func TestOriginalPropertiesPreserved(t *testing.T) {
	input := map[string]string{
		"margin":  "10px",
		"color":   "red",
		"display": "flex",
	}
	got := ExpandShorthands(input)
	if got["margin"] != "10px" {
		t.Error("original margin not preserved")
	}
	if got["color"] != "red" {
		t.Error("original color not preserved")
	}
	if got["display"] != "flex" {
		t.Error("original display not preserved")
	}
}

func TestExpandBorderWithRGB(t *testing.T) {
	got := ExpandShorthands(map[string]string{"border": "1px solid rgb(255, 0, 0)"})
	if got["border-top-width"] != "1px" {
		t.Errorf("border-top-width = %q", got["border-top-width"])
	}
	if got["border-top-style"] != "solid" {
		t.Errorf("border-top-style = %q", got["border-top-style"])
	}
	if got["border-top-color"] != "rgb(255, 0, 0)" {
		t.Errorf("border-top-color = %q", got["border-top-color"])
	}
}
