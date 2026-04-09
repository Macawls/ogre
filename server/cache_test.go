package server

import (
	"fmt"
	"sync"
	"testing"
)

func TestCacheGetSet(t *testing.T) {
	c := NewCache(1024)

	c.Set("a", []byte("hello"))
	v, ok := c.Get("a")
	if !ok {
		t.Fatal("expected key 'a' to exist")
	}
	if string(v) != "hello" {
		t.Fatalf("expected 'hello', got %q", v)
	}

	_, ok = c.Get("missing")
	if ok {
		t.Fatal("expected missing key to return false")
	}
}

func TestCacheOverwrite(t *testing.T) {
	c := NewCache(1024)

	c.Set("a", []byte("v1"))
	c.Set("a", []byte("v2"))

	v, ok := c.Get("a")
	if !ok || string(v) != "v2" {
		t.Fatalf("expected 'v2', got %q", v)
	}
	if c.Len() != 1 {
		t.Fatalf("expected len 1, got %d", c.Len())
	}
}

func TestCacheLRUEviction(t *testing.T) {
	// key "a" = 1 byte, value = 4 bytes => 5 bytes each
	c := NewCache(15)

	c.Set("a", []byte("aaaa")) // 5 bytes, total 5
	c.Set("b", []byte("bbbb")) // 5 bytes, total 10
	c.Set("c", []byte("cccc")) // 5 bytes, total 15

	// Adding d should evict a (oldest)
	c.Set("d", []byte("dddd")) // 5 bytes, would be 20, evict a => 15

	if _, ok := c.Get("a"); ok {
		t.Fatal("expected 'a' to be evicted")
	}
	if _, ok := c.Get("b"); !ok {
		t.Fatal("expected 'b' to exist")
	}
	if c.Len() != 3 {
		t.Fatalf("expected len 3, got %d", c.Len())
	}
}

func TestCacheAccessPattern(t *testing.T) {
	c := NewCache(15)

	c.Set("a", []byte("aaaa")) // 5
	c.Set("b", []byte("bbbb")) // 10
	c.Set("c", []byte("cccc")) // 15

	// Access a to move it to front; b is now the LRU
	c.Get("a")

	c.Set("d", []byte("dddd")) // evicts b

	if _, ok := c.Get("b"); ok {
		t.Fatal("expected 'b' to be evicted")
	}
	if _, ok := c.Get("a"); !ok {
		t.Fatal("expected 'a' to still exist")
	}
}

func TestCacheConcurrent(t *testing.T) {
	c := NewCache(10000)
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			key := fmt.Sprintf("key-%d", n)
			c.Set(key, []byte("value"))
			c.Get(key)
		}(i)
	}

	wg.Wait()

	if c.Len() < 0 {
		t.Fatal("unexpected negative length")
	}
}

func TestCacheZeroMaxBytes(t *testing.T) {
	c := NewCache(0)

	c.Set("a", []byte("hello"))

	if _, ok := c.Get("a"); ok {
		t.Fatal("expected no entries in zero-capacity cache")
	}
	if c.Len() != 0 {
		t.Fatalf("expected len 0, got %d", c.Len())
	}
}

func TestCacheEntryLargerThanMax(t *testing.T) {
	c := NewCache(5)

	c.Set("a", []byte("this is way too large"))

	if _, ok := c.Get("a"); ok {
		t.Fatal("expected oversized entry to not be stored")
	}
	if c.Len() != 0 {
		t.Fatalf("expected len 0, got %d", c.Len())
	}
}

func TestCacheEmptyGet(t *testing.T) {
	c := NewCache(1024)

	_, ok := c.Get("anything")
	if ok {
		t.Fatal("expected empty cache to return false")
	}
	if c.Len() != 0 {
		t.Fatalf("expected len 0, got %d", c.Len())
	}
	if c.Size() != 0 {
		t.Fatalf("expected size 0, got %d", c.Size())
	}
}

func TestCacheSizeTracking(t *testing.T) {
	c := NewCache(1024)

	c.Set("ab", []byte("cdef")) // key=2 + val=4 = 6
	if c.Size() != 6 {
		t.Fatalf("expected size 6, got %d", c.Size())
	}

	c.Set("ab", []byte("gh")) // key=2 + val=2 = 4
	if c.Size() != 4 {
		t.Fatalf("expected size 4 after overwrite, got %d", c.Size())
	}
}
