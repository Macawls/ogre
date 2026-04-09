package layout

import (
	"github.com/macawls/ogre/parse"
	"github.com/macawls/ogre/style"
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
	for _, child := range pn.Children {
		cn := buildNode(child, styles, measureText, nodeMap)
		if cn != nil {
			children = append(children, cn)
		}
	}

	ln := NewNode(s, children...)
	nodeMap[pn] = ln
	return ln
}

func mapStyle(cs *style.ComputedStyle) Style {
	return Style{
		Display:  mapDisplay(cs.Display),
		Position: mapPosition(cs.Position),

		Direction: mapDirection(cs.FlexDirection),
		Wrap:      mapWrap(cs.FlexWrap),

		AlignItems:     mapAlignItems(cs.AlignItems),
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
	}
}

func mapDisplay(d style.Display) Display {
	switch d {
	case style.DisplayNone:
		return DisplayNone
	default:
		return DisplayFlex
	}
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
