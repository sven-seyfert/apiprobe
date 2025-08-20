package exec

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/itchyny/gojq"

	"github.com/sven-seyfert/apiprobe/internal/logger"
)

// GoJQ executes the jq query given by jqArgs against inputJSON.
// Returns the encoded JSON ([]byte) of the query result or an error.
func GoJQ(ctx context.Context, jqArgs string, inputJSON []byte) ([]byte, error) {
	const defaultJQPrettifyFilter = "."

	if jqArgs == "" {
		jqArgs = defaultJQPrettifyFilter
	}

	var input any
	if err := json.Unmarshal(inputJSON, &input); err != nil {
		logger.Errorf("Failed to unmarshal input (invalid input json). Error: %v", err)

		return nil, err
	}

	code, err := compileQuery(jqArgs)
	if err != nil {
		return nil, err
	}

	results, err := runQuery(ctx, code, input)
	if err != nil {
		return nil, err
	}

	return encodeResults(results)
}

// compileQuery parses and compiles the provided jqArgs into *gojq.Code.
// Returns the compiled code or an error if parsing/compilation fails.
func compileQuery(jqArgs string) (*gojq.Code, error) {
	query, err := gojq.Parse(jqArgs)
	if err != nil {
		logger.Errorf("Failed to parse jqArgs. Error: %v", err)

		return nil, err
	}

	code, err := gojq.Compile(query)
	if err != nil {
		logger.Errorf("Failed to compile jq. Error: %v", err)

		return nil, err
	}

	return code, nil
}

// runQuery executes the compiled gojq code with the provided context and input,
// collects all produced values and handles gojq.HaltError specially.
// Returns a slice of results ([]any) or an error.
func runQuery(ctx context.Context, code *gojq.Code, input any) ([]any, error) {
	iter := code.RunWithContext(ctx, input)
	results := []any{}

	for {
		nextVal, isOk := iter.Next()
		if !isOk {
			break
		}

		// If the iterator produced a non-error value,
		// append and continue early.
		errVal, isErr := nextVal.(error)
		if !isErr {
			results = append(results, nextVal)

			continue
		}

		// Handle error values
		if halt, ok := errVal.(*gojq.HaltError); ok { //nolint:errorlint
			if halt.Value() == nil {
				break
			}

			// If HaltError carries a value,
			// append it and continue.
			if ve, okay := errVal.(interface{ Value() any }); okay {
				results = append(results, ve.Value())

				continue
			}

			logger.Errorf("Failed to run query. JQ halt error with non-nil value (type=%T).", errVal)

			return nil, errVal
		}

		logger.Errorf("Failed to run query. JQ runtime error: %s (type=%T)", safeErrorString(errVal), errVal)

		return nil, errVal
	}

	return results, nil
}

// encodeResults marshals results to indented JSON. If results contains exactly
// one element, that element is marshaled directly; otherwise the whole slice
// is marshaled. Returns the JSON bytes ([]byte) or an error.
func encodeResults(results []any) ([]byte, error) {
	var out any

	if len(results) == 1 {
		out = results[0]
	} else {
		out = results
	}

	enc, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		logger.Errorf("Failed to marshal indent. Marshal output: %v", err)

		return nil, err
	}

	return enc, nil
}

// safeErrorString returns err.Error(), but recovers and returns a fallback
// string if calling Error() panics.
func safeErrorString(err error) string {
	var errMsg string

	func() {
		defer func() {
			if rec := recover(); rec != nil {
				errMsg = fmt.Sprintf("error.Error() panicked: %v (error type=%T)", rec, err)
			}
		}()

		errMsg = err.Error()
	}()

	return errMsg
}
