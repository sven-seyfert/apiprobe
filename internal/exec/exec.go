package exec

import (
	"context"

	"github.com/sven-seyfert/apiprobe/internal/diff"
	"github.com/sven-seyfert/apiprobe/internal/fileutil"
	"github.com/sven-seyfert/apiprobe/internal/loader"
	"github.com/sven-seyfert/apiprobe/internal/logger"
	"github.com/sven-seyfert/apiprobe/internal/report"
)

// ProcessRequest executes the APIRequest (including optional test cases),
// compares the response against existing output, and triggers the webhook
// if differences are detected.
func ProcessRequest(
	ctx context.Context, idx int, req *loader.APIRequest, testCaseIndex *int,
	res *report.Result, rep *report.Report) {
	if testCaseIndex != nil {
		logger.NewLine()
		logger.Debugf("Run: %d, Test case: %d", idx, *testCaseIndex+1)
	}

	const noTestCaseIndicator = -1

	outputFile := fileutil.BuildOutputFilePath(req, testCaseIndex)

	response, statusCode, err := executeRequest(ctx, req)
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
func executeRequest(ctx context.Context, req *loader.APIRequest) ([]byte, string, error) {
	curlOutput, statusCode, err := runCurl(ctx, req)
	if err != nil {
		return nil, statusCode, err
	}

	return curlOutput, statusCode, nil
}

// formatResponse formats the curl output using jq
// and returns the filtered result.
func formatResponse(ctx context.Context, req *loader.APIRequest, response []byte) ([]byte, error) {
	jqArgs := []string{req.JqCommand}
	jqOutput, err := RunJQ(ctx, jqArgs, response)

	if err != nil {
		return nil, err
	}

	return jqOutput, nil
}
