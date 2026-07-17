package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/runapi-ai/core-sdk/go/core"
)

func TestCallbackAPIKeyResolutionPrefersFlagOverProjectConfig(t *testing.T) {
	root := projectRootFixture(t, "callback_api_key_id = \"token_config\"\n")
	c := newCLI()
	c.projectDir = root

	selection, err := c.resolveCallbackAPIKey(t.Context(), nil, "", "", "token_flag")
	if err != nil {
		t.Fatal(err)
	}
	if selection.ID != "token_flag" || selection.Source != "flag" || selection.SaveAfterValidation {
		t.Fatalf("unexpected selection: %#v", selection)
	}
}

func TestCallbackAPIKeyResolutionFlagBypassesInvalidProjectConfig(t *testing.T) {
	root := projectRootFixture(t, "api_key = \"must-not-be-read\"\n")
	c := newCLI()
	c.projectDir = root

	selection, err := c.resolveCallbackAPIKey(t.Context(), nil, "", "", "token_flag")
	if err != nil {
		t.Fatal(err)
	}
	if selection.ID != "token_flag" || selection.Source != "flag" {
		t.Fatalf("unexpected selection: %#v", selection)
	}
}

func TestCallbackAPIKeyResolutionUsesProjectConfig(t *testing.T) {
	root := projectRootFixture(t, "callback_api_key_id = \"token_config\"\n")
	c := newCLI()
	c.projectDir = root

	selection, err := c.resolveCallbackAPIKey(t.Context(), nil, "", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if selection.ID != "token_config" || selection.Source != "config" || selection.SaveAfterValidation {
		t.Fatalf("unexpected selection: %#v", selection)
	}
}

func TestCallbackAPIKeyResolutionReturnsCandidatesForNonTTY(t *testing.T) {
	root := projectRootFixture(t, "")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"api_keys":[{"id":"token_project","name":"Project key","masked_token":"runapi_abcd••••••••1234","enabled":true}]}`))
	}))
	defer server.Close()

	c := newCLI()
	c.projectDir = root
	c.stdinTTY = func() bool { return false }
	_, err := c.resolveCallbackAPIKey(t.Context(), server.Client(), server.URL, "cli-credential", "")
	var apiErr *core.Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected structured core error, got %v", err)
	}
	details, ok := apiErr.Details.(map[string]any)
	if !ok {
		t.Fatalf("expected error details, got %#v", apiErr.Details)
	}
	keys, ok := details["api_keys"].([]callbackAPIKey)
	if !ok || len(keys) != 1 || keys[0].ID != "token_project" {
		t.Fatalf("expected candidate list, got %#v", details["api_keys"])
	}
	if details["hint"] != "runapi listen <url> --callback-api-key-id token_project" {
		t.Fatalf("unexpected hint: %#v", details["hint"])
	}
	if got := exitCode(err); got != 4 {
		t.Fatalf("expected validation exit code 4, got %d", got)
	}
}

func TestCallbackAPIKeyResolutionSelectsEnabledKeyInTTY(t *testing.T) {
	root := t.TempDir()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"api_keys":[
			{"id":"token_disabled","name":"Disabled","masked_token":"runapi_old","enabled":false},
			{"id":"token_project","name":"Project key","masked_token":"runapi_abcd••••••••1234","enabled":true}
		]}`))
	}))
	defer server.Close()

	c := newCLI()
	c.projectDir = root
	c.stdinTTY = func() bool { return true }
	c.stderrTTY = func() bool { return true }
	c.selectCallbackAPIKey = func(keys []callbackAPIKey) (callbackAPIKey, error) {
		if len(keys) != 1 || keys[0].ID != "token_project" {
			t.Fatalf("selector received unexpected keys: %#v", keys)
		}
		return keys[0], nil
	}

	selection, err := c.resolveCallbackAPIKey(t.Context(), server.Client(), server.URL, "cli-credential", "")
	if err != nil {
		t.Fatal(err)
	}
	if selection.ID != "token_project" || selection.Source != "interactive" || !selection.SaveAfterValidation {
		t.Fatalf("unexpected selection: %#v", selection)
	}
	if selection.ConfigPath != filepath.Join(root, projectConfigName) {
		t.Fatalf("expected current-directory config path, got %q", selection.ConfigPath)
	}
}

func TestCallbackAPIKeyResolutionIsNonInteractiveWhenStderrIsNotTTY(t *testing.T) {
	root := projectRootFixture(t, "")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"api_keys":[{"id":"token_project","name":"Project key","masked_token":"runapi_abcd","enabled":true}]}`))
	}))
	defer server.Close()

	c := newCLI()
	c.projectDir = root
	c.stdinTTY = func() bool { return true }
	c.stderrTTY = func() bool { return false }
	c.selectCallbackAPIKey = func([]callbackAPIKey) (callbackAPIKey, error) {
		t.Fatal("selector must not run when its output is redirected")
		return callbackAPIKey{}, nil
	}

	_, err := c.resolveCallbackAPIKey(t.Context(), server.Client(), server.URL, "cli-credential", "")
	var apiErr *core.Error
	if !errors.As(err, &apiErr) || apiErr.Code != core.ErrorCode("callback_api_key_required") {
		t.Fatalf("expected structured selection error, got %v", err)
	}
}

func projectRootFixture(t *testing.T, config string) string {
	t.Helper()
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	if config != "" {
		if err := os.WriteFile(filepath.Join(root, projectConfigName), []byte(config), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}
