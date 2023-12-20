package config_test

import (
	"testing"

	"github.com/benjamonnguyen/opendoor-chat-frontend/config"
)

func TestLoadConfig(t *testing.T) {
	// want
	const wantSalt = "salt"
	const wantBackendBaseUrl = "http://localhost:8080"

	cfg := config.LoadConfig("test.env")
	if cfg.Salt != wantSalt {
		t.Errorf("cfg.Salt: got %s, want %s\n", cfg.Salt, wantSalt)
	}
	if cfg.BackendBaseUrl != wantBackendBaseUrl {
		t.Errorf("cfg.BackendBaseUrl: got %s, want %s\n", cfg.BackendBaseUrl, wantBackendBaseUrl)
	}
}
