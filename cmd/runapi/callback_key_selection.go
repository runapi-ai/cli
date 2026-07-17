package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/runapi-ai/core-sdk/go/core"
)

type callbackKeySelection struct {
	ID                  string
	Source              string
	ConfigPath          string
	SaveAfterValidation bool
}

func (c *cli) resolveCallbackAPIKey(
	ctx context.Context,
	client *http.Client,
	baseURL string,
	apiKey string,
	flagValue string,
) (callbackKeySelection, error) {
	configPath, err := projectConfigPath(c.projectDir)
	if err != nil {
		return callbackKeySelection{}, err
	}
	if id := strings.TrimSpace(flagValue); id != "" {
		return callbackKeySelection{ID: id, Source: "flag", ConfigPath: configPath}, nil
	}

	projectCfg, err := loadProjectConfigFile(configPath)
	if err != nil {
		return callbackKeySelection{}, err
	}
	if id := strings.TrimSpace(projectCfg.CallbackAPIKeyID); id != "" {
		return callbackKeySelection{ID: id, Source: "config", ConfigPath: configPath}, nil
	}
	return c.chooseCallbackAPIKey(ctx, client, baseURL, apiKey, configPath)
}

func (c *cli) chooseCallbackAPIKey(
	ctx context.Context,
	client *http.Client,
	baseURL string,
	apiKey string,
	configPath string,
) (callbackKeySelection, error) {
	response, err := fetchAPIKeys(ctx, client, baseURL, apiKey)
	if err != nil {
		return callbackKeySelection{}, err
	}
	enabled := make([]callbackAPIKey, 0, len(response.APIKeys))
	for _, key := range response.APIKeys {
		if key.Enabled {
			enabled = append(enabled, key)
		}
	}
	if c.stdinTTY == nil || c.stderrTTY == nil || !c.stdinTTY() || !c.stderrTTY() {
		details := map[string]any{"api_keys": response.APIKeys}
		if len(enabled) > 0 {
			details["hint"] = fmt.Sprintf(
				"runapi listen <url> --callback-api-key-id %s",
				enabled[0].ID,
			)
		}
		return callbackKeySelection{}, core.NewError(
			core.ErrorCode("callback_api_key_required"),
			"callback API key selection required; pass --callback-api-key-id or add callback_api_key_id to .runapi.toml",
			http.StatusUnprocessableEntity,
			"",
			details,
			nil,
		)
	}
	if len(enabled) == 0 {
		return callbackKeySelection{}, core.NewError(
			core.ErrValidation,
			"no enabled callback API keys are available",
			http.StatusUnprocessableEntity,
			"",
			map[string]any{"api_keys": response.APIKeys},
			nil,
		)
	}

	selected, err := c.selectCallbackAPIKey(enabled)
	if err != nil {
		return callbackKeySelection{}, err
	}
	return callbackKeySelection{
		ID:                  selected.ID,
		Source:              "interactive",
		ConfigPath:          configPath,
		SaveAfterValidation: true,
	}, nil
}

func (c *cli) promptForCallbackAPIKey(keys []callbackAPIKey) (callbackAPIKey, error) {
	stdin, ok := c.stdin.(io.ReadCloser)
	if !ok {
		return callbackAPIKey{}, fmt.Errorf("interactive selection requires a terminal stdin")
	}
	stdout, ok := c.stderr.(io.WriteCloser)
	if !ok {
		return callbackAPIKey{}, fmt.Errorf("interactive selection requires a terminal stderr")
	}

	prompt := promptui.Select{
		Label:  "Select callback API key",
		Items:  keys,
		Size:   min(10, len(keys)),
		Stdin:  stdin,
		Stdout: stdout,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}:",
			Active:   "▸ {{ .Name }}  {{ .ID }}  {{ .MaskedToken }}",
			Inactive: "  {{ .Name }}  {{ .ID }}  {{ .MaskedToken }}",
			Selected: "{{ .Name }}  {{ .ID }}  {{ .MaskedToken }}",
		},
	}
	index, _, err := prompt.Run()
	if err != nil {
		return callbackAPIKey{}, err
	}
	return keys[index], nil
}
