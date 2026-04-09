---
title: JSX-style Builder
description: Build templates programmatically with Go functions instead of HTML strings.
---

Ogre provides a JSX-inspired builder API for constructing templates in Go code. Instead of writing HTML strings, you compose elements with function calls.

## Basic usage

```go
e := ogre.Div(ogre.Props{Class: "flex w-full h-full bg-blue-500 items-center justify-center"},
    ogre.Div(ogre.Props{Class: "text-4xl font-bold text-white"}, "Hello World"),
)

result, err := e.Render(ogre.Options{Width: 1200, Height: 630})
```

## Elements

| Function | HTML |
|----------|------|
| `ogre.Div(props, children...)` | `<div>` |
| `ogre.Span(props, children...)` | `<span>` |
| `ogre.P(props, children...)` | `<p>` |
| `ogre.A(props, children...)` | `<a>` |
| `ogre.Img(props)` | `<img/>` |
| `ogre.Text(s)` | Text node |

## Props

```go
type Props struct {
    Style   map[string]string // Inline CSS properties
    Class   string            // Space-separated CSS/Tailwind classes
    Src     string            // Image source URL (for Img)
    Alt     string            // Alt text (for Img)
    Href    string            // Link URL (for A)
}
```

## Children

Children can be `*Element` values or `string` values. Strings are converted to text nodes.

```go
ogre.Div(ogre.Props{},
    "Plain text",                                    // string → text node
    ogre.Span(ogre.Props{Class: "font-bold"}, "Bold"), // element
)
```

## Inline styles

```go
ogre.Div(ogre.Props{
    Style: map[string]string{
        "background-image": "linear-gradient(135deg, #0f0c29, #302b63)",
        "padding":          "60px",
    },
}, "Gradient background")
```

## Tailwind classes

```go
ogre.Div(ogre.Props{Class: "flex flex-col w-full h-full bg-slate-900 p-16 justify-center"},
    ogre.Div(ogre.Props{Class: "text-5xl font-bold text-white"}, "Title"),
    ogre.Div(ogre.Props{Class: "text-xl text-slate-400 mt-4"}, "Subtitle"),
)
```

## Rendering

Three ways to render:

```go
// Direct render (creates a new renderer each time)
result, err := element.Render(ogre.Options{Width: 1200, Height: 630})

// Render with a shared renderer
r := ogre.NewRenderer()
result, err := element.RenderWith(r, ogre.Options{Width: 1200, Height: 630})

// Convert to HTML string first
html := element.ToHTML()
result, err := ogre.Render(html, ogre.Options{Width: 1200, Height: 630})
```

## Full example

```go
card := ogre.Div(ogre.Props{
    Class: "flex w-full h-full",
    Style: map[string]string{
        "background-image": "linear-gradient(135deg, #0f0c29 0%, #302b63 50%, #24243e 100%)",
    },
},
    ogre.Div(ogre.Props{Class: "flex flex-col flex-1 p-16 justify-center"},
        ogre.Span(ogre.Props{Class: "text-sm font-bold text-purple-400"},
            "ENGINEERING BLOG",
        ),
        ogre.Div(ogre.Props{Class: "text-5xl font-bold text-white mt-4"},
            "Building a Pure Go OG Image Generator",
        ),
        ogre.P(ogre.Props{Class: "text-xl text-slate-400 mt-4"},
            "Zero dependencies. Single binary. Blazing fast.",
        ),
    ),
)

result, _ := card.Render(ogre.Options{
    Width:  1200,
    Height: 630,
    Format: ogre.FormatPNG,
})
os.WriteFile("blog-card.png", result.Data, 0644)
```
