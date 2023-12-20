package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/benjamonnguyen/opendoor-chat-frontend/config"
)

func TestCreateUserIT(t *testing.T) {
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		t.SkipNow()
	}

	params := map[string]string{
		"firstName": "Jesse",
		"lastName":  "Pinkman",
		"password":  "scienceB*tch!",
		"email":     "captaincook@bb.com",
	}

	data, _ := json.Marshal(params)
	body := bytes.NewReader(data)
	req, _ := http.NewRequest(http.MethodPost, cfg.BackendBaseUrl+"/user", body)

	resp, _ := http.DefaultClient.Do(req)
	if resp.StatusCode != 201 {
		t.Fatal(resp.Status)
	}
}
