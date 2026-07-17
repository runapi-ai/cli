package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAPIKeysListJSON(t *testing.T) {
	isolateConfig(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/api/v1/cli/keys" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if got := r.Header.Get("X-API-Key"); got != "cli-credential" {
			t.Errorf("expected CLI credential header, got %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"api_keys":[{"id":"token_project","name":"Project key","masked_token":"runapi_abcd••••••••1234","enabled":true}]}`))
	}))
	defer server.Close()
	if err := saveConfig(configFile{APIKey: "cli-credential", BaseURL: server.URL}); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	c := newCLI()
	c.stdout = &stdout
	c.stderr = &bytes.Buffer{}
	c.httpClient = server.Client()
	cmd := c.command()
	cmd.SetArgs([]string{"api-keys", "list", "--json"})

	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var result apiKeysResponse
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if len(result.APIKeys) != 1 {
		t.Fatalf("expected one API key, got %d", len(result.APIKeys))
	}
	key := result.APIKeys[0]
	if key.ID != "token_project" || key.Name != "Project key" || !key.Enabled {
		t.Fatalf("unexpected API key: %#v", key)
	}
}

func TestListenExplainsHowToReplaceAnImportedAPIKey(t *testing.T) {
	isolateConfig(t)
	const message = "`runapi listen` requires the CLI credential created by `runapi login` so it can list your eligible API keys and let you select one callback source. Your current API key can keep using its existing API access, but it cannot access listener operations. Remove any `--api-key` or `RUNAPI_API_KEY` override, run `runapi login`, then retry `runapi listen`."
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/api/v1/cli/keys" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": message,
			"code":  "cli_listen_required",
		})
	}))
	defer server.Close()
	if err := saveConfig(configFile{APIKey: "imported-api-key", BaseURL: server.URL}); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	c := newCLI()
	c.stdout = &stdout
	c.stderr = &stderr
	c.httpClient = server.Client()
	c.projectDir = projectRootFixture(t, "")

	if code := c.run([]string{"listen", "localhost:3000/webhooks/runapi"}); code == 0 {
		t.Fatal("expected listen to reject an imported API key")
	}

	var result struct {
		Error struct {
			Message string `json:"message"`
			Code    string `json:"code"`
		} `json:"error"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode JSON error: %v\nstdout=%s", err, stdout.String())
	}
	if result.Error.Code != "cli_listen_required" || result.Error.Message != message {
		t.Fatalf("unexpected JSON error: %#v", result.Error)
	}
	for _, expected := range []string{"runapi login", "list your eligible API keys", "RUNAPI_API_KEY"} {
		if !strings.Contains(stderr.String(), expected) {
			t.Fatalf("expected stderr to contain %q, got %q", expected, stderr.String())
		}
	}
}
