// Package layout implements a flexbox layout engine for positioned node trees.
package layout

import "math"

// Direction specifies the main axis of a flex container.
type Direction int

const (
	Row Direction = iota
	RowReverse
	Column
	ColumnReverse
)

// Wrap specifies whether flex items wrap onto multiple lines.
type Wrap int

const (
	NoWrap Wrap = iota
	WrapWrap
	WrapReverse
)

// Justify specifies how flex items are distributed along the main axis.
type Justify int

const (
	JustifyStart Justify = iota
	JustifyEnd
	JustifyCenter
	JustifySpaceBetween
	JustifySpaceAround
	JustifySpaceEvenly
)

// Align specifies how flex items are aligned along the cross axis.
type Align int

const (
	AlignAuto Align = iota
	AlignStart
	AlignEnd
	AlignCenter
	AlignStretch
	AlignBaseline
)

// Position specifies whether a node uses relative or absolute positioning.
type Position int

const (
	// Relative positions the node relative to its normal flow position.
	Relative Position = iota
	Absolute
)

// DimensionUnit identifies the unit type of a layout dimension.
type DimensionUnit int

const (
	UnitUndefined DimensionUnit = iota
	UnitPoint
	UnitPercent
	UnitAuto
)

// Dimension is a value with a unit, used for widths, heights, and margins.
type Dimension struct {
	Value float64
	Unit  DimensionUnit
}

// Pt creates a Dimension with a point value.
func Pt(v float64) Dimension { return Dimension{Value: v, Unit: UnitPoint} }
// Pct creates a Dimension with a percentage value.
func Pct(v float64) Dimension { return Dimension{Value: v, Unit: UnitPercent} }
// Auto creates a Dimension representing the CSS auto keyword.
func Auto() Dimension { return Dimension{Unit: UnitAuto} }
// Undefined creates a Dimension with no defined value.
func Undefined() Dimension { return Dimension{Unit: UnitUndefined} }

func (d Dimension) IsAuto() bool      { return d.Unit == UnitAuto }
func (d Dimension) IsUndefined() bool  { return d.Unit == UnitUndefined }
func (d Dimension) IsDefined() bool    { return d.Unit == UnitPoint || d.Unit == UnitPercent }

func (d Dimension) Resolve(base float64) float64 {
	switch d.Unit {
	case UnitPoint:
		return d.Value
	case UnitPercent:
		if math.IsInf(base, 0) || math.IsNaN(base) {
			return math.NaN()
		}
		return d.Value / 100 * base
	default:
		return math.NaN()
	}
}

// Style holds the layout-relevant CSS properties for a node.
type Style struct {
	Display  Display
	Position Position

	Direction Direction
	Wrap      Wrap

	AlignItems     Align
	AlignSelf      Align
	AlignContent   Align
	JustifyContent Justify

	FlexGrow   float64
	FlexShrink float64
	FlexBasis  Dimension

	Width     Dimension
	Height    Dimension
	MinWidth  Dimension
	MinHeight Dimension
	MaxWidth  Dimension
	MaxHeight Dimension

	MarginTop    Dimension
	MarginRight  Dimension
	MarginBottom Dimension
	MarginLeft   Dimension

	PaddingTop    Dimension
	PaddingRight  Dimension
	PaddingBottom Dimension
	PaddingLeft   Dimension

	BorderTop    float64
	BorderRight  float64
	BorderBottom float64
	BorderLeft   float64

	Gap       float64
	RowGap    float64
	ColumnGap float64

	Top    Dimension
	Right  Dimension
	Bottom Dimension
	Left   Dimension

	AspectRatio float64
}

// Display controls whether a node participates in layout.
type Display int

const (
	DisplayFlex Display = iota
	DisplayNone
)

// Layout holds the computed position and size of a node after layout.
type Layout struct {
	X, Y          float64
	Width, Height float64
	Padding       [4]float64
	Border        [4]float64
}

// MeasureFunc measures a leaf node's intrinsic size given available dimensions.
type MeasureFunc func(maxWidth, maxHeight float64) (float64, float64)

// Node is a layout tree node with a style, children, and computed layout.
type Node struct {
	Style    Style
	Children []*Node
	Measure  MeasureFunc
	Layout   Layout
	parent   *Node
}

// NewNode creates a layout Node with the given style and children.
func NewNode(style Style, children ...*Node) *Node {
	n := &Node{Style: style, Children: children}
	for _, c := range children {
		c.parent = n
	}
	return n
}

// NewLeaf creates a leaf layout Node that uses the given measure function for sizing.
func NewLeaf(style Style, measure MeasureFunc) *Node {
	return &Node{Style: style, Measure: measure}
}

// Compute runs the flexbox layout algorithm on the tree rooted at root.
func Compute(root *Node, availableWidth, availableHeight float64) {
	computeNode(root, availableWidth, availableHeight)
}
