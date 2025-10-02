package exec

import "strings"

type token struct {
	option   string
	value    string
	hasValue bool
}

// buildCurlFormat formats a given input string into a multi-line,
// indented curl command. Returns the formatted curl command as a
// string or an empty string if input is invalid.
func buildCurlFormat(input string) string {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return ""
	}

	const (
		executable = "curl"
		indent     = "     "
		backslash  = " \\"
	)

	tokens := parseTokens(parts[1:])
	isFirstOption := true

	var lines []string

	for _, tok := range tokens {
		if !tok.hasValue {
			if isFirstOption {
				lines = append(lines, executable+" "+tok.option+backslash)
				isFirstOption = false
			} else {
				lines = append(lines, indent+tok.option+backslash)
			}

			continue
		}

		formattedValue := quoteValue(tok.option, tok.value)

		if isFirstOption {
			lines = append(lines, executable+" "+tok.option+" "+formattedValue+backslash)
			isFirstOption = false
		} else {
			lines = append(lines, indent+tok.option+" "+formattedValue+backslash)
		}
	}

	if len(lines) == 0 {
		return executable
	}

	result := strings.Join(lines, "\n")
	result = strings.TrimRight(result, backslash)

	return result
}

// parseTokens parses command line parts into tokens with options and values.
// Returns a slice of token structs.
func parseTokens(parts []string) []token {
	var tokens []token

	valueFlags := map[string]struct{}{
		"--request":         {},
		"--connect-timeout": {},
		"--max-time":        {},
		"--url":             {},
		"--write-out":       {},
		"--data":            {},
		"--user":            {},
		"--header":          {},
	}

	for idx := 0; idx < len(parts); idx++ {
		part := parts[idx]

		if !strings.HasPrefix(part, "-") {
			continue
		}

		if _, ok := valueFlags[part]; !ok {
			tokens = append(tokens, token{option: part, value: "", hasValue: false})
			continue
		}

		value := ""
		nextIndex := idx + 1

		for nextIndex < len(parts) && !strings.HasPrefix(parts[nextIndex], "-") {
			if value != "" {
				value += " "
			}

			value += parts[nextIndex]
			nextIndex++
		}

		tokens = append(tokens, token{option: part, value: value, hasValue: true})
		idx = nextIndex - 1
	}

	return tokens
}

// quoteValue adds appropriate quoting to a flag value based on the flag type.
// Returns the quoted value as a string.
func quoteValue(flag, value string) string {
	switch flag {
	case "--request", "--connect-timeout", "--max-time":
		return value
	default:
		return "'" + escapeSingleQuotes(value) + "'"
	}
}

// escapeSingleQuotes escapes single quotes in a string for shell safety.
// Returns the escaped string.
func escapeSingleQuotes(s string) string {
	return strings.ReplaceAll(s, "'", "'\\''")
}
