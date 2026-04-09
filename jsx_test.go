package ogre

import (
	"strings"
	"testing"
)

func TestDivWithText(t *testing.T) {
	e := Div(Props{Class: "flex w-full h-full bg-blue-500"}, "Hello")
	html := e.ToHTML()
	if !strings.Contains(html, "Hello") {
		t.Errorf("expected Hello in output, got %s", html)
	}
	if !strings.Contains(html, `class="flex w-full h-full bg-blue-500"`) {
		t.Errorf("expected class attribute, got %s", html)
	}
}

func TestNestedElements(t *testing.T) {
	e := Div(Props{Class: "flex flex-col w-full h-full bg-slate-900 p-16"},
		Div(Props{Class: "text-5xl font-bold text-white"}, "Title"),
		Div(Props{Class: "text-xl text-slate-400 mt-4"}, "Subtitle"),
	)
	html := e.ToHTML()
	if !strings.Contains(html, "Title") {
		t.Errorf("expected Title, got %s", html)
	}
	if !strings.Contains(html, "Subtitle") {
		t.Errorf("expected Subtitle, got %s", html)
	}
}

func TestInlineStyles(t *testing.T) {
	e := Div(Props{
		Style: map[string]string{
			"background-color": "#ff0000",
			"padding":          "20px",
		},
	}, "Red box")
	html := e.ToHTML()
	if !strings.Contains(html, "background-color") {
		t.Errorf("expected background-color, got %s", html)
	}
}

func TestImgElement(t *testing.T) {
	e := Img(Props{Src: "https://example.com/logo.png", Alt: "Logo"})
	html := e.ToHTML()
	if !strings.Contains(html, `src="https://example.com/logo.png"`) {
		t.Errorf("expected src attr, got %s", html)
	}
	if !strings.Contains(html, "/>") {
		t.Errorf("expected self-closing tag, got %s", html)
	}
}

func TestTextElement(t *testing.T) {
	e := Text("plain text")
	html := e.ToHTML()
	if html != "plain text" {
		t.Errorf("expected plain text, got %s", html)
	}
}

func TestComplexTree(t *testing.T) {
	e := Div(Props{Class: "flex w-full h-full", Style: map[string]string{"background-image": "linear-gradient(135deg, #0f0c29, #302b63)"}},
		Div(Props{Class: "flex flex-col flex-1 p-16 justify-center"},
			Span(Props{Class: "text-sm font-bold text-purple-400"}, "BLOG"),
			Div(Props{Class: "text-5xl font-bold text-white mt-4"}, "My Post Title"),
			P(Props{Class: "text-xl text-slate-400 mt-4"}, "A description of the post"),
		),
	)
	html := e.ToHTML()
	if !strings.Contains(html, "BLOG") {
		t.Errorf("expected BLOG, got %s", html)
	}
	if !strings.Contains(html, "My Post Title") {
		t.Errorf("expected title, got %s", html)
	}
}
