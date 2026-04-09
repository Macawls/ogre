package parse

import "testing"

func TestParseStyleSimple(t *testing.T) {
	m := ParseStyle("color: red; font-size: 16px")
	assertProp(t, m, "color", "red")
	assertProp(t, m, "font-size", "16px")
}

func TestParseStyleURL(t *testing.T) {
	m := ParseStyle("background: url(https://example.com/image.png)")
	assertProp(t, m, "background", "url(https://example.com/image.png)")
}

func TestParseStyleParenthesized(t *testing.T) {
	m := ParseStyle("transform: rotate(45deg); background: linear-gradient(to right, red, blue)")
	assertProp(t, m, "transform", "rotate(45deg)")
	assertProp(t, m, "background", "linear-gradient(to right, red, blue)")
}

func TestParseStyleQuotedValues(t *testing.T) {
	m := ParseStyle("font-family: 'Times New Roman', serif")
	assertProp(t, m, "font-family", "'Times New Roman', serif")
}

func TestParseStyleTrailingSemicolon(t *testing.T) {
	m := ParseStyle("color: red;")
	assertProp(t, m, "color", "red")
	if len(m) != 1 {
		t.Errorf("expected 1 property, got %d", len(m))
	}
}

func TestParseStyleExtraWhitespace(t *testing.T) {
	m := ParseStyle("  color :  red ;  font-size :  16px  ")
	assertProp(t, m, "color", "red")
	assertProp(t, m, "font-size", "16px")
}

func TestParseStyleEmpty(t *testing.T) {
	m := ParseStyle("")
	if len(m) != 0 {
		t.Errorf("expected empty map, got %d entries", len(m))
	}
}

func TestParseStyleOnlySemicolons(t *testing.T) {
	m := ParseStyle(";;;")
	if len(m) != 0 {
		t.Errorf("expected empty map, got %d entries", len(m))
	}
}

func TestParseStyleNoValue(t *testing.T) {
	m := ParseStyle("color:; font-size: 16px")
	if _, ok := m["color"]; ok {
		t.Error("expected color to be skipped (empty value)")
	}
	assertProp(t, m, "font-size", "16px")
}

func TestParseStyleCaseInsensitiveProperty(t *testing.T) {
	m := ParseStyle("Color: red; FONT-SIZE: 16px")
	assertProp(t, m, "color", "red")
	assertProp(t, m, "font-size", "16px")
}

func TestParseStylePreservesValueCase(t *testing.T) {
	m := ParseStyle("font-family: 'Times New Roman'")
	assertProp(t, m, "font-family", "'Times New Roman'")
}

func TestParseStyleSemicolonInQuotes(t *testing.T) {
	m := ParseStyle(`content: "hello; world"; color: red`)
	assertProp(t, m, "content", `"hello; world"`)
	assertProp(t, m, "color", "red")
}

func TestParseStyleColonInURL(t *testing.T) {
	m := ParseStyle("background: url(data:image/png;base64,abc); color: red")
	assertProp(t, m, "background", "url(data:image/png;base64,abc)")
	assertProp(t, m, "color", "red")
}

func TestParseStyleNestedParens(t *testing.T) {
	m := ParseStyle("filter: drop-shadow(0 0 10px rgba(0,0,0,0.5))")
	assertProp(t, m, "filter", "drop-shadow(0 0 10px rgba(0,0,0,0.5))")
}

func assertProp(t *testing.T, m map[string]string, key, want string) {
	t.Helper()
	got, ok := m[key]
	if !ok {
		t.Errorf("missing property %q", key)
		return
	}
	if got != want {
		t.Errorf("property %q = %q, want %q", key, got, want)
	}
}
