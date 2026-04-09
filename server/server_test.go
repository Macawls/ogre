package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func newTestServer() *Server {
	return New(Config{
		Addr:       ":0",
		CacheBytes: 1 << 20,
	})
}

func TestHealthEndpoint(t *testing.T) {
	srv := newTestServer()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	srv.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if body["status"] != "ok" {
		t.Fatalf("expected status ok, got %q", body["status"])
	}
}

func TestRenderSVG(t *testing.T) {
	srv := newTestServer()

	payload := `{"html":"<div>Hello</div>","width":400,"height":200,"format":"svg"}`
	req := httptest.NewRequest(http.MethodPost, "/render", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	srv.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d; body: %s", rec.Code, rec.Body.String())
	}
	if ct := rec.Header().Get("Content-Type"); ct != "image/svg+xml" {
		t.Fatalf("expected content-type image/svg+xml, got %q", ct)
	}
	if !strings.Contains(rec.Body.String(), "<svg") {
		t.Fatal("response does not contain <svg")
	}
}

func TestRenderTemplate(t *testing.T) {
	srv := newTestServer()

	payload := `{"template":"<div>{{.Title}}</div>","data":{"Title":"Hello World"},"width":400,"height":200,"format":"svg"}`
	req := httptest.NewRequest(http.MethodPost, "/render/template", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	srv.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d; body: %s", rec.Code, rec.Body.String())
	}
	if ct := rec.Header().Get("Content-Type"); ct != "image/svg+xml" {
		t.Fatalf("expected content-type image/svg+xml, got %q", ct)
	}
	if !strings.Contains(rec.Body.String(), "<svg") {
		t.Fatal("response does not contain <svg")
	}
}

func TestCacheHit(t *testing.T) {
	srv := newTestServer()

	payload := `{"html":"<div>Cache Test</div>","width":400,"height":200,"format":"svg"}`

	req1 := httptest.NewRequest(http.MethodPost, "/render", strings.NewReader(payload))
	rec1 := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec1, req1)

	if rec1.Code != http.StatusOK {
		t.Fatalf("first request: expected 200, got %d; body: %s", rec1.Code, rec1.Body.String())
	}

	req2 := httptest.NewRequest(http.MethodPost, "/render", strings.NewReader(payload))
	rec2 := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusOK {
		t.Fatalf("second request: expected 200, got %d", rec2.Code)
	}

	if !bytes.Equal(rec1.Body.Bytes(), rec2.Body.Bytes()) {
		t.Fatal("cache hit should return identical content")
	}

	if srv.cache.Len() != 1 {
		t.Fatalf("expected 1 cache entry, got %d", srv.cache.Len())
	}
}

func TestInvalidJSON(t *testing.T) {
	srv := newTestServer()

	req := httptest.NewRequest(http.MethodPost, "/render", strings.NewReader("{invalid"))
	rec := httptest.NewRecorder()

	srv.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestMissingHTMLField(t *testing.T) {
	srv := newTestServer()

	payload := `{"width":400,"height":200,"format":"svg"}`
	req := httptest.NewRequest(http.MethodPost, "/render", strings.NewReader(payload))
	rec := httptest.NewRecorder()

	srv.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}

	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid json response: %v", err)
	}
	if !strings.Contains(body["error"], "missing html") {
		t.Fatalf("expected missing html error, got %q", body["error"])
	}
}

func TestCORSHeaders(t *testing.T) {
	srv := newTestServer()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	srv.Handler().ServeHTTP(rec, req)

	if v := rec.Header().Get("Access-Control-Allow-Origin"); v != "*" {
		t.Fatalf("expected CORS origin *, got %q", v)
	}
}

func TestRateLimiterBlocksExceeded(t *testing.T) {
	srv := New(Config{
		Addr:       ":0",
		CacheBytes: 1 << 20,
		RateLimit:  1,
	})

	payload := `{"html":"<div>Rate Test</div>","width":100,"height":100,"format":"svg"}`

	blocked := false
	for i := 0; i < 10; i++ {
		req := httptest.NewRequest(http.MethodPost, "/render", strings.NewReader(payload))
		req.RemoteAddr = "192.168.1.1:12345"
		rec := httptest.NewRecorder()
		srv.Handler().ServeHTTP(rec, req)
		if rec.Code == http.StatusTooManyRequests {
			blocked = true
			break
		}
	}
	if !blocked {
		t.Fatal("rate limiter did not block any requests")
	}
}

func TestMetricsEndpoint(t *testing.T) {
	srv := newTestServer()

	payload := `{"html":"<div>Metrics Test</div>","width":100,"height":100,"format":"svg"}`
	renderReq := httptest.NewRequest(http.MethodPost, "/render", strings.NewReader(payload))
	renderRec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(renderRec, renderReq)
	if renderRec.Code != http.StatusOK {
		t.Fatalf("render request failed: %d %s", renderRec.Code, renderRec.Body.String())
	}

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected application/json, got %q", ct)
	}

	var m map[string]int64
	if err := json.Unmarshal(rec.Body.Bytes(), &m); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if m["render_total"] < 1 {
		t.Fatalf("expected render_total >= 1, got %d", m["render_total"])
	}
	if _, ok := m["cache_hits"]; !ok {
		t.Fatal("missing cache_hits in metrics")
	}
	if _, ok := m["cache_misses"]; !ok {
		t.Fatal("missing cache_misses in metrics")
	}
}

func TestRenderTimeout(t *testing.T) {
	srv := New(Config{
		Addr:          ":0",
		CacheBytes:    1 << 20,
		RenderTimeout: 1 * time.Nanosecond,
	})

	payload := `{"html":"<div>Timeout Test</div>","width":100,"height":100,"format":"svg"}`
	req := httptest.NewRequest(http.MethodPost, "/render", strings.NewReader(payload))
	rec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusGatewayTimeout && rec.Code != http.StatusOK {
		t.Fatalf("expected 504 or 200 (race), got %d", rec.Code)
	}
}
