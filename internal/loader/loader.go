package loader

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/sven-seyfert/apiprobe/internal/logger"
	"github.com/sven-seyfert/apiprobe/internal/util"
)

// APIRequest represents the structure of each API request definition
// as specified in the input JSON configuration.
type APIRequest struct {
	ID            string      `json:"id"`
	IsAuthRequest bool        `json:"isAuthRequest"`
	PreRequestID  string      `json:"preRequestId"`
	Request       Request     `json:"request"`
	TestCases     []TestCases `json:"testCases"`
	Tags          []string    `json:"tags"`
	JqCommand     string      `json:"jq"`

	// Relative JSON file path.
	JSONFilePath string `json:"-"`
}

// Request holds the HTTP-specific details for an API request.
type Request struct {
	Description string          `json:"description"`
	Method      string          `json:"method"`
	BaseURL     string          `json:"url"`
	Endpoint    string          `json:"endpoint"`
	BasicAuth   string          `json:"basicAuth"`
	Headers     []string        `json:"headers"`
	Params      []string        `json:"params"`
	PostBodyRaw json.RawMessage `json:"postBody"`

	// Target data type for the POST body format is string.
	PostBody string `json:"-"`
}

// TestCases defines the input variations for the requests.
type TestCases struct {
	Name            string          `json:"name"`
	ParamsData      string          `json:"paramsData"`
	PostBodyDataRaw json.RawMessage `json:"postBodyData"`

	// Target data type for the POST body format is string.
	PostBodyData string `json:"-"`
}

// PreparePostBody prepares the request body (empty, x-www-form-urlencoded
// or compacted JSON). Returns nil on success or an error if JSON compaction fails.
func (req *APIRequest) PreparePostBody() error {
	const emptyPostBodyLength = 2

	if len(string(req.Request.PostBodyRaw)) == emptyPostBodyLength {
		req.Request.PostBody = ""

		return nil
	}

	// Case POST body form is "x-www-form-urlencoded" which is no JSON.
	if util.ContainsSubstring(req.Request.Headers, "x-www-form-urlencoded") {
		req.Request.PostBody = util.TrimQuotes(string(req.Request.PostBodyRaw))

		return nil
	}

	var buf bytes.Buffer
	if err := json.Compact(&buf, req.Request.PostBodyRaw); err != nil {
		logger.Errorf("Failed by attempting JSON compact. Error: %v", err)

		return err
	}

	req.Request.PostBody = buf.String()

	return nil
}

// PreparePostBodyData processes the raw POST body data of all test cases.
// It normalizes the content based on header type, compacts JSON when needed,
// and sets the processed result into PostBodyData.
// Returns an error if JSON compaction fails, otherwise nil.
func (req *APIRequest) PreparePostBodyData() error {
	for idx := range req.TestCases {
		testCase := &req.TestCases[idx]

		const emptyPostBodyDataLength = 2

		if len(string(testCase.PostBodyDataRaw)) == emptyPostBodyDataLength {
			testCase.PostBodyData = ""

			continue
		}

		// Case POST body form is "x-www-form-urlencoded" which is no JSON.
		if util.ContainsSubstring(req.Request.Headers, "x-www-form-urlencoded") {
			testCase.PostBodyData = util.TrimQuotes(string(testCase.PostBodyDataRaw))

			continue
		}

		var buf bytes.Buffer
		if err := json.Compact(&buf, testCase.PostBodyDataRaw); err != nil {
			logger.Errorf("Failed by attempting JSON compact for test case %d. Error: %v", idx, err)

			return err
		}

		testCase.PostBodyData = buf.String()
	}

	return nil
}

// BuildRequestURL constructs the full request URL by concatenating the BaseURL,
// Endpoint, and optional query parameters defined in the APIRequest.
func (req *APIRequest) BuildRequestURL() string {
	var requestURL strings.Builder

	requestURL.WriteString(req.Request.BaseURL)
	requestURL.WriteString(req.Request.Endpoint)

	if len(req.Request.Params) > 0 {
		requestURL.WriteString("?")
		requestURL.WriteString(url.PathEscape(strings.Join(req.Request.Params, "&")))
	}

	return requestURL.String()
}

// CurlCmdArguments builds the command-line arguments for a curl invocation
// based on the HTTP method, URL, headers, authentication and payload
// specified in the APIRequest.
func (req *APIRequest) CurlCmdArguments() []string {
	cmdArgs := []string{
		"--request", req.Request.Method,
		"--silent", "--location", "--insecure",
		"--connect-timeout", "8",
		"--max-time", "24",
		"--url", req.BuildRequestURL(),
		"--write-out", "%{http_code}",
	}

	if req.Request.Method == http.MethodGet {
		cmdArgs = append(cmdArgs, "--get")
	}

	if req.Request.Method == http.MethodPost && req.Request.PostBody != "" {
		postBody := req.Request.PostBody

		// Encoding for POST body form "x-www-form-urlencoded".
		if util.ContainsSubstring(req.Request.Headers, "x-www-form-urlencoded") {
			postBody = url.PathEscape(req.Request.PostBody)
		}

		cmdArgs = append(cmdArgs, "--data", postBody)
	}

	if req.Request.BasicAuth != "" {
		cmdArgs = append(cmdArgs, "--user", req.Request.BasicAuth)
	}

	for _, h := range req.Request.Headers {
		cmdArgs = append(cmdArgs, "--header", h)
	}

	return cmdArgs
}

// LoadAllRequests recursively walks the input directory, parses all JSON files
// and returns APIRequest pointers.
func LoadAllRequests() ([]*APIRequest, error) {
	const inputDir = "./data/input"

	var requests []*APIRequest

	err := filepath.Walk(inputDir, func(path string, _ os.FileInfo, err error) error {
		if err != nil {
			logger.Errorf("Failed to walk path. Error: %v", err)

			return err
		}

		if filepath.Ext(path) == ".json" {
			fileRequest, loadErr := loadRequestFromFile(path, inputDir)
			if loadErr != nil {
				return loadErr
			}

			requests = append(requests, fileRequest...)
		}

		return nil
	})

	return requests, err
}

// loadRequestFromFile reads a JSON file, unmarshals it into APIRequest structs
// and assigns the file path.
func loadRequestFromFile(path string, inputDir string) ([]*APIRequest, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		logger.Errorf(`Failed to read file "%s". Error: %v`, path, err)

		return nil, err
	}

	var requestData []APIRequest

	if err = json.Unmarshal(bytes, &requestData); err != nil {
		logger.Errorf(`Failed to unmarshal JSON "%s". Error: %v`, path, err)

		return nil, err
	}

	// Store JSON file path in each request (relative to ./data/input).
	relPath, err := filepath.Rel(inputDir, path)
	if err != nil {
		logger.Errorf(`Failed to get relative path "%s". Error: %v`, path, err)

		return nil, err
	}

	request := make([]*APIRequest, len(requestData))

	for idx := range requestData {
		requestData[idx].JSONFilePath = relPath
		request[idx] = &requestData[idx]
	}

	return request, nil
}
