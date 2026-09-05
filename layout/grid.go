package layout

import "math"

func computeGrid(node *Node, in layoutRun, contentWidth, contentHeight float64) {
	s := &node.Style

	var items, absItems []*Node
	for _, child := range node.Children {
		if child.Style.Display == DisplayNone {
			continue
		}
		if child.Style.Position == Absolute {
			absItems = append(absItems, child)
			continue
		}
		items = append(items, child)
	}

	colGap := s.ColumnGap
	if colGap == 0 {
		colGap = s.Gap
	}
	rowGap := s.RowGap
	if rowGap == 0 {
		rowGap = s.Gap
	}

	colTemplate := s.GridTemplateColumns.Tracks
	rowTemplate := s.GridTemplateRows.Tracks

	placements := placeGridItems(items, colTemplate, rowTemplate, s.GridAutoFlow)

	numCols := len(colTemplate)
	numRows := len(rowTemplate)
	for _, p := range placements {
		if p.colEnd-1 > numCols {
			numCols = p.colEnd - 1
		}
		if p.rowEnd-1 > numRows {
			numRows = p.rowEnd - 1
		}
	}
	if numCols == 0 {
		numCols = 1
	}
	if numRows == 0 {
		numRows = 1
	}

	colSizes := sizeTracks(colTemplate, numCols, s.GridAutoColumns, contentWidth, colGap)
	rowSizes := sizeTracks(rowTemplate, numRows, s.GridAutoRows, contentHeight, rowGap)

	needAutoRows := hasAutoRows(rowTemplate, s.GridAutoRows, numRows)
	if needAutoRows && !in.definiteHeight {
		measureAutoRows(items, placements, colSizes, colGap, rowSizes, rowTemplate, s.GridAutoRows, contentHeight)
	}

	colOffsets := trackOffsets(colSizes, colGap)
	rowOffsets := trackOffsets(rowSizes, rowGap)

	bpLeft := node.Layout.Padding[3] + node.Layout.Border[3]
	bpTop := node.Layout.Padding[0] + node.Layout.Border[0]

	for i, child := range items {
		p := placements[i]
		x := bpLeft + colOffsets[p.colStart-1]
		w := trackSpanSize(colSizes, colGap, p.colStart-1, p.colEnd-1)
		y := bpTop + rowOffsets[p.rowStart-1]
		h := trackSpanSize(rowSizes, rowGap, p.rowStart-1, p.rowEnd-1)

		childIn := layoutRun{
			availableWidth:  w,
			availableHeight: h,
			ownerWidth:      w,
			ownerHeight:     h,
			definiteWidth:   true,
			definiteHeight:  true,
		}
		computeNode(child, childIn)
		child.Layout.X = x
		child.Layout.Y = y
		if !child.Style.Width.IsDefined() {
			child.Layout.Width = w
		}
		if !child.Style.Height.IsDefined() {
			child.Layout.Height = h
		}
	}

	if !node.Style.Height.IsDefined() && !in.definiteHeight {
		totalRows := 0.0
		for _, r := range rowSizes {
			totalRows += r
		}
		if len(rowSizes) > 1 {
			totalRows += rowGap * float64(len(rowSizes)-1)
		}
		bpH := node.Layout.Padding[0] + node.Layout.Padding[2] + node.Layout.Border[0] + node.Layout.Border[2]
		node.Layout.Height = totalRows + bpH
	}

	for _, child := range absItems {
		computeAbsolute(child, node)
	}
}

type gridSpot struct {
	colStart, colEnd int
	rowStart, rowEnd int
}

const maxImplicitTracks = 10_000

func placeGridItems(items []*Node, colTemplate, rowTemplate []TrackSize, flow GridAutoFlow) []gridSpot {
	explicitCols := len(colTemplate)
	explicitRows := len(rowTemplate)
	numCols := explicitCols
	if numCols == 0 {
		numCols = 1
	}
	numRows := explicitRows
	if numRows == 0 {
		numRows = 1
	}

	placements := make([]gridSpot, len(items))
	occupied := map[[2]int]bool{}

	type deferredItem struct {
		i                     int
		colHint, rowHint      int
		colSpan, rowSpan      int
		colFixed, rowFixed    bool
	}
	var majorDefinite []deferredItem
	var minorDefinite []deferredItem
	var fullyAuto []deferredItem

	rowMajor := flow != FlowColumn && flow != FlowColumnDense

	for i, child := range items {
		s := &child.Style
		colSpan := placementSpan(s.GridColumnStart, s.GridColumnEnd)
		rowSpan := placementSpan(s.GridRowStart, s.GridRowEnd)
		cs, ce, cOK := resolveLine(s.GridColumnStart, s.GridColumnEnd, explicitCols)
		rs, re, rOK := resolveLine(s.GridRowStart, s.GridRowEnd, explicitRows)
		if cOK && rOK {
			cs, ce = clampSpan(cs, ce)
			rs, re = clampSpan(rs, re)
			placements[i] = gridSpot{cs, ce, rs, re}
			if ce-1 > numCols {
				numCols = ce - 1
			}
			if re-1 > numRows {
				numRows = re - 1
			}
			markOccupied(occupied, cs, ce, rs, re)
			continue
		}
		d := deferredItem{i: i, colSpan: colSpan, rowSpan: rowSpan}
		if cOK {
			cs, ce = clampSpan(cs, ce)
			d.colHint, d.colSpan, d.colFixed = cs, ce-cs, true
		}
		if rOK {
			rs, re = clampSpan(rs, re)
			d.rowHint, d.rowSpan, d.rowFixed = rs, re-rs, true
		}
		switch {
		case rowMajor && d.rowFixed:
			majorDefinite = append(majorDefinite, d)
		case !rowMajor && d.colFixed:
			majorDefinite = append(majorDefinite, d)
		case rowMajor && d.colFixed:
			minorDefinite = append(minorDefinite, d)
		case !rowMajor && d.rowFixed:
			minorDefinite = append(minorDefinite, d)
		default:
			fullyAuto = append(fullyAuto, d)
		}
	}

	if rowMajor {
		for _, d := range majorDefinite {
			row := 1
			if d.rowFixed {
				row = d.rowHint
			}
			c := 1
			if d.colFixed {
				c = d.colHint
			}
			for {
				if row > maxImplicitTracks {
					break
				}
				if !d.colFixed && c+d.colSpan-1 > numCols {
					c = 1
					row++
					continue
				}
				if fits(occupied, c, c+d.colSpan, row, row+d.rowSpan) {
					placements[d.i] = gridSpot{c, c + d.colSpan, row, row + d.rowSpan}
					markOccupied(occupied, c, c+d.colSpan, row, row+d.rowSpan)
					if c+d.colSpan-1 > numCols {
						numCols = c + d.colSpan - 1
					}
					if row+d.rowSpan-1 > numRows {
						numRows = row + d.rowSpan - 1
					}
					break
				}
				if d.colFixed {
					row++
				} else {
					c++
				}
			}
		}
	} else {
		for _, d := range majorDefinite {
			col := 1
			if d.colFixed {
				col = d.colHint
			}
			r := 1
			if d.rowFixed {
				r = d.rowHint
			}
			for {
				if col > maxImplicitTracks {
					break
				}
				if !d.rowFixed && r+d.rowSpan-1 > numRows {
					r = 1
					col++
					continue
				}
				if fits(occupied, col, col+d.colSpan, r, r+d.rowSpan) {
					placements[d.i] = gridSpot{col, col + d.colSpan, r, r + d.rowSpan}
					markOccupied(occupied, col, col+d.colSpan, r, r+d.rowSpan)
					if col+d.colSpan-1 > numCols {
						numCols = col + d.colSpan - 1
					}
					if r+d.rowSpan-1 > numRows {
						numRows = r + d.rowSpan - 1
					}
					break
				}
				if d.rowFixed {
					col++
				} else {
					r++
				}
			}
		}
	}

	for _, d := range minorDefinite {
		if rowMajor {
			c := d.colHint
			r := 1
			for r <= maxImplicitTracks {
				if fits(occupied, c, c+d.colSpan, r, r+d.rowSpan) {
					placements[d.i] = gridSpot{c, c + d.colSpan, r, r + d.rowSpan}
					markOccupied(occupied, c, c+d.colSpan, r, r+d.rowSpan)
					if c+d.colSpan-1 > numCols {
						numCols = c + d.colSpan - 1
					}
					if r+d.rowSpan-1 > numRows {
						numRows = r + d.rowSpan - 1
					}
					break
				}
				r++
			}
		} else {
			r := d.rowHint
			c := 1
			for c <= maxImplicitTracks {
				if fits(occupied, c, c+d.colSpan, r, r+d.rowSpan) {
					placements[d.i] = gridSpot{c, c + d.colSpan, r, r + d.rowSpan}
					markOccupied(occupied, c, c+d.colSpan, r, r+d.rowSpan)
					if c+d.colSpan-1 > numCols {
						numCols = c + d.colSpan - 1
					}
					if r+d.rowSpan-1 > numRows {
						numRows = r + d.rowSpan - 1
					}
					break
				}
				c++
			}
		}
	}

	row, col := 1, 1
	step := func() {
		if rowMajor {
			col++
			if col > numCols {
				col = 1
				row++
			}
		} else {
			row++
			if row > numRows {
				row = 1
				col++
			}
		}
	}
	for _, d := range fullyAuto {
		colSpan, rowSpan := d.colSpan, d.rowSpan
		if colSpan > maxImplicitTracks {
			colSpan = maxImplicitTracks
		}
		if rowSpan > maxImplicitTracks {
			rowSpan = maxImplicitTracks
		}
		safety := 0
		for {
			if row > maxImplicitTracks || col > maxImplicitTracks {
				break
			}
			safety++
			if safety > maxImplicitTracks {
				break
			}
			if col+colSpan-1 > numCols {
				col = 1
				row++
			}
			if fits(occupied, col, col+colSpan, row, row+rowSpan) {
				break
			}
			step()
		}
		placements[d.i] = gridSpot{col, col + colSpan, row, row + rowSpan}
		markOccupied(occupied, col, col+colSpan, row, row+rowSpan)
		step()
	}
	return placements
}

func clampSpan(s, e int) (int, int) {
	if s < 1 {
		s = 1
	}
	if e <= s {
		e = s + 1
	}
	if e-s > maxImplicitTracks {
		e = s + maxImplicitTracks
	}
	return s, e
}

func markOccupied(occ map[[2]int]bool, cs, ce, rs, re int) {
	for r := rs; r < re; r++ {
		for c := cs; c < ce; c++ {
			occ[[2]int{c, r}] = true
		}
	}
}

func fits(occ map[[2]int]bool, cs, ce, rs, re int) bool {
	for r := rs; r < re; r++ {
		for c := cs; c < ce; c++ {
			if occ[[2]int{c, r}] {
				return false
			}
		}
	}
	return true
}

func resolveLine(start, end GridPlacement, explicitCount int) (int, int, bool) {
	s := resolveNegative(start.Start, explicitCount)
	e := resolveNegative(end.Start, explicitCount)
	span := start.Span
	if span == 0 {
		span = end.Span
	}
	switch {
	case s != 0 && e != 0:
		if e < s {
			s, e = e, s
		}
		if e == s {
			e = s + 1
		}
		return s, e, true
	case s != 0 && span > 0:
		return s, s + span, true
	case s != 0:
		return s, s + 1, true
	case e != 0 && span > 0:
		return e - span, e, true
	case e != 0:
		return e - 1, e, true
	}
	return 0, 0, false
}

func resolveNegative(line, explicitCount int) int {
	if line >= 0 {
		return line
	}
	if explicitCount <= 0 {
		return 0
	}
	resolved := explicitCount + 2 + line
	if resolved < 1 {
		return 0
	}
	return resolved
}

func placementSpan(start, end GridPlacement) int {
	if start.Span > 0 {
		return start.Span
	}
	if end.Span > 0 {
		return end.Span
	}
	if start.Start > 0 && end.Start > 0 && end.Start > start.Start {
		return end.Start - start.Start
	}
	return 1
}

func sizeTracks(template []TrackSize, count int, auto TrackSize, available, gap float64) []float64 {
	tracks := make([]TrackSize, count)
	for i := range count {
		if i < len(template) {
			tracks[i] = template[i]
		} else {
			tracks[i] = auto
			if tracks[i].Kind == TrackFixed && tracks[i].Length.IsUndefined() {
				tracks[i] = TrackSize{Kind: TrackAuto}
			}
		}
	}
	sizes := make([]float64, count)
	if count == 0 {
		return sizes
	}
	totalGap := 0.0
	if count > 1 {
		totalGap = gap * float64(count-1)
	}
	remaining := available - totalGap
	if remaining < 0 {
		remaining = 0
	}
	var totalFr float64
	for i, t := range tracks {
		switch t.Kind {
		case TrackFixed:
			v := t.Length.Resolve(available)
			if !math.IsNaN(v) {
				sizes[i] = v
				remaining -= v
			}
		case TrackFr:
			totalFr += t.Fr
		}
	}
	if remaining < 0 {
		remaining = 0
	}
	autoCount := 0
	for _, t := range tracks {
		if t.Kind == TrackAuto {
			autoCount++
		}
	}
	if totalFr > 0 {
		per := 0.0
		if math.IsInf(remaining, 0) || math.IsNaN(remaining) {
			per = 0
		} else {
			per = remaining / totalFr
			if per < 0 {
				per = 0
			}
		}
		for i, t := range tracks {
			if t.Kind == TrackFr {
				sizes[i] = per * t.Fr
			}
		}
		remaining = 0
	}
	if autoCount > 0 && remaining > 0 && !math.IsInf(remaining, 0) && !math.IsNaN(remaining) {
		per := remaining / float64(autoCount)
		for i, t := range tracks {
			if t.Kind == TrackAuto {
				sizes[i] = per
			}
		}
	}
	return sizes
}

func hasAutoRows(template []TrackSize, auto TrackSize, count int) bool {
	for _, t := range template {
		if t.Kind == TrackAuto {
			return true
		}
	}
	if count > len(template) {
		if auto.Kind == TrackAuto {
			return true
		}
		if auto.Kind == TrackFixed && auto.Length.IsUndefined() {
			return true
		}
	}
	return false
}

func measureAutoRows(items []*Node, placements []gridSpot, colSizes []float64, colGap float64, rowSizes []float64, template []TrackSize, autoRow TrackSize, contentHeight float64) {
	autoRowIdx := map[int]bool{}
	for i := range rowSizes {
		if i < len(template) {
			if template[i].Kind == TrackAuto {
				autoRowIdx[i] = true
			}
		} else if autoRow.Kind == TrackAuto || (autoRow.Kind == TrackFixed && autoRow.Length.IsUndefined()) {
			autoRowIdx[i] = true
		}
	}
	if len(autoRowIdx) == 0 {
		return
	}
	for i, child := range items {
		p := placements[i]
		w := trackSpanSize(colSizes, colGap, p.colStart-1, p.colEnd-1)
		childIn := layoutRun{
			availableWidth:  w,
			availableHeight: math.Inf(1),
			ownerWidth:      w,
			ownerHeight:     contentHeight,
			definiteWidth:   true,
			definiteHeight:  false,
		}
		computeNode(child, childIn)
		itemHeight := child.Layout.Height
		spans := p.rowEnd - p.rowStart
		for r := p.rowStart - 1; r < p.rowEnd-1; r++ {
			if autoRowIdx[r] {
				share := itemHeight / float64(spans)
				if share > rowSizes[r] {
					rowSizes[r] = share
				}
			}
		}
	}
}

func trackOffsets(sizes []float64, gap float64) []float64 {
	offsets := make([]float64, len(sizes))
	cur := 0.0
	for i, s := range sizes {
		offsets[i] = cur
		cur += s + gap
	}
	return offsets
}

func trackSpanSize(sizes []float64, gap float64, start, end int) float64 {
	if start < 0 {
		start = 0
	}
	if end > len(sizes) {
		end = len(sizes)
	}
	if end <= start {
		return 0
	}
	total := 0.0
	for i := start; i < end; i++ {
		total += sizes[i]
	}
	if end-start > 1 {
		total += gap * float64(end-start-1)
	}
	return total
}
