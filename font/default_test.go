package font

import "testing"

func TestLoadDefaults(t *testing.T) {
	m := NewManager()
	if err := m.LoadDefaults(); err != nil {
		t.Fatalf("LoadDefaults: %v", err)
	}

	tests := []struct {
		family string
		weight int
		style  string
	}{
		{"sans-serif", 400, "normal"},
		{"sans-serif", 700, "normal"},
		{"default", 400, "normal"},
	}

	for _, tt := range tests {
		face := m.Resolve(tt.family, tt.weight, tt.style)
		if face == nil {
			t.Errorf("Resolve(%q, %d, %q) returned nil", tt.family, tt.weight, tt.style)
		}
	}
}
