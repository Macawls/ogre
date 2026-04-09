package style

import (
	"testing"
)

func TestTailwindFontSizesWithLineHeight(t *testing.T) {
	tests := []struct {
		class      string
		fontSize   string
		lineHeight string
	}{
		{"text-xs", "12px", "16px"},
		{"text-sm", "14px", "20px"},
		{"text-base", "16px", "24px"},
		{"text-lg", "18px", "28px"},
		{"text-xl", "20px", "28px"},
		{"text-2xl", "24px", "32px"},
		{"text-3xl", "30px", "36px"},
		{"text-4xl", "36px", "40px"},
		{"text-5xl", "48px", "48px"},
		{"text-6xl", "60px", "60px"},
		{"text-7xl", "72px", "72px"},
		{"text-8xl", "96px", "96px"},
		{"text-9xl", "128px", "128px"},
	}
	for _, tt := range tests {
		t.Run(tt.class, func(t *testing.T) {
			result := ResolveTailwind([]string{tt.class})
			if result["font-size"] != tt.fontSize {
				t.Errorf("font-size: got %q, want %q", result["font-size"], tt.fontSize)
			}
			if result["line-height"] != tt.lineHeight {
				t.Errorf("line-height: got %q, want %q", result["line-height"], tt.lineHeight)
			}
		})
	}
}

func TestTailwindFractionalSpacing(t *testing.T) {
	tests := []struct {
		class string
		prop  string
		value string
	}{
		{"p-0.5", "padding-top", "2px"},
		{"p-1.5", "padding-top", "6px"},
		{"p-2.5", "padding-top", "10px"},
		{"p-3.5", "padding-top", "14px"},
		{"m-0.5", "margin-top", "2px"},
		{"m-1.5", "margin-top", "6px"},
		{"px-0.5", "padding-left", "2px"},
		{"py-1.5", "padding-top", "6px"},
	}
	for _, tt := range tests {
		t.Run(tt.class, func(t *testing.T) {
			result := ResolveTailwind([]string{tt.class})
			if result[tt.prop] != tt.value {
				t.Errorf("%s: got %q, want %q", tt.prop, result[tt.prop], tt.value)
			}
		})
	}
}

func TestTailwindColors(t *testing.T) {
	tests := []struct {
		class string
		prop  string
		value string
	}{
		{"text-red-500", "color", "#ef4444"},
		{"bg-blue-500", "background-color", "#3b82f6"},
		{"border-green-500", "border-color", "#22c55e"},
		{"text-white", "color", "#ffffff"},
		{"text-black", "color", "#000000"},
		{"bg-transparent", "background-color", "transparent"},
		{"bg-slate-950", "background-color", "#020617"},
	}
	for _, tt := range tests {
		t.Run(tt.class, func(t *testing.T) {
			result := ResolveTailwind([]string{tt.class})
			if result[tt.prop] != tt.value {
				t.Errorf("%s: got %q, want %q", tt.prop, result[tt.prop], tt.value)
			}
		})
	}
}

func TestTailwindFlexUtilities(t *testing.T) {
	result := ResolveTailwind([]string{"flex"})
	if result["display"] != "flex" {
		t.Errorf("flex display: got %q, want %q", result["display"], "flex")
	}

	result = ResolveTailwind([]string{"flex-col"})
	if result["flex-direction"] != "column" {
		t.Errorf("flex-col: got %q, want %q", result["flex-direction"], "column")
	}

	result = ResolveTailwind([]string{"items-center"})
	if result["align-items"] != "center" {
		t.Errorf("items-center: got %q, want %q", result["align-items"], "center")
	}

	result = ResolveTailwind([]string{"justify-between"})
	if result["justify-content"] != "space-between" {
		t.Errorf("justify-between: got %q, want %q", result["justify-content"], "space-between")
	}
}

func TestTailwindDisplay(t *testing.T) {
	tests := []struct {
		class string
		value string
	}{
		{"hidden", "none"},
		{"block", "block"},
		{"inline", "flex"},
		{"inline-flex", "flex"},
		{"grid", "flex"},
	}
	for _, tt := range tests {
		t.Run(tt.class, func(t *testing.T) {
			result := ResolveTailwind([]string{tt.class})
			if result["display"] != tt.value {
				t.Errorf("display: got %q, want %q", result["display"], tt.value)
			}
		})
	}
}

func TestTailwindArbitrary(t *testing.T) {
	tests := []struct {
		class string
		prop  string
		value string
	}{
		{"text-[32px]", "font-size", "32px"},
		{"bg-[#ff0000]", "background-color", "#ff0000"},
		{"p-[20px]", "padding-top", "20px"},
		{"w-[200px]", "width", "200px"},
		{"h-[100px]", "height", "100px"},
		{"rounded-[8px]", "border-radius", "8px"},
		{"gap-[12px]", "gap", "12px"},
	}
	for _, tt := range tests {
		t.Run(tt.class, func(t *testing.T) {
			result := ResolveTailwind([]string{tt.class})
			if result[tt.prop] != tt.value {
				t.Errorf("%s: got %q, want %q", tt.prop, result[tt.prop], tt.value)
			}
		})
	}
}

func TestTailwindMultipleClasses(t *testing.T) {
	result := ResolveTailwind([]string{"flex", "flex-col", "items-center", "p-4", "bg-blue-500", "text-white"})
	expected := map[string]string{
		"display":          "flex",
		"flex-direction":   "column",
		"align-items":      "center",
		"padding-top":      "16px",
		"padding-right":    "16px",
		"padding-bottom":   "16px",
		"padding-left":     "16px",
		"background-color": "#3b82f6",
		"color":            "#ffffff",
	}
	for prop, want := range expected {
		if result[prop] != want {
			t.Errorf("%s: got %q, want %q", prop, result[prop], want)
		}
	}
}

func TestTailwindInlineStyleWins(t *testing.T) {
	tw := ResolveTailwind([]string{"text-red-500", "p-4"})
	inline := map[string]string{
		"color":       "#00ff00",
		"padding-top": "8px",
	}
	for k, v := range inline {
		tw[k] = v
	}
	if tw["color"] != "#00ff00" {
		t.Errorf("inline color should win: got %q", tw["color"])
	}
	if tw["padding-top"] != "8px" {
		t.Errorf("inline padding should win: got %q", tw["padding-top"])
	}
	if tw["padding-right"] != "16px" {
		t.Errorf("tw padding-right should remain: got %q", tw["padding-right"])
	}
}

func TestTailwindPosition(t *testing.T) {
	tests := []struct {
		class string
		value string
	}{
		{"static", "static"},
		{"relative", "relative"},
		{"absolute", "absolute"},
		{"fixed", "absolute"},
		{"sticky", "relative"},
	}
	for _, tt := range tests {
		t.Run(tt.class, func(t *testing.T) {
			result := ResolveTailwind([]string{tt.class})
			if result["position"] != tt.value {
				t.Errorf("position: got %q, want %q", result["position"], tt.value)
			}
		})
	}
}

func TestTailwindZIndex(t *testing.T) {
	tests := []struct {
		class string
		value string
	}{
		{"z-0", "0"},
		{"z-10", "10"},
		{"z-20", "20"},
		{"z-30", "30"},
		{"z-40", "40"},
		{"z-50", "50"},
		{"z-auto", "auto"},
	}
	for _, tt := range tests {
		t.Run(tt.class, func(t *testing.T) {
			result := ResolveTailwind([]string{tt.class})
			if result["z-index"] != tt.value {
				t.Errorf("z-index: got %q, want %q", result["z-index"], tt.value)
			}
		})
	}
}

func TestTailwindAspectRatio(t *testing.T) {
	tests := []struct {
		class string
		value string
	}{
		{"aspect-square", "1"},
		{"aspect-video", "1.7778"},
		{"aspect-auto", "auto"},
	}
	for _, tt := range tests {
		t.Run(tt.class, func(t *testing.T) {
			result := ResolveTailwind([]string{tt.class})
			if result["aspect-ratio"] != tt.value {
				t.Errorf("aspect-ratio: got %q, want %q", result["aspect-ratio"], tt.value)
			}
		})
	}
}

func TestTailwindSpaceUtilities(t *testing.T) {
	result := ResolveTailwind([]string{"space-x-4"})
	if result["column-gap"] != "16px" {
		t.Errorf("space-x-4 column-gap: got %q, want %q", result["column-gap"], "16px")
	}

	result = ResolveTailwind([]string{"space-y-2"})
	if result["row-gap"] != "8px" {
		t.Errorf("space-y-2 row-gap: got %q, want %q", result["row-gap"], "8px")
	}
}

func TestTailwindLineClamp(t *testing.T) {
	for n := 1; n <= 6; n++ {
		cls := "line-clamp-" + string(rune('0'+n))
		result := ResolveTailwind([]string{cls})
		want := string(rune('0' + n))
		if result["-webkit-line-clamp"] != want {
			t.Errorf("%s: got %q, want %q", cls, result["-webkit-line-clamp"], want)
		}
		if result["overflow"] != "hidden" {
			t.Errorf("%s overflow: got %q, want %q", cls, result["overflow"], "hidden")
		}
	}
}

func TestTailwindSizeShortcut(t *testing.T) {
	result := ResolveTailwind([]string{"size-4"})
	if result["width"] != "16px" {
		t.Errorf("size-4 width: got %q, want %q", result["width"], "16px")
	}
	if result["height"] != "16px" {
		t.Errorf("size-4 height: got %q, want %q", result["height"], "16px")
	}
}

func TestTailwindFitMinMax(t *testing.T) {
	tests := []struct {
		class string
		prop  string
	}{
		{"w-fit", "width"},
		{"h-fit", "height"},
		{"w-min", "width"},
		{"w-max", "width"},
	}
	for _, tt := range tests {
		t.Run(tt.class, func(t *testing.T) {
			result := ResolveTailwind([]string{tt.class})
			if result[tt.prop] != "auto" {
				t.Errorf("%s: got %q, want %q", tt.prop, result[tt.prop], "auto")
			}
		})
	}
}
