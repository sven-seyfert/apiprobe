package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"zombiezen.com/go/sqlite"

	"github.com/sven-seyfert/apiprobe/internal/auth"
	"github.com/sven-seyfert/apiprobe/internal/config"
	"github.com/sven-seyfert/apiprobe/internal/crypto"
	"github.com/sven-seyfert/apiprobe/internal/db"
	"github.com/sven-seyfert/apiprobe/internal/exec"
	"github.com/sven-seyfert/apiprobe/internal/flags"
	"github.com/sven-seyfert/apiprobe/internal/loader"
	"github.com/sven-seyfert/apiprobe/internal/logger"
	"github.com/sven-seyfert/apiprobe/internal/report"
	"github.com/sven-seyfert/apiprobe/internal/util"
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
	filteredRequests := loader.ExcludeRequestsByID(requests, *cliFlags.Exclude)

	// Filter requests based on single id (ten character long hex hash) or by flags.
	filteredRequests, notFound := loader.FilterRequests(filteredRequests, *cliFlags.ID, *cliFlags.Tags)
	if notFound {
		return
	}

	// Merge possible pre-requests (prepend) with the filtered requests.
	preparedRequests, err := loader.MergePreRequests(requests, filteredRequests)
	if err != nil {
		logger.Fatalf("Program exits: Failed to gather pre-requests.")

		return
	}

	// Replace secrets placeholders in the requests with actual values.
	finalRequests, err := crypto.HandleSecrets(preparedRequests, conn)
	if err != nil {
		logger.Fatalf("Program exits: Failed to handle secrets in requests.")

		return
	}

	// Only once requests are loaded successfully, set up signal-cancellation context.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Initializes token store.
	tokenStore := auth.NewTokenStore()

	// Process each API request, optionally with test case variations.
	res, rep := processRequests(ctx, finalRequests, tokenStore)

	// Send notification on error case or on changes.
	notification(ctx, cfg, conn, res, rep, *cliFlags.Name)
}

// processRequests iterates over the APIRequests, executes
// each (including test cases), and writes the results. It returns
// the aggregated Result and Report.
func processRequests(
	ctx context.Context,
	requests []*loader.APIRequest,
	tokenStore *auth.TokenStore,
) (*report.Result, *report.Report) {
	res := &report.Result{}
	rep := &report.Report{}

	for idx, req := range requests {
		if ctx.Err() != nil {
			logger.Debugf("Received cancellation signal. Stopping request processing.")

			return res, rep
		}

		if idx > 0 {
			logger.NewLine()
		}

		logger.Infof(`Run: %d, Test case: %d, File: "%s"`, idx+1, 0, req.JSONFilePath)

		testCases := req.TestCases

		if req.PreRequestID != "" {
			repaceAuthTokenPlaceholderInRequestHeader(req, tokenStore)
		}

		// Execute first (main) request, regardless of whether additional test cases exist.
		exec.ProcessRequest(ctx, idx+1, req, nil, res, rep, tokenStore)

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

			exec.ProcessRequest(ctx, idx+1, &modifiedReq, &testCaseIndex, res, rep, tokenStore)
			logger.Infof("Test case: %s", testCase.Name)
		}
	}

	return res, rep
}

// repaceAuthTokenPlaceholderInRequestHeader replaces the <auth-token> placeholder
// in request headers with the corresponding token from the token store, if available.
// Returns nothing.
func repaceAuthTokenPlaceholderInRequestHeader(req *loader.APIRequest, tokenStore *auth.TokenStore) {
	const headerReplacementIndicator = "<auth-token>"

	lookupID := req.PreRequestID

	for idx, header := range req.Request.Headers {
		if !strings.Contains(header, headerReplacementIndicator) {
			continue
		}

		if token, found := tokenStore.Get(lookupID); found {
			lastTokenChars := token[util.Max(0, len(token)-12):] //nolint:mnd

			logger.Debugf(`Token "...%s" found for auth request "%s".`, lastTokenChars, lookupID)

			req.Request.Headers[idx] = strings.ReplaceAll(header, headerReplacementIndicator, token)

			break
		}

		logger.Warnf(`No token found for auth request "%s".`, lookupID)
	}
}

// notification sends a summary notification via WebEx webhook.
func notification(ctx context.Context, cfg *config.Config, conn *sqlite.Conn, res *report.Result, rep *report.Report, name string) {
	if cfg.Notification.WebEx == nil || !cfg.Notification.WebEx.Active {
		return
	}

	const reportFile = "./logs/report.json"

	hostname, _ := os.Hostname()
	hostnameMessage := fmt.Sprintf("Message from: __%s__ (hostname)", hostname)

	if res.RequestErrorCount == 0 && res.FormatResponseErrorCount == 0 && res.ChangedFilesCount == 0 {
		_ = os.Remove(reportFile)

		isHeartbeatTime, err := report.IsHeartbeatTime(cfg)
		if err != nil {
			return
		}

		if !isHeartbeatTime {
			return
		}

		if err = report.UpdateHeartbeatTime(cfg); err != nil {
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

	testRunName := ""
	if name != "" {
		testRunName = fmt.Sprintf("`%s`\n\n", name)
	}

	mdResult := fmt.Sprintf(
		"%sChanged files: __%d__\nRequest errors: __%d__\nFormat response errors: __%d__\n\n📄 _report.json_",
		testRunName,
		res.ChangedFilesCount,
		res.RequestErrorCount,
		res.FormatResponseErrorCount,
	)

	trafficLight := "🔴"
	if res.RequestErrorCount == 0 && res.FormatResponseErrorCount == 0 && res.ChangedFilesCount > 0 {
		trafficLight = "🟡"
	}

	mdMessage := fmt.Sprintf(
		"{markdown: \"#### %s %s\n%s\n\\($code)\n\n%s\"}",
		trafficLight,
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
