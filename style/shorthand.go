package style

import "strings"

func splitValues(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}

	var parts []string
	var buf strings.Builder
	depth := 0
	inSingle := false
	inDouble := false

	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c == '\'' && !inDouble:
			inSingle = !inSingle
			buf.WriteByte(c)
		case c == '"' && !inSingle:
			inDouble = !inDouble
			buf.WriteByte(c)
		case !inSingle && !inDouble && c == '(':
			depth++
			buf.WriteByte(c)
		case !inSingle && !inDouble && c == ')':
			if depth > 0 {
				depth--
			}
			buf.WriteByte(c)
		case (c == ' ' || c == '\t') && depth == 0 && !inSingle && !inDouble:
			if buf.Len() > 0 {
				parts = append(parts, buf.String())
				buf.Reset()
			}
		default:
			buf.WriteByte(c)
		}
	}
	if buf.Len() > 0 {
		parts = append(parts, buf.String())
	}
	return parts
}

// ExpandShorthands expands CSS shorthand properties (margin, padding, border, etc.) into their longhand forms.
// ExpandShorthands expands CSS shorthand properties to longhand.
func ExpandShorthands(props map[string]string) map[string]string {
	result := make(map[string]string, len(props))
	for k, v := range props {
		result[k] = v
	}

	for prop, val := range props {
		switch prop {
		case "margin":
			expandBoxSides(result, "margin", val)
		case "padding":
			expandBoxSides(result, "padding", val)
		case "border":
			expandBorder(result, val, []string{"top", "right", "bottom", "left"})
		case "border-top", "border-right", "border-bottom", "border-left":
			side := strings.TrimPrefix(prop, "border-")
			expandBorder(result, val, []string{side})
		case "border-radius":
			expandBorderRadius(result, val)
		case "border-width":
			expandBorderSideProp(result, "width", val)
		case "border-style":
			expandBorderSideProp(result, "style", val)
		case "border-color":
			expandBorderSideProp(result, "color", val)
		case "flex":
			expandFlex(result, val)
		case "gap":
			expandGap(result, val)
		case "background":
			expandBackground(result, val)
		case "font":
			expandFont(result, val)
		case "text-decoration":
			expandTextDecoration(result, val)
		case "overflow":
			expandOverflow(result, val)
		}
	}

	return result
}

func expandBoxSides(result map[string]string, prefix, val string) {
	parts := splitValues(val)
	var top, right, bottom, left string
	switch len(parts) {
	case 1:
		top, right, bottom, left = parts[0], parts[0], parts[0], parts[0]
	case 2:
		top, bottom = parts[0], parts[0]
		right, left = parts[1], parts[1]
	case 3:
		top = parts[0]
		right, left = parts[1], parts[1]
		bottom = parts[2]
	case 4:
		top, right, bottom, left = parts[0], parts[1], parts[2], parts[3]
	default:
		return
	}
	result[prefix+"-top"] = top
	result[prefix+"-right"] = right
	result[prefix+"-bottom"] = bottom
	result[prefix+"-left"] = left
}

func expandBorder(result map[string]string, val string, sides []string) {
	parts := splitValues(val)
	var width, bstyle, color string

	for _, p := range parts {
		if isBorderStyle(p) {
			bstyle = p
		} else if looksLikeLength(p) {
			width = p
		} else {
			color = p
		}
	}

	for _, side := range sides {
		if width != "" {
			result["border-"+side+"-width"] = width
		}
		if bstyle != "" {
			result["border-"+side+"-style"] = bstyle
		}
		if color != "" {
			result["border-"+side+"-color"] = color
		}
	}
}

func expandBorderRadius(result map[string]string, val string) {
	horizontal, vertical := val, ""
	if idx := strings.Index(val, "/"); idx >= 0 {
		horizontal = strings.TrimSpace(val[:idx])
		vertical = strings.TrimSpace(val[idx+1:])
	}

	hParts := splitValues(horizontal)
	tl, tr, br, bl := expandFourValues(hParts)

	if vertical == "" {
		result["border-top-left-radius"] = tl
		result["border-top-right-radius"] = tr
		result["border-bottom-right-radius"] = br
		result["border-bottom-left-radius"] = bl
		return
	}

	vParts := splitValues(vertical)
	vtl, vtr, vbr, vbl := expandFourValues(vParts)
	result["border-top-left-radius"] = tl + " " + vtl
	result["border-top-right-radius"] = tr + " " + vtr
	result["border-bottom-right-radius"] = br + " " + vbr
	result["border-bottom-left-radius"] = bl + " " + vbl
}

func expandFourValues(parts []string) (tl, tr, br, bl string) {
	switch len(parts) {
	case 1:
		return parts[0], parts[0], parts[0], parts[0]
	case 2:
		return parts[0], parts[1], parts[0], parts[1]
	case 3:
		return parts[0], parts[1], parts[2], parts[1]
	case 4:
		return parts[0], parts[1], parts[2], parts[3]
	}
	return
}

func expandBorderSideProp(result map[string]string, prop, val string) {
	parts := splitValues(val)
	var top, right, bottom, left string
	switch len(parts) {
	case 1:
		top, right, bottom, left = parts[0], parts[0], parts[0], parts[0]
	case 2:
		top, bottom = parts[0], parts[0]
		right, left = parts[1], parts[1]
	case 3:
		top = parts[0]
		right, left = parts[1], parts[1]
		bottom = parts[2]
	case 4:
		top, right, bottom, left = parts[0], parts[1], parts[2], parts[3]
	default:
		return
	}
	result["border-top-"+prop] = top
	result["border-right-"+prop] = right
	result["border-bottom-"+prop] = bottom
	result["border-left-"+prop] = left
}

func expandFlex(result map[string]string, val string) {
	switch val {
	case "none":
		result["flex-grow"] = "0"
		result["flex-shrink"] = "0"
		result["flex-basis"] = "auto"
		return
	case "auto":
		result["flex-grow"] = "1"
		result["flex-shrink"] = "1"
		result["flex-basis"] = "auto"
		return
	}

	parts := splitValues(val)
	switch len(parts) {
	case 1:
		result["flex-grow"] = parts[0]
		result["flex-shrink"] = "1"
		result["flex-basis"] = "0%"
	case 2:
		result["flex-grow"] = parts[0]
		if looksLikeLengthWithUnit(parts[1]) {
			result["flex-shrink"] = "1"
			result["flex-basis"] = parts[1]
		} else {
			result["flex-shrink"] = parts[1]
			result["flex-basis"] = "0%"
		}
	case 3:
		result["flex-grow"] = parts[0]
		result["flex-shrink"] = parts[1]
		result["flex-basis"] = parts[2]
	}
}

func expandGap(result map[string]string, val string) {
	parts := splitValues(val)
	switch len(parts) {
	case 1:
		result["row-gap"] = parts[0]
		result["column-gap"] = parts[0]
	case 2:
		result["row-gap"] = parts[0]
		result["column-gap"] = parts[1]
	}
}

func expandBackground(result map[string]string, val string) {
	lower := strings.ToLower(val)
	if strings.Contains(lower, "gradient") || strings.Contains(lower, "url(") {
		result["background-image"] = val
	} else {
		result["background-color"] = val
	}
}

func expandFont(result map[string]string, val string) {
	parts := splitValues(val)
	if len(parts) == 0 {
		return
	}

	i := 0
	for i < len(parts) {
		lower := strings.ToLower(parts[i])
		if isFontStyle(lower) {
			result["font-style"] = lower
			i++
		} else if isFontWeight(lower) {
			result["font-weight"] = lower
			i++
		} else {
			break
		}
	}

	if i >= len(parts) {
		return
	}

	sizeVal := parts[i]
	if idx := strings.Index(sizeVal, "/"); idx >= 0 {
		result["font-size"] = sizeVal[:idx]
		result["line-height"] = sizeVal[idx+1:]
	} else {
		result["font-size"] = sizeVal
	}
	i++

	if i < len(parts) {
		result["font-family"] = strings.Join(parts[i:], " ")
	}
}

func expandTextDecoration(result map[string]string, val string) {
	parts := splitValues(val)
	for _, p := range parts {
		lower := strings.ToLower(p)
		if isTextDecorationLine(lower) {
			result["text-decoration-line"] = lower
		} else if isTextDecorationStyle(lower) {
			result["text-decoration-style"] = lower
		} else {
			result["text-decoration-color"] = p
		}
	}
}

func expandOverflow(result map[string]string, val string) {
	parts := splitValues(val)
	switch len(parts) {
	case 1:
		result["overflow-x"] = parts[0]
		result["overflow-y"] = parts[0]
	case 2:
		result["overflow-x"] = parts[0]
		result["overflow-y"] = parts[1]
	}
}

func isBorderStyle(s string) bool {
	switch strings.ToLower(s) {
	case "none", "solid", "dashed", "dotted", "double", "groove", "ridge", "inset", "outset", "hidden":
		return true
	}
	return false
}

func looksLikeLength(s string) bool {
	if s == "0" {
		return true
	}
	for _, suffix := range []string{"px", "em", "rem", "%", "vw", "vh", "pt", "cm", "mm", "in"} {
		if strings.HasSuffix(s, suffix) {
			return true
		}
	}
	if len(s) > 0 && (s[0] >= '0' && s[0] <= '9' || s[0] == '.' || s[0] == '-') {
		return true
	}
	return false
}

func looksLikeLengthWithUnit(s string) bool {
	for _, suffix := range []string{"px", "em", "rem", "%", "vw", "vh", "pt", "cm", "mm", "in"} {
		if strings.HasSuffix(s, suffix) {
			return true
		}
	}
	if s == "auto" || s == "content" {
		return true
	}
	return false
}

func isFontStyle(s string) bool {
	switch s {
	case "italic", "oblique", "normal":
		return true
	}
	return false
}

func isFontWeight(s string) bool {
	switch s {
	case "bold", "bolder", "lighter", "100", "200", "300", "400", "500", "600", "700", "800", "900":
		return true
	}
	return false
}

func isTextDecorationLine(s string) bool {
	switch s {
	case "none", "underline", "overline", "line-through":
		return true
	}
	return false
}

func isTextDecorationStyle(s string) bool {
	switch s {
	case "solid", "double", "dotted", "dashed", "wavy":
		return true
	}
	return false
}
