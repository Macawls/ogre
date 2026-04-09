package font

import (
	"testing"
)

func newWrapConfig(maxWidth float64) WrapConfig {
	return WrapConfig{
		MaxWidth:      maxWidth,
		FontFace:      newMockFace(10, 12, 4),
		FontSize:      16,
		LineHeight:    20,
		LetterSpacing: 0,
		WhiteSpace:    wsNormal,
		WordBreak:     wbNormal,
	}
}

func TestWrapSimpleWordBoundary(t *testing.T) {
	cfg := newWrapConfig(110)
	lines := WrapText("hello world foo", cfg)

	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d: %+v", len(lines), lines)
	}
	if lines[0].Text != "hello world" {
		t.Errorf("line 0: expected %q, got %q", "hello world", lines[0].Text)
	}
	if lines[1].Text != "foo" {
		t.Errorf("line 1: expected %q, got %q", "foo", lines[1].Text)
	}
	if lines[0].Y != 0 {
		t.Errorf("line 0 Y: expected 0, got %f", lines[0].Y)
	}
	if lines[1].Y != 20 {
		t.Errorf("line 1 Y: expected 20, got %f", lines[1].Y)
	}
}

func TestWrapNowrap(t *testing.T) {
	cfg := newWrapConfig(50)
	cfg.WhiteSpace = wsNowrap
	lines := WrapText("hello world foo", cfg)

	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d: %+v", len(lines), lines)
	}
	if lines[0].Text != "hello world foo" {
		t.Errorf("expected %q, got %q", "hello world foo", lines[0].Text)
	}
}

func TestWrapPre(t *testing.T) {
	cfg := newWrapConfig(1000)
	cfg.WhiteSpace = wsPre
	lines := WrapText("hello  world\n  foo  bar", cfg)

	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d: %+v", len(lines), lines)
	}
	if lines[0].Text != "hello  world" {
		t.Errorf("line 0: expected %q, got %q", "hello  world", lines[0].Text)
	}
	if lines[1].Text != "  foo  bar" {
		t.Errorf("line 1: expected %q, got %q", "  foo  bar", lines[1].Text)
	}
}

func TestWrapBreakWord(t *testing.T) {
	cfg := newWrapConfig(50)
	cfg.WordBreak = wbBreakWord
	lines := WrapText("abcdefghij", cfg)

	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d: %+v", len(lines), lines)
	}
	if lines[0].Text != "abcde" {
		t.Errorf("line 0: expected %q, got %q", "abcde", lines[0].Text)
	}
	if lines[1].Text != "fghij" {
		t.Errorf("line 1: expected %q, got %q", "fghij", lines[1].Text)
	}
}

func TestWrapEmpty(t *testing.T) {
	cfg := newWrapConfig(100)
	lines := WrapText("", cfg)

	if lines != nil {
		t.Fatalf("expected nil, got %+v", lines)
	}
}

func TestWrapSingleLongWord(t *testing.T) {
	cfg := newWrapConfig(50)
	lines := WrapText("abcdefghijklmno", cfg)

	if len(lines) != 1 {
		t.Fatalf("expected 1 line (no break without break-word), got %d: %+v", len(lines), lines)
	}
	if lines[0].Text != "abcdefghijklmno" {
		t.Errorf("expected %q, got %q", "abcdefghijklmno", lines[0].Text)
	}
}

func TestWrapNewlines(t *testing.T) {
	cfg := newWrapConfig(1000)
	cfg.WhiteSpace = wsPreLine
	lines := WrapText("hello\nworld", cfg)

	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d: %+v", len(lines), lines)
	}
	if lines[0].Text != "hello" {
		t.Errorf("line 0: expected %q, got %q", "hello", lines[0].Text)
	}
	if lines[1].Text != "world" {
		t.Errorf("line 1: expected %q, got %q", "world", lines[1].Text)
	}
}

func TestWrapBreakAll(t *testing.T) {
	cfg := newWrapConfig(50)
	cfg.WordBreak = wbBreakAll
	lines := WrapText("abcdefghij", cfg)

	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d: %+v", len(lines), lines)
	}
	if lines[0].Text != "abcde" {
		t.Errorf("line 0: expected %q, got %q", "abcde", lines[0].Text)
	}
	if lines[1].Text != "fghij" {
		t.Errorf("line 1: expected %q, got %q", "fghij", lines[1].Text)
	}
}

func TestWrapPreWrap(t *testing.T) {
	cfg := newWrapConfig(120)
	cfg.WhiteSpace = wsPreWrap
	lines := WrapText("hello  world\nfoo bar baz qux", cfg)

	if len(lines) < 2 {
		t.Fatalf("expected at least 2 lines, got %d: %+v", len(lines), lines)
	}
	if lines[0].Text != "hello  world" {
		t.Errorf("line 0: expected %q, got %q", "hello  world", lines[0].Text)
	}
}
