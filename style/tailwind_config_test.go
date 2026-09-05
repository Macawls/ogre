package style

import (
	"strings"
	"testing"
)

func TestConfigColorOverride(t *testing.T) {
	cfg := &TailwindConfig{
		Colors: map[string]map[int]string{
			"brand": {500: "#ff00aa"},
		},
	}
	out := ResolveTailwindWithConfig([]string{"bg-brand-500"}, TailwindV3, cfg)
	if got := out["background-color"]; got != "#ff00aa" {
		t.Fatalf("bg-brand-500 = %q, want #ff00aa", got)
	}
}

func TestConfigColorOverrideBeatsDefaultShade(t *testing.T) {
	cfg := &TailwindConfig{
		Colors: map[string]map[int]string{
			"slate": {500: "#abcdef"},
		},
	}
	out := ResolveTailwindWithConfig([]string{"bg-slate-500"}, TailwindV3, cfg)
	if got := out["background-color"]; got != "#abcdef" {
		t.Fatalf("overridden slate-500 = %q, want #abcdef", got)
	}
	// Non-overridden shade should still fall back to default.
	out2 := ResolveTailwindWithConfig([]string{"bg-slate-400"}, TailwindV3, cfg)
	if got := out2["background-color"]; got != "#94a3b8" {
		t.Fatalf("fallback slate-400 = %q, want #94a3b8", got)
	}
}

func TestConfigSpacingOverride(t *testing.T) {
	cfg := &TailwindConfig{
		Spacing: map[string]string{
			"13":      "3.25rem",
			"sidebar": "220px",
		},
	}
	out := ResolveTailwindWithConfig([]string{"p-13", "w-sidebar"}, TailwindV3, cfg)
	if got := out["padding-top"]; got != "3.25rem" {
		t.Fatalf("p-13 → padding-top = %q, want 3.25rem", got)
	}
	if got := out["width"]; got != "220px" {
		t.Fatalf("w-sidebar → width = %q, want 220px", got)
	}
}

func TestConfigSpacingDoesNotBreakBuiltins(t *testing.T) {
	cfg := &TailwindConfig{Spacing: map[string]string{"sidebar": "220px"}}
	out := ResolveTailwindWithConfig([]string{"p-4"}, TailwindV3, cfg)
	if got := out["padding-top"]; got != "16px" {
		t.Fatalf("p-4 with cfg = %q, want 16px", got)
	}
}

func TestConfigFontOverride(t *testing.T) {
	cfg := &TailwindConfig{
		Fonts: map[string]string{"display": "\"Inter\", sans-serif"},
	}
	out := ResolveTailwindWithConfig([]string{"font-display"}, TailwindV3, cfg)
	if got := out["font-family"]; got != "\"Inter\", sans-serif" {
		t.Fatalf("font-display = %q", got)
	}
}

func TestConfigFontDoesNotShadowWeight(t *testing.T) {
	cfg := &TailwindConfig{Fonts: map[string]string{"bold": "should-not-apply"}}
	out := ResolveTailwindWithConfig([]string{"font-bold"}, TailwindV3, cfg)
	if got := out["font-weight"]; got != "700" {
		t.Fatalf("font-bold weight = %q, want 700 (static match wins)", got)
	}
	if _, present := out["font-family"]; present {
		t.Fatalf("font-bold should not also emit font-family from cfg")
	}
}

func TestConfigExtendRawClass(t *testing.T) {
	cfg := &TailwindConfig{
		Extend: map[string]map[string]string{
			"card": {
				"background-color": "#123456",
				"border-radius":    "8px",
				"padding":          "16px",
			},
		},
	}
	out := ResolveTailwindWithConfig([]string{"card"}, TailwindV3, cfg)
	if got := out["background-color"]; got != "#123456" {
		t.Fatalf("card bg = %q", got)
	}
	if got := out["border-radius"]; got != "8px" {
		t.Fatalf("card radius = %q", got)
	}
}

func TestConfigExtendBeatsStaticClass(t *testing.T) {
	cfg := &TailwindConfig{
		Extend: map[string]map[string]string{
			"flex": {"display": "grid"},
		},
	}
	out := ResolveTailwindWithConfig([]string{"flex"}, TailwindV3, cfg)
	if got := out["display"]; got != "grid" {
		t.Fatalf("extend should override static 'flex' → %q", got)
	}
}

func TestConfigNilFallsThrough(t *testing.T) {
	out := ResolveTailwindWithConfig([]string{"bg-slate-500"}, TailwindV3, nil)
	if got := out["background-color"]; got != "#64748b" {
		t.Fatalf("nil cfg default = %q, want #64748b", got)
	}
}

func TestConfigKeyStable(t *testing.T) {
	a := &TailwindConfig{
		Colors:  map[string]map[int]string{"a": {1: "x"}, "b": {2: "y"}},
		Spacing: map[string]string{"k1": "v1", "k2": "v2"},
	}
	b := &TailwindConfig{
		Colors:  map[string]map[int]string{"b": {2: "y"}, "a": {1: "x"}},
		Spacing: map[string]string{"k2": "v2", "k1": "v1"},
	}
	if a.Key() != b.Key() {
		t.Fatalf("logically equal configs have different keys:\n a=%s\n b=%s", a.Key(), b.Key())
	}
}

func TestConfigKeyDiffersOnChange(t *testing.T) {
	a := &TailwindConfig{Spacing: map[string]string{"k": "10px"}}
	b := &TailwindConfig{Spacing: map[string]string{"k": "12px"}}
	if a.Key() == b.Key() {
		t.Fatalf("different configs must have different keys")
	}
}

func TestConfigCacheIsolation(t *testing.T) {
	cfgA := &TailwindConfig{Colors: map[string]map[int]string{"slate": {500: "#aaaaaa"}}}
	cfgB := &TailwindConfig{Colors: map[string]map[int]string{"slate": {500: "#bbbbbb"}}}
	a1 := ResolveTailwindWithConfig([]string{"bg-slate-500"}, TailwindV3, cfgA)
	b1 := ResolveTailwindWithConfig([]string{"bg-slate-500"}, TailwindV3, cfgB)
	if a1["background-color"] != "#aaaaaa" || b1["background-color"] != "#bbbbbb" {
		t.Fatalf("configs cross-contaminated: a=%q b=%q", a1["background-color"], b1["background-color"])
	}
	// Round-trip: default (nil) cache must remain untouched.
	def := ResolveTailwindWithConfig([]string{"bg-slate-500"}, TailwindV3, nil)
	if def["background-color"] != "#64748b" {
		t.Fatalf("nil-cfg cache polluted: %q", def["background-color"])
	}
}

// --- CSS @theme parser ---

func TestParseThemeCSSBasic(t *testing.T) {
	css := `
@theme {
  --color-brand-500: #ff00aa;
  --color-brand-600: #ee0099;
  --spacing-sidebar: 220px;
  --spacing-13: 3.25rem;
  --font-display: "Inter", sans-serif;
  --breakpoint-3xl: 1920px;
}
`
	cfg, err := ParseTailwindThemeCSS(css)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if got := cfg.Colors["brand"][500]; got != "#ff00aa" {
		t.Errorf("brand 500 = %q", got)
	}
	if got := cfg.Colors["brand"][600]; got != "#ee0099" {
		t.Errorf("brand 600 = %q", got)
	}
	if got := cfg.Spacing["sidebar"]; got != "220px" {
		t.Errorf("spacing sidebar = %q", got)
	}
	if got := cfg.Spacing["13"]; got != "3.25rem" {
		t.Errorf("spacing 13 = %q", got)
	}
	if got := cfg.Fonts["display"]; got != "\"Inter\", sans-serif" {
		t.Errorf("font display = %q", got)
	}
	if got := cfg.Screens["3xl"]; got != "1920px" {
		t.Errorf("breakpoint 3xl = %q", got)
	}
}

func TestParseThemeCSSWithComments(t *testing.T) {
	css := `
/* leading comment */
@theme {
  /* inline comment */
  --color-brand-500: #ff00aa; /* trailing */
  --spacing-13: /* mid */ 3.25rem;
}
`
	cfg, err := ParseTailwindThemeCSS(css)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if cfg.Colors["brand"][500] != "#ff00aa" {
		t.Errorf("colors: %v", cfg.Colors)
	}
	if cfg.Spacing["13"] != "3.25rem" {
		t.Errorf("spacing: %v", cfg.Spacing)
	}
}

func TestParseThemeCSSVariants(t *testing.T) {
	css := `@theme inline { --color-a-100: red; }`
	cfg, err := ParseTailwindThemeCSS(css)
	if err != nil {
		t.Fatalf("parse @theme inline: %v", err)
	}
	if cfg.Colors["a"][100] != "red" {
		t.Errorf("got %v", cfg.Colors)
	}
}

func TestParseThemeCSSMultipleBlocks(t *testing.T) {
	css := `
@theme { --color-a-100: red; }
@theme { --color-a-200: blue; --spacing-1: 4px; }
`
	cfg, err := ParseTailwindThemeCSS(css)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if cfg.Colors["a"][100] != "red" || cfg.Colors["a"][200] != "blue" {
		t.Errorf("merge failed: %v", cfg.Colors)
	}
	if cfg.Spacing["1"] != "4px" {
		t.Errorf("spacing merge failed: %v", cfg.Spacing)
	}
}

func TestParseThemeCSSCompoundColorName(t *testing.T) {
	css := `@theme { --color-brand-primary-500: #123456; }`
	cfg, err := ParseTailwindThemeCSS(css)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if got := cfg.Colors["brand-primary"][500]; got != "#123456" {
		t.Fatalf("compound-name color = %q", got)
	}
}

func TestParseThemeCSSMalformed(t *testing.T) {
	if _, err := ParseTailwindThemeCSS(`@theme {`); err == nil {
		t.Errorf("unclosed @theme should error")
	}
	if _, err := ParseTailwindThemeCSS(`body { color: red; }`); err == nil {
		t.Errorf("no @theme block should error")
	}
}

func TestParseThemeCSSSkipsNonNamespacedVars(t *testing.T) {
	css := `@theme { --custom: whatever; --color-a-100: red; }`
	cfg, err := ParseTailwindThemeCSS(css)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if strings.Contains(cfg.Fonts["custom"]+cfg.Spacing["custom"], "whatever") {
		t.Errorf("unnamespaced var leaked: %+v", cfg)
	}
	if cfg.Colors["a"][100] != "red" {
		t.Errorf("expected --color-a-100 to still parse: %v", cfg.Colors)
	}
}

func TestParseThemeCSSIntegration(t *testing.T) {
	css := `@theme { --color-brand-500: #112233; --spacing-13: 3.25rem; }`
	cfg, err := ParseTailwindThemeCSS(css)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	out := ResolveTailwindWithConfig([]string{"bg-brand-500", "p-13"}, TailwindV3, cfg)
	if out["background-color"] != "#112233" {
		t.Errorf("bg color = %q", out["background-color"])
	}
	if out["padding-top"] != "3.25rem" {
		t.Errorf("padding = %q", out["padding-top"])
	}
}
