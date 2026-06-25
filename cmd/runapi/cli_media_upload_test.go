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
)

func TestServiceCommandUploadsLocalMediaInputsBeforeCreate(t *testing.T) {
	isolateConfig(t)
	t.Setenv("RUNAPI_API_KEY", "test-key")
	dir := t.TempDir()
	localImage := filepath.Join(dir, "input.png")
	if err := os.WriteFile(localImage, []byte("png"), 0o644); err != nil {
		t.Fatal(err)
	}

	var sawUpload, sawCreate bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "POST" && r.URL.Path == "/api/v1/files":
			sawUpload = true
			if got := r.Header.Get("Content-Type"); !strings.Contains(got, "multipart/form-data") {
				t.Fatalf("expected multipart upload, got %q", got)
			}
			if err := r.ParseMultipartForm(1024); err != nil {
				t.Fatal(err)
			}
			file, header, err := r.FormFile("file")
			if err != nil {
				t.Fatal(err)
			}
			defer file.Close()
			if header.Filename != "input.png" {
				t.Fatalf("unexpected upload filename: %s", header.Filename)
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"file_name":  "input.png",
				"url":        "https://files.runapi.ai/temp/input.png",
				"size_bytes": 3,
				"mime_type":  "image/png",
				"created_at": "2026-06-08T10:30:00Z",
				"expires_at": "2026-06-08T11:30:00Z",
			})
		case r.Method == "POST" && r.URL.Path == "/api/v1/gpt_image/edit_image":
			sawCreate = true
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatal(err)
			}
			urls := body["source_image_urls"].([]any)
			if urls[0] != "https://files.runapi.ai/temp/input.png" {
				t.Fatalf("expected uploaded URL in create body, got %#v", body)
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{"id": "task_123", "status": "processing"})
		case r.Method == "GET" && r.URL.Path == "/api/v1/gpt_image/edit_image/task_123":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{"id": "task_123", "status": "completed", "images": []any{}})
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	}))
	defer server.Close()

	input := `{"model":"gpt-image-1.5","prompt":"edit","source_image_urls":["` + localImage + `"],"aspect_ratio":"1:1","quality":"medium"}`
	c := newCLI()
	var stdout, stderr bytes.Buffer
	c.stdout = &stdout
	c.stderr = &stderr

	code := c.run([]string{"--base-url", server.URL, "gpt-image", "edit-image", "--input", input})
	if code != 0 {
		t.Fatalf("expected zero exit code, got %d stderr=%s stdout=%s", code, stderr.String(), stdout.String())
	}
	if !sawUpload || !sawCreate {
		t.Fatalf("expected upload and create requests, sawUpload=%v sawCreate=%v", sawUpload, sawCreate)
	}
	if !strings.Contains(stderr.String(), "uploaded source_image_urls[0] from "+localImage) {
		t.Fatalf("expected upload log, got %q", stderr.String())
	}
}

func TestServiceCommandPreservesNullMediaInputsWhenUploadingLocalPaths(t *testing.T) {
	isolateConfig(t)
	t.Setenv("RUNAPI_API_KEY", "test-key")
	dir := t.TempDir()
	localImage := filepath.Join(dir, "input.png")
	if err := os.WriteFile(localImage, []byte("png"), 0o644); err != nil {
		t.Fatal(err)
	}

	var createBody map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "POST" && r.URL.Path == "/api/v1/files":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"file_name":  "input.png",
				"url":        "https://files.runapi.ai/temp/input.png",
				"size_bytes": 3,
				"mime_type":  "image/png",
				"created_at": "2026-06-08T10:30:00Z",
				"expires_at": "2026-06-08T11:30:00Z",
			})
		case r.Method == "POST" && r.URL.Path == "/api/v1/gpt_image/edit_image":
			if err := json.NewDecoder(r.Body).Decode(&createBody); err != nil {
				t.Fatal(err)
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{"id": "task_789", "status": "processing"})
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	}))
	defer server.Close()

	input := `{"model":"gpt-image-1.5","prompt":"edit","source_image_urls":["https://example.com/input.png",null,"` + localImage + `"],"aspect_ratio":"1:1","quality":"medium"}`
	c := newCLI()
	var stdout, stderr bytes.Buffer
	c.stdout = &stdout
	c.stderr = &stderr

	code := c.run([]string{"--base-url", server.URL, "gpt-image", "edit-image", "--async", "--input", input})
	if code != 0 {
		t.Fatalf("expected zero exit code, got %d stderr=%s stdout=%s", code, stderr.String(), stdout.String())
	}

	urls := createBody["source_image_urls"].([]any)
	if urls[0] != "https://example.com/input.png" || urls[1] != nil || urls[2] != "https://files.runapi.ai/temp/input.png" {
		t.Fatalf("expected remote URL, null, uploaded URL in create body, got %#v", createBody)
	}
}

func TestServiceCommandDoesNotUploadCallbackOrRemoteMediaURL(t *testing.T) {
	isolateConfig(t)
	t.Setenv("RUNAPI_API_KEY", "test-key")

	var createBody map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/files" {
			t.Fatal("did not expect file upload")
		}
		if r.Method != "POST" || r.URL.Path != "/api/v1/gpt_image/edit_image" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		if err := json.NewDecoder(r.Body).Decode(&createBody); err != nil {
			t.Fatal(err)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"id": "task_456", "status": "processing"})
	}))
	defer server.Close()

	input := `{"model":"gpt-image-1.5","prompt":"edit","source_image_urls":["https://example.com/input.png"],"aspect_ratio":"1:1","quality":"medium","callback_url":"./callback.json"}`
	c := newCLI()
	var stdout, stderr bytes.Buffer
	c.stdout = &stdout
	c.stderr = &stderr

	code := c.run([]string{"--base-url", server.URL, "gpt-image", "edit-image", "--async", "--input", input})
	if code != 0 {
		t.Fatalf("expected zero exit code, got %d stderr=%s stdout=%s", code, stderr.String(), stdout.String())
	}
	if createBody["callback_url"] != "./callback.json" {
		t.Fatalf("expected callback_url to stay untouched, got %#v", createBody)
	}
	urls := createBody["source_image_urls"].([]any)
	if urls[0] != "https://example.com/input.png" {
		t.Fatalf("expected remote media URL to stay untouched, got %#v", createBody)
	}
	if strings.Contains(stderr.String(), "uploaded ") {
		t.Fatalf("did not expect upload log, got %q", stderr.String())
	}
}

func TestServiceCommandStopsWhenMediaUploadFails(t *testing.T) {
	isolateConfig(t)
	t.Setenv("RUNAPI_API_KEY", "test-key")
	dir := t.TempDir()
	localImage := filepath.Join(dir, "input.png")
	if err := os.WriteFile(localImage, []byte("png"), 0o644); err != nil {
		t.Fatal(err)
	}

	var sawCreate bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/files":
			http.Error(w, `{"error":{"message":"upload failed"}}`, http.StatusUnprocessableEntity)
		case "/api/v1/gpt_image/edit_image":
			sawCreate = true
			t.Fatal("did not expect create after failed upload")
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	}))
	defer server.Close()

	input := `{"model":"gpt-image-1.5","prompt":"edit","source_image_urls":["` + localImage + `"],"aspect_ratio":"1:1","quality":"medium"}`
	c := newCLI()
	var stdout, stderr bytes.Buffer
	c.stdout = &stdout
	c.stderr = &stderr

	code := c.run([]string{"--base-url", server.URL, "gpt-image", "edit-image", "--async", "--input", input})
	if code == 0 {
		t.Fatal("expected non-zero exit code")
	}
	if sawCreate {
		t.Fatal("create request should not be sent after upload failure")
	}
	if !strings.Contains(stdout.String(), "upload failed") {
		t.Fatalf("expected upload error, got stdout=%s stderr=%s", stdout.String(), stderr.String())
	}
}

func TestServiceCommandDoesNotUploadWhenInputDecodeFails(t *testing.T) {
	isolateConfig(t)
	t.Setenv("RUNAPI_API_KEY", "test-key")
	dir := t.TempDir()
	localImage := filepath.Join(dir, "input.png")
	if err := os.WriteFile(localImage, []byte("png"), 0o644); err != nil {
		t.Fatal(err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("did not expect request before input decode succeeds: %s %s", r.Method, r.URL.Path)
	}))
	defer server.Close()

	input := `{"model":"gpt-image-1.5","prompt":"edit","source_image_urls":["` + localImage + `"],"aspect_ratio":["1:1"],"quality":"medium"}`
	c := newCLI()
	var stdout, stderr bytes.Buffer
	c.stdout = &stdout
	c.stderr = &stderr

	code := c.run([]string{"--base-url", server.URL, "gpt-image", "edit-image", "--async", "--input", input})
	if code == 0 {
		t.Fatal("expected non-zero exit code")
	}
	if !strings.Contains(stdout.String(), "input must be valid JSON for the selected action") {
		t.Fatalf("expected input decode error, got stdout=%s stderr=%s", stdout.String(), stderr.String())
	}
	if strings.Contains(stderr.String(), "uploaded ") {
		t.Fatalf("did not expect upload log, got %q", stderr.String())
	}
}
