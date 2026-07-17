package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	runapi "github.com/runapi-ai/cli/internal/runapi"
	"github.com/runapi-ai/core-sdk/go/core"
	"github.com/spf13/cobra"
)

type callbackAPIKey struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	MaskedToken string `json:"masked_token"`
	Enabled     bool   `json:"enabled"`
}

type apiKeysResponse struct {
	APIKeys []callbackAPIKey `json:"api_keys"`
}

func (c *cli) apiKeysCommand() *cobra.Command {
	apiKeys := &cobra.Command{
		Use:   "api-keys",
		Short: "Inspect API keys available to CLI listener operations",
		Args:  cobra.NoArgs,
	}
	var outputJSON bool
	list := &cobra.Command{
		Use:   "list",
		Short: "List callback API key candidates",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			apiKey, baseURL, err := c.listenerCredentials()
			if err != nil {
				return err
			}
			response, err := fetchAPIKeys(cmd.Context(), c.listenHTTPClient(), baseURL, apiKey)
			if err != nil {
				return err
			}
			if outputJSON {
				return c.writeJSON(response)
			}

			writer := tabwriter.NewWriter(c.stdout, 0, 4, 2, ' ', 0)
			_, _ = fmt.Fprintln(writer, "NAME\tID\tMASKED KEY\tENABLED")
			for _, key := range response.APIKeys {
				_, _ = fmt.Fprintf(writer, "%s\t%s\t%s\t%t\n", key.Name, key.ID, key.MaskedToken, key.Enabled)
			}
			return writer.Flush()
		},
	}
	list.Flags().BoolVar(&outputJSON, "json", false, "Write the API key list as JSON")
	apiKeys.AddCommand(list)
	return apiKeys
}

func fetchAPIKeys(ctx context.Context, client *http.Client, baseURL, apiKey string) (apiKeysResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, strings.TrimRight(baseURL, "/")+"/api/v1/cli/keys", nil)
	if err != nil {
		return apiKeysResponse{}, err
	}
	req.Header.Set("X-API-Key", apiKey)
	req.Header.Set("User-Agent", core.CLIUserAgent(runapi.Version))

	resp, err := client.Do(req)
	if err != nil {
		return apiKeysResponse{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return apiKeysResponse{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return apiKeysResponse{}, cliAPIError(resp, body)
	}

	var response apiKeysResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return apiKeysResponse{}, err
	}
	return response, nil
}

func (c *cli) listenerCredentials() (string, string, error) {
	cfg, err := loadConfig()
	if err != nil {
		return "", "", err
	}
	apiKey := firstNonEmpty(
		strings.TrimSpace(c.apiKeyFlag),
		strings.TrimSpace(os.Getenv("RUNAPI_API_KEY")),
		strings.TrimSpace(cfg.APIKey),
	)
	if apiKey == "" {
		return "", "", core.NewError(
			core.ErrAuthentication,
			"API key required (--api-key, RUNAPI_API_KEY, or runapi login)",
			http.StatusUnauthorized,
			"",
			nil,
			nil,
		)
	}
	baseURL := strings.TrimRight(firstNonEmpty(
		strings.TrimSpace(c.baseURLFlag),
		strings.TrimSpace(os.Getenv("RUNAPI_BASE_URL")),
		strings.TrimSpace(cfg.BaseURL),
		core.DefaultBaseURL,
	), "/")
	return apiKey, baseURL, nil
}

func (c *cli) listenHTTPClient() *http.Client {
	if c.httpClient != nil {
		return c.httpClient
	}
	return &http.Client{Timeout: 30 * time.Second}
}

func cliAPIError(response *http.Response, body []byte) error {
	err := core.ErrorFromResponse(response, body)
	apiErr, ok := err.(*core.Error)
	if !ok {
		return err
	}
	var payload struct {
		Code string `json:"code"`
	}
	if json.Unmarshal(body, &payload) == nil && payload.Code != "" {
		apiErr.Code = core.ErrorCode(payload.Code)
	}
	return apiErr
}
