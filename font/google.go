package font

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var googleFontsMap = map[string]string{
	"inter":                "Inter",
	"roboto":               "Roboto",
	"open sans":            "Open+Sans",
	"lato":                 "Lato",
	"montserrat":           "Montserrat",
	"poppins":              "Poppins",
	"raleway":              "Raleway",
	"nunito":               "Nunito",
	"playfair display":     "Playfair+Display",
	"merriweather":         "Merriweather",
	"source sans pro":      "Source+Sans+Pro",
	"source code pro":      "Source+Code+Pro",
	"fira code":            "Fira+Code",
	"jetbrains mono":       "JetBrains+Mono",
	"dm sans":              "DM+Sans",
	"dm mono":              "DM+Mono",
	"geist":                "Geist",
	"geist mono":           "Geist+Mono",
	"ibm plex sans":        "IBM+Plex+Sans",
	"ibm plex mono":        "IBM+Plex+Mono",
	"work sans":            "Work+Sans",
	"space grotesk":        "Space+Grotesk",
	"space mono":           "Space+Mono",
	"ubuntu":               "Ubuntu",
	"ubuntu mono":          "Ubuntu+Mono",
	"noto sans":            "Noto+Sans",
	"noto serif":           "Noto+Serif",
	"outfit":               "Outfit",
	"manrope":              "Manrope",
	"plus jakarta sans":    "Plus+Jakarta+Sans",
	"bricolage grotesque":  "Bricolage+Grotesque",
	"sora":                 "Sora",
	"lexend":               "Lexend",
	"archivo":              "Archivo",
	"rubik":                "Rubik",
	"karla":                "Karla",
	"cabin":                "Cabin",
	"josefin sans":         "Josefin+Sans",
	"quicksand":            "Quicksand",
	"barlow":               "Barlow",
	"mulish":               "Mulish",
}

var fontURLRegexp = regexp.MustCompile(`src:\s*url\(([^)]+)\)`)

// GoogleFontURL returns the Google Fonts CSS URL for the given family and weight, or empty if unknown.
// GoogleFontURL builds a Google Fonts CSS API URL.
func GoogleFontURL(family string, weight int) string {
	slug, ok := googleFontsMap[strings.ToLower(family)]
	if !ok {
		slug = strings.ReplaceAll(family, " ", "+")
	}
	return fmt.Sprintf("https://fonts.googleapis.com/css2?family=%s:wght@%d", slug, weight)
}

func parseFontURL(css string) string {
	m := fontURLRegexp.FindStringSubmatch(css)
	if len(m) < 2 {
		return ""
	}
	return strings.TrimSpace(m[1])
}

// FetchGoogleFont downloads a Google Font by family and weight, using the cache.
// FetchGoogleFont downloads a font from Google Fonts CDN.
func FetchGoogleFont(family string, weight int, cache *FontCache) ([]byte, error) {
	cssURL := GoogleFontURL(family, weight)
	if cssURL == "" {
		return nil, fmt.Errorf("unknown Google Font family: %q", family)
	}

	cssKey := cssURL
	var cssBody []byte

	cache.mu.RLock()
	if data, ok := cache.mem[cache.cacheKey(cssKey)]; ok {
		cssBody = data
	}
	cache.mu.RUnlock()

	if cssBody == nil {
		client := &http.Client{Timeout: 10 * time.Second}
		req, err := http.NewRequest("GET", cssURL, nil)
		if err != nil {
			return nil, fmt.Errorf("build CSS request: %w", err)
		}
		req.Header.Set("User-Agent", "ogre/1.0")

		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("fetch Google Fonts CSS for %q: %w", family, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("fetch Google Fonts CSS for %q: status %d", family, resp.StatusCode)
		}

		cssBody, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("read Google Fonts CSS for %q: %w", family, err)
		}

		key := cache.cacheKey(cssKey)
		cache.mu.Lock()
		cache.mem[key] = cssBody
		cache.mu.Unlock()
	}

	fontURL := parseFontURL(string(cssBody))
	if fontURL == "" {
		return nil, fmt.Errorf("no font URL found in Google Fonts CSS for %q", family)
	}

	return cache.Fetch(fontURL)
}
