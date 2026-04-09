package render

import (
	"bytes"
	"testing"

	"github.com/macawls/ogre/style"
)

func TestRenderJPEG_MagicBytes(t *testing.T) {
	tree, styles := buildTestTree(&style.ComputedStyle{
		BackgroundColor: style.Color{R: 100, G: 150, B: 200, A: 1},
	}, 200, 100)

	data, err := RenderJPEG(tree, styles, nil, 200, 100, 90)
	if err != nil {
		t.Fatalf("RenderJPEG returned error: %v", err)
	}

	if len(data) < 2 {
		t.Fatal("output too short to contain JPEG header")
	}

	magic := []byte{0xFF, 0xD8}
	if !bytes.HasPrefix(data, magic) {
		t.Errorf("output does not start with JPEG magic bytes, got %x", data[:2])
	}
}

func TestRenderJPEG_ReasonableSize(t *testing.T) {
	tree, styles := buildTestTree(&style.ComputedStyle{
		BackgroundColor: style.Color{R: 100, G: 150, B: 200, A: 1},
	}, 400, 300)

	data, err := RenderJPEG(tree, styles, nil, 400, 300, 80)
	if err != nil {
		t.Fatalf("RenderJPEG returned error: %v", err)
	}

	if len(data) == 0 {
		t.Fatal("JPEG output is empty")
	}
	if len(data) > 500000 {
		t.Errorf("JPEG output unexpectedly large: %d bytes", len(data))
	}
}

func TestRenderJPEG_DefaultQuality(t *testing.T) {
	tree, styles := buildTestTree(&style.ComputedStyle{
		BackgroundColor: style.Color{R: 255, G: 0, B: 0, A: 1},
	}, 100, 100)

	data, err := RenderJPEG(tree, styles, nil, 100, 100, 0)
	if err != nil {
		t.Fatalf("RenderJPEG with quality 0 returned error: %v", err)
	}

	if len(data) < 2 {
		t.Fatal("output too short")
	}
	if data[0] != 0xFF || data[1] != 0xD8 {
		t.Error("invalid JPEG magic bytes with default quality")
	}
}
