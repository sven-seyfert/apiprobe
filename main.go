package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/sven-seyfert/apiprobe/internal/auth"
	"github.com/sven-seyfert/apiprobe/internal/config"
	"github.com/sven-seyfert/apiprobe/internal/crypto"
	"github.com/sven-seyfert/apiprobe/internal/db"
	"github.com/sven-seyfert/apiprobe/internal/exec"
	"github.com/sven-seyfert/apiprobe/internal/flags"
	"github.com/sven-seyfert/apiprobe/internal/loader"
	"github.com/sven-seyfert/apiprobe/internal/logger"
	"github.com/sven-seyfert/apiprobe/internal/report"

	"zombiezen.com/go/sqlite"
)

// main initializes the logger and database, parses command-line flags loads
// configuration and seeds the database. It then loads and filters API request
// definitions, injects secrets, establishes a cancellation-aware context,
// processes each request and finally sends notifications based on errors
// or detected changes.
func main() {
	dbConn, cliFlags, err := initializeServices()
	if err != nil {
		logger.Fatalf("Program exits: %v", err)

		return
	}
	defer dbConn.Close()

	cfg, err := config.Load("./config/apiprobe.json")
	if err != nil {
		logger.Fatalf("Program exits: Failed to load config file.")

		return
	}

	// Handle command-line flags.
	complete, err := flags.IsNewID(*cliFlags.NewID)
	if complete || err != nil {
		return
	}

	complete, err = flags.IsNewFile(*cliFlags.NewFile)
	if complete || err != nil {
		return
	}

	complete, err = flags.IsAddSecret(*cliFlags.AddSecret, dbConn)
	if complete || err != nil {
		return
	}

	// Fill database with default seed data.
	err = db.InsertSeedData(dbConn)
	if err != nil {
		logger.Fatalf("Program exits: Failed to fill database with seed default data.")

		return
	}

	// Load requests from JSON files in the input directory.
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

	// Prepare the requests by compacting the JSON POST body,
	// handling "x-www-form-urlencoded" and POST body test cases.
	for _, req := range preparedRequests {
		if err = req.PreparePostBody(); err != nil {
			logger.Fatalf("Program exits: Failed to prepare the POST body.")
		}

		if err = req.PreparePostBodyData(); err != nil {
			logger.Fatalf("Program exits: Failed to prepare the POST body test cases.")
		}
	}

	// Replace secrets placeholders in the requests with actual values.
	finalRequests, err := crypto.HandleSecrets(preparedRequests, dbConn)
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
	res, rep := processRequests(ctx, finalRequests, tokenStore, cfg.DebugMode)

	// Send notification on error case or on changes.
	report.Notification(ctx, cfg, dbConn, res, rep, *cliFlags.Name)
}

// initializeServices initializes logger, database and CLI flags.
// Returns database connection, CLI flags and error if initialization fails.
func initializeServices() (*sqlite.Conn, *flags.CLIFlags, error) {
	if err := logger.Init(); err != nil {
		return nil, nil, errors.Join(errors.New("failed to initialize logger: "), err)
	}

	conn, err := db.Init()
	if err != nil {
		return nil, nil, errors.Join(errors.New("failed to initialize database: "), err)
	}

	cliFlags := flags.Init()

	return conn, cliFlags, nil
}

// processRequests iterates over the APIRequests, executes
// each (including test cases), and writes the results. It returns
// the aggregated Result and Report.
func processRequests(
	ctx context.Context,
	requests []*loader.APIRequest,
	tokenStore *auth.TokenStore,
	debugMode bool,
) (*report.Result, *report.Report) {
	res := &report.Result{}
	rep := &report.Report{}

	for idx, req := range requests {
		if ctx.Err() != nil {
			logger.Debugf("Received cancellation signal. Stopping request processing.")

			return res, rep
		}

		if !req.IsActive {
			continue
		}

		if idx > 0 {
			logger.NewLine()
		}

		logger.Infof(`Run: %d, Test case: %d, File: "%s"`, idx+1, 0, req.JSONFilePath)

		if req.PreRequestID != "" {
			auth.RepaceAuthTokenPlaceholderInRequestHeader(req, tokenStore)
		}

		// Execute first (main) request, regardless of whether additional test cases exist.
		exec.ProcessFirstRequest(ctx, idx+1, req, nil, res, rep, tokenStore, debugMode)

		// Execute additional requests of the same JSON definition file,
		// depending on the number of defined test cases.
		exec.ProcessTestCasesRequests(ctx, req, idx, res, rep, tokenStore, debugMode)
	}

	return res, rep
}
