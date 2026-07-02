package main

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestLogout(t *testing.T) {
	isolateConfig(t)

	// Write config with api_key and base_url
	if err := saveConfig(configFile{APIKey: "my-key", BaseURL: "https://staging.runapi.ai", CreatedAt: "2026-04-08T00:00:00Z"}); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	c := newCLI()
	c.stdout = &stdout

	cmd := c.command()
	cmd.SetArgs([]string{"logout"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	// Verify output
	var result map[string]any
	json.Unmarshal(stdout.Bytes(), &result)
	if result["logged_out"] != true {
		t.Fatalf("expected logged_out: true, got %v", result)
	}

	// Verify config: api_key cleared, base_url preserved
	cfg, err := loadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.APIKey != "" {
		t.Fatalf("expected api_key cleared, got %q", cfg.APIKey)
	}
	if cfg.BaseURL != "https://staging.runapi.ai" {
		t.Fatalf("expected base_url preserved, got %q", cfg.BaseURL)
	}
}

func TestLogoutPreservesBaseURL(t *testing.T) {
	isolateConfig(t)

	if err := saveConfig(configFile{APIKey: "key", BaseURL: "https://custom.runapi.ai"}); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	c := newCLI()
	c.stdout = &stdout

	cmd := c.command()
	cmd.SetArgs([]string{"logout"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	cfg, err := loadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.APIKey != "" {
		t.Fatalf("expected api_key cleared, got %q", cfg.APIKey)
	}
	if cfg.BaseURL != "https://custom.runapi.ai" {
		t.Fatalf("expected base_url 'https://custom.runapi.ai', got %q", cfg.BaseURL)
	}
}

func TestLogoutNotLoggedIn(t *testing.T) {
	isolateConfig(t)

	var stdout bytes.Buffer
	c := newCLI()
	c.stdout = &stdout

	cmd := c.command()
	cmd.SetArgs([]string{"logout"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var result map[string]any
	json.Unmarshal(stdout.Bytes(), &result)
	if result["logged_out"] != true {
		t.Fatalf("expected logged_out: true, got %v", result)
	}
}
