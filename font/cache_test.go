package font

import "testing"

func TestGlyphCacheGetSet(t *testing.T) {
	c := newGlyphCache()
	p := GlyphPath{D: "M0 0", Advance: 5}
	c.Set("f", 'a', 16, p)
	got, ok := c.Get("f", 'a', 16)
	if !ok || got.D != "M0 0" || got.Advance != 5 {
		t.Fatalf("Get returned %+v ok=%v; want %+v true", got, ok, p)
	}
	if _, ok := c.Get("f", 'b', 16); ok {
		t.Fatal("Get returned hit for missing key")
	}
}

func TestGlyphCacheEviction(t *testing.T) {
	c := newGlyphCache()
	for i := 0; i < glyphCacheCapacity+10; i++ {
		c.Set("f", rune(i), 16, GlyphPath{Advance: float64(i)})
	}
	if c.order.Len() != glyphCacheCapacity {
		t.Fatalf("order len = %d; want %d", c.order.Len(), glyphCacheCapacity)
	}
	// The first 10 inserted should have been evicted.
	if _, ok := c.Get("f", rune(0), 16); ok {
		t.Fatal("expected rune 0 to be evicted")
	}
	// Recent entries should still be present.
	if _, ok := c.Get("f", rune(glyphCacheCapacity+5), 16); !ok {
		t.Fatal("expected recent entry to be present")
	}
}

func TestGlyphCacheMoveToFront(t *testing.T) {
	c := newGlyphCache()
	// Fill to capacity.
	for i := 0; i < glyphCacheCapacity; i++ {
		c.Set("f", rune(i), 16, GlyphPath{Advance: float64(i)})
	}
	// Touch entry 0 so it becomes most-recently-used.
	if _, ok := c.Get("f", rune(0), 16); !ok {
		t.Fatal("entry 0 missing before touch")
	}
	// Insert one more; the least-recently-used (now rune 1, not 0) should evict.
	c.Set("f", rune(glyphCacheCapacity), 16, GlyphPath{})
	if _, ok := c.Get("f", rune(0), 16); !ok {
		t.Fatal("touched entry 0 was evicted; MoveToFront failed")
	}
	if _, ok := c.Get("f", rune(1), 16); ok {
		t.Fatal("entry 1 should have been evicted as LRU")
	}
}

func BenchmarkGlyphCacheHit(b *testing.B) {
	c := newGlyphCache()
	c.Set("f", 'a', 16, GlyphPath{D: "M0 0", Advance: 5})
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, ok := c.Get("f", 'a', 16); !ok {
			b.Fatal("miss")
		}
	}
}

func BenchmarkGlyphCacheSet(b *testing.B) {
	c := newGlyphCache()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		c.Set("f", rune(i&0xFFFF), 16, GlyphPath{Advance: float64(i)})
	}
}
