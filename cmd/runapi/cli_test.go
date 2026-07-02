package main

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	runapi "github.com/runapi-ai/cli/internal/runapi"
	"github.com/runapi-ai/core-sdk/go/core"
	"github.com/runapi-ai/core-sdk/go/option"
)

// isolateConfig points RUNAPI_CONFIG_DIR at a fresh temp dir
// so configFilePath() returns a path inside it, fully isolated
// from the real user config regardless of HOME or XDG_CONFIG_HOME.
func isolateConfig(t *testing.T) {
	t.Helper()
	t.Setenv("RUNAPI_CONFIG_DIR", t.TempDir())
}

func TestReadInputRejectsMutuallyExclusiveFlags(t *testing.T) {
	cli := newCLI()
	_, err := cli.readInput(`{"a":1}`, "payload.json")
	if err == nil {
		t.Fatal("expected validation error")
	}
	if !core.IsValidation(err) {
		t.Fatal("expected validation error code")
	}
}

func TestLoadConfigFromUserConfigDir(t *testing.T) {
	isolateConfig(t)
	path, err := configFilePath()
	if err != nil {
		t.Fatal(err)
	}
	configDir := filepath.Dir(path)
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(`{"api_key":"from-config","base_url":"https://runapi.ai"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := loadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.APIKey != "from-config" || cfg.BaseURL != "https://runapi.ai" {
		t.Fatalf("unexpected config: %#v", cfg)
	}
}

func TestRunAPIBaseURLEnvOverridesConfig(t *testing.T) {
	isolateConfig(t)
	t.Setenv("RUNAPI_BASE_URL", "http://env.runapi.ai")
	path, err := configFilePath()
	if err != nil {
		t.Fatal(err)
	}
	configDir := filepath.Dir(path)
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(`{"api_key":"from-config","base_url":"http://config.example"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	cli := newCLI()
	cli.newClient = func(args ...option.ClientOption) (*runapi.Client, error) {
		resolved, err := option.ResolveClientOptions(args...)
		if err != nil {
			t.Fatal(err)
		}
		if resolved.BaseURL != "http://env.runapi.ai" {
			t.Fatalf("expected env base url, got %q", resolved.BaseURL)
		}
		return &runapi.Client{}, nil
	}

	cmd := cli.command()
	cmd.SetArgs([]string{"version"})
	cmd.SetContext(context.Background())
	_, _, _, cancel, err := cli.clientForCommand(cmd)
	if cancel != nil {
		cancel()
	}
	if err != nil {
		t.Fatal(err)
	}
}

func TestRunVersionOutputsJSON(t *testing.T) {
	cli := newCLI()
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cli.stdout = stdout
	cli.stderr = stderr
	code := cli.run([]string{"version"})
	if code != 0 {
		t.Fatalf("expected zero exit code, got %d", code)
	}
	var payload map[string]string
	if err := json.Unmarshal(stdout.Bytes(), &payload); err != nil {
		t.Fatalf("stdout is not json: %v", err)
	}
	if payload["version"] == "" {
		t.Fatalf("missing version payload: %s", stdout.String())
	}
}

func TestAccountBalanceUsesConfigAPIKey(t *testing.T) {
	isolateConfig(t)
	path, err := configFilePath()
	if err != nil {
		t.Fatal(err)
	}
	configDir := filepath.Dir(path)
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(`{"api_key":"from-config"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	cli := newCLI()
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cli.stdout = stdout
	cli.stderr = stderr
	cli.newClient = func(args ...option.ClientOption) (*runapi.Client, error) {
		resolved, err := option.ResolveClientOptions(args...)
		if err != nil {
			t.Fatal(err)
		}
		found := resolved.APIKey == "from-config"
		if !found {
			t.Fatal("expected config api key option")
		}
		return &runapi.Client{}, nil
	}
	cmd := cli.command()
	cmd.SetArgs([]string{"version"})
	cmd.SetContext(context.Background())
	_, _, _, cancel, err := cli.clientForCommand(cmd)
	if cancel != nil {
		cancel()
	}
	if err != nil {
		t.Fatal(err)
	}
}
