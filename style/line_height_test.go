package style

import (
	"testing"

	"github.com/macawls/ogre/v3/parse"
)

// buildTree resolves a small tree: parent -> children, each with its own style
// map, so inheritance can be checked without going through the HTML parser.
func buildTree(t *testing.T, parentStyle map[string]string, childStyles ...map[string]string) (*ComputedStyle, []*ComputedStyle) {
	t.Helper()

	root := &parse.Node{Type: parse.ElementNode, Tag: "div", Style: parentStyle}
	for _, cs := range childStyles {
		root.Children = append(root.Children, &parse.Node{
			Type:  parse.ElementNode,
			Tag:   "div",
			Style: cs,
		})
	}

	result := Resolve(root, 1200, 630, TailwindV3)
	children := make([]*ComputedStyle, 0, len(root.Children))
	for _, c := range root.Children {
		children = append(children, result[c])
	}
	return result[root], children
}

func TestLineHeightDefault(t *testing.T) {
	root, kids := buildTree(t, map[string]string{"font-size": "20px"}, map[string]string{})

	if !approxEqual(root.LineHeight, 24, 0.01) {
		t.Errorf("root line-height: got %.2f, want 24 (1.2 x 20)", root.LineHeight)
	}
	if !approxEqual(kids[0].LineHeight, 24, 0.01) {
		t.Errorf("child line-height: got %.2f, want 24 inherited", kids[0].LineHeight)
	}
}

func TestLineHeightUnitlessInheritsAsMultiplier(t *testing.T) {
	root, kids := buildTree(t,
		map[string]string{"font-size": "20px", "line-height": "2.0"},
		map[string]string{"font-size": "10px"},
		map[string]string{"font-size": "30px"},
		map[string]string{},
	)

	if !approxEqual(root.LineHeight, 40, 0.01) {
		t.Errorf("root line-height: got %.2f, want 40", root.LineHeight)
	}
	if !approxEqual(root.LineHeightFactor, 2.0, 0.001) {
		t.Errorf("root line-height factor: got %.3f, want 2.0", root.LineHeightFactor)
	}
	// A unitless line-height inherits as the number, so each child resolves it
	// against its own font-size: 2.0 x 10 = 20, not the ancestor's 40.
	if !approxEqual(kids[0].LineHeight, 20, 0.01) {
		t.Errorf("child font-size 10 line-height: got %.2f, want 20 (2.0 x 10)", kids[0].LineHeight)
	}
	if !approxEqual(kids[1].LineHeight, 60, 0.01) {
		t.Errorf("child font-size 30 line-height: got %.2f, want 60 (2.0 x 30)", kids[1].LineHeight)
	}
	if !approxEqual(kids[2].LineHeight, 40, 0.01) {
		t.Errorf("child without declaration line-height: got %.2f, want 40", kids[2].LineHeight)
	}
}

func TestLineHeightPxInheritsAsLength(t *testing.T) {
	root, kids := buildTree(t,
		map[string]string{"font-size": "20px", "line-height": "60px"},
		map[string]string{"font-size": "10px"},
	)

	if !approxEqual(root.LineHeight, 60, 0.01) {
		t.Errorf("root line-height: got %.2f, want 60", root.LineHeight)
	}
	if !approxEqual(root.LineHeightFactor, 0, 0.001) {
		t.Errorf("root line-height factor: got %.3f, want 0 for a length", root.LineHeightFactor)
	}
	// A length inherits as computed px, so the child keeps 60 despite 10px text.
	if !approxEqual(kids[0].LineHeight, 60, 0.01) {
		t.Errorf("child line-height: got %.2f, want 60 inherited as a length", kids[0].LineHeight)
	}
}

func TestLineHeightExplicitOverride(t *testing.T) {
	root, kids := buildTree(t,
		map[string]string{"font-size": "20px", "line-height": "2.0"},
		map[string]string{"line-height": "3.0"},
		map[string]string{"line-height": "50px"},
	)

	if !approxEqual(root.LineHeight, 40, 0.01) {
		t.Errorf("root line-height: got %.2f, want 40", root.LineHeight)
	}
	if !approxEqual(kids[0].LineHeight, 60, 0.01) {
		t.Errorf("child line-height 3.0: got %.2f, want 60 (3.0 x 20)", kids[0].LineHeight)
	}
	if !approxEqual(kids[1].LineHeight, 50, 0.01) {
		t.Errorf("child line-height 50px: got %.2f, want 50", kids[1].LineHeight)
	}
}

func TestLineHeightPercent(t *testing.T) {
	root, _ := buildTree(t, map[string]string{"font-size": "20px", "line-height": "150%"})

	if !approxEqual(root.LineHeight, 30, 0.01) {
		t.Errorf("line-height 150%%: got %.2f, want 30", root.LineHeight)
	}
	if !approxEqual(root.LineHeightFactor, 0, 0.001) {
		t.Errorf("line-height 150%% factor: got %.3f, want 0", root.LineHeightFactor)
	}
}
