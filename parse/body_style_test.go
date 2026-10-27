package parse

import "testing"

const fullDocWithBodyStyle = `<!DOCTYPE html><html><head><meta charset="utf-8"></head>
<body style="background:#ff0000;font-family:'DejaVu Sans';color:#00ff00;padding:8px;">
<div>text</div>
</body></html>`

// A full document declares its font, background and base colour on <body>.
// ParseFragment drops the <body> element together with those attributes, so
// without special handling every such document renders with built-in defaults.
func TestFullDocumentKeepsBodyStyle(t *testing.T) {
	root, err := Parse(fullDocWithBodyStyle)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}

	if got := root.Style["background"]; got != "#ff0000" {
		t.Errorf("root background = %q, want %q (body style was dropped)", got, "#ff0000")
	}
	if got := root.Style["font-family"]; got != "'DejaVu Sans'" {
		t.Errorf("root font-family = %q, want %q (body style was dropped)", got, "'DejaVu Sans'")
	}
	if got := root.Style["color"]; got != "#00ff00" {
		t.Errorf("root color = %q, want %q (body style was dropped)", got, "#00ff00")
	}

	// The body's children must survive the document path unchanged.
	if len(root.Children) != 1 || root.Children[0].Tag != "div" {
		t.Fatalf("root children = %d, want one <div>", len(root.Children))
	}
}

// A fragment has no <body> element, so the fragment path must stay in charge:
// the root is the fragment's first element, not a synthesised body.
func TestFragmentUnaffected(t *testing.T) {
	root, err := Parse(`<div style="color:#123456">x</div>`)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if root.Tag != "div" {
		t.Errorf("root tag = %q, want div", root.Tag)
	}
	if got := root.Style["color"]; got != "#123456" {
		t.Errorf("root color = %q, want %q", got, "#123456")
	}
}

// <body> without any attribute carries nothing to preserve, so the fragment
// path must remain in charge and the tree must not gain a wrapper element.
func TestBodyWithoutAttrsUsesFragmentPath(t *testing.T) {
	root, err := Parse(`<body><div>only child</div></body>`)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if root.Tag != "div" {
		t.Errorf("root tag = %q, want div (no wrapper expected)", root.Tag)
	}
}
