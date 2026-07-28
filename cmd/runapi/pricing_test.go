package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestPricingCommandsAreRegistered(t *testing.T) {
	c := newCLI()
	c.stdout, c.stderr = &bytes.Buffer{}, &bytes.Buffer{}
	cmd := c.command()
	cmd.SetArgs([]string{"pricing", "quote", "--help"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	output := c.stdout.(*bytes.Buffer).String()
	for _, flag := range []string{"--service", "--action", "--model", "--params", "--params-file"} {
		if !bytes.Contains([]byte(output), []byte(flag)) {
			t.Fatalf("expected %s in help: %s", flag, output)
		}
	}
}

func TestPricingListIsPublicAndWritesJSON(t *testing.T) {
	isolateConfig(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/price_schedules" || r.URL.Query().Get("service") != "suno" {
			t.Fatalf("unexpected request: %s", r.URL.String())
		}
		if got := r.Header.Get("Authorization"); got != "" {
			t.Fatalf("expected no authorization, got %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"as_of":"2026-07-23T00:00:00.000000Z","price_schedules":[]}`))
	}))
	defer server.Close()
	c := newCLI()
	c.stdout, c.stderr = &bytes.Buffer{}, &bytes.Buffer{}
	code := c.run([]string{"--base-url", server.URL, "pricing", "list", "--service", "suno"})
	if code != 0 {
		t.Fatalf("expected success, got %d: %s", code, c.stderr.(*bytes.Buffer).String())
	}
	var response map[string]any
	if err := json.Unmarshal(c.stdout.(*bytes.Buffer).Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response["as_of"] == nil {
		t.Fatalf("unexpected response: %#v", response)
	}
}

func TestPricingQuoteAcceptsOptionalInlineAndFileParams(t *testing.T) {
	isolateConfig(t)
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if r.URL.Path != "/api/v1/price_quotes" || r.Method != http.MethodPost {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body["service"] != "suno" || body["action"] != "convert_audio" {
			t.Fatalf("unexpected body: %#v", body)
		}
		params, ok := body["params"].(map[string]any)
		if !ok {
			t.Fatalf("expected params object: %#v", body)
		}
		if requests == 1 && len(params) != 0 {
			t.Fatalf("expected omitted params to send an empty object: %#v", params)
		}
		if requests > 1 && params["quality"] != "high" {
			t.Fatalf("unexpected params: %#v", params)
		}
		_, _ = w.Write([]byte(`{"price_quote":{"service":"suno","action":"convert_audio","pricing_status":"available","currency":"USD","reservation_amount_cents":10,"estimate_basis":"exact","as_of":"2026-07-23T00:00:00.000000Z"}}`))
	}))
	defer server.Close()
	file := t.TempDir() + "/params.json"
	if err := os.WriteFile(file, []byte(`{"quality":"high"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	for _, paramsArgs := range [][]string{nil, {"--params", `{"quality":"high"}`}, {"--params-file", file}, {"--params-file", "-"}} {
		c := newCLI()
		c.stdout, c.stderr = &bytes.Buffer{}, &bytes.Buffer{}
		if len(paramsArgs) == 2 && paramsArgs[1] == "-" {
			c.stdin = strings.NewReader(`{"quality":"high"}`)
		}
		args := append([]string{"--base-url", server.URL, "pricing", "quote", "--service", "suno", "--action", "convert_audio"}, paramsArgs...)
		if code := c.run(args); code != 0 {
			t.Fatalf("expected success, got %d: %s", code, c.stderr.(*bytes.Buffer).String())
		}
	}
	if requests != 4 {
		t.Fatalf("expected four quote requests, got %d", requests)
	}
}

func TestPricingQuoteRejectsConflictingParamInputs(t *testing.T) {
	isolateConfig(t)
	c := newCLI()
	c.stdout, c.stderr = &bytes.Buffer{}, &bytes.Buffer{}
	code := c.run([]string{"pricing", "quote", "--service", "suno", "--action", "convert_audio", "--params", `{}`, "--params-file", "params.json"})
	if code == 0 {
		t.Fatal("expected conflicting parameter inputs to fail")
	}
	var response map[string]map[string]any
	if err := json.Unmarshal(c.stdout.(*bytes.Buffer).Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response["error"]["message"] != "--params and --params-file are mutually exclusive" {
		t.Fatalf("unexpected error: %#v", response)
	}
}

func TestPricingQuoteSurfacesProtectedContextAuthenticationFailure(t *testing.T) {
	isolateConfig(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"code":"protected_pricing_context_not_found","message":"Protected pricing context not found"}}`))
	}))
	defer server.Close()
	c := newCLI()
	c.stdout, c.stderr = &bytes.Buffer{}, &bytes.Buffer{}
	code := c.run([]string{"--base-url", server.URL, "pricing", "quote", "--service", "grok_imagine", "--action", "extend_video", "--params", `{"source_task_id":"task_123"}`})
	if code == 0 {
		t.Fatal("expected protected-context failure")
	}
	var response map[string]map[string]any
	if err := json.Unmarshal(c.stdout.(*bytes.Buffer).Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response["error"]["code"] != "protected_pricing_context_not_found" {
		t.Fatalf("unexpected error: %#v", response)
	}
}
