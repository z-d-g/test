package markdown

import "testing"

func TestFindClosingDelimiter(t *testing.T) {
	tests := []struct {
		name      string
		line      string
		start     int
		delimiter string
		want      int
	}{
		{"bold", "**hello**", 2, "**", 7},
		{"italic", "*hello*", 1, "*", 6},
		{"code", "`code`", 1, "`", 5},
		{"no close", "**hello", 2, "**", -1},
		{"empty delim", "hello", 0, "", -1},
		{"escaped", `bold \* star`, 7, "*", -1},
		{"underscore bold", "__hello__", 2, "__", 7},
		{"tilde", "~~strike~~", 2, "~~", 8},
		{"link text", "[text](url)", 1, "]", 5},
		{"link url", "[text](url)", 7, ")", 10},
		{"plain italic", "*ab*", 1, "*", 3},
		{"no match italic", "*abc", 1, "*", -1},
		{"backtick in code", "`a`b`", 1, "`", 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FindClosingDelimiter(tt.line, tt.start, tt.delimiter); got != tt.want {
				t.Errorf("FindClosingDelimiter(%q, %d, %q) = %d, want %d",
					tt.line, tt.start, tt.delimiter, got, tt.want)
			}
		})
	}
}
