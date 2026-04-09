package parse

import (
	"testing"
)

func TestSimpleDivWithText(t *testing.T) {
	node, err := Parse(`<div>hello</div>`)
	if err != nil {
		t.Fatal(err)
	}
	if node.Tag != "div" {
		t.Fatalf("expected div, got %s", node.Tag)
	}
	if len(node.Children) != 1 {
		t.Fatalf("expected 1 child, got %d", len(node.Children))
	}
	if node.Children[0].Type != TextNode || node.Children[0].Text != "hello" {
		t.Fatalf("expected text node 'hello', got %+v", node.Children[0])
	}
}

func TestNestedElements(t *testing.T) {
	node, err := Parse(`<div><span><p>deep</p></span></div>`)
	if err != nil {
		t.Fatal(err)
	}
	if node.Tag != "div" {
		t.Fatalf("expected div, got %s", node.Tag)
	}
	span := node.Children[0]
	if span.Tag != "span" {
		t.Fatalf("expected span, got %s", span.Tag)
	}
	p := span.Children[0]
	if p.Tag != "p" {
		t.Fatalf("expected p, got %s", p.Tag)
	}
	if p.Children[0].Text != "deep" {
		t.Fatalf("expected 'deep', got %s", p.Children[0].Text)
	}
}

func TestStyleAttributeExtraction(t *testing.T) {
	node, err := Parse(`<div style="color: red; font-size: 16px;"></div>`)
	if err != nil {
		t.Fatal(err)
	}
	if node.Style["color"] != "red" {
		t.Fatalf("expected color=red, got %s", node.Style["color"])
	}
	if node.Style["font-size"] != "16px" {
		t.Fatalf("expected font-size=16px, got %s", node.Style["font-size"])
	}
}

func TestMultipleRootElements(t *testing.T) {
	node, err := Parse(`<div>a</div><span>b</span>`)
	if err != nil {
		t.Fatal(err)
	}
	if node.Tag != "div" {
		t.Fatalf("expected wrapper div, got %s", node.Tag)
	}
	if node.Style["display"] != "flex" {
		t.Fatalf("expected display:flex on wrapper, got %s", node.Style["display"])
	}
	if len(node.Children) != 2 {
		t.Fatalf("expected 2 children, got %d", len(node.Children))
	}
	if node.Children[0].Tag != "div" {
		t.Fatalf("expected first child div, got %s", node.Children[0].Tag)
	}
	if node.Children[1].Tag != "span" {
		t.Fatalf("expected second child span, got %s", node.Children[1].Tag)
	}
}

func TestWhitespaceOnlyTextNodesSkipped(t *testing.T) {
	node, err := Parse(`<div>   </div>`)
	if err != nil {
		t.Fatal(err)
	}
	if len(node.Children) != 0 {
		t.Fatalf("expected 0 children (whitespace skipped), got %d", len(node.Children))
	}
}

func TestMixedTextAndWhitespace(t *testing.T) {
	node, err := Parse(`<div>  hello  </div>`)
	if err != nil {
		t.Fatal(err)
	}
	if len(node.Children) != 1 {
		t.Fatalf("expected 1 child, got %d", len(node.Children))
	}
	if node.Children[0].Text != "hello" {
		t.Fatalf("expected 'hello', got '%s'", node.Children[0].Text)
	}
}

func TestClassAttribute(t *testing.T) {
	node, err := Parse(`<div class="foo bar baz"></div>`)
	if err != nil {
		t.Fatal(err)
	}
	if len(node.Classes) != 3 {
		t.Fatalf("expected 3 classes, got %d", len(node.Classes))
	}
	if node.Classes[0] != "foo" || node.Classes[1] != "bar" || node.Classes[2] != "baz" {
		t.Fatalf("unexpected classes: %v", node.Classes)
	}
}

func TestAttributes(t *testing.T) {
	node, err := Parse(`<img src="test.png" alt="test"/>`)
	if err != nil {
		t.Fatal(err)
	}
	if node.Attrs["src"] != "test.png" {
		t.Fatalf("expected src=test.png, got %s", node.Attrs["src"])
	}
	if node.Attrs["alt"] != "test" {
		t.Fatalf("expected alt=test, got %s", node.Attrs["alt"])
	}
}
