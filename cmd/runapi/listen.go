package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	runapi "github.com/runapi-ai/cli/internal/runapi"
	"github.com/runapi-ai/core-sdk/go/core"
	"github.com/spf13/cobra"
)

var errSessionExpired = errors.New("session expired")

type listenSecretResponse struct {
	ListenSecret string `json:"listen_secret"`
}

type createSessionResponse struct {
	SessionToken string `json:"session_token"`
}

type listenEvent struct {
	ID         int64             `json:"id"`
	SignedBody string            `json:"signed_body"`
	Headers    map[string]string `json:"headers"`
}

type pollResponse struct {
	Events []listenEvent `json:"events"`
}

func (c *cli) listenCommand() *cobra.Command {
	var forwardTo string
	cmd := &cobra.Command{
		Use:   "listen",
		Short: "Receive RunAPI callbacks locally",
		Long:  "Forward task callbacks to a local URL. Tasks with a callback_url are also copied to the listener while still being delivered to that URL.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, _ := loadConfig()
			apiKey := firstNonEmpty(strings.TrimSpace(c.apiKeyFlag), strings.TrimSpace(os.Getenv("RUNAPI_API_KEY")), strings.TrimSpace(cfg.APIKey))
			baseURL := firstNonEmpty(strings.TrimSpace(c.baseURLFlag), strings.TrimSpace(os.Getenv("RUNAPI_BASE_URL")), strings.TrimSpace(cfg.BaseURL), core.DefaultBaseURL)
			if apiKey == "" {
				return fmt.Errorf("API key required (--api-key, RUNAPI_API_KEY, or runapi login)")
			}
			forwardTo, err := normalizeForwardURL(forwardTo)
			if err != nil {
				return err
			}

			hc := &http.Client{Timeout: 30 * time.Second}
			ctx, cancel := context.WithCancel(cmd.Context())
			defer cancel()

			secret, err := fetchListenSecret(ctx, hc, baseURL, apiKey)
			if err != nil {
				return fmt.Errorf("failed to fetch listen secret: %w", err)
			}

			sessionToken, err := createListenSession(ctx, hc, baseURL, apiKey)
			if err != nil {
				return fmt.Errorf("failed to create listen session: %w", err)
			}
			defer destroyListenSession(context.Background(), hc, baseURL, apiKey, sessionToken)

			fmt.Fprintf(c.stderr, "Ready! Webhook signing secret: %s\n", secret)
			if forwardTo != "" {
				fmt.Fprintf(c.stderr, "Forwarding to %s\n", forwardTo)
			}

			var lastID int64
			backoff := time.Second
			for ctx.Err() == nil {
				resp, err := pollEvents(ctx, hc, baseURL, apiKey, sessionToken, lastID)
				if err != nil {
					if errors.Is(err, errSessionExpired) {
						fmt.Fprintf(c.stderr, "[session expired] creating new session...\n")
						sessionToken, err = createListenSession(ctx, hc, baseURL, apiKey)
						if err != nil {
							fmt.Fprintf(c.stderr, "[error] %v, retrying in %v...\n", err, backoff)
							time.Sleep(backoff)
							if backoff < 30*time.Second {
								backoff *= 2
							}
							continue
						}
						lastID = 0
						backoff = time.Second
						continue
					}
					fmt.Fprintf(c.stderr, "[error] %v, retrying in %v...\n", err, backoff)
					time.Sleep(backoff)
					if backoff < 30*time.Second {
						backoff *= 2
					}
					continue
				}
				backoff = time.Second
				retryEvent := false
				for _, ev := range resp.Events {
					var m map[string]interface{}
					_ = json.Unmarshal([]byte(ev.SignedBody), &m)
					taskID, _ := m["id"].(string)
					status, _ := m["status"].(string)
					if forwardTo != "" {
						code, fwdErr := forwardEvent(ctx, hc, forwardTo, ev)
						if fwdErr != nil {
							fmt.Fprintf(c.stderr, "[%d] task=%s status=%s → %s error: %v\n", ev.ID, taskID, status, forwardTo, fwdErr)
							retryEvent = true
							break
						} else {
							fmt.Fprintf(c.stderr, "[%d] task=%s status=%s → %s %d\n", ev.ID, taskID, status, forwardTo, code)
						}
					} else {
						fmt.Fprintf(c.stderr, "[%d] task=%s status=%s\n", ev.ID, taskID, status)
					}
					fmt.Fprintln(c.stdout, ev.SignedBody)
					if err := ackEvent(ctx, hc, baseURL, apiKey, sessionToken, ev.ID); err != nil {
						if errors.Is(err, errSessionExpired) {
							fmt.Fprintf(c.stderr, "[session expired] creating new session...\n")
							sessionToken, err = createListenSession(ctx, hc, baseURL, apiKey)
							if err != nil {
								fmt.Fprintf(c.stderr, "[error] %v, retrying in %v...\n", err, backoff)
							} else {
								lastID = 0
								backoff = time.Second
							}
						} else {
							fmt.Fprintf(c.stderr, "[%d] ack error: %v\n", ev.ID, err)
						}
						retryEvent = true
						break
					}
					lastID = ev.ID
				}
				if retryEvent {
					time.Sleep(1 * time.Second)
					continue
				}
				if len(resp.Events) == 0 {
					time.Sleep(1 * time.Second)
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&forwardTo, "forward-to", "", "Local URL to POST callbacks to")
	return cmd
}

func fetchListenSecret(ctx context.Context, hc *http.Client, baseURL, apiKey string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/api/v1/cli/listen_secret", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("X-API-Key", apiKey)
	req.Header.Set("User-Agent", core.CLIUserAgent(runapi.Version))

	resp, err := hc.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("GET listen_secret: %d %s", resp.StatusCode, body)
	}
	var result listenSecretResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.ListenSecret, nil
}

func createListenSession(ctx context.Context, hc *http.Client, baseURL, apiKey string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+"/api/v1/cli/listen_sessions", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("X-API-Key", apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", core.CLIUserAgent(runapi.Version))

	resp, err := hc.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("POST listen_sessions: %d %s", resp.StatusCode, body)
	}
	var result createSessionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.SessionToken, nil
}

func destroyListenSession(ctx context.Context, hc *http.Client, baseURL, apiKey, sessionToken string) {
	url := fmt.Sprintf("%s/api/v1/cli/listen_sessions/%s", baseURL, sessionToken)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return
	}
	req.Header.Set("X-API-Key", apiKey)
	req.Header.Set("User-Agent", core.CLIUserAgent(runapi.Version))

	resp, err := hc.Do(req)
	if err != nil {
		return
	}
	resp.Body.Close()
}

func pollEvents(ctx context.Context, hc *http.Client, baseURL, apiKey, sessionToken string, lastID int64) (*pollResponse, error) {
	url := fmt.Sprintf("%s/api/v1/cli/listen_events?last_id=%d", baseURL, lastID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-API-Key", apiKey)
	req.Header.Set("X-Session-Token", sessionToken)
	req.Header.Set("User-Agent", core.CLIUserAgent(runapi.Version))

	resp, err := hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, errSessionExpired
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GET listen_events: %d %s", resp.StatusCode, body)
	}
	var result pollResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func ackEvent(ctx context.Context, hc *http.Client, baseURL, apiKey, sessionToken string, eventID int64) error {
	url := fmt.Sprintf("%s/api/v1/cli/listen_events/%d/ack", baseURL, eventID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-API-Key", apiKey)
	req.Header.Set("X-Session-Token", sessionToken)
	req.Header.Set("User-Agent", core.CLIUserAgent(runapi.Version))

	resp, err := hc.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return errSessionExpired
	}
	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("POST listen_events/%d/ack: %d %s", eventID, resp.StatusCode, strings.TrimSpace(string(body)))
	}
	return nil
}

func forwardEvent(ctx context.Context, hc *http.Client, targetURL string, ev listenEvent) (int, error) {
	fwdCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(fwdCtx, http.MethodPost, targetURL, bytes.NewReader([]byte(ev.SignedBody)))
	if err != nil {
		return 0, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range ev.Headers {
		req.Header.Set(k, v)
	}

	resp, err := hc.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		detail := strings.TrimSpace(string(body))
		if detail != "" {
			return resp.StatusCode, fmt.Errorf("local endpoint returned %d: %s", resp.StatusCode, detail)
		}
		return resp.StatusCode, fmt.Errorf("local endpoint returned %d", resp.StatusCode)
	}
	return resp.StatusCode, nil
}

func normalizeForwardURL(raw string) (string, error) {
	target := strings.TrimSpace(raw)
	if target == "" {
		return "", nil
	}
	if !strings.Contains(target, "://") {
		target = "http://" + target
	}

	parsed, err := url.Parse(target)
	if err != nil {
		return "", fmt.Errorf("invalid --forward-to URL: %w", err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return "", fmt.Errorf("--forward-to must use http or https")
	}
	if parsed.Host == "" {
		return "", fmt.Errorf("--forward-to must include a host")
	}
	return parsed.String(), nil
}
