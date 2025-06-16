package config

import (
	"encoding/json"
	"os"

	"github.com/sven-seyfert/apiprobe/internal/logger"
)

const Version = "APIProbe 📡 v0.3.0 - 2025-06-13"

type Notification struct {
	WebEx *struct {
		Active      bool   `json:"active"`
		WebhookURL  string `json:"webhookUrl"`
		SpaceSecret string `json:"spaceSecret"`
	} `json:"webEx"`
}

type Config struct {
	Notification Notification `json:"notification"`
}

func Load(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		logger.Errorf(`Failure opening config file "%s". Error: %v`, filePath, err)

		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	var cfg Config

	if err := decoder.Decode(&cfg); err != nil {
		logger.Errorf(`Failure parsing config file "%s". Error: %v`, filePath, err)

		return nil, err
	}

	return &cfg, nil
}
