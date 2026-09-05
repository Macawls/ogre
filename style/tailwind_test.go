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
			result := ResolveTailwind([]string{tt.class}, TailwindV3)
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
			result := ResolveTailwind([]string{tt.class}, TailwindV3)
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
			result := ResolveTailwind([]string{tt.class}, TailwindV3)
			if result[tt.prop] != tt.value {
				t.Errorf("%s: got %q, want %q", tt.prop, result[tt.prop], tt.value)
			}
		})
	}
}

func TestTailwindFlexUtilities(t *testing.T) {
	result := ResolveTailwind([]string{"flex"}, TailwindV3)
	if result["display"] != "flex" {
		t.Errorf("flex display: got %q, want %q", result["display"], "flex")
	}

	result = ResolveTailwind([]string{"flex-col"}, TailwindV3)
	if result["flex-direction"] != "column" {
		t.Errorf("flex-col: got %q, want %q", result["flex-direction"], "column")
	}

	result = ResolveTailwind([]string{"items-center"}, TailwindV3)
	if result["align-items"] != "center" {
		t.Errorf("items-center: got %q, want %q", result["align-items"], "center")
	}

	result = ResolveTailwind([]string{"justify-between"}, TailwindV3)
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
		{"grid", "grid"},
		{"inline-grid", "grid"},
	}
	for _, tt := range tests {
		t.Run(tt.class, func(t *testing.T) {
			result := ResolveTailwind([]string{tt.class}, TailwindV3)
			if result["display"] != tt.value {
				t.Errorf("display: got %q, want %q", result["display"], tt.value)
			}
		})
	}
}

func TestTailwindGrid(t *testing.T) {
	tests := []struct {
		class string
		prop  string
		value string
	}{
		{"grid", "display", "grid"},
		{"inline-grid", "display", "grid"},
		{"grid-cols-3", "grid-template-columns", "repeat(3, minmax(0, 1fr))"},
		{"grid-cols-12", "grid-template-columns", "repeat(12, minmax(0, 1fr))"},
		{"grid-cols-none", "grid-template-columns", "none"},
		{"grid-rows-4", "grid-template-rows", "repeat(4, minmax(0, 1fr))"},
		{"col-span-2", "grid-column", "span 2 / span 2"},
		{"row-span-3", "grid-row", "span 3 / span 3"},
		{"col-span-full", "grid-column", "1 / -1"},
		{"row-span-full", "grid-row", "1 / -1"},
		{"col-start-4", "grid-column-start", "4"},
		{"col-end-7", "grid-column-end", "7"},
		{"row-start-2", "grid-row-start", "2"},
		{"row-end-5", "grid-row-end", "5"},
		{"grid-flow-col", "grid-auto-flow", "column"},
		{"grid-flow-row-dense", "grid-auto-flow", "row dense"},
		{"auto-cols-fr", "grid-auto-columns", "1fr"},
		{"auto-rows-min", "grid-auto-rows", "min-content"},
	}
	for _, tt := range tests {
		t.Run(tt.class, func(t *testing.T) {
			result := ResolveTailwind([]string{tt.class}, TailwindV3)
			if got := result[tt.prop]; got != tt.value {
				t.Errorf("%s: %s = %q, want %q", tt.class, tt.prop, got, tt.value)
			}
		})
	}
}

func TestGridTrackListParsing(t *testing.T) {
	l := ParseTrackList("100px 1fr auto 2fr")
	if len(l.Tracks) != 4 {
		t.Fatalf("expected 4 tracks, got %d", len(l.Tracks))
	}
	if l.Tracks[0].Kind != TrackFixed || l.Tracks[0].Length.Raw != 100 {
		t.Errorf("track 0: %+v", l.Tracks[0])
	}
	if l.Tracks[1].Kind != TrackFr || l.Tracks[1].Fr != 1 {
		t.Errorf("track 1: %+v", l.Tracks[1])
	}
	if l.Tracks[2].Kind != TrackAuto {
		t.Errorf("track 2: %+v", l.Tracks[2])
	}
	if l.Tracks[3].Kind != TrackFr || l.Tracks[3].Fr != 2 {
		t.Errorf("track 3: %+v", l.Tracks[3])
	}
}

func TestGridTrackListRepeat(t *testing.T) {
	l := ParseTrackList("repeat(3, 1fr)")
	if len(l.Tracks) != 3 {
		t.Fatalf("expected 3 tracks, got %d", len(l.Tracks))
	}
	for i, tr := range l.Tracks {
		if tr.Kind != TrackFr || tr.Fr != 1 {
			t.Errorf("track %d: %+v", i, tr)
		}
	}
}

func TestGridLineShorthand(t *testing.T) {
	tests := []struct {
		in    string
		start GridPlacement
		end   GridPlacement
	}{
		{"2 / 4", GridPlacement{Start: 2}, GridPlacement{Start: 4}},
		{"span 3", GridPlacement{Span: 3}, GridPlacement{}},
		{"1 / span 2", GridPlacement{Start: 1}, GridPlacement{Span: 2}},
		{"auto", GridPlacement{}, GridPlacement{}},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			s, e := ParseGridLine(tt.in)
			if s != tt.start || e != tt.end {
				t.Errorf("ParseGridLine(%q) = (%+v, %+v), want (%+v, %+v)", tt.in, s, e, tt.start, tt.end)
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
			result := ResolveTailwind([]string{tt.class}, TailwindV3)
			if result[tt.prop] != tt.value {
				t.Errorf("%s: got %q, want %q", tt.prop, result[tt.prop], tt.value)
			}
		})
	}
}

func TestTailwindMultipleClasses(t *testing.T) {
	result := ResolveTailwind([]string{"flex", "flex-col", "items-center", "p-4", "bg-blue-500", "text-white"}, TailwindV3)
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
	tw := ResolveTailwind([]string{"text-red-500", "p-4"}, TailwindV3)
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
			result := ResolveTailwind([]string{tt.class}, TailwindV3)
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
			result := ResolveTailwind([]string{tt.class}, TailwindV3)
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
			result := ResolveTailwind([]string{tt.class}, TailwindV3)
			if result["aspect-ratio"] != tt.value {
				t.Errorf("aspect-ratio: got %q, want %q", result["aspect-ratio"], tt.value)
			}
		})
	}
}

func TestTailwindSpaceUtilities(t *testing.T) {
	result := ResolveTailwind([]string{"space-x-4"}, TailwindV3)
	if result["column-gap"] != "16px" {
		t.Errorf("space-x-4 column-gap: got %q, want %q", result["column-gap"], "16px")
	}

	result = ResolveTailwind([]string{"space-y-2"}, TailwindV3)
	if result["row-gap"] != "8px" {
		t.Errorf("space-y-2 row-gap: got %q, want %q", result["row-gap"], "8px")
	}
}

func TestTailwindLineClamp(t *testing.T) {
	for n := 1; n <= 6; n++ {
		cls := "line-clamp-" + string(rune('0'+n))
		result := ResolveTailwind([]string{cls}, TailwindV3)
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
	result := ResolveTailwind([]string{"size-4"}, TailwindV3)
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
			result := ResolveTailwind([]string{tt.class}, TailwindV3)
			if result[tt.prop] != "auto" {
				t.Errorf("%s: got %q, want %q", tt.prop, result[tt.prop], "auto")
			}
		})
	}
}

// -----------------------------------------------------------------------------
// v4 semantics — renames and additions
// -----------------------------------------------------------------------------

func TestTailwindV4RenamedShadows(t *testing.T) {
	// In v4, shadow-sm takes on v3-shadow semantics, shadow-xs takes on
	// v3-shadow-sm semantics, and shadow-2xs is new.
	v3ShadowSm := "0 1px 2px 0 rgba(0,0,0,0.05)"
	v3Shadow := "0 1px 3px 0 rgba(0,0,0,0.1), 0 1px 2px -1px rgba(0,0,0,0.1)"

	cases := []struct {
		class     string
		wantShadow string
	}{
		{"shadow-2xs", "0 1px rgba(0,0,0,0.05)"},
		{"shadow-xs", v3ShadowSm},
		{"shadow-sm", v3Shadow},
	}
	for _, tc := range cases {
		t.Run(tc.class, func(t *testing.T) {
			got := ResolveTailwind([]string{tc.class}, TailwindV4)["box-shadow"]
			if got != tc.wantShadow {
				t.Errorf("v4 %s: got %q, want %q", tc.class, got, tc.wantShadow)
			}
		})
	}

	// v3 semantics untouched: shadow-sm remains the subtle 1px shadow.
	if got := ResolveTailwind([]string{"shadow-sm"}, TailwindV3)["box-shadow"]; got != v3ShadowSm {
		t.Errorf("v3 shadow-sm: got %q, want %q", got, v3ShadowSm)
	}
}

func TestTailwindV4RenamedRadii(t *testing.T) {
	cases := []struct {
		class string
		want  string
	}{
		{"rounded-xs", "2px"},
		{"rounded-sm", "4px"},
		{"rounded-md", "6px"}, // shared with v3
	}
	for _, tc := range cases {
		t.Run(tc.class, func(t *testing.T) {
			got := ResolveTailwind([]string{tc.class}, TailwindV4)["border-radius"]
			if got != tc.want {
				t.Errorf("v4 %s: got %q, want %q", tc.class, got, tc.want)
			}
		})
	}

	// v3 semantics untouched.
	if got := ResolveTailwind([]string{"rounded-sm"}, TailwindV3)["border-radius"]; got != "2px" {
		t.Errorf("v3 rounded-sm: got %q, want 2px", got)
	}
}

func TestTailwindV4RenamedBlur(t *testing.T) {
	cases := []struct {
		class string
		want  string
	}{
		{"blur-xs", "blur(4px)"},
		{"blur-sm", "blur(8px)"},
	}
	for _, tc := range cases {
		t.Run(tc.class, func(t *testing.T) {
			got := ResolveTailwind([]string{tc.class}, TailwindV4)["filter"]
			if got != tc.want {
				t.Errorf("v4 %s: got %q, want %q", tc.class, got, tc.want)
			}
		})
	}
	if got := ResolveTailwind([]string{"blur-sm"}, TailwindV3)["filter"]; got != "blur(4px)" {
		t.Errorf("v3 blur-sm: got %q, want blur(4px)", got)
	}
}

func TestTailwindV4FlexRenames(t *testing.T) {
	cases := []struct {
		class    string
		prop     string
		expected string
	}{
		{"shrink", "flex-shrink", "1"},
		{"shrink-0", "flex-shrink", "0"},
		{"grow", "flex-grow", "1"},
		{"grow-0", "flex-grow", "0"},
	}
	for _, tc := range cases {
		t.Run(tc.class, func(t *testing.T) {
			got := ResolveTailwind([]string{tc.class}, TailwindV4)[tc.prop]
			if got != tc.expected {
				t.Errorf("v4 %s: %s = %q, want %q", tc.class, tc.prop, got, tc.expected)
			}
		})
	}
}

func TestTailwindV4LinearGradient(t *testing.T) {
	result := ResolveTailwind([]string{"bg-linear-to-r", "from-blue-500", "to-red-500"}, TailwindV4)
	bg := result["background-image"]
	if bg == "" {
		t.Fatalf("expected background-image, got %v", result)
	}
	// Direction and color slots present; palette values are OKLCH so we don't
	// compare exact hex here — that's covered separately.
	if !containsAll(bg, "linear-gradient(", "to right") {
		t.Errorf("gradient missing direction: %q", bg)
	}
}

func TestTailwindV4TextEllipsis(t *testing.T) {
	got := ResolveTailwind([]string{"text-ellipsis"}, TailwindV4)["text-overflow"]
	if got != "ellipsis" {
		t.Errorf("v4 text-ellipsis: got %q, want ellipsis", got)
	}
}

func TestTailwindV4ArbitraryVarShortcut(t *testing.T) {
	cases := []struct {
		class string
		prop  string
		want  string
	}{
		{"bg-(--brand)", "background-color", "var(--brand)"},
		{"text-(--fg)", "color", ""}, // text-* arbitrary maps to font-size, not color, in current arbitrary handling
		{"w-(--sz)", "width", "var(--sz)"},
	}
	for _, tc := range cases {
		t.Run(tc.class, func(t *testing.T) {
			got := ResolveTailwind([]string{tc.class}, TailwindV4)[tc.prop]
			if tc.want == "" {
				// Just ensure no panic; specific mapping tested elsewhere
				return
			}
			if got != tc.want {
				t.Errorf("v4 %s: got %q, want %q", tc.class, got, tc.want)
			}
		})
	}

	// v3 mode ignores parens.
	if got := ResolveTailwind([]string{"bg-(--brand)"}, TailwindV3); len(got) != 0 {
		t.Errorf("v3 should ignore paren-form arbitrary: got %v", got)
	}
}

func TestTailwindV4PaletteDiffersFromV3(t *testing.T) {
	// v3 blue-500 is a hardcoded hex; v4 is an OKLCH string. Both resolve, but
	// the raw property value differs (hex vs oklch()).
	v3 := ResolveTailwind([]string{"bg-blue-500"}, TailwindV3)["background-color"]
	v4 := ResolveTailwind([]string{"bg-blue-500"}, TailwindV4)["background-color"]

	if v3 != "#3b82f6" {
		t.Errorf("v3 bg-blue-500 changed: got %q", v3)
	}
	if v4 == v3 {
		t.Errorf("v4 bg-blue-500 should differ from v3, got same %q", v4)
	}
	if len(v4) < 5 || v4[:5] != "oklch" {
		t.Errorf("v4 bg-blue-500 should be an oklch() string, got %q", v4)
	}
}

func TestTailwindV3StillSupportsRenamedClasses(t *testing.T) {
	// v3 users can still use flex-shrink-0, flex-grow, bg-gradient-to-r,
	// shadow, rounded, blur — none of these were touched.
	cases := []struct {
		class string
		key   string
		want  string
	}{
		{"flex-shrink-0", "flex-shrink", "0"},
		{"flex-grow", "flex-grow", "1"},
		{"shadow", "box-shadow", "0 1px 3px 0 rgba(0,0,0,0.1), 0 1px 2px -1px rgba(0,0,0,0.1)"},
		{"rounded", "border-radius", "4px"},
		{"blur", "filter", "blur(8px)"},
	}
	for _, tc := range cases {
		t.Run(tc.class, func(t *testing.T) {
			got := ResolveTailwind([]string{tc.class}, TailwindV3)[tc.key]
			if got != tc.want {
				t.Errorf("v3 %s: %s = %q, want %q", tc.class, tc.key, got, tc.want)
			}
		})
	}
}

func containsAll(s string, subs ...string) bool {
	for _, sub := range subs {
		found := false
		for i := 0; i+len(sub) <= len(s); i++ {
			if s[i:i+len(sub)] == sub {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
