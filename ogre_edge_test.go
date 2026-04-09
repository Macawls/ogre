package ogre

import (
	"bytes"
	"encoding/binary"
	"encoding/xml"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestRenderEmptyHTML(t *testing.T) {
	result, err := Render("", Options{Width: 400, Height: 300})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || len(result.Data) == 0 {
		t.Fatal("expected non-empty result")
	}
	if !strings.Contains(string(result.Data), "<svg") {
		t.Error("output does not contain <svg")
	}
}

func TestRenderWhitespaceOnly(t *testing.T) {
	result, err := Render("   ", Options{Width: 400, Height: 300})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || len(result.Data) == 0 {
		t.Fatal("expected non-empty result")
	}
	if !strings.Contains(string(result.Data), "<svg") {
		t.Error("output does not contain <svg")
	}
}

func TestRenderSingleDiv(t *testing.T) {
	result, err := Render("<div></div>", Options{Width: 400, Height: 300})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || len(result.Data) == 0 {
		t.Fatal("expected non-empty result")
	}
	if !strings.Contains(string(result.Data), "<svg") {
		t.Error("output does not contain <svg")
	}
}

func TestRenderTextOnly(t *testing.T) {
	result, err := Render("Hello", Options{Width: 400, Height: 300})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || len(result.Data) == 0 {
		t.Fatal("expected non-empty result")
	}
	svg := string(result.Data)
	if !strings.Contains(svg, "<svg") {
		t.Error("output does not contain <svg")
	}
}

func TestRenderLongText(t *testing.T) {
	longText := strings.Repeat("A", 10000)
	html := fmt.Sprintf(`<div style="width:800px">%s</div>`, longText)

	start := time.Now()
	result, err := Render(html, Options{Width: 1200, Height: 630})
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || len(result.Data) == 0 {
		t.Fatal("expected non-empty result")
	}
	if elapsed > 5*time.Second {
		t.Errorf("render took %v, want < 5s", elapsed)
	}
}

func TestRenderManyWords(t *testing.T) {
	words := make([]string, 500)
	for i := range words {
		words[i] = fmt.Sprintf("word%d", i)
	}
	html := fmt.Sprintf(`<div style="width:600px">%s</div>`, strings.Join(words, " "))

	result, err := Render(html, Options{Width: 1200, Height: 630})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || len(result.Data) == 0 {
		t.Fatal("expected non-empty result")
	}
}

func TestRenderDeepNesting(t *testing.T) {
	var b strings.Builder
	for i := 0; i < 50; i++ {
		b.WriteString("<div>")
	}
	b.WriteString("deep")
	for i := 0; i < 50; i++ {
		b.WriteString("</div>")
	}

	result, err := Render(b.String(), Options{Width: 400, Height: 300})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || len(result.Data) == 0 {
		t.Fatal("expected non-empty result")
	}
}

func TestRenderZeroWidth(t *testing.T) {
	result, err := Render(`<div style="background:red">test</div>`, Options{Width: 0, Height: 300})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Width != 1200 {
		t.Errorf("width = %d, want 1200 (default)", result.Width)
	}
}

func TestRenderZeroHeight(t *testing.T) {
	result, err := Render(`<div style="background:red">test</div>`, Options{Width: 400, Height: 0})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Height != 630 {
		t.Errorf("height = %d, want 630 (default)", result.Height)
	}
}

func TestRenderLargeSize(t *testing.T) {
	result, err := Render(`<div style="background:blue;width:100%;height:100%">big</div>`, Options{Width: 4000, Height: 4000})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || len(result.Data) == 0 {
		t.Fatal("expected non-empty result")
	}
	if result.Width != 4000 || result.Height != 4000 {
		t.Errorf("dimensions = %dx%d, want 4000x4000", result.Width, result.Height)
	}
}

func TestRenderSmallSize(t *testing.T) {
	result, err := Render(`<div style="background:green">tiny</div>`, Options{Width: 10, Height: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || len(result.Data) == 0 {
		t.Fatal("expected non-empty result")
	}
	if result.Width != 10 || result.Height != 10 {
		t.Errorf("dimensions = %dx%d, want 10x10", result.Width, result.Height)
	}
}

func TestRenderUnclosedTags(t *testing.T) {
	result, err := Render("<div><p>text", Options{Width: 400, Height: 300})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || len(result.Data) == 0 {
		t.Fatal("expected non-empty result")
	}
}

func TestRenderBadNesting(t *testing.T) {
	result, err := Render("<div><p></div></p>", Options{Width: 400, Height: 300})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || len(result.Data) == 0 {
		t.Fatal("expected non-empty result")
	}
}

func TestRenderInvalidStyle(t *testing.T) {
	result, err := Render(`<div style="color: notacolor; font-size: abc">text</div>`, Options{Width: 400, Height: 300})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || len(result.Data) == 0 {
		t.Fatal("expected non-empty result")
	}
}

func TestRenderHTMLEntities(t *testing.T) {
	result, err := Render(`<div>&amp; &lt; &gt;</div>`, Options{Width: 400, Height: 300})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || len(result.Data) == 0 {
		t.Fatal("expected non-empty result")
	}
	svg := string(result.Data)
	if strings.Contains(svg, "&&") {
		t.Error("SVG contains unescaped ampersand")
	}
}

func TestRenderUnicode(t *testing.T) {
	result, err := Render(`<div>你好世界 こんにちは 안녕하세요</div>`, Options{Width: 800, Height: 400})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || len(result.Data) == 0 {
		t.Fatal("expected non-empty result")
	}
}

func TestRenderSpecialChars(t *testing.T) {
	result, err := Render(`<div>"quotes" 'single' \backslash\ &amp; ampersand</div>`, Options{Width: 400, Height: 300})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || len(result.Data) == 0 {
		t.Fatal("expected non-empty result")
	}
}

func TestRenderNegativeMargin(t *testing.T) {
	result, err := Render(`<div style="margin:-10px;background:red">negative margin</div>`, Options{Width: 400, Height: 300})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || len(result.Data) == 0 {
		t.Fatal("expected non-empty result")
	}
}

func TestRenderPercentageOnly(t *testing.T) {
	result, err := Render(`<div style="width:100%;height:100%;background:blue"><div style="width:50%;height:50%;background:red"></div></div>`, Options{Width: 400, Height: 300})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || len(result.Data) == 0 {
		t.Fatal("expected non-empty result")
	}
}

func TestRenderAutoMargins(t *testing.T) {
	result, err := Render(`<div style="width:100%;height:100%"><div style="width:100px;height:100px;margin:auto;background:green"></div></div>`, Options{Width: 400, Height: 300})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || len(result.Data) == 0 {
		t.Fatal("expected non-empty result")
	}
}

func TestRenderFlexGrowZero(t *testing.T) {
	result, err := Render(`<div style="display:flex;width:400px;height:200px"><div style="flex-grow:0;width:50px;background:red"></div><div style="flex-grow:0;width:50px;background:blue"></div></div>`, Options{Width: 400, Height: 300})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || len(result.Data) == 0 {
		t.Fatal("expected non-empty result")
	}
}

func TestRenderTailwindOnly(t *testing.T) {
	result, err := Render(`<div class="flex items-center justify-center w-full h-full bg-blue-500"><div class="text-white text-4xl font-bold p-8">Tailwind</div></div>`, Options{Width: 800, Height: 400})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || len(result.Data) == 0 {
		t.Fatal("expected non-empty result")
	}
}

func TestRenderInvalidTailwind(t *testing.T) {
	result, err := Render(`<div class="not-a-real-class another-fake-one">content</div>`, Options{Width: 400, Height: 300})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || len(result.Data) == 0 {
		t.Fatal("expected non-empty result")
	}
}

func TestRenderTailwindOverride(t *testing.T) {
	result, err := Render(`<div class="p-4" style="padding:0">override test</div>`, Options{Width: 400, Height: 300})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || len(result.Data) == 0 {
		t.Fatal("expected non-empty result")
	}
}

func TestRenderPNGValid(t *testing.T) {
	result, err := Render(`<div style="width:200px;height:100px;background:red">PNG test</div>`, Options{Width: 400, Height: 300, Format: FormatPNG})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || len(result.Data) < 8 {
		t.Fatal("expected non-empty PNG data")
	}
	pngMagic := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	if !bytes.HasPrefix(result.Data, pngMagic) {
		t.Errorf("PNG data does not start with magic bytes, got %v", result.Data[:8])
	}
	if result.ContentType != "image/png" {
		t.Errorf("content type = %q, want image/png", result.ContentType)
	}
}

func TestRenderSVGValid(t *testing.T) {
	result, err := Render(`<div style="width:200px;height:100px;background:blue">SVG test</div>`, Options{Width: 400, Height: 300})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || len(result.Data) == 0 {
		t.Fatal("expected non-empty SVG data")
	}
	if result.ContentType != "image/svg+xml" {
		t.Errorf("content type = %q, want image/svg+xml", result.ContentType)
	}
	decoder := xml.NewDecoder(bytes.NewReader(result.Data))
	for {
		_, err := decoder.Token()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			t.Fatalf("SVG is not valid XML: %v", err)
		}
	}
}

func TestRenderPNGDimensions(t *testing.T) {
	width, height := 400, 300
	result, err := Render(`<div style="width:100%;height:100%;background:green">dim test</div>`, Options{Width: width, Height: height, Format: FormatPNG})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || len(result.Data) < 24 {
		t.Fatal("PNG data too short")
	}
	if result.Width != width || result.Height != height {
		t.Errorf("result dimensions = %dx%d, want %dx%d", result.Width, result.Height, width, height)
	}
	pngWidth := binary.BigEndian.Uint32(result.Data[16:20])
	pngHeight := binary.BigEndian.Uint32(result.Data[20:24])
	if int(pngWidth) != width {
		t.Errorf("PNG IHDR width = %d, want %d", pngWidth, width)
	}
	if int(pngHeight) != height {
		t.Errorf("PNG IHDR height = %d, want %d", pngHeight, height)
	}
}

func TestRenderConcurrent(t *testing.T) {
	var wg sync.WaitGroup
	colors := []string{"red", "blue", "green"}
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			html := fmt.Sprintf(`<div style="background:%s;width:100%%;height:100%%">Test %d</div>`,
				colors[n%3], n)
			result, err := Render(html, Options{Width: 400, Height: 200})
			if err != nil {
				t.Errorf("concurrent render %d failed: %v", n, err)
				return
			}
			if len(result.Data) == 0 {
				t.Errorf("concurrent render %d empty", n)
			}
		}(i)
	}
	wg.Wait()
}
