package style

import (
	"math"
	"testing"

	"github.com/macawls/ogre/parse"
)

func approxEqual(a, b, epsilon float64) bool {
	return math.Abs(a-b) < epsilon
}

func TestDefaultDiv(t *testing.T) {
	root := &parse.Node{
		Type:  parse.ElementNode,
		Tag:   "div",
		Style: map[string]string{},
	}
	result := Resolve(root, 1200, 630)
	cs := result[root]

	if cs.Display != DisplayFlex {
		t.Errorf("div display: got %v, want flex", cs.Display)
	}
	if cs.FlexDirection != FlexDirectionRow {
		t.Errorf("div flex-direction: got %v, want row", cs.FlexDirection)
	}
	if cs.Position != PositionRelative {
		t.Errorf("div position: got %v, want relative", cs.Position)
	}
	if cs.BoxSizing != BoxSizingBorderBox {
		t.Errorf("div box-sizing: got %v, want border-box", cs.BoxSizing)
	}
	if cs.FontSize != 16 {
		t.Errorf("root font-size: got %v, want 16", cs.FontSize)
	}
	if cs.FontWeight != 400 {
		t.Errorf("root font-weight: got %v, want 400", cs.FontWeight)
	}
	if cs.Color != (Color{0, 0, 0, 1}) {
		t.Errorf("root color: got %v, want black", cs.Color)
	}
}

func TestDefaultP(t *testing.T) {
	root := &parse.Node{
		Type:  parse.ElementNode,
		Tag:   "div",
		Style: map[string]string{},
		Children: []*parse.Node{
			{
				Type:  parse.ElementNode,
				Tag:   "p",
				Style: map[string]string{},
			},
		},
	}
	result := Resolve(root, 1200, 630)
	p := result[root.Children[0]]

	if p.Display != DisplayFlex {
		t.Errorf("p display: got %v, want flex", p.Display)
	}
	if !approxEqual(p.MarginTop.Raw, 16, 0.01) {
		t.Errorf("p margin-top: got %v, want 16 (1em at 16px)", p.MarginTop.Raw)
	}
	if !approxEqual(p.MarginBottom.Raw, 16, 0.01) {
		t.Errorf("p margin-bottom: got %v, want 16 (1em at 16px)", p.MarginBottom.Raw)
	}
}

func TestDefaultH1(t *testing.T) {
	root := &parse.Node{
		Type:  parse.ElementNode,
		Tag:   "div",
		Style: map[string]string{},
		Children: []*parse.Node{
			{
				Type:  parse.ElementNode,
				Tag:   "h1",
				Style: map[string]string{},
			},
		},
	}
	result := Resolve(root, 1200, 630)
	h1 := result[root.Children[0]]

	if !approxEqual(h1.FontSize, 32, 0.01) {
		t.Errorf("h1 font-size: got %v, want 32 (2em at 16px)", h1.FontSize)
	}
	if h1.FontWeight != 700 {
		t.Errorf("h1 font-weight: got %v, want 700", h1.FontWeight)
	}
	if !approxEqual(h1.MarginTop.Raw, 32*0.67, 0.01) {
		t.Errorf("h1 margin-top: got %v, want %v", h1.MarginTop.Raw, 32*0.67)
	}
}

func TestInheritance(t *testing.T) {
	root := &parse.Node{
		Type:  parse.ElementNode,
		Tag:   "div",
		Style: map[string]string{"color": "red"},
		Children: []*parse.Node{
			{
				Type:  parse.ElementNode,
				Tag:   "div",
				Style: map[string]string{},
			},
		},
	}
	result := Resolve(root, 1200, 630)
	child := result[root.Children[0]]

	expected := Color{255, 0, 0, 1}
	if child.Color != expected {
		t.Errorf("inherited color: got %v, want %v", child.Color, expected)
	}
}

func TestShorthandExpansionAndResolution(t *testing.T) {
	root := &parse.Node{
		Type: parse.ElementNode,
		Tag:  "div",
		Style: map[string]string{
			"margin": "10px 20px",
		},
	}
	result := Resolve(root, 1200, 630)
	cs := result[root]

	if !approxEqual(cs.MarginTop.Raw, 10, 0.01) {
		t.Errorf("margin-top: got %v, want 10", cs.MarginTop.Raw)
	}
	if !approxEqual(cs.MarginRight.Raw, 20, 0.01) {
		t.Errorf("margin-right: got %v, want 20", cs.MarginRight.Raw)
	}
	if !approxEqual(cs.MarginBottom.Raw, 10, 0.01) {
		t.Errorf("margin-bottom: got %v, want 10", cs.MarginBottom.Raw)
	}
	if !approxEqual(cs.MarginLeft.Raw, 20, 0.01) {
		t.Errorf("margin-left: got %v, want 20", cs.MarginLeft.Raw)
	}
}

func TestUnitResolutionEm(t *testing.T) {
	root := &parse.Node{
		Type: parse.ElementNode,
		Tag:  "div",
		Style: map[string]string{
			"font-size": "20px",
			"padding":   "2em",
		},
	}
	result := Resolve(root, 1200, 630)
	cs := result[root]

	if !approxEqual(cs.PaddingTop, 40, 0.01) {
		t.Errorf("padding-top (2em at 20px): got %v, want 40", cs.PaddingTop)
	}
}

func TestUnitResolutionRem(t *testing.T) {
	root := &parse.Node{
		Type: parse.ElementNode,
		Tag:  "div",
		Style: map[string]string{
			"font-size": "20px",
			"padding":   "2rem",
		},
	}
	result := Resolve(root, 1200, 630)
	cs := result[root]

	if !approxEqual(cs.PaddingTop, 32, 0.01) {
		t.Errorf("padding-top (2rem, root=16px): got %v, want 32", cs.PaddingTop)
	}
}

func TestUnitResolutionPercent(t *testing.T) {
	root := &parse.Node{
		Type: parse.ElementNode,
		Tag:  "div",
		Style: map[string]string{
			"width": "50%",
		},
	}
	result := Resolve(root, 1200, 630)
	cs := result[root]

	if cs.Width.Unit != UnitPercent {
		t.Errorf("width unit: got %v, want percent", cs.Width.Unit)
	}
	if !approxEqual(cs.Width.Raw, 50, 0.01) {
		t.Errorf("width raw: got %v, want 50", cs.Width.Raw)
	}
}

func TestOverrideDefault(t *testing.T) {
	root := &parse.Node{
		Type: parse.ElementNode,
		Tag:  "div",
		Style: map[string]string{
			"display":        "none",
			"flex-direction":  "column",
		},
	}
	result := Resolve(root, 1200, 630)
	cs := result[root]

	if cs.Display != DisplayNone {
		t.Errorf("display override: got %v, want none", cs.Display)
	}
	if cs.FlexDirection != FlexDirectionColumn {
		t.Errorf("flex-direction override: got %v, want column", cs.FlexDirection)
	}
}

func TestOverrideInherited(t *testing.T) {
	root := &parse.Node{
		Type:  parse.ElementNode,
		Tag:   "div",
		Style: map[string]string{"color": "red"},
		Children: []*parse.Node{
			{
				Type:  parse.ElementNode,
				Tag:   "div",
				Style: map[string]string{"color": "blue"},
			},
		},
	}
	result := Resolve(root, 1200, 630)
	child := result[root.Children[0]]

	expected := Color{0, 0, 255, 1}
	if child.Color != expected {
		t.Errorf("overridden color: got %v, want %v", child.Color, expected)
	}
}

func TestCSSVariableBasic(t *testing.T) {
	root := &parse.Node{
		Type: parse.ElementNode,
		Tag:  "div",
		Style: map[string]string{
			"--primary": "#ff0000",
			"color":     "var(--primary)",
		},
	}
	result := Resolve(root, 1200, 630)
	cs := result[root]

	expected := Color{255, 0, 0, 1}
	if cs.Color != expected {
		t.Errorf("var(--primary) color: got %v, want %v", cs.Color, expected)
	}
}

func TestCSSVariableFallback(t *testing.T) {
	root := &parse.Node{
		Type: parse.ElementNode,
		Tag:  "div",
		Style: map[string]string{
			"color": "var(--missing, blue)",
		},
	}
	result := Resolve(root, 1200, 630)
	cs := result[root]

	expected := Color{0, 0, 255, 1}
	if cs.Color != expected {
		t.Errorf("var fallback: got %v, want %v", cs.Color, expected)
	}
}

func TestCSSVariableInheritance(t *testing.T) {
	child := &parse.Node{
		Type: parse.ElementNode,
		Tag:  "div",
		Style: map[string]string{
			"color": "var(--theme-color)",
		},
	}
	root := &parse.Node{
		Type: parse.ElementNode,
		Tag:  "div",
		Style: map[string]string{
			"--theme-color": "green",
		},
		Children: []*parse.Node{child},
	}
	result := Resolve(root, 1200, 630)
	cs := result[child]

	expected := Color{0, 128, 0, 1}
	if cs.Color != expected {
		t.Errorf("inherited var: got %v, want %v", cs.Color, expected)
	}
}

func TestCSSVariableOverride(t *testing.T) {
	child := &parse.Node{
		Type: parse.ElementNode,
		Tag:  "div",
		Style: map[string]string{
			"--color": "#0000ff",
			"color":   "var(--color)",
		},
	}
	root := &parse.Node{
		Type: parse.ElementNode,
		Tag:  "div",
		Style: map[string]string{
			"--color": "#ff0000",
		},
		Children: []*parse.Node{child},
	}
	result := Resolve(root, 1200, 630)
	cs := result[child]

	expected := Color{0, 0, 255, 1}
	if cs.Color != expected {
		t.Errorf("overridden var: got %v, want %v", cs.Color, expected)
	}
}

func TestCSSVariableNested(t *testing.T) {
	root := &parse.Node{
		Type: parse.ElementNode,
		Tag:  "div",
		Style: map[string]string{
			"--size":   "20px",
			"padding":  "var(--size)",
		},
	}
	result := Resolve(root, 1200, 630)
	cs := result[root]

	if !approxEqual(cs.PaddingTop, 20, 0.01) {
		t.Errorf("var(--size) padding: got %v, want 20", cs.PaddingTop)
	}
}

func TestResolveVarFunction(t *testing.T) {
	vars := map[string]string{
		"--a": "10px",
		"--b": "red",
	}
	tests := []struct {
		input string
		want  string
	}{
		{"var(--a)", "10px"},
		{"var(--b)", "red"},
		{"var(--missing, 5px)", "5px"},
		{"var(--missing)", ""},
		{"solid var(--a) var(--b)", "solid 10px red"},
		{"no-var-here", "no-var-here"},
	}
	for _, tt := range tests {
		got := resolveVar(tt.input, vars)
		if got != tt.want {
			t.Errorf("resolveVar(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestHTMLDefaultA(t *testing.T) {
	root := &parse.Node{
		Type:  parse.ElementNode,
		Tag:   "div",
		Style: map[string]string{},
		Children: []*parse.Node{
			{Type: parse.ElementNode, Tag: "a", Style: map[string]string{}},
		},
	}
	result := Resolve(root, 1200, 630)
	a := result[root.Children[0]]

	expected := Color{0, 0, 238, 1}
	if a.Color != expected {
		t.Errorf("a color: got %v, want %v", a.Color, expected)
	}
	if a.TextDecorationLine != TextDecorationUnderline {
		t.Errorf("a text-decoration: got %v, want underline", a.TextDecorationLine)
	}
}

func TestHTMLDefaultBlockquote(t *testing.T) {
	root := &parse.Node{
		Type:  parse.ElementNode,
		Tag:   "div",
		Style: map[string]string{},
		Children: []*parse.Node{
			{Type: parse.ElementNode, Tag: "blockquote", Style: map[string]string{}},
		},
	}
	result := Resolve(root, 1200, 630)
	bq := result[root.Children[0]]

	if !approxEqual(bq.MarginLeft.Raw, 40, 0.01) {
		t.Errorf("blockquote margin-left: got %v, want 40", bq.MarginLeft.Raw)
	}
	if !approxEqual(bq.MarginRight.Raw, 40, 0.01) {
		t.Errorf("blockquote margin-right: got %v, want 40", bq.MarginRight.Raw)
	}
}

func TestHTMLDefaultUL(t *testing.T) {
	root := &parse.Node{
		Type:  parse.ElementNode,
		Tag:   "div",
		Style: map[string]string{},
		Children: []*parse.Node{
			{Type: parse.ElementNode, Tag: "ul", Style: map[string]string{}},
		},
	}
	result := Resolve(root, 1200, 630)
	ul := result[root.Children[0]]

	if !approxEqual(ul.PaddingLeft, 40, 0.01) {
		t.Errorf("ul padding-left: got %v, want 40", ul.PaddingLeft)
	}
}

func TestHTMLDefaultSummary(t *testing.T) {
	root := &parse.Node{
		Type:  parse.ElementNode,
		Tag:   "div",
		Style: map[string]string{},
		Children: []*parse.Node{
			{Type: parse.ElementNode, Tag: "summary", Style: map[string]string{}},
		},
	}
	result := Resolve(root, 1200, 630)
	s := result[root.Children[0]]

	if s.FontWeight != 700 {
		t.Errorf("summary font-weight: got %v, want 700", s.FontWeight)
	}
}

func TestHTMLDefaultSup(t *testing.T) {
	root := &parse.Node{
		Type:  parse.ElementNode,
		Tag:   "div",
		Style: map[string]string{},
		Children: []*parse.Node{
			{Type: parse.ElementNode, Tag: "sup", Style: map[string]string{}},
		},
	}
	result := Resolve(root, 1200, 630)
	sup := result[root.Children[0]]

	if !approxEqual(sup.FontSize, 16*0.83, 0.1) {
		t.Errorf("sup font-size: got %v, want %v", sup.FontSize, 16*0.83)
	}
}

func TestHTMLDefaultCenter(t *testing.T) {
	root := &parse.Node{
		Type:  parse.ElementNode,
		Tag:   "div",
		Style: map[string]string{},
		Children: []*parse.Node{
			{Type: parse.ElementNode, Tag: "center", Style: map[string]string{}},
		},
	}
	result := Resolve(root, 1200, 630)
	c := result[root.Children[0]]

	if c.TextAlign != TextAlignCenter {
		t.Errorf("center text-align: got %v, want center", c.TextAlign)
	}
}

func TestHTMLDefaultDel(t *testing.T) {
	root := &parse.Node{
		Type:  parse.ElementNode,
		Tag:   "div",
		Style: map[string]string{},
		Children: []*parse.Node{
			{Type: parse.ElementNode, Tag: "del", Style: map[string]string{}},
		},
	}
	result := Resolve(root, 1200, 630)
	del := result[root.Children[0]]

	if del.TextDecorationLine != TextDecorationLineThrough {
		t.Errorf("del text-decoration: got %v, want line-through", del.TextDecorationLine)
	}
}
