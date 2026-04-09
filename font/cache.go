package font

import "sync"

type glyphKey struct {
	fontName string
	r        rune
	size     float64
}

type glyphCache struct {
	mu    sync.RWMutex
	paths map[glyphKey]GlyphPath
}

func newGlyphCache() *glyphCache {
	return &glyphCache{paths: make(map[glyphKey]GlyphPath)}
}

func (c *glyphCache) Get(fontName string, r rune, size float64) (GlyphPath, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	p, ok := c.paths[glyphKey{fontName, r, size}]
	return p, ok
}

func (c *glyphCache) Set(fontName string, r rune, size float64, path GlyphPath) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.paths[glyphKey{fontName, r, size}] = path
}
