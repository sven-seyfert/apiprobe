package exec

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/sven-seyfert/apiprobe/internal/logger"
)

// RunJQ invokes the external 'jq' binary on the given JSON input.
// If no jq command is provided, it defaults to the identity filter '.'.
// Returns the filtered JSON bytes or an error if execution fails.
func RunJQ(ctx context.Context, jqArgs []string, input []byte) ([]byte, error) {
	const defaultJQPrettifyFilter = "."

	if len(jqArgs) == 0 {
		jqArgs = []string{defaultJQPrettifyFilter}
	}

	cmd := exec.CommandContext(ctx, "./lib/jq.exe", jqArgs...)
	cmd.Stdin = bytes.NewReader(input)

	var stderr bytes.Buffer

	cmd.Stderr = &stderr

	output, err := cmd.Output()
	if err != nil && len(output) == 0 {
		jqError := strings.ReplaceAll(stderr.String(), "\n", " ")

		logger.Errorf("JQ execution failed. Error: %v", jqError)

		return nil, fmt.Errorf("jq error: %w", err)
	}

	return output, nil
}

// TODO: Add functionality to replace JSON values for
// specific keys (dynamic content like datetimes).
