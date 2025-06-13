package util

import (
	"strings"
)

// ReplaceQueryParam returns a copy of the given params slice, with the
// key from testCaseValue replaced if present, or appended otherwise.
func ReplaceQueryParam(params []string, testCaseValue string) []string {
	const subStringCount = 2

	keyToReplace := strings.SplitN(testCaseValue, "=", subStringCount)[0]
	replaced := false

	newParams := make([]string, len(params))
	copy(newParams, params)

	for idx, param := range newParams {
		parts := strings.SplitN(param, "=", subStringCount)

		if len(parts) == subStringCount && parts[0] == keyToReplace {
			newParams[idx] = testCaseValue
			replaced = true

			break
		}
	}

	if !replaced {
		newParams = append(newParams, testCaseValue) //nolint:makezero
	}

	return newParams
}
