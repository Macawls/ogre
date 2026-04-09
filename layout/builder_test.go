package layout

import (
	"testing"

	"github.com/macawls/ogre/parse"
	"github.com/macawls/ogre/style"
)

func dummyMeasure(pn *parse.Node, text string, cs *style.ComputedStyle, maxWidth float64) (float64, float64) {
	w := float64(len(text)) * 8
	if maxWidth > 0 && w > maxWidth {
		w = maxWidth
	}
	return w, 16
}

func TestBuildTreeBasic(t *testing.T) {
	root := &parse.Node{
		Type: parse.ElementNode,
		Tag:  "div",
	}
	styles := map[*parse.Node]*style.ComputedStyle{
		root: {
			Display:       style.DisplayFlex,
			FlexDirection: style.FlexDirectionRow,
		},
	}

	tree := BuildTree(root, styles, dummyMeasure)
	if tree.Root == nil {
		t.Fatal("expected non-nil root")
	}
	if tree.Root.Style.Display != DisplayFlex {
		t.Errorf("expected DisplayFlex, got %d", tree.Root.Style.Display)
	}
	if tree.Root.Style.Direction != Row {
		t.Errorf("expected Row, got %d", tree.Root.Style.Direction)
	}
	if tree.NodeMap[root] != tree.Root {
		t.Error("NodeMap should map DOM root to layout root")
	}
}

func TestBuildTreeDisplayNone(t *testing.T) {
	child := &parse.Node{Type: parse.ElementNode, Tag: "div"}
	root := &parse.Node{
		Type:     parse.ElementNode,
		Tag:      "div",
		Children: []*parse.Node{child},
	}
	styles := map[*parse.Node]*style.ComputedStyle{
		root:  {Display: style.DisplayFlex},
		child: {Display: style.DisplayNone},
	}

	tree := BuildTree(root, styles, dummyMeasure)
	if len(tree.Root.Children) != 0 {
		t.Errorf("expected 0 children for display:none child, got %d", len(tree.Root.Children))
	}
}

func TestBuildTreeTextNode(t *testing.T) {
	textNode := &parse.Node{
		Type: parse.TextNode,
		Text: "hello",
	}
	root := &parse.Node{
		Type:     parse.ElementNode,
		Tag:      "div",
		Children: []*parse.Node{textNode},
	}
	cs := &style.ComputedStyle{Display: style.DisplayFlex}
	styles := map[*parse.Node]*style.ComputedStyle{
		root:     cs,
		textNode: cs,
	}

	tree := BuildTree(root, styles, dummyMeasure)
	if len(tree.Root.Children) != 1 {
		t.Fatalf("expected 1 child, got %d", len(tree.Root.Children))
	}
	leaf := tree.Root.Children[0]
	if leaf.Measure == nil {
		t.Fatal("text node should have MeasureFunc")
	}
	w, h := leaf.Measure(100, 100)
	if w != 40 || h != 16 {
		t.Errorf("expected measure (40,16), got (%g,%g)", w, h)
	}
}

func TestMapDimension(t *testing.T) {
	tests := []struct {
		input style.Value
		want  Dimension
	}{
		{style.Value{Raw: 100, Unit: style.UnitPx}, Pt(100)},
		{style.Value{Raw: 50, Unit: style.UnitPercent}, Pct(50)},
		{style.Value{Unit: style.UnitAuto}, Auto()},
		{style.Value{Unit: style.UnitNone}, Undefined()},
		{style.Value{}, Pt(0)},
	}
	for _, tt := range tests {
		got := mapDimension(tt.input)
		if got != tt.want {
			t.Errorf("mapDimension(%+v) = %+v, want %+v", tt.input, got, tt.want)
		}
	}
}

func TestMapStyleProperties(t *testing.T) {
	cs := &style.ComputedStyle{
		Display:        style.DisplayFlex,
		Position:       style.PositionAbsolute,
		FlexDirection:  style.FlexDirectionColumn,
		FlexWrap:       style.FlexWrapWrap,
		JustifyContent: style.JustifyContentCenter,
		AlignItems:     style.AlignItemsCenter,
		AlignSelf:      style.AlignSelfStretch,
		AlignContent:   style.AlignContentSpaceBetween,
		FlexGrow:       2,
		FlexShrink:     0.5,
		Width:          style.Value{Raw: 200, Unit: style.UnitPx},
		Height:         style.Value{Raw: 50, Unit: style.UnitPercent},
		PaddingTop:     10,
		PaddingRight:   20,
		BorderTopWidth: 1,
		Gap:            5,
		AspectRatio:    1.5,
	}
	s := mapStyle(cs)

	if s.Display != DisplayFlex {
		t.Errorf("Display: got %d, want %d", s.Display, DisplayFlex)
	}
	if s.Position != Absolute {
		t.Errorf("Position: got %d, want %d", s.Position, Absolute)
	}
	if s.Direction != Column {
		t.Errorf("Direction: got %d, want %d", s.Direction, Column)
	}
	if s.Wrap != WrapWrap {
		t.Errorf("Wrap: got %d, want %d", s.Wrap, WrapWrap)
	}
	if s.JustifyContent != JustifyCenter {
		t.Errorf("JustifyContent: got %d, want %d", s.JustifyContent, JustifyCenter)
	}
	if s.AlignItems != AlignCenter {
		t.Errorf("AlignItems: got %d, want %d", s.AlignItems, AlignCenter)
	}
	if s.AlignSelf != AlignStretch {
		t.Errorf("AlignSelf: got %d, want %d", s.AlignSelf, AlignStretch)
	}
	if s.AlignContent != AlignAuto {
		t.Errorf("AlignContent: got %d, want %d", s.AlignContent, AlignAuto)
	}
	if s.FlexGrow != 2 {
		t.Errorf("FlexGrow: got %g, want 2", s.FlexGrow)
	}
	if s.FlexShrink != 0.5 {
		t.Errorf("FlexShrink: got %g, want 0.5", s.FlexShrink)
	}
	if s.Width != Pt(200) {
		t.Errorf("Width: got %+v, want Pt(200)", s.Width)
	}
	if s.Height != Pct(50) {
		t.Errorf("Height: got %+v, want Pct(50)", s.Height)
	}
	if s.PaddingTop != Pt(10) {
		t.Errorf("PaddingTop: got %+v, want Pt(10)", s.PaddingTop)
	}
	if s.PaddingRight != Pt(20) {
		t.Errorf("PaddingRight: got %+v, want Pt(20)", s.PaddingRight)
	}
	if s.BorderTop != 1 {
		t.Errorf("BorderTop: got %g, want 1", s.BorderTop)
	}
	if s.Gap != 5 {
		t.Errorf("Gap: got %g, want 5", s.Gap)
	}
	if s.AspectRatio != 1.5 {
		t.Errorf("AspectRatio: got %g, want 1.5", s.AspectRatio)
	}
}

func TestBuildTreeNestedChildren(t *testing.T) {
	grandchild := &parse.Node{Type: parse.ElementNode, Tag: "div"}
	child := &parse.Node{
		Type:     parse.ElementNode,
		Tag:      "div",
		Children: []*parse.Node{grandchild},
	}
	root := &parse.Node{
		Type:     parse.ElementNode,
		Tag:      "div",
		Children: []*parse.Node{child},
	}
	cs := &style.ComputedStyle{Display: style.DisplayFlex}
	styles := map[*parse.Node]*style.ComputedStyle{
		root:       cs,
		child:      cs,
		grandchild: cs,
	}

	tree := BuildTree(root, styles, dummyMeasure)
	if len(tree.Root.Children) != 1 {
		t.Fatalf("root: expected 1 child, got %d", len(tree.Root.Children))
	}
	if len(tree.Root.Children[0].Children) != 1 {
		t.Fatalf("child: expected 1 child, got %d", len(tree.Root.Children[0].Children))
	}
	if tree.NodeMap[grandchild] == nil {
		t.Error("grandchild should be in NodeMap")
	}
}

func TestComputeLayout(t *testing.T) {
	root := &parse.Node{
		Type: parse.ElementNode,
		Tag:  "div",
	}
	styles := map[*parse.Node]*style.ComputedStyle{
		root: {
			Display: style.DisplayFlex,
			Width:   style.Value{Raw: 800, Unit: style.UnitPx},
			Height:  style.Value{Raw: 600, Unit: style.UnitPx},
		},
	}

	tree := ComputeLayout(root, styles, 800, 600, dummyMeasure)
	if tree.Root == nil {
		t.Fatal("expected non-nil root")
	}
}
