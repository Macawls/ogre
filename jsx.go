package ogre

import (
	"fmt"
	"maps"
	"strings"

	"github.com/macawls/ogre/parse"
)

type Element struct {
	node *parse.Node
}

type Props struct {
	Style   map[string]string
	Class   string
	Src     string
	Alt     string
	Href    string
	Display string
}

func Div(props Props, children ...any) *Element {
	return el("div", props, children)
}

func Span(props Props, children ...any) *Element {
	return el("span", props, children)
}

func P(props Props, children ...any) *Element {
	return el("p", props, children)
}

func Img(props Props) *Element {
	e := el("img", props, nil)
	if props.Src != "" {
		e.node.Attrs["src"] = props.Src
	}
	if props.Alt != "" {
		e.node.Attrs["alt"] = props.Alt
	}
	return e
}

func A(props Props, children ...any) *Element {
	e := el("a", props, children)
	if props.Href != "" {
		e.node.Attrs["href"] = props.Href
	}
	return e
}

func Text(s string) *Element {
	return &Element{node: &parse.Node{Type: parse.TextNode, Text: s}}
}

func el(tag string, props Props, children []any) *Element {
	node := &parse.Node{
		Type:  parse.ElementNode,
		Tag:   tag,
		Attrs: make(map[string]string),
		Style: make(map[string]string),
	}

	if props.Class != "" {
		node.Classes = strings.Fields(props.Class)
	}

	maps.Copy(node.Style, props.Style)

	for _, child := range children {
		switch c := child.(type) {
		case *Element:
			if c != nil && c.node != nil {
				node.Children = append(node.Children, c.node)
			}
		case string:
			if c != "" {
				node.Children = append(node.Children, &parse.Node{
					Type: parse.TextNode,
					Text: c,
				})
			}
		}
	}

	return &Element{node: node}
}

func (e *Element) ToHTML() string {
	if e == nil || e.node == nil {
		return ""
	}
	var b strings.Builder
	writeNode(&b, e.node)
	return b.String()
}

func (e *Element) Render(opts Options) (*Result, error) {
	return Render(e.ToHTML(), opts)
}

func (e *Element) RenderWith(r *Renderer, opts Options) (*Result, error) {
	return r.Render(e.ToHTML(), opts)
}

func writeNode(b *strings.Builder, n *parse.Node) {
	if n.Type == parse.TextNode {
		b.WriteString(n.Text)
		return
	}

	b.WriteByte('<')
	b.WriteString(n.Tag)

	if len(n.Classes) > 0 {
		b.WriteString(` class="`)
		b.WriteString(strings.Join(n.Classes, " "))
		b.WriteByte('"')
	}

	if len(n.Style) > 0 {
		b.WriteString(` style="`)
		first := true
		for k, v := range n.Style {
			if !first {
				b.WriteByte(';')
			}
			fmt.Fprintf(b, "%s:%s", k, v)
			first = false
		}
		b.WriteByte('"')
	}

	for k, v := range n.Attrs {
		fmt.Fprintf(b, ` %s="%s"`, k, v)
	}

	if n.Tag == "img" || n.Tag == "br" || n.Tag == "hr" {
		b.WriteString("/>")
		return
	}

	b.WriteByte('>')

	for _, child := range n.Children {
		writeNode(b, child)
	}

	b.WriteString("</")
	b.WriteString(n.Tag)
	b.WriteByte('>')
}
