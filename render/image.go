package render

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/macawls/ogre/parse"
	"github.com/macawls/ogre/style"
)

var imageCache sync.Map

// RenderImage generates the corresponding output format.
// RenderImage generates SVG image element from src URL.
func RenderImage(src string, cs *style.ComputedStyle, x, y, w, h float64) string {
	dataURI, ok := resolveImageSource(src)
	if !ok {
		return renderBrokenImage(x, y, w, h)
	}

	par := objectFitToPreserveAspectRatio(cs.ObjectFit)

	ox, oy := parseObjectPosition(cs.ObjectPosition, w, h)
	imgX := x + ox
	imgY := y + oy

	return fmt.Sprintf(`<image href="%s" x="%.4g" y="%.4g" width="%.4g" height="%.4g" preserveAspectRatio="%s"/>`,
		xmlEscape(dataURI), imgX, imgY, w, h, par)
}

type cachedImage struct {
	dataURI string
}

func fetchImage(url string) (string, bool) {
	if v, ok := imageCache.Load(url); ok {
		return v.(*cachedImage).dataURI, true
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", false
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", false
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", false
	}
	ct := resp.Header.Get("Content-Type")
	if ct == "" {
		ct = "image/png"
	}
	if idx := strings.Index(ct, ";"); idx != -1 {
		ct = strings.TrimSpace(ct[:idx])
	}
	dataURI := fmt.Sprintf("data:%s;base64,%s", ct, base64.StdEncoding.EncodeToString(body))
	imageCache.Store(url, &cachedImage{dataURI: dataURI})
	return dataURI, true
}

func resolveImageSource(src string) (string, bool) {
	if strings.HasPrefix(src, "data:") {
		return src, true
	}
	if strings.HasPrefix(src, "http://") || strings.HasPrefix(src, "https://") {
		return fetchImage(src)
	}
	return "", false
}

func objectFitToPreserveAspectRatio(fit style.ObjectFit) string {
	switch fit {
	case style.ObjectFitContain:
		return "xMidYMid meet"
	case style.ObjectFitCover:
		return "xMidYMid slice"
	case style.ObjectFitFill:
		return "none"
	case style.ObjectFitScaleDown:
		return "xMidYMid meet"
	case style.ObjectFitNone:
		return "xMidYMid meet"
	default:
		return "none"
	}
}

func parseObjectPosition(pos string, w, h float64) (float64, float64) {
	pos = strings.TrimSpace(pos)
	if pos == "" {
		return 0, 0
	}
	parts := strings.Fields(pos)
	ox := parsePosValue(parts[0], w)
	oy := 0.0
	if len(parts) > 1 {
		oy = parsePosValue(parts[1], h)
	}
	return ox, oy
}

func parsePosValue(s string, size float64) float64 {
	s = strings.TrimSpace(s)
	switch s {
	case "center":
		return 0
	case "left", "top":
		return -size / 2
	case "right", "bottom":
		return size / 2
	}
	if strings.HasSuffix(s, "%") {
		var pct float64
		fmt.Sscanf(s[:len(s)-1], "%f", &pct)
		return (pct/100 - 0.5) * size
	}
	if strings.HasSuffix(s, "px") {
		var v float64
		fmt.Sscanf(s[:len(s)-2], "%f", &v)
		return v
	}
	return 0
}

func SerializeSVGNode(n *parse.Node) string {
	var b strings.Builder
	serializeNode(&b, n)
	return b.String()
}

func serializeNode(b *strings.Builder, n *parse.Node) {
	if n.Type == parse.TextNode {
		b.WriteString(xmlEscape(n.Text))
		return
	}
	b.WriteByte('<')
	b.WriteString(n.Tag)
	for k, v := range n.Attrs {
		fmt.Fprintf(b, ` %s="%s"`, k, xmlEscape(v))
	}
	for k, v := range n.Style {
		fmt.Fprintf(b, ` %s="%s"`, k, xmlEscape(v))
	}
	if len(n.Children) == 0 {
		b.WriteString("/>")
		return
	}
	b.WriteByte('>')
	for _, c := range n.Children {
		serializeNode(b, c)
	}
	fmt.Fprintf(b, "</%s>", n.Tag)
}

func RenderInlineSVG(pn *parse.Node, x, y, w, h float64) string {
	var b strings.Builder
	b.WriteString(`<svg xmlns="http://www.w3.org/2000/svg"`)
	if vb, ok := pn.Attrs["viewBox"]; ok {
		fmt.Fprintf(&b, ` viewBox="%s"`, xmlEscape(vb))
	}
	fmt.Fprintf(&b, ` x="%.4g" y="%.4g" width="%.4g" height="%.4g"`, x, y, w, h)
	for k, v := range pn.Attrs {
		switch k {
		case "viewBox", "width", "height", "xmlns":
			continue
		}
		fmt.Fprintf(&b, ` %s="%s"`, k, xmlEscape(v))
	}
	b.WriteByte('>')
	for _, c := range pn.Children {
		serializeNode(&b, c)
	}
	b.WriteString("</svg>")
	return b.String()
}

func renderBrokenImage(x, y, w, h float64) string {
	var b strings.Builder
	fmt.Fprintf(&b, `<rect x="%.4g" y="%.4g" width="%.4g" height="%.4g" fill="#f0f0f0" stroke="#cccccc" stroke-width="1"/>`, x, y, w, h)
	cx := x + w/2
	cy := y + h/2
	fontSize := 12.0
	if w < 60 || h < 30 {
		fontSize = 8
	}
	fmt.Fprintf(&b, `<text x="%.4g" y="%.4g" font-size="%.4g" fill="#999999" text-anchor="middle" dominant-baseline="central">&#x1F5BC;</text>`,
		cx, cy, fontSize)
	return b.String()
}
