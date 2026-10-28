package font

import (
	"math"
	"strconv"
	"strings"
	"testing"
)

// pathXRange returns the horizontal extent of an absolute-coordinate SVG path.
// Ogre emits only M/L/Q/C/H/V/Z in absolute form, which is all this handles.
func pathXRange(t *testing.T, d string) (minX, maxX float64) {
	t.Helper()

	minX, maxX = math.Inf(1), math.Inf(-1)
	seen := false

	addX := func(x float64) {
		seen = true
		minX = math.Min(minX, x)
		maxX = math.Max(maxX, x)
	}

	// Split the path into segments, each one command letter plus its numbers.
	var seg strings.Builder
	segments := make([]string, 0, 64)
	for _, r := range d {
		if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
			if seg.Len() > 0 {
				segments = append(segments, seg.String())
				seg.Reset()
			}
		}
		seg.WriteRune(r)
	}
	if seg.Len() > 0 {
		segments = append(segments, seg.String())
	}

	for _, s := range segments {
		cmd := rune(s[0])
		fields := strings.FieldsFunc(s[1:], func(r rune) bool {
			return r == ' ' || r == ',' || r == '\t' || r == '\n'
		})
		nums := make([]float64, 0, len(fields))
		for _, f := range fields {
			if f == "" {
				continue
			}
			if v, err := strconv.ParseFloat(f, 64); err == nil {
				nums = append(nums, v)
			}
		}
		if len(nums) == 0 {
			continue
		}

		switch cmd {
		case 'M', 'L', 'T':
			for i := 0; i+1 < len(nums); i += 2 {
				addX(nums[i])
			}
		case 'Q', 'S':
			for i := 0; i+3 < len(nums); i += 4 {
				addX(nums[i])
				addX(nums[i+2])
			}
		case 'C':
			for i := 0; i+5 < len(nums); i += 6 {
				addX(nums[i])
				addX(nums[i+2])
				addX(nums[i+4])
			}
		case 'H':
			for _, v := range nums {
				addX(v)
			}
		default:
			t.Fatalf("unsupported path command %q in %q", string(cmd), s)
		}
	}

	if !seen {
		t.Fatalf("no coordinates parsed from path %q", d[:minInt(60, len(d))])
	}
	return minX, maxX
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func TestTextToPathVariationsSpaced_ZeroSpacingMatchesUnspaced(t *testing.T) {
	mgr := NewManager()
	if err := mgr.LoadDefaults(); err != nil {
		t.Fatalf("load defaults: %v", err)
	}

	plain, plainAdv := TextToPathVariations(mgr, "ABCDEFGHIJ", "sans-serif", 400, "normal", 20, nil)
	spaced, spacedAdv := TextToPathVariationsSpaced(mgr, "ABCDEFGHIJ", "sans-serif", 400, "normal", 20, nil, 0)

	if spaced != plain {
		t.Error("zero letter-spacing produced a different path than the unspaced call")
	}
	if math.Abs(spacedAdv-plainAdv) > 0.01 {
		t.Errorf("zero letter-spacing advance: got %.2f, want %.2f", spacedAdv, plainAdv)
	}
}

func TestTextToPathVariationsSpaced_WidensPath(t *testing.T) {
	mgr := NewManager()
	if err := mgr.LoadDefaults(); err != nil {
		t.Fatalf("load defaults: %v", err)
	}

	const text = "ABCDEFGHIJ"
	const spacing = 10.0

	plain, plainAdv := TextToPathVariations(mgr, text, "sans-serif", 400, "normal", 20, nil)
	spaced, spacedAdv := TextToPathVariationsSpaced(mgr, text, "sans-serif", 400, "normal", 20, nil, spacing)

	if spaced == plain {
		t.Fatal("letter-spacing did not change the rendered path")
	}

	// The convention matches Measurer.StringWidth: spacing sits between
	// characters, so n-1 gaps for n characters.
	wantAdv := plainAdv + spacing*float64(len([]rune(text))-1)
	if math.Abs(spacedAdv-wantAdv) > 0.01 {
		t.Errorf("spaced advance: got %.2f, want %.2f", spacedAdv, wantAdv)
	}

	plainMinX, plainMaxX := pathXRange(t, plain)
	spacedMinX, spacedMaxX := pathXRange(t, spaced)

	plainW := plainMaxX - plainMinX
	spacedW := spacedMaxX - spacedMinX
	if math.Abs(spacedW-plainW-(wantAdv-plainAdv)) > 0.5 {
		t.Errorf("path grew by %.2f, want %.2f", spacedW-plainW, wantAdv-plainAdv)
	}
}

func TestTextToPathVariationsSpaced_SingleCharUnaffected(t *testing.T) {
	mgr := NewManager()
	if err := mgr.LoadDefaults(); err != nil {
		t.Fatalf("load defaults: %v", err)
	}

	plain, _ := TextToPathVariations(mgr, "A", "sans-serif", 400, "normal", 20, nil)
	spaced, _ := TextToPathVariationsSpaced(mgr, "A", "sans-serif", 400, "normal", 20, nil, 12)

	if plain != spaced {
		t.Error("a single character must not be shifted by letter-spacing")
	}
}
