package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

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
