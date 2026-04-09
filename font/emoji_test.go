package font

import "testing"

func TestIsEmoji_ZWJ(t *testing.T) {
	if !IsEmoji(0x200D) {
		t.Fatal("ZWJ (U+200D) should be treated as emoji")
	}
}

func TestSplitEmoji_ZWJSequence(t *testing.T) {
	// 👨‍💻 = U+1F468 U+200D U+1F4BB
	input := "\U0001F468\u200D\U0001F4BB"
	segments := SplitEmoji(input)
	if len(segments) != 1 {
		t.Fatalf("expected 1 segment for ZWJ sequence, got %d: %v", len(segments), segments)
	}
	if !segments[0].IsEmoji {
		t.Fatal("ZWJ sequence should be marked as emoji")
	}
	if segments[0].Text != input {
		t.Fatalf("expected %q, got %q", input, segments[0].Text)
	}
}

func TestSplitEmoji_FamilyZWJ(t *testing.T) {
	// 👨‍👩‍👧‍👦 = U+1F468 U+200D U+1F469 U+200D U+1F467 U+200D U+1F466
	input := "\U0001F468\u200D\U0001F469\u200D\U0001F467\u200D\U0001F466"
	segments := SplitEmoji(input)
	if len(segments) != 1 {
		t.Fatalf("expected 1 segment for family ZWJ sequence, got %d", len(segments))
	}
	if !segments[0].IsEmoji {
		t.Fatal("family ZWJ sequence should be marked as emoji")
	}
}

func TestSplitEmoji_MixedTextAndEmoji(t *testing.T) {
	input := "Hello 👋 World"
	segments := SplitEmoji(input)
	if len(segments) != 3 {
		t.Fatalf("expected 3 segments, got %d: %v", len(segments), segments)
	}
	if segments[0].IsEmoji || segments[0].Text != "Hello " {
		t.Fatalf("segment 0: expected text 'Hello ', got %q (emoji=%v)", segments[0].Text, segments[0].IsEmoji)
	}
	if !segments[1].IsEmoji {
		t.Fatal("segment 1 should be emoji")
	}
	if segments[2].IsEmoji || segments[2].Text != " World" {
		t.Fatalf("segment 2: expected text ' World', got %q (emoji=%v)", segments[2].Text, segments[2].IsEmoji)
	}
}
