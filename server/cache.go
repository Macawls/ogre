package server

import (
	"container/list"
	"sync"
)

type cacheEntry struct {
	key   string
	value []byte
}

func (e *cacheEntry) size() int64 {
	return int64(len(e.key)) + int64(len(e.value))
}

// Cache is a thread-safe LRU byte cache with a configurable size limit.
type Cache struct {
	maxBytes  int64
	usedBytes int64
	mu        sync.Mutex
	ll        *list.List
	items     map[string]*list.Element
}

// NewCache creates a Cache with the given maximum byte capacity.
func NewCache(maxBytes int64) *Cache {
	return &Cache{
		maxBytes: maxBytes,
		ll:       list.New(),
		items:    make(map[string]*list.Element),
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if el, ok := c.items[key]; ok {
		c.ll.MoveToFront(el)
		return el.Value.(*cacheEntry).value, true
	}
	return nil, false
}

func (c *Cache) Set(key string, value []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry := &cacheEntry{key: key, value: value}

	if entry.size() > c.maxBytes {
		return
	}

	if el, ok := c.items[key]; ok {
		old := el.Value.(*cacheEntry)
		c.usedBytes -= old.size()
		el.Value = entry
		c.usedBytes += entry.size()
		c.ll.MoveToFront(el)
	} else {
		el := c.ll.PushFront(entry)
		c.items[key] = el
		c.usedBytes += entry.size()
	}

	for c.usedBytes > c.maxBytes && c.ll.Len() > 0 {
		back := c.ll.Back()
		evicted := back.Value.(*cacheEntry)
		c.ll.Remove(back)
		delete(c.items, evicted.key)
		c.usedBytes -= evicted.size()
	}
}

func (c *Cache) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.ll.Len()
}

func (c *Cache) Size() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.usedBytes
}
