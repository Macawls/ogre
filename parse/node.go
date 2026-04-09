// Package parse converts HTML strings into a node tree with inline styles.
package parse

// NodeType distinguishes between element and text nodes.
type NodeType int

const (
	// ElementNode represents an HTML element.
	ElementNode NodeType = iota
	// TextNode represents a text content node.
	TextNode
)

// Node represents a single element or text node in the parsed HTML tree.
type Node struct {
	Type     NodeType
	Tag      string
	Attrs    map[string]string
	Style    map[string]string
	Classes  []string
	Children []*Node
	Text     string
}

// CountNodes returns the total number of nodes in the subtree rooted at n.
func CountNodes(n *Node) int {
	if n == nil {
		return 0
	}
	count := 1
	for _, c := range n.Children {
		count += CountNodes(c)
	}
	return count
}
