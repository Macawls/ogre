package parse

import (
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// Parse parses an HTML fragment and returns the root Node.
// Parse parses an HTML string into a node tree.
func Parse(htmlStr string) (*Node, error) {
	nodes, err := html.ParseFragment(strings.NewReader(htmlStr), &html.Node{
		Type:     html.ElementNode,
		Data:     "body",
		DataAtom: atom.Body,
	})
	if err != nil {
		return nil, err
	}

	var roots []*Node
	for _, n := range nodes {
		if converted := convertNode(n); converted != nil {
			roots = append(roots, converted)
		}
	}

	if len(roots) == 0 {
		return &Node{Type: ElementNode, Tag: "div"}, nil
	}
	if len(roots) == 1 {
		return roots[0], nil
	}

	return &Node{
		Type:     ElementNode,
		Tag:      "div",
		Style:    map[string]string{"display": "flex"},
		Children: roots,
	}, nil
}

func convertNode(n *html.Node) *Node {
	switch n.Type {
	case html.TextNode:
		text := strings.TrimSpace(n.Data)
		if text == "" {
			return nil
		}
		return &Node{Type: TextNode, Text: text}

	case html.ElementNode:
		if n.Data == "br" {
			return &Node{Type: TextNode, Text: "\n"}
		}

		node := &Node{
			Type:  ElementNode,
			Tag:   n.Data,
			Attrs: make(map[string]string),
			Style: make(map[string]string),
		}

		for _, attr := range n.Attr {
			if attr.Key == "style" {
				parseStyleAttr(attr.Val, node.Style)
			} else if attr.Key == "class" {
				classes := strings.Fields(attr.Val)
				if len(classes) > 0 {
					node.Classes = classes
				}
			} else {
				node.Attrs[attr.Key] = attr.Val
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if child := convertNode(c); child != nil {
				node.Children = append(node.Children, child)
			}
		}

		return node

	default:
		return nil
	}
}

func parseStyleAttr(style string, m map[string]string) {
	for _, decl := range strings.Split(style, ";") {
		decl = strings.TrimSpace(decl)
		if decl == "" {
			continue
		}
		parts := strings.SplitN(decl, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		if key != "" && val != "" {
			m[key] = val
		}
	}
}
