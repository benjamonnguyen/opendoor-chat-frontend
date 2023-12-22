package config_test

import (
	"testing"

	"github.com/benjamonnguyen/opendoor-chat-frontend/config"
)

func TestLoadConfig(t *testing.T) {
	// want
	const wantBackendApiKey = "backendApiKey"
	const wantBackendBaseUrl = "http://localhost:8080"

	cfg, _ := config.LoadConfig("test.env")
	if cfg.BackendApiKey != wantBackendApiKey {
		t.Errorf("cfg.Salt: got %s, want %s\n", cfg.BackendApiKey, wantBackendApiKey)
	}
	if cfg.BackendBaseUrl != wantBackendBaseUrl {
		t.Errorf("cfg.BackendBaseUrl: got %s, want %s\n", cfg.BackendBaseUrl, wantBackendBaseUrl)
	}
}
