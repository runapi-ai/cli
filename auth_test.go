package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	runapi "github.com/runapi-ai/cli/internal/runapi"
	"github.com/runapi-ai/core-sdk/go/option"
)

func TestAuthStatusFromConfig(t *testing.T) {
	isolateConfig(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/me" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"id": 42, "name": "Gary", "email": "gary@example.com",
				"account": map[string]any{"id": 1, "name": "Personal"},
			})
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	if err := saveConfig(configFile{APIKey: "valid-key", BaseURL: server.URL}); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	c := newCLI()
	c.stdout = &stdout
	c.stderr = &bytes.Buffer{}

	cmd := c.command()
	cmd.SetArgs([]string{"auth", "status"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var result map[string]any
	json.Unmarshal(stdout.Bytes(), &result)
	if result["authenticated"] != true {
		t.Fatalf("expected authenticated: true, got %v", result)
	}
	if result["source"] != "config" {
		t.Fatalf("expected source 'config', got %v", result["source"])
	}
	user := result["user"].(map[string]any)
	if user["email"] != "gary@example.com" {
		t.Fatalf("expected email 'gary@example.com', got %v", user["email"])
	}
	if user["name"] != "Gary" {
		t.Fatalf("expected name 'Gary', got %v", user["name"])
	}
	// Should not include account (user-scoped)
	if _, ok := user["account"]; ok {
		t.Fatal("expected no account in user output")
	}
}

func TestAuthStatusFromEnv(t *testing.T) {
	isolateConfig(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"id": 1, "name": "Test", "email": "env@example.com",
			"account": map[string]any{"id": 1, "name": "Acme"},
		})
	}))
	defer server.Close()

	t.Setenv("RUNAPI_API_KEY", "env-key")
	t.Setenv("RUNAPI_BASE_URL", server.URL)

	var stdout bytes.Buffer
	c := newCLI()
	c.stdout = &stdout
	c.stderr = &bytes.Buffer{}

	cmd := c.command()
	cmd.SetArgs([]string{"auth", "status"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var result map[string]any
	json.Unmarshal(stdout.Bytes(), &result)
	if result["source"] != "env" {
		t.Fatalf("expected source 'env', got %v", result["source"])
	}
	if result["authenticated"] != true {
		t.Fatalf("expected authenticated: true, got %v", result)
	}
}

func TestAuthStatusNoToken(t *testing.T) {
	isolateConfig(t)
	t.Setenv("RUNAPI_API_KEY", "")
	t.Setenv("RUNAPI_BASE_URL", "")

	// Write an empty config to guarantee no leaked API key from other tests
	if err := saveConfig(configFile{}); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	c := newCLI()
	c.stdout = &stdout
	c.stderr = &bytes.Buffer{}
	c.newClient = func(opts ...option.ClientOption) (*runapi.Client, error) {
		t.Fatal("unexpected network call: newClient should not be called when no API key is set")
		return nil, nil
	}

	cmd := c.command()
	cmd.SetArgs([]string{"auth", "status"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var result map[string]any
	json.Unmarshal(stdout.Bytes(), &result)
	if result["authenticated"] != false {
		t.Fatalf("expected authenticated: false, got %v", result)
	}
	if result["source"] != "none" {
		t.Fatalf("expected source 'none', got %v", result["source"])
	}
}

func TestAuthStatusInvalidToken(t *testing.T) {
	isolateConfig(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":"unauthorized"}`))
	}))
	defer server.Close()

	if err := saveConfig(configFile{APIKey: "bad-key", BaseURL: server.URL}); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	c := newCLI()
	c.stdout = &stdout
	c.stderr = &stderr

	cmd := c.command()
	cmd.SetArgs([]string{"auth", "status"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var result map[string]any
	json.Unmarshal(stdout.Bytes(), &result)
	if result["authenticated"] != false {
		t.Fatalf("expected authenticated: false, got %v", result)
	}
	if result["source"] != "config" {
		t.Fatalf("expected source 'config', got %v", result["source"])
	}
	if result["error"] != "token expired or revoked" {
		t.Fatalf("expected error message, got %v", result["error"])
	}
}

func TestAuthImportTokenVerifiesAndWritesConfig(t *testing.T) {
	isolateConfig(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer good-key" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error":"bad token"}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"id": 1, "name": "Test", "email": "import@example.com",
			"account": map[string]any{"id": 1, "name": "Acme"},
		})
	}))
	defer server.Close()

	var stdout bytes.Buffer
	c := newCLI()
	c.stdout = &stdout
	c.stderr = &bytes.Buffer{}

	cmd := c.command()
	cmd.SetArgs([]string{"auth", "import-token", "--token", "good-key", "--base-url", server.URL})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result map[string]any
	json.Unmarshal(stdout.Bytes(), &result)
	if result["imported"] != true {
		t.Fatalf("expected imported: true, got %v", result)
	}
	if result["verified"] != true {
		t.Fatalf("expected verified: true, got %v", result)
	}
	if result["base_url"] != server.URL {
		t.Fatalf("expected base_url %q, got %v", server.URL, result["base_url"])
	}

	cfg, err := loadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.APIKey != "good-key" {
		t.Fatalf("expected APIKey saved, got %q", cfg.APIKey)
	}
	if cfg.BaseURL != server.URL {
		t.Fatalf("expected BaseURL saved, got %q", cfg.BaseURL)
	}
	if cfg.CreatedAt == "" {
		t.Fatalf("expected CreatedAt set")
	}
}

func TestAuthImportTokenRejectsInvalidToken(t *testing.T) {
	isolateConfig(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":"unauthorized"}`))
	}))
	defer server.Close()

	c := newCLI()
	c.stdout = &bytes.Buffer{}
	c.stderr = &bytes.Buffer{}

	cmd := c.command()
	cmd.SetArgs([]string{"auth", "import-token", "--token", "bad-key", "--base-url", server.URL})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error, got nil")
	}

	// Config must not have been written.
	if cfg, _ := loadConfig(); cfg.APIKey != "" {
		t.Fatalf("config should be untouched, got APIKey=%q", cfg.APIKey)
	}
}

func TestAuthImportTokenSkipVerify(t *testing.T) {
	isolateConfig(t)
	c := newCLI()
	c.stdout = &bytes.Buffer{}
	c.stderr = &bytes.Buffer{}
	c.newClient = func(opts ...option.ClientOption) (*runapi.Client, error) {
		t.Fatal("--skip-verify must not build a client")
		return nil, nil
	}

	cmd := c.command()
	cmd.SetArgs([]string{"auth", "import-token", "--token", "yolo", "--skip-verify"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cfg, _ := loadConfig()
	if cfg.APIKey != "yolo" {
		t.Fatalf("expected APIKey=yolo, got %q", cfg.APIKey)
	}
}

func TestAuthImportTokenReadsStdin(t *testing.T) {
	isolateConfig(t)
	c := newCLI()
	c.stdin = strings.NewReader("piped-key\n")
	c.stdout = &bytes.Buffer{}
	c.stderr = &bytes.Buffer{}

	cmd := c.command()
	cmd.SetArgs([]string{"auth", "import-token", "--token", "-", "--skip-verify"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cfg, _ := loadConfig()
	if cfg.APIKey != "piped-key" {
		t.Fatalf("expected APIKey=piped-key, got %q", cfg.APIKey)
	}
}

func TestAuthImportTokenFromEnv(t *testing.T) {
	isolateConfig(t)
	t.Setenv("RUNAPI_API_KEY", "env-key")
	c := newCLI()
	c.stdout = &bytes.Buffer{}
	c.stderr = &bytes.Buffer{}

	cmd := c.command()
	cmd.SetArgs([]string{"auth", "import-token", "--skip-verify"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cfg, _ := loadConfig()
	if cfg.APIKey != "env-key" {
		t.Fatalf("expected APIKey=env-key, got %q", cfg.APIKey)
	}
}

func TestAuthImportTokenRequiresToken(t *testing.T) {
	isolateConfig(t)
	t.Setenv("RUNAPI_API_KEY", "")
	// Force HOME isolation so the test does not pick up the developer's real config.
	t.Setenv("HOME", t.TempDir())
	os.Unsetenv("XDG_CONFIG_HOME")

	c := newCLI()
	c.stdout = &bytes.Buffer{}
	c.stderr = &bytes.Buffer{}

	cmd := c.command()
	cmd.SetArgs([]string{"auth", "import-token", "--skip-verify"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error when no token supplied")
	}
}
