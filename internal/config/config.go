package config

import (
	"encoding/json"
	"os"

	"github.com/sven-seyfert/apiprobe/internal/logger"
)

const Version = "APIProbe 📡 v0.15.0 - 2025-09-30"

type Heartbeat struct {
	IntervalInHours   int    `json:"intervalInHours"`
	LastHeartbeatTime string `json:"lastHeartbeatTime"`
}

type Notification struct {
	WebEx *struct {
		Active     bool   `json:"active"`
		WebhookURL string `json:"webhookUrl"`
		Space      string `json:"space"`
	} `json:"webEx"`
}

type Config struct {
	DebugMode    bool         `json:"debugMode"`
	Heartbeat    Heartbeat    `json:"heartbeat"`
	Notification Notification `json:"notification"`
}

// Load opens the JSON configuration file, decodes its contents into
// a Config struct and returns the loaded configuration or an error.
func Load(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		logger.Errorf(`Failure opening config file "%s". Error: %v`, filePath, err)

		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	var cfg Config

	if err = decoder.Decode(&cfg); err != nil {
		logger.Errorf(`Failure parsing config file "%s". Error: %v`, filePath, err)

		return nil, err
	}

	return &cfg, nil
}
