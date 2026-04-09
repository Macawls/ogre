package font

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestFontCache_FetchFromServer(t *testing.T) {
	content := []byte("fake-font-data-for-testing")
	hits := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		w.Write(content)
	}))
	defer srv.Close()

	dir := t.TempDir()
	cache := NewFontCache(dir)

	data, err := cache.Fetch(srv.URL + "/test.ttf")
	if err != nil {
		t.Fatalf("first fetch: %v", err)
	}
	if string(data) != string(content) {
		t.Fatalf("got %q, want %q", data, content)
	}
	if hits != 1 {
		t.Fatalf("expected 1 server hit, got %d", hits)
	}

	data, err = cache.Fetch(srv.URL + "/test.ttf")
	if err != nil {
		t.Fatalf("second fetch: %v", err)
	}
	if string(data) != string(content) {
		t.Fatalf("got %q, want %q", data, content)
	}
	if hits != 1 {
		t.Fatalf("expected 1 server hit after mem cache, got %d", hits)
	}
}

func TestFontCache_DiskCacheHit(t *testing.T) {
	content := []byte("disk-cached-font-data")
	hits := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		w.Write(content)
	}))
	defer srv.Close()

	dir := t.TempDir()
	url := srv.URL + "/font.ttf"

	cache1 := NewFontCache(dir)
	if _, err := cache1.Fetch(url); err != nil {
		t.Fatalf("initial fetch: %v", err)
	}
	if hits != 1 {
		t.Fatalf("expected 1 hit, got %d", hits)
	}

	cache2 := NewFontCache(dir)
	data, err := cache2.Fetch(url)
	if err != nil {
		t.Fatalf("disk cache fetch: %v", err)
	}
	if string(data) != string(content) {
		t.Fatalf("got %q, want %q", data, content)
	}
	if hits != 1 {
		t.Fatalf("expected no new server hit, got %d", hits)
	}
}

func TestFontCache_NoDiskDir(t *testing.T) {
	content := []byte("mem-only-data")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(content)
	}))
	defer srv.Close()

	cache := NewFontCache("")
	data, err := cache.Fetch(srv.URL + "/font.ttf")
	if err != nil {
		t.Fatalf("fetch: %v", err)
	}
	if string(data) != string(content) {
		t.Fatalf("got %q, want %q", data, content)
	}
}

func TestFontCache_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	cache := NewFontCache("")
	_, err := cache.Fetch(srv.URL + "/missing.ttf")
	if err == nil {
		t.Fatal("expected error for 404 response")
	}
}

func TestFontCache_CacheKeyIsDeterministic(t *testing.T) {
	cache := NewFontCache("")
	k1 := cache.cacheKey("https://example.com/font.ttf")
	k2 := cache.cacheKey("https://example.com/font.ttf")
	if k1 != k2 {
		t.Fatalf("cache keys differ: %q vs %q", k1, k2)
	}
	k3 := cache.cacheKey("https://example.com/other.ttf")
	if k1 == k3 {
		t.Fatal("different URLs should produce different keys")
	}
}

func TestFontCache_DiskFileWritten(t *testing.T) {
	content := []byte("persist-me")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(content)
	}))
	defer srv.Close()

	dir := t.TempDir()
	cache := NewFontCache(dir)
	url := srv.URL + "/test.ttf"

	if _, err := cache.Fetch(url); err != nil {
		t.Fatalf("fetch: %v", err)
	}

	key := cache.cacheKey(url)
	stored, err := os.ReadFile(filepath.Join(dir, key))
	if err != nil {
		t.Fatalf("read disk cache: %v", err)
	}
	if string(stored) != string(content) {
		t.Fatalf("disk content %q, want %q", stored, content)
	}
}
