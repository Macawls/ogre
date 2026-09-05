package style

import (
	"strconv"
	"strings"
)

// ParseTrackList parses a CSS grid-template-columns / -rows value into a
// TrackList. Supports fixed lengths, percentages, `auto`, `Nfr`, and
// `repeat(N, ...)`. Not yet supported (Phase A limitations): minmax(),
// fit-content(), min-content, max-content, repeat(auto-fill/auto-fit, ...),
// and named line brackets.
func ParseTrackList(s string) TrackList {
	tokens := splitTrackTokens(s)
	var tracks []TrackSize
	for _, tok := range tokens {
		if strings.HasPrefix(tok, "repeat(") {
			tracks = append(tracks, expandRepeat(tok)...)
			continue
		}
		if t, ok := parseTrackSize(tok); ok {
			tracks = append(tracks, t)
		}
	}
	return TrackList{Tracks: tracks}
}

func splitTrackTokens(s string) []string {
	var out []string
	depth := 0
	start := 0
	for i, r := range s {
		switch r {
		case '(':
			depth++
		case ')':
			depth--
		case ' ':
			if depth == 0 {
				if start < i {
					out = append(out, s[start:i])
				}
				start = i + 1
			}
		}
		_ = r
	}
	if start < len(s) {
		out = append(out, s[start:])
	}
	return out
}

func parseTrackSize(tok string) (TrackSize, bool) {
	tok = strings.TrimSpace(tok)
	if tok == "" {
		return TrackSize{}, false
	}
	if tok == "auto" || tok == "min-content" || tok == "max-content" {
		return TrackSize{Kind: TrackAuto}, true
	}
	if strings.HasPrefix(tok, "minmax(") && strings.HasSuffix(tok, ")") {
		inner := tok[len("minmax(") : len(tok)-1]
		parts := splitCSVOuter(inner)
		if len(parts) == 2 {
			if t, ok := parseTrackSize(strings.TrimSpace(parts[1])); ok {
				return t, true
			}
		}
		return TrackSize{Kind: TrackAuto}, true
	}
	if strings.HasPrefix(tok, "fit-content(") && strings.HasSuffix(tok, ")") {
		return TrackSize{Kind: TrackAuto}, true
	}
	if strings.HasSuffix(tok, "fr") {
		coeff, err := strconv.ParseFloat(tok[:len(tok)-2], 64)
		if err != nil {
			return TrackSize{}, false
		}
		return TrackSize{Kind: TrackFr, Fr: coeff}, true
	}
	v := ParseValue(tok)
	if v.Unit == UnitNone || v.Unit == UnitAuto {
		return TrackSize{Kind: TrackAuto}, true
	}
	return TrackSize{Kind: TrackFixed, Length: v}, true
}

func splitCSVOuter(s string) []string {
	var parts []string
	depth := 0
	start := 0
	for i, r := range s {
		switch r {
		case '(':
			depth++
		case ')':
			depth--
		case ',':
			if depth == 0 {
				parts = append(parts, s[start:i])
				start = i + 1
			}
		}
	}
	parts = append(parts, s[start:])
	return parts
}

func expandRepeat(tok string) []TrackSize {
	inner := strings.TrimSuffix(strings.TrimPrefix(tok, "repeat("), ")")
	parts := splitCSVOuter(inner)
	if len(parts) < 2 {
		return nil
	}
	countStr := strings.TrimSpace(parts[0])
	patternStr := strings.TrimSpace(strings.Join(parts[1:], ","))
	count, err := strconv.Atoi(countStr)
	if err != nil || count <= 0 {
		return nil
	}
	pattern := splitTrackTokens(patternStr)
	var one []TrackSize
	for _, p := range pattern {
		if t, ok := parseTrackSize(p); ok {
			one = append(one, t)
		}
	}
	out := make([]TrackSize, 0, len(one)*count)
	for range count {
		out = append(out, one...)
	}
	return out
}

// ParseGridLine parses a single grid-column or grid-row shorthand into
// (start, end) placements. Supports `N`, `N / M`, `span N`, `N / span M`,
// and `span M / N`.
func ParseGridLine(s string) (GridPlacement, GridPlacement) {
	s = strings.TrimSpace(s)
	var start, end GridPlacement
	if s == "" {
		return start, end
	}
	if slash := strings.IndexByte(s, '/'); slash >= 0 {
		start = parseGridToken(strings.TrimSpace(s[:slash]))
		end = parseGridToken(strings.TrimSpace(s[slash+1:]))
		return start, end
	}
	start = parseGridToken(s)
	return start, end
}

func parseGridToken(s string) GridPlacement {
	s = strings.TrimSpace(s)
	if s == "" || s == "auto" {
		return GridPlacement{}
	}
	if strings.HasPrefix(s, "span ") {
		n, err := strconv.Atoi(strings.TrimSpace(s[5:]))
		if err != nil || n <= 0 {
			return GridPlacement{}
		}
		return GridPlacement{Span: n}
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return GridPlacement{}
	}
	return GridPlacement{Start: n}
}

// ParseGridAutoFlow parses a grid-auto-flow value.
func ParseGridAutoFlow(s string) GridAutoFlow {
	s = strings.TrimSpace(s)
	switch s {
	case "column":
		return FlowColumn
	case "row dense", "dense":
		return FlowRowDense
	case "column dense":
		return FlowColumnDense
	}
	return FlowRow
}
