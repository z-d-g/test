package markdown

import "testing"

func TestIsCodeFence(t *testing.T) {
	tests := []struct {
		line string
		want bool
	}{
		{"```go", true},
		{"```", true},
		{"~~~python", true},
		{"~~~", true},
		{"  ```  ", true},
		{"plain text", false},
		{"`single`", false},
		{"", false},
	}
	for _, tt := range tests {
		if got := IsCodeFence(tt.line); got != tt.want {
			t.Errorf("IsCodeFence(%q) = %v, want %v", tt.line, got, tt.want)
		}
	}
}

func TestIsHorizontalRule(t *testing.T) {
	tests := []struct {
		line string
		want bool
	}{
		{"---", true},
		{"***", true},
		{"___", true},
		{"  ---  ", true},
		{"--", false},
		{"- -", false},
		{"text", false},
		{"", false},
	}
	for _, tt := range tests {
		if got := IsHorizontalRule(tt.line); got != tt.want {
			t.Errorf("IsHorizontalRule(%q) = %v, want %v", tt.line, got, tt.want)
		}
	}
}

func TestIsListLine(t *testing.T) {
	tests := []struct {
		line string
		want bool
	}{
		{"- item", true},
		{"* item", true},
		{"+ item", true},
		{"1. item", true},
		{"42. item", true},
		{"-item", false},
		{"*", false},
		{"1.", false},
		{"plain text", false},
		{"", false},
	}
	for _, tt := range tests {
		if got := IsListLine(tt.line); got != tt.want {
			t.Errorf("IsListLine(%q) = %v, want %v", tt.line, got, tt.want)
		}
	}
}

func TestIsHeadingLine(t *testing.T) {
	tests := []struct {
		line string
		want bool
	}{
		{"# H1", true},
		{"## H2", true},
		{"###### H6", true},
		{"#NoSpace", false},
		{"####### seven", false},
		{"plain", false},
		{"", false},
		{"  # Indented", true},
	}
	for _, tt := range tests {
		if got := IsHeadingLine(tt.line); got != tt.want {
			t.Errorf("IsHeadingLine(%q) = %v, want %v", tt.line, got, tt.want)
		}
	}
}

func TestIsBlockquoteLine(t *testing.T) {
	tests := []struct {
		line string
		want bool
	}{
		{"> quote", true},
		{">> nested", true},
		{"plain", false},
		{"", false},
	}
	for _, tt := range tests {
		if got := IsBlockquoteLine(tt.line); got != tt.want {
			t.Errorf("IsBlockquoteLine(%q) = %v, want %v", tt.line, got, tt.want)
		}
	}
}

func TestIsTableLine(t *testing.T) {
	tests := []struct {
		line string
		want bool
	}{
		{"| a | b |", true},
		{"a|b|c", false},
		{"a|b", false},
		{"a | b | c", true},
		{"|---|---|", true},
		{"echo foo", false},
		{"plain", false},
		{"", false},
	}
	for _, tt := range tests {
		if got := IsTableLine(tt.line); got != tt.want {
			t.Errorf("IsTableLine(%q) = %v, want %v", tt.line, got, tt.want)
		}
	}
}

func TestIsTableSeparatorLine(t *testing.T) {
	tests := []struct {
		line string
		want bool
	}{
		{"|---|---|", true},
		{"| --- | --- |", true},
		{"|:---:|:---|", true},
		{"| a | b |", false},
		{"plain", false},
		{"", false},
	}
	for _, tt := range tests {
		if got := IsTableSeparatorLine(tt.line); got != tt.want {
			t.Errorf("IsTableSeparatorLine(%q) = %v, want %v", tt.line, got, tt.want)
		}
	}
}

func TestIsEmptyLine(t *testing.T) {
	tests := []struct {
		line string
		want bool
	}{
		{"", true},
		{"   ", true},
		{"\t", true},
		{"text", false},
	}
	for _, tt := range tests {
		if got := IsEmptyLine(tt.line); got != tt.want {
			t.Errorf("IsEmptyLine(%q) = %v, want %v", tt.line, got, tt.want)
		}
	}
}

func TestClassifyLine(t *testing.T) {
	tests := []struct {
		line          string
		isInCodeBlock bool
		want          int
	}{
		{"plain text", false, LineNormal},
		{"  ```  ", false, LineCodeFence},
		{"~~~", false, LineCodeFence},
		{"> quote", false, LineBlockQuote},
		{"# H1", false, LineHeading1},
		{"## H2", false, LineHeading2},
		{"### H3", false, LineHeading3},
		{"#### H4", false, LineHeading4},
		{"##### H5", false, LineHeading5},
		{"###### H6", false, LineHeading6},
		{"plain in code", true, LineCodeContent},
		{"# not heading in code", true, LineCodeContent},
	}
	for _, tt := range tests {
		if got := ClassifyLine(tt.line, tt.isInCodeBlock); got != tt.want {
			t.Errorf("ClassifyLine(%q, %v) = %d, want %d", tt.line, tt.isInCodeBlock, got, tt.want)
		}
	}
}

func TestCountBlockquoteDepth(t *testing.T) {
	tests := []struct {
		line string
		want int
	}{
		{"> quote", 1},
		{">> nested", 2},
		{"> > spaced", 2},
		{"plain", 0},
		{"", 0},
	}
	for _, tt := range tests {
		if got := CountBlockquoteDepth(tt.line); got != tt.want {
			t.Errorf("CountBlockquoteDepth(%q) = %d, want %d", tt.line, got, tt.want)
		}
	}
}

func TestCountLeadingHashes(t *testing.T) {
	tests := []struct {
		line string
		want int
	}{
		{"# one", 1},
		{"## two", 2},
		{"###### six", 6},
		{"plain", 0},
		{"", 0},
	}
	for _, tt := range tests {
		if got := CountLeadingHashes(tt.line); got != tt.want {
			t.Errorf("CountLeadingHashes(%q) = %d, want %d", tt.line, got, tt.want)
		}
	}
}
