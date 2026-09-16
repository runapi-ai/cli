package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	runapi "github.com/runapi-ai/cli/internal/runapi"
)

func TestRunShowsUpdateNoticeOnInteractiveStderrOncePerInterval(t *testing.T) {
	isolateConfig(t)
	previousVersion := runapi.Version
	runapi.Version = "1.2.3"
	t.Cleanup(func() { runapi.Version = previousVersion })

	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if r.URL.Path != "/cli/latest.json" {
			t.Fatalf("unexpected update path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"version":"1.3.0"}`))
	}))
	defer server.Close()

	run := func() (string, string) {
		var stdout, stderr bytes.Buffer
		c := newCLI()
		c.stdout = &stdout
		c.stderr = &stderr
		c.stderrTTY = func() bool { return true }
		c.updateBaseURL = server.URL
		c.httpClient = server.Client()
		if code := c.run([]string{"agent", "list-targets"}); code != 0 {
			t.Fatalf("expected success, got %d: %s", code, stderr.String())
		}
		if !json.Valid(stdout.Bytes()) {
			t.Fatalf("stdout is not JSON: %s", stdout.String())
		}
		return stdout.String(), stderr.String()
	}

	stdout, stderr := run()
	if strings.Contains(stdout, "1.3.0") {
		t.Fatalf("update notice must not change stdout: %s", stdout)
	}
	if !strings.Contains(stderr, "RunAPI CLI 1.3.0 is available (current: 1.2.3).") || !strings.Contains(stderr, "Update with:") {
		t.Fatalf("missing update notice on stderr: %s", stderr)
	}

	_, secondStderr := run()
	if secondStderr != "" {
		t.Fatalf("expected cached update check to stay quiet, got: %s", secondStderr)
	}
	if requests != 1 {
		t.Fatalf("expected one update request, got %d", requests)
	}
}

func TestUpdateCommandDetectsInstallPath(t *testing.T) {
	for _, test := range []struct {
		name       string
		executable string
		want       string
	}{
		{name: "homebrew macos", executable: "/opt/homebrew/Cellar/runapi/0.14.0/bin/runapi", want: brewUpdateCommand},
		{name: "homebrew linux", executable: "/home/linuxbrew/.linuxbrew/Cellar/runapi/0.14.0/bin/runapi", want: brewUpdateCommand},
		{name: "go default", executable: "/home/test/go/bin/runapi", want: goUpdateCommand},
		{name: "script install", executable: "/usr/local/bin/runapi", want: shUpdateCommand + " -s -- --dir '/usr/local/bin'"},
		{name: "script install custom directory", executable: "/home/test/tools/bin/runapi", want: shUpdateCommand + " -s -- --dir '/home/test/tools/bin'"},
		{name: "script install quoted directory", executable: "/home/test/tools' bin/runapi", want: shUpdateCommand + " -s -- --dir '/home/test/tools'\\'' bin'"},
		{name: "windows go install", executable: `C:\Users\test\go\bin\runapi.exe`, want: goUpdateCommand},
		{name: "windows precompiled", executable: `C:\Tools\runapi.exe`, want: "Download the matching Windows archive from https://github.com/runapi-ai/cli/releases/latest"},
	} {
		t.Run(test.name, func(t *testing.T) {
			if got := updateCommandForPath(test.executable); got != test.want {
				t.Fatalf("updateCommandForPath(%q) = %q, want %q", test.executable, got, test.want)
			}
		})
	}
}

func TestUpdateCommandUsesResolvedSymlinkDirectory(t *testing.T) {
	targetDir := filepath.Join(t.TempDir(), "target bin")
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(targetDir, "runapi")
	if err := os.WriteFile(target, nil, 0o755); err != nil {
		t.Fatal(err)
	}
	linkDir := filepath.Join(t.TempDir(), "entry-bin")
	if err := os.MkdirAll(linkDir, 0o755); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(linkDir, "runapi")
	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}

	resolvedTargetDir, err := filepath.EvalSymlinks(targetDir)
	if err != nil {
		t.Fatal(err)
	}
	want := shUpdateCommand + " -s -- --dir '" + resolvedTargetDir + "'"
	if got := updateCommandForPath(link); got != want {
		t.Fatalf("updateCommandForPath(%q) = %q, want %q", link, got, want)
	}
}

func TestUpdateCommandDetectsConfiguredGoInstallPath(t *testing.T) {
	t.Setenv("GOBIN", "/tmp/runapi-bin")
	if got := updateCommandForPath("/tmp/runapi-bin/runapi"); got != goUpdateCommand {
		t.Fatalf("expected GOBIN binary to use go install, got %q", got)
	}

	t.Setenv("GOBIN", "")
	t.Setenv("GOPATH", "/tmp/go-work:/tmp/second-go-work")
	if got := updateCommandForPath("/tmp/second-go-work/bin/runapi"); got != goUpdateCommand {
		t.Fatalf("expected GOPATH binary to use go install, got %q", got)
	}
}

func TestRunDoesNotCheckForUpdatesOutsideInteractiveCommands(t *testing.T) {
	previousVersion := runapi.Version
	runapi.Version = "1.2.3"
	t.Cleanup(func() { runapi.Version = previousVersion })

	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		_, _ = w.Write([]byte(`{"version":"1.3.0"}`))
	}))
	defer server.Close()

	tests := []struct {
		name      string
		args      []string
		stderrTTY bool
	}{
		{name: "non-interactive", args: []string{"agent", "list-targets"}},
		{name: "quiet", args: []string{"--quiet", "agent", "list-targets"}, stderrTTY: true},
		{name: "version", args: []string{"version"}, stderrTTY: true},
		{name: "help", args: []string{"--help"}, stderrTTY: true},
		{name: "nested help", args: []string{"agent", "list-targets", "--help"}, stderrTTY: true},
		{name: "standalone root help", args: []string{"help"}, stderrTTY: true},
		{name: "standalone nested help", args: []string{"help", "agent", "list-targets"}, stderrTTY: true},
		{name: "root displays help", args: nil, stderrTTY: true},
		{name: "command group displays help", args: []string{"agent"}, stderrTTY: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			isolateConfig(t)
			statePath, err := updateStatePath()
			if err != nil {
				t.Fatal(err)
			}
			requestsBefore := requests
			c := newCLI()
			c.stdout = &bytes.Buffer{}
			c.stderr = &bytes.Buffer{}
			c.stderrTTY = func() bool { return test.stderrTTY }
			c.updateBaseURL = server.URL
			c.httpClient = server.Client()
			if code := c.run(test.args); code != 0 {
				t.Fatalf("expected success, got %d", code)
			}
			if requests != requestsBefore {
				t.Fatalf("expected no update request, got %d", requests-requestsBefore)
			}
			if _, err := os.Stat(statePath); !os.IsNotExist(err) {
				t.Fatalf("expected no update cache, stat error = %v", err)
			}
		})
	}
	if requests != 0 {
		t.Fatalf("expected no update requests, got %d", requests)
	}
}

func TestNewerReleaseVersion(t *testing.T) {
	tests := []struct {
		candidate string
		current   string
		newer     bool
	}{
		{candidate: "2.0.0", current: "1.9.9", newer: true},
		{candidate: "1.3.0", current: "1.2.9", newer: true},
		{candidate: "1.2.4", current: "1.2.3", newer: true},
		{candidate: "1.2.3", current: "1.2.3"},
		{candidate: "1.2.2", current: "1.2.3"},
		{candidate: "1.3.0-beta.1", current: "1.2.3"},
		{candidate: "latest", current: "1.2.3"},
		{candidate: "1.3.0", current: "dev"},
	}
	for _, test := range tests {
		if got := newerReleaseVersion(test.candidate, test.current); got != test.newer {
			t.Errorf("newerReleaseVersion(%q, %q) = %t, want %t", test.candidate, test.current, got, test.newer)
		}
	}
}
