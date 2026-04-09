package font

import (
	"testing"
)

func makeManager(faces map[string][]*Face) *Manager {
	m := NewManager()
	m.families = faces
	return m
}

func TestNewManager(t *testing.T) {
	m := NewManager()
	if m.families == nil {
		t.Fatal("families map not initialized")
	}
	if len(m.Families()) != 0 {
		t.Fatal("expected no families")
	}
}

func TestResolveExactMatch(t *testing.T) {
	m := makeManager(map[string][]*Face{
		"sans": {
			{Name: "sans", Weight: 400, Style: "normal"},
			{Name: "sans", Weight: 700, Style: "normal"},
		},
	})

	f := m.Resolve("sans", 400, "normal")
	if f == nil || f.Weight != 400 {
		t.Fatal("expected exact match at 400")
	}

	f = m.Resolve("sans", 700, "normal")
	if f == nil || f.Weight != 700 {
		t.Fatal("expected exact match at 700")
	}
}

func TestResolve400And500Interchangeable(t *testing.T) {
	m := makeManager(map[string][]*Face{
		"sans": {
			{Name: "sans", Weight: 500, Style: "normal"},
		},
	})

	f := m.Resolve("sans", 400, "normal")
	if f == nil || f.Weight != 500 {
		t.Fatal("expected 500 when requesting 400")
	}

	m2 := makeManager(map[string][]*Face{
		"sans": {
			{Name: "sans", Weight: 400, Style: "normal"},
		},
	})

	f = m2.Resolve("sans", 500, "normal")
	if f == nil || f.Weight != 400 {
		t.Fatal("expected 400 when requesting 500")
	}
}

func TestResolveBelow400PrefersLighter(t *testing.T) {
	m := makeManager(map[string][]*Face{
		"sans": {
			{Name: "sans", Weight: 100, Style: "normal"},
			{Name: "sans", Weight: 200, Style: "normal"},
			{Name: "sans", Weight: 600, Style: "normal"},
		},
	})

	f := m.Resolve("sans", 300, "normal")
	if f == nil || f.Weight != 200 {
		t.Fatalf("expected 200 (lighter), got %d", f.Weight)
	}
}

func TestResolveBelow400FallsToHeavier(t *testing.T) {
	m := makeManager(map[string][]*Face{
		"sans": {
			{Name: "sans", Weight: 600, Style: "normal"},
			{Name: "sans", Weight: 700, Style: "normal"},
		},
	})

	f := m.Resolve("sans", 300, "normal")
	if f == nil || f.Weight != 600 {
		t.Fatalf("expected 600 (nearest heavier), got %d", f.Weight)
	}
}

func TestResolveAbove500PrefersHeavier(t *testing.T) {
	m := makeManager(map[string][]*Face{
		"sans": {
			{Name: "sans", Weight: 400, Style: "normal"},
			{Name: "sans", Weight: 700, Style: "normal"},
			{Name: "sans", Weight: 900, Style: "normal"},
		},
	})

	f := m.Resolve("sans", 600, "normal")
	if f == nil || f.Weight != 700 {
		t.Fatalf("expected 700 (nearest heavier), got %d", f.Weight)
	}
}

func TestResolveAbove500FallsToLighter(t *testing.T) {
	m := makeManager(map[string][]*Face{
		"sans": {
			{Name: "sans", Weight: 300, Style: "normal"},
			{Name: "sans", Weight: 400, Style: "normal"},
		},
	})

	f := m.Resolve("sans", 600, "normal")
	if f == nil || f.Weight != 400 {
		t.Fatalf("expected 400 (nearest lighter), got %d", f.Weight)
	}
}

func TestResolveFallbackToDefault(t *testing.T) {
	m := makeManager(map[string][]*Face{
		"default": {
			{Name: "default", Weight: 400, Style: "normal"},
		},
	})

	f := m.Resolve("missing", 400, "normal")
	if f == nil || f.Name != "default" {
		t.Fatal("expected fallback to default family")
	}
}

func TestResolveNilWhenEmpty(t *testing.T) {
	m := NewManager()
	f := m.Resolve("missing", 400, "normal")
	if f != nil {
		t.Fatal("expected nil when no families loaded")
	}
}

func TestFamilies(t *testing.T) {
	m := makeManager(map[string][]*Face{
		"sans":  {{Name: "sans", Weight: 400, Style: "normal"}},
		"serif": {{Name: "serif", Weight: 400, Style: "normal"}},
	})

	fams := m.Families()
	if len(fams) != 2 {
		t.Fatalf("expected 2 families, got %d", len(fams))
	}
	if fams[0] != "sans" || fams[1] != "serif" {
		t.Fatalf("expected [sans serif], got %v", fams)
	}
}

func TestResolveStylePreference(t *testing.T) {
	m := makeManager(map[string][]*Face{
		"sans": {
			{Name: "sans", Weight: 400, Style: "normal"},
			{Name: "sans", Weight: 400, Style: "italic"},
			{Name: "sans", Weight: 700, Style: "normal"},
			{Name: "sans", Weight: 700, Style: "italic"},
		},
	})

	f := m.Resolve("sans", 400, "italic")
	if f == nil || f.Style != "italic" || f.Weight != 400 {
		t.Fatal("expected italic 400 exact match")
	}
}

func TestLoadFontInvalidData(t *testing.T) {
	m := NewManager()
	err := m.LoadFont(FontSource{
		Name: "bad",
		Data: []byte("not a font"),
	})
	if err == nil {
		t.Fatal("expected error for invalid font data")
	}
}
