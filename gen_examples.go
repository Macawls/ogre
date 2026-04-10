//go:build ignore

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/macawls/ogre"
	"github.com/macawls/ogre/font"
)

func readTemplate(content, name string) string {
	key := name + ": `"
	start := strings.Index(content, key)
	if start < 0 {
		panic("template not found: " + name)
	}
	start += len(key)
	end := strings.Index(content[start:], "`,")
	if end < 0 {
		end = strings.Index(content[start:], "`}")
	}
	if end < 0 {
		panic("template end not found: " + name)
	}
	return content[start : start+end]
}

func main() {
	data, err := os.ReadFile("site/src/data/templates.ts")
	if err != nil {
		panic(err)
	}
	content := string(data)

	names := []string{"blog", "event", "product", "repo", "emoji", "rtl", "tailwind"}

	renderer := ogre.NewRenderer()
	cache := font.NewFontCache(os.TempDir())

	arabicData, err := font.FetchGoogleFont("Noto Sans Arabic", 400, cache)
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: could not fetch Arabic font: %v\n", err)
	} else {
		renderer.LoadFont(ogre.FontSource{Name: "Noto Sans Arabic", Weight: 400, Data: arabicData})
		bold, _ := font.FetchGoogleFont("Noto Sans Arabic", 700, cache)
		if bold != nil {
			renderer.LoadFont(ogre.FontSource{Name: "Noto Sans Arabic", Weight: 700, Data: bold})
		}
	}

	outDir := "site/public/examples"
	for _, name := range names {
		html := readTemplate(content, name)
		for _, format := range []ogre.Format{ogre.FormatSVG, ogre.FormatPNG, ogre.FormatJPEG} {
			quality := 0
			if format == ogre.FormatJPEG {
				quality = 95
			}
			result, err := renderer.Render(html, ogre.Options{Width: 1200, Height: 630, Format: format, Quality: quality})
			if err != nil {
				fmt.Fprintf(os.Stderr, "error rendering %s/%s: %v\n", name, format, err)
				continue
			}
			ext := strings.ToLower(string(format))
			if ext == "jpeg" {
				ext = "jpg"
			}
			path := fmt.Sprintf("%s/ex-%s.%s", outDir, name, ext)
			if err := os.WriteFile(path, result.Data, 0644); err != nil {
				fmt.Fprintf(os.Stderr, "error writing %s: %v\n", path, err)
				continue
			}
			fmt.Printf("wrote %s\n", path)
		}
	}
}
