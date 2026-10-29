package layout

import (
	"github.com/macawls/ogre/v3/parse"
	"github.com/macawls/ogre/v3/style"
)

type LayoutTree struct {
	Root    *Node
	NodeMap map[*parse.Node]*Node
}

type MeasureTextFunc func(pn *parse.Node, text string, cs *style.ComputedStyle, maxWidth float64) (float64, float64)

func BuildTree(root *parse.Node, styles map[*parse.Node]*style.ComputedStyle, measureText MeasureTextFunc) *LayoutTree {
	tree := &LayoutTree{
		NodeMap: make(map[*parse.Node]*Node),
	}
	tree.Root = buildNode(root, styles, measureText, tree.NodeMap)
	return tree
}

func buildNode(pn *parse.Node, styles map[*parse.Node]*style.ComputedStyle, measureText MeasureTextFunc, nodeMap map[*parse.Node]*Node) *Node {
	cs := styles[pn]
	if cs == nil {
		cs = &style.ComputedStyle{}
	}

	if cs.Display == style.DisplayNone {
		return nil
	}

	s := mapStyle(cs)

	if pn.Type == parse.TextNode {
		ln := NewLeaf(s, func(maxWidth, maxHeight float64) (float64, float64) {
			return measureText(pn, pn.Text, cs, maxWidth)
		})
		nodeMap[pn] = ln
		return ln
	}

	var children []*Node
	for i, child := range pn.Children {
		var cn *Node
		if child.Type == parse.TextNode && child.LineBreak {
			// A <br> ends the current line, and the line box itself comes from
			// the text that follows, so the break normally contributes no
			// height. It does open one when nothing follows it on that line: a
			// trailing break, or one immediately followed by another break,
			// leaves an empty line, exactly as the browser does.
			cn = buildBreakNode(child, styles, measureText, nodeMap, opensLineBox(pn.Children, i, styles))
		} else {
			cn = buildNode(child, styles, measureText, nodeMap)
		}
		if cn != nil {
			children = append(children, cn)
		}
	}

	ln := NewNode(s, children...)
	nodeMap[pn] = ln
	return ln
}

// buildBreakNode lays out a <br>. It occupies zero width and, unless it opens a
// line box of its own, zero height: the next text leaf supplies the line box it
// starts.
func buildBreakNode(pn *parse.Node, styles map[*parse.Node]*style.ComputedStyle, measureText MeasureTextFunc, nodeMap map[*parse.Node]*Node, opensLine bool) *Node {
	cs := styles[pn]
	if cs == nil {
		cs = &style.ComputedStyle{}
	}
	if cs.Display == style.DisplayNone {
		return nil
	}

	s := mapStyle(cs)
	ln := NewLeaf(s, func(maxWidth, maxHeight float64) (float64, float64) {
		if !opensLine {
			return 0, 0
		}
		_, h := measureText(pn, " ", cs, maxWidth)
		return 0, h
	})
	nodeMap[pn] = ln
	return ln
}

// opensLineBox reports whether the <br> at index i opens a line box rather than
// only ending the current one. That is the case exactly when the next rendered
// sibling is itself a break: the run's last break is followed by the text that
// supplies the line it starts, and a trailing break opens nothing at all (the
// browser drops the empty line box at the end of a block). So a run of N breaks
// advances N-1 empty lines plus the line of whatever follows.
func opensLineBox(children []*parse.Node, i int, styles map[*parse.Node]*style.ComputedStyle) bool {
	for _, sib := range children[i+1:] {
		if cs := styles[sib]; cs != nil && cs.Display == style.DisplayNone {
			continue
		}
		return sib.Type == parse.TextNode && sib.LineBreak
	}
	return false
}

func mapStyle(cs *style.ComputedStyle) Style {
	direction := mapDirection(cs.FlexDirection)
	alignItems := mapAlignItems(cs.AlignItems)
	if cs.Display == style.DisplayBlock {
		direction = Column
		if cs.AlignItems == style.AlignItemsAuto {
			alignItems = AlignStretch
		}
	}
	return Style{
		Display:  mapDisplay(cs.Display),
		Position: mapPosition(cs.Position),

		Direction: direction,
		Wrap:      mapWrap(cs.FlexWrap),

		AlignItems:     alignItems,
		AlignSelf:      mapAlignSelf(cs.AlignSelf),
		AlignContent:   mapAlignContent(cs.AlignContent),
		JustifyContent: mapJustify(cs.JustifyContent),

		FlexGrow:   cs.FlexGrow,
		FlexShrink: cs.FlexShrink,
		FlexBasis:  mapDimension(cs.FlexBasis),

		Width:     mapDimension(cs.Width),
		Height:    mapDimension(cs.Height),
		MinWidth:  mapDimension(cs.MinWidth),
		MinHeight: mapDimension(cs.MinHeight),
		MaxWidth:  mapDimension(cs.MaxWidth),
		MaxHeight: mapDimension(cs.MaxHeight),

		MarginTop:    mapDimension(cs.MarginTop),
		MarginRight:  mapDimension(cs.MarginRight),
		MarginBottom: mapDimension(cs.MarginBottom),
		MarginLeft:   mapDimension(cs.MarginLeft),

		PaddingTop:    Pt(cs.PaddingTop),
		PaddingRight:  Pt(cs.PaddingRight),
		PaddingBottom: Pt(cs.PaddingBottom),
		PaddingLeft:   Pt(cs.PaddingLeft),

		BorderTop:    cs.BorderTopWidth,
		BorderRight:  cs.BorderRightWidth,
		BorderBottom: cs.BorderBottomWidth,
		BorderLeft:   cs.BorderLeftWidth,

		Gap:       cs.Gap,
		RowGap:    cs.RowGap,
		ColumnGap: cs.ColumnGap,

		Top:    mapDimension(cs.Top),
		Right:  mapDimension(cs.Right),
		Bottom: mapDimension(cs.Bottom),
		Left:   mapDimension(cs.Left),

		AspectRatio: cs.AspectRatio,

		GridTemplateColumns: mapTrackList(cs.GridTemplateColumns),
		GridTemplateRows:    mapTrackList(cs.GridTemplateRows),
		GridAutoColumns:     mapTrackSize(cs.GridAutoColumns),
		GridAutoRows:        mapTrackSize(cs.GridAutoRows),
		GridAutoFlow:        mapGridAutoFlow(cs.GridAutoFlow),

		GridColumnStart: mapGridPlacement(cs.GridColumnStart),
		GridColumnEnd:   mapGridPlacement(cs.GridColumnEnd),
		GridRowStart:    mapGridPlacement(cs.GridRowStart),
		GridRowEnd:      mapGridPlacement(cs.GridRowEnd),
	}
}

func mapDisplay(d style.Display) Display {
	switch d {
	case style.DisplayNone:
		return DisplayNone
	case style.DisplayGrid:
		return DisplayGrid
	default:
		return DisplayFlex
	}
}

func mapTrackSize(t style.TrackSize) TrackSize {
	kind := TrackAuto
	switch t.Kind {
	case style.TrackFixed:
		kind = TrackFixed
	case style.TrackFr:
		kind = TrackFr
	}
	length := Undefined()
	switch t.Length.Unit {
	case style.UnitPx:
		length = Pt(t.Length.Raw)
	case style.UnitPercent:
		length = Pct(t.Length.Raw)
	case style.UnitAuto:
		length = Auto()
	}
	return TrackSize{Kind: kind, Length: length, Fr: t.Fr}
}

func mapTrackList(l style.TrackList) TrackList {
	if len(l.Tracks) == 0 {
		return TrackList{}
	}
	out := make([]TrackSize, len(l.Tracks))
	for i, t := range l.Tracks {
		out[i] = mapTrackSize(t)
	}
	return TrackList{Tracks: out}
}

func mapGridAutoFlow(f style.GridAutoFlow) GridAutoFlow {
	switch f {
	case style.FlowColumn:
		return FlowColumn
	case style.FlowRowDense:
		return FlowRowDense
	case style.FlowColumnDense:
		return FlowColumnDense
	}
	return FlowRow
}

func mapGridPlacement(p style.GridPlacement) GridPlacement {
	return GridPlacement{Start: p.Start, End: p.End, Span: p.Span}
}

func mapPosition(p style.Position) Position {
	switch p {
	case style.PositionAbsolute:
		return Absolute
	default:
		return Relative
	}
}

func mapDirection(d style.FlexDirection) Direction {
	switch d {
	case style.FlexDirectionRow:
		return Row
	case style.FlexDirectionRowReverse:
		return RowReverse
	case style.FlexDirectionColumn:
		return Column
	case style.FlexDirectionColumnReverse:
		return ColumnReverse
	default:
		return Row
	}
}

func mapWrap(w style.FlexWrap) Wrap {
	switch w {
	case style.FlexWrapWrap:
		return WrapWrap
	case style.FlexWrapWrapReverse:
		return WrapReverse
	default:
		return NoWrap
	}
}

func mapJustify(j style.JustifyContent) Justify {
	switch j {
	case style.JustifyContentFlexStart:
		return JustifyStart
	case style.JustifyContentFlexEnd:
		return JustifyEnd
	case style.JustifyContentCenter:
		return JustifyCenter
	case style.JustifyContentSpaceBetween:
		return JustifySpaceBetween
	case style.JustifyContentSpaceAround:
		return JustifySpaceAround
	case style.JustifyContentSpaceEvenly:
		return JustifySpaceEvenly
	default:
		return JustifyStart
	}
}

func mapAlignItems(a style.AlignItems) Align {
	switch a {
	case style.AlignItemsAuto:
		return AlignAuto
	case style.AlignItemsFlexStart:
		return AlignStart
	case style.AlignItemsFlexEnd:
		return AlignEnd
	case style.AlignItemsCenter:
		return AlignCenter
	case style.AlignItemsStretch:
		return AlignStretch
	case style.AlignItemsBaseline:
		return AlignBaseline
	default:
		return AlignAuto
	}
}

func mapAlignSelf(a style.AlignSelf) Align {
	switch a {
	case style.AlignSelfAuto:
		return AlignAuto
	case style.AlignSelfFlexStart:
		return AlignStart
	case style.AlignSelfFlexEnd:
		return AlignEnd
	case style.AlignSelfCenter:
		return AlignCenter
	case style.AlignSelfStretch:
		return AlignStretch
	case style.AlignSelfBaseline:
		return AlignBaseline
	default:
		return AlignAuto
	}
}

func mapAlignContent(a style.AlignContent) Align {
	switch a {
	case style.AlignContentAuto:
		return AlignAuto
	case style.AlignContentFlexStart:
		return AlignStart
	case style.AlignContentFlexEnd:
		return AlignEnd
	case style.AlignContentCenter:
		return AlignCenter
	case style.AlignContentStretch:
		return AlignStretch
	case style.AlignContentBaseline:
		return AlignBaseline
	default:
		return AlignAuto
	}
}

func mapDimension(v style.Value) Dimension {
	switch v.Unit {
	case style.UnitPx:
		return Pt(v.Raw)
	case style.UnitPercent:
		return Pct(v.Raw)
	case style.UnitAuto:
		return Auto()
	case style.UnitNone:
		return Undefined()
	default:
		return Undefined()
	}
}

func ComputeLayout(root *parse.Node, styles map[*parse.Node]*style.ComputedStyle, width, height float64, measureText MeasureTextFunc) *LayoutTree {
	tree := BuildTree(root, styles, measureText)
	if tree.Root != nil {
		Compute(tree.Root, width, height)
	}
	return tree
}
