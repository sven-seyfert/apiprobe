package util

import "strings"

// TrimQuotes removes leading and trailing double quotes and trailing
// CRLF from the given string. Returns the cleaned string.
func TrimQuotes(value string) string {
	value = strings.TrimPrefix(value, `"`)
	value = strings.TrimSuffix(value, "\r\n")
	value = strings.TrimSuffix(value, `"`)

	return value
}

// ContainsSubstring checks if any string in the slice contains the given substring.
// Returns true if found, otherwise false.
func ContainsSubstring(slice []string, substr string) bool {
	for _, value := range slice {
		if strings.Contains(value, substr) {
			return true
		}
	}

	return false
}
