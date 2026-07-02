package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"
)

func TestAgentListTargets(t *testing.T) {
	c, out, _ := newTestCLIWithBuffers()
	if rc := c.run([]string{"agent", "list-targets"}); rc != 0 {
		t.Fatalf("expected exit 0, got %d", rc)
	}
	var payload struct {
		Skill   string            `json:"skill"`
		Targets map[string]string `json:"targets"`
	}
	if err := json.Unmarshal(out.Bytes(), &payload); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if payload.Skill != "runapi-cli" {
		t.Fatalf("expected skill runapi-cli, got %q", payload.Skill)
	}
	for _, k := range []string{"claude", "codex", "gemini", "openclaw", "hermes"} {
		if _, ok := payload.Targets[k]; !ok {
			t.Fatalf("missing target %q", k)
		}
	}
}

func TestNormalizeSkillTag(t *testing.T) {
	cases := []struct {
		in      string
		want    string
		wantErr bool
	}{
		{in: "v0.1.0", want: "v0.1.0"},
		{in: "0.1.0", want: "v0.1.0"},
		{in: "../etc", wantErr: true},
	}
	for _, tc := range cases {
		got, err := normalizeSkillTag(tc.in)
		if tc.wantErr {
			if err == nil {
				t.Fatalf("input %q: expected error, got %q", tc.in, got)
			}
			continue
		}
		if err != nil {
			t.Fatalf("input %q: unexpected error: %v", tc.in, err)
		}
		if got != tc.want {
			t.Fatalf("input %q: want %q, got %q", tc.in, tc.want, got)
		}
	}
}

func TestDecodeSkillTagsKeepsStableReleaseTags(t *testing.T) {
	tags, err := decodeSkillTags(strings.NewReader(`[
		{"name":"v0.2.9"},
		{"name":"v0.2.10"},
		{"name":"v0.3.0-rc.1"},
		{"name":"v0.2.10"},
		{"name":"not-a-release"}
	]`))
	if err != nil {
		t.Fatal(err)
	}
	got := strings.Join(tags, ",")
	want := "v0.2.9,v0.2.10"
	if got != want {
		t.Fatalf("want %s, got %s", want, got)
	}
}

func TestCompareSkillTagsSortsSemver(t *testing.T) {
	tags := []string{"v0.2.9", "v0.10.0", "v0.2.10", "v1.0.0"}
	sort.Slice(tags, func(i, j int) bool {
		return compareSkillTags(tags[i], tags[j]) > 0
	})
	got := strings.Join(tags, ",")
	want := "v1.0.0,v0.10.0,v0.2.10,v0.2.9"
	if got != want {
		t.Fatalf("want %s, got %s", want, got)
	}
}

func TestInstallSkillWritesFiles(t *testing.T) {
	tmp := t.TempDir()
	repo := "runapi-ai/cli-skill"
	tag := "v9.9.9"

	archive := buildSkillTarball(t, map[string]string{
		"cli-skill-9.9.9/skills/runapi-cli/SKILL.md":            "# skill content\n",
		"cli-skill-9.9.9/skills/runapi-cli/references/notes.md": "notes\n",
		"cli-skill-9.9.9/README.md":                             "ignored\n",
	})

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/" + repo + "/archive/refs/tags/" + tag + ".tar.gz"
		if r.URL.Path != expectedPath {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/gzip")
		w.Write(archive)
	}))
	defer srv.Close()

	c, out, _ := newTestCLIWithBuffers()
	c.archiveBaseURL = srv.URL
	c.httpClient = &http.Client{Timeout: 5 * time.Second}

	rc := c.run([]string{
		"agent", "install-skill",
		"--target-dir", tmp,
		"--version", tag,
		"--source", repo,
	})
	if rc != 0 {
		t.Fatalf("install-skill exit %d", rc)
	}

	skillRoot := filepath.Join(tmp, "runapi-cli")
	if _, err := os.Stat(filepath.Join(skillRoot, "SKILL.md")); err != nil {
		t.Fatalf("SKILL.md not written: %v", err)
	}
	if _, err := os.Stat(filepath.Join(skillRoot, "references", "notes.md")); err != nil {
		t.Fatalf("references/notes.md not written: %v", err)
	}
	if _, err := os.Stat(filepath.Join(skillRoot, "..", "README.md")); err == nil {
		t.Fatalf("non-skill file should not be written outside skill dir")
	}

	var payload struct {
		Installed bool   `json:"installed"`
		Files     int    `json:"files"`
		TargetDir string `json:"target_dir"`
	}
	if err := json.Unmarshal(out.Bytes(), &payload); err != nil {
		t.Fatalf("decode stdout: %v", err)
	}
	if !payload.Installed || payload.Files != 2 {
		t.Fatalf("unexpected payload: %+v", payload)
	}
}

func TestInstallSkillRejectsExistingWithoutForce(t *testing.T) {
	tmp := t.TempDir()
	existing := filepath.Join(tmp, "runapi-cli")
	if err := os.MkdirAll(existing, 0o755); err != nil {
		t.Fatal(err)
	}
	c, _, _ := newTestCLIWithBuffers()
	rc := c.run([]string{
		"agent", "install-skill",
		"--target-dir", tmp,
		"--version", "v9.9.9",
	})
	if rc == 0 {
		t.Fatalf("expected non-zero exit when skill already exists")
	}
}

func TestInstallSkillForceOverwrites(t *testing.T) {
	tmp := t.TempDir()
	stale := filepath.Join(tmp, "runapi-cli", "stale.md")
	if err := os.MkdirAll(filepath.Dir(stale), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(stale, []byte("stale"), 0o644); err != nil {
		t.Fatal(err)
	}

	repo := "runapi-ai/cli-skill"
	tag := "v9.9.9"
	archive := buildSkillTarball(t, map[string]string{
		"cli-skill-9.9.9/skills/runapi-cli/SKILL.md": "fresh\n",
	})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(archive)
	}))
	defer srv.Close()

	c, _, _ := newTestCLIWithBuffers()
	c.archiveBaseURL = srv.URL
	c.httpClient = &http.Client{Timeout: 5 * time.Second}

	rc := c.run([]string{
		"agent", "install-skill",
		"--target-dir", tmp,
		"--version", tag,
		"--source", repo,
		"--force",
	})
	if rc != 0 {
		t.Fatalf("force install exit %d", rc)
	}
	if _, err := os.Stat(stale); !os.IsNotExist(err) {
		t.Fatalf("stale file should be removed: err=%v", err)
	}
	if data, err := os.ReadFile(filepath.Join(tmp, "runapi-cli", "SKILL.md")); err != nil || string(data) != "fresh\n" {
		t.Fatalf("expected fresh SKILL.md, got data=%q err=%v", string(data), err)
	}
}

func TestInstallSkillMissingArchive(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}))
	defer srv.Close()

	tmp := t.TempDir()
	c, _, _ := newTestCLIWithBuffers()
	c.archiveBaseURL = srv.URL
	c.httpClient = &http.Client{Timeout: 5 * time.Second}

	rc := c.run([]string{
		"agent", "install-skill",
		"--target-dir", tmp,
		"--version", "v9.9.9",
	})
	if rc == 0 {
		t.Fatalf("expected non-zero exit when archive missing")
	}
}

func TestInstallSkillDefaultsToLatestStableTag(t *testing.T) {
	tmp := t.TempDir()
	repo := "runapi-ai/cli-skill"
	latestTag := "v0.2.10"
	archive := buildSkillTarball(t, map[string]string{
		"cli-skill-0.2.10/skills/runapi-cli/SKILL.md": "fallback\n",
	})

	var requested []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requested = append(requested, r.URL.Path)
		switch r.URL.Path {
		case "/repos/" + repo + "/tags":
			if got := r.Header.Get("Accept"); got != "application/vnd.github+json" {
				http.Error(w, "json accept required", http.StatusUnsupportedMediaType)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`[
				{"name":"v0.2.7"},
				{"name":"v0.2.9"},
				{"name":"v0.2.10"},
				{"name":"v0.3.0-rc.1"},
				{"name":"not-a-release"}
			]`))
		case "/" + repo + "/archive/refs/tags/" + latestTag + ".tar.gz":
			if got := r.Header.Get("Accept"); got != "application/octet-stream" {
				t.Errorf("archive Accept header = %q", got)
				http.Error(w, "archive accept required", http.StatusUnsupportedMediaType)
				return
			}
			w.Header().Set("Content-Type", "application/gzip")
			w.Write(archive)
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	c, out, errBuf := newTestCLIWithBuffers()
	c.archiveBaseURL = srv.URL
	c.httpClient = &http.Client{Timeout: 5 * time.Second}

	rc := c.run([]string{
		"agent", "install-skill",
		"--target-dir", tmp,
		"--source", repo,
	})
	if rc != 0 {
		t.Fatalf("install-skill exit %d stderr=%s", rc, errBuf.String())
	}

	if data, err := os.ReadFile(filepath.Join(tmp, "runapi-cli", "SKILL.md")); err != nil || string(data) != "fallback\n" {
		t.Fatalf("expected latest skill, got data=%q err=%v", string(data), err)
	}

	var payload struct {
		Version string `json:"version"`
		Source  string `json:"source"`
	}
	if err := json.Unmarshal(out.Bytes(), &payload); err != nil {
		t.Fatalf("decode stdout: %v", err)
	}
	if payload.Version != latestTag || payload.Source != repo {
		t.Fatalf("unexpected payload: %+v", payload)
	}

	want := []string{
		"/repos/" + repo + "/tags",
		"/" + repo + "/archive/refs/tags/" + latestTag + ".tar.gz",
	}
	if strings.Join(requested, "\n") != strings.Join(want, "\n") {
		t.Fatalf("unexpected request sequence:\n%s", strings.Join(requested, "\n"))
	}
}

func TestInstallSkillExplicitVersionUsesPinnedTag(t *testing.T) {
	tmp := t.TempDir()
	repo := "runapi-ai/cli-skill"
	tag := "v0.2.11"
	archive := buildSkillTarball(t, map[string]string{
		"cli-skill-0.2.11/skills/runapi-cli/SKILL.md": "matching\n",
	})

	var requested []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requested = append(requested, r.URL.Path)
		switch r.URL.Path {
		case "/" + repo + "/archive/refs/tags/" + tag + ".tar.gz":
			if got := r.Header.Get("Accept"); got != "application/octet-stream" {
				t.Errorf("archive Accept header = %q", got)
				http.Error(w, "archive accept required", http.StatusUnsupportedMediaType)
				return
			}
			w.Header().Set("Content-Type", "application/gzip")
			w.Write(archive)
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	c, out, errBuf := newTestCLIWithBuffers()
	c.archiveBaseURL = srv.URL
	c.httpClient = &http.Client{Timeout: 5 * time.Second}

	rc := c.run([]string{
		"agent", "install-skill",
		"--target-dir", tmp,
		"--source", repo,
		"--version", tag,
	})
	if rc != 0 {
		t.Fatalf("install-skill exit %d stderr=%s", rc, errBuf.String())
	}
	if data, err := os.ReadFile(filepath.Join(tmp, "runapi-cli", "SKILL.md")); err != nil || string(data) != "matching\n" {
		t.Fatalf("expected matching skill, got data=%q err=%v", string(data), err)
	}

	var payload struct {
		Version string `json:"version"`
	}
	if err := json.Unmarshal(out.Bytes(), &payload); err != nil {
		t.Fatalf("decode stdout: %v", err)
	}
	if payload.Version != tag {
		t.Fatalf("expected version %s, got %+v", tag, payload)
	}

	want := []string{
		"/" + repo + "/archive/refs/tags/" + tag + ".tar.gz",
	}
	if strings.Join(requested, "\n") != strings.Join(want, "\n") {
		t.Fatalf("unexpected request sequence:\n%s", strings.Join(requested, "\n"))
	}
}

func TestInstallSkillExplicitVersionDoesNotFallback(t *testing.T) {
	tmp := t.TempDir()
	repo := "runapi-ai/cli-skill"
	tag := "v0.2.11"

	var requested []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requested = append(requested, r.URL.Path)
		http.NotFound(w, r)
	}))
	defer srv.Close()

	c, _, _ := newTestCLIWithBuffers()
	c.archiveBaseURL = srv.URL
	c.httpClient = &http.Client{Timeout: 5 * time.Second}

	rc := c.run([]string{
		"agent", "install-skill",
		"--target-dir", tmp,
		"--source", repo,
		"--version", tag,
	})
	if rc == 0 {
		t.Fatalf("expected non-zero exit when explicit archive is missing")
	}

	want := []string{"/" + repo + "/archive/refs/tags/" + tag + ".tar.gz"}
	if strings.Join(requested, "\n") != strings.Join(want, "\n") {
		t.Fatalf("unexpected request sequence:\n%s", strings.Join(requested, "\n"))
	}
}

func TestInstallSkillRejectsPathTraversalInArchive(t *testing.T) {
	repo := "runapi-ai/cli-skill"
	tag := "v9.9.9"
	archive := buildSkillTarball(t, map[string]string{
		"cli-skill-9.9.9/skills/runapi-cli/../../etc/evil": "bad\n",
	})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(archive)
	}))
	defer srv.Close()

	tmp := t.TempDir()
	c, _, _ := newTestCLIWithBuffers()
	c.archiveBaseURL = srv.URL
	c.httpClient = &http.Client{Timeout: 5 * time.Second}

	rc := c.run([]string{
		"agent", "install-skill",
		"--target-dir", tmp,
		"--version", tag,
		"--source", repo,
	})
	if rc == 0 {
		t.Fatalf("expected non-zero exit on path traversal")
	}
}

func TestInstallSkillRejectsHTTPArchiveBase(t *testing.T) {
	tmp := t.TempDir()
	c, _, _ := newTestCLIWithBuffers()
	c.archiveBaseURL = "http://attacker.example/" // non-loopback http
	rc := c.run([]string{
		"agent", "install-skill",
		"--target-dir", tmp,
		"--version", "v9.9.9",
	})
	if rc == 0 {
		t.Fatalf("expected non-zero exit when archive URL is plain http on non-loopback")
	}
}

func TestInstallSkillRejectsOversizedArchive(t *testing.T) {
	// Build a tarball where skills/runapi-cli/big.bin is bigger than the
	// extract budget so the size cap trips. Keep the buffer modest by
	// overriding the package-level constant via a test-only helper.
	prevCap := skillMaxExtractBytes
	defer func() { setSkillMaxExtractBytesForTest(prevCap) }()
	setSkillMaxExtractBytesForTest(1024) // 1 KiB cap

	big := make([]byte, 4096)
	for i := range big {
		big[i] = 'x'
	}
	archive := buildSkillTarball(t, map[string]string{
		"cli-skill-9.9.9/skills/runapi-cli/big.bin": string(big),
	})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(archive)
	}))
	defer srv.Close()

	tmp := t.TempDir()
	c, _, _ := newTestCLIWithBuffers()
	c.archiveBaseURL = srv.URL
	c.httpClient = &http.Client{Timeout: 5 * time.Second}

	rc := c.run([]string{
		"agent", "install-skill",
		"--target-dir", tmp,
		"--version", "v9.9.9",
	})
	if rc == 0 {
		t.Fatalf("expected non-zero exit when archive exceeds the extract budget")
	}
}

func TestUninstallSkillRemovesDir(t *testing.T) {
	tmp := t.TempDir()
	skill := filepath.Join(tmp, "runapi-cli")
	if err := os.MkdirAll(skill, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(skill, "SKILL.md"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	c, _, _ := newTestCLIWithBuffers()
	rc := c.run([]string{"agent", "uninstall-skill", "--target-dir", tmp})
	if rc != 0 {
		t.Fatalf("uninstall exit %d", rc)
	}
	if _, err := os.Stat(skill); !os.IsNotExist(err) {
		t.Fatalf("skill dir should be gone, err=%v", err)
	}
}

// --- helpers ---

func newTestCLIWithBuffers() (*cli, *bytes.Buffer, *bytes.Buffer) {
	out := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}
	c := newCLI()
	c.stdout = out
	c.stderr = errBuf
	return c, out, errBuf
}

func buildSkillTarball(t *testing.T, files map[string]string) []byte {
	t.Helper()
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	for name, contents := range files {
		hdr := &tar.Header{
			Name: name,
			Mode: 0o644,
			Size: int64(len(contents)),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			t.Fatal(err)
		}
		if _, err := tw.Write([]byte(contents)); err != nil {
			t.Fatal(err)
		}
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}
