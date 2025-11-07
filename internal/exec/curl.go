package exec

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/sven-seyfert/apiprobe/internal/loader"
	"github.com/sven-seyfert/apiprobe/internal/logger"
)

// runCurl executes an external 'curl' command with specified timeouts
// and write-out flags, captures its stdout, splits the HTTP status code
// and returns the response body if the status code is 2xx;
// otherwise, returns an error.
func runCurl(ctx context.Context, req *loader.APIRequest, debugMode bool) ([]byte, string, string, error) {
	cmdArgs := req.CurlCmdArguments()

	var stdout, stderr bytes.Buffer

	cmd := exec.CommandContext(ctx, "./lib/curl.exe", cmdArgs...)

	if debugMode {
		fmt.Printf("\n%s\n\n", buildCurlFormat(cmd.String())) //nolint:forbidigo
	}

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	logger.Debugf(`Executing endpoint request "%s"`, req.Request.Endpoint)
	logger.Infof(`Description: "%s"`, req.Request.Description)

	start := time.Now()
	err := cmd.Run()
	duration := time.Since(start)
	rawOutput := stdout.Bytes()

	if err != nil {
		logger.Errorf("Curl execution failed. Error: %v", err)
		logger.Errorf("StdOut response: %s", stdout.String())
		logger.Errorf("StdErr response: %s", stderr.String())

		return nil, "", stderr.String(), fmt.Errorf("curl error: %w", err)
	}

	body, statusCode, err := extractStatusCode(rawOutput)
	if err != nil {
		return nil, "", stdout.String(), err
	}

	logger.Debugf("Status: %s, Duration: %dms", statusCode, duration.Milliseconds())

	if !strings.HasPrefix(statusCode, "2") {
		logger.Warnf("Non-2xx status code received: status %s", statusCode)
		logger.Warnf("Curl stdout response: %s", body)
		logger.Warnf("StdOut response: %s", stdout.String())
		logger.Warnf("StdErr response: %s", stderr.String())

		return nil, statusCode, string(body), fmt.Errorf("status %s", statusCode)
	}

	return body, statusCode, "", nil
}

// extractStatusCode splits the raw output from curl (where the last
// three bytes encode the HTTP status code) into the response body
// and the status code string.
func extractStatusCode(output []byte) ([]byte, string, error) {
	const HTTPCodeLength = 3

	if len(output) < HTTPCodeLength {
		logger.Warnf("Output too short to contain status code: only %d bytes", len(output))

		return nil, "", fmt.Errorf("only %d bytes", len(output))
	}

	body := output[:len(output)-HTTPCodeLength]
	statusCode := string(output[len(output)-HTTPCodeLength:])

	return body, statusCode, nil
}
