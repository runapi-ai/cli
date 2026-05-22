package main

import (
	"context"
	"io"
	"os"
	"strings"
	"time"

	"github.com/runapi-ai/core-sdk/go/core"
	"github.com/runapi-ai/core-sdk/go/option"
	"github.com/spf13/cobra"
)

func (c *cli) authCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Manage authentication",
	}
	cmd.AddCommand(c.authStatusCommand())
	cmd.AddCommand(c.authImportTokenCommand())
	return cmd
}

func (c *cli) authStatusCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Check authentication status",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.runAuthStatus(cmd)
		},
	}
}

func (c *cli) runAuthStatus(cmd *cobra.Command) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	configPath, _ := configFilePath()

	// Determine auth source
	source := "none"
	apiKey := ""
	if strings.TrimSpace(c.apiKeyFlag) != "" {
		source = "flag"
		apiKey = c.apiKeyFlag
	} else if v := strings.TrimSpace(os.Getenv("RUNAPI_API_KEY")); v != "" {
		source = "env"
		apiKey = v
	} else if strings.TrimSpace(cfg.APIKey) != "" {
		source = "config"
		apiKey = cfg.APIKey
	}

	if apiKey == "" {
		c.logf("No API key found")
		return c.writeJSON(map[string]any{
			"authenticated": false,
			"source":        source,
		})
	}

	// Verify token by calling /api/v1/me
	client, callOpts, ctx, cancel, err := c.clientForCommand(cmd)
	if err != nil {
		return err
	}
	defer cancel()

	info, err := client.Account.Info(ctx, callOpts...)
	if err != nil {
		if core.IsAuthentication(err) {
			c.logf("Token is invalid or expired (source: %s)", source)
			return c.writeJSON(map[string]any{
				"authenticated": false,
				"source":        source,
				"error":         "token expired or revoked",
			})
		}
		return err
	}

	c.logf("✓ Authenticated as %s (source: %s)", info.Email, source)

	result := map[string]any{
		"authenticated": true,
		"source":        source,
		"user": map[string]any{
			"id":    info.ID,
			"name":  info.Name,
			"email": info.Email,
		},
	}
	if source == "config" {
		result["config_path"] = configPath
	}
	return c.writeJSON(result)
}

func (c *cli) authImportTokenCommand() *cobra.Command {
	var token, baseURL string
	var skipVerify bool
	cmd := &cobra.Command{
		Use:   "import-token",
		Short: "Save an API key to the local config (headless setup)",
		Long: "import-token writes an API key to the config file so subsequent commands do not need\n" +
			"RUNAPI_API_KEY in the environment.\n\n" +
			"The token is read from --token, RUNAPI_API_KEY, or stdin (\"--token -\"). Unless\n" +
			"--skip-verify is set, the token is verified by calling /api/v1/me first; an invalid\n" +
			"token does not touch the config.",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.runAuthImportToken(cmd.Context(), token, baseURL, skipVerify)
		},
	}
	cmd.Flags().StringVar(&token, "token", "", `API key value. Use "-" to read from stdin.`)
	cmd.Flags().StringVar(&baseURL, "base-url", "", "Optional base URL to persist alongside the token")
	cmd.Flags().BoolVar(&skipVerify, "skip-verify", false, "Skip the /api/v1/me verification before writing the config")
	return cmd
}

func (c *cli) runAuthImportToken(ctx context.Context, tokenFlag, baseURLFlag string, skipVerify bool) error {
	token, err := c.resolveImportToken(tokenFlag)
	if err != nil {
		return err
	}

	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	baseURL := strings.TrimRight(firstNonEmpty(
		strings.TrimSpace(baseURLFlag),
		strings.TrimSpace(c.baseURLFlag),
		strings.TrimSpace(os.Getenv("RUNAPI_BASE_URL")),
		strings.TrimSpace(cfg.BaseURL),
		core.DefaultBaseURL,
	), "/")

	if !skipVerify {
		if err := c.verifyImportedToken(ctx, token, baseURL); err != nil {
			return err
		}
	}

	cfg.APIKey = token
	if strings.TrimSpace(baseURLFlag) != "" {
		cfg.BaseURL = strings.TrimRight(strings.TrimSpace(baseURLFlag), "/")
	}
	cfg.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	if err := saveConfig(cfg); err != nil {
		return err
	}

	configPath, _ := configFilePath()
	c.logf("✓ token saved to %s", configPath)
	return c.writeJSON(map[string]any{
		"imported":    true,
		"verified":    !skipVerify,
		"config_path": configPath,
		"base_url":    cfg.BaseURL,
	})
}

func (c *cli) resolveImportToken(flag string) (string, error) {
	v := strings.TrimSpace(flag)
	if v == "-" {
		data, err := io.ReadAll(io.LimitReader(c.stdin, 4096))
		if err != nil {
			return "", err
		}
		v = strings.TrimSpace(string(data))
	}
	if v == "" {
		v = strings.TrimSpace(os.Getenv("RUNAPI_API_KEY"))
	}
	if v == "" {
		return "", core.NewError(core.ErrValidation,
			"no token supplied: pass --token, set RUNAPI_API_KEY, or pipe via '--token -'",
			422, "", nil, nil)
	}
	return v, nil
}

func (c *cli) verifyImportedToken(ctx context.Context, token, baseURL string) error {
	opts := []option.ClientOption{
		option.WithAPIKey(token),
		option.WithBaseURL(baseURL),
	}
	client, err := c.newClient(opts...)
	if err != nil {
		return err
	}
	verifyCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if _, err := client.Account.Info(verifyCtx); err != nil {
		return err
	}
	return nil
}
