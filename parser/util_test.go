package parser

import (
	"testing"
)

func Test_readEscape(t *testing.T) {
	table := []struct {
		scenario string
		input    []byte
		expected string
	}{
		{
			scenario: "alarm bell",
			input:    []byte{'a', 'b', 'c', '\\', 'a', 'd', 'e', 'f'},
			expected: "abc\adef",
		},
		{
			scenario: "backspace",
			input:    []byte{'a', 'b', 'c', '\\', 'b', 'd', 'e', 'f'},
			expected: "abc\bdef",
		},
		{
			scenario: "form-feed",
			input:    []byte{'a', 'b', 'c', '\\', 'f', 'd', 'e', 'f'},
			expected: "abc\fdef",
		},
		{
			scenario: "newline",
			input:    []byte{'a', 'b', 'c', '\\', 'n', 'd', 'e', 'f'},
			expected: "abc\ndef",
		},
		{
			scenario: "carriage return",
			input:    []byte{'a', 'b', 'c', '\\', 'r', 'd', 'e', 'f'},
			expected: "abc\rdef",
		},
		{
			scenario: "horizontal tap",
			input:    []byte{'a', 'b', 'c', '\\', 't', 'd', 'e', 'f'},
			expected: "abc\tdef",
		},
		{
			scenario: "vertical tap",
			input:    []byte{'a', 'b', 'c', '\\', 'v', 'd', 'e', 'f'},
			expected: "abc\vdef",
		},
		{
			scenario: "double quote",
			input:    []byte{'a', 'b', 'c', '\\', '"', 'd', 'e', 'f'},
			expected: "abc\"def",
		},
		{
			scenario: "backslash",
			input:    []byte{'a', 'b', 'c', '\\', '\\', 'd', 'e', 'f'},
			expected: "abc\\def",
		},
		{
			scenario: "multiple escapes",
			input:    []byte{'a', 'b', 'c', '\\', 'b', '\\', '\\', 'd', 'e', 'f'},
			expected: "abc\b\\def",
		},
		{
			scenario: "ending with backslash",
			input:    []byte{'a', 'b', 'c', 'd', 'e', 'f', '\\', '\\'},
			expected: "abcdef\\",
		},
	}

	for _, test := range table {
		v := readEscape(string(test.input))

		if v != test.expected {
			t.Errorf("expected %q, got %q", test.expected, v)
		}
	}
}
