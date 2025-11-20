package report

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/sven-seyfert/apiprobe/internal/config"
	"github.com/sven-seyfert/apiprobe/internal/logger"

	"zombiezen.com/go/sqlite"
)

// sendWebExNotifications sends a notification to WebEx based on the result and report.
// It sends either a heartbeat or a report payload depending on the result counts.
func sendWebExNotifications(
	ctx context.Context,
	cfg *config.Config,
	conn *sqlite.Conn,
	res *Result,
	rep *Report,
	runName string,
	channelName string,
) {
	const notificationTool = "WebEx"

	webhookURL, exists := cfg.Notification.WebEx.Webhooks[channelName]
	if !exists {
		logger.Warnf(notificationTool+" webhook channel '%s' not found in config. No notification sent.", channelName)

		return
	}

	reportFilePath := buildReportFilePath()
	hostname, _ := os.Hostname()
	hostnameMessage := fmt.Sprintf("Message from: __%s__ (hostname)", hostname)

	if res.RequestErrorCount == 0 && res.FormatResponseErrorCount == 0 && res.ChangedFilesCount == 0 {
		webhookPayload := buildWebExHeartbeatPayload(reportFilePath, cfg, hostnameMessage)
		if webhookPayload != nil {
			sendNotification(ctx, conn, webhookURL, webhookPayload, notificationTool)
		}

		return
	}

	if err := rep.SaveToFile(reportFilePath); err != nil {
		logger.Errorf("Error on save file. Error: %v", err)

		return
	}

	data, err := os.ReadFile(reportFilePath)
	if err != nil {
		logger.Errorf("Error on read file. Error: %v", err)

		return
	}

	webhookPayload := buildWebExReportPayload(res, runName, reportFilePath, data, hostnameMessage)

	sendNotification(ctx, conn, webhookURL, webhookPayload, notificationTool)
}

// buildWebExHeartbeatPayload creates the payload for a heartbeat notification.
// Returns the payload as a byte slice, or nil if heartbeat should not be sent.
func buildWebExHeartbeatPayload(reportFilePath string, cfg *config.Config, hostnameMessage string) []byte {
	_ = os.Remove(reportFilePath)

	isHeartbeatTime, err := IsHeartbeatTime(cfg)
	if err != nil {
		return nil
	}

	if !isHeartbeatTime {
		return nil
	}

	if err = UpdateHeartbeatTime(cfg); err != nil {
		return nil
	}

	mdMessage := fmt.Sprintf(
		`{"markdown":"#### 💙 %s\nHeartbeat: __still alive__\n\n%s"}`,
		config.Version,
		hostnameMessage,
	)

	return []byte(mdMessage)
}

// buildWebExReportPayload creates the payload for a report notification
// including result details and the report file content. Returns the
// payload as a byte slice.
func buildWebExReportPayload(res *Result, runName string, reportFilePath string, data []byte, hostnameMessage string) []byte {
	trafficLight := "🔴"
	if res.RequestErrorCount == 0 && res.FormatResponseErrorCount == 0 && res.ChangedFilesCount > 0 {
		trafficLight = "🟡"
	}

	testRunName := ""
	if runName != "" {
		testRunName = fmt.Sprintf("`%s`\n\n", runName)
	}

	mdResult := fmt.Sprintf(
		"%sFiles with changed content: __%d__\nRequest errors: __%d__\nFormat response errors: __%d__\n\n📄 _%s_",
		testRunName,
		res.ChangedFilesCount,
		res.RequestErrorCount,
		res.FormatResponseErrorCount,
		reportFilePath,
	)

	mdCodeBlock := fmt.Sprintf("```json\n%s\n```", data)

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

	return webhookPayload
}
