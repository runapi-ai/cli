package main

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/runapi-ai/core-sdk/go/core"
	"github.com/spf13/cobra"
)

func (c *cli) loginCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "login",
		Short: "Authenticate with RunAPI via browser",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.runLogin(cmd)
		},
	}
}

func (c *cli) runLogin(cmd *cobra.Command) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	baseURL := strings.TrimRight(firstNonEmpty(
		strings.TrimSpace(c.baseURLFlag),
		strings.TrimSpace(os.Getenv("RUNAPI_BASE_URL")),
		strings.TrimSpace(cfg.BaseURL),
		core.DefaultBaseURL,
	), "/")

	state, err := generateState()
	if err != nil {
		return err
	}
	verifier, challenge, err := generatePKCE()
	if err != nil {
		return err
	}
	displayCode := generateDisplayCode()

	// Start local callback server
	listener, err := startListener()
	if err != nil {
		return err
	}
	defer listener.Close()

	port := listener.Addr().(*net.TCPAddr).Port
	hostname, _ := os.Hostname()
	if hostname == "" {
		hostname = "CLI"
	}

	// Build authorization URL
	authURL := fmt.Sprintf("%s/cli/authorize?state=%s&code_challenge=%s&code_challenge_method=S256&redirect_port=%d&display_code=%s&hostname=%s",
		baseURL,
		url.QueryEscape(state), url.QueryEscape(challenge), port,
		url.QueryEscape(displayCode), url.QueryEscape(hostname))

	c.logf("Opening browser for authorization...")
	c.logf("Your pairing code is: %s", displayCode)
	c.logf("")
	c.logf("Confirm this code matches in your browser before authorizing.")
	if err := openBrowser(authURL); err != nil {
		c.logf("If browser doesn't open, visit: %s", authURL)
	}
	c.logf("Waiting for authorization...")

	// Wait for callback
	ctx, cancel := context.WithTimeout(cmd.Context(), 5*time.Minute)
	defer cancel()

	result, err := waitForCallback(ctx, listener, state, baseURL)
	if err != nil {
		return err
	}

	// Exchange code for token
	exchangeURL := fmt.Sprintf("%s/api/v1/cli/exchange", baseURL)

	tokenResp, err := exchangeCode(ctx, exchangeURL, result.code, verifier, port, hostname)
	if err != nil {
		return err
	}

	// Save to config — preserve existing fields, only update api_key and created_at
	cfg.APIKey = tokenResp.Key
	cfg.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	if err := saveConfig(cfg); err != nil {
		return err
	}

	configPath, _ := configFilePath()
	c.logf("✓ Authenticated as %s", tokenResp.User.Email)
	c.logf("Token saved to %s", configPath)

	return c.writeJSON(map[string]any{
		"authenticated": true,
		"user":          map[string]string{"email": tokenResp.User.Email},
		"config_path":   configPath,
	})
}

type callbackResult struct {
	code string
}

func waitForCallback(ctx context.Context, listener net.Listener, expectedState, corsOrigin string) (*callbackResult, error) {
	resultCh := make(chan *callbackResult, 1)
	errCh := make(chan error, 1)

	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		// Allow cross-origin fetch from Rails callback page
		w.Header().Set("Access-Control-Allow-Origin", corsOrigin)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if errParam := r.URL.Query().Get("error"); errParam != "" {
			w.WriteHeader(http.StatusOK)
			errCh <- core.NewError(core.ErrAuthentication, "Authorization cancelled by user", 0, "", nil, nil)
			return
		}

		code := r.URL.Query().Get("code")
		state := r.URL.Query().Get("state")

		if state != expectedState {
			w.WriteHeader(http.StatusBadRequest)
			errCh <- core.NewError(core.ErrAuthentication, "State mismatch", 0, "", nil, nil)
			return
		}

		w.WriteHeader(http.StatusOK)
		resultCh <- &callbackResult{code: code}
	})

	server := &http.Server{Handler: mux}
	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	var result *callbackResult
	var callbackErr error
	select {
	case result = <-resultCh:
	case callbackErr = <-errCh:
	case <-ctx.Done():
		callbackErr = core.NewError(core.ErrTimeout, "Authorization timed out (5 minutes)", 0, "", nil, ctx.Err())
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer shutdownCancel()
	server.Shutdown(shutdownCtx)

	return result, callbackErr
}

type exchangeResponse struct {
	Key       string `json:"key"`
	TokenType string `json:"token_type"`
	Name      string `json:"name"`
	User      struct {
		Email string `json:"email"`
	} `json:"user"`
	Error string `json:"error"`
}

func exchangeCode(ctx context.Context, exchangeURL, code, verifier string, port int, hostname string) (*exchangeResponse, error) {
	body, _ := json.Marshal(map[string]any{
		"code":          code,
		"code_verifier": verifier,
		"redirect_port": port,
	})

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, exchangeURL, strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-CLI-Hostname", hostname)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, core.NewError(core.ErrNetwork, fmt.Sprintf("exchange request failed: %v", err), 0, "", nil, err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		var errResp exchangeResponse
		json.Unmarshal(respBody, &errResp)
		msg := errResp.Error
		if msg == "" {
			msg = fmt.Sprintf("exchange failed with status %d", resp.StatusCode)
		}
		return nil, core.NewError(core.ErrAuthentication, msg, resp.StatusCode, "", nil, nil)
	}

	var tokenResp exchangeResponse
	if err := json.Unmarshal(respBody, &tokenResp); err != nil {
		return nil, err
	}
	return &tokenResp, nil
}

func generateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func generatePKCE() (verifier, challenge string, err error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", "", err
	}
	verifier = base64.RawURLEncoding.EncodeToString(b)
	h := sha256.Sum256([]byte(verifier))
	challenge = base64.RawURLEncoding.EncodeToString(h[:])
	return verifier, challenge, nil
}

// pairingWords is a curated word list for human-readable pairing codes.
// 64 short, distinct, easy-to-read English words.
var pairingWords = []string{
	"apple", "brave", "cloud", "delta", "eagle", "flame", "grace", "haven",
	"ivory", "jewel", "knack", "lemon", "maple", "noble", "ocean", "pearl",
	"quest", "ridge", "solar", "tiger", "ultra", "vivid", "waltz", "xenon",
	"yacht", "zephyr", "amber", "blaze", "cedar", "dusk", "ember", "frost",
	"gleam", "heron", "iris", "jade", "kite", "lunar", "mango", "nexus",
	"opal", "prism", "quilt", "river", "spark", "thorn", "unity", "vault",
	"wings", "pixel", "aspen", "birch", "coral", "drift", "epoch", "forge",
	"grove", "haze", "inlet", "jasper", "kelp", "lotus", "marsh", "north",
}

func generateDisplayCode() string {
	words := make([]string, 4)
	for i := range words {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(pairingWords))))
		words[i] = pairingWords[n.Int64()]
	}
	return strings.Join(words, "-")
}

func openBrowser(url string) error {
	switch runtime.GOOS {
	case "darwin":
		return exec.Command("open", url).Start()
	case "linux":
		return exec.Command("xdg-open", url).Start()
	case "windows":
		return exec.Command("cmd", "/c", "start", "", url).Start()
	default:
		return fmt.Errorf("unsupported platform %s", runtime.GOOS)
	}
}

func startListener() (net.Listener, error) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, core.NewError(core.ErrNetwork, "failed to start local callback server", 0, "", nil, err)
	}
	return ln, nil
}
