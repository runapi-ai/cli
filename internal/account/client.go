// Package account provides access to RunAPI account info and balance endpoints.
package account

import (
	"context"

	"github.com/runapi-ai/core-sdk/go/core"
	"github.com/runapi-ai/core-sdk/go/option"
)

// InfoResponse is the response from the account info endpoint.
type InfoResponse struct {
	ID      int64         `json:"id"`
	Name    string        `json:"name"`
	Email   string        `json:"email"`
	Account AccountRecord `json:"account"`
}

// AccountRecord contains the account details nested in InfoResponse.
type AccountRecord struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// BalanceResponse is the response from the account balance endpoint.
type BalanceResponse struct {
	BalanceCents    int64 `json:"balance_cents"`
	SpentCentsToday int64 `json:"spent_cents_today"`
	SpentCentsTotal int64 `json:"spent_cents_total"`
}

// Client provides account info and balance operations.
type Client struct {
	http core.HTTPClient
}

// NewClient creates an account client with the given options.
func NewClient(opts ...option.ClientOption) (*Client, error) {
	resolved, err := option.ResolveClientOptions(opts...)
	if err != nil {
		return nil, err
	}
	httpClient, err := core.NewHTTPClient(resolved)
	if err != nil {
		return nil, err
	}
	return NewClientWithHTTP(httpClient), nil
}

// NewClientWithHTTP creates an account client with a pre-configured HTTP transport.
func NewClientWithHTTP(httpClient core.HTTPClient) *Client {
	return &Client{http: httpClient}
}

// Info fetches the current user and account record.
func (c *Client) Info(ctx context.Context, opts ...option.RequestOption) (*InfoResponse, error) {
	requestOptions, _ := option.ResolveRequestOptions(opts...)
	return core.GetJSON[InfoResponse](ctx, c.http, "/api/v1/me", requestOptions)
}

// Balance fetches account balance and spend counters.
func (c *Client) Balance(ctx context.Context, opts ...option.RequestOption) (*BalanceResponse, error) {
	requestOptions, _ := option.ResolveRequestOptions(opts...)
	return core.GetJSON[BalanceResponse](ctx, c.http, "/api/v1/me/balance", requestOptions)
}
