package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFilesCreateFromURLWritesJSONResponse(t *testing.T) {
	isolateConfig(t)
	t.Setenv("RUNAPI_API_KEY", "test-key")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" || r.URL.Path != "/api/v1/files" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		if got := r.Header.Get("Content-Type"); !strings.Contains(got, "application/json") {
			t.Fatalf("expected JSON content type, got %q", got)
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		source := body["source"].(map[string]any)
		if source["type"] != "url" || source["url"] != "https://cdn.runapi.ai/public/samples/mask.png" {
			t.Fatalf("unexpected body: %#v", body)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"file_name":  "image.png",
			"url":        "https://file.runapi.ai/temp/image.png",
			"size_bytes": 204800,
			"mime_type":  "image/png",
			"created_at": "2026-06-08T10:30:00Z",
			"expires_at": "2026-06-08T11:30:00Z",
		})
	}))
	defer server.Close()

	c := newCLI()
	var stdout, stderr bytes.Buffer
	c.stdout = &stdout
	c.stderr = &stderr
	code := c.run([]string{"--base-url", server.URL, "files", "create", "--url", "https://cdn.runapi.ai/public/samples/mask.png", "--file-name", "image.png"})
	if code != 0 {
		t.Fatalf("expected zero exit code, got %d stderr=%s stdout=%s", code, stderr.String(), stdout.String())
	}
	var payload map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &payload); err != nil {
		t.Fatal(err)
	}
	if payload["url"] != "https://file.runapi.ai/temp/image.png" || payload["expires_at"] == "" {
		t.Fatalf("unexpected payload: %#v", payload)
	}
}

func TestFilesCreateFromPathCanPrintURLOnly(t *testing.T) {
	isolateConfig(t)
	t.Setenv("RUNAPI_API_KEY", "test-key")
	dir := t.TempDir()
	path := filepath.Join(dir, "image.png")
	if err := os.WriteFile(path, []byte("png"), 0o644); err != nil {
		t.Fatal(err)
	}

	var putReceived []byte
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "POST" && r.URL.Path == "/api/v1/files/prepare":
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatal(err)
			}
			if body["filename"] != "image.png" || body["checksum"] == "" {
				t.Fatalf("unexpected prepare body: %#v", body)
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"signed_id":  "signed-blob-id",
				"upload_url": server.URL + "/put-target",
				"headers":    map[string]string{"Content-Type": "application/octet-stream"},
			})
		case r.Method == "PUT" && r.URL.Path == "/put-target":
			putReceived, _ = io.ReadAll(r.Body)
			w.WriteHeader(http.StatusOK)
		case r.Method == "POST" && r.URL.Path == "/api/v1/files/confirm":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"file_name":  "image.png",
				"url":        "https://file.runapi.ai/temp/image.png",
				"size_bytes": 3,
				"mime_type":  "image/png",
				"created_at": "2026-06-08T10:30:00Z",
				"expires_at": "2026-06-08T11:30:00Z",
			})
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	}))
	defer server.Close()

	c := newCLI()
	var stdout bytes.Buffer
	c.stdout = &stdout
	c.stderr = &bytes.Buffer{}
	code := c.run([]string{"--base-url", server.URL, "files", "create", path, "--file-name", "image.png", "--url-only"})
	if code != 0 {
		t.Fatalf("expected zero exit code, got %d", code)
	}
	if string(putReceived) != "png" {
		t.Fatalf("expected file bytes at the upload URL, got %q", putReceived)
	}
	if got := strings.TrimSpace(stdout.String()); got != "https://file.runapi.ai/temp/image.png" {
		t.Fatalf("unexpected stdout: %q", got)
	}
}

func TestFilesCreateRejectsMultipleSources(t *testing.T) {
	isolateConfig(t)
	t.Setenv("RUNAPI_API_KEY", "test-key")
	c := newCLI()
	c.stdout = &bytes.Buffer{}
	c.stderr = &bytes.Buffer{}

	code := c.run([]string{"files", "create", "image.png", "--url", "https://cdn.runapi.ai/public/samples/mask.png"})
	if code == 0 {
		t.Fatal("expected non-zero exit code")
	}
	if !strings.Contains(c.stdout.(*bytes.Buffer).String(), "exactly one source") {
		t.Fatalf("expected validation error, got %s", c.stdout.(*bytes.Buffer).String())
	}
}

func TestFilesCreateDoesNotExposeContentTypeOverride(t *testing.T) {
	isolateConfig(t)
	t.Setenv("RUNAPI_API_KEY", "test-key")
	c := newCLI()
	c.stdout = &bytes.Buffer{}
	c.stderr = &bytes.Buffer{}

	code := c.run([]string{"files", "create", "image.png", "--content-type", "application/x-javascript"})
	if code == 0 {
		t.Fatal("expected non-zero exit code")
	}
	if !strings.Contains(c.stdout.(*bytes.Buffer).String(), "unknown flag: --content-type") {
		t.Fatalf("expected unknown flag error, got %s", c.stdout.(*bytes.Buffer).String())
	}
}
