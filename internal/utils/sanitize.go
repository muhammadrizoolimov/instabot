package utils

import (
	"html"
	"strings"
	"unicode"
)

// SanitizeForTelegram removes or replaces characters that can cause Telegram API errors
func SanitizeForTelegram(s string) string {
	// First, unescape HTML entities
	s = html.UnescapeString(s)

	// Map function to filter problematic characters
	filtered := strings.Map(func(r rune) rune {
		// Keep alphanumeric, spaces, and common punctuation
		if unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsSpace(r) {
			return r
		}

		// Keep safe punctuation
		safe := ".,!?-_()[]'\"&:;@#"
		if strings.ContainsRune(safe, r) {
			return r
		}

		// Replace problematic characters with space
		return ' '
	}, s)

	// Clean up multiple spaces
	filtered = strings.Join(strings.Fields(filtered), " ")

	// Trim to reasonable length (Telegram has limits)
	if len(filtered) > 200 {
		filtered = filtered[:200] + "..."
	}

	return strings.TrimSpace(filtered)
}

// FormatHTML formats text for Telegram HTML parsing
func FormatHTML(text string) string {
	// Escape special HTML characters
	text = strings.ReplaceAll(text, "&", "&amp;")
	text = strings.ReplaceAll(text, "<", "&lt;")
	text = strings.ReplaceAll(text, ">", "&gt;")
	return text
}

// TruncateString truncates string to maxLen and adds ellipsis
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
