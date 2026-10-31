package layout

import (
	"math"
	"testing"
)

func assertLayout(t *testing.T, node *Node, name string, x, y, w, h float64) {
	t.Helper()
	l := node.Layout
	if math.Abs(l.X-x) > 0.1 || math.Abs(l.Y-y) > 0.1 || math.Abs(l.Width-w) > 0.1 || math.Abs(l.Height-h) > 0.1 {
		t.Errorf("%s: got (%.1f, %.1f, %.1f, %.1f), want (%.1f, %.1f, %.1f, %.1f)", name, l.X, l.Y, l.Width, l.Height, x, y, w, h)
	}
}

func TestSingleChildFillsContainer(t *testing.T) {
	child := NewNode(Style{FlexGrow: 1})
	root := NewNode(Style{
		Width:  Pt(400),
		Height: Pt(300),
	}, child)

	Compute(root, 400, 300)

	assertLayout(t, root, "root", 0, 0, 400, 300)
	assertLayout(t, child, "child", 0, 0, 400, 300)
}

func TestRowDirection(t *testing.T) {
	c1 := NewNode(Style{Width: Pt(100)})
	c2 := NewNode(Style{Width: Pt(150)})
	c3 := NewNode(Style{Width: Pt(80)})
	root := NewNode(Style{
		Width:     Pt(400),
		Height:    Pt(100),
		Direction: Row,
	}, c1, c2, c3)

	Compute(root, 400, 100)

	assertLayout(t, root, "root", 0, 0, 400, 100)
	assertLayout(t, c1, "c1", 0, 0, 100, 100)
	assertLayout(t, c2, "c2", 100, 0, 150, 100)
	assertLayout(t, c3, "c3", 250, 0, 80, 100)
}

func TestColumnDirection(t *testing.T) {
	c1 := NewNode(Style{Height: Pt(50)})
	c2 := NewNode(Style{Height: Pt(70)})
	c3 := NewNode(Style{Height: Pt(30)})
	root := NewNode(Style{
		Width:     Pt(200),
		Height:    Pt(300),
		Direction: Column,
	}, c1, c2, c3)

	Compute(root, 200, 300)

	assertLayout(t, root, "root", 0, 0, 200, 300)
	assertLayout(t, c1, "c1", 0, 0, 200, 50)
	assertLayout(t, c2, "c2", 0, 50, 200, 70)
	assertLayout(t, c3, "c3", 0, 120, 200, 30)
}

func TestFlexGrow(t *testing.T) {
	c1 := NewNode(Style{FlexGrow: 1})
	c2 := NewNode(Style{FlexGrow: 2})
	c3 := NewNode(Style{FlexGrow: 1})
	root := NewNode(Style{
		Width:     Pt(400),
		Height:    Pt(100),
		Direction: Row,
	}, c1, c2, c3)

	Compute(root, 400, 100)

	assertLayout(t, c1, "c1", 0, 0, 100, 100)
	assertLayout(t, c2, "c2", 100, 0, 200, 100)
	assertLayout(t, c3, "c3", 300, 0, 100, 100)
}

func TestFlexShrink(t *testing.T) {
	// Use FlexBasis instead of Width to avoid the explicit Width overriding
	// the flex-resolved size during recursive computeNode.
	c1 := NewNode(Style{FlexBasis: Pt(200), FlexShrink: 1})
	c2 := NewNode(Style{FlexBasis: Pt(200), FlexShrink: 2})
	c3 := NewNode(Style{FlexBasis: Pt(200), FlexShrink: 1})
	root := NewNode(Style{
		Width:     Pt(400),
		Height:    Pt(100),
		Direction: Row,
	}, c1, c2, c3)

	Compute(root, 400, 100)

	// Total base = 600, container = 400, overflow = 200
	// Scaled shrink: c1=200*1=200, c2=200*2=400, c3=200*1=200, total=800
	// c1: 200 - 200/800*200 = 150
	// c2: 200 - 400/800*200 = 100
	// c3: 200 - 200/800*200 = 150
	assertLayout(t, c1, "c1", 0, 0, 150, 100)
	assertLayout(t, c2, "c2", 150, 0, 100, 100)
	assertLayout(t, c3, "c3", 250, 0, 150, 100)
}

func TestJustifyContentStart(t *testing.T) {
	children := make([]*Node, 3)
	for i := range children {
		children[i] = NewNode(Style{Width: Pt(50)})
	}
	root := NewNode(Style{
		Width:          Pt(400),
		Height:         Pt(100),
		Direction:      Row,
		JustifyContent: JustifyStart,
	}, children...)

	Compute(root, 400, 100)

	assertLayout(t, children[0], "c0", 0, 0, 50, 100)
	assertLayout(t, children[1], "c1", 50, 0, 50, 100)
	assertLayout(t, children[2], "c2", 100, 0, 50, 100)
}

func TestJustifyContentEnd(t *testing.T) {
	children := make([]*Node, 3)
	for i := range children {
		children[i] = NewNode(Style{Width: Pt(50)})
	}
	root := NewNode(Style{
		Width:          Pt(400),
		Height:         Pt(100),
		Direction:      Row,
		JustifyContent: JustifyEnd,
	}, children...)

	Compute(root, 400, 100)

	assertLayout(t, children[0], "c0", 250, 0, 50, 100)
	assertLayout(t, children[1], "c1", 300, 0, 50, 100)
	assertLayout(t, children[2], "c2", 350, 0, 50, 100)
}

func TestJustifyContentCenter(t *testing.T) {
	children := make([]*Node, 3)
	for i := range children {
		children[i] = NewNode(Style{Width: Pt(50)})
	}
	root := NewNode(Style{
		Width:          Pt(400),
		Height:         Pt(100),
		Direction:      Row,
		JustifyContent: JustifyCenter,
	}, children...)

	Compute(root, 400, 100)

	assertLayout(t, children[0], "c0", 125, 0, 50, 100)
	assertLayout(t, children[1], "c1", 175, 0, 50, 100)
	assertLayout(t, children[2], "c2", 225, 0, 50, 100)
}

func TestJustifyContentSpaceBetween(t *testing.T) {
	children := make([]*Node, 3)
	for i := range children {
		children[i] = NewNode(Style{Width: Pt(50)})
	}
	root := NewNode(Style{
		Width:          Pt(400),
		Height:         Pt(100),
		Direction:      Row,
		JustifyContent: JustifySpaceBetween,
	}, children...)

	Compute(root, 400, 100)

	assertLayout(t, children[0], "c0", 0, 0, 50, 100)
	assertLayout(t, children[1], "c1", 175, 0, 50, 100)
	assertLayout(t, children[2], "c2", 350, 0, 50, 100)
}

func TestJustifyContentSpaceAround(t *testing.T) {
	children := make([]*Node, 3)
	for i := range children {
		children[i] = NewNode(Style{Width: Pt(50)})
	}
	root := NewNode(Style{
		Width:          Pt(400),
		Height:         Pt(100),
		Direction:      Row,
		JustifyContent: JustifySpaceAround,
	}, children...)

	Compute(root, 400, 100)

	assertLayout(t, children[0], "c0", 41.7, 0, 50, 100)
	assertLayout(t, children[1], "c1", 175, 0, 50, 100)
	assertLayout(t, children[2], "c2", 308.3, 0, 50, 100)
}

func TestJustifyContentSpaceEvenly(t *testing.T) {
	children := make([]*Node, 3)
	for i := range children {
		children[i] = NewNode(Style{Width: Pt(50)})
	}
	root := NewNode(Style{
		Width:          Pt(400),
		Height:         Pt(100),
		Direction:      Row,
		JustifyContent: JustifySpaceEvenly,
	}, children...)

	Compute(root, 400, 100)

	assertLayout(t, children[0], "c0", 62.5, 0, 50, 100)
	assertLayout(t, children[1], "c1", 175, 0, 50, 100)
	assertLayout(t, children[2], "c2", 287.5, 0, 50, 100)
}

func TestAlignItemsStart(t *testing.T) {
	c1 := NewNode(Style{Width: Pt(50), Height: Pt(30)})
	c2 := NewNode(Style{Width: Pt(50), Height: Pt(60)})
	c3 := NewNode(Style{Width: Pt(50), Height: Pt(40)})
	root := NewNode(Style{
		Width:      Pt(300),
		Height:     Pt(100),
		Direction:  Row,
		AlignItems: AlignStart,
	}, c1, c2, c3)

	Compute(root, 300, 100)

	assertLayout(t, c1, "c1", 0, 0, 50, 30)
	assertLayout(t, c2, "c2", 50, 0, 50, 60)
	assertLayout(t, c3, "c3", 100, 0, 50, 40)
}

func TestAlignItemsEnd(t *testing.T) {
	c1 := NewNode(Style{Width: Pt(50), Height: Pt(30)})
	c2 := NewNode(Style{Width: Pt(50), Height: Pt(60)})
	c3 := NewNode(Style{Width: Pt(50), Height: Pt(40)})
	root := NewNode(Style{
		Width:      Pt(300),
		Height:     Pt(100),
		Direction:  Row,
		AlignItems: AlignEnd,
	}, c1, c2, c3)

	Compute(root, 300, 100)

	assertLayout(t, c1, "c1", 0, 70, 50, 30)
	assertLayout(t, c2, "c2", 50, 40, 50, 60)
	assertLayout(t, c3, "c3", 100, 60, 50, 40)
}

func TestAlignItemsCenter(t *testing.T) {
	c1 := NewNode(Style{Width: Pt(50), Height: Pt(30)})
	c2 := NewNode(Style{Width: Pt(50), Height: Pt(60)})
	c3 := NewNode(Style{Width: Pt(50), Height: Pt(40)})
	root := NewNode(Style{
		Width:      Pt(300),
		Height:     Pt(100),
		Direction:  Row,
		AlignItems: AlignCenter,
	}, c1, c2, c3)

	Compute(root, 300, 100)

	assertLayout(t, c1, "c1", 0, 35, 50, 30)
	assertLayout(t, c2, "c2", 50, 20, 50, 60)
	assertLayout(t, c3, "c3", 100, 30, 50, 40)
}

func TestAlignItemsStretch(t *testing.T) {
	c1 := NewNode(Style{Width: Pt(50)})
	c2 := NewNode(Style{Width: Pt(50)})
	root := NewNode(Style{
		Width:      Pt(300),
		Height:     Pt(100),
		Direction:  Row,
		AlignItems: AlignStretch,
	}, c1, c2)

	Compute(root, 300, 100)

	assertLayout(t, c1, "c1", 0, 0, 50, 100)
	assertLayout(t, c2, "c2", 50, 0, 50, 100)
}

func TestFlexWrap(t *testing.T) {
	c1 := NewNode(Style{Width: Pt(150), Height: Pt(40)})
	c2 := NewNode(Style{Width: Pt(150), Height: Pt(40)})
	c3 := NewNode(Style{Width: Pt(150), Height: Pt(40)})
	root := NewNode(Style{
		Width:        Pt(300),
		Height:       Pt(200),
		Direction:    Row,
		Wrap:         WrapWrap,
		AlignContent: AlignStretch,
	}, c1, c2, c3)

	Compute(root, 300, 200)

	// Line 1: c1(150) + c2(150) = 300
	// Line 2: c3(150)
	// With align-content:stretch, remaining 120 (200-40-40) split across 2 lines: +60 each
	// Line 1 cross = 100, line 2 cross = 100
	assertLayout(t, c1, "c1", 0, 0, 150, 40)
	assertLayout(t, c2, "c2", 150, 0, 150, 40)
	assertLayout(t, c3, "c3", 0, 100, 150, 40)
}

func TestGap(t *testing.T) {
	c1 := NewNode(Style{Width: Pt(100)})
	c2 := NewNode(Style{Width: Pt(100)})
	c3 := NewNode(Style{Width: Pt(100)})
	root := NewNode(Style{
		Width:     Pt(400),
		Height:    Pt(100),
		Direction: Row,
		Gap:       20,
	}, c1, c2, c3)

	Compute(root, 400, 100)

	assertLayout(t, c1, "c1", 0, 0, 100, 100)
	assertLayout(t, c2, "c2", 120, 0, 100, 100)
	assertLayout(t, c3, "c3", 240, 0, 100, 100)
}

func TestPadding(t *testing.T) {
	child := NewNode(Style{Width: Pt(100)})
	root := NewNode(Style{
		Width:         Pt(400),
		Height:        Pt(300),
		Direction:     Row,
		PaddingTop:    Pt(20),
		PaddingLeft:   Pt(30),
		PaddingRight:  Pt(30),
		PaddingBottom: Pt(20),
	}, child)

	Compute(root, 400, 300)

	assertLayout(t, root, "root", 0, 0, 400, 300)
	assertLayout(t, child, "child", 30, 20, 100, 260)
}

func TestMargin(t *testing.T) {
	c1 := NewNode(Style{
		Width:       Pt(100),
		MarginRight: Pt(20),
	})
	c2 := NewNode(Style{
		Width:      Pt(100),
		MarginLeft: Pt(10),
	})
	root := NewNode(Style{
		Width:     Pt(400),
		Height:    Pt(100),
		Direction: Row,
	}, c1, c2)

	Compute(root, 400, 100)

	assertLayout(t, c1, "c1", 0, 0, 100, 100)
	assertLayout(t, c2, "c2", 130, 0, 100, 100)
}

func TestAutoMarginPushRight(t *testing.T) {
	c1 := NewNode(Style{
		Width:      Pt(50),
		MarginLeft: Auto(),
	})
	root := NewNode(Style{
		Width:     Pt(400),
		Height:    Pt(100),
		Direction: Row,
	}, c1)

	Compute(root, 400, 100)

	assertLayout(t, c1, "c1", 350, 0, 50, 100)
}

func TestNestedFlex(t *testing.T) {
	innerChild1 := NewNode(Style{FlexGrow: 1})
	innerChild2 := NewNode(Style{FlexGrow: 1})
	inner := NewNode(Style{
		Width:     Pt(200),
		Height:    Pt(100),
		Direction: Row,
	}, innerChild1, innerChild2)

	outer := NewNode(Style{
		Width:     Pt(400),
		Height:    Pt(100),
		Direction: Row,
	}, inner)

	Compute(outer, 400, 100)

	assertLayout(t, outer, "outer", 0, 0, 400, 100)
	assertLayout(t, inner, "inner", 0, 0, 200, 100)
	assertLayout(t, innerChild1, "innerChild1", 0, 0, 100, 100)
	assertLayout(t, innerChild2, "innerChild2", 100, 0, 100, 100)
}

func TestMeasureFunc(t *testing.T) {
	leaf := NewLeaf(Style{}, func(maxW, maxH float64) (float64, float64) {
		return 80, 20
	})
	root := NewNode(Style{
		Width:      Pt(400),
		Height:     Pt(100),
		Direction:  Row,
		AlignItems: AlignStart,
	}, leaf)

	Compute(root, 400, 100)

	assertLayout(t, leaf, "leaf", 0, 0, 80, 20)
}

func TestMeasureFuncWithStretch(t *testing.T) {
	leaf := NewLeaf(Style{}, func(maxW, maxH float64) (float64, float64) {
		return 80, 20
	})
	root := NewNode(Style{
		Width:     Pt(400),
		Height:    Pt(100),
		Direction: Row,
	}, leaf)

	Compute(root, 400, 100)

	assertLayout(t, leaf, "leaf", 0, 0, 80, 100)
}

func TestAbsolutePositioning(t *testing.T) {
	child := NewNode(Style{
		Position: Absolute,
		Width:    Pt(50),
		Height:   Pt(50),
		Top:      Pt(10),
		Left:     Pt(10),
	})
	root := NewNode(Style{
		Width:  Pt(400),
		Height: Pt(300),
	}, child)

	Compute(root, 400, 300)

	assertLayout(t, root, "root", 0, 0, 400, 300)
	assertLayout(t, child, "child", 10, 10, 50, 50)
}

func TestPercentageWidth(t *testing.T) {
	child := NewNode(Style{
		FlexBasis: Pct(50),
	})
	root := NewNode(Style{
		Width:     Pt(400),
		Height:    Pt(200),
		Direction: Row,
	}, child)

	Compute(root, 400, 200)

	assertLayout(t, child, "child", 0, 0, 200, 200)
}

// Regression for the double-percent-resolution bug: a flex child with an
// explicit Width: Pct(50) should occupy 50% of the parent, not 25%. Before
// the fix, computeNode re-resolved s.Width against the child's already-
// allocated slot (which was itself 50% of the parent), producing 25%.
func TestPercentageWidthOnFlexChild(t *testing.T) {
	left := NewNode(Style{Width: Pct(50)})
	right := NewNode(Style{Width: Pct(50)})
	root := NewNode(Style{
		Width:     Pt(400),
		Height:    Pt(200),
		Direction: Row,
	}, left, right)

	Compute(root, 400, 200)

	assertLayout(t, left, "left", 0, 0, 200, 200)
	assertLayout(t, right, "right", 200, 0, 200, 200)
}

// Regression for the flex-child height-shrink bug: a flex-column child with
// no explicit height should stretch to the parent's cross-axis size when
// align-items:stretch (the default) is in effect. Before the fix, computeNode
// shrunk Layout.Height to the sum of its children's main sizes even though
// the parent had already stretched it.
func TestFlexChildStretchNoShrink(t *testing.T) {
	inner1 := NewNode(Style{Width: Pt(100), Height: Pt(50)})
	inner2 := NewNode(Style{Width: Pt(100), Height: Pt(50)})
	col := NewNode(Style{
		Width:     Pct(50),
		Direction: Column,
	}, inner1, inner2)
	root := NewNode(Style{
		Width:     Pt(400),
		Height:    Pt(400),
		Direction: Row,
	}, col)

	Compute(root, 400, 400)

	assertLayout(t, col, "col", 0, 0, 200, 400)
}

func TestPercentageHeight(t *testing.T) {
	child := NewNode(Style{
		FlexBasis: Pct(50),
		Height:    Pct(25),
	})
	root := NewNode(Style{
		Width:      Pt(400),
		Height:     Pt(200),
		Direction:  Row,
		AlignItems: AlignStart,
	}, child)

	Compute(root, 400, 200)

	assertLayout(t, child, "child", 0, 0, 200, 50)
}

func TestMinMaxConstraints(t *testing.T) {
	child := NewNode(Style{
		FlexGrow: 1,
		MaxWidth: Pt(250),
	})
	root := NewNode(Style{
		Width:     Pt(400),
		Height:    Pt(100),
		Direction: Row,
	}, child)

	Compute(root, 400, 100)

	assertLayout(t, child, "child", 0, 0, 250, 100)
}

func TestAspectRatio(t *testing.T) {
	child := NewNode(Style{
		Width:       Pt(100),
		AspectRatio: 2.0,
	})
	root := NewNode(Style{
		Width:      Pt(400),
		Height:     Pt(300),
		Direction:  Row,
		AlignItems: AlignStart,
	}, child)

	Compute(root, 400, 300)

	assertLayout(t, child, "child", 0, 0, 100, 50)
}

// Regression for #40. An absolutely positioned box whose size is auto was laid
// out against a zero-sized box: a single bottom/right offset displaced the box
// by its own size (the offset landed on the wrong edge), and Layout.Width and
// Layout.Height stayed 0, so the box's own background and border painted
// nothing while its content overflowed the empty box. The failure was silent.
//
// Every want below is Chromium's getBoundingClientRect() for the same
// declarations on a 400x200 container, so the numbers come from the browser
// rather than from this engine's own output.
func TestAbsoluteAutoSize(t *testing.T) {
	// text stands in for a 14px line of two digits: 18x16 in Chromium.
	text := func(s Style) *Node {
		s.Position = Absolute
		return NewLeaf(s, func(maxW, maxH float64) (float64, float64) {
			return 18, 16
		})
	}
	// boxed stands in for an element with 4px padding wrapping a 40x30 child.
	boxed := func(s Style) *Node {
		s.Position = Absolute
		s.PaddingTop, s.PaddingRight = Pt(4), Pt(4)
		s.PaddingBottom, s.PaddingLeft = Pt(4), Pt(4)
		return NewNode(s, NewNode(Style{Width: Pt(40), Height: Pt(30)}))
	}

	tests := []struct {
		name  string
		child *Node
		want  [4]float64 // x, y, w, h
	}{
		{"bottom+left, auto size", text(Style{Left: Pt(20), Bottom: Pt(20)}), [4]float64{20, 164, 18, 16}},
		{"bottom+left, auto height", text(Style{Left: Pt(20), Bottom: Pt(20), Width: Pt(60)}), [4]float64{20, 164, 60, 16}},
		{"top+left, auto size", text(Style{Left: Pt(20), Top: Pt(20)}), [4]float64{20, 20, 18, 16}},
		{"top+right, auto size", text(Style{Top: Pt(20), Right: Pt(20)}), [4]float64{362, 20, 18, 16}},
		{"left+right+bottom, auto height", text(Style{Left: Pt(20), Right: Pt(20), Bottom: Pt(20)}), [4]float64{20, 164, 360, 16}},
		{"bottom+padding, auto size", boxed(Style{Left: Pt(20), Bottom: Pt(20)}), [4]float64{20, 142, 48, 38}},
		// Both offsets on an axis resolve the size, which already worked; only
		// the auto axis was wrong. Kept so a fix cannot pass by breaking it.
		{"top+bottom+left, auto width", text(Style{Top: Pt(20), Bottom: Pt(20), Left: Pt(20)}), [4]float64{20, 20, 18, 160}},
		// Control: an explicit size already positioned correctly, which is why
		// the only absolute fixture in the corpus never caught this defect.
		{"explicit size, bottom+right", NewNode(Style{Position: Absolute, Right: Pt(10), Bottom: Pt(10), Width: Pt(100), Height: Pt(50)}), [4]float64{290, 140, 100, 50}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			root := NewNode(Style{Width: Pt(400), Height: Pt(200)}, tc.child)
			Compute(root, 400, 200)

			assertLayout(t, tc.child, tc.name, tc.want[0], tc.want[1], tc.want[2], tc.want[3])
		})
	}
}
