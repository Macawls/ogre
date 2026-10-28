package layout

import "math"

// computeTable lays out a node whose display is table: rows stacked vertically,
// cells in equal or preferred-width columns, every row as tall as its tallest
// cell. It is the minimum needed to reproduce what a browser does with the
// markup people actually write - <table><tr><td> with padding, a width on the
// table and optionally a width on a cell - and deliberately does not implement
// table-layout:fixed, border-collapse, rowspan or colspan.
func computeTable(node *Node, in layoutRun, contentWidth, contentHeight float64) {
	rows := collectTableRows(node)
	if len(rows) == 0 {
		return
	}

	// Cells per row, skipping none-displayed ones. Column count is the widest row.
	rowCells := make([][]*Node, len(rows))
	numCols := 0
	for i, row := range rows {
		for _, child := range row.Children {
			if child.Style.Display == DisplayNone {
				continue
			}
			rowCells[i] = append(rowCells[i], child)
		}
		if len(rowCells[i]) > numCols {
			numCols = len(rowCells[i])
		}
	}
	if numCols == 0 {
		return
	}

	colWidths := tableColumnWidths(rowCells, numCols, contentWidth)
	rowHeights := tableRowHeights(rowCells, colWidths, contentWidth, contentHeight)

	bpLeft := node.Layout.Padding[3] + node.Layout.Border[3]
	bpTop := node.Layout.Padding[0] + node.Layout.Border[0]

	y := bpTop
	for i, row := range rows {
		x := bpLeft
		for c, cell := range rowCells[i] {
			w := colWidths[c]
			h := rowHeights[i]
			cellIn := layoutRun{
				availableWidth:  w,
				availableHeight: h,
				ownerWidth:      w,
				ownerHeight:     h,
				definiteWidth:   true,
				definiteHeight:  true,
			}
			computeNode(cell, cellIn)
			// Cell offsets are relative to the row, which already carries its own
			// position: setting the absolute y here would add the row's offset twice
			// once the renderer walks down the tree.
			cell.Layout.X = x - bpLeft
			cell.Layout.Y = 0
			cell.Layout.Width = w
			cell.Layout.Height = h
			x += w
		}
		row.Layout.X = bpLeft
		row.Layout.Y = y
		row.Layout.Width = sumFloats(colWidths)
		row.Layout.Height = rowHeights[i]
		y += rowHeights[i]
	}

	for _, row := range rows {
		for _, child := range row.Children {
			if child.Style.Position == Absolute {
				computeAbsolute(child, row)
			}
		}
	}

	// A table never stretches to a definite container height: its height is the sum
	// of its rows, exactly as in CSS. The flex pass upstream may have handed down a
	// definite height; ignore it and correct the node.
	_ = in
	if !node.Style.Height.IsDefined() {
		bpH := node.Layout.Padding[0] + node.Layout.Padding[2] + node.Layout.Border[0] + node.Layout.Border[2]
		node.Layout.Height = sumFloats(rowHeights) + bpH
	}
}

// collectTableRows returns the table's rows in document order. Rows may be direct
// children of the table or wrapped in a row group (thead/tbody/tfoot), which is
// transparent for layout purposes.
func collectTableRows(table *Node) []*Node {
	var rows []*Node
	var walk func(n *Node)
	walk = func(n *Node) {
		for _, child := range n.Children {
			switch child.Style.Display {
			case DisplayNone:
				continue
			case DisplayTableRow:
				rows = append(rows, child)
			case DisplayTableRowGroup:
				walk(child)
			}
		}
	}
	walk(table)
	return rows
}

// tableColumnWidths resolves one width per column. A cell's explicit width is
// treated as that column's preferred width; the remaining space is shared equally
// between the auto columns, and if the preferred widths do not fit the table the
// columns are scaled down proportionally, which is what a browser's auto table
// layout converges to for this markup.
func tableColumnWidths(rowCells [][]*Node, numCols int, contentWidth float64) []float64 {
	widths := make([]float64, numCols)
	preferred := make([]float64, numCols)
	hasPreferred := make([]bool, numCols)

	for _, cells := range rowCells {
		for c, cell := range cells {
			if c >= numCols {
				continue
			}
			if cell.Style.Width.IsDefined() {
				if w := resolveNodeSize(cell.Style.Width, contentWidth); !math.IsNaN(w) && w > preferred[c] {
					preferred[c] = w
					hasPreferred[c] = true
				}
			}
		}
	}

	// Columns without an explicit width take their content's intrinsic width as a
	// preference, so a table narrower than its content still gives wider columns
	// more room than narrow ones.
	for c := 0; c < numCols; c++ {
		if hasPreferred[c] {
			continue
		}
		intrinsic := 0.0
		for _, cells := range rowCells {
			if c >= len(cells) {
				continue
			}
			// estimateSize returns the cell's max-content width: text is measured
			// unconstrained, block children stack at the widest of them. Laying the
			// cell out with an infinite available width instead would make a
			// stretch-aligned cell adopt that infinity as its own width.
			if w := estimateSize(cells[c], true, contentWidth, 0); !math.IsNaN(w) && w > intrinsic {
				intrinsic = w
			}
		}
		preferred[c] = intrinsic
	}

	total := sumFloats(preferred)
	switch {
	case total <= 0:
		for c := range widths {
			widths[c] = contentWidth / float64(numCols)
		}
	case total > contentWidth:
		for c := range widths {
			widths[c] = preferred[c] * contentWidth / total
		}
	default:
		// Everything fits: keep preferred widths and share the slack equally.
		slack := (contentWidth - total) / float64(numCols)
		for c := range widths {
			widths[c] = preferred[c] + slack
		}
	}
	return widths
}

// tableRowHeights resolves one height per row: the tallest cell in it, measured
// at that column's resolved width.
func tableRowHeights(rowCells [][]*Node, colWidths []float64, contentWidth, contentHeight float64) []float64 {
	heights := make([]float64, len(rowCells))
	for i, cells := range rowCells {
		for c, cell := range cells {
			if c >= len(colWidths) {
				continue
			}
			w := colWidths[c]
			cellIn := layoutRun{
				availableWidth:  w,
				availableHeight: math.Inf(1),
				ownerWidth:      w,
				ownerHeight:     contentHeight,
				definiteWidth:   true,
				definiteHeight:  false,
			}
			computeNode(cell, cellIn)
			if cell.Layout.Height > heights[i] {
				heights[i] = cell.Layout.Height
			}
		}
	}
	return heights
}

func sumFloats(v []float64) float64 {
	total := 0.0
	for _, x := range v {
		total += x
	}
	return total
}
