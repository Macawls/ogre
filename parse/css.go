package parse

import "strings"

// ParseStyle parses a CSS style string into a map of property-value pairs.
// ParseStyle parses a CSS inline style string into property-value pairs.
func ParseStyle(s string) map[string]string {
	result := make(map[string]string)
	for _, decl := range splitDeclarations(s) {
		decl = strings.TrimSpace(decl)
		if decl == "" {
			continue
		}
		idx := findFirstColon(decl)
		if idx < 0 {
			continue
		}
		prop := strings.TrimSpace(decl[:idx])
		val := strings.TrimSpace(decl[idx+1:])
		if prop == "" || val == "" {
			continue
		}
		result[strings.ToLower(prop)] = val
	}
	return result
}

func splitDeclarations(s string) []string {
	var parts []string
	var buf strings.Builder
	depth := 0
	inSingle := false
	inDouble := false

	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c == '\'' && !inDouble:
			inSingle = !inSingle
			buf.WriteByte(c)
		case c == '"' && !inSingle:
			inDouble = !inDouble
			buf.WriteByte(c)
		case !inSingle && !inDouble && c == '(':
			depth++
			buf.WriteByte(c)
		case !inSingle && !inDouble && c == ')':
			if depth > 0 {
				depth--
			}
			buf.WriteByte(c)
		case c == ';' && depth == 0 && !inSingle && !inDouble:
			parts = append(parts, buf.String())
			buf.Reset()
		default:
			buf.WriteByte(c)
		}
	}
	if buf.Len() > 0 {
		parts = append(parts, buf.String())
	}
	return parts
}

func findFirstColon(s string) int {
	depth := 0
	inSingle := false
	inDouble := false

	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c == '\'' && !inDouble:
			inSingle = !inSingle
		case c == '"' && !inSingle:
			inDouble = !inDouble
		case !inSingle && !inDouble && c == '(':
			depth++
		case !inSingle && !inDouble && c == ')':
			if depth > 0 {
				depth--
			}
		case c == ':' && depth == 0 && !inSingle && !inDouble:
			return i
		}
	}
	return -1
}
