package report

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/sven-seyfert/apiprobe/internal/config"
	"github.com/sven-seyfert/apiprobe/internal/crypto"
	"github.com/sven-seyfert/apiprobe/internal/db"
	"github.com/sven-seyfert/apiprobe/internal/logger"

	"zombiezen.com/go/sqlite"
)

// Notification sends summary notifications via WebEx and MS Teams webhooks.
// It selects the notification channel and triggers the appropriate send function.
func Notification(
	ctx context.Context,
	cfg *config.Config,
	conn *sqlite.Conn,
	res *Result,
	rep *Report,
	runName string,
	notifyChannel string,
) {
	if notifyChannel == "" {
		notifyChannel = "default"
	}

	if cfg.Notification.WebEx != nil && cfg.Notification.WebEx.Active {
		sendWebExNotifications(ctx, cfg, conn, res, rep, runName, notifyChannel)
	}

	if cfg.Notification.MSTeams != nil && cfg.Notification.MSTeams.Active {
		sendMSTeamsNotifications(ctx, cfg, conn, res, rep, runName, notifyChannel)
	}
}

// buildReportFilePath generates a timestamped JSON file path for saving reports.
// The format is ./reports/YYYY-MM-DD-HH-MM-SS.mmm.json.
func buildReportFilePath() string {
	const reportsPath = "./reports"
	const ext = "json"

	now := time.Now()
	timestamp := now.Format("2006-01-02-15-04-05.000")

	return fmt.Sprintf("%s/%s.%s", reportsPath, timestamp, ext)
}

// sendNotification sends the given JSON payload to the configured incoming webhook URL.
// It handles secret replacement in the webhook URL and logs the result.
func sendNotification(
	ctx context.Context,
	conn *sqlite.Conn,
	webhookURL string,
	webhookPayload []byte,
	notificationTool string,
) {
	const secretPrefix = "<secret-"
	const secretSuffix = ">"

	url := webhookURL

	if strings.Contains(webhookURL, secretPrefix) {
		urlSecret := crypto.ExtractSecretHash(webhookURL)
		urlIdentifier, _ := db.SelectHash(conn, urlSecret)
		webhookIdentifier := crypto.Deobfuscate(urlIdentifier)
		url = strings.Replace(webhookURL, secretPrefix+urlSecret+secretSuffix, webhookIdentifier, 1)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(webhookPayload))
	if err != nil {
		logger.Errorf("Error creating new request. Error: %v", err)

		return
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Errorf("Error sending request. Error: %v", err)

		return
	}
	defer resp.Body.Close()

	logger.Infof(notificationTool+" notification sent successfully (status: %d)", resp.StatusCode)
}
