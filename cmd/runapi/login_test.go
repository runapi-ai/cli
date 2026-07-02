package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestGenerateState(t *testing.T) {
	state, err := generateState()
	if err != nil {
		t.Fatal(err)
	}
	if len(state) != 64 {
		t.Fatalf("expected 64 hex chars, got %d", len(state))
	}
	if _, err := hex.DecodeString(state); err != nil {
		t.Fatalf("expected valid hex, got error: %v", err)
	}
}

func TestGeneratePKCE(t *testing.T) {
	verifier, challenge, err := generatePKCE()
	if err != nil {
		t.Fatal(err)
	}
	if len(verifier) == 0 {
		t.Fatal("verifier is empty")
	}
	if len(challenge) == 0 {
		t.Fatal("challenge is empty")
	}

	// Verify SHA256 relationship
	h := sha256.Sum256([]byte(verifier))
	expected := base64.RawURLEncoding.EncodeToString(h[:])
	if challenge != expected {
		t.Fatalf("challenge mismatch: got %q, want %q", challenge, expected)
	}
}

func TestSaveConfig(t *testing.T) {
	isolateConfig(t)

	cfg := configFile{
		APIKey:  "test-key",
		BaseURL: "https://staging.runapi.ai",
	}
	if err := saveConfig(cfg); err != nil {
		t.Fatal(err)
	}

	path, _ := configFilePath()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	var loaded configFile
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatal(err)
	}
	if loaded.APIKey != "test-key" {
		t.Fatalf("expected api_key 'test-key', got %q", loaded.APIKey)
	}
	if loaded.BaseURL != "https://staging.runapi.ai" {
		t.Fatalf("expected base_url preserved, got %q", loaded.BaseURL)
	}

	// Check file permissions
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("expected file permissions 0600, got %o", info.Mode().Perm())
	}
}

func TestSaveConfigPreservesBaseURL(t *testing.T) {
	isolateConfig(t)

	// Write initial config with custom base_url
	initial := configFile{APIKey: "old-key", BaseURL: "https://staging.runapi.ai"}
	if err := saveConfig(initial); err != nil {
		t.Fatal(err)
	}

	// Load, update api_key only, save
	loaded, err := loadConfig()
	if err != nil {
		t.Fatal(err)
	}
	loaded.APIKey = "new-key"
	if err := saveConfig(loaded); err != nil {
		t.Fatal(err)
	}

	// Verify base_url preserved
	final, err := loadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if final.APIKey != "new-key" {
		t.Fatalf("expected api_key 'new-key', got %q", final.APIKey)
	}
	if final.BaseURL != "https://staging.runapi.ai" {
		t.Fatalf("expected base_url preserved, got %q", final.BaseURL)
	}
}

func TestLoginFlowEndToEnd(t *testing.T) {
	isolateConfig(t)

	// Mock exchange server
	exchangeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v1/cli/exchange" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{
			"key":        "test-token-abc123",
			"token_type": "standard",
			"name":       "RunAPI CLI on test",
			"user":       map[string]string{"email": "developer@runapi.ai"},
		})
	}))
	defer exchangeServer.Close()

	// Start the local callback server
	listener, err := startListener()
	if err != nil {
		t.Fatal(err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close() // Free port for reuse in test

	verifier, _, err := generatePKCE()
	if err != nil {
		t.Fatal(err)
	}

	resp, err := exchangeCode(
		t.Context(),
		exchangeServer.URL+"/api/v1/cli/exchange",
		"test-code", verifier, port, "test-host",
	)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Key != "test-token-abc123" {
		t.Fatalf("expected token 'test-token-abc123', got %q", resp.Key)
	}
	if resp.User.Email != "developer@runapi.ai" {
		t.Fatalf("expected email 'developer@runapi.ai', got %q", resp.User.Email)
	}
}

func TestLoginCancelCallback(t *testing.T) {
	listener, err := startListener()
	if err != nil {
		t.Fatal(err)
	}
	port := listener.Addr().(*net.TCPAddr).Port

	state, _ := generateState()

	// Launch waitForCallback in a goroutine
	ctx := t.Context()
	errCh := make(chan error, 1)
	go func() {
		_, err := waitForCallback(ctx, listener, state, "http://localhost:3000")
		errCh <- err
	}()

	// Simulate cancel callback from browser
	cancelURL := fmt.Sprintf("http://127.0.0.1:%d/callback?error=access_denied&state=%s", port, state)
	resp, err := http.Get(cancelURL)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	// Should get an error back
	cbErr := <-errCh
	if cbErr == nil {
		t.Fatal("expected error from cancel callback")
	}
	if cbErr.Error() != "Authorization cancelled by user" {
		t.Fatalf("unexpected error: %v", cbErr)
	}
}

func TestLoginStateMismatch(t *testing.T) {
	listener, err := startListener()
	if err != nil {
		t.Fatal(err)
	}
	port := listener.Addr().(*net.TCPAddr).Port

	ctx := t.Context()
	errCh := make(chan error, 1)
	go func() {
		_, err := waitForCallback(ctx, listener, "expected-state", "http://localhost:3000")
		errCh <- err
	}()

	// Send callback with wrong state
	url := fmt.Sprintf("http://127.0.0.1:%d/callback?code=abc&state=wrong-state", port)
	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	cbErr := <-errCh
	if cbErr == nil {
		t.Fatal("expected error from state mismatch")
	}
}

func TestLoginTimeout(t *testing.T) {
	listener, err := startListener()
	if err != nil {
		t.Fatal(err)
	}

	state, _ := generateState()

	// Use a very short timeout context
	ctx, cancel := context.WithTimeout(t.Context(), 50*time.Millisecond)
	defer cancel()

	_, err = waitForCallback(ctx, listener, state, "http://localhost:3000")
	if err == nil {
		t.Fatal("expected timeout error")
	}
	if !strings.Contains(err.Error(), "timed out") {
		t.Fatalf("expected timeout error message, got: %v", err)
	}

	// Verify exit code mapping
	code := exitCode(err)
	if code != 5 {
		t.Fatalf("expected exit code 5 for timeout, got %d", code)
	}
}
