package report

import (
	"context"
	"fmt"
	"os"

	"github.com/sven-seyfert/apiprobe/internal/config"
	"github.com/sven-seyfert/apiprobe/internal/logger"

	"zombiezen.com/go/sqlite"
)

// sendMSTeamsNotifications sends a notification to MS Teams based on the result and report.
// It sends either a heartbeat or a report payload depending on the result counts.
func sendMSTeamsNotifications(
	ctx context.Context,
	cfg *config.Config,
	conn *sqlite.Conn,
	res *Result,
	rep *Report,
	runName string,
	channelName string,
) {
	const notificationTool = "MS Teams"

	webhookURL, exists := cfg.Notification.MSTeams.Webhooks[channelName]
	if !exists {
		logger.Warnf(notificationTool+" webhook channel '%s' not found in config. No notification sent.", channelName)

		return
	}

	reportFilePath := buildReportFilePath()
	hostname, _ := os.Hostname()
	hostnameMessage := fmt.Sprintf("Message from: __%s__ (hostname)", hostname)

	if res.RequestErrorCount == 0 && res.FormatResponseErrorCount == 0 && res.ChangedFilesCount == 0 {
		webhookPayload := buildMSTeamsHeartbeatPayload(reportFilePath, cfg, hostnameMessage)
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

	webhookPayload := buildMSTeamsReportPayload(res, runName, reportFilePath, data, hostnameMessage)

	sendNotification(ctx, conn, webhookURL, webhookPayload, notificationTool)
}

func buildMSTeamsHeartbeatPayload(reportFilePath string, cfg *config.Config, hostnameMessage string) []byte {
	// TODO
	fmt.Println(reportFilePath)
	fmt.Println(cfg)
	fmt.Println(hostnameMessage)

	return nil
}

func buildMSTeamsReportPayload(res *Result, runName string, reportFilePath string, data []byte, hostnameMessage string) []byte {
	// TODO
	fmt.Println(res)
	fmt.Println(runName)
	fmt.Println(reportFilePath)
	fmt.Println(data)
	fmt.Println(hostnameMessage)

	return nil
}
