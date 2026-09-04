package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	runapi "github.com/runapi-ai/cli/internal/runapi"
	"github.com/runapi-ai/core-sdk/go/core"
	"github.com/runapi-ai/core-sdk/go/option"
)

func TestRootHelpIsCapabilityNeutral(t *testing.T) {
	c := newCLI()
	c.stdout = &bytes.Buffer{}
	c.stderr = &bytes.Buffer{}

	cmd := c.command()
	cmd.SetArgs([]string{"--help"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	output := c.stdout.(*bytes.Buffer).String()
	if strings.Contains(output, "Suno, Veo 3.1, and Nano Banana operations") {
		t.Fatalf("expected capability-neutral root help, got:\n%s", output)
	}
	if !strings.Contains(output, "typed SDK-backed RunAPI operations") {
		t.Fatalf("expected neutral CLI scope description, got:\n%s", output)
	}
}

func TestMidjourneyShortenPromptResumesAcceptedTask(t *testing.T) {
	isolateConfig(t)
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		switch requests {
		case 1:
			if r.Method != http.MethodPost || r.URL.Path != "/api/v1/midjourney/shorten_prompt" || r.Header.Get("Idempotency-Key") == "" {
				t.Fatalf("unexpected create request: %s %s", r.Method, r.URL.Path)
			}
			w.Header().Set("Location", "/api/v1/tasks/tsk_cli_123")
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(http.StatusAccepted)
			_, _ = w.Write([]byte(`{"id":"tsk_cli_123","status":"processing"}`))
		case 2:
			w.Header().Set("Retry-After", "0")
			_, _ = w.Write([]byte(`{"id":"tsk_cli_123","status":"processing"}`))
		case 3:
			_, _ = w.Write([]byte(`{"id":"tsk_cli_123","status":"completed","response":{"status":200,"content_type":"application/json","headers":{},"body":{"prompts":["Short prompt"]}}}`))
		default:
			t.Fatalf("unexpected request %d", requests)
		}
	}))
	defer server.Close()

	c := newCLI()
	c.stdout = &bytes.Buffer{}
	c.stderr = &bytes.Buffer{}
	if code := c.run([]string{"--api-key", "test-key", "--base-url", server.URL, "--poll-interval", "1ms", "midjourney", "shorten-prompt", "--input", `{"prompt":"A detailed prompt"}`}); code != 0 {
		t.Fatalf("expected success, got %d: %s", code, c.stderr.(*bytes.Buffer).String())
	}
	var output struct {
		Prompts []string `json:"prompts"`
	}
	if err := json.Unmarshal(c.stdout.(*bytes.Buffer).Bytes(), &output); err != nil {
		t.Fatal(err)
	}
	if len(output.Prompts) != 1 || output.Prompts[0] != "Short prompt" {
		t.Fatalf("unexpected terminal output: %s", c.stdout.(*bytes.Buffer).String())
	}
}

func TestHybridActionSpecsUseRunLifecycle(t *testing.T) {
	actions := [][2]string{
		{"fish-audio", "create-voice"},
		{"fish-audio", "text-to-speech"},
		{"midjourney", "get-seed"},
		{"midjourney", "image-to-prompt"},
		{"midjourney", "shorten-prompt"},
		{"openai-transcription", "speech-to-text"},
		{"openai-tts", "text-to-speech"},
	}
	for _, action := range actions {
		t.Run(strings.Join(action[:], "/"), func(t *testing.T) {
			spec, err := findActionSpec(action[0], action[1])
			if err != nil {
				t.Fatal(err)
			}
			if !spec.isHybrid || spec.isAsync || spec.run == nil {
				t.Fatalf("expected hybrid run lifecycle: %#v", spec)
			}
			if !strings.Contains(describeAction(spec), "resume an accepted Task") {
				t.Fatalf("unexpected action description: %q", describeAction(spec))
			}
		})
	}
}

func TestFilesUploadHelpRejectsRemovedCommand(t *testing.T) {
	for _, args := range [][]string{
		{"files", "upload", "--help"},
		{"files", "--api-key", "test-key", "upload", "--help"},
		{"files", "--api-key=test-key", "upload", "-h"},
		{"files", "--quiet", "upload", "--help"},
	} {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			c := newCLI()
			c.stdout = &bytes.Buffer{}
			c.stderr = &bytes.Buffer{}

			if code := c.run(args); code == 0 {
				t.Fatal("expected removed files upload command help to return nonzero")
			}

			output := c.stderr.(*bytes.Buffer).String()
			if !strings.Contains(output, "runapi files create") {
				t.Fatalf("expected canonical files create guidance, got:\n%s", output)
			}
		})
	}
}

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

func TestPrintErrorOmitsMissingHTTPCodeAndPreservesExplicitCode(t *testing.T) {
	for _, test := range []struct {
		name     string
		code     core.ErrorCode
		expected any
	}{
		{"missing", "", nil},
		{"explicit", "source_task_not_ready", "source_task_not_ready"},
	} {
		t.Run(test.name, func(t *testing.T) {
			cli := newCLI()
			stdout := &bytes.Buffer{}
			cli.stdout = stdout
			cli.stderr = &bytes.Buffer{}
			cli.printError(core.NewError(test.code, "wait", 409, "", nil, nil))

			var payload map[string]map[string]any
			if err := json.Unmarshal(stdout.Bytes(), &payload); err != nil {
				t.Fatal(err)
			}
			value, present := payload["error"]["code"]
			if test.expected == nil && present {
				t.Fatalf("expected code to be omitted, got %v", value)
			}
			if test.expected != nil && value != test.expected {
				t.Fatalf("unexpected code: %v", value)
			}
		})
	}
}

func TestExitCodeClassifiesHTTPErrorByStatus(t *testing.T) {
	err := core.NewError("source_task_not_ready", "wait", 409, "", nil, nil)
	if code := exitCode(err); code != 4 {
		t.Fatalf("expected conflict exit code 4, got %d", code)
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
