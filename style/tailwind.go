package style

import (
	"fmt"
	"strconv"
	"strings"
)

var tailwindColors = map[string]map[int]string{
	"slate": {
		50: "#f8fafc", 100: "#f1f5f9", 200: "#e2e8f0", 300: "#cbd5e1",
		400: "#94a3b8", 500: "#64748b", 600: "#475569", 700: "#334155",
		800: "#1e293b", 900: "#0f172a", 950: "#020617",
	},
	"gray": {
		50: "#f9fafb", 100: "#f3f4f6", 200: "#e5e7eb", 300: "#d1d5db",
		400: "#9ca3af", 500: "#6b7280", 600: "#4b5563", 700: "#374151",
		800: "#1f2937", 900: "#111827", 950: "#030712",
	},
	"zinc": {
		50: "#fafafa", 100: "#f4f4f5", 200: "#e4e4e7", 300: "#d4d4d8",
		400: "#a1a1aa", 500: "#71717a", 600: "#52525b", 700: "#3f3f46",
		800: "#27272a", 900: "#18181b", 950: "#09090b",
	},
	"neutral": {
		50: "#fafafa", 100: "#f5f5f5", 200: "#e5e5e5", 300: "#d4d4d4",
		400: "#a3a3a3", 500: "#737373", 600: "#525252", 700: "#404040",
		800: "#262626", 900: "#171717", 950: "#0a0a0a",
	},
	"stone": {
		50: "#fafaf9", 100: "#f5f5f4", 200: "#e7e5e4", 300: "#d6d3d1",
		400: "#a8a29e", 500: "#78716c", 600: "#57534e", 700: "#44403c",
		800: "#292524", 900: "#1c1917", 950: "#0c0a09",
	},
	"red": {
		50: "#fef2f2", 100: "#fee2e2", 200: "#fecaca", 300: "#fca5a5",
		400: "#f87171", 500: "#ef4444", 600: "#dc2626", 700: "#b91c1c",
		800: "#991b1b", 900: "#7f1d1d", 950: "#450a0a",
	},
	"orange": {
		50: "#fff7ed", 100: "#ffedd5", 200: "#fed7aa", 300: "#fdba74",
		400: "#fb923c", 500: "#f97316", 600: "#ea580c", 700: "#c2410c",
		800: "#9a3412", 900: "#7c2d12", 950: "#431407",
	},
	"amber": {
		50: "#fffbeb", 100: "#fef3c7", 200: "#fde68a", 300: "#fcd34d",
		400: "#fbbf24", 500: "#f59e0b", 600: "#d97706", 700: "#b45309",
		800: "#92400e", 900: "#78350f", 950: "#451a03",
	},
	"yellow": {
		50: "#fefce8", 100: "#fef9c3", 200: "#fef08a", 300: "#fde047",
		400: "#facc15", 500: "#eab308", 600: "#ca8a04", 700: "#a16207",
		800: "#854d0e", 900: "#713f12", 950: "#422006",
	},
	"lime": {
		50: "#f7fee7", 100: "#ecfccb", 200: "#d9f99d", 300: "#bef264",
		400: "#a3e635", 500: "#84cc16", 600: "#65a30d", 700: "#4d7c0f",
		800: "#3f6212", 900: "#365314", 950: "#1a2e05",
	},
	"green": {
		50: "#f0fdf4", 100: "#dcfce7", 200: "#bbf7d0", 300: "#86efac",
		400: "#4ade80", 500: "#22c55e", 600: "#16a34a", 700: "#15803d",
		800: "#166534", 900: "#14532d", 950: "#052e16",
	},
	"emerald": {
		50: "#ecfdf5", 100: "#d1fae5", 200: "#a7f3d0", 300: "#6ee7b7",
		400: "#34d399", 500: "#10b981", 600: "#059669", 700: "#047857",
		800: "#065f46", 900: "#064e3b", 950: "#022c22",
	},
	"teal": {
		50: "#f0fdfa", 100: "#ccfbf1", 200: "#99f6e4", 300: "#5eead4",
		400: "#2dd4bf", 500: "#14b8a6", 600: "#0d9488", 700: "#0f766e",
		800: "#115e59", 900: "#134e4a", 950: "#042f2e",
	},
	"cyan": {
		50: "#ecfeff", 100: "#cffafe", 200: "#a5f3fc", 300: "#67e8f9",
		400: "#22d3ee", 500: "#06b6d4", 600: "#0891b2", 700: "#0e7490",
		800: "#155e75", 900: "#164e63", 950: "#083344",
	},
	"sky": {
		50: "#f0f9ff", 100: "#e0f2fe", 200: "#bae6fd", 300: "#7dd3fc",
		400: "#38bdf8", 500: "#0ea5e9", 600: "#0284c7", 700: "#0369a1",
		800: "#075985", 900: "#0c4a6e", 950: "#082f49",
	},
	"blue": {
		50: "#eff6ff", 100: "#dbeafe", 200: "#bfdbfe", 300: "#93c5fd",
		400: "#60a5fa", 500: "#3b82f6", 600: "#2563eb", 700: "#1d4ed8",
		800: "#1e40af", 900: "#1e3a8a", 950: "#172554",
	},
	"indigo": {
		50: "#eef2ff", 100: "#e0e7ff", 200: "#c7d2fe", 300: "#a5b4fc",
		400: "#818cf8", 500: "#6366f1", 600: "#4f46e5", 700: "#4338ca",
		800: "#3730a3", 900: "#312e81", 950: "#1e1b4b",
	},
	"violet": {
		50: "#f5f3ff", 100: "#ede9fe", 200: "#ddd6fe", 300: "#c4b5fd",
		400: "#a78bfa", 500: "#8b5cf6", 600: "#7c3aed", 700: "#6d28d9",
		800: "#5b21b6", 900: "#4c1d95", 950: "#2e1065",
	},
	"purple": {
		50: "#faf5ff", 100: "#f3e8ff", 200: "#e9d5ff", 300: "#d8b4fe",
		400: "#c084fc", 500: "#a855f7", 600: "#9333ea", 700: "#7e22ce",
		800: "#6b21a8", 900: "#581c87", 950: "#3b0764",
	},
	"fuchsia": {
		50: "#fdf4ff", 100: "#fae8ff", 200: "#f5d0fe", 300: "#f0abfc",
		400: "#e879f9", 500: "#d946ef", 600: "#c026d3", 700: "#a21caf",
		800: "#86198f", 900: "#701a75", 950: "#4a044e",
	},
	"pink": {
		50: "#fdf2f8", 100: "#fce7f3", 200: "#fbcfe8", 300: "#f9a8d4",
		400: "#f472b6", 500: "#ec4899", 600: "#db2777", 700: "#be185d",
		800: "#9d174d", 900: "#831843", 950: "#500724",
	},
	"rose": {
		50: "#fff1f2", 100: "#ffe4e6", 200: "#fecdd3", 300: "#fda4af",
		400: "#fb7185", 500: "#f43f5e", 600: "#e11d48", 700: "#be123c",
		800: "#9f1239", 900: "#881337", 950: "#4c0519",
	},
}

func parseSpacing(s string) (string, bool) {
	switch s {
	case "px":
		return "1px", true
	case "0":
		return "0px", true
	case "0.5":
		return "2px", true
	case "1.5":
		return "6px", true
	case "2.5":
		return "10px", true
	case "3.5":
		return "14px", true
	}
	n, err := strconv.Atoi(s)
	if err != nil || n < 0 || n > 96 {
		return "", false
	}
	return fmt.Sprintf("%dpx", n*4), true
}

func resolveColorClass(prefix, rest string) (string, string, bool) {
	var prop string
	switch prefix {
	case "text":
		prop = "color"
	case "bg":
		prop = "background-color"
	case "border":
		prop = "border-color"
	default:
		return "", "", false
	}

	switch rest {
	case "white":
		return prop, "#ffffff", true
	case "black":
		return prop, "#000000", true
	case "transparent":
		return prop, "transparent", true
	}

	lastDash := strings.LastIndex(rest, "-")
	if lastDash < 0 {
		return "", "", false
	}
	colorName := rest[:lastDash]
	shadeStr := rest[lastDash+1:]
	shade, err := strconv.Atoi(shadeStr)
	if err != nil {
		return "", "", false
	}
	shades, ok := tailwindColors[colorName]
	if !ok {
		return "", "", false
	}
	hex, ok := shades[shade]
	if !ok {
		return "", "", false
	}
	return prop, hex, true
}

// ResolveTailwind converts a list of Tailwind CSS class names into a CSS property map.
// ResolveTailwind converts Tailwind utility classes to CSS properties.
func ResolveTailwind(classes []string) map[string]string {
	result := make(map[string]string)
	for _, cls := range classes {
		resolveTailwindClass(cls, result)
	}
	return result
}

func resolveTailwindClass(cls string, out map[string]string) {
	if strings.Contains(cls, "[") && strings.HasSuffix(cls, "]") {
		resolveArbitrary(cls, out)
		return
	}

	switch cls {
	// Flexbox display & direction
	case "flex":
		out["display"] = "flex"
	case "flex-row":
		out["flex-direction"] = "row"
	case "flex-col":
		out["flex-direction"] = "column"
	case "flex-wrap":
		out["flex-wrap"] = "wrap"
	case "flex-nowrap":
		out["flex-wrap"] = "nowrap"
	case "flex-1":
		out["flex-grow"] = "1"
		out["flex-shrink"] = "1"
		out["flex-basis"] = "0%"
	case "flex-auto":
		out["flex-grow"] = "1"
		out["flex-shrink"] = "1"
		out["flex-basis"] = "auto"
	case "flex-initial":
		out["flex-grow"] = "0"
		out["flex-shrink"] = "1"
		out["flex-basis"] = "auto"
	case "flex-none":
		out["flex-grow"] = "0"
		out["flex-shrink"] = "0"
		out["flex-basis"] = "auto"
	case "flex-grow":
		out["flex-grow"] = "1"
	case "flex-grow-0":
		out["flex-grow"] = "0"
	case "flex-shrink":
		out["flex-shrink"] = "1"
	case "flex-shrink-0":
		out["flex-shrink"] = "0"

	// Align items
	case "items-start":
		out["align-items"] = "flex-start"
	case "items-end":
		out["align-items"] = "flex-end"
	case "items-center":
		out["align-items"] = "center"
	case "items-stretch":
		out["align-items"] = "stretch"
	case "items-baseline":
		out["align-items"] = "baseline"

	// Justify content
	case "justify-start":
		out["justify-content"] = "flex-start"
	case "justify-end":
		out["justify-content"] = "flex-end"
	case "justify-center":
		out["justify-content"] = "center"
	case "justify-between":
		out["justify-content"] = "space-between"
	case "justify-around":
		out["justify-content"] = "space-around"
	case "justify-evenly":
		out["justify-content"] = "space-evenly"

	// Align self
	case "self-auto":
		out["align-self"] = "auto"
	case "self-start":
		out["align-self"] = "flex-start"
	case "self-end":
		out["align-self"] = "flex-end"
	case "self-center":
		out["align-self"] = "center"
	case "self-stretch":
		out["align-self"] = "stretch"

	// Align content
	case "content-start":
		out["align-content"] = "flex-start"
	case "content-end":
		out["align-content"] = "flex-end"
	case "content-center":
		out["align-content"] = "center"
	case "content-between":
		out["align-content"] = "space-between"
	case "content-around":
		out["align-content"] = "space-around"
	case "content-stretch":
		out["align-content"] = "stretch"

	// Typography - font size
	case "text-xs":
		out["font-size"] = "12px"
		out["line-height"] = "16px"
	case "text-sm":
		out["font-size"] = "14px"
		out["line-height"] = "20px"
	case "text-base":
		out["font-size"] = "16px"
		out["line-height"] = "24px"
	case "text-lg":
		out["font-size"] = "18px"
		out["line-height"] = "28px"
	case "text-xl":
		out["font-size"] = "20px"
		out["line-height"] = "28px"
	case "text-2xl":
		out["font-size"] = "24px"
		out["line-height"] = "32px"
	case "text-3xl":
		out["font-size"] = "30px"
		out["line-height"] = "36px"
	case "text-4xl":
		out["font-size"] = "36px"
		out["line-height"] = "40px"
	case "text-5xl":
		out["font-size"] = "48px"
		out["line-height"] = "48px"
	case "text-6xl":
		out["font-size"] = "60px"
		out["line-height"] = "60px"
	case "text-7xl":
		out["font-size"] = "72px"
		out["line-height"] = "72px"
	case "text-8xl":
		out["font-size"] = "96px"
		out["line-height"] = "96px"
	case "text-9xl":
		out["font-size"] = "128px"
		out["line-height"] = "128px"

	// Font weight
	case "font-thin":
		out["font-weight"] = "100"
	case "font-light":
		out["font-weight"] = "300"
	case "font-normal":
		out["font-weight"] = "400"
	case "font-medium":
		out["font-weight"] = "500"
	case "font-semibold":
		out["font-weight"] = "600"
	case "font-bold":
		out["font-weight"] = "700"
	case "font-extrabold":
		out["font-weight"] = "800"
	case "font-black":
		out["font-weight"] = "900"

	// Text alignment
	case "text-left":
		out["text-align"] = "left"
	case "text-center":
		out["text-align"] = "center"
	case "text-right":
		out["text-align"] = "right"
	case "text-justify":
		out["text-align"] = "justify"

	// Font style
	case "italic":
		out["font-style"] = "italic"
	case "not-italic":
		out["font-style"] = "normal"

	// Text transform
	case "uppercase":
		out["text-transform"] = "uppercase"
	case "lowercase":
		out["text-transform"] = "lowercase"
	case "capitalize":
		out["text-transform"] = "capitalize"
	case "normal-case":
		out["text-transform"] = "none"

	// Text decoration
	case "underline":
		out["text-decoration-line"] = "underline"
	case "overline":
		out["text-decoration-line"] = "overline"
	case "line-through":
		out["text-decoration-line"] = "line-through"
	case "no-underline":
		out["text-decoration-line"] = "none"

	// Leading (line-height) named
	case "leading-none":
		out["line-height"] = "1"
	case "leading-tight":
		out["line-height"] = "1.25"
	case "leading-normal":
		out["line-height"] = "1.5"
	case "leading-loose":
		out["line-height"] = "2"

	// Tracking (letter-spacing)
	case "tracking-tighter":
		out["letter-spacing"] = "-0.05em"
	case "tracking-tight":
		out["letter-spacing"] = "-0.025em"
	case "tracking-normal":
		out["letter-spacing"] = "0"
	case "tracking-wide":
		out["letter-spacing"] = "0.025em"
	case "tracking-wider":
		out["letter-spacing"] = "0.05em"
	case "tracking-widest":
		out["letter-spacing"] = "0.1em"

	// Truncate
	case "truncate":
		out["overflow"] = "hidden"
		out["text-overflow"] = "ellipsis"
		out["white-space"] = "nowrap"

	// Whitespace
	case "whitespace-normal":
		out["white-space"] = "normal"
	case "whitespace-nowrap":
		out["white-space"] = "nowrap"
	case "whitespace-pre":
		out["white-space"] = "pre"
	case "whitespace-pre-wrap":
		out["white-space"] = "pre-wrap"

	// Width special values
	case "w-full":
		out["width"] = "100%"
	case "w-screen":
		out["width"] = "100vw"
	case "w-auto":
		out["width"] = "auto"
	case "w-px":
		out["width"] = "1px"
	case "w-1/2":
		out["width"] = "50%"
	case "w-1/3":
		out["width"] = "33.333333%"
	case "w-2/3":
		out["width"] = "66.666667%"
	case "w-1/4":
		out["width"] = "25%"
	case "w-3/4":
		out["width"] = "75%"
	case "w-1/5":
		out["width"] = "20%"
	case "w-2/5":
		out["width"] = "40%"
	case "w-3/5":
		out["width"] = "60%"
	case "w-4/5":
		out["width"] = "80%"
	case "w-1/6":
		out["width"] = "16.666667%"
	case "w-5/6":
		out["width"] = "83.333333%"
	case "w-1/12":
		out["width"] = "8.333333%"
	case "w-5/12":
		out["width"] = "41.666667%"
	case "w-7/12":
		out["width"] = "58.333333%"
	case "w-11/12":
		out["width"] = "91.666667%"

	// Height special values
	case "h-px":
		out["height"] = "1px"
	case "h-full":
		out["height"] = "100%"
	case "h-screen":
		out["height"] = "100vh"
	case "h-auto":
		out["height"] = "auto"
	case "h-1/2":
		out["height"] = "50%"
	case "h-1/3":
		out["height"] = "33.333333%"
	case "h-2/3":
		out["height"] = "66.666667%"
	case "h-1/4":
		out["height"] = "25%"
	case "h-3/4":
		out["height"] = "75%"
	case "h-1/5":
		out["height"] = "20%"
	case "h-2/5":
		out["height"] = "40%"
	case "h-3/5":
		out["height"] = "60%"
	case "h-4/5":
		out["height"] = "80%"
	case "h-1/6":
		out["height"] = "16.666667%"
	case "h-5/6":
		out["height"] = "83.333333%"

	// Min/max width
	case "min-w-0":
		out["min-width"] = "0px"
	case "min-w-full":
		out["min-width"] = "100%"
	case "max-w-sm":
		out["max-width"] = "384px"
	case "max-w-md":
		out["max-width"] = "448px"
	case "max-w-lg":
		out["max-width"] = "512px"
	case "max-w-xl":
		out["max-width"] = "576px"
	case "max-w-2xl":
		out["max-width"] = "672px"
	case "max-w-full":
		out["max-width"] = "100%"
	case "max-w-none":
		out["max-width"] = "none"

	// Min/max height
	case "min-h-0":
		out["min-height"] = "0px"
	case "min-h-full":
		out["min-height"] = "100%"
	case "min-h-screen":
		out["min-height"] = "100vh"
	case "max-h-full":
		out["max-height"] = "100%"
	case "max-h-screen":
		out["max-height"] = "100vh"
	case "max-h-none":
		out["max-height"] = "none"

	// Margin auto
	case "m-auto":
		out["margin-top"] = "auto"
		out["margin-right"] = "auto"
		out["margin-bottom"] = "auto"
		out["margin-left"] = "auto"
	case "mx-auto":
		out["margin-left"] = "auto"
		out["margin-right"] = "auto"
	case "my-auto":
		out["margin-top"] = "auto"
		out["margin-bottom"] = "auto"

	// Border shorthand
	case "border":
		out["border-top-width"] = "1px"
		out["border-right-width"] = "1px"
		out["border-bottom-width"] = "1px"
		out["border-left-width"] = "1px"
	case "border-0":
		out["border-top-width"] = "0px"
		out["border-right-width"] = "0px"
		out["border-bottom-width"] = "0px"
		out["border-left-width"] = "0px"
	case "border-2":
		out["border-top-width"] = "2px"
		out["border-right-width"] = "2px"
		out["border-bottom-width"] = "2px"
		out["border-left-width"] = "2px"
	case "border-4":
		out["border-top-width"] = "4px"
		out["border-right-width"] = "4px"
		out["border-bottom-width"] = "4px"
		out["border-left-width"] = "4px"
	case "border-8":
		out["border-top-width"] = "8px"
		out["border-right-width"] = "8px"
		out["border-bottom-width"] = "8px"
		out["border-left-width"] = "8px"

	// Border style
	case "border-solid":
		out["border-top-style"] = "solid"
		out["border-right-style"] = "solid"
		out["border-bottom-style"] = "solid"
		out["border-left-style"] = "solid"
	case "border-dashed":
		out["border-top-style"] = "dashed"
		out["border-right-style"] = "dashed"
		out["border-bottom-style"] = "dashed"
		out["border-left-style"] = "dashed"
	case "border-dotted":
		out["border-top-style"] = "dotted"
		out["border-right-style"] = "dotted"
		out["border-bottom-style"] = "dotted"
		out["border-left-style"] = "dotted"

	// Border radius
	case "rounded-none":
		out["border-radius"] = "0"
	case "rounded-sm":
		out["border-radius"] = "2px"
	case "rounded":
		out["border-radius"] = "4px"
	case "rounded-md":
		out["border-radius"] = "6px"
	case "rounded-lg":
		out["border-radius"] = "8px"
	case "rounded-xl":
		out["border-radius"] = "12px"
	case "rounded-2xl":
		out["border-radius"] = "16px"
	case "rounded-3xl":
		out["border-radius"] = "24px"
	case "rounded-full":
		out["border-radius"] = "9999px"

	// Shadow
	case "shadow-sm":
		out["box-shadow"] = "0 1px 2px 0 rgba(0,0,0,0.05)"
	case "shadow":
		out["box-shadow"] = "0 1px 3px 0 rgba(0,0,0,0.1), 0 1px 2px -1px rgba(0,0,0,0.1)"
	case "shadow-md":
		out["box-shadow"] = "0 4px 6px -1px rgba(0,0,0,0.1), 0 2px 4px -2px rgba(0,0,0,0.1)"
	case "shadow-lg":
		out["box-shadow"] = "0 10px 15px -3px rgba(0,0,0,0.1), 0 4px 6px -4px rgba(0,0,0,0.1)"
	case "shadow-xl":
		out["box-shadow"] = "0 20px 25px -5px rgba(0,0,0,0.1), 0 8px 10px -6px rgba(0,0,0,0.1)"
	case "shadow-2xl":
		out["box-shadow"] = "0 25px 50px -12px rgba(0,0,0,0.25)"
	case "shadow-none":
		out["box-shadow"] = "none"

	// Display
	case "hidden":
		out["display"] = "none"
	case "block":
		out["display"] = "block"
	case "inline":
		out["display"] = "flex"
	case "inline-flex":
		out["display"] = "flex"
	case "grid":
		out["display"] = "flex"

	// Layout
	case "overflow-hidden":
		out["overflow"] = "hidden"
	case "overflow-visible":
		out["overflow"] = "visible"
	case "relative":
		out["position"] = "relative"
	case "absolute":
		out["position"] = "absolute"
	case "static":
		out["position"] = "static"
	case "fixed":
		out["position"] = "absolute"
	case "sticky":
		out["position"] = "relative"

	// Z-index
	case "z-0":
		out["z-index"] = "0"
	case "z-10":
		out["z-index"] = "10"
	case "z-20":
		out["z-index"] = "20"
	case "z-30":
		out["z-index"] = "30"
	case "z-40":
		out["z-index"] = "40"
	case "z-50":
		out["z-index"] = "50"
	case "z-auto":
		out["z-index"] = "auto"

	// Aspect ratio
	case "aspect-square":
		out["aspect-ratio"] = "1"
	case "aspect-video":
		out["aspect-ratio"] = "1.7778"
	case "aspect-auto":
		out["aspect-ratio"] = "auto"

	// Filters
	case "blur-none":
		out["filter"] = "blur(0)"
	case "blur-sm":
		out["filter"] = "blur(4px)"
	case "blur":
		out["filter"] = "blur(8px)"
	case "blur-md":
		out["filter"] = "blur(12px)"
	case "blur-lg":
		out["filter"] = "blur(16px)"
	case "blur-xl":
		out["filter"] = "blur(24px)"
	case "blur-2xl":
		out["filter"] = "blur(40px)"
	case "blur-3xl":
		out["filter"] = "blur(64px)"
	case "brightness-0":
		out["filter"] = "brightness(0)"
	case "brightness-50":
		out["filter"] = "brightness(.5)"
	case "brightness-75":
		out["filter"] = "brightness(.75)"
	case "brightness-90":
		out["filter"] = "brightness(.9)"
	case "brightness-95":
		out["filter"] = "brightness(.95)"
	case "brightness-100":
		out["filter"] = "brightness(1)"
	case "brightness-105":
		out["filter"] = "brightness(1.05)"
	case "brightness-110":
		out["filter"] = "brightness(1.1)"
	case "brightness-125":
		out["filter"] = "brightness(1.25)"
	case "brightness-150":
		out["filter"] = "brightness(1.5)"
	case "brightness-200":
		out["filter"] = "brightness(2)"
	case "grayscale-0":
		out["filter"] = "grayscale(0)"
	case "grayscale":
		out["filter"] = "grayscale(100%)"

	// Transforms
	case "rotate-0":
		out["transform"] = "rotate(0deg)"
	case "rotate-1":
		out["transform"] = "rotate(1deg)"
	case "rotate-2":
		out["transform"] = "rotate(2deg)"
	case "rotate-3":
		out["transform"] = "rotate(3deg)"
	case "rotate-6":
		out["transform"] = "rotate(6deg)"
	case "rotate-12":
		out["transform"] = "rotate(12deg)"
	case "rotate-45":
		out["transform"] = "rotate(45deg)"
	case "rotate-90":
		out["transform"] = "rotate(90deg)"
	case "rotate-180":
		out["transform"] = "rotate(180deg)"
	case "scale-0":
		out["transform"] = "scale(0)"
	case "scale-50":
		out["transform"] = "scale(.5)"
	case "scale-75":
		out["transform"] = "scale(.75)"
	case "scale-90":
		out["transform"] = "scale(.9)"
	case "scale-95":
		out["transform"] = "scale(.95)"
	case "scale-100":
		out["transform"] = "scale(1)"
	case "scale-105":
		out["transform"] = "scale(1.05)"
	case "scale-110":
		out["transform"] = "scale(1.1)"
	case "scale-125":
		out["transform"] = "scale(1.25)"
	case "scale-150":
		out["transform"] = "scale(1.5)"

	// Sizing shortcuts
	case "w-fit":
		out["width"] = "auto"
	case "h-fit":
		out["height"] = "auto"
	case "w-min":
		out["width"] = "auto"
	case "w-max":
		out["width"] = "auto"

	default:
		resolveTailwindDynamic(cls, out)
	}
}

func resolveTailwindDynamic(cls string, out map[string]string) {
	parts := strings.SplitN(cls, "-", 2)
	if len(parts) < 2 {
		return
	}
	prefix := parts[0]
	rest := parts[1]

	switch prefix {
	case "p":
		if v, ok := parseSpacing(rest); ok {
			out["padding-top"] = v
			out["padding-right"] = v
			out["padding-bottom"] = v
			out["padding-left"] = v
		}
	case "px":
		if v, ok := parseSpacing(rest); ok {
			out["padding-left"] = v
			out["padding-right"] = v
		}
	case "py":
		if v, ok := parseSpacing(rest); ok {
			out["padding-top"] = v
			out["padding-bottom"] = v
		}
	case "pt":
		if v, ok := parseSpacing(rest); ok {
			out["padding-top"] = v
		}
	case "pr":
		if v, ok := parseSpacing(rest); ok {
			out["padding-right"] = v
		}
	case "pb":
		if v, ok := parseSpacing(rest); ok {
			out["padding-bottom"] = v
		}
	case "pl":
		if v, ok := parseSpacing(rest); ok {
			out["padding-left"] = v
		}
	case "m":
		if v, ok := parseSpacing(rest); ok {
			out["margin-top"] = v
			out["margin-right"] = v
			out["margin-bottom"] = v
			out["margin-left"] = v
		}
	case "mx":
		if v, ok := parseSpacing(rest); ok {
			out["margin-left"] = v
			out["margin-right"] = v
		}
	case "my":
		if v, ok := parseSpacing(rest); ok {
			out["margin-top"] = v
			out["margin-bottom"] = v
		}
	case "mt":
		if v, ok := parseSpacing(rest); ok {
			out["margin-top"] = v
		}
	case "mr":
		if v, ok := parseSpacing(rest); ok {
			out["margin-right"] = v
		}
	case "mb":
		if v, ok := parseSpacing(rest); ok {
			out["margin-bottom"] = v
		}
	case "ml":
		if v, ok := parseSpacing(rest); ok {
			out["margin-left"] = v
		}
	case "gap":
		if strings.HasPrefix(rest, "x-") {
			if v, ok := parseSpacing(rest[2:]); ok {
				out["column-gap"] = v
			}
		} else if strings.HasPrefix(rest, "y-") {
			if v, ok := parseSpacing(rest[2:]); ok {
				out["row-gap"] = v
			}
		} else if v, ok := parseSpacing(rest); ok {
			out["gap"] = v
			out["row-gap"] = v
			out["column-gap"] = v
		}
	case "space":
		if strings.HasPrefix(rest, "x-") {
			if v, ok := parseSpacing(rest[2:]); ok {
				out["column-gap"] = v
			}
		} else if strings.HasPrefix(rest, "y-") {
			if v, ok := parseSpacing(rest[2:]); ok {
				out["row-gap"] = v
			}
		}
	case "w":
		if v, ok := parseSpacing(rest); ok {
			out["width"] = v
		}
	case "h":
		if v, ok := parseSpacing(rest); ok {
			out["height"] = v
		}
	case "size":
		if v, ok := parseSpacing(rest); ok {
			out["width"] = v
			out["height"] = v
		}
	case "top":
		if v, ok := parseSpacing(rest); ok {
			out["top"] = v
		}
	case "right":
		if v, ok := parseSpacing(rest); ok {
			out["right"] = v
		}
	case "bottom":
		if v, ok := parseSpacing(rest); ok {
			out["bottom"] = v
		}
	case "left":
		if v, ok := parseSpacing(rest); ok {
			out["left"] = v
		}
	case "inset":
		if v, ok := parseSpacing(rest); ok {
			out["top"] = v
			out["right"] = v
			out["bottom"] = v
			out["left"] = v
		}
	case "opacity":
		if n, err := strconv.Atoi(rest); err == nil {
			out["opacity"] = fmt.Sprintf("%.2f", float64(n)/100.0)
		}
	case "leading":
		if v, ok := parseSpacing(rest); ok {
			out["line-height"] = v
		}
	case "line":
		if strings.HasPrefix(rest, "clamp-") {
			if n, err := strconv.Atoi(rest[6:]); err == nil && n >= 1 && n <= 6 {
				out["-webkit-line-clamp"] = strconv.Itoa(n)
				out["overflow"] = "hidden"
			}
		}
	case "translate":
		if strings.HasPrefix(rest, "x-") {
			if v, ok := parseSpacing(rest[2:]); ok {
				out["transform"] = "translateX(" + v + ")"
			}
		} else if strings.HasPrefix(rest, "y-") {
			if v, ok := parseSpacing(rest[2:]); ok {
				out["transform"] = "translateY(" + v + ")"
			}
		}
	case "rotate":
		if n, err := strconv.Atoi(rest); err == nil {
			out["transform"] = fmt.Sprintf("rotate(%ddeg)", n)
		}
	case "scale":
		if strings.HasPrefix(rest, "x-") {
			if n, err := strconv.Atoi(rest[2:]); err == nil {
				out["transform"] = fmt.Sprintf("scaleX(%.2f)", float64(n)/100.0)
			}
		} else if strings.HasPrefix(rest, "y-") {
			if n, err := strconv.Atoi(rest[2:]); err == nil {
				out["transform"] = fmt.Sprintf("scaleY(%.2f)", float64(n)/100.0)
			}
		} else if n, err := strconv.Atoi(rest); err == nil {
			out["transform"] = fmt.Sprintf("scale(%.2f)", float64(n)/100.0)
		}
	case "skew":
		if strings.HasPrefix(rest, "x-") {
			if n, err := strconv.Atoi(rest[2:]); err == nil {
				out["transform"] = fmt.Sprintf("skewX(%ddeg)", n)
			}
		} else if strings.HasPrefix(rest, "y-") {
			if n, err := strconv.Atoi(rest[2:]); err == nil {
				out["transform"] = fmt.Sprintf("skewY(%ddeg)", n)
			}
		}
	case "text", "bg":
		resolveColorPrefix(prefix, rest, out)
	case "border":
		resolveBorderDynamic(rest, out)
	}
}

func resolveColorPrefix(prefix, rest string, out map[string]string) {
	prop, val, ok := resolveColorClass(prefix, rest)
	if ok {
		out[prop] = val
	}
}

func resolveBorderDynamic(rest string, out map[string]string) {
	if prop, val, ok := resolveColorClass("border", rest); ok {
		out[prop] = val
		return
	}

	sides := map[string]struct {
		prop string
	}{
		"t-": {prop: "border-top-width"},
		"r-": {prop: "border-right-width"},
		"b-": {prop: "border-bottom-width"},
		"l-": {prop: "border-left-width"},
	}
	for sidePrefix, info := range sides {
		if strings.HasPrefix(rest, sidePrefix) {
			numStr := rest[len(sidePrefix):]
			if n, err := strconv.Atoi(numStr); err == nil {
				out[info.prop] = fmt.Sprintf("%dpx", n)
				return
			}
		}
	}
}

func resolveArbitrary(cls string, out map[string]string) {
	bracketIdx := strings.Index(cls, "[")
	if bracketIdx < 0 {
		return
	}
	prefix := cls[:bracketIdx]
	if strings.HasSuffix(prefix, "-") {
		prefix = prefix[:len(prefix)-1]
	}
	value := cls[bracketIdx+1 : len(cls)-1]

	switch prefix {
	case "text":
		out["font-size"] = value
	case "bg":
		out["background-color"] = value
	case "w":
		out["width"] = value
	case "h":
		out["height"] = value
	case "p":
		out["padding-top"] = value
		out["padding-right"] = value
		out["padding-bottom"] = value
		out["padding-left"] = value
	case "m":
		out["margin-top"] = value
		out["margin-right"] = value
		out["margin-bottom"] = value
		out["margin-left"] = value
	case "rounded":
		out["border-radius"] = value
	case "gap":
		out["gap"] = value
		out["row-gap"] = value
		out["column-gap"] = value
	case "top":
		out["top"] = value
	case "right":
		out["right"] = value
	case "bottom":
		out["bottom"] = value
	case "left":
		out["left"] = value
	case "border":
		out["border-top-width"] = value
		out["border-right-width"] = value
		out["border-bottom-width"] = value
		out["border-left-width"] = value
	case "opacity":
		out["opacity"] = value
	case "leading":
		out["line-height"] = value
	case "tracking":
		out["letter-spacing"] = value
	case "rotate":
		out["transform"] = "rotate(" + value + ")"
	case "scale":
		out["transform"] = "scale(" + value + ")"
	case "translate":
		out["transform"] = "translate(" + value + ")"
	case "blur":
		out["filter"] = "blur(" + value + ")"
	case "brightness":
		out["filter"] = "brightness(" + value + ")"
	case "grayscale":
		out["filter"] = "grayscale(" + value + ")"
	}
}
