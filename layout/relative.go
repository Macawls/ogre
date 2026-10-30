package layout

import "math"

// applyRelativeOffsets translates every position:relative box by its
// left/right/top/bottom offsets.
//
// It runs after the whole tree has been laid out, because a relative offset is
// a paint-time shift, not a layout change: it must not move the box's siblings,
// and it must not change the space the parent reserved for it. Layout
// coordinates are parent-relative, so shifting a box's own origin carries its
// entire subtree with it — nested children and absolute descendants alike,
// which are placed against this box.
func applyRelativeOffsets(node *Node) {
	if node.Style.Display == DisplayNone {
		return
	}

	cl := &node.Layout
	// The containing block for a percentage offset is this node's padding box,
	// the same box absolute children resolve against.
	paddingBoxW := cl.Width - cl.Border[1] - cl.Border[3]
	paddingBoxH := cl.Height - cl.Border[0] - cl.Border[2]

	for _, child := range node.Children {
		if child.Style.Position == PositionRelative {
			child.Layout.X += relativeOffset(child.Style.Left, child.Style.Right, paddingBoxW)
			child.Layout.Y += relativeOffset(child.Style.Top, child.Style.Bottom, paddingBoxH)
		}
		applyRelativeOffsets(child)
	}
}

// relativeOffset resolves one axis of a relative offset. The start edge (left
// or top) wins when both are given; the end edge (right or bottom) pulls the
// box back. A percentage resolves against the containing block, and an auto or
// undefined value means no offset at all.
func relativeOffset(start, end Dimension, base float64) float64 {
	if v := start.Resolve(base); !math.IsNaN(v) {
		return v
	}
	if v := end.Resolve(base); !math.IsNaN(v) {
		return -v
	}
	return 0
}
