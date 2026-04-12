package layout

import "math"

type flexItem struct {
	node        *Node
	baseSize    float64
	mainSize    float64
	crossSize   float64
	frozen      bool
	mainOffset  float64
	crossOffset float64
}

type flexLine struct {
	items       []*flexItem
	crossSize   float64
	crossOffset float64
}

func computeNode(node *Node, availableWidth, availableHeight float64) {
	s := &node.Style

	paddingTop := resolveOr0(s.PaddingTop, availableHeight)
	paddingRight := resolveOr0(s.PaddingRight, availableWidth)
	paddingBottom := resolveOr0(s.PaddingBottom, availableHeight)
	paddingLeft := resolveOr0(s.PaddingLeft, availableWidth)

	node.Layout.Padding = [4]float64{paddingTop, paddingRight, paddingBottom, paddingLeft}
	node.Layout.Border = [4]float64{s.BorderTop, s.BorderRight, s.BorderBottom, s.BorderLeft}

	bpH := paddingTop + paddingBottom + s.BorderTop + s.BorderBottom
	bpW := paddingLeft + paddingRight + s.BorderLeft + s.BorderRight

	containerWidth := resolveNodeSize(s.Width, availableWidth)
	if math.IsNaN(containerWidth) {
		containerWidth = availableWidth
	}
	containerHeight := resolveNodeSize(s.Height, availableHeight)
	if math.IsNaN(containerHeight) {
		containerHeight = availableHeight
	}

	containerWidth = clampSize(containerWidth, s.MinWidth, s.MaxWidth, availableWidth)
	containerHeight = clampSize(containerHeight, s.MinHeight, s.MaxHeight, availableHeight)

	if s.AspectRatio > 0 {
		if s.Width.IsDefined() && !s.Height.IsDefined() {
			containerHeight = containerWidth / s.AspectRatio
		} else if s.Height.IsDefined() && !s.Width.IsDefined() {
			containerWidth = containerHeight * s.AspectRatio
		}
	}

	node.Layout.Width = containerWidth
	node.Layout.Height = containerHeight

	contentWidth := containerWidth - bpW
	contentHeight := containerHeight - bpH
	if contentWidth < 0 {
		contentWidth = 0
	}
	if contentHeight < 0 {
		contentHeight = 0
	}

	if node.Measure != nil {
		w, h := node.Measure(contentWidth, contentHeight)
		if !node.Style.Width.IsDefined() {
			node.Layout.Width = w + bpW
		}
		if !node.Style.Height.IsDefined() {
			node.Layout.Height = h + bpH
		}
		return
	}

	if len(node.Children) == 0 {
		return
	}

	isRow := s.Direction == Row || s.Direction == RowReverse
	isReverse := s.Direction == RowReverse || s.Direction == ColumnReverse

	var mainSize, crossSize float64
	if isRow {
		mainSize = contentWidth
		crossSize = contentHeight
	} else {
		mainSize = contentHeight
		crossSize = contentWidth
	}

	gap := resolveGap(s, isRow)

	var flexChildren []*Node
	var absChildren []*Node
	for _, child := range node.Children {
		if child.Style.Display == DisplayNone {
			continue
		}
		if child.Style.Position == Absolute {
			absChildren = append(absChildren, child)
			continue
		}
		flexChildren = append(flexChildren, child)
	}

	items := make([]*flexItem, len(flexChildren))
	for i, child := range flexChildren {
		base := determineFlexBaseSize(child, mainSize, isRow, containerWidth, containerHeight)
		items[i] = &flexItem{
			node:     child,
			baseSize: base,
			mainSize: base,
		}
	}

	lines := collectIntoLines(items, mainSize, s.Wrap, gap)

	for _, line := range lines {
		resolveFlexibleLengths(line, mainSize, gap)
	}

	for _, line := range lines {
		for _, item := range line.items {
			cs := determineCrossSize(item, crossSize, isRow, containerWidth, containerHeight)
			item.crossSize = cs
		}
		maxCross := 0.0
		for _, item := range line.items {
			mc := item.crossSize + marginCross(item.node, crossSize, isRow)
			if mc > maxCross {
				maxCross = mc
			}
		}
		line.crossSize = maxCross
	}

	applyAlignContentStretch(lines, crossSize, s.AlignContent)

	for _, line := range lines {
		applyAlignStretch(line, node.Style.AlignItems, isRow, containerWidth, containerHeight)
	}

	for _, line := range lines {
		for _, item := range line.items {
			cs := &item.node.Style
			mTop := resolveOr0(cs.MarginTop, containerHeight)
			mRight := resolveOr0(cs.MarginRight, containerWidth)
			mBottom := resolveOr0(cs.MarginBottom, containerHeight)
			mLeft := resolveOr0(cs.MarginLeft, containerWidth)

			align := resolveAlign(item.node, s.AlignItems)

			var hasCrossDim bool
			if isRow {
				hasCrossDim = item.node.Style.Height.IsDefined()
			} else {
				hasCrossDim = item.node.Style.Width.IsDefined()
			}

			effectiveCross := line.crossSize
			if align != AlignStretch && !hasCrossDim {
				effectiveCross = item.crossSize + marginCross(item.node, line.crossSize, isRow)
			}

			var availCross float64
			if isRow {
				availCross = effectiveCross - mTop - mBottom
			} else {
				availCross = effectiveCross - mLeft - mRight
			}
			if availCross < 0 {
				availCross = 0
			}

			var iw, ih float64
			if isRow {
				iw = item.mainSize - mLeft - mRight
				ih = availCross
			} else {
				iw = availCross
				ih = item.mainSize - mTop - mBottom
			}
			if iw < 0 {
				iw = 0
			}
			if ih < 0 {
				ih = 0
			}
			computeNode(item.node, iw, ih)

			if item.node.Measure != nil {
				if align == AlignStretch {
					if isRow && !item.node.Style.Height.IsDefined() {
						item.node.Layout.Height = item.crossSize
					} else if !isRow && !item.node.Style.Width.IsDefined() {
						item.node.Layout.Width = item.crossSize
					}
				}
			}
		}
	}

	for _, line := range lines {
		mainAxisAlignment(line, mainSize, s.JustifyContent, gap, isRow)
	}

	crossAxisAlignment(lines, crossSize, s.AlignItems, s.AlignContent, isRow)

	if !s.Height.IsDefined() {
		if isRow {
			var totalCross float64
			for _, line := range lines {
				totalCross += line.crossSize
			}
			if totalCross+bpH < node.Layout.Height {
				node.Layout.Height = totalCross + bpH
			}
		} else {
			var totalMain float64
			for _, line := range lines {
				lineMain := 0.0
				for i, item := range line.items {
					lineMain += item.mainSize + marginStart(item.node, mainSize, isRow) + marginEnd(item.node, mainSize, isRow)
					if i > 0 {
						lineMain += gap
					}
				}
				if lineMain > totalMain {
					totalMain = lineMain
				}
			}
			if totalMain+bpH < node.Layout.Height {
				node.Layout.Height = totalMain + bpH
			}
		}
	}

	insetX := paddingLeft + s.BorderLeft
	insetY := paddingTop + s.BorderTop

	for _, line := range lines {
		for _, item := range line.items {
			if isRow {
				item.node.Layout.X = insetX + item.mainOffset
				item.node.Layout.Y = insetY + item.crossOffset
			} else {
				item.node.Layout.X = insetX + item.crossOffset
				item.node.Layout.Y = insetY + item.mainOffset
			}
		}
	}

	if isReverse {
		for _, line := range lines {
			for _, item := range line.items {
				ms := marginStart(item.node, mainSize, isRow)
				reversed := mainSize - item.mainOffset - item.mainSize - ms + marginEnd(item.node, mainSize, isRow)
				newOffset := reversed - marginEnd(item.node, mainSize, isRow) + ms
				if isRow {
					item.node.Layout.X = insetX + newOffset
				} else {
					item.node.Layout.Y = insetY + newOffset
				}
			}
		}
	}

	if s.Wrap == WrapReverse {
		totalCross := 0.0
		for _, line := range lines {
			totalCross += line.crossSize
		}
		for _, line := range lines {
			for _, item := range line.items {
				inLineCross := item.crossOffset - line.crossOffset
				reversed := totalCross - line.crossOffset - line.crossSize + inLineCross
				if isRow {
					item.node.Layout.Y = insetY + reversed
				} else {
					item.node.Layout.X = insetX + reversed
				}
			}
		}
	}

	for _, child := range absChildren {
		computeAbsolute(child, node)
	}
}

func resolveGap(s *Style, isRow bool) float64 {
	if isRow {
		if s.ColumnGap > 0 {
			return s.ColumnGap
		}
		return s.Gap
	}
	if s.RowGap > 0 {
		return s.RowGap
	}
	return s.Gap
}

func resolveOr0(d Dimension, base float64) float64 {
	v := d.Resolve(base)
	if math.IsNaN(v) {
		return 0
	}
	return v
}

func resolveNodeSize(d Dimension, base float64) float64 {
	if d.Unit == UnitPoint {
		return d.Value
	}
	if d.Unit == UnitPercent {
		if math.IsInf(base, 0) || math.IsNaN(base) {
			return math.NaN()
		}
		return d.Value / 100 * base
	}
	return math.NaN()
}

func clampSize(size float64, minD, maxD Dimension, base float64) float64 {
	if math.IsNaN(size) {
		return size
	}
	minV := resolveNodeSize(minD, base)
	maxV := resolveNodeSize(maxD, base)
	if !math.IsNaN(minV) && size < minV {
		size = minV
	}
	if !math.IsNaN(maxV) && size > maxV {
		size = maxV
	}
	return size
}

func determineFlexBaseSize(child *Node, mainSize float64, isRow bool, containerW, containerH float64) float64 {
	s := &child.Style

	if s.FlexBasis.IsDefined() {
		v := s.FlexBasis.Resolve(mainSize)
		if !math.IsNaN(v) {
			var minBP float64
			if isRow {
				minBP = resolveOr0(s.PaddingLeft, containerW) + resolveOr0(s.PaddingRight, containerW) + s.BorderLeft + s.BorderRight
			} else {
				minBP = resolveOr0(s.PaddingTop, containerH) + resolveOr0(s.PaddingBottom, containerH) + s.BorderTop + s.BorderBottom
			}
			if v < minBP {
				v = minBP
			}
			return v
		}
	}

	var sizeD Dimension
	var base float64
	if isRow {
		sizeD = s.Width
		base = containerW
	} else {
		sizeD = s.Height
		base = containerH
	}
	if sizeD.IsDefined() {
		v := sizeD.Resolve(base)
		if !math.IsNaN(v) {
			return v
		}
	}

	if child.Measure != nil {
		var w, h float64
		if isRow {
			w, h = child.Measure(mainSize, math.Inf(1))
		} else {
			w, h = child.Measure(math.Inf(1), mainSize)
		}
		if isRow {
			return w
		}
		return h
	}

	return estimateSize(child, isRow, containerW, containerH)
}

func estimateSize(node *Node, wantWidth bool, containerW, containerH float64) float64 {
	s := &node.Style

	var sizeD Dimension
	var base float64
	if wantWidth {
		sizeD = s.Width
		base = containerW
	} else {
		sizeD = s.Height
		base = containerH
	}
	if sizeD.IsDefined() {
		v := sizeD.Resolve(base)
		if !math.IsNaN(v) {
			return v
		}
	}

	if node.Measure != nil {
		w, h := node.Measure(containerW, math.Inf(1))
		if wantWidth {
			return w
		}
		return h
	}

	nodeIsRow := s.Direction == Row || s.Direction == RowReverse

	var padding float64
	if wantWidth {
		padding = resolveOr0(s.PaddingLeft, containerW) + resolveOr0(s.PaddingRight, containerW) + s.BorderLeft + s.BorderRight
	} else {
		padding = resolveOr0(s.PaddingTop, containerH) + resolveOr0(s.PaddingBottom, containerH) + s.BorderTop + s.BorderBottom
	}

	mainIsWanted := nodeIsRow == wantWidth
	gap := resolveGap(s, nodeIsRow)
	total := 0.0
	count := 0

	for _, child := range node.Children {
		if child.Style.Display == DisplayNone || child.Style.Position == Absolute {
			continue
		}

		childSize := estimateSize(child, wantWidth, containerW, containerH)

		var ms float64
		if wantWidth {
			ms = resolveOr0(child.Style.MarginLeft, containerW) + resolveOr0(child.Style.MarginRight, containerW)
		} else {
			ms = resolveOr0(child.Style.MarginTop, containerH) + resolveOr0(child.Style.MarginBottom, containerH)
		}

		if mainIsWanted {
			total += childSize + ms
			if count > 0 {
				total += gap
			}
		} else {
			if childSize+ms > total {
				total = childSize + ms
			}
		}
		count++
	}
	return total + padding
}

func collectIntoLines(items []*flexItem, mainSize float64, wrap Wrap, gap float64) []*flexLine {
	if len(items) == 0 {
		return []*flexLine{{}}
	}

	if wrap == NoWrap {
		return []*flexLine{{items: items}}
	}

	var lines []*flexLine
	var current []*flexItem
	lineMain := 0.0

	for _, item := range items {
		itemSize := item.baseSize
		gapAdd := 0.0
		if len(current) > 0 {
			gapAdd = gap
		}
		if len(current) > 0 && lineMain+gapAdd+itemSize > mainSize {
			lines = append(lines, &flexLine{items: current})
			current = nil
			lineMain = 0
			gapAdd = 0
		}
		current = append(current, item)
		lineMain += gapAdd + itemSize
	}
	if len(current) > 0 {
		lines = append(lines, &flexLine{items: current})
	}

	return lines
}

func resolveFlexibleLengths(line *flexLine, mainSize float64, gap float64) {
	items := line.items
	if len(items) == 0 {
		return
	}

	totalGap := gap * float64(len(items)-1)

	sumBases := 0.0
	for _, item := range items {
		sumBases += item.baseSize
	}

	growing := sumBases+totalGap <= mainSize

	for _, item := range items {
		gs := flexGrow(item.node)
		ss := flexShrink(item.node)

		if growing && gs == 0 {
			item.frozen = true
			item.mainSize = item.baseSize
		} else if !growing && ss == 0 {
			item.frozen = true
			item.mainSize = item.baseSize
		}
	}

	for iter := 0; iter < 10; iter++ {
		allFrozen := true
		for _, item := range items {
			if !item.frozen {
				allFrozen = false
				break
			}
		}
		if allFrozen {
			break
		}

		freeSpace := mainSize - totalGap
		for _, item := range items {
			if item.frozen {
				freeSpace -= item.mainSize
			} else {
				freeSpace -= item.baseSize
			}
		}

		if growing {
			totalFactor := 0.0
			for _, item := range items {
				if !item.frozen {
					totalFactor += flexGrow(item.node)
				}
			}
			if totalFactor > 0 {
				for _, item := range items {
					if !item.frozen {
						ratio := flexGrow(item.node) / totalFactor
						item.mainSize = item.baseSize + ratio*freeSpace
					}
				}
			}
		} else {
			totalScaled := 0.0
			for _, item := range items {
				if !item.frozen {
					totalScaled += flexShrink(item.node) * item.baseSize
				}
			}
			if totalScaled > 0 {
				for _, item := range items {
					if !item.frozen {
						scaled := flexShrink(item.node) * item.baseSize
						ratio := scaled / totalScaled
						item.mainSize = item.baseSize + ratio*freeSpace
					}
				}
			}
		}

		adjustments := 0.0
		for _, item := range items {
			if item.frozen {
				continue
			}
			clamped := applyMainMinMax(item)
			adjustments += clamped - item.mainSize
		}

		if adjustments == 0 {
			for _, item := range items {
				if !item.frozen {
					item.mainSize = applyMainMinMax(item)
				}
				item.frozen = true
			}
		} else {
			for _, item := range items {
				if item.frozen {
					continue
				}
				clamped := applyMainMinMax(item)
				if clamped != item.mainSize {
					item.mainSize = clamped
					item.frozen = true
				}
			}
		}
	}

	for _, item := range items {
		if item.mainSize < 0 {
			item.mainSize = 0
		}
	}
}

func applyMainMinMax(item *flexItem) float64 {
	s := &item.node.Style
	size := item.mainSize

	minV := resolveNodeSize(s.MinWidth, math.Inf(1))
	maxV := resolveNodeSize(s.MaxWidth, math.Inf(1))

	if !math.IsNaN(minV) && size < minV {
		size = minV
	}
	if !math.IsNaN(maxV) && size > maxV {
		size = maxV
	}
	return size
}

func determineCrossSize(item *flexItem, crossSize float64, isRow bool, containerW, containerH float64) float64 {
	s := &item.node.Style

	var crossD Dimension
	var base float64
	if isRow {
		crossD = s.Height
		base = containerH
	} else {
		crossD = s.Width
		base = containerW
	}

	if crossD.IsDefined() {
		v := crossD.Resolve(base)
		if !math.IsNaN(v) {
			return v
		}
	}

	if s.AspectRatio > 0 {
		if isRow {
			return item.mainSize / s.AspectRatio
		}
		return item.mainSize * s.AspectRatio
	}

	if item.node.Measure != nil {
		var w, h float64
		if isRow {
			w, h = item.node.Measure(item.mainSize, crossSize)
		} else {
			w, h = item.node.Measure(crossSize, item.mainSize)
		}
		if isRow {
			return h
		}
		return w
	}

	return estimateSize(item.node, !isRow, containerW, containerH)
}

func applyAlignContentStretch(lines []*flexLine, crossSize float64, ac Align) {
	if (ac != AlignStretch && ac != AlignAuto) || math.IsInf(crossSize, 0) || math.IsNaN(crossSize) {
		return
	}
	totalLineCross := 0.0
	for _, line := range lines {
		totalLineCross += line.crossSize
	}
	remaining := crossSize - totalLineCross
	if remaining <= 0 || len(lines) == 0 {
		return
	}
	extra := remaining / float64(len(lines))
	for _, line := range lines {
		line.crossSize += extra
	}
}

func applyAlignStretch(line *flexLine, alignItems Align, isRow bool, containerW, containerH float64) {
	for _, item := range line.items {
		align := resolveAlign(item.node, alignItems)
		if align != AlignStretch {
			continue
		}

		var crossD Dimension
		if isRow {
			crossD = item.node.Style.Height
		} else {
			crossD = item.node.Style.Width
		}
		if crossD.IsDefined() {
			continue
		}

		mc := marginCross(item.node, line.crossSize, isRow)
		item.crossSize = line.crossSize - mc
		if item.crossSize < 0 {
			item.crossSize = 0
		}
	}
}

func mainAxisAlignment(line *flexLine, mainSize float64, justify Justify, gap float64, isRow bool) {
	items := line.items
	n := len(items)
	if n == 0 {
		return
	}

	usedMain := 0.0
	for _, item := range items {
		usedMain += item.mainSize + marginMain(item.node, mainSize, isRow)
	}
	totalGap := gap * float64(n-1)
	freeSpace := mainSize - usedMain - totalGap
	if freeSpace < 0 {
		freeSpace = 0
	}

	autoMargins := 0
	for _, item := range items {
		if isRow {
			if item.node.Style.MarginLeft.IsAuto() {
				autoMargins++
			}
			if item.node.Style.MarginRight.IsAuto() {
				autoMargins++
			}
		} else {
			if item.node.Style.MarginTop.IsAuto() {
				autoMargins++
			}
			if item.node.Style.MarginBottom.IsAuto() {
				autoMargins++
			}
		}
	}

	if autoMargins > 0 {
		perMargin := freeSpace / float64(autoMargins)
		offset := 0.0
		for i, item := range items {
			ms := marginStartAutoResolved(item.node, mainSize, isRow, perMargin)
			me := marginEndAutoResolved(item.node, mainSize, isRow, perMargin)
			item.mainOffset = offset + ms
			offset = item.mainOffset + item.mainSize + me
			if i < n-1 {
				offset += gap
			}
		}
		return
	}

	var startOffset, spacing float64
	switch justify {
	case JustifyStart:
		startOffset = 0
		spacing = gap
	case JustifyEnd:
		startOffset = freeSpace
		spacing = gap
	case JustifyCenter:
		startOffset = freeSpace / 2
		spacing = gap
	case JustifySpaceBetween:
		startOffset = 0
		if n > 1 {
			spacing = gap + freeSpace/float64(n-1)
		}
	case JustifySpaceAround:
		perItem := freeSpace / float64(n)
		startOffset = perItem / 2
		spacing = gap + perItem
	case JustifySpaceEvenly:
		sp := freeSpace / float64(n+1)
		startOffset = sp
		spacing = gap + sp
	}

	offset := startOffset
	for i, item := range items {
		ms := marginStart(item.node, mainSize, isRow)
		item.mainOffset = offset + ms
		offset = item.mainOffset + item.mainSize + marginEnd(item.node, mainSize, isRow)
		if i < n-1 {
			offset += spacing
		}
	}
}

func crossAxisAlignment(lines []*flexLine, crossSize float64, alignItems Align, alignContent Align, isRow bool) {
	n := len(lines)
	if n == 0 {
		return
	}

	totalLineCross := 0.0
	for _, line := range lines {
		totalLineCross += line.crossSize
	}
	freeSpace := crossSize - totalLineCross
	if math.IsInf(freeSpace, 0) || math.IsNaN(freeSpace) || freeSpace < 0 {
		freeSpace = 0
	}

	var lineStart, lineSpacing float64
	switch alignContent {
	case AlignStart, AlignAuto, AlignStretch:
		lineStart = 0
	case AlignEnd:
		lineStart = freeSpace
	case AlignCenter:
		lineStart = freeSpace / 2
	case AlignBaseline:
		lineStart = 0
	}

	offset := lineStart
	for _, line := range lines {
		line.crossOffset = offset

		for _, item := range line.items {
			align := resolveAlign(item.node, alignItems)
			cms := crossMarginStart(item, crossSize, isRow)
			cme := crossMarginEnd(item, crossSize, isRow)
			availCross := line.crossSize - cms - cme

			switch align {
			case AlignStart:
				item.crossOffset = line.crossOffset + cms
			case AlignEnd:
				item.crossOffset = line.crossOffset + cms + availCross - item.crossSize
			case AlignCenter:
				item.crossOffset = line.crossOffset + cms + (availCross-item.crossSize)/2
			default:
				item.crossOffset = line.crossOffset + cms
			}
		}

		offset += line.crossSize + lineSpacing
	}
}

func resolveAlign(node *Node, alignItems Align) Align {
	if node.Style.AlignSelf != AlignAuto {
		return node.Style.AlignSelf
	}
	if alignItems == AlignAuto {
		return AlignStretch
	}
	return alignItems
}

func flexGrow(node *Node) float64 {
	return node.Style.FlexGrow
}

func flexShrink(node *Node) float64 {
	if node.Style.FlexShrink == 0 && node.Style.FlexGrow == 0 && node.Style.FlexBasis.IsUndefined() {
		return 1
	}
	return node.Style.FlexShrink
}

func marginMain(node *Node, mainSize float64, isRow bool) float64 {
	return marginStart(node, mainSize, isRow) + marginEnd(node, mainSize, isRow)
}

func marginStart(node *Node, mainSize float64, isRow bool) float64 {
	var d Dimension
	if isRow {
		d = node.Style.MarginLeft
	} else {
		d = node.Style.MarginTop
	}
	if d.IsAuto() {
		return 0
	}
	return resolveOr0(d, mainSize)
}

func marginEnd(node *Node, mainSize float64, isRow bool) float64 {
	var d Dimension
	if isRow {
		d = node.Style.MarginRight
	} else {
		d = node.Style.MarginBottom
	}
	if d.IsAuto() {
		return 0
	}
	return resolveOr0(d, mainSize)
}

func marginCross(node *Node, crossSize float64, isRow bool) float64 {
	return crossMarginStartNode(node, crossSize, isRow) + crossMarginEndNode(node, crossSize, isRow)
}

func crossMarginStartNode(node *Node, crossSize float64, isRow bool) float64 {
	var d Dimension
	if isRow {
		d = node.Style.MarginTop
	} else {
		d = node.Style.MarginLeft
	}
	if d.IsAuto() {
		return 0
	}
	return resolveOr0(d, crossSize)
}

func crossMarginEndNode(node *Node, crossSize float64, isRow bool) float64 {
	var d Dimension
	if isRow {
		d = node.Style.MarginBottom
	} else {
		d = node.Style.MarginRight
	}
	if d.IsAuto() {
		return 0
	}
	return resolveOr0(d, crossSize)
}

func crossMarginStart(item *flexItem, crossSize float64, isRow bool) float64 {
	return crossMarginStartNode(item.node, crossSize, isRow)
}

func crossMarginEnd(item *flexItem, crossSize float64, isRow bool) float64 {
	return crossMarginEndNode(item.node, crossSize, isRow)
}

func marginStartAutoResolved(node *Node, mainSize float64, isRow bool, autoValue float64) float64 {
	var d Dimension
	if isRow {
		d = node.Style.MarginLeft
	} else {
		d = node.Style.MarginTop
	}
	if d.IsAuto() {
		return autoValue
	}
	return resolveOr0(d, mainSize)
}

func marginEndAutoResolved(node *Node, mainSize float64, isRow bool, autoValue float64) float64 {
	var d Dimension
	if isRow {
		d = node.Style.MarginRight
	} else {
		d = node.Style.MarginBottom
	}
	if d.IsAuto() {
		return autoValue
	}
	return resolveOr0(d, mainSize)
}
