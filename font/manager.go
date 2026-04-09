// Package font handles font loading, text measurement, and glyph path extraction.
package font

import (
	"fmt"
	"math"
	"sort"
	"sync"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

// FontSource describes a font to register, with raw data and metadata.
type FontSource struct {
	Name   string
	Weight int
	Style  string
	Data   []byte
	URL    string
}

// Face holds a parsed OpenType font along with its name, weight, and style.
type Face struct {
	Font    *opentype.Font
	RawData []byte
	Name    string
	Weight  int
	Style   string
}

type faceKey struct {
	name   string
	weight int
	style  string
	size   float64
}

// Manager stores loaded font families and provides font resolution and face creation.
type Manager struct {
	mu       sync.RWMutex
	families map[string][]*Face
	glyphs   *glyphCache
	facesMu  sync.RWMutex
	faces    map[faceKey]font.Face
}

// NewManager creates an empty font Manager.
func NewManager() *Manager {
	return &Manager{
		families: make(map[string][]*Face),
		glyphs:   newGlyphCache(),
		faces:    make(map[faceKey]font.Face),
	}
}

// LoadFont parses and registers a font from the given source.
func (m *Manager) LoadFont(src FontSource) error {
	data := src.Data
	if IsWOFF2(data) {
		return fmt.Errorf("WOFF2 not yet supported, convert to TTF/OTF")
	}
	if IsWOFF(data) {
		var err error
		data, err = DecompressWOFF(data)
		if err != nil {
			return fmt.Errorf("decompress WOFF %q: %w", src.Name, err)
		}
	}

	f, err := opentype.Parse(data)
	if err != nil {
		return fmt.Errorf("parse font %q: %w", src.Name, err)
	}

	weight := src.Weight
	if weight == 0 {
		weight = 400
	}

	style := src.Style
	if style == "" {
		style = "normal"
	}

	face := &Face{
		Font:    f,
		RawData: data,
		Name:    src.Name,
		Weight:  weight,
		Style:   style,
	}

	m.mu.Lock()
	m.families[src.Name] = append(m.families[src.Name], face)
	m.mu.Unlock()
	return nil
}

// Resolve finds the best matching Face for the given family, weight, and style.
func (m *Manager) Resolve(family string, weight int, style string) *Face {
	m.mu.RLock()
	defer m.mu.RUnlock()
	faces := m.families[family]
	if len(faces) == 0 {
		faces = m.families["default"]
	}
	if len(faces) == 0 {
		return nil
	}

	for _, f := range faces {
		if f.Weight == weight && f.Style == style {
			return f
		}
	}

	if weight == 400 {
		for _, f := range faces {
			if f.Weight == 500 && f.Style == style {
				return f
			}
		}
	}
	if weight == 500 {
		for _, f := range faces {
			if f.Weight == 400 && f.Style == style {
				return f
			}
		}
	}

	styleFaces := make([]*Face, 0)
	for _, f := range faces {
		if f.Style == style {
			styleFaces = append(styleFaces, f)
		}
	}
	if len(styleFaces) == 0 {
		styleFaces = faces
	}

	sort.Slice(styleFaces, func(i, j int) bool {
		return styleFaces[i].Weight < styleFaces[j].Weight
	})

	if weight < 400 {
		var best *Face
		minDist := math.MaxInt
		for _, f := range styleFaces {
			if f.Weight <= weight {
				d := weight - f.Weight
				if d < minDist {
					minDist = d
					best = f
				}
			}
		}
		if best != nil {
			return best
		}
		for _, f := range styleFaces {
			if f.Weight > weight {
				return f
			}
		}
	}

	var best *Face
	minDist := math.MaxInt
	for _, f := range styleFaces {
		if f.Weight >= weight {
			d := f.Weight - weight
			if d < minDist {
				minDist = d
				best = f
			}
		}
	}
	if best != nil {
		return best
	}
	for i := len(styleFaces) - 1; i >= 0; i-- {
		if styleFaces[i].Weight < weight {
			return styleFaces[i]
		}
	}

	return styleFaces[0]
}

// NewFace returns a font.Face for the given Face at the specified size, using a cache.
func (m *Manager) NewFace(f *Face, size float64) (font.Face, error) {
	key := faceKey{
		name:   f.Name,
		weight: f.Weight,
		style:  f.Style,
		size:   size,
	}

	m.facesMu.RLock()
	if cached, ok := m.faces[key]; ok {
		m.facesMu.RUnlock()
		return cached, nil
	}
	m.facesMu.RUnlock()

	face, err := opentype.NewFace(f.Font, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingNone,
	})
	if err != nil {
		return nil, err
	}

	m.facesMu.Lock()
	m.faces[key] = face
	m.facesMu.Unlock()
	return face, nil
}

func (m *Manager) CachedGlyphPath(fontName string, r rune, size float64, f *opentype.Font) (GlyphPath, error) {
	if p, ok := m.glyphs.Get(fontName, r, size); ok {
		return p, nil
	}
	p, err := GlyphToPath(f, r, size)
	if err != nil {
		return p, err
	}
	m.glyphs.Set(fontName, r, size, p)
	return p, nil
}

// HasFamily reports whether a font family with the given name is loaded.
func (m *Manager) HasFamily(name string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.families[name]) > 0
}

// Families returns a sorted list of all loaded font family names.
func (m *Manager) Families() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	names := make([]string, 0, len(m.families))
	for name := range m.families {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
