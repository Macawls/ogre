package layout

import "testing"

func gridChild(styles ...func(*Style)) *Node {
	s := Style{Display: DisplayFlex}
	for _, fn := range styles {
		fn(&s)
	}
	return NewNode(s)
}

func TestGridFixedTwoByTwo(t *testing.T) {
	c1 := gridChild()
	c2 := gridChild()
	c3 := gridChild()
	c4 := gridChild()

	root := NewNode(Style{
		Display: DisplayGrid,
		Width:   Pt(400),
		Height:  Pt(200),
		GridTemplateColumns: TrackList{Tracks: []TrackSize{
			{Kind: TrackFixed, Length: Pt(200)},
			{Kind: TrackFixed, Length: Pt(200)},
		}},
		GridTemplateRows: TrackList{Tracks: []TrackSize{
			{Kind: TrackFixed, Length: Pt(100)},
			{Kind: TrackFixed, Length: Pt(100)},
		}},
	}, c1, c2, c3, c4)

	Compute(root, 400, 200)

	assertLayout(t, c1, "c1", 0, 0, 200, 100)
	assertLayout(t, c2, "c2", 200, 0, 200, 100)
	assertLayout(t, c3, "c3", 0, 100, 200, 100)
	assertLayout(t, c4, "c4", 200, 100, 200, 100)
}

func TestGridFrDistribution(t *testing.T) {
	a := gridChild()
	b := gridChild()
	c := gridChild()

	root := NewNode(Style{
		Display: DisplayGrid,
		Width:   Pt(600),
		Height:  Pt(100),
		GridTemplateColumns: TrackList{Tracks: []TrackSize{
			{Kind: TrackFr, Fr: 1},
			{Kind: TrackFr, Fr: 2},
			{Kind: TrackFr, Fr: 3},
		}},
		GridTemplateRows: TrackList{Tracks: []TrackSize{
			{Kind: TrackFixed, Length: Pt(100)},
		}},
	}, a, b, c)

	Compute(root, 600, 100)

	assertLayout(t, a, "a", 0, 0, 100, 100)
	assertLayout(t, b, "b", 100, 0, 200, 100)
	assertLayout(t, c, "c", 300, 0, 300, 100)
}

func TestGridGap(t *testing.T) {
	c1 := gridChild()
	c2 := gridChild()

	root := NewNode(Style{
		Display:   DisplayGrid,
		Width:     Pt(210),
		Height:    Pt(100),
		ColumnGap: 10,
		GridTemplateColumns: TrackList{Tracks: []TrackSize{
			{Kind: TrackFixed, Length: Pt(100)},
			{Kind: TrackFixed, Length: Pt(100)},
		}},
		GridTemplateRows: TrackList{Tracks: []TrackSize{
			{Kind: TrackFixed, Length: Pt(100)},
		}},
	}, c1, c2)

	Compute(root, 210, 100)

	assertLayout(t, c1, "c1", 0, 0, 100, 100)
	assertLayout(t, c2, "c2", 110, 0, 100, 100)
}

func TestGridExplicitPlacement(t *testing.T) {
	a := gridChild(func(s *Style) {
		s.GridColumnStart = GridPlacement{Start: 2}
		s.GridColumnEnd = GridPlacement{Start: 4}
		s.GridRowStart = GridPlacement{Start: 1}
		s.GridRowEnd = GridPlacement{Start: 2}
	})
	b := gridChild(func(s *Style) {
		s.GridColumnStart = GridPlacement{Start: 1}
		s.GridRowStart = GridPlacement{Start: 2}
	})

	root := NewNode(Style{
		Display: DisplayGrid,
		Width:   Pt(300),
		Height:  Pt(200),
		GridTemplateColumns: TrackList{Tracks: []TrackSize{
			{Kind: TrackFixed, Length: Pt(100)},
			{Kind: TrackFixed, Length: Pt(100)},
			{Kind: TrackFixed, Length: Pt(100)},
		}},
		GridTemplateRows: TrackList{Tracks: []TrackSize{
			{Kind: TrackFixed, Length: Pt(100)},
			{Kind: TrackFixed, Length: Pt(100)},
		}},
	}, a, b)

	Compute(root, 300, 200)

	assertLayout(t, a, "a", 100, 0, 200, 100)
	assertLayout(t, b, "b", 0, 100, 100, 100)
}

func TestGridSpan(t *testing.T) {
	a := gridChild(func(s *Style) {
		s.GridColumnStart = GridPlacement{Span: 2}
	})
	b := gridChild()

	root := NewNode(Style{
		Display: DisplayGrid,
		Width:   Pt(300),
		Height:  Pt(100),
		GridTemplateColumns: TrackList{Tracks: []TrackSize{
			{Kind: TrackFixed, Length: Pt(100)},
			{Kind: TrackFixed, Length: Pt(100)},
			{Kind: TrackFixed, Length: Pt(100)},
		}},
		GridTemplateRows: TrackList{Tracks: []TrackSize{
			{Kind: TrackFixed, Length: Pt(100)},
		}},
	}, a, b)

	Compute(root, 300, 100)

	assertLayout(t, a, "a", 0, 0, 200, 100)
	assertLayout(t, b, "b", 200, 0, 100, 100)
}

func TestGridAutoPlacementRowFlow(t *testing.T) {
	items := make([]*Node, 4)
	for i := range items {
		items[i] = gridChild()
	}

	root := NewNode(Style{
		Display: DisplayGrid,
		Width:   Pt(400),
		Height:  Pt(200),
		GridTemplateColumns: TrackList{Tracks: []TrackSize{
			{Kind: TrackFixed, Length: Pt(200)},
			{Kind: TrackFixed, Length: Pt(200)},
		}},
		GridTemplateRows: TrackList{Tracks: []TrackSize{
			{Kind: TrackFixed, Length: Pt(100)},
			{Kind: TrackFixed, Length: Pt(100)},
		}},
	}, items...)

	Compute(root, 400, 200)

	assertLayout(t, items[0], "0", 0, 0, 200, 100)
	assertLayout(t, items[1], "1", 200, 0, 200, 100)
	assertLayout(t, items[2], "2", 0, 100, 200, 100)
	assertLayout(t, items[3], "3", 200, 100, 200, 100)
}

func TestGridColSpanFull(t *testing.T) {
	a := gridChild(func(s *Style) {
		s.GridColumnStart = GridPlacement{Start: 1}
		s.GridColumnEnd = GridPlacement{Start: -1}
	})
	root := NewNode(Style{
		Display: DisplayGrid,
		Width:   Pt(300),
		Height:  Pt(100),
		GridTemplateColumns: TrackList{Tracks: []TrackSize{
			{Kind: TrackFixed, Length: Pt(100)},
			{Kind: TrackFixed, Length: Pt(100)},
			{Kind: TrackFixed, Length: Pt(100)},
		}},
		GridTemplateRows: TrackList{Tracks: []TrackSize{{Kind: TrackFixed, Length: Pt(100)}}},
	}, a)
	Compute(root, 300, 100)
	assertLayout(t, a, "col-span-full", 0, 0, 300, 100)
}

func TestGridRowOnlyPlacement(t *testing.T) {
	a := gridChild(func(s *Style) { s.GridRowStart = GridPlacement{Start: 2} })
	b := gridChild()
	root := NewNode(Style{
		Display: DisplayGrid,
		Width:   Pt(200),
		Height:  Pt(200),
		GridTemplateColumns: TrackList{Tracks: []TrackSize{
			{Kind: TrackFixed, Length: Pt(100)},
			{Kind: TrackFixed, Length: Pt(100)},
		}},
		GridTemplateRows: TrackList{Tracks: []TrackSize{
			{Kind: TrackFixed, Length: Pt(100)},
			{Kind: TrackFixed, Length: Pt(100)},
		}},
	}, a, b)
	Compute(root, 200, 200)
	assertLayout(t, a, "row-only", 0, 100, 100, 100)
	assertLayout(t, b, "auto", 0, 0, 100, 100)
}

func TestGridNegativeLineEnd(t *testing.T) {
	a := gridChild(func(s *Style) {
		s.GridColumnStart = GridPlacement{Start: 2}
		s.GridColumnEnd = GridPlacement{Start: -1}
	})
	root := NewNode(Style{
		Display: DisplayGrid,
		Width:   Pt(400),
		Height:  Pt(100),
		GridTemplateColumns: TrackList{Tracks: []TrackSize{
			{Kind: TrackFixed, Length: Pt(100)},
			{Kind: TrackFixed, Length: Pt(100)},
			{Kind: TrackFixed, Length: Pt(100)},
			{Kind: TrackFixed, Length: Pt(100)},
		}},
		GridTemplateRows: TrackList{Tracks: []TrackSize{{Kind: TrackFixed, Length: Pt(100)}}},
	}, a)
	Compute(root, 400, 100)
	assertLayout(t, a, "2 / -1", 100, 0, 300, 100)
}

func TestGridEndLessThanStartSwaps(t *testing.T) {
	a := gridChild(func(s *Style) {
		s.GridColumnStart = GridPlacement{Start: 3}
		s.GridColumnEnd = GridPlacement{Start: 1}
	})
	root := NewNode(Style{
		Display: DisplayGrid,
		Width:   Pt(300),
		Height:  Pt(100),
		GridTemplateColumns: TrackList{Tracks: []TrackSize{
			{Kind: TrackFixed, Length: Pt(100)},
			{Kind: TrackFixed, Length: Pt(100)},
			{Kind: TrackFixed, Length: Pt(100)},
		}},
		GridTemplateRows: TrackList{Tracks: []TrackSize{{Kind: TrackFixed, Length: Pt(100)}}},
	}, a)
	Compute(root, 300, 100)
	assertLayout(t, a, "3 / 1 → 1..3", 0, 0, 200, 100)
}

func TestGridUnboundedSpanIsCapped(t *testing.T) {
	a := gridChild(func(s *Style) {
		s.GridColumnStart = GridPlacement{Start: 1}
		s.GridColumnEnd = GridPlacement{Start: 100000}
	})
	root := NewNode(Style{
		Display: DisplayGrid,
		Width:   Pt(300),
		Height:  Pt(100),
		GridTemplateColumns: TrackList{Tracks: []TrackSize{
			{Kind: TrackFixed, Length: Pt(100)},
			{Kind: TrackFixed, Length: Pt(100)},
			{Kind: TrackFixed, Length: Pt(100)},
		}},
		GridTemplateRows: TrackList{Tracks: []TrackSize{{Kind: TrackFixed, Length: Pt(100)}}},
	}, a)
	Compute(root, 300, 100)
	if a.Layout.Y > 200 {
		t.Fatalf("unbounded span caused runaway placement: y=%v", a.Layout.Y)
	}
}

func TestGridPercentTracks(t *testing.T) {
	a := gridChild()
	b := gridChild()

	root := NewNode(Style{
		Display: DisplayGrid,
		Width:   Pt(400),
		Height:  Pt(100),
		GridTemplateColumns: TrackList{Tracks: []TrackSize{
			{Kind: TrackFixed, Length: Pct(25)},
			{Kind: TrackFixed, Length: Pct(75)},
		}},
		GridTemplateRows: TrackList{Tracks: []TrackSize{
			{Kind: TrackFixed, Length: Pt(100)},
		}},
	}, a, b)

	Compute(root, 400, 100)

	assertLayout(t, a, "a", 0, 0, 100, 100)
	assertLayout(t, b, "b", 100, 0, 300, 100)
}
