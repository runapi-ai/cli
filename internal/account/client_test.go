package account

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/runapi-ai/core-sdk/go/core"
)

type stubHTTPClient struct {
	method string
	path   string
}

func (s *stubHTTPClient) Request(_ context.Context, method, path string, _ *core.HTTPRequestOptions) (json.RawMessage, error) {
	s.method = method
	s.path = path
	return json.RawMessage(`{"id":1,"name":"test","email":"test@example.com","account":{"id":2,"name":"acme"}}`), nil
}

func TestInfoSendsCorrectRequest(t *testing.T) {
	stub := &stubHTTPClient{}
	client := NewClientWithHTTP(stub)
	resp, err := client.Info(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if stub.method != "GET" || stub.path != "/api/v1/me" {
		t.Fatalf("unexpected request: %s %s", stub.method, stub.path)
	}
	if resp.Email != "test@example.com" {
		t.Fatalf("unexpected email: %s", resp.Email)
	}
	if resp.Account.Name != "acme" {
		t.Fatalf("unexpected account name: %s", resp.Account.Name)
	}
}

type balanceStubHTTPClient struct {
	method string
	path   string
}

func (s *balanceStubHTTPClient) Request(_ context.Context, method, path string, _ *core.HTTPRequestOptions) (json.RawMessage, error) {
	s.method = method
	s.path = path
	return json.RawMessage(`{"balance_cents":5000,"spent_cents_today":100,"spent_cents_total":2000}`), nil
}

func TestBalanceSendsCorrectRequest(t *testing.T) {
	stub := &balanceStubHTTPClient{}
	client := NewClientWithHTTP(stub)
	resp, err := client.Balance(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if stub.method != "GET" || stub.path != "/api/v1/me/balance" {
		t.Fatalf("unexpected request: %s %s", stub.method, stub.path)
	}
	if resp.BalanceCents != 5000 {
		t.Fatalf("unexpected balance: %d", resp.BalanceCents)
	}
	if resp.SpentCentsToday != 100 {
		t.Fatalf("unexpected spent today: %d", resp.SpentCentsToday)
	}
}
