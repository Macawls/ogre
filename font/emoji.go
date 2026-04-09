package font

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

type TextSegment struct {
	Text    string
	IsEmoji bool
}

func IsEmoji(r rune) bool {
	switch {
	case r >= 0x1F600 && r <= 0x1F64F:
		return true
	case r >= 0x1F300 && r <= 0x1F5FF:
		return true
	case r >= 0x1F680 && r <= 0x1F6FF:
		return true
	case r >= 0x1F1E0 && r <= 0x1F1FF:
		return true
	case r >= 0x2600 && r <= 0x27BF:
		return true
	case r >= 0x1F900 && r <= 0x1F9FF:
		return true
	case r >= 0x1FA00 && r <= 0x1FA6F:
		return true
	case r >= 0x1FA70 && r <= 0x1FAFF:
		return true
	case r >= 0x2700 && r <= 0x27BF:
		return true
	case r >= 0xFE00 && r <= 0xFE0F:
		return true
	case r == 0x200D:
		return true
	case r == 0x2764 || r == 0x2B50 || r == 0x2705 || r == 0x2B55 || r == 0x2934 || r == 0x2935:
		return true
	case r == 0x2139 || r == 0x2328 || r == 0x23CF:
		return true
	case r >= 0x23E9 && r <= 0x23F3:
		return true
	case r >= 0x23F8 && r <= 0x23FA:
		return true
	case r == 0x20E3:
		return true
	case r >= 0xE0020 && r <= 0xE007F:
		return true
	}
	return false
}

func SplitEmoji(text string) []TextSegment {
	if text == "" {
		return nil
	}

	var segments []TextSegment
	runes := []rune(text)
	i := 0

	for i < len(runes) {
		if IsEmoji(runes[i]) {
			start := i
			for i < len(runes) && IsEmoji(runes[i]) {
				i++
			}
			segments = append(segments, TextSegment{
				Text:    string(runes[start:i]),
				IsEmoji: true,
			})
		} else {
			start := i
			for i < len(runes) && !IsEmoji(runes[i]) {
				i++
			}
			segments = append(segments, TextSegment{
				Text:    string(runes[start:i]),
				IsEmoji: false,
			})
		}
	}

	return segments
}

func emojiCodepoints(emoji string) string {
	var parts []string
	for _, r := range emoji {
		if r == 0xFE0F {
			continue
		}
		parts = append(parts, fmt.Sprintf("%x", r))
	}
	return strings.Join(parts, "-")
}

type EmojiStyle string

const (
	EmojiTwemoji  EmojiStyle = "twemoji"
	EmojiOpenMoji EmojiStyle = "openmoji"
	EmojiNoto     EmojiStyle = "noto"
	EmojiNone     EmojiStyle = "none"
)

func TwemojiURL(emoji string) string {
	return "https://cdn.jsdelivr.net/gh/twitter/twemoji@latest/assets/svg/" + emojiCodepoints(emoji) + ".svg"
}

func TwemojiPNGURL(emoji string) string {
	return "https://cdn.jsdelivr.net/gh/twitter/twemoji@latest/assets/72x72/" + emojiCodepoints(emoji) + ".png"
}

func OpenMojiURL(emoji string) string {
	cp := strings.ToUpper(emojiCodepoints(emoji))
	return "https://cdn.jsdelivr.net/gh/hfg-gmuend/openmoji@latest/color/svg/" + cp + ".svg"
}

func OpenMojiPNGURL(emoji string) string {
	cp := strings.ToUpper(emojiCodepoints(emoji))
	return "https://cdn.jsdelivr.net/gh/hfg-gmuend/openmoji@latest/color/72x72/" + cp + ".png"
}

func NotoEmojiURL(emoji string) string {
	cp := emojiCodepoints(emoji)
	return "https://cdn.jsdelivr.net/gh/googlefonts/noto-emoji@main/svg/emoji_u" + strings.ReplaceAll(cp, "-", "_") + ".svg"
}

func EmojiSVGURL(emoji string, style EmojiStyle) string {
	switch style {
	case EmojiOpenMoji:
		return OpenMojiURL(emoji)
	case EmojiNoto:
		return NotoEmojiURL(emoji)
	default:
		return TwemojiURL(emoji)
	}
}

func EmojiPNGURL(emoji string, style EmojiStyle) string {
	switch style {
	case EmojiOpenMoji:
		return OpenMojiPNGURL(emoji)
	default:
		return TwemojiPNGURL(emoji)
	}
}

type EmojiProvider struct {
	Style    EmojiStyle
	cache    map[string][]byte
	pngCache map[string]image.Image
	mu       sync.RWMutex
}

func NewEmojiProvider() *EmojiProvider {
	return &EmojiProvider{
		Style:    EmojiTwemoji,
		cache:    make(map[string][]byte),
		pngCache: make(map[string]image.Image),
	}
}

func NewEmojiProviderWithStyle(style EmojiStyle) *EmojiProvider {
	return &EmojiProvider{
		Style:    style,
		cache:    make(map[string][]byte),
		pngCache: make(map[string]image.Image),
	}
}

func (p *EmojiProvider) FetchSVG(emoji string) ([]byte, error) {
	key := "svg:" + emoji
	p.mu.RLock()
	if data, ok := p.cache[key]; ok {
		p.mu.RUnlock()
		return data, nil
	}
	p.mu.RUnlock()

	url := EmojiSVGURL(emoji, p.Style)
	data, err := fetchURL(url)
	if err != nil {
		return nil, err
	}

	p.mu.Lock()
	p.cache[key] = data
	p.mu.Unlock()

	return data, nil
}

func (p *EmojiProvider) FetchPNG(emoji string) (image.Image, error) {
	p.mu.RLock()
	if img, ok := p.pngCache[emoji]; ok {
		p.mu.RUnlock()
		return img, nil
	}
	p.mu.RUnlock()

	url := EmojiPNGURL(emoji, p.Style)
	data, err := fetchURL(url)
	if err != nil {
		return nil, err
	}

	img, err := png.Decode(strings.NewReader(string(data)))
	if err != nil {
		return nil, fmt.Errorf("decode emoji PNG %q: %w", emoji, err)
	}

	p.mu.Lock()
	p.pngCache[emoji] = img
	p.mu.Unlock()

	return img, nil
}

func (p *EmojiProvider) SVGURL(emoji string) string {
	return EmojiSVGURL(emoji, p.Style)
}

func fetchURL(url string) ([]byte, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetch %q: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch %q: status %d", url, resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read %q: %w", url, err)
	}
	return data, nil
}
