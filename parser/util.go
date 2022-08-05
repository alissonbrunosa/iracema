package parser

import (
	"strings"
)

func readEscape(s string) string {
	var buf strings.Builder
	buf.Grow(len(s))

	for i := 0; i < len(s); i++ {
		ch := s[i]
		if ch != '\\' {
			buf.WriteByte(ch)
			continue
		}

		next := i + 1
		switch s[next] {
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

		case '"':
			buf.WriteByte('"')

		default:
			buf.WriteByte('\\')
		}
		i++
	}

	return buf.String()
}
