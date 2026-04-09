package style

// Display represents the CSS display property.
type Display int

const (
	DisplayFlex Display = iota
	DisplayNone
	DisplayBlock
	DisplayContents
)

func (d Display) String() string {
	switch d {
	case DisplayFlex:
		return "flex"
	case DisplayNone:
		return "none"
	case DisplayBlock:
		return "block"
	case DisplayContents:
		return "contents"
	}
	return "flex"
}

// ParseDisplay parses a CSS display value string.
func ParseDisplay(s string) Display {
	switch s {
	case "flex":
		return DisplayFlex
	case "none":
		return DisplayNone
	case "block":
		return DisplayBlock
	case "contents":
		return DisplayContents
	}
	return DisplayFlex
}

// Position represents the CSS position property.
type Position int

const (
	PositionStatic Position = iota
	PositionRelative
	PositionAbsolute
)

func (p Position) String() string {
	switch p {
	case PositionStatic:
		return "static"
	case PositionRelative:
		return "relative"
	case PositionAbsolute:
		return "absolute"
	}
	return "static"
}

// ParsePosition parses a CSS position value string.
func ParsePosition(s string) Position {
	switch s {
	case "static":
		return PositionStatic
	case "relative":
		return PositionRelative
	case "absolute":
		return PositionAbsolute
	}
	return PositionStatic
}

// FlexDirection represents the CSS flex-direction property.
type FlexDirection int

const (
	FlexDirectionRow FlexDirection = iota
	FlexDirectionRowReverse
	FlexDirectionColumn
	FlexDirectionColumnReverse
)

func (f FlexDirection) String() string {
	switch f {
	case FlexDirectionRow:
		return "row"
	case FlexDirectionRowReverse:
		return "row-reverse"
	case FlexDirectionColumn:
		return "column"
	case FlexDirectionColumnReverse:
		return "column-reverse"
	}
	return "row"
}

// ParseFlexDirection parses a CSS flex-direction value string.
func ParseFlexDirection(s string) FlexDirection {
	switch s {
	case "row":
		return FlexDirectionRow
	case "row-reverse":
		return FlexDirectionRowReverse
	case "column":
		return FlexDirectionColumn
	case "column-reverse":
		return FlexDirectionColumnReverse
	}
	return FlexDirectionRow
}

// FlexWrap represents the CSS flex-wrap property.
type FlexWrap int

const (
	FlexWrapNoWrap FlexWrap = iota
	FlexWrapWrap
	FlexWrapWrapReverse
)

func (f FlexWrap) String() string {
	switch f {
	case FlexWrapNoWrap:
		return "nowrap"
	case FlexWrapWrap:
		return "wrap"
	case FlexWrapWrapReverse:
		return "wrap-reverse"
	}
	return "nowrap"
}

// ParseFlexWrap parses a CSS flex-wrap value string.
func ParseFlexWrap(s string) FlexWrap {
	switch s {
	case "nowrap":
		return FlexWrapNoWrap
	case "wrap":
		return FlexWrapWrap
	case "wrap-reverse":
		return FlexWrapWrapReverse
	}
	return FlexWrapNoWrap
}

// AlignItems represents the CSS align-items property.
type AlignItems int

const (
	AlignItemsAuto AlignItems = iota
	AlignItemsFlexStart
	AlignItemsFlexEnd
	AlignItemsCenter
	AlignItemsStretch
	AlignItemsBaseline
	AlignItemsSpaceBetween
	AlignItemsSpaceAround
)

func (a AlignItems) String() string {
	switch a {
	case AlignItemsAuto:
		return "auto"
	case AlignItemsFlexStart:
		return "flex-start"
	case AlignItemsFlexEnd:
		return "flex-end"
	case AlignItemsCenter:
		return "center"
	case AlignItemsStretch:
		return "stretch"
	case AlignItemsBaseline:
		return "baseline"
	case AlignItemsSpaceBetween:
		return "space-between"
	case AlignItemsSpaceAround:
		return "space-around"
	}
	return "auto"
}

// ParseAlignItems parses a CSS align-items value string.
func ParseAlignItems(s string) AlignItems {
	switch s {
	case "auto":
		return AlignItemsAuto
	case "flex-start":
		return AlignItemsFlexStart
	case "flex-end":
		return AlignItemsFlexEnd
	case "center":
		return AlignItemsCenter
	case "stretch":
		return AlignItemsStretch
	case "baseline":
		return AlignItemsBaseline
	case "space-between":
		return AlignItemsSpaceBetween
	case "space-around":
		return AlignItemsSpaceAround
	}
	return AlignItemsAuto
}

// AlignSelf represents the CSS align-self property.
type AlignSelf int

const (
	AlignSelfAuto AlignSelf = iota
	AlignSelfFlexStart
	AlignSelfFlexEnd
	AlignSelfCenter
	AlignSelfStretch
	AlignSelfBaseline
	AlignSelfSpaceBetween
	AlignSelfSpaceAround
)

func (a AlignSelf) String() string {
	switch a {
	case AlignSelfAuto:
		return "auto"
	case AlignSelfFlexStart:
		return "flex-start"
	case AlignSelfFlexEnd:
		return "flex-end"
	case AlignSelfCenter:
		return "center"
	case AlignSelfStretch:
		return "stretch"
	case AlignSelfBaseline:
		return "baseline"
	case AlignSelfSpaceBetween:
		return "space-between"
	case AlignSelfSpaceAround:
		return "space-around"
	}
	return "auto"
}

// ParseAlignSelf parses a CSS align-self value string.
func ParseAlignSelf(s string) AlignSelf {
	switch s {
	case "auto":
		return AlignSelfAuto
	case "flex-start":
		return AlignSelfFlexStart
	case "flex-end":
		return AlignSelfFlexEnd
	case "center":
		return AlignSelfCenter
	case "stretch":
		return AlignSelfStretch
	case "baseline":
		return AlignSelfBaseline
	case "space-between":
		return AlignSelfSpaceBetween
	case "space-around":
		return AlignSelfSpaceAround
	}
	return AlignSelfAuto
}

// AlignContent represents the CSS align-content property.
type AlignContent int

const (
	AlignContentAuto AlignContent = iota
	AlignContentFlexStart
	AlignContentFlexEnd
	AlignContentCenter
	AlignContentStretch
	AlignContentBaseline
	AlignContentSpaceBetween
	AlignContentSpaceAround
)

func (a AlignContent) String() string {
	switch a {
	case AlignContentAuto:
		return "auto"
	case AlignContentFlexStart:
		return "flex-start"
	case AlignContentFlexEnd:
		return "flex-end"
	case AlignContentCenter:
		return "center"
	case AlignContentStretch:
		return "stretch"
	case AlignContentBaseline:
		return "baseline"
	case AlignContentSpaceBetween:
		return "space-between"
	case AlignContentSpaceAround:
		return "space-around"
	}
	return "auto"
}

// ParseAlignContent parses a CSS align-content value string.
func ParseAlignContent(s string) AlignContent {
	switch s {
	case "auto":
		return AlignContentAuto
	case "flex-start":
		return AlignContentFlexStart
	case "flex-end":
		return AlignContentFlexEnd
	case "center":
		return AlignContentCenter
	case "stretch":
		return AlignContentStretch
	case "baseline":
		return AlignContentBaseline
	case "space-between":
		return AlignContentSpaceBetween
	case "space-around":
		return AlignContentSpaceAround
	}
	return AlignContentAuto
}

// JustifyContent represents the CSS justify-content property.
type JustifyContent int

const (
	JustifyContentFlexStart JustifyContent = iota
	JustifyContentFlexEnd
	JustifyContentCenter
	JustifyContentSpaceBetween
	JustifyContentSpaceAround
	JustifyContentSpaceEvenly
)

func (j JustifyContent) String() string {
	switch j {
	case JustifyContentFlexStart:
		return "flex-start"
	case JustifyContentFlexEnd:
		return "flex-end"
	case JustifyContentCenter:
		return "center"
	case JustifyContentSpaceBetween:
		return "space-between"
	case JustifyContentSpaceAround:
		return "space-around"
	case JustifyContentSpaceEvenly:
		return "space-evenly"
	}
	return "flex-start"
}

// ParseJustifyContent parses a CSS justify-content value string.
func ParseJustifyContent(s string) JustifyContent {
	switch s {
	case "flex-start":
		return JustifyContentFlexStart
	case "flex-end":
		return JustifyContentFlexEnd
	case "center":
		return JustifyContentCenter
	case "space-between":
		return JustifyContentSpaceBetween
	case "space-around":
		return JustifyContentSpaceAround
	case "space-evenly":
		return JustifyContentSpaceEvenly
	}
	return JustifyContentFlexStart
}

// TextAlign represents the CSS text-align property.
type TextAlign int

const (
	TextAlignLeft TextAlign = iota
	TextAlignRight
	TextAlignCenter
	TextAlignJustify
	TextAlignStart
	TextAlignEnd
)

func (t TextAlign) String() string {
	switch t {
	case TextAlignLeft:
		return "left"
	case TextAlignRight:
		return "right"
	case TextAlignCenter:
		return "center"
	case TextAlignJustify:
		return "justify"
	case TextAlignStart:
		return "start"
	case TextAlignEnd:
		return "end"
	}
	return "left"
}

// ParseTextAlign parses a CSS text-align value string.
func ParseTextAlign(s string) TextAlign {
	switch s {
	case "left":
		return TextAlignLeft
	case "right":
		return TextAlignRight
	case "center":
		return TextAlignCenter
	case "justify":
		return TextAlignJustify
	case "start":
		return TextAlignStart
	case "end":
		return TextAlignEnd
	}
	return TextAlignLeft
}

// WhiteSpace represents the CSS white-space property.
type WhiteSpace int

const (
	WhiteSpaceNormal WhiteSpace = iota
	WhiteSpaceNoWrap
	WhiteSpacePre
	WhiteSpacePreWrap
	WhiteSpacePreLine
)

func (w WhiteSpace) String() string {
	switch w {
	case WhiteSpaceNormal:
		return "normal"
	case WhiteSpaceNoWrap:
		return "nowrap"
	case WhiteSpacePre:
		return "pre"
	case WhiteSpacePreWrap:
		return "pre-wrap"
	case WhiteSpacePreLine:
		return "pre-line"
	}
	return "normal"
}

// ParseWhiteSpace parses a CSS white-space value string.
func ParseWhiteSpace(s string) WhiteSpace {
	switch s {
	case "normal":
		return WhiteSpaceNormal
	case "nowrap":
		return WhiteSpaceNoWrap
	case "pre":
		return WhiteSpacePre
	case "pre-wrap":
		return WhiteSpacePreWrap
	case "pre-line":
		return WhiteSpacePreLine
	}
	return WhiteSpaceNormal
}

// WordBreak represents the CSS word-break property.
type WordBreak int

const (
	WordBreakNormal WordBreak = iota
	WordBreakAll
	WordBreakWord
	WordBreakKeepAll
)

func (w WordBreak) String() string {
	switch w {
	case WordBreakNormal:
		return "normal"
	case WordBreakAll:
		return "break-all"
	case WordBreakWord:
		return "break-word"
	case WordBreakKeepAll:
		return "keep-all"
	}
	return "normal"
}

// ParseWordBreak parses a CSS word-break value string.
func ParseWordBreak(s string) WordBreak {
	switch s {
	case "normal":
		return WordBreakNormal
	case "break-all":
		return WordBreakAll
	case "break-word":
		return WordBreakWord
	case "keep-all":
		return WordBreakKeepAll
	}
	return WordBreakNormal
}

// Overflow represents the CSS overflow property.
type Overflow int

const (
	OverflowVisible Overflow = iota
	OverflowHidden
)

func (o Overflow) String() string {
	switch o {
	case OverflowVisible:
		return "visible"
	case OverflowHidden:
		return "hidden"
	}
	return "visible"
}

// ParseOverflow parses a CSS overflow value string.
func ParseOverflow(s string) Overflow {
	switch s {
	case "visible":
		return OverflowVisible
	case "hidden":
		return OverflowHidden
	}
	return OverflowVisible
}

// ObjectFit represents the CSS object-fit property.
type ObjectFit int

const (
	ObjectFitFill ObjectFit = iota
	ObjectFitContain
	ObjectFitCover
	ObjectFitScaleDown
	ObjectFitNone
)

func (o ObjectFit) String() string {
	switch o {
	case ObjectFitFill:
		return "fill"
	case ObjectFitContain:
		return "contain"
	case ObjectFitCover:
		return "cover"
	case ObjectFitScaleDown:
		return "scale-down"
	case ObjectFitNone:
		return "none"
	}
	return "fill"
}

// ParseObjectFit parses a CSS object-fit value string.
func ParseObjectFit(s string) ObjectFit {
	switch s {
	case "fill":
		return ObjectFitFill
	case "contain":
		return ObjectFitContain
	case "cover":
		return ObjectFitCover
	case "scale-down":
		return ObjectFitScaleDown
	case "none":
		return ObjectFitNone
	}
	return ObjectFitFill
}

// BorderStyle represents the CSS border-style property.
type BorderStyle int

const (
	BorderStyleNone BorderStyle = iota
	BorderStyleSolid
	BorderStyleDashed
	BorderStyleDotted
	BorderStyleDouble
)

func (b BorderStyle) String() string {
	switch b {
	case BorderStyleNone:
		return "none"
	case BorderStyleSolid:
		return "solid"
	case BorderStyleDashed:
		return "dashed"
	case BorderStyleDotted:
		return "dotted"
	case BorderStyleDouble:
		return "double"
	}
	return "none"
}

// ParseBorderStyle parses a CSS border-style value string.
func ParseBorderStyle(s string) BorderStyle {
	switch s {
	case "none":
		return BorderStyleNone
	case "solid":
		return BorderStyleSolid
	case "dashed":
		return BorderStyleDashed
	case "dotted":
		return BorderStyleDotted
	case "double":
		return BorderStyleDouble
	}
	return BorderStyleNone
}

// TextTransform represents the CSS text-transform property.
type TextTransform int

const (
	TextTransformNone TextTransform = iota
	TextTransformUppercase
	TextTransformLowercase
	TextTransformCapitalize
)

func (t TextTransform) String() string {
	switch t {
	case TextTransformNone:
		return "none"
	case TextTransformUppercase:
		return "uppercase"
	case TextTransformLowercase:
		return "lowercase"
	case TextTransformCapitalize:
		return "capitalize"
	}
	return "none"
}

// ParseTextTransform parses a CSS text-transform value string.
func ParseTextTransform(s string) TextTransform {
	switch s {
	case "none":
		return TextTransformNone
	case "uppercase":
		return TextTransformUppercase
	case "lowercase":
		return TextTransformLowercase
	case "capitalize":
		return TextTransformCapitalize
	}
	return TextTransformNone
}

// TextDecorationLine represents the CSS text-decoration-line property.
type TextDecorationLine int

const (
	TextDecorationNone TextDecorationLine = iota
	TextDecorationUnderline
	TextDecorationOverline
	TextDecorationLineThrough
)

func (t TextDecorationLine) String() string {
	switch t {
	case TextDecorationNone:
		return "none"
	case TextDecorationUnderline:
		return "underline"
	case TextDecorationOverline:
		return "overline"
	case TextDecorationLineThrough:
		return "line-through"
	}
	return "none"
}

// ParseTextDecorationLine parses a CSS text-decoration-line value string.
func ParseTextDecorationLine(s string) TextDecorationLine {
	switch s {
	case "none":
		return TextDecorationNone
	case "underline":
		return TextDecorationUnderline
	case "overline":
		return TextDecorationOverline
	case "line-through":
		return TextDecorationLineThrough
	}
	return TextDecorationNone
}

// BoxSizing represents the CSS box-sizing property.
type BoxSizing int

const (
	BoxSizingContentBox BoxSizing = iota
	BoxSizingBorderBox
)

func (b BoxSizing) String() string {
	switch b {
	case BoxSizingContentBox:
		return "content-box"
	case BoxSizingBorderBox:
		return "border-box"
	}
	return "content-box"
}

// ParseBoxSizing parses a CSS box-sizing value string.
func ParseBoxSizing(s string) BoxSizing {
	switch s {
	case "content-box":
		return BoxSizingContentBox
	case "border-box":
		return BoxSizingBorderBox
	}
	return BoxSizingContentBox
}
