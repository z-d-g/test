package render

import "strings"

// renderImageAlt renders an image with its alt text.
func renderImageAlt(alt string, styles *styleCache) string {
	var b strings.Builder
	b.WriteString(styles.imageFunc("⊞ "))
	if alt != "" {
		b.WriteString(styles.imageFunc(alt))
	} else {
		b.WriteString(styles.imageFunc("[image]"))
	}
	return b.String()
}
