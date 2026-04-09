package test

import (
	"bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/macawls/ogre"
)

var (
	updateRef = flag.Bool("update-ref", false, "regenerate satori reference images using bun")
	showDiff  = flag.Bool("show-diff", false, "write diff images to test/output/")
)

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

func TestCompareWithSatori(t *testing.T) {
	if *updateRef {
		t.Log("Regenerating satori reference images...")
		// This would shell out to bun — skip in CI, run manually
		t.Skip("Run 'cd test/satori-reference && bun run generate.ts' manually")
	}

	fixtures, err := filepath.Glob("fixtures/*.html")
	if err != nil {
		t.Fatal(err)
	}
	if len(fixtures) == 0 {
		t.Fatal("no fixtures found in test/fixtures/")
	}

	os.MkdirAll("output", 0755)

	for _, fixture := range fixtures {
		name := strings.TrimSuffix(filepath.Base(fixture), ".html")
		t.Run(name, func(t *testing.T) {
			htmlBytes, err := os.ReadFile(fixture)
			if err != nil {
				t.Fatal(err)
			}

			result, err := ogre.Render(string(htmlBytes), ogre.Options{
				Width:  1200,
				Height: 630,
			})
			if err != nil {
				t.Fatalf("ogre.Render failed: %v", err)
			}

			os.WriteFile(filepath.Join("output", name+".svg"), result.Data, 0644)

			refSVG := filepath.Join("reference", name+".svg")
			if _, err := os.Stat(refSVG); os.IsNotExist(err) {
				t.Skipf("no reference SVG for %s (run satori-reference/generate.ts first)", name)
				return
			}

			refData, err := os.ReadFile(refSVG)
			if err != nil {
				t.Fatal(err)
			}

			report := compareSVGs(string(refData), string(result.Data), name)
			if report.score < 0.3 {
				t.Errorf("FAIL [score=%.1f%%] %s", report.score*100, report.summary)
			} else if report.score < 0.7 {
				t.Logf("WARN [score=%.1f%%] %s", report.score*100, report.summary)
			} else {
				t.Logf("PASS [score=%.1f%%]", report.score*100)
			}

			refPNG := filepath.Join("reference", name+".png")
			if _, err := os.Stat(refPNG); err == nil {
				satoriReport := comparePNGs(t, htmlBytes, refPNG, name)

				takumiPNG := filepath.Join("reference", name+".takumi.png")
				if _, err := os.Stat(takumiPNG); err == nil {
					takumiReport := comparePNGs(t, htmlBytes, takumiPNG, name)
					t.Logf("PNG vs Satori: %.1f%% | vs Takumi: %.1f%%",
						satoriReport.matchPercent*100, takumiReport.matchPercent*100)
				} else {
					t.Logf("PNG vs Satori: %.1f%% | vs Takumi: N/A (no reference)",
						satoriReport.matchPercent*100)
				}
			}
		})
	}
}

type svgReport struct {
	score   float64
	summary string
}

type svgElement struct {
	Name     string
	Attrs    map[string]string
	Text     string
	Children []svgElement
}

var trackedElements = map[string]bool{
	"rect": true, "text": true, "path": true,
	"linearGradient": true, "radialGradient": true,
	"clipPath": true, "filter": true, "g": true, "image": true,
	"stop": true, "circle": true,
}

func parseSVGTree(data string) ([]svgElement, error) {
	dec := xml.NewDecoder(strings.NewReader(data))
	var stack []*svgElement
	var root []svgElement

	for {
		tok, err := dec.Token()
		if err != nil {
			break
		}
		switch t := tok.(type) {
		case xml.StartElement:
			el := svgElement{
				Name:  t.Name.Local,
				Attrs: make(map[string]string),
			}
			for _, a := range t.Attr {
				el.Attrs[a.Name.Local] = a.Value
			}
			stack = append(stack, &el)
		case xml.CharData:
			text := strings.TrimSpace(string(t))
			if text != "" && len(stack) > 0 {
				stack[len(stack)-1].Text += text
			}
		case xml.EndElement:
			if len(stack) == 0 {
				continue
			}
			el := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if len(stack) > 0 {
				stack[len(stack)-1].Children = append(stack[len(stack)-1].Children, *el)
			} else {
				root = append(root, *el)
			}
		}
	}
	return root, nil
}

func collectElements(elements []svgElement) map[string][]svgElement {
	result := make(map[string][]svgElement)
	var walk func([]svgElement)
	walk = func(els []svgElement) {
		for _, el := range els {
			if trackedElements[el.Name] {
				result[el.Name] = append(result[el.Name], el)
			}
			walk(el.Children)
		}
	}
	walk(elements)
	return result
}

func parseFloat(s string) float64 {
	s = strings.TrimSuffix(s, "%")
	s = strings.TrimSuffix(s, "px")
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

func attrDiff(ref, act map[string]string, keys []string, tolerance float64) (matched, total int) {
	for _, k := range keys {
		rv, rok := ref[k]
		av, aok := act[k]
		if !rok {
			continue
		}
		total++
		if !aok {
			continue
		}
		rfv := parseFloat(rv)
		afv := parseFloat(av)
		if rfv == 0 && afv == 0 {
			matched++
		} else if math.Abs(rfv-afv) <= tolerance {
			matched++
		} else if rv == av {
			matched++
		}
	}
	return
}

func colorClose(a, b string) bool {
	if a == b {
		return true
	}
	a = strings.TrimSpace(strings.ToLower(a))
	b = strings.TrimSpace(strings.ToLower(b))
	if a == b {
		return true
	}
	parseHex := func(s string) (r, g, bl int, ok bool) {
		s = strings.TrimPrefix(s, "#")
		if len(s) == 6 {
			rr, _ := strconv.ParseInt(s[0:2], 16, 32)
			gg, _ := strconv.ParseInt(s[2:4], 16, 32)
			bb, _ := strconv.ParseInt(s[4:6], 16, 32)
			return int(rr), int(gg), int(bb), true
		}
		return 0, 0, 0, false
	}
	parseRGBA := func(s string) (r, g, bl int, ok bool) {
		s = strings.TrimPrefix(s, "rgba(")
		s = strings.TrimPrefix(s, "rgb(")
		s = strings.TrimSuffix(s, ")")
		parts := strings.Split(s, ",")
		if len(parts) < 3 {
			return 0, 0, 0, false
		}
		rr, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
		gg, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
		bb, _ := strconv.Atoi(strings.TrimSpace(parts[2]))
		return rr, gg, bb, true
	}

	r1, g1, b1, ok1 := parseHex(a)
	if !ok1 {
		r1, g1, b1, ok1 = parseRGBA(a)
	}
	r2, g2, b2, ok2 := parseHex(b)
	if !ok2 {
		r2, g2, b2, ok2 = parseRGBA(b)
	}
	if !ok1 || !ok2 {
		return false
	}
	dist := math.Abs(float64(r1-r2)) + math.Abs(float64(g1-g2)) + math.Abs(float64(b1-b2))
	return dist < 10
}

func isVisualRect(el svgElement) bool {
	fill := el.Attrs["fill"]
	if fill == "" || fill == "none" || fill == "transparent" {
		return false
	}
	if strings.HasPrefix(fill, "url(") {
		return true
	}
	if fill == "#fff" || fill == "#ffffff" || fill == "white" {
		return false
	}
	return true
}

func isRoundedRectPath(el svgElement) bool {
	d := el.Attrs["d"]
	return strings.Contains(d, " a") || strings.Contains(d, " A")
}

func filterVisualRects(rects []svgElement) []svgElement {
	var out []svgElement
	for _, r := range rects {
		if isVisualRect(r) {
			out = append(out, r)
		}
	}
	return out
}

func collectPathFills(paths []svgElement) map[string]int {
	fills := make(map[string]int)
	for _, p := range paths {
		f := strings.ToLower(strings.TrimSpace(p.Attrs["fill"]))
		if f != "" && f != "none" {
			fills[f]++
		}
	}
	return fills
}

func compareSVGs(reference, actual, name string) svgReport {
	refTree, refErr := parseSVGTree(reference)
	actTree, actErr := parseSVGTree(actual)
	if refErr != nil || actErr != nil || len(refTree) == 0 || len(actTree) == 0 {
		return compareSVGsFallback(reference, actual, name)
	}

	refEls := collectElements(refTree)
	actEls := collectElements(actTree)

	var issues []string

	refPaths := refEls["path"]
	actPaths := actEls["path"]
	actRects := actEls["rect"]

	refRoundedRectPaths := 0
	refTextPaths := 0
	for _, p := range refPaths {
		if isRoundedRectPath(p) {
			refRoundedRectPaths++
		} else {
			refTextPaths++
		}
	}

	actRoundedRects := 0
	for _, r := range actRects {
		if r.Attrs["rx"] != "" || r.Attrs["ry"] != "" {
			actRoundedRects++
		}
	}

	pathsEquivalent := false
	if refRoundedRectPaths > 0 && actRoundedRects > 0 {
		ratio := float64(actRoundedRects) / float64(refRoundedRectPaths)
		if ratio >= 0.5 && ratio <= 2.0 {
			pathsEquivalent = true
		}
	}
	if len(refPaths) > 0 && len(actPaths) > 0 {
		pathsEquivalent = true
	}

	countScore := 0.0
	countTotal := 0
	for _, tag := range []string{"rect", "text", "linearGradient", "radialGradient", "filter", "image"} {
		rc := len(refEls[tag])
		ac := len(actEls[tag])
		if tag == "rect" {
			rc = len(filterVisualRects(refEls["rect"]))
			ac = len(filterVisualRects(actEls["rect"]))
		}
		if rc == 0 && ac == 0 {
			continue
		}
		countTotal++
		if rc == 0 && ac > 0 {
			countScore += 0.8
			continue
		}
		if rc > 0 && ac == 0 {
			issues = append(issues, fmt.Sprintf("missing <%s> (ref=%d, actual=0)", tag, rc))
			continue
		}
		ratio := float64(ac) / float64(rc)
		if ratio > 1 {
			ratio = 1.0 / ratio
		}
		countScore += ratio
	}

	if len(refPaths) > 0 || len(actPaths) > 0 {
		countTotal++
		if pathsEquivalent || len(actPaths) > 0 {
			refFills := collectPathFills(refPaths)
			actFills := collectPathFills(actPaths)
			matched := 0
			total := 0
			for fill := range refFills {
				total++
				if actFills[fill] > 0 {
					matched++
				} else {
					for af := range actFills {
						if colorClose(fill, af) {
							matched++
							break
						}
					}
				}
			}
			if total > 0 {
				countScore += float64(matched) / float64(total)
			} else {
				countScore += 1.0
			}
		} else if len(refPaths) > 0 && len(actPaths) == 0 {
			if refRoundedRectPaths > 0 && actRoundedRects > 0 {
				countScore += 0.8
			} else {
				issues = append(issues, fmt.Sprintf("missing <path> (ref=%d, actual=0)", len(refPaths)))
			}
		}
	}

	if countTotal > 0 {
		countScore /= float64(countTotal)
	} else {
		countScore = 1.0
	}

	attrMatched := 0
	attrTotal := 0

	refVisualRects := filterVisualRects(refEls["rect"])
	actVisualRects := filterVisualRects(actEls["rect"])
	n := len(refVisualRects)
	if n > len(actVisualRects) {
		n = len(actVisualRects)
	}
	for i := 0; i < n; i++ {
		m, t := attrDiff(refVisualRects[i].Attrs, actVisualRects[i].Attrs,
			[]string{"x", "y", "width", "height"}, 3)
		attrMatched += m
		attrTotal += t
		rf := refVisualRects[i].Attrs["fill"]
		af := actVisualRects[i].Attrs["fill"]
		if rf != "" {
			attrTotal++
			if colorClose(rf, af) || rf == af {
				attrMatched++
			} else if strings.HasPrefix(rf, "url(") && strings.HasPrefix(af, "url(") {
				attrMatched++
			}
		}
	}

	refTexts := refEls["text"]
	actTexts := actEls["text"]
	n = len(refTexts)
	if n > len(actTexts) {
		n = len(actTexts)
	}
	for i := 0; i < n; i++ {
		m, t := attrDiff(refTexts[i].Attrs, actTexts[i].Attrs,
			[]string{"x", "y"}, 5)
		attrMatched += m
		attrTotal += t
		rfs := refTexts[i].Attrs["font-size"]
		afs := actTexts[i].Attrs["font-size"]
		if rfs != "" {
			attrTotal++
			if math.Abs(parseFloat(rfs)-parseFloat(afs)) < 2 {
				attrMatched++
			}
		}
	}

	refStops := refEls["stop"]
	actStops := actEls["stop"]
	n = len(refStops)
	if n > len(actStops) {
		n = len(actStops)
	}
	for i := 0; i < n; i++ {
		rOff := refStops[i].Attrs["offset"]
		aOff := actStops[i].Attrs["offset"]
		if rOff != "" {
			attrTotal++
			if math.Abs(parseFloat(rOff)-parseFloat(aOff)) < 2 {
				attrMatched++
			}
		}
		rCol := refStops[i].Attrs["stop-color"]
		aCol := actStops[i].Attrs["stop-color"]
		if rCol != "" {
			attrTotal++
			if colorClose(rCol, aCol) {
				attrMatched++
			}
		}
	}

	attrScore := 1.0
	if attrTotal > 0 {
		attrScore = float64(attrMatched) / float64(attrTotal)
	}

	textMatched := 0
	textTotal := 0
	n = len(refTexts)
	if n > len(actTexts) {
		n = len(actTexts)
	}
	for i := 0; i < n; i++ {
		textTotal++
		rt := strings.TrimSpace(refTexts[i].Text)
		at := strings.TrimSpace(actTexts[i].Text)
		if rt == at {
			textMatched++
		} else if strings.EqualFold(rt, at) {
			textMatched++
		} else if rt != "" && strings.Contains(at, rt) {
			textMatched++
		}
	}
	if len(refTexts) > len(actTexts) {
		textTotal += len(refTexts) - len(actTexts)
		issues = append(issues, fmt.Sprintf("text count mismatch: ref=%d actual=%d", len(refTexts), len(actTexts)))
	}
	textScore := 1.0
	if textTotal > 0 {
		textScore = float64(textMatched) / float64(textTotal)
	}

	veScore := visualEquivalenceScore(refEls, actEls, reference, actual)

	score := 0.2*countScore + 0.25*attrScore + 0.25*textScore + 0.3*veScore
	if score < 0 {
		score = 0
	}
	if score > 1 {
		score = 1
	}

	if attrScore < 1 && attrTotal > 0 {
		issues = append(issues, fmt.Sprintf("attribute accuracy: %d/%d", attrMatched, attrTotal))
	}

	summary := "OK"
	if len(issues) > 0 {
		summary = strings.Join(issues, "; ")
	}

	return svgReport{score: score, summary: summary}
}

func visualEquivalenceScore(refEls, actEls map[string][]svgElement, refRaw, actRaw string) float64 {
	score := 0.0
	checks := 0

	checks++
	refBg := ""
	actBg := ""
	for _, r := range refEls["rect"] {
		fill := r.Attrs["fill"]
		if r.Attrs["width"] == "1200" && r.Attrs["height"] == "630" && fill != "" && fill != "#fff" && fill != "white" {
			refBg = fill
			break
		}
	}
	for _, r := range actEls["rect"] {
		fill := r.Attrs["fill"]
		if r.Attrs["width"] == "1200" && r.Attrs["height"] == "630" && fill != "" {
			actBg = fill
			break
		}
	}
	if refBg != "" && actBg != "" {
		if refBg == actBg || colorClose(refBg, actBg) {
			score += 1.0
		} else if strings.HasPrefix(refBg, "url(") && strings.HasPrefix(actBg, "url(") {
			score += 1.0
		} else {
			score += 0.5
		}
	} else {
		score += 0.8
	}

	checks++
	refPathFills := collectPathFills(refEls["path"])
	actPathFills := collectPathFills(actEls["path"])
	fillMatched := 0
	fillTotal := 0
	for fill := range refPathFills {
		fillTotal++
		if actPathFills[fill] > 0 {
			fillMatched++
		} else {
			for af := range actPathFills {
				if colorClose(fill, af) {
					fillMatched++
					break
				}
			}
		}
	}
	if fillTotal > 0 {
		score += float64(fillMatched) / float64(fillTotal)
	} else {
		score += 1.0
	}

	checks++
	refGrads := len(refEls["linearGradient"]) + len(refEls["radialGradient"])
	actGrads := len(actEls["linearGradient"]) + len(actEls["radialGradient"])
	if refGrads == 0 && actGrads == 0 {
		score += 1.0
	} else if refGrads > 0 && actGrads > 0 {
		ratio := float64(actGrads) / float64(refGrads)
		if ratio > 1 {
			ratio = 1.0 / ratio
		}
		score += ratio
	} else {
		score += 0.0
	}

	checks++
	refVisual := len(filterVisualRects(refEls["rect"])) + len(refEls["path"])
	actVisual := len(filterVisualRects(actEls["rect"])) + len(actEls["path"])
	if refVisual == 0 && actVisual == 0 {
		score += 1.0
	} else if refVisual > 0 && actVisual > 0 {
		ratio := float64(actVisual) / float64(refVisual)
		if ratio > 1 {
			ratio = 1.0 / ratio
		}
		if ratio > 0.5 {
			score += ratio
		} else {
			score += ratio * 0.5
		}
	}

	return score / float64(checks)
}

func compareSVGsFallback(reference, actual, name string) svgReport {
	var issues []string
	score := 1.0

	refHasRect := strings.Contains(reference, "<rect")
	actHasRect := strings.Contains(actual, "<rect")
	if refHasRect && !actHasRect {
		issues = append(issues, "missing <rect> elements")
		score -= 0.2
	}

	refHasText := strings.Contains(reference, "<text")
	actHasText := strings.Contains(actual, "<text")
	if refHasText && !actHasText {
		issues = append(issues, "missing <text> elements")
		score -= 0.3
	}

	refHasGradient := strings.Contains(reference, "Gradient")
	actHasGradient := strings.Contains(actual, "Gradient")
	if refHasGradient && !actHasGradient {
		issues = append(issues, "missing gradient definitions")
		score -= 0.2
	}

	refHasClip := strings.Contains(reference, "clipPath")
	actHasClip := strings.Contains(actual, "clipPath")
	if refHasClip && !actHasClip {
		issues = append(issues, "missing clip-path")
		score -= 0.1
	}

	refHasFilter := strings.Contains(reference, "<filter")
	actHasFilter := strings.Contains(actual, "<filter")
	if refHasFilter && !actHasFilter {
		issues = append(issues, "missing SVG filters (shadow)")
		score -= 0.1
	}

	refHasOpacity := strings.Contains(reference, `opacity="`)
	actHasOpacity := strings.Contains(actual, `opacity="`)
	if refHasOpacity && !actHasOpacity {
		issues = append(issues, "missing opacity")
		score -= 0.1
	}

	refRectCount := strings.Count(reference, "<rect")
	actRectCount := strings.Count(actual, "<rect")
	if refRectCount > 0 && actRectCount > 0 {
		ratio := float64(actRectCount) / float64(refRectCount)
		if ratio < 0.5 || ratio > 2.0 {
			issues = append(issues, fmt.Sprintf("rect count mismatch: ref=%d actual=%d", refRectCount, actRectCount))
			score -= 0.1
		}
	}

	refTextCount := strings.Count(reference, "<text")
	actTextCount := strings.Count(actual, "<text")
	if refTextCount > 0 && actTextCount > 0 {
		ratio := float64(actTextCount) / float64(refTextCount)
		if ratio < 0.5 || ratio > 2.0 {
			issues = append(issues, fmt.Sprintf("text count mismatch: ref=%d actual=%d", refTextCount, actTextCount))
			score -= 0.1
		}
	}

	if score < 0 {
		score = 0
	}

	summary := "OK"
	if len(issues) > 0 {
		summary = strings.Join(issues, "; ")
	}

	return svgReport{score: score, summary: summary}
}

type pngReport struct {
	matchPercent   float64
	matchingPixels int
	totalPixels    int
}

func comparePNGs(t *testing.T, htmlBytes []byte, refPath, name string) pngReport {
	t.Helper()

	pngResult, err := ogre.Render(string(htmlBytes), ogre.Options{
		Width: 1200, Height: 630, Format: ogre.FormatPNG,
	})
	if err != nil {
		t.Logf("PNG render failed: %v", err)
		return pngReport{}
	}

	os.WriteFile(filepath.Join("output", name+".png"), pngResult.Data, 0644)

	actualImg, err := png.Decode(bytes.NewReader(pngResult.Data))
	if err != nil {
		t.Logf("failed to decode ogre PNG: %v", err)
		return pngReport{}
	}

	refFile, err := os.Open(refPath)
	if err != nil {
		t.Logf("failed to open reference PNG: %v", err)
		return pngReport{}
	}
	defer refFile.Close()

	refImg, err := png.Decode(refFile)
	if err != nil {
		t.Logf("failed to decode reference PNG: %v", err)
		return pngReport{}
	}

	bounds := refImg.Bounds()
	total := bounds.Dx() * bounds.Dy()
	matching := 0
	threshold := 50.0

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			d := pixelDist(refImg.At(x, y), actualImg.At(x, y))
			if d < threshold {
				matching++
			}
		}
	}

	if *showDiff {
		diffPath := filepath.Join("output", name+"_diff.png")
		if err := writeDiffImage(refImg, actualImg, diffPath); err != nil {
			t.Logf("failed to write diff image: %v", err)
		}
	}

	return pngReport{
		matchPercent:   float64(matching) / float64(total),
		matchingPixels: matching,
		totalPixels:    total,
	}
}

func pixelDist(c1, c2 color.Color) float64 {
	r1, g1, b1, _ := c1.RGBA()
	r2, g2, b2, _ := c2.RGBA()
	dr := float64(r1>>8) - float64(r2>>8)
	dg := float64(g1>>8) - float64(g2>>8)
	db := float64(b1>>8) - float64(b2>>8)
	return math.Sqrt(dr*dr + dg*dg + db*db)
}

func writeDiffImage(ref, actual image.Image, path string) error {
	bounds := ref.Bounds()
	diff := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			d := pixelDist(ref.At(x, y), actual.At(x, y))
			if d < 10 {
				diff.Set(x, y, color.RGBA{0, 0, 0, 255})
			} else {
				intensity := uint8(math.Min(d, 255))
				diff.Set(x, y, color.RGBA{intensity, 0, 0, 255})
			}
		}
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	var buf bytes.Buffer
	png.Encode(&buf, diff)
	_, err = f.Write(buf.Bytes())
	return err
}
