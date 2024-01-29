package utils

import "regexp"

func SanitizeFileName(input string) string {
	// Replace special characters with underscores
	re := regexp.MustCompile(`[^\w]`)
	sanitized := re.ReplaceAllString(input, "_")
	return sanitized
}