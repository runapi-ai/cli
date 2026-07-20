package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestListenInteractiveSelectionWritesProjectConfigAfterSessionValidation(t *testing.T) {
	isolateConfig(t)
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/v1/cli/keys":
			_, _ = w.Write([]byte(`{"api_keys":[{"id":"token_project","name":"Project key","masked_token":"runapi_abcd••••••••1234","enabled":true}]}`))
		case r.Method == http.MethodPost && r.URL.Path == "/api/v1/cli/listen_sessions":
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{
				"session_token":"session_123",
				"listen_secret":"secret_123",
				"callback_api_key":{"id":"token_project","name":"Project key","masked_token":"runapi_abcd••••••••1234"}
			}`))
		case r.Method == http.MethodGet && r.URL.Path == "/api/v1/cli/listen_events":
			_, _ = w.Write([]byte(`{"events":[]}`))
			time.AfterFunc(10*time.Millisecond, cancel)
		case r.Method == http.MethodDelete && strings.HasPrefix(r.URL.Path, "/api/v1/cli/listen_sessions/"):
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()
	if err := saveConfig(configFile{APIKey: "cli-credential", BaseURL: server.URL}); err != nil {
		t.Fatal(err)
	}

	root := projectRootFixture(t, "")
	var stderr bytes.Buffer
	c := newCLI()
	c.stdout = &bytes.Buffer{}
	c.stderr = &stderr
	c.httpClient = server.Client()
	c.projectDir = root
	c.stdinTTY = func() bool { return true }
	c.stderrTTY = func() bool { return true }
	c.selectCallbackAPIKey = func(keys []callbackAPIKey) (callbackAPIKey, error) { return keys[0], nil }
	cmd := c.command()
	cmd.SetContext(ctx)
	cmd.SetArgs([]string{"listen"})

	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	projectCfg, configPath, err := loadProjectConfig(root)
	if err != nil {
		t.Fatal(err)
	}
	if projectCfg.CallbackAPIKeyID != "token_project" {
		t.Fatalf("expected selected key in project config, got %#v", projectCfg)
	}
	for _, expected := range []string{"Project key", "token_project", "runapi_abcd", configPath, "secret_123"} {
		if !strings.Contains(stderr.String(), expected) {
			t.Fatalf("expected startup output to contain %q, got %q", expected, stderr.String())
		}
	}
}

func TestListenDoesNotWriteProjectConfigWhenSessionValidationFails(t *testing.T) {
	isolateConfig(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/api/v1/cli/keys" {
			_, _ = w.Write([]byte(`{"api_keys":[{"id":"token_project","name":"Project key","masked_token":"runapi_abcd","enabled":true}]}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"Callback API key not found","code":"not_found"}`))
	}))
	defer server.Close()
	if err := saveConfig(configFile{APIKey: "cli-credential", BaseURL: server.URL}); err != nil {
		t.Fatal(err)
	}

	root := projectRootFixture(t, "")
	c := newCLI()
	c.stdout = &bytes.Buffer{}
	c.stderr = &bytes.Buffer{}
	c.httpClient = server.Client()
	c.projectDir = root
	c.stdinTTY = func() bool { return true }
	c.stderrTTY = func() bool { return true }
	c.selectCallbackAPIKey = func(keys []callbackAPIKey) (callbackAPIKey, error) { return keys[0], nil }
	cmd := c.command()
	cmd.SetArgs([]string{"listen"})

	if err := cmd.Execute(); err == nil {
		t.Fatal("expected listen session validation failure")
	}
	if _, err := os.Stat(filepath.Join(root, projectConfigName)); !os.IsNotExist(err) {
		t.Fatalf("expected no project config after failed validation, got %v", err)
	}
}

func TestListenDestroysSessionWhenProjectConfigWriteFails(t *testing.T) {
	isolateConfig(t)
	root := projectRootFixture(t, "")
	configPath := filepath.Join(root, projectConfigName)
	destroyed := make(chan struct{}, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/v1/cli/keys":
			_, _ = w.Write([]byte(`{"api_keys":[{"id":"token_project","name":"Project key","masked_token":"runapi_abcd","enabled":true}]}`))
		case r.Method == http.MethodPost && r.URL.Path == "/api/v1/cli/listen_sessions":
			if err := os.Mkdir(configPath, 0o755); err != nil {
				t.Errorf("create blocking config directory: %v", err)
			}
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{
				"session_token":"session_123",
				"listen_secret":"secret_123",
				"callback_api_key":{"id":"token_project","name":"Project key","masked_token":"runapi_abcd"}
			}`))
		case r.Method == http.MethodDelete && r.URL.Path == "/api/v1/cli/listen_sessions/session_123":
			destroyed <- struct{}{}
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()
	if err := saveConfig(configFile{APIKey: "cli-credential", BaseURL: server.URL}); err != nil {
		t.Fatal(err)
	}

	c := newCLI()
	c.stdout = &bytes.Buffer{}
	c.stderr = &bytes.Buffer{}
	c.httpClient = server.Client()
	c.projectDir = root
	c.stdinTTY = func() bool { return true }
	c.stderrTTY = func() bool { return true }
	c.selectCallbackAPIKey = func(keys []callbackAPIKey) (callbackAPIKey, error) { return keys[0], nil }
	cmd := c.command()
	cmd.SetArgs([]string{"listen"})

	if err := cmd.Execute(); err == nil {
		t.Fatal("expected project config write failure")
	}
	select {
	case <-destroyed:
	default:
		t.Fatal("expected validated listen session to be destroyed")
	}
}

func TestListenReselectsWhenCommittedCallbackAPIKeyIsUnavailable(t *testing.T) {
	isolateConfig(t)
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	var attemptedIDs []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/v1/cli/keys":
			_, _ = w.Write([]byte(`{"api_keys":[{"id":"token_member_b","name":"Member B key","masked_token":"runapi_bbbb","enabled":true}]}`))
		case r.Method == http.MethodPost && r.URL.Path == "/api/v1/cli/listen_sessions":
			var payload map[string]string
			_ = json.NewDecoder(r.Body).Decode(&payload)
			attemptedIDs = append(attemptedIDs, payload["callback_api_key_id"])
			if payload["callback_api_key_id"] == "token_member_a" {
				w.WriteHeader(http.StatusNotFound)
				_, _ = w.Write([]byte(`{"error":"Not found","code":"not_found"}`))
				return
			}
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{
				"session_token":"session_b",
				"listen_secret":"secret_b",
				"callback_api_key":{"id":"token_member_b","name":"Member B key","masked_token":"runapi_bbbb"}
			}`))
		case r.Method == http.MethodGet && r.URL.Path == "/api/v1/cli/listen_events":
			_, _ = w.Write([]byte(`{"events":[]}`))
			time.AfterFunc(10*time.Millisecond, cancel)
		case r.Method == http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()
	if err := saveConfig(configFile{APIKey: "cli-credential", BaseURL: server.URL}); err != nil {
		t.Fatal(err)
	}

	root := projectRootFixture(t, "callback_api_key_id = \"token_member_a\"\n")
	c := newCLI()
	c.stdout = &bytes.Buffer{}
	c.stderr = &bytes.Buffer{}
	c.httpClient = server.Client()
	c.projectDir = root
	c.stdinTTY = func() bool { return true }
	c.stderrTTY = func() bool { return true }
	c.selectCallbackAPIKey = func(keys []callbackAPIKey) (callbackAPIKey, error) { return keys[0], nil }
	cmd := c.command()
	cmd.SetContext(ctx)
	cmd.SetArgs([]string{"listen"})

	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if strings.Join(attemptedIDs, ",") != "token_member_a,token_member_b" {
		t.Fatalf("expected committed key then replacement key, got %v", attemptedIDs)
	}
	projectCfg, _, err := loadProjectConfig(root)
	if err != nil {
		t.Fatal(err)
	}
	if projectCfg.CallbackAPIKeyID != "token_member_b" {
		t.Fatalf("expected replacement key persisted, got %#v", projectCfg)
	}
}

func TestListenPrintSecretUsesSelectedCallbackAPIKeyAndExits(t *testing.T) {
	isolateConfig(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/api/v1/cli/listen_secret" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if got := r.URL.Query().Get("callback_api_key_id"); got != "token_project" {
			t.Errorf("expected selected callback API key, got %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"listen_secret":"secret_123",
			"callback_api_key":{"id":"token_project","name":"Project key","masked_token":"runapi_abcd••••••••1234"}
		}`))
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
	c.projectDir = projectRootFixture(t, "")
	cmd := c.command()
	cmd.SetArgs([]string{
		"listen", "localhost:3000/webhooks/runapi",
		"--callback-api-key-id", "token_project",
		"--print-secret",
	})

	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if stdout.String() != "secret_123\n" {
		t.Fatalf("expected secret only on stdout, got %q", stdout.String())
	}
}

func TestListenRotateSecretUsesSelectedCallbackAPIKeyAndExits(t *testing.T) {
	isolateConfig(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch || r.URL.Path != "/api/v1/cli/listen_secret" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if got := r.URL.Query().Get("callback_api_key_id"); got != "token_project" {
			t.Errorf("expected selected callback API key, got %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"listen_secret":"rotated_secret_123",
			"callback_api_key":{"id":"token_project","name":"Project key","masked_token":"runapi_abcd••••••••1234"}
		}`))
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
	c.projectDir = projectRootFixture(t, "")
	cmd := c.command()
	cmd.SetArgs([]string{
		"listen",
		"--callback-api-key-id", "token_project",
		"--rotate-secret",
	})

	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if stdout.String() != "rotated_secret_123\n" {
		t.Fatalf("expected rotated secret only on stdout, got %q", stdout.String())
	}
}

func TestListenStopsWhenSecretChangesDuringSessionRecreation(t *testing.T) {
	isolateConfig(t)
	var createCount atomic.Int32
	var destroyedNewSession atomic.Bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/v1/cli/keys":
			_, _ = w.Write([]byte(`{"api_keys":[{"id":"token_project","name":"Project key","masked_token":"runapi_abcd","enabled":true}]}`))
		case r.Method == http.MethodPost && r.URL.Path == "/api/v1/cli/listen_sessions":
			count := createCount.Add(1)
			w.WriteHeader(http.StatusCreated)
			if count == 1 {
				_, _ = w.Write([]byte(`{
					"session_token":"session_old",
					"listen_secret":"secret_old",
					"callback_api_key":{"id":"token_project","name":"Project key","masked_token":"runapi_abcd"}
				}`))
				return
			}
			_, _ = w.Write([]byte(`{
				"session_token":"session_new",
				"listen_secret":"secret_new",
				"callback_api_key":{"id":"token_project","name":"Project key","masked_token":"runapi_abcd"}
			}`))
		case r.Method == http.MethodGet && r.URL.Path == "/api/v1/cli/listen_events":
			if r.Header.Get("X-Session-Token") == "session_old" {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusGone)
			_, _ = w.Write([]byte(`{"error":"unexpected continued listener","code":"callback_api_key_unusable"}`))
		case r.Method == http.MethodDelete && strings.HasSuffix(r.URL.Path, "/session_new"):
			destroyedNewSession.Store(true)
			w.WriteHeader(http.StatusNoContent)
		case r.Method == http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()
	if err := saveConfig(configFile{APIKey: "cli-credential", BaseURL: server.URL}); err != nil {
		t.Fatal(err)
	}

	c := newCLI()
	c.stdout = &bytes.Buffer{}
	c.stderr = &bytes.Buffer{}
	c.httpClient = server.Client()
	c.projectDir = projectRootFixture(t, "callback_api_key_id = \"token_project\"\n")
	cmd := c.command()
	cmd.SetArgs([]string{"listen"})

	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "listen signing secret changed") {
		t.Fatalf("expected rotated-secret error, got %v", err)
	}
	if got := createCount.Load(); got != 2 {
		t.Fatalf("expected one initial and one replacement session, got %d", got)
	}
	if !destroyedNewSession.Load() {
		t.Fatal("expected replacement session with the new secret to be destroyed")
	}
}

func TestCreateListenSessionSendsSelectedCallbackAPIKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v1/cli/listen_sessions" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		var payload map[string]string
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Errorf("decode request: %v", err)
		}
		if payload["callback_api_key_id"] != "token_project" {
			t.Errorf("expected callback API key id, got %q", payload["callback_api_key_id"])
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{
			"session_token":"session_123",
			"listen_secret":"secret_123",
			"callback_api_key":{
				"id":"token_project",
				"name":"Project key",
				"masked_token":"runapi_abcd••••••••1234"
			}
		}`))
	}))
	defer server.Close()

	result, err := createListenSession(
		context.Background(), server.Client(), server.URL, "cli-credential", "token_project",
	)
	if err != nil {
		t.Fatal(err)
	}
	if result.SessionToken != "session_123" || result.ListenSecret != "secret_123" {
		t.Fatalf("unexpected session response: %#v", result)
	}
	if result.CallbackAPIKey.ID != "token_project" || result.CallbackAPIKey.Name != "Project key" {
		t.Fatalf("unexpected callback API key: %#v", result.CallbackAPIKey)
	}
}

func TestForwardEventTreatsNon2xxAsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("boom"))
	}))
	defer server.Close()

	code, err := forwardEvent(context.Background(), server.Client(), server.URL, listenEvent{SignedBody: `{"id":"task_1"}`})
	if err == nil {
		t.Fatal("expected non-2xx response to return an error")
	}
	if code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", code)
	}
	if !strings.Contains(err.Error(), "boom") {
		t.Fatalf("expected response body in error, got %q", err.Error())
	}
}

func TestNormalizeForwardURLDefaultsLocalhostToHTTP(t *testing.T) {
	got, err := normalizeForwardURL("localhost:3000/inbound_webhooks/runapi")
	if err != nil {
		t.Fatalf("normalizeForwardURL returned error: %v", err)
	}
	if got != "http://localhost:3000/inbound_webhooks/runapi" {
		t.Fatalf("expected http localhost URL, got %q", got)
	}
}

func TestNormalizeForwardURLPreservesExplicitScheme(t *testing.T) {
	got, err := normalizeForwardURL("https://your-domain.com/webhooks/runapi")
	if err != nil {
		t.Fatalf("normalizeForwardURL returned error: %v", err)
	}
	if got != "https://your-domain.com/webhooks/runapi" {
		t.Fatalf("expected explicit scheme to be preserved, got %q", got)
	}
}

func TestAckEventPostsAckEndpoint(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/cli/listen_events/42/ack" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-API-Key"); got != "key_123" {
			t.Errorf("expected X-API-Key, got %q", got)
		}
		if got := r.Header.Get("X-Session-Token"); got != "session_123" {
			t.Errorf("expected X-Session-Token, got %q", got)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	if err := ackEvent(context.Background(), server.Client(), server.URL, "key_123", "session_123", 42); err != nil {
		t.Fatalf("ackEvent returned error: %v", err)
	}
}

func TestAckEventTreats404AsSessionExpired(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	err := ackEvent(context.Background(), server.Client(), server.URL, "key_123", "session_123", 42)
	if err != errSessionExpired {
		t.Fatalf("expected errSessionExpired, got %v", err)
	}
}

func TestPollEventsTreats410AsUnusableCallbackAPIKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusGone)
		_, _ = w.Write([]byte(`{"error":"Callback API key is no longer usable","code":"callback_api_key_unusable"}`))
	}))
	defer server.Close()

	_, err := pollEvents(context.Background(), server.Client(), server.URL, "key_123", "session_123", 0)
	if !errors.Is(err, errCallbackAPIKeyUnusable) {
		t.Fatalf("expected errCallbackAPIKeyUnusable, got %v", err)
	}
}

func TestAckEventTreats410AsUnusableCallbackAPIKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusGone)
		_, _ = w.Write([]byte(`{"error":"Callback API key is no longer usable","code":"callback_api_key_unusable"}`))
	}))
	defer server.Close()

	err := ackEvent(context.Background(), server.Client(), server.URL, "key_123", "session_123", 42)
	if !errors.Is(err, errCallbackAPIKeyUnusable) {
		t.Fatalf("expected errCallbackAPIKeyUnusable, got %v", err)
	}
}

func TestListenStopsWhenCredentialAccessIsRevoked(t *testing.T) {
	tests := []struct {
		name        string
		status      int
		revokeOnAck bool
	}{
		{name: "poll unauthorized", status: http.StatusUnauthorized},
		{name: "poll forbidden", status: http.StatusForbidden},
		{name: "ack unauthorized", status: http.StatusUnauthorized, revokeOnAck: true},
		{name: "ack forbidden", status: http.StatusForbidden, revokeOnAck: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isolateConfig(t)
			var pollCount atomic.Int32
			var ackCount atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				switch {
				case r.Method == http.MethodPost && r.URL.Path == "/api/v1/cli/listen_sessions":
					w.WriteHeader(http.StatusCreated)
					_, _ = w.Write([]byte(`{
						"session_token":"session_123",
						"listen_secret":"secret_123",
						"callback_api_key":{"id":"token_project","name":"Project key","masked_token":"runapi_abcd"}
					}`))
				case r.Method == http.MethodGet && r.URL.Path == "/api/v1/cli/listen_events":
					pollCount.Add(1)
					if tt.revokeOnAck {
						_, _ = w.Write([]byte(`{"events":[{"id":42,"signed_body":"{\"id\":\"task_123\",\"status\":\"completed\"}","headers":{}}]}`))
						return
					}
					w.WriteHeader(tt.status)
					_, _ = w.Write([]byte(`{"error":"Listener access revoked","code":"cli_listen_required"}`))
				case r.Method == http.MethodPost && r.URL.Path == "/api/v1/cli/listen_events/42/ack":
					ackCount.Add(1)
					w.WriteHeader(tt.status)
					_, _ = w.Write([]byte(`{"error":"Listener access revoked","code":"cli_listen_required"}`))
				case r.Method == http.MethodDelete && r.URL.Path == "/api/v1/cli/listen_sessions/session_123":
					w.WriteHeader(http.StatusNoContent)
				default:
					t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
					w.WriteHeader(http.StatusNotFound)
				}
			}))
			defer server.Close()
			if err := saveConfig(configFile{APIKey: "cli-credential", BaseURL: server.URL}); err != nil {
				t.Fatal(err)
			}

			ctx, cancel := context.WithTimeout(t.Context(), 200*time.Millisecond)
			defer cancel()
			c := newCLI()
			c.stdout = &bytes.Buffer{}
			c.stderr = &bytes.Buffer{}
			c.httpClient = server.Client()
			c.projectDir = projectRootFixture(t, "callback_api_key_id = \"token_project\"\n")
			cmd := c.command()
			cmd.SetContext(ctx)
			cmd.SetArgs([]string{"listen"})

			if err := cmd.Execute(); err == nil {
				t.Fatal("expected revoked listener access to stop the command")
			}
			if got := pollCount.Load(); got != 1 {
				t.Fatalf("expected one poll before exit, got %d", got)
			}
			wantAckCount := int32(0)
			if tt.revokeOnAck {
				wantAckCount = 1
			}
			if got := ackCount.Load(); got != wantAckCount {
				t.Fatalf("expected %d ack requests before exit, got %d", wantAckCount, got)
			}
		})
	}
}
