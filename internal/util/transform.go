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

// Min returns the smaller of two integer values.
func Min(a, b int) int {
	if a < b {
		return a
	}

	return b
}
