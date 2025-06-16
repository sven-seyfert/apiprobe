package report

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/sven-seyfert/apiprobe/internal/crypto"
	"github.com/sven-seyfert/apiprobe/internal/db"
	"github.com/sven-seyfert/apiprobe/internal/loader"
	"github.com/sven-seyfert/apiprobe/internal/logger"
	"zombiezen.com/go/sqlite"
)

type Result struct {
	RequestErrorCount        int
	FormatResponseErrorCount int
	ChangedFilesCount        int
}

// IncreaseRequestErrorCount increments the Result counter for failed HTTP requests.
func (res *Result) IncreaseRequestErrorCount() {
	res.RequestErrorCount++
}

// IncreaseFormatErrorCount increments the Result counter for JSON formatting or jq errors.
func (res *Result) IncreaseFormatErrorCount() {
	res.FormatResponseErrorCount++
}

// IncreaseChangedFilesCount increments the Result counter for the number of output files that have changed.
func (res *Result) IncreaseChangedFilesCount() {
	res.ChangedFilesCount++
}

type Request struct {
	ID             string `json:"id"`
	Description    string `json:"description"`
	Endpoint       string `json:"endpoint"`
	StatusCode     string `json:"statusCode"`
	OutputFilePath string `json:"outputFilePath"`
	TestCase       string `json:"testCase"`
}

type Report struct {
	Requests []Request `json:"issues"`
}

// AddReportData records a single API request’s result, its ID, description,
// endpoint, status code, output file path and test case into the Report.
func (r *Report) AddReportData(req *loader.APIRequest, statusCode string, outputFilePath string, testCaseIndex int) {
	testCase := ""

	if testCaseIndex != -1 {
		testCase = req.TestCases[testCaseIndex]
	}

	request := Request{
		ID:             req.HexHash,
		Description:    req.Description,
		Endpoint:       req.Endpoint,
		StatusCode:     statusCode,
		OutputFilePath: outputFilePath,
		TestCase:       testCase,
	}

	r.Requests = append(r.Requests, request)
}

func (r *Report) SaveToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		logger.Errorf("Failure on create file. Error: %v", err)

		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Prettify json.

	if err := encoder.Encode(r); err != nil {
		logger.Errorf("Failure on write file. Error: %v", err)

		return err
	}

	return nil
}

// WebExWebhookNotification sends the given JSON payload to the configured
// WebEx incoming webhook URL.
func WebExWebhookNotification(ctx context.Context, conn *sqlite.Conn,
	webhookURL string, spaceSecret string,
	webhookPayload []byte) {
	spaceSecret = crypto.ExtractSecretHash(spaceSecret)

	spaceIdentifier, _ := db.SelectHash(conn, spaceSecret)
	url := webhookURL + crypto.Deobfuscate(spaceIdentifier)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(webhookPayload))
	if err != nil {
		logger.Errorf("Error on new request. Error: %v", err)

		return
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Errorf("Error on send request. Error: %v", err)

		return
	}
	defer resp.Body.Close()
}
