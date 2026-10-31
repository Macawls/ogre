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

	// An auto size is shrink-to-fit against the space the declared offsets
	// leave over, not against the whole padding box.
	availW := absoluteAvailable(paddingBoxW, left, right)
	availH := absoluteAvailable(paddingBoxH, top, bottom)

	// Measure before positioning. A measured leaf reports its content size plus
	// border and padding, which is its shrink-to-fit size; anything else is
	// estimated from its children the same way flex estimates a base size.
	if child.Measure != nil && (math.IsNaN(w) || math.IsNaN(h)) {
		computeNode(child, layoutRun{
			availableWidth:  availW,
			availableHeight: availH,
			ownerWidth:      paddingBoxW,
			ownerHeight:     paddingBoxH,
			definiteWidth:   true,
			definiteHeight:  true,
		})
	}

	if math.IsNaN(w) {
		if child.Measure != nil {
			w = child.Layout.Width
		} else {
			w = estimateSize(child, true, paddingBoxW, paddingBoxH)
		}
		if w > availW {
			w = availW
		}
		w = clampSize(w, s.MinWidth, s.MaxWidth, paddingBoxW)
	}

	if math.IsNaN(h) {
		if child.Measure != nil {
			h = child.Layout.Height
		} else {
			h = estimateSize(child, false, paddingBoxW, paddingBoxH)
		}
		if h > availH {
			h = availH
		}
		h = clampSize(h, s.MinHeight, s.MaxHeight, paddingBoxH)
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

	computeNode(child, layoutRun{
		availableWidth:  w,
		availableHeight: h,
		ownerWidth:      paddingBoxW,
		ownerHeight:     paddingBoxH,
		definiteWidth:   true,
		definiteHeight:  true,
	})

	// computeNode measures a leaf and overwrites Layout.Width/Height with the
	// content size, which is narrower and shorter than a size resolved from
	// both offsets or from min/max. For an absolutely positioned box the
	// resolved size is the one to keep.
	child.Layout.Width = w
	child.Layout.Height = h
}

// absoluteAvailable returns the space left on one axis once the declared
// offsets are subtracted from the containing block's padding box.
func absoluteAvailable(paddingBox, start, end float64) float64 {
	avail := paddingBox
	if !math.IsNaN(start) {
		avail -= start
	}
	if !math.IsNaN(end) {
		avail -= end
	}
	if avail < 0 {
		avail = 0
	}
	return avail
}
