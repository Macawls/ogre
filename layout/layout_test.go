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
