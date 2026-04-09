package layout

import "math"

func computeAbsolute(child *Node, container *Node) {
	s := &child.Style
	cl := &container.Layout

	paddingBoxW := cl.Width - cl.Border[1] - cl.Border[3]
	paddingBoxH := cl.Height - cl.Border[0] - cl.Border[2]

	w := resolveNodeSize(s.Width, paddingBoxW)
	h := resolveNodeSize(s.Height, paddingBoxH)

	w = clampSize(w, s.MinWidth, s.MaxWidth, paddingBoxW)
	h = clampSize(h, s.MinHeight, s.MaxHeight, paddingBoxH)

	left := resolveNodeSize(s.Left, paddingBoxW)
	right := resolveNodeSize(s.Right, paddingBoxW)
	top := resolveNodeSize(s.Top, paddingBoxH)
	bottom := resolveNodeSize(s.Bottom, paddingBoxH)

	if math.IsNaN(w) && !math.IsNaN(left) && !math.IsNaN(right) {
		w = paddingBoxW - left - right
		if w < 0 {
			w = 0
		}
		w = clampSize(w, s.MinWidth, s.MaxWidth, paddingBoxW)
	}

	if math.IsNaN(h) && !math.IsNaN(top) && !math.IsNaN(bottom) {
		h = paddingBoxH - top - bottom
		if h < 0 {
			h = 0
		}
		h = clampSize(h, s.MinHeight, s.MaxHeight, paddingBoxH)
	}

	if s.AspectRatio > 0 {
		if !math.IsNaN(w) && math.IsNaN(h) {
			h = w / s.AspectRatio
		} else if !math.IsNaN(h) && math.IsNaN(w) {
			w = h * s.AspectRatio
		}
	}

	if math.IsNaN(w) {
		w = 0
	}
	if math.IsNaN(h) {
		h = 0
	}

	var x, y float64

	borderLeft := cl.Border[3]
	borderTop := cl.Border[0]

	if !math.IsNaN(left) {
		x = borderLeft + left
	} else if !math.IsNaN(right) {
		x = borderLeft + paddingBoxW - right - w
	} else {
		x = borderLeft
	}

	if !math.IsNaN(top) {
		y = borderTop + top
	} else if !math.IsNaN(bottom) {
		y = borderTop + paddingBoxH - bottom - h
	} else {
		y = borderTop
	}

	child.Layout.X = x
	child.Layout.Y = y
	child.Layout.Width = w
	child.Layout.Height = h

	computeNode(child, w, h)
}
