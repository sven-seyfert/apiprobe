package util //nolint:revive

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
// The comparison is case-insensitive and returns true if found, otherwise false.
func ContainsSubstring(slice []string, substr string) bool {
	for _, value := range slice {
		if strings.Contains(strings.ToLower(value), strings.ToLower(substr)) {
			return true
		}
	}

	return false
}
