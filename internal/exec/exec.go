package exec

import (
	"context"
	"strings"

	"github.com/sven-seyfert/apiprobe/internal/auth"
	"github.com/sven-seyfert/apiprobe/internal/diff"
	"github.com/sven-seyfert/apiprobe/internal/fileutil"
	"github.com/sven-seyfert/apiprobe/internal/loader"
	"github.com/sven-seyfert/apiprobe/internal/logger"
	"github.com/sven-seyfert/apiprobe/internal/report"
	"github.com/sven-seyfert/apiprobe/internal/util"
)

// ProcessRequest executes the APIRequest (including optional test cases),
// compares the response against existing output, and triggers the webhook
// if differences are detected.
func ProcessRequest(
	ctx context.Context,
	idx int,
	req *loader.APIRequest,
	testCaseIndex *int,
	res *report.Result,
	rep *report.Report,
	tokenStore *auth.TokenStore,
	debugMode bool,
) {
	if testCaseIndex != nil {
		logger.NewLine()
		logger.Debugf("Run: %d, Test case: %d", idx, *testCaseIndex+1)
	}

	const noTestCaseIndicator = -1

	outputFile := fileutil.BuildOutputFilePath(req, testCaseIndex)

	response, statusCode, err := executeRequest(ctx, req, debugMode)
	if err != nil {
		logger.Errorf(`Failed endpoint request "%s": %v`, req.Request.Endpoint, err)
		res.IncreaseRequestErrorCount()

		if testCaseIndex != nil {
			rep.AddReportData(req, statusCode, outputFile, *testCaseIndex)
		} else {
			rep.AddReportData(req, statusCode, outputFile, noTestCaseIndicator)
		}

		return
	}

	result, err := formatResponse(ctx, req, response)
	if err != nil {
		logger.Errorf("Failed processing JSON query by JQ. Error: %v", err)
		res.IncreaseFormatErrorCount()

		if testCaseIndex != nil {
			rep.AddReportData(req, statusCode, outputFile, *testCaseIndex)
		} else {
			rep.AddReportData(req, statusCode, outputFile, noTestCaseIndicator)
		}

		return
	}

	if req.IsAuthRequest {
		addAuthTokenToTokenStore(result, tokenStore, req)

		logger.Debugf("No output file will be written (unnecessary), because generic token result.")

		return
	}

	hasChanged, err := diff.HasFileContentChanged(result, outputFile)
	if err != nil {
		logger.Errorf("%v", err)

		return
	}

	if !hasChanged {
		return
	}

	res.IncreaseChangedFilesCount()

	if testCaseIndex != nil {
		rep.AddReportData(req, statusCode, outputFile, *testCaseIndex)
	} else {
		rep.AddReportData(req, statusCode, outputFile, noTestCaseIndicator)
	}
}

// executeRequest wraps runCurl to perform the HTTP request defined by APIRequest
// and returns the raw response body and status code.
func executeRequest(ctx context.Context, req *loader.APIRequest, debugMode bool) ([]byte, string, error) {
	curlOutput, statusCode, err := runCurl(ctx, req, debugMode)
	if err != nil {
		return nil, statusCode, err
	}

	return curlOutput, statusCode, nil
}

// formatResponse formats the curl output using jq
// and returns the filtered result.
func formatResponse(ctx context.Context, req *loader.APIRequest, response []byte) ([]byte, error) {
	// If response is not JSON ("content-type: application/json"),
	// it's plain text and therefore there is no need for jq formatting.
	if !strings.HasPrefix(string(response), "{") && !strings.HasPrefix(string(response), "[") {
		return response, nil
	}

	jqOutput, err := GoJQ(ctx, req.JqCommand, response)
	if err != nil {
		return nil, err
	}

	return jqOutput, nil
}

// addAuthTokenToTokenStore attempts to add the token to the provided token store
// using the request ID as the key. Returns nothing.
func addAuthTokenToTokenStore(result []byte, tokenStore *auth.TokenStore, req *loader.APIRequest) {
	token := util.TrimQuotes(string(result))
	lastTokenChars := token[util.Max(0, len(token)-12):] //nolint:mnd

	if added := tokenStore.Add(req.ID, token); added {
		logger.Debugf(`Token "...%s" for auth request "%s" added to token store.`, lastTokenChars, req.ID)
	} else {
		logger.Warnf(`Token "...%s" for auth request "%s" already exists in token store.`, lastTokenChars, req.ID)
	}
}
