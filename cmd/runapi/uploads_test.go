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

func TestUploadsCommandsCoverLifecycle(t *testing.T) {
	isolateConfig(t)
	t.Setenv("RUNAPI_API_KEY", "test-key")
	var calls []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		if strings.HasSuffix(r.URL.Path, "/parts") {
			if err := r.ParseMultipartForm(1 << 20); err != nil {
				t.Fatal(err)
			}
			_, _ = w.Write([]byte(`{"id":"part_123","object":"upload.part","created_at":1,"upload_id":"upload_123"}`))
			return
		}
		if strings.HasSuffix(r.URL.Path, "/complete") {
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatal(err)
			}
			if body["part_ids"].([]any)[0] != "part_123" {
				t.Fatalf("unexpected completion body: %#v", body)
			}
		}
		_, _ = w.Write([]byte(`{"id":"upload_123","object":"upload","bytes":3,"created_at":1,"filename":"data.bin","purpose":"user_data","status":"pending","expires_at":2}`))
	}))
	defer server.Close()
	part := filepath.Join(t.TempDir(), "part.bin")
	if err := os.WriteFile(part, []byte("abc"), 0o600); err != nil {
		t.Fatal(err)
	}

	commands := [][]string{
		{"--base-url", server.URL, "uploads", "create", "--bytes", "3", "--filename", "data.bin", "--mime-type", "application/octet-stream"},
		{"--base-url", server.URL, "uploads", "add-part", "upload_123", part},
		{"--base-url", server.URL, "uploads", "complete", "upload_123", "--part-id", "part_123"},
		{"--base-url", server.URL, "uploads", "cancel", "upload_123"},
	}
	for _, args := range commands {
		c := newCLI()
		c.stdout, c.stderr = &bytes.Buffer{}, &bytes.Buffer{}
		if code := c.run(args); code != 0 {
			t.Fatalf("command %v failed with %d: %s", args, code, c.stdout.(*bytes.Buffer).String())
		}
	}
	want := []string{"/v1/uploads", "/v1/uploads/upload_123/parts", "/v1/uploads/upload_123/complete", "/v1/uploads/upload_123/cancel"}
	if strings.Join(calls, ",") != strings.Join(want, ",") {
		t.Fatalf("unexpected calls: %#v", calls)
	}
}
