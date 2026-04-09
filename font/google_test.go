package font

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGoogleFontURL_Inter400(t *testing.T) {
	url := GoogleFontURL("Inter", 400)
	want := "https://fonts.googleapis.com/css2?family=Inter:wght@400"
	if url != want {
		t.Errorf("got %q, want %q", url, want)
	}
}

func TestGoogleFontURL_CaseInsensitive(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"inter", "https://fonts.googleapis.com/css2?family=Inter:wght@400"},
		{"INTER", "https://fonts.googleapis.com/css2?family=Inter:wght@400"},
		{"Inter", "https://fonts.googleapis.com/css2?family=Inter:wght@400"},
		{"Open Sans", "https://fonts.googleapis.com/css2?family=Open+Sans:wght@400"},
		{"open sans", "https://fonts.googleapis.com/css2?family=Open+Sans:wght@400"},
		{"OPEN SANS", "https://fonts.googleapis.com/css2?family=Open+Sans:wght@400"},
	}
	for _, tt := range tests {
		got := GoogleFontURL(tt.input, 400)
		if got != tt.want {
			t.Errorf("GoogleFontURL(%q, 400) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestGoogleFontURL_UnknownFamily(t *testing.T) {
	url := GoogleFontURL("NotARealFont", 400)
	expected := "https://fonts.googleapis.com/css2?family=NotARealFont:wght@400"
	if url != expected {
		t.Errorf("expected %q, got %q", expected, url)
	}
}

func TestGoogleFontURL_ArabicFont(t *testing.T) {
	url := GoogleFontURL("Noto Sans Arabic", 400)
	expected := "https://fonts.googleapis.com/css2?family=Noto+Sans+Arabic:wght@400"
	if url != expected {
		t.Errorf("expected %q, got %q", expected, url)
	}
}

func TestGoogleFontURL_Weight700(t *testing.T) {
	url := GoogleFontURL("Roboto", 700)
	want := "https://fonts.googleapis.com/css2?family=Roboto:wght@700"
	if url != want {
		t.Errorf("got %q, want %q", url, want)
	}
}

func TestParseFontURL(t *testing.T) {
	css := `@font-face {
  font-family: 'Inter';
  font-style: normal;
  font-weight: 400;
  src: url(https://fonts.gstatic.com/s/inter/v18/abc123.ttf) format('truetype');
}`
	got := parseFontURL(css)
	want := "https://fonts.gstatic.com/s/inter/v18/abc123.ttf"
	if got != want {
		t.Errorf("parseFontURL() = %q, want %q", got, want)
	}
}

func TestParseFontURL_NoMatch(t *testing.T) {
	got := parseFontURL("body { color: red; }")
	if got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestFetchGoogleFont_MockServer(t *testing.T) {
	fontData := []byte("fake-ttf-data")
	fontServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(fontData)
	}))
	defer fontServer.Close()

	cssBody := fmt.Sprintf(`@font-face {
  font-family: 'Inter';
  src: url(%s/inter.ttf) format('truetype');
}`, fontServer.URL)

	cssServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(cssBody))
	}))
	defer cssServer.Close()

	origMap := make(map[string]string)
	for k, v := range googleFontsMap {
		origMap[k] = v
	}

	googleFontsMap["testfont"] = "TestFont"
	defer func() {
		delete(googleFontsMap, "testfont")
	}()

	cache := NewFontCache("")

	cssURL := fmt.Sprintf("%s/css?family=TestFont:wght@400", cssServer.URL)
	key := cache.cacheKey(cssURL)
	cache.mu.Lock()
	cache.mem[key] = []byte(cssBody)
	cache.mu.Unlock()

	oldURL := GoogleFontURL("testfont", 400)
	_ = oldURL

	cache.mu.Lock()
	delete(cache.mem, key)
	realKey := cache.cacheKey(GoogleFontURL("testfont", 400))
	cache.mem[realKey] = []byte(cssBody)
	cache.mu.Unlock()

	data, err := FetchGoogleFont("testfont", 400, cache)
	if err != nil {
		t.Fatalf("FetchGoogleFont() error: %v", err)
	}
	if string(data) != string(fontData) {
		t.Errorf("got %q, want %q", data, fontData)
	}
}

func TestFetchGoogleFont_UnknownFamily(t *testing.T) {
	cache := NewFontCache("")
	_, err := FetchGoogleFont("NotARealFont", 400, cache)
	if err == nil {
		t.Error("expected error for unknown family")
	}
}
