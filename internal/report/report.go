package report

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/sven-seyfert/apiprobe/internal/config"
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
	const noTestCaseIndicator = -1

	testCase := ""

	if testCaseIndex != noTestCaseIndicator {
		testCase = req.TestCases[testCaseIndex].Name
	}

	request := Request{
		ID:             req.ID,
		Description:    req.Request.Description,
		Endpoint:       req.Request.Endpoint,
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
	encoder.SetIndent("", "    ")

	if err := encoder.Encode(r); err != nil {
		logger.Errorf("Failure on write file. Error: %v", err)

		return err
	}

	return nil
}

// IsHeartbeatTime checks whether enough time has passed
// since the last heartbeat.
func IsHeartbeatTime(cfg *config.Config) (bool, error) {
	lastHeartbeatTime := cfg.Heartbeat.LastHeartbeatTime
	if lastHeartbeatTime == "" {
		return true, nil
	}

	threshold := time.Hour * time.Duration(cfg.Heartbeat.IntervalInHours)

	lastTime, err := time.Parse(time.RFC3339, lastHeartbeatTime)
	if err != nil {
		logger.Errorf(`Invalid datetime "%s". Error: %v\n`, lastHeartbeatTime, err)

		return false, err
	}

	diff := time.Since(lastTime)

	return diff >= threshold, nil
}

// UpdateHeartbeatTime writes the current UTC time (RFC3339) into
// cfg.Heartbeat.LastHeartbeatTime and persists the entire cfg back
// to the config JSON file.
func UpdateHeartbeatTime(cfg *config.Config) error {
	cfg.Heartbeat.LastHeartbeatTime = time.Now().UTC().Format(time.RFC3339)

	file, err := os.Create("./config/config.json")
	if err != nil {
		logger.Errorf("Failure on create file. Error: %v", err)

		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)

	encoder.SetIndent("", "    ")
	encoder.SetEscapeHTML(false)

	if err := encoder.Encode(cfg); err != nil {
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
	url := webhookURL + spaceSecret

	const secretPrefix = "<secret-"

	if strings.Contains(spaceSecret, secretPrefix) {
		spaceSecret = crypto.ExtractSecretHash(spaceSecret)
		spaceIdentifier, _ := db.SelectHash(conn, spaceSecret)
		url = webhookURL + crypto.Deobfuscate(spaceIdentifier)
	}

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
