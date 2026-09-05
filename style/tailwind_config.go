package style

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"sort"
	"strconv"
	"sync"
)

// TailwindConfig is immutable after first Resolve. Mutating fields after
// Key() has been called silently produces stale cache hits; Clone() first.
type TailwindConfig struct {
	Colors  map[string]map[int]string    `json:"colors,omitempty"`
	Spacing map[string]string            `json:"spacing,omitempty"`
	Fonts   map[string]string            `json:"fonts,omitempty"`
	Screens map[string]string            `json:"screens,omitempty"`
	Extend  map[string]map[string]string `json:"extend,omitempty"`

	keyOnce sync.Once `json:"-"`
	cached  string    `json:"-"`
}

func (c *TailwindConfig) IsEmpty() bool {
	if c == nil {
		return true
	}
	return len(c.Colors) == 0 && len(c.Spacing) == 0 && len(c.Fonts) == 0 && len(c.Screens) == 0 && len(c.Extend) == 0
}

func (c *TailwindConfig) Clone() *TailwindConfig {
	if c == nil {
		return nil
	}
	out := &TailwindConfig{}
	if len(c.Colors) > 0 {
		out.Colors = make(map[string]map[int]string, len(c.Colors))
		for k, shades := range c.Colors {
			m := make(map[int]string, len(shades))
			for sk, sv := range shades {
				m[sk] = sv
			}
			out.Colors[k] = m
		}
	}
	if len(c.Spacing) > 0 {
		out.Spacing = make(map[string]string, len(c.Spacing))
		for k, v := range c.Spacing {
			out.Spacing[k] = v
		}
	}
	if len(c.Fonts) > 0 {
		out.Fonts = make(map[string]string, len(c.Fonts))
		for k, v := range c.Fonts {
			out.Fonts[k] = v
		}
	}
	if len(c.Screens) > 0 {
		out.Screens = make(map[string]string, len(c.Screens))
		for k, v := range c.Screens {
			out.Screens[k] = v
		}
	}
	if len(c.Extend) > 0 {
		out.Extend = make(map[string]map[string]string, len(c.Extend))
		for k, props := range c.Extend {
			m := make(map[string]string, len(props))
			for pk, pv := range props {
				m[pk] = pv
			}
			out.Extend[k] = m
		}
	}
	return out
}

func (c *TailwindConfig) Key() string {
	if c == nil {
		return ""
	}
	c.keyOnce.Do(func() {
		h := sha256.New()
		writeStringMap := func(label string, m map[string]string) {
			binary.Write(h, binary.LittleEndian, uint32(len(m)))
			h.Write([]byte(label))
			keys := make([]string, 0, len(m))
			for k := range m {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				h.Write([]byte(k))
				h.Write([]byte{0})
				h.Write([]byte(m[k]))
				h.Write([]byte{0})
			}
		}
		writeShadeMap := func(m map[string]map[int]string) {
			binary.Write(h, binary.LittleEndian, uint32(len(m)))
			names := make([]string, 0, len(m))
			for k := range m {
				names = append(names, k)
			}
			sort.Strings(names)
			for _, name := range names {
				h.Write([]byte(name))
				h.Write([]byte{0})
				shades := m[name]
				binary.Write(h, binary.LittleEndian, uint32(len(shades)))
				keys := make([]int, 0, len(shades))
				for k := range shades {
					keys = append(keys, k)
				}
				sort.Ints(keys)
				for _, k := range keys {
					h.Write([]byte(strconv.Itoa(k)))
					h.Write([]byte{0})
					h.Write([]byte(shades[k]))
					h.Write([]byte{0})
				}
			}
		}
		writeExtend := func(m map[string]map[string]string) {
			binary.Write(h, binary.LittleEndian, uint32(len(m)))
			names := make([]string, 0, len(m))
			for k := range m {
				names = append(names, k)
			}
			sort.Strings(names)
			for _, name := range names {
				h.Write([]byte(name))
				h.Write([]byte{0})
				writeStringMap("props", m[name])
			}
		}

		writeShadeMap(c.Colors)
		writeStringMap("spacing", c.Spacing)
		writeStringMap("fonts", c.Fonts)
		writeStringMap("screens", c.Screens)
		writeExtend(c.Extend)

		c.cached = hex.EncodeToString(h.Sum(nil))
	})
	return c.cached
}

func (c *TailwindConfig) spacingOverride(s string) (string, bool) {
	if c == nil {
		return "", false
	}
	v, ok := c.Spacing[s]
	return v, ok
}

func (c *TailwindConfig) colorOverride(palette string, shade int) (string, bool) {
	if c == nil {
		return "", false
	}
	shades, ok := c.Colors[palette]
	if !ok {
		return "", false
	}
	v, ok := shades[shade]
	return v, ok
}

func (c *TailwindConfig) fontOverride(key string) (string, bool) {
	if c == nil {
		return "", false
	}
	v, ok := c.Fonts[key]
	return v, ok
}

func (c *TailwindConfig) extendOverride(cls string) (map[string]string, bool) {
	if c == nil {
		return nil, false
	}
	m, ok := c.Extend[cls]
	return m, ok
}

var configCaches sync.Map

func configCacheFor(cfg *TailwindConfig, v TailwindVersion) *sync.Map {
	if cfg == nil {
		return cacheFor(v)
	}
	k := cfg.Key() + string(v)
	if existing, ok := configCaches.Load(k); ok {
		return existing.(*sync.Map)
	}
	fresh := &sync.Map{}
	actual, _ := configCaches.LoadOrStore(k, fresh)
	return actual.(*sync.Map)
}
