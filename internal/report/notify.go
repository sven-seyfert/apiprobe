package report

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"zombiezen.com/go/sqlite"

	"github.com/sven-seyfert/apiprobe/internal/config"
	"github.com/sven-seyfert/apiprobe/internal/crypto"
	"github.com/sven-seyfert/apiprobe/internal/db"
	"github.com/sven-seyfert/apiprobe/internal/logger"
)

// Notification sends a summary notification via WebEx webhook.
func Notification(
	ctx context.Context,
	cfg *config.Config,
	conn *sqlite.Conn,
	res *Result,
	rep *Report,
	name string,
) {
	if cfg.Notification.WebEx == nil || !cfg.Notification.WebEx.Active {
		return
	}

	const reportFile = "./logs/report.json"

	hostname, _ := os.Hostname()
	hostnameMessage := fmt.Sprintf("Message from: __%s__ (hostname)", hostname)

	if res.RequestErrorCount == 0 && res.FormatResponseErrorCount == 0 && res.ChangedFilesCount == 0 {
		_ = os.Remove(reportFile)

		isHeartbeatTime, err := IsHeartbeatTime(cfg)
		if err != nil {
			return
		}

		if !isHeartbeatTime {
			return
		}

		if err = UpdateHeartbeatTime(cfg); err != nil {
			return
		}

		mdMessage := fmt.Sprintf(
			`{"markdown":"#### 💙 %s\nHeartbeat: __still alive__\n\n%s"}`,
			config.Version,
			hostnameMessage,
		)

		webhookPayload := []byte(mdMessage)

		webExWebhookNotification(ctx, conn,
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
		"#### %s %s\n%s\n%s\n\n%s",
		trafficLight,
		config.Version,
		mdResult,
		mdCodeBlock,
		hostnameMessage,
	)

	payload := map[string]string{
		"markdown": mdMessage,
	}

	webhookPayload, _ := json.Marshal(payload)

	webExWebhookNotification(ctx, conn,
		cfg.Notification.WebEx.WebhookURL,
		cfg.Notification.WebEx.Space,
		webhookPayload)
}

// webExWebhookNotification sends the given JSON payload to the configured
// WebEx incoming webhook URL.
func webExWebhookNotification(
	ctx context.Context,
	conn *sqlite.Conn,
	webhookURL string,
	spaceSecret string,
	webhookPayload []byte,
) {
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
