package font

import (
	"container/list"
	"sync"
)

// glyphCacheCapacity bounds the glyph path cache. ASCII + Latin-Extended +
// common emoji fit comfortably; CJK-heavy renders will evict but that avoids
// unbounded growth on servers rendering diverse text.
const glyphCacheCapacity = 4096

type glyphKey struct {
	fontName string
	r        rune
	size     float64
}

type glyphEntry struct {
	key  glyphKey
	path GlyphPath
}

type glyphCache struct {
	mu    sync.Mutex
	index map[glyphKey]*list.Element
	order *list.List
}

func newGlyphCache() *glyphCache {
	return &glyphCache{
		index: make(map[glyphKey]*list.Element),
		order: list.New(),
	}
}

func (c *glyphCache) Get(fontName string, r rune, size float64) (GlyphPath, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	el, ok := c.index[glyphKey{fontName, r, size}]
	if !ok {
		return GlyphPath{}, false
	}
	c.order.MoveToFront(el)
	return el.Value.(*glyphEntry).path, true
}

func (c *glyphCache) Set(fontName string, r rune, size float64, path GlyphPath) {
	key := glyphKey{fontName, r, size}
	c.mu.Lock()
	defer c.mu.Unlock()
	if el, ok := c.index[key]; ok {
		el.Value.(*glyphEntry).path = path
		c.order.MoveToFront(el)
		return
	}
	el := c.order.PushFront(&glyphEntry{key: key, path: path})
	c.index[key] = el
	if c.order.Len() > glyphCacheCapacity {
		oldest := c.order.Back()
		if oldest != nil {
			c.order.Remove(oldest)
			delete(c.index, oldest.Value.(*glyphEntry).key)
		}
	}
}
