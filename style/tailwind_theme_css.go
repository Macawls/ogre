package style

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseTailwindThemeCSS parses a Tailwind v4 @theme block into a
// TailwindConfig. It recognises the standard v4 theme namespaces:
//
//	--color-<name>-<shade>: <value>;
//	--spacing-<key>:        <value>;
//	--font-<key>:           <value>;
//	--breakpoint-<key>:     <value>;
//
// Any @theme variant is accepted (@theme, @theme inline, @theme static,
// @theme reference, @theme default); the CSS outside those blocks is
// ignored. C-style comments are stripped before parsing.
func ParseTailwindThemeCSS(css string) (*TailwindConfig, error) {
	cfg := &TailwindConfig{
		Colors:  map[string]map[int]string{},
		Spacing: map[string]string{},
		Fonts:   map[string]string{},
		Screens: map[string]string{},
	}
	stripped := stripCSSComments(css)
	i := 0
	for i < len(stripped) {
		at := strings.Index(stripped[i:], "@theme")
		if at < 0 {
			break
		}
		start := i + at + len("@theme")
		braceOpen := strings.IndexByte(stripped[start:], '{')
		if braceOpen < 0 {
			return nil, fmt.Errorf("@theme block missing '{'")
		}
		openAt := start + braceOpen
		closeAt, err := matchBrace(stripped, openAt)
		if err != nil {
			return nil, err
		}
		body := stripped[openAt+1 : closeAt]
		if err := parseThemeBody(body, cfg); err != nil {
			return nil, err
		}
		i = closeAt + 1
	}
	if isEmptyConfig(cfg) {
		return nil, fmt.Errorf("no @theme block found or block was empty")
	}
	return cfg, nil
}

func stripCSSComments(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for i := 0; i < len(s); {
		if i+1 < len(s) && s[i] == '/' && s[i+1] == '*' {
			end := strings.Index(s[i+2:], "*/")
			if end < 0 {
				break
			}
			i += 2 + end + 2
			continue
		}
		b.WriteByte(s[i])
		i++
	}
	return b.String()
}

func matchBrace(s string, openIdx int) (int, error) {
	depth := 1
	for i := openIdx + 1; i < len(s); i++ {
		switch s[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return i, nil
			}
		}
	}
	return 0, fmt.Errorf("unbalanced '{' in @theme block")
}

func parseThemeBody(body string, cfg *TailwindConfig) error {
	for _, decl := range strings.Split(body, ";") {
		decl = strings.TrimSpace(decl)
		if decl == "" {
			continue
		}
		colon := strings.IndexByte(decl, ':')
		if colon < 0 {
			continue
		}
		key := strings.TrimSpace(decl[:colon])
		val := strings.TrimSpace(decl[colon+1:])
		if !strings.HasPrefix(key, "--") {
			continue
		}
		assignThemeVar(key[2:], val, cfg)
	}
	return nil
}

func assignThemeVar(key, val string, cfg *TailwindConfig) {
	switch {
	case strings.HasPrefix(key, "color-"):
		body := key[len("color-"):]
		name, shade, ok := splitLastShade(body)
		if !ok {
			return
		}
		if cfg.Colors[name] == nil {
			cfg.Colors[name] = map[int]string{}
		}
		cfg.Colors[name][shade] = val
	case strings.HasPrefix(key, "spacing-"):
		cfg.Spacing[key[len("spacing-"):]] = val
	case strings.HasPrefix(key, "font-"):
		cfg.Fonts[key[len("font-"):]] = val
	case strings.HasPrefix(key, "breakpoint-"):
		cfg.Screens[key[len("breakpoint-"):]] = val
	}
}

func splitLastShade(s string) (string, int, bool) {
	lastDash := strings.LastIndexByte(s, '-')
	if lastDash < 0 {
		return "", 0, false
	}
	shade, err := strconv.Atoi(s[lastDash+1:])
	if err != nil {
		return "", 0, false
	}
	return s[:lastDash], shade, true
}

func isEmptyConfig(c *TailwindConfig) bool {
	return len(c.Colors) == 0 && len(c.Spacing) == 0 && len(c.Fonts) == 0 && len(c.Screens) == 0 && len(c.Extend) == 0
}
