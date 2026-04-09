package font

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// FontCache caches font data in memory and optionally on disk.
type FontCache struct {
	dir string
	mu  sync.RWMutex
	mem map[string][]byte
}

// NewFontCache creates a FontCache that stores files in the given directory.
func NewFontCache(cacheDir string) *FontCache {
	return &FontCache{
		dir: cacheDir,
		mem: make(map[string][]byte),
	}
}

func (c *FontCache) cacheKey(url string) string {
	h := sha256.Sum256([]byte(url))
	return fmt.Sprintf("%x", h)
}

// Fetch retrieves font data from the cache or downloads it from the URL.
func (c *FontCache) Fetch(url string) ([]byte, error) {
	key := c.cacheKey(url)

	c.mu.RLock()
	if data, ok := c.mem[key]; ok {
		c.mu.RUnlock()
		return data, nil
	}
	c.mu.RUnlock()

	if c.dir != "" {
		path := filepath.Join(c.dir, key)
		if data, err := os.ReadFile(path); err == nil {
			c.mu.Lock()
			c.mem[key] = data
			c.mu.Unlock()
			return data, nil
		}
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetch font %q: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch font %q: status %d", url, resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read font %q: %w", url, err)
	}

	c.mu.Lock()
	c.mem[key] = data
	c.mu.Unlock()

	if c.dir != "" {
		_ = os.MkdirAll(c.dir, 0o755)
		_ = os.WriteFile(filepath.Join(c.dir, key), data, 0o644)
	}

	return data, nil
}
