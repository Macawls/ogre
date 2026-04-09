package render

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/macawls/ogre/style"
)

func TestRenderImageDataURI(t *testing.T) {
	cs := &style.ComputedStyle{ObjectFit: style.ObjectFitContain}
	src := "data:image/png;base64,iVBORw0KGgo="
	result := RenderImage(src, cs, 10, 20, 100, 50)

	if !strings.Contains(result, "<image") {
		t.Fatal("expected <image element")
	}
	if !strings.Contains(result, `href="data:image/png;base64,iVBORw0KGgo="`) {
		t.Fatalf("expected data URI in href, got: %s", result)
	}
	if !strings.Contains(result, `preserveAspectRatio="xMidYMid meet"`) {
		t.Fatalf("expected contain aspect ratio, got: %s", result)
	}
}

func TestRenderImageHTTP(t *testing.T) {
	imgData := []byte{0x89, 0x50, 0x4E, 0x47}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		w.Write(imgData)
	}))
	defer srv.Close()

	cs := &style.ComputedStyle{ObjectFit: style.ObjectFitCover}
	result := RenderImage(srv.URL+"/test.png", cs, 0, 0, 200, 100)

	expected := base64.StdEncoding.EncodeToString(imgData)
	if !strings.Contains(result, expected) {
		t.Fatalf("expected base64 encoded data, got: %s", result)
	}
	if !strings.Contains(result, `preserveAspectRatio="xMidYMid slice"`) {
		t.Fatalf("expected cover aspect ratio, got: %s", result)
	}
}

func TestRenderImageBroken(t *testing.T) {
	cs := &style.ComputedStyle{}
	result := RenderImage("http://127.0.0.1:1/nonexistent.png", cs, 10, 20, 100, 50)

	if !strings.Contains(result, "<rect") {
		t.Fatal("expected placeholder rect for broken image")
	}
	if !strings.Contains(result, `fill="#f0f0f0"`) {
		t.Fatalf("expected placeholder fill, got: %s", result)
	}
}

func TestObjectFitMapping(t *testing.T) {
	tests := []struct {
		fit    style.ObjectFit
		expect string
	}{
		{style.ObjectFitContain, "xMidYMid meet"},
		{style.ObjectFitCover, "xMidYMid slice"},
		{style.ObjectFitFill, "none"},
		{style.ObjectFitScaleDown, "xMidYMid meet"},
	}
	for _, tt := range tests {
		got := objectFitToPreserveAspectRatio(tt.fit)
		if got != tt.expect {
			t.Errorf("ObjectFit %v: expected %q, got %q", tt.fit, tt.expect, got)
		}
	}
}

func TestObjectPosition(t *testing.T) {
	ox, oy := parseObjectPosition("10px 20px", 200, 100)
	if ox != 10 || oy != 20 {
		t.Errorf("expected (10, 20), got (%.4g, %.4g)", ox, oy)
	}
}

func TestRenderImageHTTPFail(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	cs := &style.ComputedStyle{}
	result := RenderImage(srv.URL+"/missing.png", cs, 0, 0, 100, 50)

	if !strings.Contains(result, "<rect") {
		t.Fatal("expected broken image placeholder for 404")
	}
}

func TestRenderImagePositionOffset(t *testing.T) {
	cs := &style.ComputedStyle{
		ObjectFit:      style.ObjectFitContain,
		ObjectPosition: "10px 5px",
	}
	src := "data:image/png;base64,AAAA"
	result := RenderImage(src, cs, 100, 200, 50, 50)

	if !strings.Contains(result, fmt.Sprintf(`x="%.4g"`, 110.0)) {
		t.Fatalf("expected x offset of 110, got: %s", result)
	}
}
