package style

import (
	"strconv"
	"strings"

	"github.com/macawls/ogre/parse"
)

const defaultRootFontSize = 16.0

var inheritedProperties = map[string]bool{
	"color":                  true,
	"font-family":            true,
	"font-size":              true,
	"font-weight":            true,
	"font-style":             true,
	"line-height":            true,
	"letter-spacing":         true,
	"text-align":             true,
	"text-transform":         true,
	"text-decoration-line":   true,
	"text-decoration-color":  true,
	"text-decoration-style":  true,
	"white-space":            true,
	"word-break":             true,
	"text-shadow":            true,
}

var tagDefaults = map[string]map[string]string{
	"div": {
		"display":        "flex",
		"flex-direction":  "row",
	},
	"p": {
		"display":       "flex",
		"margin-top":    "1em",
		"margin-bottom": "1em",
	},
	"h1": {
		"font-size":     "2em",
		"font-weight":   "700",
		"margin-top":    "0.67em",
		"margin-bottom": "0.67em",
	},
	"h2": {
		"font-size":     "1.5em",
		"font-weight":   "700",
		"margin-top":    "0.83em",
		"margin-bottom": "0.83em",
	},
	"h3": {
		"font-size":     "1.17em",
		"font-weight":   "700",
		"margin-top":    "1em",
		"margin-bottom": "1em",
	},
	"h4": {
		"font-size":     "1em",
		"font-weight":   "700",
		"margin-top":    "1.33em",
		"margin-bottom": "1.33em",
	},
	"h5": {
		"font-size":     "0.83em",
		"font-weight":   "700",
		"margin-top":    "1.67em",
		"margin-bottom": "1.67em",
	},
	"h6": {
		"font-size":     "0.67em",
		"font-weight":   "700",
		"margin-top":    "2.33em",
		"margin-bottom": "2.33em",
	},
	"strong": {
		"font-weight": "700",
	},
	"b": {
		"font-weight": "700",
	},
	"em": {
		"font-style": "italic",
	},
	"i": {
		"font-style": "italic",
	},
	"u": {
		"text-decoration-line": "underline",
	},
	"s": {
		"text-decoration-line": "line-through",
	},
	"code": {
		"font-family": "monospace",
	},
	"pre": {
		"font-family": "monospace",
	},
	"small": {
		"font-size": "0.83em",
	},
	"big": {
		"font-size": "1.17em",
	},
	"mark": {
		"background-color": "yellow",
	},
	"hr": {
		"border-top-width": "1px",
		"border-top-style": "solid",
		"width":            "100%",
	},
	"a": {
		"color":                "#0000ee",
		"text-decoration-line": "underline",
	},
	"blockquote": {
		"margin-left":  "40px",
		"margin-right": "40px",
	},
	"center": {
		"text-align": "center",
	},
	"ul": {
		"padding-left": "40px",
	},
	"ol": {
		"padding-left": "40px",
	},
	"li": {
		"display": "flex",
	},
	"section":    {},
	"article":    {},
	"header":     {},
	"footer":     {},
	"nav":        {},
	"main":       {},
	"aside":      {},
	"figure":     {},
	"figcaption": {},
	"details":    {},
	"summary": {
		"font-weight": "700",
	},
	"sup": {
		"font-size":      "0.83em",
		"vertical-align": "super",
	},
	"sub": {
		"font-size":      "0.83em",
		"vertical-align": "sub",
	},
	"del": {
		"text-decoration-line": "line-through",
	},
	"ins": {
		"text-decoration-line": "underline",
	},
	"abbr": {
		"text-decoration-line": "underline",
	},
}

// Resolve computes styles for every node in the tree, applying inheritance and Tailwind classes.
// Resolve walks the node tree and produces computed styles for each node.
func Resolve(root *parse.Node, viewportWidth, viewportHeight float64) map[*parse.Node]*ComputedStyle {
	result := make(map[*parse.Node]*ComputedStyle)
	resolveNode(root, nil, nil, defaultRootFontSize, viewportWidth, viewportHeight, result)
	return result
}

func resolveNode(node *parse.Node, parent *ComputedStyle, parentVars map[string]string, rootFontSize, viewportWidth, viewportHeight float64, result map[*parse.Node]*ComputedStyle) {
	props := mergeProps(node)
	vars, resolved := resolveVariables(props, parentVars)
	cs := resolveStyle(resolved, parent, rootFontSize, viewportWidth, viewportHeight)
	result[node] = cs

	for _, child := range node.Children {
		resolveNode(child, cs, vars, rootFontSize, viewportWidth, viewportHeight, result)
	}
}

func resolveVariables(props map[string]string, parentVars map[string]string) (map[string]string, map[string]string) {
	vars := make(map[string]string)
	for k, v := range parentVars {
		vars[k] = v
	}

	for k, v := range props {
		if strings.HasPrefix(k, "--") {
			vars[k] = v
		}
	}

	resolved := make(map[string]string, len(props))
	for k, v := range props {
		resolved[k] = resolveVar(v, vars)
	}
	return vars, resolved
}

func resolveVar(value string, vars map[string]string) string {
	for {
		idx := strings.Index(value, "var(")
		if idx == -1 {
			return value
		}
		depth := 0
		end := -1
		for i := idx + 4; i < len(value); i++ {
			switch value[i] {
			case '(':
				depth++
			case ')':
				if depth == 0 {
					end = i
				} else {
					depth--
				}
			}
			if end != -1 {
				break
			}
		}
		if end == -1 {
			return value
		}

		inner := strings.TrimSpace(value[idx+4 : end])
		var name, fallback string
		commaIdx := strings.Index(inner, ",")
		if commaIdx != -1 {
			name = strings.TrimSpace(inner[:commaIdx])
			fallback = strings.TrimSpace(inner[commaIdx+1:])
		} else {
			name = inner
		}

		replacement := fallback
		if v, ok := vars[name]; ok {
			replacement = v
		}

		value = value[:idx] + replacement + value[end+1:]
	}
}

func mergeProps(node *parse.Node) map[string]string {
	props := make(map[string]string)

	props["display"] = "flex"
	props["position"] = "relative"
	props["box-sizing"] = "border-box"

	if defaults, ok := tagDefaults[node.Tag]; ok {
		for k, v := range defaults {
			props[k] = v
		}
	}

	if len(node.Classes) > 0 {
		tw := ResolveTailwind(node.Classes)
		for k, v := range tw {
			props[k] = v
		}
	}

	expanded := ExpandShorthands(node.Style)
	for k, v := range expanded {
		props[k] = v
	}

	return props
}

func resolveStyle(props map[string]string, parent *ComputedStyle, rootFontSize, viewportWidth, viewportHeight float64) *ComputedStyle {
	cs := NewComputedStyle()

	if parent == nil {
		cs.FontSize = defaultRootFontSize
		cs.FontWeight = 400
		cs.FontFamily = "sans-serif"
		cs.Color = Color{0, 0, 0, 1}
		cs.LineHeight = 1.2 * defaultRootFontSize
	} else {
		for prop := range inheritedProperties {
			if _, explicit := props[prop]; !explicit {
				inheritProperty(cs, parent, prop)
			}
		}
		if cs.FontSize == 0 {
			cs.FontSize = parent.FontSize
		}
		if cs.FontWeight == 0 {
			cs.FontWeight = parent.FontWeight
		}
		if cs.FontFamily == "" {
			cs.FontFamily = parent.FontFamily
		}
		if cs.LineHeight == 0 {
			cs.LineHeight = parent.LineHeight
		}
	}

	parentFontSize := cs.FontSize
	if parent != nil {
		parentFontSize = parent.FontSize
	}

	if v, ok := props["font-size"]; ok {
		cs.FontSize = resolveFontSize(v, parentFontSize, rootFontSize, viewportWidth, viewportHeight)
	}

	ctx := ResolveContext{
		ParentFontSize: cs.FontSize,
		RootFontSize:   rootFontSize,
		ViewportWidth:  viewportWidth,
		ViewportHeight: viewportHeight,
	}

	if v, ok := props["line-height"]; ok {
		cs.LineHeight = resolveLineHeight(v, cs.FontSize, ctx)
	} else {
		cs.LineHeight = 1.2 * cs.FontSize
	}

	cs.Display = ParseDisplay(getOr(props, "display", "flex"))
	cs.Position = ParsePosition(getOr(props, "position", "relative"))
	cs.BoxSizing = ParseBoxSizing(getOr(props, "box-sizing", "border-box"))

	cs.FlexDirection = ParseFlexDirection(getOr(props, "flex-direction", "row"))
	cs.FlexWrap = ParseFlexWrap(getOr(props, "flex-wrap", "nowrap"))
	cs.AlignItems = ParseAlignItems(getOr(props, "align-items", ""))
	cs.AlignSelf = ParseAlignSelf(getOr(props, "align-self", ""))
	cs.AlignContent = ParseAlignContent(getOr(props, "align-content", ""))
	cs.JustifyContent = ParseJustifyContent(getOr(props, "justify-content", ""))

	if v, ok := props["flex-grow"]; ok {
		cs.FlexGrow = parseFloat(v)
	}
	if v, ok := props["flex-shrink"]; ok {
		cs.FlexShrink = parseFloat(v)
	}
	cs.FlexBasis = resolveValue(props, "flex-basis", ctx)

	cs.Width = resolveValue(props, "width", ctx)
	cs.Height = resolveValue(props, "height", ctx)
	cs.MinWidth = resolveValue(props, "min-width", ctx)
	cs.MinHeight = resolveValue(props, "min-height", ctx)
	cs.MaxWidth = resolveValue(props, "max-width", ctx)
	cs.MaxHeight = resolveValue(props, "max-height", ctx)

	if v, ok := props["aspect-ratio"]; ok {
		cs.AspectRatio = parseAspectRatio(v)
	}

	cs.Top = resolveValue(props, "top", ctx)
	cs.Right = resolveValue(props, "right", ctx)
	cs.Bottom = resolveValue(props, "bottom", ctx)
	cs.Left = resolveValue(props, "left", ctx)

	cs.MarginTop = resolveValue(props, "margin-top", ctx)
	cs.MarginRight = resolveValue(props, "margin-right", ctx)
	cs.MarginBottom = resolveValue(props, "margin-bottom", ctx)
	cs.MarginLeft = resolveValue(props, "margin-left", ctx)

	cs.PaddingTop = resolveToFloat(props, "padding-top", ctx)
	cs.PaddingRight = resolveToFloat(props, "padding-right", ctx)
	cs.PaddingBottom = resolveToFloat(props, "padding-bottom", ctx)
	cs.PaddingLeft = resolveToFloat(props, "padding-left", ctx)

	cs.BorderTopWidth = resolveToFloat(props, "border-top-width", ctx)
	cs.BorderRightWidth = resolveToFloat(props, "border-right-width", ctx)
	cs.BorderBottomWidth = resolveToFloat(props, "border-bottom-width", ctx)
	cs.BorderLeftWidth = resolveToFloat(props, "border-left-width", ctx)

	cs.BorderTopStyle = ParseBorderStyle(getOr(props, "border-top-style", ""))
	cs.BorderRightStyle = ParseBorderStyle(getOr(props, "border-right-style", ""))
	cs.BorderBottomStyle = ParseBorderStyle(getOr(props, "border-bottom-style", ""))
	cs.BorderLeftStyle = ParseBorderStyle(getOr(props, "border-left-style", ""))

	cs.BorderTopColor = resolveColor(props, "border-top-color", cs.Color)
	cs.BorderRightColor = resolveColor(props, "border-right-color", cs.Color)
	cs.BorderBottomColor = resolveColor(props, "border-bottom-color", cs.Color)
	cs.BorderLeftColor = resolveColor(props, "border-left-color", cs.Color)

	cs.BorderTopLeftRadius = resolveToFloat(props, "border-top-left-radius", ctx)
	cs.BorderTopRightRadius = resolveToFloat(props, "border-top-right-radius", ctx)
	cs.BorderBottomLeftRadius = resolveToFloat(props, "border-bottom-left-radius", ctx)
	cs.BorderBottomRightRadius = resolveToFloat(props, "border-bottom-right-radius", ctx)

	cs.BackgroundColor = resolveColor(props, "background-color", Color{0, 0, 0, 0})
	cs.BackgroundImage = getOr(props, "background-image", "")
	cs.BackgroundSize = getOr(props, "background-size", "")
	cs.BackgroundPosition = getOr(props, "background-position", "")
	cs.BackgroundRepeat = getOr(props, "background-repeat", "")

	if v, ok := props["font-family"]; ok {
		cs.FontFamily = strings.Trim(v, "\"'")
	}
	if v, ok := props["font-weight"]; ok {
		cs.FontWeight = parseFontWeight(v)
	}
	if v, ok := props["font-style"]; ok {
		cs.FontStyle = v
	}

	if v, ok := props["color"]; ok {
		c, err := ParseColor(v)
		if err == nil {
			cs.Color = c
		}
	}

	if v, ok := props["letter-spacing"]; ok {
		cs.LetterSpacing = ParseValue(v).Resolve(ctx)
	}

	cs.Direction = getOr(props, "direction", "")
	cs.TextAlign = ParseTextAlign(getOr(props, "text-align", ""))
	cs.TextTransform = ParseTextTransform(getOr(props, "text-transform", ""))
	cs.TextDecorationLine = ParseTextDecorationLine(getOr(props, "text-decoration-line", ""))
	cs.TextDecorationColor = resolveColor(props, "text-decoration-color", cs.Color)
	cs.TextDecorationStyle = getOr(props, "text-decoration-style", "")
	cs.WhiteSpace = ParseWhiteSpace(getOr(props, "white-space", ""))
	cs.WordBreak = ParseWordBreak(getOr(props, "word-break", ""))
	cs.TextOverflow = getOr(props, "text-overflow", "")
	cs.TextShadow = getOr(props, "text-shadow", "")

	if v, ok := props["line-clamp"]; ok {
		cs.LineClamp = parseInt(v)
	} else if v, ok := props["-webkit-line-clamp"]; ok {
		cs.LineClamp = parseInt(v)
	}

	if v, ok := props["opacity"]; ok {
		cs.Opacity = parseFloat(v)
	}
	cs.Overflow = ParseOverflow(getOr(props, "overflow", getOr(props, "overflow-x", "")))
	cs.BoxShadow = getOr(props, "box-shadow", "")
	cs.Transform = getOr(props, "transform", "")
	cs.TransformOrigin = getOr(props, "transform-origin", "")
	cs.ObjectFit = ParseObjectFit(getOr(props, "object-fit", ""))
	cs.ObjectPosition = getOr(props, "object-position", "")
	cs.Filter = getOr(props, "filter", "")
	cs.ClipPath = getOr(props, "clip-path", "")

	if v, ok := props["gap"]; ok {
		g := ParseValue(v).Resolve(ctx)
		cs.Gap = g
		cs.RowGap = g
		cs.ColumnGap = g
	}
	if v, ok := props["row-gap"]; ok {
		cs.RowGap = ParseValue(v).Resolve(ctx)
	}
	if v, ok := props["column-gap"]; ok {
		cs.ColumnGap = ParseValue(v).Resolve(ctx)
	}

	return cs
}

func inheritProperty(cs *ComputedStyle, parent *ComputedStyle, prop string) {
	switch prop {
	case "color":
		cs.Color = parent.Color
	case "font-family":
		cs.FontFamily = parent.FontFamily
	case "font-size":
		cs.FontSize = parent.FontSize
	case "font-weight":
		cs.FontWeight = parent.FontWeight
	case "font-style":
		cs.FontStyle = parent.FontStyle
	case "line-height":
		cs.LineHeight = parent.LineHeight
	case "letter-spacing":
		cs.LetterSpacing = parent.LetterSpacing
	case "text-align":
		cs.Direction = parent.Direction
		cs.TextAlign = parent.TextAlign
	case "text-transform":
		cs.TextTransform = parent.TextTransform
	case "text-decoration-line":
		cs.TextDecorationLine = parent.TextDecorationLine
	case "text-decoration-color":
		cs.TextDecorationColor = parent.TextDecorationColor
	case "text-decoration-style":
		cs.TextDecorationStyle = parent.TextDecorationStyle
	case "white-space":
		cs.WhiteSpace = parent.WhiteSpace
	case "word-break":
		cs.WordBreak = parent.WordBreak
	case "text-shadow":
		cs.TextShadow = parent.TextShadow
	}
}

func resolveValue(props map[string]string, prop string, ctx ResolveContext) Value {
	v, ok := props[prop]
	if !ok {
		return Value{Unit: UnitNone}
	}
	parsed := ParseValue(v)
	if parsed.Unit == UnitAuto || parsed.Unit == UnitNone {
		return parsed
	}
	if parsed.Unit == UnitPercent {
		return parsed
	}
	resolved := parsed.Resolve(ctx)
	return Value{Raw: resolved, Unit: UnitPx}
}

func resolveToFloat(props map[string]string, prop string, ctx ResolveContext) float64 {
	v, ok := props[prop]
	if !ok {
		return 0
	}
	return ParseValue(v).Resolve(ctx)
}

func resolveColor(props map[string]string, prop string, fallback Color) Color {
	v, ok := props[prop]
	if !ok {
		return fallback
	}
	c, err := ParseColor(v)
	if err != nil {
		return fallback
	}
	if c.A == -1 {
		return fallback
	}
	return c
}

func resolveFontSize(v string, parentFontSize, rootFontSize, viewportWidth, viewportHeight float64) float64 {
	parsed := ParseValue(v)
	ctx := ResolveContext{
		ParentFontSize: parentFontSize,
		RootFontSize:   rootFontSize,
		ViewportWidth:  viewportWidth,
		ViewportHeight: viewportHeight,
	}
	return parsed.Resolve(ctx)
}

func resolveLineHeight(v string, fontSize float64, ctx ResolveContext) float64 {
	if f, err := strconv.ParseFloat(v, 64); err == nil {
		if !strings.ContainsAny(v, "pxemr%vwh") {
			return f * fontSize
		}
	}
	return ParseValue(v).Resolve(ctx)
}

func getOr(props map[string]string, key, fallback string) string {
	if v, ok := props[key]; ok {
		return v
	}
	return fallback
}

func parseFloat(s string) float64 {
	f, _ := strconv.ParseFloat(strings.TrimSpace(s), 64)
	return f
}

func parseFontWeight(s string) int {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "bold":
		return 700
	case "bolder":
		return 700
	case "lighter":
		return 300
	case "normal":
		return 400
	}
	w, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil {
		return 400
	}
	return w
}

func parseInt(s string) int {
	v, _ := strconv.Atoi(strings.TrimSpace(s))
	return v
}

func parseAspectRatio(s string) float64 {
	s = strings.TrimSpace(s)
	if strings.Contains(s, "/") {
		parts := strings.SplitN(s, "/", 2)
		num := parseFloat(parts[0])
		den := parseFloat(parts[1])
		if den != 0 {
			return num / den
		}
		return 0
	}
	return parseFloat(s)
}

