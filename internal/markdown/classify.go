package markdown

import "strings"

// IsCodeFence reports whether line is a fenced code block delimiter.
func IsCodeFence(line string) bool {
	trimmed := strings.TrimSpace(line)
	return strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~")
}

// CodeFenceChar returns the fence character (` or ~) of a code fence line,
// or 0 if the line is not a fence.
func CodeFenceChar(line string) byte {
	trimmed := strings.TrimSpace(line)
	if len(trimmed) < 3 {
		return 0
	}
	switch trimmed[0] {
	case '`', '~':
		return trimmed[0]
	default:
		return 0
	}
}

// IsHorizontalRule reports whether line is a horizontal rule.
func IsHorizontalRule(line string) bool {
	trimmed := strings.TrimSpace(line)
	return trimmed == "---" || trimmed == "***" || trimmed == "___"
}

// IsListLine reports whether line starts a list item.
func IsListLine(line string) bool {
	trimmed := strings.TrimSpace(line)
	if len(trimmed) > 1 && ((trimmed[0] == '-' || trimmed[0] == '*' || trimmed[0] == '+') && trimmed[1] == ' ') {
		return true
	}
	if len(trimmed) > 2 && trimmed[0] >= '0' && trimmed[0] <= '9' {
		i := 1
		for i < len(trimmed) && trimmed[i] >= '0' && trimmed[i] <= '9' {
			i++
		}
		if i < len(trimmed)-1 && trimmed[i] == '.' && trimmed[i+1] == ' ' {
			return true
		}
	}
	return false
}

// IsHeadingLine reports whether line is a heading.
func IsHeadingLine(line string) bool {
	trimmed := strings.TrimSpace(line)
	if len(trimmed) == 0 {
		return false
	}
	if trimmed[0] == '#' {
		hashCount := 0
		for i := 0; i < len(trimmed) && i < 6; i++ {
			if trimmed[i] == '#' {
				hashCount++
			} else {
				break
			}
		}
		return hashCount > 0 && hashCount < len(trimmed) && trimmed[hashCount] == ' '
	}
	return false
}

// IsBlockquoteLine reports whether line starts a blockquote.
func IsBlockquoteLine(line string) bool {
	trimmed := strings.TrimSpace(line)
	return strings.HasPrefix(trimmed, ">")
}

// IsTableLine reports whether line is part of a table.
// Requires pipes with surrounding whitespace or at line boundaries,
// and at least 2 pipe-delimited fields.
func IsTableLine(line string) bool {
	if strings.Count(line, "|") < 2 {
		return false
	}
	// Lines starting/ending with | are likely tables
	trimmed := strings.TrimSpace(line)
	if strings.HasPrefix(trimmed, "|") || strings.HasSuffix(trimmed, "|") {
		return true
	}
	// For mid-line pipes, require whitespace around at least one pipe
	// to avoid matching shell pipelines like `foo | bar | baz`
	for i := 0; i < len(line)-1; i++ {
		if line[i] == '|' && i > 0 && i < len(line)-1 {
			if (line[i-1] == ' ' || line[i-1] == '\t') &&
				(line[i+1] == ' ' || line[i+1] == '\t' || line[i+1] == '-') {
				return true
			}
		}
	}
	return false
}

// IsTableSeparatorLine reports whether line is a table separator (e.g., |---|---|).
func IsTableSeparatorLine(line string) bool {
	trimmed := strings.TrimSpace(line)
	if !strings.Contains(trimmed, "|") {
		return false
	}
	// Remove all pipe chars, whitespace, dashes, and colons
	// If anything remains, it's not a separator
	for _, r := range trimmed {
		if r != '|' && r != '-' && r != ':' && r != ' ' && r != '\t' {
			return false
		}
	}
	// Must have at least one dash
	return strings.Contains(trimmed, "-")
}

// IsEmptyLine reports whether line is blank.
func IsEmptyLine(line string) bool {
	return len(strings.TrimSpace(line)) == 0
}

func ClassifyLine(line string, isInCodeBlock bool) int {
	if isInCodeBlock {
		return LineCodeContent
	}

	trimmed := strings.TrimSpace(line)

	if strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~") {
		return LineCodeFence
	}

	if strings.HasPrefix(trimmed, ">") {
		return LineBlockQuote
	}

	hashCount := CountLeadingHashes(line)
	if hashCount > 0 && len(trimmed) > hashCount && trimmed[hashCount] == ' ' {
		return hashCount
	}

	return LineNormal
}

// CountBlockquoteDepth returns the nesting depth of blockquote markers (> > → 2).
func CountBlockquoteDepth(line string) int {
	count := 0
	rest := strings.TrimSpace(line)
	for strings.HasPrefix(rest, ">") {
		count++
		rest = strings.TrimPrefix(rest, ">")
		rest = strings.TrimSpace(rest)
	}
	return count
}

// CountLeadingHashes returns the number of leading # characters (up to 6).
func CountLeadingHashes(line string) int {
	count := 0
	trimmed := strings.TrimSpace(line)
	for i := 0; i < len(trimmed) && i < 6; i++ {
		if trimmed[i] == '#' {
			count++
		} else {
			break
		}
	}
	return count
}
