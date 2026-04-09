package ogre

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandlerRender(t *testing.T) {
	r := NewRenderer()
	h := r.Handler(HandlerConfig{})

	body := `{"html":"<div style=\"background-color:red;width:100%;height:100%\">Hello</div>"}`
	req := httptest.NewRequest(http.MethodPost, "/og", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	if ct := rec.Header().Get("Content-Type"); ct != "image/png" {
		t.Errorf("expected image/png, got %s", ct)
	}
	if rec.Body.Len() == 0 {
		t.Error("expected non-empty body")
	}
}

func TestHandlerMissingHTML(t *testing.T) {
	r := NewRenderer()
	h := r.Handler(HandlerConfig{})

	req := httptest.NewRequest(http.MethodPost, "/og", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestHandlerSVGFormat(t *testing.T) {
	r := NewRenderer()
	h := r.Handler(HandlerConfig{Format: FormatSVG})

	body := `{"html":"<div style=\"background-color:blue;width:100%;height:100%\">Hi</div>"}`
	req := httptest.NewRequest(http.MethodPost, "/og", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	if ct := rec.Header().Get("Content-Type"); ct != "image/svg+xml" {
		t.Errorf("expected image/svg+xml, got %s", ct)
	}
}

func TestHandlerJPEGQuality(t *testing.T) {
	r := NewRenderer()
	h := r.Handler(HandlerConfig{Format: FormatJPEG, Quality: 50})

	body := `{"html":"<div style=\"background-color:red;width:100%;height:100%\">Hi</div>"}`
	req := httptest.NewRequest(http.MethodPost, "/og", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	if ct := rec.Header().Get("Content-Type"); ct != "image/jpeg" {
		t.Errorf("expected image/jpeg, got %s", ct)
	}
}

func TestHandlerQualityOverride(t *testing.T) {
	r := NewRenderer()
	h := r.Handler(HandlerConfig{Format: FormatJPEG, Quality: 90})

	body := `{"html":"<div style=\"background-color:red;width:100%;height:100%\">Hi</div>","quality":10}`
	req := httptest.NewRequest(http.MethodPost, "/og", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestHandlerErrorSurfaced(t *testing.T) {
	r := NewRenderer()
	h := r.Handler(HandlerConfig{})

	body := `{"template":"{{.Invalid"}`
	req := httptest.NewRequest(http.MethodPost, "/og", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
	resp := rec.Body.String()
	if !strings.Contains(resp, "invalid template") {
		t.Errorf("expected template error details, got: %s", resp)
	}
}

func TestHandlerInvalidJSON(t *testing.T) {
	r := NewRenderer()
	h := r.Handler(HandlerConfig{})

	req := httptest.NewRequest(http.MethodPost, "/og", strings.NewReader(`not json`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}
