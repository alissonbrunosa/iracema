package parser

import (
	"strings"
)

func unescapeString(s string) string {
	var buf strings.Builder

	for i := 0; i < len(s); i++ {
		ch := s[i]
		if ch != '\\' {
			buf.WriteByte(ch)
			continue
		}

		i++
		if i >= len(s) {
			break
		}

		switch s[i] {
		case 'a':
			buf.WriteByte('\a')
		case 'b':
			buf.WriteByte('\b')
		case 'f':
			buf.WriteByte('\f')
		case 'n':
			buf.WriteByte('\n')
		case 'r':
			buf.WriteByte('\r')
		case 't':
			buf.WriteByte('\t')
		case 'v':
			buf.WriteByte('\v')
		case '\\':
			buf.WriteByte('\\')
		case '"':
			buf.WriteByte('"')
		default:
			// pretty sure we should never get here,
			// since we check all valid escape character
			panic("Damn! That's a bug!")
		}
	}

	return buf.String()
}
