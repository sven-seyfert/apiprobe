package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/sven-seyfert/apiprobe/internal/config"
	"github.com/sven-seyfert/apiprobe/internal/crypto"
	"github.com/sven-seyfert/apiprobe/internal/db"
	"github.com/sven-seyfert/apiprobe/internal/exec"
	"github.com/sven-seyfert/apiprobe/internal/flags"
	"github.com/sven-seyfert/apiprobe/internal/loader"
	"github.com/sven-seyfert/apiprobe/internal/logger"
	"github.com/sven-seyfert/apiprobe/internal/report"
	"github.com/sven-seyfert/apiprobe/internal/util"
	"zombiezen.com/go/sqlite"
)

// main initializes the logger and database, parses command-line flags loads
// configuration and seeds the database. It then loads and filters API request
// definitions, injects secrets, establishes a cancellation-aware context,
// processes each request and finally sends notifications based on errors
// or detected changes.
func main() {
	if err := logger.Init(); err != nil {
		logger.Fatalf("Program exits: Failed to initialize logger.")

		return
	}

	conn, err := db.Init()
	if err != nil {
		logger.Fatalf("Program exits: Failed to initialize database.")

		return
	}
	defer conn.Close()

	cliFlags := flags.Init()

	// Load config file and values.
	cfg, err := config.Load("./config/config.json")
	if err != nil {
		logger.Fatalf("Program exits: Failed to load config file.")

		return
	}

	// Handle command-line flags.
	complete, err := flags.IsNewID(*cliFlags.NewID)
	if complete || err != nil {
		return
	}

	complete, err = flags.IsAddSecret(*cliFlags.AddSecret, conn)
	if complete || err != nil {
		return
	}

	// Fill database with default seed data.
	err = db.InsertSeedData(conn)
	if err != nil {
		logger.Fatalf("Program exits: Failed to fill database with seed default data.")

		return
	}

	// LoadAllRequests loads all API request definitions from JSON files in the input directory.
	requests, err := loader.LoadAllRequests()
	if err != nil {
		logger.Fatalf("Program exits: Failed to load API request definitions.")

		return
	}

	// Exclude requests based on IDs.
	requests = loader.ExcludeRequestsByID(requests, *cliFlags.Exclude)

	// Filter requests based on single id (ten character long hex hash) or by flags.
	filteredRequests, notFound := filterRequests(requests, *cliFlags.ID, *cliFlags.Tags)
	if notFound {
		return
	}

	// Replace secrets placeholders in the requests with actual values.
	preparedRequests, err := crypto.HandleSecrets(filteredRequests, conn)
	if err != nil {
		logger.Fatalf("Program exits: Failed to handle secrets in requests.")

		return
	}

	// Only once requests are loaded successfully, set up signal-cancellation context.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Process each API request, optionally with test case variations.
	res, rep := processRequests(ctx, preparedRequests)

	// Send notification on error case or on changes.
	notification(ctx, cfg, conn, res, rep)
}

// filterRequests filters the given slice of APIRequest by the '--id'
// and '--tags' flags. It returns a slice of matching requests and a
// boolean flag that is true if no requests matched the filters.
func filterRequests(requests []*loader.APIRequest, id string, tags string) ([]*loader.APIRequest, bool) { //nolint:varnamelen
	if len(requests) == 0 {
		logger.Warnf(`No requests found.`)

		return requests, true
	}

	// Filter requests by ID.
	if id != "" {
		if req := loader.FilterByID(requests, id); req != nil {
			return []*loader.APIRequest{req}, false
		}

		logger.Warnf(`No request with id (hex hash) "%s" found.`, id)

		return requests, true
	}

	// Or filter requests by tags.
	if tags != "" {
		tagsList := strings.Split(tags, ",")
		wantedTags := make([]string, 0, len(tagsList))

		for _, tag := range tagsList {
			tag = strings.TrimSpace(tag)
			if tag != "" {
				wantedTags = append(wantedTags, tag)
			}
		}

		filteredRequests := loader.FilterByTags(requests, wantedTags)
		if len(filteredRequests) > 0 {
			return filteredRequests, false
		}

		logger.Warnf(`No requests found for tags "%s".`, tags)

		return requests, true
	}

	// Or use the fallback (return all requests).
	return requests, false
}

// processRequests iterates over the filtered APIRequests, executes
// each (including test cases), and writes the results. It returns
// the aggregated Result and Report.
func processRequests(ctx context.Context, filteredRequests []*loader.APIRequest) (*report.Result, *report.Report) {
	res := &report.Result{} //nolint:exhaustruct
	rep := &report.Report{} //nolint:exhaustruct

	for idx, req := range filteredRequests {
		if ctx.Err() != nil {
			logger.Debugf("Received cancellation signal. Stopping request processing.")

			return res, rep
		}

		if idx > 0 {
			logger.NewLine()
		}

		logger.Infof(`Run: %d, Test case: %d, File: "%s"`, idx+1, 0, req.JSONFilePath)

		testCases := req.TestCases

		// Execute first (main) request, regardless of whether additional test cases exist.
		exec.ProcessRequest(ctx, idx+1, req, nil, res, rep)

		// Execute additional requests depending on the number of defined test cases.
		for testCaseIndex, testCase := range testCases {
			if testCase.ParamsData == "" && testCase.PostBodyData == "" {
				continue
			}

			modifiedReq := *req

			if testCase.ParamsData != "" {
				modifiedReq.Request.Params = util.ReplaceQueryParam(req.Request.Params, testCase.ParamsData)
			}

			if testCase.PostBodyData != "" {
				modifiedReq.Request.PostBody = testCase.PostBodyData
			}

			exec.ProcessRequest(ctx, idx+1, &modifiedReq, &testCaseIndex, res, rep)
			logger.Infof("Test case: %s", testCase.Name)
		}
	}

	return res, rep
}

// notification sends a summary notification via WebEx webhook.
func notification(ctx context.Context, cfg *config.Config, conn *sqlite.Conn, res *report.Result, rep *report.Report) {
	if cfg.Notification.WebEx == nil || !cfg.Notification.WebEx.Active {
		return
	}

	const reportFile = "./logs/report.json"

	hostname, _ := os.Hostname()
	hostnameMessage := fmt.Sprintf("_Message from **%s** (hostname)_", hostname)

	if res.RequestErrorCount == 0 && res.FormatResponseErrorCount == 0 && res.ChangedFilesCount == 0 {
		_ = os.Remove(reportFile)

		isHeartbeatTime, err := report.IsHeartbeatTime(cfg)
		if err != nil {
			return
		}

		if !isHeartbeatTime {
			return
		}

		if err := report.UpdateHeartbeatTime(cfg); err != nil {
			return
		}

		mdMessage := fmt.Sprintf(
			`{"markdown":"#### 💙 %s\nHeartbeat: __still alive__\n\n%s"}`,
			config.Version,
			hostnameMessage,
		)

		webhookPayload := []byte(mdMessage)

		report.WebExWebhookNotification(ctx, conn,
			cfg.Notification.WebEx.WebhookURL,
			cfg.Notification.WebEx.Space,
			webhookPayload)

		return
	}

	if err := rep.SaveToFile(reportFile); err != nil {
		logger.Errorf("Error on save file. Error: %v", err)

		return
	}

	data, err := os.ReadFile(reportFile)
	if err != nil {
		logger.Errorf("Error on read file. Error: %v", err)

		return
	}

	mdCodeBlock := fmt.Sprintf("```json\n%s\n```", data)

	mdResult := fmt.Sprintf(
		"Changed files: __%d__\nRequest errors: __%d__\nFormat response errors: __%d__\n\n📄 _report.json_",
		res.ChangedFilesCount,
		res.RequestErrorCount,
		res.FormatResponseErrorCount,
	)

	mdMessage := fmt.Sprintf(
		"{markdown: \"#### 🔴 %s\n%s\n\\($code)\n\n%s\"}",
		config.Version,
		mdResult,
		hostnameMessage,
	)

	jqArgs := []string{
		"-nr",
		"--arg",
		"code", mdCodeBlock,
		mdMessage,
	}

	webhookPayload, err := exec.RunJQ(ctx, jqArgs, nil)
	if err != nil {
		return
	}

	report.WebExWebhookNotification(ctx, conn,
		cfg.Notification.WebEx.WebhookURL,
		cfg.Notification.WebEx.Space,
		webhookPayload)
}
