package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"time"

	runapi "github.com/runapi-ai/cli/internal/runapi"
	"github.com/runapi-ai/core-sdk/go/core"
	"github.com/runapi-ai/core-sdk/go/files"
	"github.com/runapi-ai/core-sdk/go/option"
	"github.com/runapi-ai/elevenlabs-sdk/go/elevenlabs"
	"github.com/runapi-ai/flux-2-sdk/go/flux2"
	"github.com/runapi-ai/flux-kontext-sdk/go/fluxkontext"
	"github.com/runapi-ai/gemini-omni-sdk/go/geminiomni"
	"github.com/runapi-ai/gpt-4o-image-sdk/go/gpt4oimage"
	"github.com/runapi-ai/gpt-image-2-sdk/go/gptimage2"
	"github.com/runapi-ai/gpt-image-sdk/go/gptimage"
	"github.com/runapi-ai/grok-imagine-sdk/go/grokimagine"
	"github.com/runapi-ai/hailuo-sdk/go/hailuo"
	"github.com/runapi-ai/happyhorse-sdk/go/happyhorse"
	"github.com/runapi-ai/ideogram-v3-sdk/go/ideogramv3"
	"github.com/runapi-ai/imagen-4-sdk/go/imagen4"
	"github.com/runapi-ai/infinitetalk-sdk/go/infinitetalk"
	"github.com/runapi-ai/kling-sdk/go/kling"
	"github.com/runapi-ai/luma-sdk/go/luma"
	"github.com/runapi-ai/nano-banana-sdk/go/nanobanana"
	"github.com/runapi-ai/omnihuman-sdk/go/omnihuman"
	"github.com/runapi-ai/qwen-2-sdk/go/qwen2"
	"github.com/runapi-ai/recraft-sdk/go/recraft"
	"github.com/runapi-ai/runway-aleph-sdk/go/runwayaleph"
	"github.com/runapi-ai/runway-sdk/go/runway"
	"github.com/runapi-ai/seedance-sdk/go/seedance"
	"github.com/runapi-ai/seedream-sdk/go/seedream"
	"github.com/runapi-ai/suno-sdk/go/suno"
	"github.com/runapi-ai/topaz-sdk/go/topaz"
	"github.com/runapi-ai/veo-3.1-sdk/go/veo31"
	volcenginelipsync "github.com/runapi-ai/volcengine-lip-sync-sdk/go/volcenginelipsync"
	"github.com/runapi-ai/wan-sdk/go/wan"
	"github.com/runapi-ai/z-image-sdk/go/zimage"
	"github.com/spf13/cobra"
)

type configFile struct {
	APIKey    string `json:"api_key,omitempty"`
	BaseURL   string `json:"base_url,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
}

type cli struct {
	stdout io.Writer
	stderr io.Writer
	stdin  io.Reader

	apiKeyFlag   string
	baseURLFlag  string
	timeout      time.Duration
	async        bool
	pollInterval time.Duration
	quiet        bool
	newClient    func(...option.ClientOption) (*runapi.Client, error)

	archiveBaseURL string
	httpClient     *http.Client
}

type actionSpec struct {
	service     string
	action      string
	isAsync     bool
	inputFields string
	decode      func([]byte) (any, error)
	create      func(context.Context, *runapi.Client, any, []option.RequestOption) (*core.TaskCreateResponse, error)
	run         func(context.Context, *runapi.Client, any, []option.RequestOption) (any, error)
	get         func(context.Context, *runapi.Client, string, []option.RequestOption) (core.TaskResponse, error)
}

func newCLI() *cli {
	return &cli{
		stdout:    os.Stdout,
		stderr:    os.Stderr,
		stdin:     os.Stdin,
		newClient: runapi.NewClient,
	}
}

func (c *cli) run(args []string) int {
	cmd := c.command()
	cmd.SetArgs(args)
	if err := cmd.Execute(); err != nil {
		c.printError(err)
		return exitCode(err)
	}
	return 0
}

func (c *cli) command() *cobra.Command {
	root := &cobra.Command{
		Use:           "runapi",
		Short:         "RunAPI command-line client for typed SDK-backed API calls",
		Long:          "RunAPI CLI is Agent first: a JSON-first command-line client for account, Suno, Veo 3.1, and Nano Banana operations. It reads credentials from flags, environment variables, or config, writes data to stdout, and logs progress to stderr.",
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	root.SetOut(c.stdout)
	root.SetErr(c.stderr)
	root.PersistentFlags().StringVar(&c.apiKeyFlag, "api-key", "", "API key. Overrides RUNAPI_API_KEY.")
	root.PersistentFlags().StringVar(&c.baseURLFlag, "base-url", "", "API base URL. Overrides RUNAPI_BASE_URL, then config, then defaults to https://runapi.ai.")
	root.PersistentFlags().DurationVar(&c.timeout, "timeout", core.DefaultTimeout, "Overall command timeout. Also used as the default per-request timeout and max wait.")
	root.PersistentFlags().BoolVar(&c.async, "async", false, "For async actions, submit the task and return immediately instead of polling until completion.")
	root.PersistentFlags().DurationVar(&c.pollInterval, "poll-interval", 3*time.Second, "Polling interval for async actions and the wait command.")
	root.PersistentFlags().BoolVar(&c.quiet, "quiet", false, "Silence stderr progress logs. JSON responses are still written to stdout.")

	root.AddCommand(c.loginCommand())
	root.AddCommand(c.logoutCommand())
	root.AddCommand(c.authCommand())
	root.AddCommand(c.versionCommand())
	root.AddCommand(c.accountCommand())
	root.AddCommand(c.filesCommand())
	root.AddCommand(c.serviceCommand("suno"))
	root.AddCommand(c.serviceCommand("veo-3-1"))
	root.AddCommand(c.serviceCommand("nano-banana"))
	root.AddCommand(c.serviceCommand("imagen-4"))
	root.AddCommand(c.serviceCommand("seedance"))
	root.AddCommand(c.serviceCommand("seedream"))
	root.AddCommand(c.serviceCommand("runway"))
	root.AddCommand(c.serviceCommand("runway-aleph"))
	root.AddCommand(c.serviceCommand("kling"))
	root.AddCommand(c.serviceCommand("flux-kontext"))
	root.AddCommand(c.serviceCommand("flux-2"))
	root.AddCommand(c.serviceCommand("gemini-omni"))
	root.AddCommand(c.serviceCommand("qwen-2"))
	root.AddCommand(c.serviceCommand("recraft"))
	root.AddCommand(c.serviceCommand("z-image"))
	root.AddCommand(c.serviceCommand("ideogram-v3"))
	root.AddCommand(c.serviceCommand("elevenlabs"))
	root.AddCommand(c.serviceCommand("infinitetalk"))
	root.AddCommand(c.serviceCommand("omnihuman"))
	root.AddCommand(c.serviceCommand("wan"))
	root.AddCommand(c.serviceCommand("luma"))
	root.AddCommand(c.serviceCommand("hailuo"))
	root.AddCommand(c.serviceCommand("volcengine-lip-sync"))
	root.AddCommand(c.serviceCommand("happyhorse"))
	root.AddCommand(c.serviceCommand("gpt-image"))
	root.AddCommand(c.serviceCommand("gpt-image-2"))
	root.AddCommand(c.serviceCommand("gpt-4o-image"))
	root.AddCommand(c.serviceCommand("grok-imagine"))
	root.AddCommand(c.serviceCommand("topaz"))
	root.AddCommand(c.getCommand())
	root.AddCommand(c.waitCommand())
	root.AddCommand(c.listenCommand())
	root.AddCommand(c.agentCommand())
	return root
}

func (c *cli) filesCommand() *cobra.Command {
	var sourceURL, base64Data, fileName string
	var urlOnly bool
	filesCmd := &cobra.Command{Use: "files", Short: "Temporary file upload operations", Args: cobra.NoArgs}
	create := &cobra.Command{
		Use:   "create [path]",
		Short: "Upload a temporary file",
		Long:  "Upload a temporary file from a local path, a URL, or base64 data. The returned URL expires after one hour.",
		Args: func(cmd *cobra.Command, args []string) error {
			sourceCount := 0
			if len(args) > 0 {
				sourceCount++
			}
			if strings.TrimSpace(sourceURL) != "" {
				sourceCount++
			}
			if strings.TrimSpace(base64Data) != "" {
				sourceCount++
			}
			if sourceCount != 1 {
				return core.NewError(core.ErrValidation, "exactly one source is required: path, --url, or --base64", 422, "", nil, nil)
			}
			if len(args) > 1 {
				return core.NewError(core.ErrValidation, "only one file path is supported", 422, "", nil, nil)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			client, callOpts, ctx, cancel, err := c.clientForCommand(cmd)
			if err != nil {
				return err
			}
			defer cancel()

			params := files.CreateParams{FileName: fileName}
			switch {
			case len(args) == 1:
				params.File = args[0]
			case strings.TrimSpace(sourceURL) != "":
				params.Source = files.Source{Type: "url", URL: sourceURL}
			default:
				params.Source = files.Source{Type: "base64", Data: base64Data}
			}

			response, err := client.Files.Create(ctx, params, callOpts...)
			if err != nil {
				return err
			}
			if urlOnly {
				_, err = fmt.Fprintln(c.stdout, response.URL)
				return err
			}
			return c.writeJSON(response)
		},
	}
	create.Flags().StringVar(&sourceURL, "url", "", "Remote file URL source.")
	create.Flags().StringVar(&base64Data, "base64", "", "Base64 encoded file data source.")
	create.Flags().StringVar(&fileName, "file-name", "", "Optional file name.")
	create.Flags().BoolVar(&urlOnly, "url-only", false, "Print only the uploaded file URL.")
	filesCmd.AddCommand(create)
	return filesCmd
}

func (c *cli) versionCommand() *cobra.Command {
	return &cobra.Command{Use: "version", Short: "Print the CLI version as JSON", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, args []string) error {
		return c.writeJSON(map[string]string{"version": runapi.Version})
	}}
}

func (c *cli) accountCommand() *cobra.Command {
	account := &cobra.Command{Use: "account", Short: "Account-related commands", Args: cobra.NoArgs}
	account.AddCommand(&cobra.Command{Use: "info", Short: "Fetch the current user and account record", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, args []string) error {
		client, callOpts, ctx, cancel, err := c.clientForCommand(cmd)
		if err != nil {
			return err
		}
		defer cancel()
		response, err := client.Account.Info(ctx, callOpts...)
		if err != nil {
			return err
		}
		return c.writeJSON(response)
	}})
	account.AddCommand(&cobra.Command{Use: "balance", Short: "Fetch account balance and spend counters", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, args []string) error {
		client, callOpts, ctx, cancel, err := c.clientForCommand(cmd)
		if err != nil {
			return err
		}
		defer cancel()
		response, err := client.Account.Balance(ctx, callOpts...)
		if err != nil {
			return err
		}
		return c.writeJSON(response)
	}})
	return account
}

func (c *cli) serviceCommand(service string) *cobra.Command {
	serviceCmd := &cobra.Command{Use: service, Short: fmt.Sprintf("%s service actions", service), Args: cobra.NoArgs}
	for _, spec := range allSpecs {
		if spec.service != service {
			continue
		}
		var input, inputFile string
		long := describeAction(spec)
		inputFields := composeInputFields(spec)
		if inputFields != "" {
			long += "\n\n" + inputFields
		}
		actionCmd := &cobra.Command{Use: spec.action, Short: describeAction(spec), Long: long, Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, args []string) error {
			payload, err := c.readInput(input, inputFile)
			if err != nil {
				return err
			}
			if _, err := spec.decode(payload); err != nil {
				return err
			}
			var client *runapi.Client
			var callOpts []option.RequestOption
			var ctx context.Context
			var cancel context.CancelFunc
			defer func() {
				if cancel != nil {
					cancel()
				}
			}()
			ensureClient := func() error {
				if client != nil {
					return nil
				}
				var err error
				client, callOpts, ctx, cancel, err = c.clientForCommand(cmd)
				return err
			}
			uploadFile := func(value string) (string, error) {
				if err := ensureClient(); err != nil {
					return "", err
				}
				response, err := client.Files.Create(ctx, files.CreateParams{File: value}, callOpts...)
				if err != nil {
					return "", err
				}
				return response.URL, nil
			}

			mediaFields := mediaInputFieldsForSpec(spec)
			payload, err = c.autoUploadMediaInputs(mediaFields, payload, uploadFile)
			if err != nil {
				return err
			}

			params, err := spec.decode(payload)
			if err != nil {
				return err
			}
			if err := ensureClient(); err != nil {
				return err
			}

			if !spec.isAsync {
				response, err := spec.run(ctx, client, params, callOpts)
				if err != nil {
					return err
				}
				return c.writeJSON(response)
			}

			if c.async {
				response, err := c.createTask(ctx, spec, client, params, payload, mediaFields, callOpts)
				if err != nil {
					return err
				}
				return c.writeJSON(response)
			}

			created, err := c.createTask(ctx, spec, client, params, payload, mediaFields, callOpts)
			if err != nil {
				return err
			}
			c.logf("submitted id=%s", created.ID)
			response, err := c.waitFor(ctx, spec, client, created.ID, callOpts)
			if err != nil {
				return err
			}
			return c.writeJSON(response)
		}}
		actionCmd.Flags().StringVar(&input, "input", "", "Inline JSON object payload for this action. Mutually exclusive with --input-file.")
		actionCmd.Flags().StringVar(&inputFile, "input-file", "", "Path to a JSON payload file. Use '-' to read JSON from stdin. Mutually exclusive with --input.")
		serviceCmd.AddCommand(actionCmd)
	}
	return serviceCmd
}

func (c *cli) createTask(ctx context.Context, spec actionSpec, client *runapi.Client, params any, payload []byte, mediaFields []string, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
	if !mediaFieldsContainJSONNull(mediaFields, payload) {
		return spec.create(ctx, client, params, opts)
	}

	path, err := validateCreateParamsAndPath(ctx, spec, params, opts)
	if err != nil {
		return nil, err
	}
	body := core.CompactParams(params)
	overlayRawMediaFieldsWithNull(body, mediaFields, payload)
	return client.CreateTaskRaw(ctx, path, body, opts...)
}

type createValidationHTTPClient struct {
	path string
}

func (c *createValidationHTTPClient) Request(_ context.Context, _ string, path string, _ *core.HTTPRequestOptions) (json.RawMessage, error) {
	c.path = path
	return json.RawMessage(`{"id":"validation","status":"processing"}`), nil
}

func validateCreateParamsAndPath(ctx context.Context, spec actionSpec, params any, opts []option.RequestOption) (string, error) {
	stub := &createValidationHTTPClient{}
	_, err := spec.create(ctx, runapi.NewClientWithHTTP(stub), params, opts)
	if err != nil {
		return "", err
	}
	return stub.path, nil
}

func (c *cli) getCommand() *cobra.Command {
	var service, action string
	cmd := &cobra.Command{Use: "get <task-id>", Short: "Fetch the current state of an async task", Long: "Get retrieves a task result for a specific service/action pair without waiting for completion.", RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return cmd.Help()
		}
		if err := requireServiceAction(service, action); err != nil {
			return err
		}
		spec, err := findActionSpec(service, action)
		if err != nil {
			return err
		}
		if !spec.isAsync {
			return core.NewError(core.ErrValidation, "service/action does not support get", 422, "", nil, nil)
		}
		client, callOpts, ctx, cancel, err := c.clientForCommand(cmd)
		if err != nil {
			return err
		}
		defer cancel()
		response, err := spec.get(ctx, client, args[0], callOpts)
		if err != nil {
			return err
		}
		return c.writeJSON(response)
	}}
	cmd.Flags().StringVar(&service, "service", "", "Service name, for example suno, veo-3-1, or nano-banana.")
	cmd.Flags().StringVar(&action, "action", "", "Action name within the selected service, for example generate or extend.")
	return cmd
}

func (c *cli) waitCommand() *cobra.Command {
	var service, action string
	cmd := &cobra.Command{Use: "wait <task-id>", Short: "Poll an async task until completion", Long: "Wait repeatedly fetches a task result for a specific service/action pair until it completes, fails, or times out.", RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return cmd.Help()
		}
		if err := requireServiceAction(service, action); err != nil {
			return err
		}
		spec, err := findActionSpec(service, action)
		if err != nil {
			return err
		}
		if !spec.isAsync {
			return core.NewError(core.ErrValidation, "service/action does not support wait", 422, "", nil, nil)
		}
		client, callOpts, ctx, cancel, err := c.clientForCommand(cmd)
		if err != nil {
			return err
		}
		defer cancel()
		response, err := c.waitFor(ctx, spec, client, args[0], callOpts)
		if err != nil {
			return err
		}
		return c.writeJSON(response)
	}}
	cmd.Flags().StringVar(&service, "service", "", "Service name, for example suno, veo-3-1, or nano-banana.")
	cmd.Flags().StringVar(&action, "action", "", "Action name within the selected service, for example generate or extend.")
	return cmd
}

func (c *cli) clientForCommand(cmd *cobra.Command) (*runapi.Client, []option.RequestOption, context.Context, context.CancelFunc, error) {
	cfg, err := loadConfig()
	if err != nil {
		return nil, nil, nil, nil, err
	}
	apiKey := firstNonEmpty(strings.TrimSpace(c.apiKeyFlag), strings.TrimSpace(os.Getenv("RUNAPI_API_KEY")), strings.TrimSpace(cfg.APIKey))
	baseURL := firstNonEmpty(strings.TrimSpace(c.baseURLFlag), strings.TrimSpace(os.Getenv("RUNAPI_BASE_URL")), strings.TrimSpace(cfg.BaseURL), core.DefaultBaseURL)
	client, err := c.newClient(option.WithAPIKey(apiKey), option.WithBaseURL(baseURL), option.WithTimeout(c.timeout), option.WithUserAgent(core.CLIUserAgent(runapi.Version)))
	if err != nil {
		return nil, nil, nil, nil, err
	}
	ctx, cancel := context.WithTimeout(cmd.Context(), c.timeout)
	callOpts := []option.RequestOption{option.WithPollInterval(c.pollInterval), option.WithMaxWait(c.timeout)}
	return client, callOpts, ctx, cancel, nil
}

func mediaInputFieldsForSpec(spec actionSpec) []string {
	value, err := spec.decode([]byte(`{}`))
	if err != nil || value == nil {
		return nil
	}
	fields := collectJSONFields(reflect.TypeOf(value))
	mediaFields := make([]string, 0)
	for _, field := range fields {
		if isMediaURLInputField(field) {
			mediaFields = append(mediaFields, field.name)
		}
	}
	return mediaFields
}

// autoUploadMediaInputs uploads any local file paths found in the spec's media
// URL fields and rewrites the payload to carry the returned hosted URLs. It
// decodes the payload once; upload lazily creates the client only when a local
// path is actually found, and the original payload bytes are returned untouched
// when nothing needs uploading.
func (c *cli) autoUploadMediaInputs(mediaFields []string, payload []byte, upload func(value string) (string, error)) ([]byte, error) {
	object, ok := decodeRawJSONObject(payload)
	if !ok {
		return payload, nil
	}
	uploadedFiles := map[string]string{}
	uploadOnce := func(value string) (string, error) {
		cacheKey, err := filepath.Abs(value)
		if err != nil {
			cacheKey = value
		}
		if url, ok := uploadedFiles[cacheKey]; ok {
			return url, nil
		}
		url, err := upload(value)
		if err != nil {
			return "", err
		}
		uploadedFiles[cacheKey] = url
		return url, nil
	}
	changed := false
	for _, field := range mediaFields {
		raw, ok := object[field]
		if !ok {
			continue
		}
		next, fieldChanged, err := c.autoUploadMediaValue(field, raw, uploadOnce)
		if err != nil {
			return nil, err
		}
		if fieldChanged {
			object[field] = next
			changed = true
		}
	}
	if !changed {
		return payload, nil
	}
	next, err := json.Marshal(object)
	if err != nil {
		return nil, core.NewError(core.ErrValidation, "failed to encode uploaded input", 422, "", nil, err)
	}
	return next, nil
}

func mediaFieldsContainJSONNull(mediaFields []string, payload []byte) bool {
	object, ok := decodeRawJSONObject(payload)
	if !ok {
		return false
	}
	for _, field := range mediaFields {
		raw, ok := object[field]
		if ok && rawContainsJSONNull(raw) {
			return true
		}
	}
	return false
}

func overlayRawMediaFieldsWithNull(body map[string]any, mediaFields []string, payload []byte) {
	object, ok := decodeRawJSONObject(payload)
	if !ok {
		return
	}
	for _, field := range mediaFields {
		raw, ok := object[field]
		if ok && rawContainsJSONNull(raw) {
			body[field] = raw
		}
	}
}

func rawContainsJSONNull(raw json.RawMessage) bool {
	trimmed := bytes.TrimSpace(raw)
	if bytes.Equal(trimmed, []byte("null")) {
		return true
	}

	var values []json.RawMessage
	if err := json.Unmarshal(raw, &values); err != nil {
		return false
	}
	for _, value := range values {
		if bytes.Equal(bytes.TrimSpace(value), []byte("null")) {
			return true
		}
	}
	return false
}

func (c *cli) autoUploadMediaValue(field string, raw json.RawMessage, upload func(value string) (string, error)) (json.RawMessage, bool, error) {
	var single *string
	if err := json.Unmarshal(raw, &single); err == nil && single != nil {
		next, changed, err := c.autoUploadMediaString(field, *single, upload)
		if err != nil || !changed {
			return raw, false, err
		}
		encoded, err := json.Marshal(next)
		return encoded, true, err
	}

	var values []*string
	if err := json.Unmarshal(raw, &values); err == nil {
		changed := false
		for i, value := range values {
			if value == nil {
				continue
			}
			fieldRef := fmt.Sprintf("%s[%d]", field, i)
			next, valueChanged, err := c.autoUploadMediaString(fieldRef, *value, upload)
			if err != nil {
				return raw, false, err
			}
			if valueChanged {
				values[i] = &next
				changed = true
			}
		}
		if !changed {
			return raw, false, nil
		}
		encoded, err := json.Marshal(values)
		return encoded, true, err
	}

	return raw, false, nil
}

func (c *cli) autoUploadMediaString(fieldRef, value string, upload func(value string) (string, error)) (string, bool, error) {
	if !isReadableRegularLocalPath(value) {
		return value, false, nil
	}
	url, err := upload(value)
	if err != nil {
		return value, false, err
	}
	c.logf("uploaded %s from %s", fieldRef, value)
	return url, true, nil
}

func decodeRawJSONObject(payload []byte) (map[string]json.RawMessage, bool) {
	var object map[string]json.RawMessage
	if err := json.Unmarshal(payload, &object); err != nil || object == nil {
		return nil, false
	}
	return object, true
}

func isReadableRegularLocalPath(value string) bool {
	if value == "" || isHTTPURL(value) {
		return false
	}
	info, err := os.Stat(value)
	if err != nil || !info.Mode().IsRegular() {
		return false
	}
	file, err := os.Open(value)
	if err != nil {
		return false
	}
	_ = file.Close()
	return true
}

func isHTTPURL(value string) bool {
	lower := strings.ToLower(value)
	return strings.HasPrefix(lower, "http://") || strings.HasPrefix(lower, "https://")
}

func (c *cli) readInput(input, inputFile string) ([]byte, error) {
	if strings.TrimSpace(input) != "" && strings.TrimSpace(inputFile) != "" {
		return nil, core.NewError(core.ErrValidation, "--input and --input-file are mutually exclusive", 422, "", nil, nil)
	}
	if strings.TrimSpace(input) != "" {
		return []byte(input), nil
	}
	if strings.TrimSpace(inputFile) != "" {
		if inputFile == "-" {
			return io.ReadAll(c.stdin)
		}
		return os.ReadFile(inputFile)
	}
	if isReaderTTY(c.stdin) {
		return nil, core.NewError(core.ErrValidation, "input is required; use --input, --input-file, or stdin", 422, "", nil, nil)
	}
	return io.ReadAll(c.stdin)
}

func (c *cli) waitFor(ctx context.Context, spec actionSpec, client *runapi.Client, taskID string, callOpts []option.RequestOption) (core.TaskResponse, error) {
	timer := time.NewTimer(0)
	defer timer.Stop()
	<-timer.C

	for {
		response, err := spec.get(ctx, client, taskID, callOpts)
		if err != nil {
			return nil, err
		}
		status := core.NormalizeStatus(response.GetStatus())
		c.logf("polling status=%s", status)
		switch status {
		case "completed":
			return response, nil
		case "failed":
			message := response.GetError()
			if message == "" {
				message = "Task failed"
			}
			return nil, core.NewError(core.ErrTaskFailed, message, 0, "", response, nil)
		}
		timer.Reset(c.pollInterval)
		select {
		case <-ctx.Done():
			return nil, core.NewError(core.ErrTaskTimeout, "Task polling timed out", 0, "", nil, ctx.Err())
		case <-timer.C:
		}
	}
}

func (c *cli) writeJSON(value any) error {
	encoder := json.NewEncoder(c.stdout)
	encoder.SetEscapeHTML(false)
	return encoder.Encode(value)
}

func (c *cli) printError(err error) {
	if apiErr, ok := errors.AsType[*core.Error](err); ok {
		_ = c.writeJSON(map[string]any{"error": map[string]any{"message": apiErr.Message, "code": apiErr.Code, "status": apiErr.Status, "details": apiErr.Details}})
	} else {
		_ = c.writeJSON(map[string]any{"error": map[string]any{"message": err.Error()}})
	}
	c.logf(core.FormatError(err))
}

func (c *cli) logf(format string, args ...any) {
	if c.quiet {
		return
	}
	fmt.Fprintf(c.stderr, format+"\n", args...)
}

func loadConfig() (configFile, error) {
	path, err := configFilePath()
	if err != nil {
		return configFile{}, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return configFile{}, nil
		}
		return configFile{}, err
	}
	var cfg configFile
	if err := json.Unmarshal(data, &cfg); err != nil {
		return configFile{}, core.NewError(core.ErrValidation, "invalid config file", 422, "", nil, err)
	}
	return cfg, nil
}

func saveConfig(cfg configFile) error {
	path, err := configFilePath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o600)
}

func configFilePath() (string, error) {
	if dir := os.Getenv("RUNAPI_CONFIG_DIR"); dir != "" {
		return filepath.Join(dir, "config.json"), nil
	}
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "runapi", "config.json"), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "runapi", "config.json"), nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func isReaderTTY(reader io.Reader) bool {
	file, ok := reader.(*os.File)
	if !ok {
		return false
	}
	info, err := file.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}

func exitCode(err error) int {
	apiErr, ok := errors.AsType[*core.Error](err)
	if !ok {
		return 1
	}
	switch apiErr.Code {
	case core.ErrAuthentication:
		return 2
	case core.ErrInsufficientCredits:
		return 3
	case core.ErrValidation, core.ErrNotFound, core.ErrConflict:
		return 4
	case core.ErrTimeout, core.ErrTaskTimeout:
		return 5
	case core.ErrRateLimit:
		return 6
	case core.ErrTaskFailed:
		return 7
	default:
		return 1
	}
}

func requireServiceAction(service, action string) error {
	if strings.TrimSpace(service) == "" || strings.TrimSpace(action) == "" {
		return core.NewError(core.ErrValidation, "--service and --action are required", 422, "", nil, nil)
	}
	return nil
}

func findActionSpec(service, action string) (actionSpec, error) {
	for _, spec := range allSpecs {
		if spec.service == service && spec.action == action {
			return spec, nil
		}
	}
	return actionSpec{}, core.NewError(core.ErrValidation, "unsupported service/action", 422, "", nil, nil)
}

func decodeInto[T any](data []byte) (any, error) {
	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		return nil, core.NewError(core.ErrValidation, "input must be valid JSON for the selected action", 422, "", nil, err)
	}
	return value, nil
}

func inputFieldsFor[T any]() string {
	var zero T
	fields := collectJSONFields(reflect.TypeOf(zero))
	if len(fields) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("Input fields (JSON):\n")
	seen := map[string]bool{}
	for _, f := range fields {
		if seen[f.name] {
			continue
		}
		seen[f.name] = true
		fmt.Fprintf(&b, "  %-24s %-10s %s\n", f.name, f.typ, f.help)
	}
	return b.String()
}

func composeInputFields(spec actionSpec) string {
	return appendGeneratedContractHelp(spec.inputFields, spec.service, spec.action)
}

type generatedContractField struct {
	Enum []any
}

type generatedContractAction struct {
	Models        []string
	FieldsByModel map[string]map[string]generatedContractField
}

func appendGeneratedContractHelp(inputFields, service, action string) string {
	if inputFields == "" {
		return ""
	}
	contract, ok := generatedContract[service+"/"+action]
	if !ok {
		return inputFields
	}

	lines := strings.Split(inputFields, "\n")
	for i, line := range lines {
		field := strings.Fields(strings.TrimSpace(line))
		if len(field) == 0 || field[0] == "Input" {
			continue
		}
		sentence := generatedContractHelpSentenceFor(contract, field[0])
		if sentence == "" || strings.Contains(line, "Accepted values") {
			continue
		}
		lines[i] = appendHelpSentence(line, sentence)
	}
	return strings.Join(lines, "\n")
}

func generatedContractHelpSentenceFor(contract generatedContractAction, field string) string {
	if field == "model" {
		if len(contract.Models) == 0 {
			return ""
		}
		return fmt.Sprintf("Accepted values: %s.", strings.Join(contract.Models, ", "))
	}

	valuesByModel := generatedContractValuesByModelFor(contract, field)
	if len(valuesByModel) == 0 {
		return ""
	}

	commonValues := commonGeneratedContractValues(generatedContractFieldModelKeys(contract), valuesByModel)
	if len(commonValues) > 0 {
		return fmt.Sprintf("Accepted values: %s.", strings.Join(commonValues, ", "))
	}

	return fmt.Sprintf("Accepted values by model: %s.", formatGeneratedContractValuesByModel(contract.Models, valuesByModel))
}

func generatedContractValuesByModelFor(contract generatedContractAction, field string) map[string][]string {
	valuesByModel := map[string][]string{}
	for _, model := range generatedContractFieldModelKeys(contract) {
		fields := contract.FieldsByModel[model]
		if fieldContract, ok := fields[field]; ok && len(fieldContract.Enum) > 0 {
			valuesByModel[model] = generatedContractEnumStrings(fieldContract.Enum)
		}
	}
	return valuesByModel
}

func generatedContractEnumStrings(values []any) []string {
	strings := make([]string, 0, len(values))
	for _, value := range values {
		strings = append(strings, fmt.Sprint(value))
	}
	return strings
}

func generatedContractFieldModelKeys(contract generatedContractAction) []string {
	if len(contract.Models) > 0 {
		return contract.Models
	}

	keys := make([]string, 0, len(contract.FieldsByModel))
	for model := range contract.FieldsByModel {
		keys = append(keys, model)
	}
	sort.Strings(keys)
	return keys
}

func commonGeneratedContractValues(models []string, valuesByModel map[string][]string) []string {
	if len(models) == 0 || len(valuesByModel) != len(models) {
		return nil
	}

	var common []string
	for _, model := range models {
		values := valuesByModel[model]
		if len(values) == 0 {
			return nil
		}
		if common == nil {
			common = values
			continue
		}
		if strings.Join(common, "\x00") != strings.Join(values, "\x00") {
			return nil
		}
	}
	return common
}

func formatGeneratedContractValuesByModel(models []string, valuesByModel map[string][]string) string {
	groups := []string{}
	seen := map[string]int{}

	for _, model := range models {
		values := valuesByModel[model]
		if len(values) == 0 {
			continue
		}

		key := strings.Join(values, "\x00")
		groupIndex, ok := seen[key]
		if !ok {
			seen[key] = len(groups)
			groups = append(groups, fmt.Sprintf("%s: %s", model, strings.Join(values, ", ")))
			continue
		}

		prefix, valuesText, _ := strings.Cut(groups[groupIndex], ": ")
		groups[groupIndex] = fmt.Sprintf("%s, %s: %s", prefix, model, valuesText)
	}

	return strings.Join(groups, "; ")
}

func appendHelpSentence(line, sentence string) string {
	trimmed := strings.TrimRight(line, " ")
	if strings.HasSuffix(strings.TrimSpace(trimmed), ".") {
		return trimmed + " " + sentence
	}
	return trimmed + ". " + sentence
}

type jsonField struct {
	name string
	typ  string
	help string
}

func collectJSONFields(t reflect.Type) []jsonField {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil
	}
	var fields []jsonField
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		if !sf.IsExported() {
			continue
		}
		if sf.Anonymous {
			fields = append(fields, collectJSONFields(sf.Type)...)
			continue
		}
		tag := sf.Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}
		name, _, _ := strings.Cut(tag, ",")
		typ := jsonTypeName(sf.Type)
		help := sf.Tag.Get("help")
		fields = append(fields, jsonField{name: name, typ: typ, help: help})
	}
	return fields
}

func isMediaURLInputField(field jsonField) bool {
	return isStringURLField(field) && !isNonMediaURLField(field.name, field.help)
}

func isStringURLField(field jsonField) bool {
	if !(strings.HasSuffix(field.name, "_url") || strings.HasSuffix(field.name, "_urls") || strings.HasSuffix(field.name, "_url_list")) {
		return false
	}
	switch field.typ {
	case "string", "[]string", "[2]string":
		return true
	default:
		return false
	}
}

func isNonMediaURLField(name, help string) bool {
	lowerName := strings.ToLower(name)
	lowerHelp := strings.ToLower(help)
	nonMediaTokens := []string{"callback", "webhook", "completion callback", "notification", "origin", "result", "stream", "cdn"}
	for _, token := range nonMediaTokens {
		if strings.Contains(lowerName, token) || strings.Contains(lowerHelp, token) {
			return true
		}
	}
	return false
}

func jsonTypeName(t reflect.Type) string {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	switch t.Kind() {
	case reflect.String:
		return "string"
	case reflect.Bool:
		return "bool"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return "int"
	case reflect.Float32, reflect.Float64:
		return "float"
	case reflect.Slice:
		return "[]" + jsonTypeName(t.Elem())
	case reflect.Array:
		return fmt.Sprintf("[%d]%s", t.Len(), jsonTypeName(t.Elem()))
	default:
		return t.Name()
	}
}

func describeAction(spec actionSpec) string {
	if spec.isAsync {
		return fmt.Sprintf("Run the %s %s action", spec.service, spec.action)
	}
	return fmt.Sprintf("Run the synchronous %s %s action", spec.service, spec.action)
}

var allSpecs = []actionSpec{
	newSunoTextToMusicSpec(), newSunoExtendMusicSpec(), newSunoGenerateArtworkSpec(), newSunoCoverAudioSpec(),
	newSunoAddInstrumentalSpec(), newSunoAddVocalsSpec(), newSunoSeparateAudioStemsSpec(), newSunoGenerateMidiSpec(), newSunoConvertAudioSpec(),
	newSunoVisualizeMusicSpec(), newSunoGenerateLyricsSpec(), newSunoGetTimestampedLyricsSpec(), newSunoReplaceSectionSpec(), newSunoCreateMashupSpec(),
	newSunoTextToSoundSpec(), newSunoVoiceToValidationPhraseSpec(), newSunoRegenerateValidationPhraseSpec(), newSunoGenerateVoiceSpec(), newSunoCheckVoiceSpec(),
	newSunoGeneratePersonaSpec(), newSunoBoostStyleSpec(),
	newVeo31TextToVideoSpec(), newVeo31ExtendVideoSpec(), newVeo31UpscaleVideoSpec(),
	newNanoBananaTextToImageSpec(), newNanoBananaEditImageSpec(), newImagen4TextToImageSpec(), newImagen4RemixImageSpec(),
	newSeedanceTextToVideoSpec(),
	newSeedreamTextToImageSpec(), newSeedreamEditImageSpec(),
	newRunwayTextToVideoSpec(), newRunwayExtendVideoSpec(), newRunwayAlephEditVideoSpec(),
	newKlingTextToVideoSpec(), newKlingAvatarSpec(), newKlingImageToVideoSpec(), newKlingMotionControlSpec(),
	newFluxKontextTextToImageSpec(), newFlux2TextToImageSpec(), newFlux2RemixImageSpec(),
	newGeminiOmniCreateAudioSpec(), newGeminiOmniCreateCharacterSpec(), newGeminiOmniTextToVideoSpec(),
	newQwen2TextToImageSpec(), newQwen2RemixImageSpec(), newQwen2EditSpec(),
	newRecraftUpscaleSpec(), newRecraftBackgroundRemovalSpec(), newZImageTextToImageSpec(),
	newIdeogramV3TextToImageSpec(), newIdeogramV3EditImageSpec(), newIdeogramV3RemixImageSpec(), newIdeogramV3ReframeImageSpec(),
	newElevenlabsSpeechSpec(), newElevenlabsDialogueSpec(), newElevenlabsSoundEffectSpec(), newElevenlabsTranscriptionSpec(), newElevenlabsAudioIsolationSpec(),
	newInfiniteTalkAudioToVideoSpec(),
	newOmniHumanAudioToVideoSpec(), newOmniHumanHumanIdentificationSpec(), newOmniHumanSubjectDetectionSpec(),
	newWanTextToVideoSpec(), newWanImageToVideoSpec(), newWanSpeechToVideoSpec(),
	newWanAnimateSpec(), newWanTextToImageSpec(), newWanEditVideoSpec(),
	newLumaModifySpec(),
	newHailuoTextToVideoSpec(), newHailuoImageToVideoSpec(),
	newVolcengineLipSyncVideoSpec(),
	newHappyHorseTextToVideoSpec(), newHappyHorseImageToVideoSpec(), newHappyHorseEditVideoSpec(),
	newGptImageTextToImageSpec(), newGptImageEditImageSpec(), newGptImage2TextToImageSpec(), newGptImage2EditImageSpec(), newGpt4oImageTextToImageSpec(),
	newGrokImagineTextToVideoSpec(), newGrokImagineImageToVideoSpec(), newGrokImagineTextToImageSpec(), newGrokImagineEditImageSpec(),
	newGrokImagineExtendSpec(), newGrokImagineUpscaleSpec(),
	newTopazUpscaleImageSpec(), newTopazUpscaleVideoSpec(),
}

func newSunoTextToMusicSpec() actionSpec {
	return actionSpec{service: "suno", action: "text-to-music", isAsync: true, inputFields: inputFieldsFor[suno.TextToMusicParams](), decode: decodeInto[suno.TextToMusicParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.Suno.TextToMusic.Create(ctx, params.(suno.TextToMusicParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.Suno.TextToMusic.Run(ctx, params.(suno.TextToMusicParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.Suno.TextToMusic.Get(ctx, id, opts...)
	}}
}

func newSunoExtendMusicSpec() actionSpec {
	return actionSpec{service: "suno", action: "extend-music", isAsync: true, inputFields: inputFieldsFor[suno.ExtendMusicParams](), decode: decodeInto[suno.ExtendMusicParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.Suno.ExtendMusic.Create(ctx, params.(suno.ExtendMusicParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.Suno.ExtendMusic.Run(ctx, params.(suno.ExtendMusicParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.Suno.ExtendMusic.Get(ctx, id, opts...)
	}}
}

func newSunoGenerateArtworkSpec() actionSpec {
	return actionSpec{service: "suno", action: "generate-artwork", isAsync: true, inputFields: inputFieldsFor[suno.GenerateArtworkParams](), decode: decodeInto[suno.GenerateArtworkParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.Suno.GenerateArtwork.Create(ctx, params.(suno.GenerateArtworkParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.Suno.GenerateArtwork.Run(ctx, params.(suno.GenerateArtworkParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.Suno.GenerateArtwork.Get(ctx, id, opts...)
	}}
}

func newSunoCoverAudioSpec() actionSpec {
	return actionSpec{service: "suno", action: "cover-audio", isAsync: true, inputFields: inputFieldsFor[suno.CoverAudioParams](), decode: decodeInto[suno.CoverAudioParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.Suno.CoverAudio.Create(ctx, params.(suno.CoverAudioParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.Suno.CoverAudio.Run(ctx, params.(suno.CoverAudioParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.Suno.CoverAudio.Get(ctx, id, opts...)
	}}
}

func newSunoAddInstrumentalSpec() actionSpec {
	return actionSpec{service: "suno", action: "add-instrumental", isAsync: true, inputFields: inputFieldsFor[suno.AddInstrumentalParams](), decode: decodeInto[suno.AddInstrumentalParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.Suno.AddInstrumental.Create(ctx, params.(suno.AddInstrumentalParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.Suno.AddInstrumental.Run(ctx, params.(suno.AddInstrumentalParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.Suno.AddInstrumental.Get(ctx, id, opts...)
	}}
}

func newSunoAddVocalsSpec() actionSpec {
	return actionSpec{service: "suno", action: "add-vocals", isAsync: true, inputFields: inputFieldsFor[suno.AddVocalsParams](), decode: decodeInto[suno.AddVocalsParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.Suno.AddVocals.Create(ctx, params.(suno.AddVocalsParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.Suno.AddVocals.Run(ctx, params.(suno.AddVocalsParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.Suno.AddVocals.Get(ctx, id, opts...)
	}}
}

func newSunoSeparateAudioStemsSpec() actionSpec {
	return actionSpec{service: "suno", action: "separate-audio-stems", isAsync: true, inputFields: inputFieldsFor[suno.SeparateAudioStemsParams](), decode: decodeInto[suno.SeparateAudioStemsParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.Suno.SeparateAudioStems.Create(ctx, params.(suno.SeparateAudioStemsParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.Suno.SeparateAudioStems.Run(ctx, params.(suno.SeparateAudioStemsParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.Suno.SeparateAudioStems.Get(ctx, id, opts...)
	}}
}

func newSunoGenerateMidiSpec() actionSpec {
	return actionSpec{service: "suno", action: "generate-midi", isAsync: true, inputFields: inputFieldsFor[suno.GenerateMidiParams](), decode: decodeInto[suno.GenerateMidiParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.Suno.GenerateMidi.Create(ctx, params.(suno.GenerateMidiParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.Suno.GenerateMidi.Run(ctx, params.(suno.GenerateMidiParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.Suno.GenerateMidi.Get(ctx, id, opts...)
	}}
}

func newSunoConvertAudioSpec() actionSpec {
	return actionSpec{service: "suno", action: "convert-audio", isAsync: true, inputFields: inputFieldsFor[suno.ConvertAudioParams](), decode: decodeInto[suno.ConvertAudioParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.Suno.ConvertAudio.Create(ctx, params.(suno.ConvertAudioParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.Suno.ConvertAudio.Run(ctx, params.(suno.ConvertAudioParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.Suno.ConvertAudio.Get(ctx, id, opts...)
	}}
}

func newSunoVisualizeMusicSpec() actionSpec {
	return actionSpec{service: "suno", action: "visualize-music", isAsync: true, inputFields: inputFieldsFor[suno.VisualizeMusicParams](), decode: decodeInto[suno.VisualizeMusicParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.Suno.VisualizeMusic.Create(ctx, params.(suno.VisualizeMusicParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.Suno.VisualizeMusic.Run(ctx, params.(suno.VisualizeMusicParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.Suno.VisualizeMusic.Get(ctx, id, opts...)
	}}
}

func newSunoGenerateLyricsSpec() actionSpec {
	return actionSpec{service: "suno", action: "generate-lyrics", isAsync: true, inputFields: inputFieldsFor[suno.GenerateLyricsParams](), decode: decodeInto[suno.GenerateLyricsParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.Suno.GenerateLyrics.Create(ctx, params.(suno.GenerateLyricsParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.Suno.GenerateLyrics.Run(ctx, params.(suno.GenerateLyricsParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.Suno.GenerateLyrics.Get(ctx, id, opts...)
	}}
}

func newSunoGetTimestampedLyricsSpec() actionSpec {
	return actionSpec{service: "suno", action: "get-timestamped-lyrics", isAsync: false, inputFields: inputFieldsFor[suno.GetTimestampedLyricsParams](), decode: decodeInto[suno.GetTimestampedLyricsParams], run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.Suno.GetTimestampedLyrics.Run(ctx, params.(suno.GetTimestampedLyricsParams), opts...)
	}}
}

func newSunoReplaceSectionSpec() actionSpec {
	return actionSpec{service: "suno", action: "replace-section", isAsync: true, inputFields: inputFieldsFor[suno.ReplaceSectionParams](), decode: decodeInto[suno.ReplaceSectionParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.Suno.ReplaceSection.Create(ctx, params.(suno.ReplaceSectionParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.Suno.ReplaceSection.Run(ctx, params.(suno.ReplaceSectionParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.Suno.ReplaceSection.Get(ctx, id, opts...)
	}}
}

func newSunoCreateMashupSpec() actionSpec {
	return actionSpec{service: "suno", action: "create-mashup", isAsync: true, inputFields: inputFieldsFor[suno.CreateMashupParams](), decode: decodeInto[suno.CreateMashupParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.Suno.CreateMashup.Create(ctx, params.(suno.CreateMashupParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.Suno.CreateMashup.Run(ctx, params.(suno.CreateMashupParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.Suno.CreateMashup.Get(ctx, id, opts...)
	}}
}

func newSunoTextToSoundSpec() actionSpec {
	return actionSpec{service: "suno", action: "text-to-sound", isAsync: true, inputFields: inputFieldsFor[suno.TextToSoundParams](), decode: decodeInto[suno.TextToSoundParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.Suno.TextToSound.Create(ctx, params.(suno.TextToSoundParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.Suno.TextToSound.Run(ctx, params.(suno.TextToSoundParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.Suno.TextToSound.Get(ctx, id, opts...)
	}}
}

func newSunoVoiceToValidationPhraseSpec() actionSpec {
	return actionSpec{service: "suno", action: "voice-to-validation-phrase", isAsync: true, inputFields: inputFieldsFor[suno.VoiceToValidationPhraseParams](), decode: decodeInto[suno.VoiceToValidationPhraseParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.Suno.VoiceToValidationPhrase.Create(ctx, params.(suno.VoiceToValidationPhraseParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.Suno.VoiceToValidationPhrase.Run(ctx, params.(suno.VoiceToValidationPhraseParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.Suno.VoiceToValidationPhrase.Get(ctx, id, opts...)
	}}
}

func newSunoRegenerateValidationPhraseSpec() actionSpec {
	return actionSpec{service: "suno", action: "regenerate-validation-phrase", isAsync: true, inputFields: inputFieldsFor[suno.RegenerateValidationPhraseParams](), decode: decodeInto[suno.RegenerateValidationPhraseParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.Suno.RegenerateValidationPhrase.Create(ctx, params.(suno.RegenerateValidationPhraseParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.Suno.RegenerateValidationPhrase.Run(ctx, params.(suno.RegenerateValidationPhraseParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.Suno.RegenerateValidationPhrase.Get(ctx, id, opts...)
	}}
}

func newSunoGenerateVoiceSpec() actionSpec {
	return actionSpec{service: "suno", action: "generate-voice", isAsync: true, inputFields: inputFieldsFor[suno.GenerateVoiceParams](), decode: decodeInto[suno.GenerateVoiceParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.Suno.GenerateVoice.Create(ctx, params.(suno.GenerateVoiceParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.Suno.GenerateVoice.Run(ctx, params.(suno.GenerateVoiceParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.Suno.GenerateVoice.Get(ctx, id, opts...)
	}}
}

func newSunoCheckVoiceSpec() actionSpec {
	return actionSpec{service: "suno", action: "check-voice", isAsync: false, inputFields: inputFieldsFor[suno.CheckVoiceParams](), decode: decodeInto[suno.CheckVoiceParams], run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.Suno.CheckVoice.Run(ctx, params.(suno.CheckVoiceParams), opts...)
	}}
}

func newSunoGeneratePersonaSpec() actionSpec {
	return actionSpec{service: "suno", action: "generate-persona", isAsync: false, inputFields: inputFieldsFor[suno.GeneratePersonaParams](), decode: decodeInto[suno.GeneratePersonaParams], run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.Suno.GeneratePersona.Run(ctx, params.(suno.GeneratePersonaParams), opts...)
	}}
}

func newSunoBoostStyleSpec() actionSpec {
	return actionSpec{service: "suno", action: "boost-style", isAsync: false, inputFields: inputFieldsFor[suno.BoostStyleParams](), decode: decodeInto[suno.BoostStyleParams], run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.Suno.BoostStyle.Run(ctx, params.(suno.BoostStyleParams), opts...)
	}}
}

func newVeo31TextToVideoSpec() actionSpec {
	return actionSpec{service: "veo-3-1", action: "text-to-video", isAsync: true, inputFields: inputFieldsFor[veo31.TextToVideoParams](), decode: decodeInto[veo31.TextToVideoParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.Veo31.TextToVideo.Create(ctx, params.(veo31.TextToVideoParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.Veo31.TextToVideo.Run(ctx, params.(veo31.TextToVideoParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.Veo31.TextToVideo.Get(ctx, id, opts...)
	}}
}
func newVeo31ExtendVideoSpec() actionSpec {
	return actionSpec{service: "veo-3-1", action: "extend-video", isAsync: true, inputFields: inputFieldsFor[veo31.ExtendVideoParams](), decode: decodeInto[veo31.ExtendVideoParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.Veo31.ExtendVideo.Create(ctx, params.(veo31.ExtendVideoParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.Veo31.ExtendVideo.Run(ctx, params.(veo31.ExtendVideoParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.Veo31.ExtendVideo.Get(ctx, id, opts...)
	}}
}
func newVeo31UpscaleVideoSpec() actionSpec {
	return actionSpec{service: "veo-3-1", action: "upscale-video", isAsync: true, inputFields: inputFieldsFor[veo31.UpscaleVideoParams](), decode: decodeInto[veo31.UpscaleVideoParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.Veo31.UpscaleVideo.Create(ctx, params.(veo31.UpscaleVideoParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.Veo31.UpscaleVideo.Run(ctx, params.(veo31.UpscaleVideoParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.Veo31.UpscaleVideo.Get(ctx, id, opts...)
	}}
}
func newRunwayTextToVideoSpec() actionSpec {
	return actionSpec{service: "runway", action: "text-to-video", isAsync: true, inputFields: inputFieldsFor[runway.TextToVideoParams](), decode: decodeInto[runway.TextToVideoParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.Runway.TextToVideo.Create(ctx, params.(runway.TextToVideoParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.Runway.TextToVideo.Run(ctx, params.(runway.TextToVideoParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.Runway.TextToVideo.Get(ctx, id, opts...)
	}}
}

func newRunwayExtendVideoSpec() actionSpec {
	return actionSpec{service: "runway", action: "extend-video", isAsync: true, inputFields: inputFieldsFor[runway.ExtendVideoParams](), decode: decodeInto[runway.ExtendVideoParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.Runway.ExtendVideo.Create(ctx, params.(runway.ExtendVideoParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.Runway.ExtendVideo.Run(ctx, params.(runway.ExtendVideoParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.Runway.ExtendVideo.Get(ctx, id, opts...)
	}}
}

func newRunwayAlephEditVideoSpec() actionSpec {
	return actionSpec{service: "runway-aleph", action: "edit-video", isAsync: true, inputFields: inputFieldsFor[runwayaleph.EditVideoParams](), decode: decodeInto[runwayaleph.EditVideoParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.RunwayAleph.EditVideo.Create(ctx, params.(runwayaleph.EditVideoParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.RunwayAleph.EditVideo.Run(ctx, params.(runwayaleph.EditVideoParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.RunwayAleph.EditVideo.Get(ctx, id, opts...)
	}}
}
func newNanoBananaTextToImageSpec() actionSpec {
	return actionSpec{service: "nano-banana", action: "text-to-image", isAsync: true, inputFields: inputFieldsFor[nanobanana.TextToImageParams](), decode: decodeInto[nanobanana.TextToImageParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.NanoBanana.TextToImage.Create(ctx, params.(nanobanana.TextToImageParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.NanoBanana.TextToImage.Run(ctx, params.(nanobanana.TextToImageParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.NanoBanana.TextToImage.Get(ctx, id, opts...)
	}}
}
func newNanoBananaEditImageSpec() actionSpec {
	return actionSpec{service: "nano-banana", action: "edit-image", isAsync: true, inputFields: inputFieldsFor[nanobanana.EditImageParams](), decode: decodeInto[nanobanana.EditImageParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.NanoBanana.EditImage.Create(ctx, params.(nanobanana.EditImageParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.NanoBanana.EditImage.Run(ctx, params.(nanobanana.EditImageParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.NanoBanana.EditImage.Get(ctx, id, opts...)
	}}
}
func newImagen4TextToImageSpec() actionSpec {
	return actionSpec{service: "imagen-4", action: "text-to-image", isAsync: true, inputFields: inputFieldsFor[imagen4.TextToImageParams](), decode: decodeInto[imagen4.TextToImageParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.Imagen4.TextToImage.Create(ctx, params.(imagen4.TextToImageParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.Imagen4.TextToImage.Run(ctx, params.(imagen4.TextToImageParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.Imagen4.TextToImage.Get(ctx, id, opts...)
	}}
}
func newImagen4RemixImageSpec() actionSpec {
	return actionSpec{service: "imagen-4", action: "remix-image", isAsync: true, inputFields: inputFieldsFor[imagen4.RemixImageParams](), decode: decodeInto[imagen4.RemixImageParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.Imagen4.RemixImage.Create(ctx, params.(imagen4.RemixImageParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.Imagen4.RemixImage.Run(ctx, params.(imagen4.RemixImageParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.Imagen4.RemixImage.Get(ctx, id, opts...)
	}}
}
func newSeedanceTextToVideoSpec() actionSpec {
	return actionSpec{service: "seedance", action: "text-to-video", isAsync: true, inputFields: inputFieldsFor[seedance.TextToVideoParams](), decode: decodeInto[seedance.TextToVideoParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.Seedance.TextToVideo.Create(ctx, params.(seedance.TextToVideoParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.Seedance.TextToVideo.Run(ctx, params.(seedance.TextToVideoParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.Seedance.TextToVideo.Get(ctx, id, opts...)
	}}
}
func newSeedreamTextToImageSpec() actionSpec {
	return actionSpec{service: "seedream", action: "text-to-image", isAsync: true, inputFields: inputFieldsFor[seedream.TextToImageParams](), decode: decodeInto[seedream.TextToImageParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.Seedream.TextToImage.Create(ctx, params.(seedream.TextToImageParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.Seedream.TextToImage.Run(ctx, params.(seedream.TextToImageParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.Seedream.TextToImage.Get(ctx, id, opts...)
	}}
}
func newSeedreamEditImageSpec() actionSpec {
	return actionSpec{service: "seedream", action: "edit-image", isAsync: true, inputFields: inputFieldsFor[seedream.EditImageParams](), decode: decodeInto[seedream.EditImageParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.Seedream.EditImage.Create(ctx, params.(seedream.EditImageParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.Seedream.EditImage.Run(ctx, params.(seedream.EditImageParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.Seedream.EditImage.Get(ctx, id, opts...)
	}}
}
func newKlingTextToVideoSpec() actionSpec {
	return actionSpec{service: "kling", action: "text-to-video", isAsync: true, inputFields: inputFieldsFor[kling.TextToVideoParams](), decode: decodeInto[kling.TextToVideoParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.Kling.TextToVideo.Create(ctx, params.(kling.TextToVideoParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.Kling.TextToVideo.Run(ctx, params.(kling.TextToVideoParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.Kling.TextToVideo.Get(ctx, id, opts...)
	}}
}
func newKlingAvatarSpec() actionSpec {
	return actionSpec{service: "kling", action: "avatar", isAsync: true, inputFields: inputFieldsFor[kling.AiAvatarParams](), decode: decodeInto[kling.AiAvatarParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.Kling.AiAvatar.Create(ctx, params.(kling.AiAvatarParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.Kling.AiAvatar.Run(ctx, params.(kling.AiAvatarParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.Kling.AiAvatar.Get(ctx, id, opts...)
	}}
}
func newKlingImageToVideoSpec() actionSpec {
	return actionSpec{service: "kling", action: "image-to-video", isAsync: true, inputFields: inputFieldsFor[kling.ImageToVideoParams](), decode: decodeInto[kling.ImageToVideoParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.Kling.ImageToVideo.Create(ctx, params.(kling.ImageToVideoParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.Kling.ImageToVideo.Run(ctx, params.(kling.ImageToVideoParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.Kling.ImageToVideo.Get(ctx, id, opts...)
	}}
}
func newKlingMotionControlSpec() actionSpec {
	return actionSpec{service: "kling", action: "motion-control", isAsync: true, inputFields: inputFieldsFor[kling.MotionControlParams](), decode: decodeInto[kling.MotionControlParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.Kling.MotionControl.Create(ctx, params.(kling.MotionControlParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.Kling.MotionControl.Run(ctx, params.(kling.MotionControlParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.Kling.MotionControl.Get(ctx, id, opts...)
	}}
}
func newFluxKontextTextToImageSpec() actionSpec {
	return actionSpec{service: "flux-kontext", action: "text-to-image", isAsync: true, inputFields: inputFieldsFor[fluxkontext.TextToImageParams](), decode: decodeInto[fluxkontext.TextToImageParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.FluxKontext.TextToImage.Create(ctx, params.(fluxkontext.TextToImageParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.FluxKontext.TextToImage.Run(ctx, params.(fluxkontext.TextToImageParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.FluxKontext.TextToImage.Get(ctx, id, opts...)
	}}
}
func newFlux2TextToImageSpec() actionSpec {
	return actionSpec{service: "flux-2", action: "text-to-image", isAsync: true, inputFields: inputFieldsFor[flux2.TextToImageParams](), decode: decodeInto[flux2.TextToImageParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.Flux2.TextToImage.Create(ctx, params.(flux2.TextToImageParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.Flux2.TextToImage.Run(ctx, params.(flux2.TextToImageParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.Flux2.TextToImage.Get(ctx, id, opts...)
	}}
}
func newFlux2RemixImageSpec() actionSpec {
	return actionSpec{service: "flux-2", action: "remix-image", isAsync: true, inputFields: inputFieldsFor[flux2.RemixImageParams](), decode: decodeInto[flux2.RemixImageParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.Flux2.RemixImage.Create(ctx, params.(flux2.RemixImageParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.Flux2.RemixImage.Run(ctx, params.(flux2.RemixImageParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.Flux2.RemixImage.Get(ctx, id, opts...)
	}}
}

func newGeminiOmniCreateAudioSpec() actionSpec {
	return actionSpec{service: "gemini-omni", action: "create-audio", isAsync: false, inputFields: inputFieldsFor[geminiomni.CreateAudioParams](), decode: decodeInto[geminiomni.CreateAudioParams], run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.GeminiOmni.CreateAudio.Run(ctx, params.(geminiomni.CreateAudioParams), opts...)
	}}
}

func newGeminiOmniCreateCharacterSpec() actionSpec {
	return actionSpec{service: "gemini-omni", action: "create-character", isAsync: false, inputFields: inputFieldsFor[geminiomni.CreateCharacterParams](), decode: decodeInto[geminiomni.CreateCharacterParams], run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.GeminiOmni.CreateCharacter.Run(ctx, params.(geminiomni.CreateCharacterParams), opts...)
	}}
}

func newGeminiOmniTextToVideoSpec() actionSpec {
	return actionSpec{service: "gemini-omni", action: "text-to-video", isAsync: true, inputFields: inputFieldsFor[geminiomni.TextToVideoParams](), decode: decodeInto[geminiomni.TextToVideoParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.GeminiOmni.TextToVideo.Create(ctx, params.(geminiomni.TextToVideoParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.GeminiOmni.TextToVideo.Run(ctx, params.(geminiomni.TextToVideoParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.GeminiOmni.TextToVideo.Get(ctx, id, opts...)
	}}
}

func newQwen2TextToImageSpec() actionSpec {
	return actionSpec{service: "qwen-2", action: "text-to-image", isAsync: true, inputFields: inputFieldsFor[qwen2.TextToImageParams](), decode: decodeInto[qwen2.TextToImageParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.Qwen2.TextToImage.Create(ctx, params.(qwen2.TextToImageParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.Qwen2.TextToImage.Run(ctx, params.(qwen2.TextToImageParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.Qwen2.TextToImage.Get(ctx, id, opts...)
	}}
}
func newQwen2RemixImageSpec() actionSpec {
	return actionSpec{service: "qwen-2", action: "remix-image", isAsync: true, inputFields: inputFieldsFor[qwen2.RemixImageParams](), decode: decodeInto[qwen2.RemixImageParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.Qwen2.RemixImage.Create(ctx, params.(qwen2.RemixImageParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.Qwen2.RemixImage.Run(ctx, params.(qwen2.RemixImageParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.Qwen2.RemixImage.Get(ctx, id, opts...)
	}}
}
func newQwen2EditSpec() actionSpec {
	return actionSpec{service: "qwen-2", action: "edit-image", isAsync: true, inputFields: inputFieldsFor[qwen2.EditImageParams](), decode: decodeInto[qwen2.EditImageParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.Qwen2.EditImage.Create(ctx, params.(qwen2.EditImageParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.Qwen2.EditImage.Run(ctx, params.(qwen2.EditImageParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.Qwen2.EditImage.Get(ctx, id, opts...)
	}}
}
func newRecraftUpscaleSpec() actionSpec {
	return actionSpec{service: "recraft", action: "upscale-image", isAsync: true, inputFields: inputFieldsFor[recraft.UpscaleImageParams](), decode: decodeInto[recraft.UpscaleImageParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.Recraft.UpscaleImage.Create(ctx, params.(recraft.UpscaleImageParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.Recraft.UpscaleImage.Run(ctx, params.(recraft.UpscaleImageParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.Recraft.UpscaleImage.Get(ctx, id, opts...)
	}}
}
func newRecraftBackgroundRemovalSpec() actionSpec {
	return actionSpec{service: "recraft", action: "remove-background", isAsync: true, inputFields: inputFieldsFor[recraft.RemoveBackgroundParams](), decode: decodeInto[recraft.RemoveBackgroundParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.Recraft.RemoveBackground.Create(ctx, params.(recraft.RemoveBackgroundParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.Recraft.RemoveBackground.Run(ctx, params.(recraft.RemoveBackgroundParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.Recraft.RemoveBackground.Get(ctx, id, opts...)
	}}
}
func newZImageTextToImageSpec() actionSpec {
	return actionSpec{service: "z-image", action: "text-to-image", isAsync: true, inputFields: inputFieldsFor[zimage.TextToImageParams](), decode: decodeInto[zimage.TextToImageParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.ZImage.TextToImage.Create(ctx, params.(zimage.TextToImageParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.ZImage.TextToImage.Run(ctx, params.(zimage.TextToImageParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.ZImage.TextToImage.Get(ctx, id, opts...)
	}}
}
func newIdeogramV3TextToImageSpec() actionSpec {
	return actionSpec{service: "ideogram-v3", action: "text-to-image", isAsync: true, inputFields: inputFieldsFor[ideogramv3.TextToImageParams](), decode: decodeInto[ideogramv3.TextToImageParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.IdeogramV3.TextToImage.Create(ctx, params.(ideogramv3.TextToImageParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.IdeogramV3.TextToImage.Run(ctx, params.(ideogramv3.TextToImageParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.IdeogramV3.TextToImage.Get(ctx, id, opts...)
	}}
}
func newIdeogramV3EditImageSpec() actionSpec {
	return actionSpec{service: "ideogram-v3", action: "edit-image", isAsync: true, inputFields: inputFieldsFor[ideogramv3.EditImageParams](), decode: decodeInto[ideogramv3.EditImageParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.IdeogramV3.EditImage.Create(ctx, params.(ideogramv3.EditImageParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.IdeogramV3.EditImage.Run(ctx, params.(ideogramv3.EditImageParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.IdeogramV3.EditImage.Get(ctx, id, opts...)
	}}
}
func newIdeogramV3RemixImageSpec() actionSpec {
	return actionSpec{service: "ideogram-v3", action: "remix-image", isAsync: true, inputFields: inputFieldsFor[ideogramv3.RemixImageParams](), decode: decodeInto[ideogramv3.RemixImageParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.IdeogramV3.RemixImage.Create(ctx, params.(ideogramv3.RemixImageParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.IdeogramV3.RemixImage.Run(ctx, params.(ideogramv3.RemixImageParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.IdeogramV3.RemixImage.Get(ctx, id, opts...)
	}}
}
func newIdeogramV3ReframeImageSpec() actionSpec {
	return actionSpec{service: "ideogram-v3", action: "reframe-image", isAsync: true, inputFields: inputFieldsFor[ideogramv3.ReframeImageParams](), decode: decodeInto[ideogramv3.ReframeImageParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.IdeogramV3.ReframeImage.Create(ctx, params.(ideogramv3.ReframeImageParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.IdeogramV3.ReframeImage.Run(ctx, params.(ideogramv3.ReframeImageParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.IdeogramV3.ReframeImage.Get(ctx, id, opts...)
	}}
}

func newElevenlabsSpeechSpec() actionSpec {
	return actionSpec{service: "elevenlabs", action: "text-to-speech", isAsync: true, inputFields: inputFieldsFor[elevenlabs.TextToSpeechParams](), decode: decodeInto[elevenlabs.TextToSpeechParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.Elevenlabs.TextToSpeech.Create(ctx, params.(elevenlabs.TextToSpeechParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.Elevenlabs.TextToSpeech.Run(ctx, params.(elevenlabs.TextToSpeechParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.Elevenlabs.TextToSpeech.Get(ctx, id, opts...)
	}}
}

func newElevenlabsDialogueSpec() actionSpec {
	return actionSpec{service: "elevenlabs", action: "text-to-dialogue", isAsync: true, inputFields: inputFieldsFor[elevenlabs.TextToDialogueParams](), decode: decodeInto[elevenlabs.TextToDialogueParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.Elevenlabs.TextToDialogue.Create(ctx, params.(elevenlabs.TextToDialogueParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.Elevenlabs.TextToDialogue.Run(ctx, params.(elevenlabs.TextToDialogueParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.Elevenlabs.TextToDialogue.Get(ctx, id, opts...)
	}}
}

func newElevenlabsSoundEffectSpec() actionSpec {
	return actionSpec{service: "elevenlabs", action: "text-to-sound", isAsync: true, inputFields: inputFieldsFor[elevenlabs.TextToSoundParams](), decode: decodeInto[elevenlabs.TextToSoundParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.Elevenlabs.TextToSound.Create(ctx, params.(elevenlabs.TextToSoundParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.Elevenlabs.TextToSound.Run(ctx, params.(elevenlabs.TextToSoundParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.Elevenlabs.TextToSound.Get(ctx, id, opts...)
	}}
}

func newElevenlabsTranscriptionSpec() actionSpec {
	return actionSpec{service: "elevenlabs", action: "speech-to-text", isAsync: true, inputFields: inputFieldsFor[elevenlabs.SpeechToTextParams](), decode: decodeInto[elevenlabs.SpeechToTextParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.Elevenlabs.SpeechToText.Create(ctx, params.(elevenlabs.SpeechToTextParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.Elevenlabs.SpeechToText.Run(ctx, params.(elevenlabs.SpeechToTextParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.Elevenlabs.SpeechToText.Get(ctx, id, opts...)
	}}
}

func newElevenlabsAudioIsolationSpec() actionSpec {
	return actionSpec{service: "elevenlabs", action: "isolate-audio", isAsync: true, inputFields: inputFieldsFor[elevenlabs.IsolateAudioParams](), decode: decodeInto[elevenlabs.IsolateAudioParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.Elevenlabs.IsolateAudio.Create(ctx, params.(elevenlabs.IsolateAudioParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.Elevenlabs.IsolateAudio.Run(ctx, params.(elevenlabs.IsolateAudioParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.Elevenlabs.IsolateAudio.Get(ctx, id, opts...)
	}}
}

func newInfiniteTalkAudioToVideoSpec() actionSpec {
	return actionSpec{service: "infinitetalk", action: "audio-to-video", isAsync: true, inputFields: inputFieldsFor[infinitetalk.AudioToVideoParams](), decode: decodeInto[infinitetalk.AudioToVideoParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.InfiniteTalk.AudioToVideo.Create(ctx, params.(infinitetalk.AudioToVideoParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.InfiniteTalk.AudioToVideo.Run(ctx, params.(infinitetalk.AudioToVideoParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.InfiniteTalk.AudioToVideo.Get(ctx, id, opts...)
	}}
}

func newOmniHumanAudioToVideoSpec() actionSpec {
	return actionSpec{service: "omnihuman", action: "audio-to-video", isAsync: true, inputFields: inputFieldsFor[omnihuman.AudioToVideoParams](), decode: decodeInto[omnihuman.AudioToVideoParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.OmniHuman.AudioToVideo.Create(ctx, params.(omnihuman.AudioToVideoParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.OmniHuman.AudioToVideo.Run(ctx, params.(omnihuman.AudioToVideoParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.OmniHuman.AudioToVideo.Get(ctx, id, opts...)
	}}
}

func newOmniHumanHumanIdentificationSpec() actionSpec {
	return actionSpec{service: "omnihuman", action: "human-identification", isAsync: true, inputFields: inputFieldsFor[omnihuman.HumanIdentificationParams](), decode: decodeInto[omnihuman.HumanIdentificationParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.OmniHuman.HumanIdentification.Create(ctx, params.(omnihuman.HumanIdentificationParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.OmniHuman.HumanIdentification.Run(ctx, params.(omnihuman.HumanIdentificationParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.OmniHuman.HumanIdentification.Get(ctx, id, opts...)
	}}
}

func newOmniHumanSubjectDetectionSpec() actionSpec {
	return actionSpec{service: "omnihuman", action: "subject-detection", isAsync: true, inputFields: inputFieldsFor[omnihuman.SubjectDetectionParams](), decode: decodeInto[omnihuman.SubjectDetectionParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.OmniHuman.SubjectDetection.Create(ctx, params.(omnihuman.SubjectDetectionParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.OmniHuman.SubjectDetection.Run(ctx, params.(omnihuman.SubjectDetectionParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.OmniHuman.SubjectDetection.Get(ctx, id, opts...)
	}}
}

func newWanTextToVideoSpec() actionSpec {
	return actionSpec{service: "wan", action: "text-to-video", isAsync: true, inputFields: inputFieldsFor[wan.TextToVideoParams](), decode: decodeInto[wan.TextToVideoParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.Wan.TextToVideo.Create(ctx, params.(wan.TextToVideoParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.Wan.TextToVideo.Run(ctx, params.(wan.TextToVideoParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.Wan.TextToVideo.Get(ctx, id, opts...)
	}}
}

func newWanImageToVideoSpec() actionSpec {
	return actionSpec{service: "wan", action: "image-to-video", isAsync: true, inputFields: inputFieldsFor[wan.ImageToVideoParams](), decode: decodeInto[wan.ImageToVideoParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.Wan.ImageToVideo.Create(ctx, params.(wan.ImageToVideoParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.Wan.ImageToVideo.Run(ctx, params.(wan.ImageToVideoParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.Wan.ImageToVideo.Get(ctx, id, opts...)
	}}
}

func newWanSpeechToVideoSpec() actionSpec {
	return actionSpec{service: "wan", action: "speech-to-video", isAsync: true, inputFields: inputFieldsFor[wan.SpeechToVideoParams](), decode: decodeInto[wan.SpeechToVideoParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.Wan.SpeechToVideo.Create(ctx, params.(wan.SpeechToVideoParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.Wan.SpeechToVideo.Run(ctx, params.(wan.SpeechToVideoParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.Wan.SpeechToVideo.Get(ctx, id, opts...)
	}}
}

func newWanAnimateSpec() actionSpec {
	return actionSpec{service: "wan", action: "animate", isAsync: true, inputFields: inputFieldsFor[wan.AnimateParams](), decode: decodeInto[wan.AnimateParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.Wan.Animate.Create(ctx, params.(wan.AnimateParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.Wan.Animate.Run(ctx, params.(wan.AnimateParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.Wan.Animate.Get(ctx, id, opts...)
	}}
}

func newWanTextToImageSpec() actionSpec {
	return actionSpec{service: "wan", action: "text-to-image", isAsync: true, inputFields: inputFieldsFor[wan.TextToImageParams](), decode: decodeInto[wan.TextToImageParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.Wan.TextToImage.Create(ctx, params.(wan.TextToImageParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.Wan.TextToImage.Run(ctx, params.(wan.TextToImageParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.Wan.TextToImage.Get(ctx, id, opts...)
	}}
}

func newWanEditVideoSpec() actionSpec {
	return actionSpec{service: "wan", action: "edit-video", isAsync: true, inputFields: inputFieldsFor[wan.EditVideoParams](), decode: decodeInto[wan.EditVideoParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.Wan.EditVideo.Create(ctx, params.(wan.EditVideoParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.Wan.EditVideo.Run(ctx, params.(wan.EditVideoParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.Wan.EditVideo.Get(ctx, id, opts...)
	}}
}

func newLumaModifySpec() actionSpec {
	return actionSpec{service: "luma", action: "modify-video", isAsync: true, inputFields: inputFieldsFor[luma.ModifyVideoParams](), decode: decodeInto[luma.ModifyVideoParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.Luma.ModifyVideo.Create(ctx, params.(luma.ModifyVideoParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.Luma.ModifyVideo.Run(ctx, params.(luma.ModifyVideoParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.Luma.ModifyVideo.Get(ctx, id, opts...)
	}}
}

func newGptImageTextToImageSpec() actionSpec {
	return actionSpec{service: "gpt-image", action: "text-to-image", isAsync: true, inputFields: inputFieldsFor[gptimage.TextToImageParams](), decode: decodeInto[gptimage.TextToImageParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.GptImage.TextToImage.Create(ctx, params.(gptimage.TextToImageParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.GptImage.TextToImage.Run(ctx, params.(gptimage.TextToImageParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.GptImage.TextToImage.Get(ctx, id, opts...)
	}}
}

func newGptImageEditImageSpec() actionSpec {
	return actionSpec{service: "gpt-image", action: "edit-image", isAsync: true, inputFields: inputFieldsFor[gptimage.EditImageParams](), decode: decodeInto[gptimage.EditImageParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.GptImage.EditImage.Create(ctx, params.(gptimage.EditImageParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.GptImage.EditImage.Run(ctx, params.(gptimage.EditImageParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.GptImage.EditImage.Get(ctx, id, opts...)
	}}
}

func newGptImage2TextToImageSpec() actionSpec {
	return actionSpec{service: "gpt-image-2", action: "text-to-image", isAsync: true, inputFields: inputFieldsFor[gptimage2.TextToImageParams](), decode: decodeInto[gptimage2.TextToImageParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.GptImage2.TextToImage.Create(ctx, params.(gptimage2.TextToImageParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.GptImage2.TextToImage.Run(ctx, params.(gptimage2.TextToImageParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.GptImage2.TextToImage.Get(ctx, id, opts...)
	}}
}

func newGptImage2EditImageSpec() actionSpec {
	return actionSpec{service: "gpt-image-2", action: "edit-image", isAsync: true, inputFields: inputFieldsFor[gptimage2.EditImageParams](), decode: decodeInto[gptimage2.EditImageParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.GptImage2.EditImage.Create(ctx, params.(gptimage2.EditImageParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.GptImage2.EditImage.Run(ctx, params.(gptimage2.EditImageParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.GptImage2.EditImage.Get(ctx, id, opts...)
	}}
}

func newGpt4oImageTextToImageSpec() actionSpec {
	return actionSpec{service: "gpt-4o-image", action: "text-to-image", isAsync: true, inputFields: inputFieldsFor[gpt4oimage.TextToImageParams](), decode: decodeInto[gpt4oimage.TextToImageParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.Gpt4oImage.TextToImage.Create(ctx, params.(gpt4oimage.TextToImageParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.Gpt4oImage.TextToImage.Run(ctx, params.(gpt4oimage.TextToImageParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.Gpt4oImage.TextToImage.Get(ctx, id, opts...)
	}}
}

func newHailuoTextToVideoSpec() actionSpec {
	return actionSpec{service: "hailuo", action: "text-to-video", isAsync: true, inputFields: inputFieldsFor[hailuo.TextToVideoParams](), decode: decodeInto[hailuo.TextToVideoParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.Hailuo.TextToVideo.Create(ctx, params.(hailuo.TextToVideoParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.Hailuo.TextToVideo.Run(ctx, params.(hailuo.TextToVideoParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.Hailuo.TextToVideo.Get(ctx, id, opts...)
	}}
}

func newHailuoImageToVideoSpec() actionSpec {
	return actionSpec{service: "hailuo", action: "image-to-video", isAsync: true, inputFields: inputFieldsFor[hailuo.ImageToVideoParams](), decode: decodeInto[hailuo.ImageToVideoParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.Hailuo.ImageToVideo.Create(ctx, params.(hailuo.ImageToVideoParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.Hailuo.ImageToVideo.Run(ctx, params.(hailuo.ImageToVideoParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.Hailuo.ImageToVideo.Get(ctx, id, opts...)
	}}
}

func newVolcengineLipSyncVideoSpec() actionSpec {
	return actionSpec{service: "volcengine-lip-sync", action: "lip-sync-video", isAsync: true, inputFields: inputFieldsFor[volcenginelipsync.LipSyncVideoParams](), decode: decodeInto[volcenginelipsync.LipSyncVideoParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.VolcengineLipSync.LipSyncVideo.Create(ctx, params.(volcenginelipsync.LipSyncVideoParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.VolcengineLipSync.LipSyncVideo.Run(ctx, params.(volcenginelipsync.LipSyncVideoParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.VolcengineLipSync.LipSyncVideo.Get(ctx, id, opts...)
	}}
}

func newHappyHorseTextToVideoSpec() actionSpec {
	return actionSpec{service: "happyhorse", action: "text-to-video", isAsync: true, inputFields: inputFieldsFor[happyhorse.TextToVideoParams](), decode: decodeInto[happyhorse.TextToVideoParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.HappyHorse.TextToVideo.Create(ctx, params.(happyhorse.TextToVideoParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.HappyHorse.TextToVideo.Run(ctx, params.(happyhorse.TextToVideoParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.HappyHorse.TextToVideo.Get(ctx, id, opts...)
	}}
}

func newHappyHorseImageToVideoSpec() actionSpec {
	return actionSpec{service: "happyhorse", action: "image-to-video", isAsync: true, inputFields: inputFieldsFor[happyhorse.ImageToVideoParams](), decode: decodeInto[happyhorse.ImageToVideoParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.HappyHorse.ImageToVideo.Create(ctx, params.(happyhorse.ImageToVideoParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.HappyHorse.ImageToVideo.Run(ctx, params.(happyhorse.ImageToVideoParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.HappyHorse.ImageToVideo.Get(ctx, id, opts...)
	}}
}

func newHappyHorseEditVideoSpec() actionSpec {
	return actionSpec{service: "happyhorse", action: "edit-video", isAsync: true, inputFields: inputFieldsFor[happyhorse.EditVideoParams](), decode: decodeInto[happyhorse.EditVideoParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.HappyHorse.EditVideo.Create(ctx, params.(happyhorse.EditVideoParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.HappyHorse.EditVideo.Run(ctx, params.(happyhorse.EditVideoParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.HappyHorse.EditVideo.Get(ctx, id, opts...)
	}}
}

func newGrokImagineTextToVideoSpec() actionSpec {
	return actionSpec{service: "grok-imagine", action: "text-to-video", isAsync: true, inputFields: inputFieldsFor[grokimagine.TextToVideoParams](), decode: decodeInto[grokimagine.TextToVideoParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.GrokImagine.TextToVideo.Create(ctx, params.(grokimagine.TextToVideoParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.GrokImagine.TextToVideo.Run(ctx, params.(grokimagine.TextToVideoParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.GrokImagine.TextToVideo.Get(ctx, id, opts...)
	}}
}

func newGrokImagineImageToVideoSpec() actionSpec {
	return actionSpec{service: "grok-imagine", action: "image-to-video", isAsync: true, inputFields: inputFieldsFor[grokimagine.ImageToVideoParams](), decode: decodeInto[grokimagine.ImageToVideoParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.GrokImagine.ImageToVideo.Create(ctx, params.(grokimagine.ImageToVideoParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.GrokImagine.ImageToVideo.Run(ctx, params.(grokimagine.ImageToVideoParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.GrokImagine.ImageToVideo.Get(ctx, id, opts...)
	}}
}

func newGrokImagineTextToImageSpec() actionSpec {
	return actionSpec{service: "grok-imagine", action: "text-to-image", isAsync: true, inputFields: inputFieldsFor[grokimagine.TextToImageParams](), decode: decodeInto[grokimagine.TextToImageParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.GrokImagine.TextToImage.Create(ctx, params.(grokimagine.TextToImageParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.GrokImagine.TextToImage.Run(ctx, params.(grokimagine.TextToImageParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.GrokImagine.TextToImage.Get(ctx, id, opts...)
	}}
}

func newGrokImagineEditImageSpec() actionSpec {
	return actionSpec{service: "grok-imagine", action: "edit-image", isAsync: true, inputFields: inputFieldsFor[grokimagine.EditImageParams](), decode: decodeInto[grokimagine.EditImageParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.GrokImagine.EditImage.Create(ctx, params.(grokimagine.EditImageParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.GrokImagine.EditImage.Run(ctx, params.(grokimagine.EditImageParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.GrokImagine.EditImage.Get(ctx, id, opts...)
	}}
}

func newGrokImagineExtendSpec() actionSpec {
	return actionSpec{service: "grok-imagine", action: "extend", isAsync: true, inputFields: inputFieldsFor[grokimagine.ExtendParams](), decode: decodeInto[grokimagine.ExtendParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.GrokImagine.Extensions.Create(ctx, params.(grokimagine.ExtendParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.GrokImagine.Extensions.Run(ctx, params.(grokimagine.ExtendParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.GrokImagine.Extensions.Get(ctx, id, opts...)
	}}
}

func newGrokImagineUpscaleSpec() actionSpec {
	return actionSpec{service: "grok-imagine", action: "upscale-image", isAsync: true, inputFields: inputFieldsFor[grokimagine.UpscaleParams](), decode: decodeInto[grokimagine.UpscaleParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.GrokImagine.Upscales.Create(ctx, params.(grokimagine.UpscaleParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.GrokImagine.Upscales.Run(ctx, params.(grokimagine.UpscaleParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.GrokImagine.Upscales.Get(ctx, id, opts...)
	}}
}

func newTopazUpscaleImageSpec() actionSpec {
	return actionSpec{service: "topaz", action: "upscale-image", isAsync: true, inputFields: inputFieldsFor[topaz.UpscaleImageParams](), decode: decodeInto[topaz.UpscaleImageParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.Topaz.UpscaleImage.Create(ctx, params.(topaz.UpscaleImageParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.Topaz.UpscaleImage.Run(ctx, params.(topaz.UpscaleImageParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.Topaz.UpscaleImage.Get(ctx, id, opts...)
	}}
}

func newTopazUpscaleVideoSpec() actionSpec {
	return actionSpec{service: "topaz", action: "upscale-video", isAsync: true, inputFields: inputFieldsFor[topaz.UpscaleVideoParams](), decode: decodeInto[topaz.UpscaleVideoParams], create: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (*core.TaskCreateResponse, error) {
		return client.Topaz.UpscaleVideo.Create(ctx, params.(topaz.UpscaleVideoParams), opts...)
	}, run: func(ctx context.Context, client *runapi.Client, params any, opts []option.RequestOption) (any, error) {
		return client.Topaz.UpscaleVideo.Run(ctx, params.(topaz.UpscaleVideoParams), opts...)
	}, get: func(ctx context.Context, client *runapi.Client, id string, opts []option.RequestOption) (core.TaskResponse, error) {
		return client.Topaz.UpscaleVideo.Get(ctx, id, opts...)
	}}
}
