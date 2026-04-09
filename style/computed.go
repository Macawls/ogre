package style

// ComputedStyle holds all resolved CSS properties for a single node.
type ComputedStyle struct {
	Display  Display
	Position Position
	Top      Value
	Right    Value
	Bottom   Value
	Left     Value

	FlexDirection  FlexDirection
	FlexWrap       FlexWrap
	FlexGrow       float64
	FlexShrink     float64
	FlexBasis      Value
	AlignItems     AlignItems
	AlignSelf      AlignSelf
	AlignContent   AlignContent
	JustifyContent JustifyContent
	Gap            float64
	RowGap         float64
	ColumnGap      float64

	Width     Value
	Height    Value
	MinWidth  Value
	MinHeight Value
	MaxWidth  Value
	MaxHeight Value
	AspectRatio float64

	MarginTop    Value
	MarginRight  Value
	MarginBottom Value
	MarginLeft   Value
	PaddingTop    float64
	PaddingRight  float64
	PaddingBottom float64
	PaddingLeft   float64
	BorderTopWidth    float64
	BorderRightWidth  float64
	BorderBottomWidth float64
	BorderLeftWidth   float64

	BorderTopStyle    BorderStyle
	BorderRightStyle  BorderStyle
	BorderBottomStyle BorderStyle
	BorderLeftStyle   BorderStyle
	BorderTopColor    Color
	BorderRightColor  Color
	BorderBottomColor Color
	BorderLeftColor   Color
	BorderTopLeftRadius     float64
	BorderTopRightRadius    float64
	BorderBottomLeftRadius  float64
	BorderBottomRightRadius float64

	BackgroundColor    Color
	BackgroundImage    string
	BackgroundSize     string
	BackgroundPosition string
	BackgroundRepeat   string

	FontFamily          string
	FontSize            float64
	FontWeight          int
	FontStyle           string
	Color               Color
	LineHeight          float64
	LetterSpacing       float64
	Direction           string
	TextAlign           TextAlign
	TextTransform       TextTransform
	TextDecorationLine  TextDecorationLine
	TextDecorationColor Color
	TextDecorationStyle string
	WhiteSpace          WhiteSpace
	WordBreak           WordBreak
	TextOverflow        string
	TextShadow          string
	LineClamp           int

	Opacity         float64
	Overflow        Overflow
	BoxShadow       string
	Transform       string
	TransformOrigin string
	ObjectFit       ObjectFit
	ObjectPosition  string
	Filter          string
	ClipPath        string

	BoxSizing BoxSizing
}

// NewComputedStyle returns a ComputedStyle with default values.
func NewComputedStyle() *ComputedStyle {
	none := Value{Unit: UnitNone}
	return &ComputedStyle{
		Width: none, Height: none,
		MinWidth: none, MinHeight: none,
		MaxWidth: none, MaxHeight: none,
		Top: none, Right: none, Bottom: none, Left: none,
		MarginTop: none, MarginRight: none, MarginBottom: none, MarginLeft: none,
		FlexBasis: none,
		Opacity:   1,
		FlexShrink: 1,
	}
}
