package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	runapi "github.com/runapi-ai/cli/internal/runapi"
	"github.com/runapi-ai/core-sdk/go/core"
	"github.com/spf13/cobra"
)

var (
	errSessionExpired         = errors.New("session expired")
	errCallbackAPIKeyUnusable = errors.New("callback API key is no longer usable")
	errListenSecretChanged    = errors.New("listen signing secret changed")
)

const (
	listenerEmptyPollInterval    = 15 * time.Second
	listenerMaxEmptyPollInterval = 30 * time.Second
	listenerPollJitterFraction   = 0.10
)

var randomListenerPollJitter = rand.Float64
var waitForListenerContext = waitForContext

type listenSecretResponse struct {
	ListenSecret   string         `json:"listen_secret"`
	CallbackAPIKey callbackAPIKey `json:"callback_api_key"`
}

type createSessionResponse struct {
	SessionToken   string         `json:"session_token"`
	ListenSecret   string         `json:"listen_secret"`
	CallbackAPIKey callbackAPIKey `json:"callback_api_key"`
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
	var forwardTo, callbackAPIKeyID string
	var printSecret, rotateSecret bool
	cmd := &cobra.Command{
		Use:   "listen [url]",
		Short: "Receive RunAPI callbacks locally",
		Long:  "Forward task callbacks to a local URL. Tasks with a callback_url are also copied to the listener while still being delivered to that URL.",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				if strings.TrimSpace(forwardTo) != "" {
					return core.NewError(
						core.ErrValidation,
						"listen URL and --forward-to are mutually exclusive",
						http.StatusUnprocessableEntity,
						"",
						nil,
						nil,
					)
				}
				forwardTo = args[0]
			}
			forwardTo, err := normalizeForwardURL(forwardTo)
			if err != nil {
				return err
			}
			apiKey, baseURL, err := c.listenerCredentials()
			if err != nil {
				return err
			}

			hc := c.listenHTTPClient()
			ctx, cancel := context.WithCancel(cmd.Context())
			defer cancel()

			selection, err := c.resolveCallbackAPIKey(ctx, hc, baseURL, apiKey, callbackAPIKeyID)
			if err != nil {
				return err
			}
			if printSecret {
				secret, err := c.fetchSecretResolvingStaleConfig(ctx, hc, baseURL, apiKey, &selection)
				if err != nil {
					return err
				}
				_, err = fmt.Fprintln(c.stdout, secret.ListenSecret)
				return err
			}
			if rotateSecret {
				secret, err := c.rotateSecretResolvingStaleConfig(ctx, hc, baseURL, apiKey, &selection)
				if err != nil {
					return err
				}
				_, err = fmt.Fprintln(c.stdout, secret.ListenSecret)
				return err
			}

			session, err := c.createSessionResolvingStaleConfig(ctx, hc, baseURL, apiKey, &selection)
			if err != nil {
				return fmt.Errorf("failed to create listen session: %w", err)
			}
			sessionToken := session.SessionToken
			listenSecret := session.ListenSecret
			defer func() {
				destroyListenSession(context.Background(), hc, baseURL, apiKey, sessionToken)
			}()
			if selection.SaveAfterValidation {
				if err := saveProjectConfig(selection.ConfigPath, projectConfig{CallbackAPIKeyID: session.CallbackAPIKey.ID}); err != nil {
					return err
				}
			}
			fmt.Fprintf(
				c.stderr,
				"Callback API key: %s (%s, %s)\n",
				session.CallbackAPIKey.Name,
				session.CallbackAPIKey.ID,
				session.CallbackAPIKey.MaskedToken,
			)
			if selection.ConfigPath != "" {
				fmt.Fprintf(c.stderr, "Project config: %s\n", selection.ConfigPath)
			}
			fmt.Fprintf(c.stderr, "Ready! Listen Signing Secret: %s\n", session.ListenSecret)
			if forwardTo != "" {
				fmt.Fprintf(c.stderr, "Forwarding to %s\n", forwardTo)
			}

			var lastID int64
			backoff := time.Second
			consecutiveEmptyPolls := 0
			for ctx.Err() == nil {
				resp, err := pollEvents(ctx, hc, baseURL, apiKey, sessionToken, lastID)
				if err != nil {
					if isTerminalListenerError(err) {
						return err
					}
					if errors.Is(err, errSessionExpired) {
						fmt.Fprintf(c.stderr, "[session expired] creating new session...\n")
						newToken, retriable, recreateErr := c.recreateSession(
							ctx, hc, baseURL, apiKey, selection.ID, listenSecret,
						)
						if recreateErr != nil {
							if !retriable {
								return recreateErr
							}
							retryWait := listenerRetryWait(recreateErr, listenerWaitDuration(backoff))
							fmt.Fprintf(c.stderr, "[error] %v, retrying in %v...\n", recreateErr, retryWait)
							if !waitForListenerContext(ctx, retryWait) {
								return nil
							}
							backoff = nextBackoff(backoff)
							continue
						}
						sessionToken = newToken
						lastID = 0
						backoff = time.Second
						continue
					}
					retryWait := listenerRetryWait(err, listenerWaitDuration(backoff))
					fmt.Fprintf(c.stderr, "[error] %v, retrying in %v...\n", err, retryWait)
					if !waitForListenerContext(ctx, retryWait) {
						return nil
					}
					backoff = nextBackoff(backoff)
					continue
				}
				backoff = time.Second
				retryEvent := false
				retryEventWait := time.Second
				for _, ev := range resp.Events {
					var m map[string]interface{}
					_ = json.Unmarshal([]byte(ev.SignedBody), &m)
					taskID, _ := m["id"].(string)
					status, _ := m["status"].(string)
					if err := ackEvent(ctx, hc, baseURL, apiKey, sessionToken, ev.ID); err != nil {
						if isTerminalListenerError(err) {
							return err
						}
						if errors.Is(err, errSessionExpired) {
							fmt.Fprintf(c.stderr, "[session expired] creating new session...\n")
							newToken, retriable, recreateErr := c.recreateSession(
								ctx, hc, baseURL, apiKey, selection.ID, listenSecret,
							)
							if recreateErr != nil {
								if !retriable {
									return recreateErr
								}
								retryEventWait = listenerRetryWait(recreateErr, listenerWaitDuration(backoff))
								fmt.Fprintf(c.stderr, "[error] %v, retrying in %v...\n", recreateErr, retryEventWait)
							} else {
								sessionToken = newToken
								lastID = 0
								backoff = time.Second
							}
						} else {
							fmt.Fprintf(c.stderr, "[%d] ack error: %v\n", ev.ID, err)
							retryEventWait = listenerRetryWait(err, listenerWaitDuration(time.Second))
						}
						retryEvent = true
						break
					}
					lastID = ev.ID
					if forwardTo != "" {
						code, fwdErr := forwardEvent(ctx, hc, forwardTo, ev)
						if fwdErr != nil {
							if code != 0 {
								fmt.Fprintf(c.stderr, "[%d] task=%s status=%s → %s %d error: %v\n", ev.ID, taskID, status, forwardTo, code, fwdErr)
							} else {
								fmt.Fprintf(c.stderr, "[%d] task=%s status=%s → %s error: %v\n", ev.ID, taskID, status, forwardTo, fwdErr)
							}
						} else {
							fmt.Fprintf(c.stderr, "[%d] task=%s status=%s → %s %d\n", ev.ID, taskID, status, forwardTo, code)
						}
					} else {
						fmt.Fprintf(c.stderr, "[%d] task=%s status=%s\n", ev.ID, taskID, status)
					}
					fmt.Fprintln(c.stdout, ev.SignedBody)
				}
				if retryEvent {
					consecutiveEmptyPolls = 0
					if !waitForListenerContext(ctx, retryEventWait) {
						return nil
					}
					continue
				}
				if len(resp.Events) == 0 {
					consecutiveEmptyPolls++
					if !waitForListenerContext(ctx, listenerPollWait(consecutiveEmptyPolls)) {
						return nil
					}
				} else {
					// A non-empty response is drained immediately on the next
					// request, so event delivery never inherits idle backoff.
					consecutiveEmptyPolls = 0
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&forwardTo, "forward-to", "", "Local URL to POST callbacks to")
	cmd.Flags().StringVar(&callbackAPIKeyID, "callback-api-key-id", "", "Stable ID of the API key whose task callbacks should be received")
	cmd.Flags().BoolVar(&printSecret, "print-secret", false, "Print the selected API key's Listen Signing Secret and exit")
	cmd.Flags().BoolVar(&rotateSecret, "rotate-secret", false, "Rotate the selected API key's Listen Signing Secret, invalidate its listeners, print the new secret, and exit")
	cmd.MarkFlagsMutuallyExclusive("print-secret", "rotate-secret")
	return cmd
}

func waitForContext(ctx context.Context, duration time.Duration) bool {
	timer := time.NewTimer(duration)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
}

func listenerPollWait(consecutiveEmptyPolls int) time.Duration {
	base := listenerEmptyPollInterval
	if consecutiveEmptyPolls > 1 {
		base = listenerMaxEmptyPollInterval
	}
	return listenerWaitDuration(base)
}

func listenerWaitDuration(base time.Duration) time.Duration {
	jitter := (randomListenerPollJitter()*2 - 1) * listenerPollJitterFraction
	return time.Duration(float64(base) * (1 + jitter))
}

func listenerRetryWait(err error, fallback time.Duration) time.Duration {
	apiErr, ok := errors.AsType[*core.Error](err)
	if ok && core.IsRateLimit(apiErr) && apiErr.RetryAfter > 0 {
		return apiErr.RetryAfter
	}
	return fallback
}

func isListenerAccessRevoked(err error) bool {
	apiErr, ok := errors.AsType[*core.Error](err)
	return ok && (apiErr.Status == http.StatusUnauthorized || apiErr.Status == http.StatusForbidden)
}

// isTerminalListenerError reports errors that must stop the listener outright
// rather than trigger a session retry: the selected callback key became unusable
// (410 Gone) or the CLI credential's listen access was revoked (401/403).
func isTerminalListenerError(err error) bool {
	return errors.Is(err, errCallbackAPIKeyUnusable) || isListenerAccessRevoked(err)
}

// nextBackoff doubles the retry delay up to a 30s ceiling.
func nextBackoff(current time.Duration) time.Duration {
	if current < 30*time.Second {
		return current * 2
	}
	return current
}

// createVerifiedSession creates a listener session for keyID and asserts the
// server echoed that exact callback key back, guarding against a server that
// silently resolved a different one.
func createVerifiedSession(ctx context.Context, hc *http.Client, baseURL, apiKey, keyID string) (*createSessionResponse, error) {
	session, err := createListenSession(ctx, hc, baseURL, apiKey, keyID)
	if err != nil {
		return nil, err
	}
	if session.CallbackAPIKey.ID != keyID {
		return nil, fmt.Errorf("listen session resolved unexpected callback API key %q", session.CallbackAPIKey.ID)
	}
	return session, nil
}

// fetchVerifiedSecret fetches the listen secret for keyID and asserts the server
// echoed that exact callback key back.
func fetchVerifiedSecret(ctx context.Context, hc *http.Client, baseURL, apiKey, keyID string) (*listenSecretResponse, error) {
	secret, err := fetchListenSecret(ctx, hc, baseURL, apiKey, keyID)
	if err != nil {
		return nil, err
	}
	if secret.CallbackAPIKey.ID != keyID {
		return nil, fmt.Errorf("listen secret resolved unexpected callback API key %q", secret.CallbackAPIKey.ID)
	}
	return secret, nil
}

// rotateVerifiedSecret rotates the listen secret for keyID and asserts the
// server echoed that exact callback key back.
func rotateVerifiedSecret(ctx context.Context, hc *http.Client, baseURL, apiKey, keyID string) (*listenSecretResponse, error) {
	secret, err := rotateListenSecret(ctx, hc, baseURL, apiKey, keyID)
	if err != nil {
		return nil, err
	}
	if secret.CallbackAPIKey.ID != keyID {
		return nil, fmt.Errorf("listen secret rotation resolved unexpected callback API key %q", secret.CallbackAPIKey.ID)
	}
	return secret, nil
}

// recreateSession makes a fresh session after the current one expired mid-run,
// shared by the poll and ack loops. Deliberately unlike startup, a committed key
// that has gone missing is NOT re-prompted (that would silently switch keys under
// a long-running listener); it returns a non-retriable error so the caller exits.
// retriable reports whether the caller should back off and try again instead.
func (c *cli) recreateSession(ctx context.Context, hc *http.Client, baseURL, apiKey, keyID, listenSecret string) (token string, retriable bool, err error) {
	session, createErr := createVerifiedSession(ctx, hc, baseURL, apiKey, keyID)
	if createErr != nil {
		if core.IsNotFound(createErr) {
			return "", false, fmt.Errorf("callback API key is no longer available: %w", createErr)
		}
		if isTerminalListenerError(createErr) {
			return "", false, createErr
		}
		return "", true, createErr
	}
	if session.ListenSecret != listenSecret {
		destroyListenSession(context.Background(), hc, baseURL, apiKey, session.SessionToken)
		return "", false, fmt.Errorf("%w; update the local verifier and restart the listener", errListenSecretChanged)
	}
	return session.SessionToken, false, nil
}

// maybeReselectStaleConfig re-runs interactive selection once when the key that
// was committed to .runapi.toml (Source == "config") is the one reported missing.
// It returns retry=true (with *selection replaced) when a fresh key was picked,
// so the caller can retry its operation. An explicit --callback-api-key-id flag
// or an already-interactive pick is never rewritten. Only used at startup: a
// long-running listener deliberately does not re-prompt mid-stream (see the poll
// loop) so it never silently switches keys under the user.
func (c *cli) maybeReselectStaleConfig(ctx context.Context, hc *http.Client, baseURL, apiKey string, selection *callbackKeySelection, cause error) (retry bool, err error) {
	if cause == nil || selection.Source != "config" || !core.IsNotFound(cause) {
		return false, nil
	}
	newSelection, chooseErr := c.chooseCallbackAPIKey(ctx, hc, baseURL, apiKey, selection.ConfigPath)
	if chooseErr != nil {
		return false, chooseErr
	}
	*selection = newSelection
	return true, nil
}

// createSessionResolvingStaleConfig creates a verified session, reselecting once
// if a committed config key has gone missing.
func (c *cli) createSessionResolvingStaleConfig(ctx context.Context, hc *http.Client, baseURL, apiKey string, selection *callbackKeySelection) (*createSessionResponse, error) {
	session, err := createVerifiedSession(ctx, hc, baseURL, apiKey, selection.ID)
	if retry, rerr := c.maybeReselectStaleConfig(ctx, hc, baseURL, apiKey, selection, err); rerr != nil {
		return nil, rerr
	} else if retry {
		return createVerifiedSession(ctx, hc, baseURL, apiKey, selection.ID)
	}
	return session, err
}

// fetchSecretResolvingStaleConfig fetches a verified secret, reselecting once if
// a committed config key has gone missing.
func (c *cli) fetchSecretResolvingStaleConfig(ctx context.Context, hc *http.Client, baseURL, apiKey string, selection *callbackKeySelection) (*listenSecretResponse, error) {
	secret, err := fetchVerifiedSecret(ctx, hc, baseURL, apiKey, selection.ID)
	if retry, rerr := c.maybeReselectStaleConfig(ctx, hc, baseURL, apiKey, selection, err); rerr != nil {
		return nil, rerr
	} else if retry {
		return fetchVerifiedSecret(ctx, hc, baseURL, apiKey, selection.ID)
	}
	return secret, err
}

// rotateSecretResolvingStaleConfig rotates a verified secret, reselecting once
// if a committed config key has gone missing.
func (c *cli) rotateSecretResolvingStaleConfig(ctx context.Context, hc *http.Client, baseURL, apiKey string, selection *callbackKeySelection) (*listenSecretResponse, error) {
	secret, err := rotateVerifiedSecret(ctx, hc, baseURL, apiKey, selection.ID)
	if retry, rerr := c.maybeReselectStaleConfig(ctx, hc, baseURL, apiKey, selection, err); rerr != nil {
		return nil, rerr
	} else if retry {
		return rotateVerifiedSecret(ctx, hc, baseURL, apiKey, selection.ID)
	}
	return secret, err
}

// callbackKeyUnusableError wraps a 410 Gone response as errCallbackAPIKeyUnusable.
func callbackKeyUnusableError(resp *http.Response, body []byte) error {
	return fmt.Errorf("%w: %v", errCallbackAPIKeyUnusable, cliAPIError(resp, body))
}

func fetchListenSecret(ctx context.Context, hc *http.Client, baseURL, apiKey, callbackAPIKeyID string) (*listenSecretResponse, error) {
	return requestListenSecret(ctx, hc, http.MethodGet, baseURL, apiKey, callbackAPIKeyID)
}

func rotateListenSecret(ctx context.Context, hc *http.Client, baseURL, apiKey, callbackAPIKeyID string) (*listenSecretResponse, error) {
	return requestListenSecret(ctx, hc, http.MethodPatch, baseURL, apiKey, callbackAPIKeyID)
}

func requestListenSecret(ctx context.Context, hc *http.Client, method, baseURL, apiKey, callbackAPIKeyID string) (*listenSecretResponse, error) {
	endpoint, err := url.Parse(baseURL + "/api/v1/cli/listen_secret")
	if err != nil {
		return nil, err
	}
	query := endpoint.Query()
	query.Set("callback_api_key_id", callbackAPIKeyID)
	endpoint.RawQuery = query.Encode()
	req, err := http.NewRequestWithContext(ctx, method, endpoint.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-API-Key", apiKey)
	req.Header.Set("User-Agent", core.CLIUserAgent(runapi.Version))

	resp, err := hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, cliAPIError(resp, body)
	}
	var result listenSecretResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func createListenSession(ctx context.Context, hc *http.Client, baseURL, apiKey, callbackAPIKeyID string) (*createSessionResponse, error) {
	body, err := json.Marshal(struct {
		CallbackAPIKeyID string `json:"callback_api_key_id,omitempty"`
	}{CallbackAPIKeyID: callbackAPIKeyID})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+"/api/v1/cli/listen_sessions", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-API-Key", apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", core.CLIUserAgent(runapi.Version))

	resp, err := hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, cliAPIError(resp, body)
	}
	var result createSessionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
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
	if resp.StatusCode == http.StatusGone {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return nil, callbackKeyUnusableError(resp, body)
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return nil, cliAPIError(resp, body)
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
	if resp.StatusCode == http.StatusGone {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return callbackKeyUnusableError(resp, body)
	}
	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return cliAPIError(resp, body)
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

	forwardingClient := *hc
	forwardingClient.CheckRedirect = func(_ *http.Request, _ []*http.Request) error {
		return http.ErrUseLastResponse
	}
	resp, err := forwardingClient.Do(req)
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
