package style

import (
	"testing"

	"github.com/macawls/ogre/v3/parse"
)

func TestTailwindBackgroundClipText(t *testing.T) {
	result := ResolveTailwind([]string{"bg-clip-text"}, TailwindV3)
	if result["background-clip"] != "text" {
		t.Errorf("bg-clip-text: got %q, want %q", result["background-clip"], "text")
	}

	result = ResolveTailwind([]string{"text-transparent"}, TailwindV3)
	if result["color"] != "transparent" {
		t.Errorf("text-transparent color: got %q, want transparent", result["color"])
	}
}

func TestBackgroundClipTextPropagatesToTextNode(t *testing.T) {
	text := &parse.Node{Type: parse.TextNode, Text: "hi"}
	root := &parse.Node{
		Type: parse.ElementNode,
		Tag:  "div",
		Style: map[string]string{
			"background-image": "linear-gradient(to right,#f00,#00f)",
			"background-clip":  "text",
		},
		Children: []*parse.Node{text},
	}

	result := Resolve(root, 1200, 630, TailwindV3)

	elem := result[root]
	if elem.BackgroundClip != "text" {
		t.Fatalf("element BackgroundClip: got %q, want text", elem.BackgroundClip)
	}

	child := result[text]
	if child.BackgroundClip != "text" {
		t.Errorf("text node BackgroundClip: got %q, want text", child.BackgroundClip)
	}
	if child.BackgroundImage != elem.BackgroundImage {
		t.Errorf("text node BackgroundImage: got %q, want %q", child.BackgroundImage, elem.BackgroundImage)
	}
}

func TestBackgroundClipTextPropagatesThroughWrapperElement(t *testing.T) {
	text := &parse.Node{Type: parse.TextNode, Text: "hi"}
	span := &parse.Node{Type: parse.ElementNode, Tag: "span", Children: []*parse.Node{text}}
	root := &parse.Node{
		Type: parse.ElementNode,
		Tag:  "div",
		Style: map[string]string{
			"background-image": "linear-gradient(to right,#f00,#00f)",
			"background-clip":  "text",
		},
		Children: []*parse.Node{span},
	}

	result := Resolve(root, 1200, 630, TailwindV3)

	elem := result[root]
	if got := result[span]; got.BackgroundClip != "text" || got.BackgroundImage != elem.BackgroundImage {
		t.Fatalf("span BackgroundClip/Image: got %q/%q, want text/%q", got.BackgroundClip, got.BackgroundImage, elem.BackgroundImage)
	}
	if got := result[text]; got.BackgroundClip != "text" || got.BackgroundImage != elem.BackgroundImage {
		t.Fatalf("nested text BackgroundClip/Image: got %q/%q, want text/%q", got.BackgroundClip, got.BackgroundImage, elem.BackgroundImage)
	}
}

func TestBackgroundClipTextRespectsChildOverride(t *testing.T) {
	overriddenText := &parse.Node{Type: parse.TextNode, Text: "plain"}
	span := &parse.Node{
		Type: parse.ElementNode,
		Tag:  "span",
		Style: map[string]string{
			"background-image": "linear-gradient(to right,#0f0,#0ff)",
		},
		Children: []*parse.Node{overriddenText},
	}
	root := &parse.Node{
		Type: parse.ElementNode,
		Tag:  "div",
		Style: map[string]string{
			"background-image": "linear-gradient(to right,#f00,#00f)",
			"background-clip":  "text",
		},
		Children: []*parse.Node{span},
	}

	result := Resolve(root, 1200, 630, TailwindV3)

	if got := result[span].BackgroundImage; got != "linear-gradient(to right,#0f0,#0ff)" {
		t.Errorf("span with own background-image should not be overridden by parent clip: got %q", got)
	}
}

func TestBackgroundImageNotInheritedWithoutClipText(t *testing.T) {
	text := &parse.Node{Type: parse.TextNode, Text: "hi"}
	root := &parse.Node{
		Type: parse.ElementNode,
		Tag:  "div",
		Style: map[string]string{
			"background-image": "linear-gradient(to right,#f00,#00f)",
		},
		Children: []*parse.Node{text},
	}

	result := Resolve(root, 1200, 630, TailwindV3)
	if result[text].BackgroundImage != "" {
		t.Errorf("text node BackgroundImage leaked without clip: got %q", result[text].BackgroundImage)
	}
}
