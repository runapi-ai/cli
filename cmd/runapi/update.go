package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	runapi "github.com/runapi-ai/cli/internal/runapi"
	"github.com/runapi-ai/core-sdk/go/core"
	"github.com/spf13/cobra"
)

const (
	updateCheckInterval = 24 * time.Hour
	updateCheckTimeout  = time.Second
)

type updateCheckState struct {
	CheckedAt time.Time `json:"checked_at"`
}

type releaseManifest struct {
	Version string `json:"version"`
}

func (c *cli) maybeNotifyUpdate(command *cobra.Command, args []string) {
	if c.quiet || command == nil || command.Name() == "version" || helpOnlyCommand(command, args) {
		return
	}
	if c.stderrTTY == nil || !c.stderrTTY() || !validReleaseVersion(runapi.Version) {
		return
	}

	now := time.Now()
	statePath, err := updateStatePath()
	if err != nil {
		return
	}
	if state, err := loadUpdateCheckState(statePath); err == nil && now.Sub(state.CheckedAt) >= 0 && now.Sub(state.CheckedAt) < updateCheckInterval {
		return
	}

	latest, fetchErr := c.fetchLatestCLIVersion()
	_ = saveUpdateCheckState(statePath, updateCheckState{CheckedAt: now})
	if fetchErr != nil || !newerReleaseVersion(latest, runapi.Version) {
		return
	}

	fmt.Fprintf(c.stderr, "\nRunAPI CLI %s is available (current: %s).\nUpdate with: %s\n", latest, runapi.Version, updateCommand())
}

const (
	brewUpdateCommand    = "brew upgrade runapi-ai/tap/runapi"
	goUpdateCommand      = "go install github.com/runapi-ai/cli/cmd/runapi@latest"
	shUpdateCommand      = "curl -fsSL https://runapi.ai/cli/install.sh | sh"
	windowsUpdateCommand = "Download the matching Windows archive from https://github.com/runapi-ai/cli/releases/latest"
)

func helpOnlyCommand(command *cobra.Command, args []string) bool {
	if command.Name() == "help" || command.Flags().Changed("help") {
		return true
	}
	// Cobra renders usage for the root and non-runnable command groups when no
	// runnable command is selected, including `runapi` and `runapi agent`.
	return len(args) == 0 || (command.Run == nil && command.RunE == nil && len(command.Commands()) > 0)
}

func updateCommand() string {
	executable, err := os.Executable()
	if err != nil {
		return shUpdateCommand
	}
	if resolved, err := filepath.EvalSymlinks(executable); err == nil {
		executable = resolved
	}
	return updateCommandForPath(executable)
}

func updateCommandForPath(executable string) string {
	if resolved, err := filepath.EvalSymlinks(executable); err == nil {
		executable = resolved
	}
	path := normalizeExecutablePath(executable)
	if strings.Contains(path, "/Cellar/") || strings.Contains(path, "/Linuxbrew/Cellar/") {
		return brewUpdateCommand
	}

	if isGoInstallPath(path) {
		return goUpdateCommand
	}
	if strings.HasSuffix(strings.ToLower(path), "/runapi.exe") {
		return windowsUpdateCommand
	}
	if runtime.GOOS == "windows" {
		return windowsUpdateCommand
	}
	if strings.HasSuffix(path, "/runapi") {
		return shUpdateCommand + " -s -- --dir " + shellQuote(filepath.Dir(filepath.FromSlash(path)))
	}
	return shUpdateCommand
}

func isGoInstallPath(executable string) bool {
	path := normalizeExecutablePath(executable)
	name := strings.ToLower(filepath.Base(path))
	if name != "runapi" && name != "runapi.exe" {
		return false
	}
	for _, root := range goInstallRoots() {
		rootPath := normalizeExecutablePath(root)
		if samePath(path, filepath.Join(rootPath, name)) || strings.EqualFold(path, normalizeExecutablePath(filepath.Join(root, name))) {
			return true
		}
	}
	return strings.HasSuffix(path, "/go/bin/"+name)
}

func normalizeExecutablePath(path string) string {
	return strings.TrimRight(strings.ReplaceAll(filepath.ToSlash(path), "\\", "/"), "/")
}

func shellQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "'\\''") + "'"
}

func goInstallRoots() []string {
	roots := make([]string, 0, 3)
	if gobin := strings.TrimSpace(os.Getenv("GOBIN")); gobin != "" {
		roots = append(roots, gobin)
	}
	if gopath := strings.TrimSpace(os.Getenv("GOPATH")); gopath != "" {
		for _, root := range filepath.SplitList(gopath) {
			if strings.TrimSpace(root) != "" {
				roots = append(roots, filepath.Join(root, "bin"))
			}
		}
	}
	return roots
}

func samePath(left, right string) bool {
	left, right = filepath.Clean(left), filepath.Clean(right)
	if left == right {
		return true
	}
	leftInfo, leftErr := os.Stat(left)
	rightInfo, rightErr := os.Stat(right)
	return leftErr == nil && rightErr == nil && os.SameFile(leftInfo, rightInfo)
}

func (c *cli) fetchLatestCLIVersion() (string, error) {
	baseURL := strings.TrimRight(c.resolvedUpdateBaseURL(), "/")
	ctx, cancel := context.WithTimeout(context.Background(), updateCheckTimeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/cli/latest.json", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", core.CLIUserAgent(runapi.Version))

	client := c.httpClient
	if client == nil {
		client = &http.Client{}
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("update check returned HTTP %d", resp.StatusCode)
	}

	var manifest releaseManifest
	if err := json.NewDecoder(io.LimitReader(resp.Body, 64<<10)).Decode(&manifest); err != nil {
		return "", err
	}
	if !validReleaseVersion(manifest.Version) {
		return "", fmt.Errorf("invalid latest CLI version")
	}
	return manifest.Version, nil
}

func (c *cli) resolvedUpdateBaseURL() string {
	if value := strings.TrimSpace(c.updateBaseURL); value != "" {
		return value
	}
	if value := strings.TrimSpace(c.baseURLFlag); value != "" {
		return value
	}
	if value := strings.TrimSpace(os.Getenv("RUNAPI_BASE_URL")); value != "" {
		return value
	}
	if cfg, err := loadConfig(); err == nil && strings.TrimSpace(cfg.BaseURL) != "" {
		return cfg.BaseURL
	}
	return core.DefaultBaseURL
}

func updateStatePath() (string, error) {
	configPath, err := configFilePath()
	if err != nil {
		return "", err
	}
	return filepath.Join(filepath.Dir(configPath), "update-check.json"), nil
}

func loadUpdateCheckState(path string) (updateCheckState, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return updateCheckState{}, err
	}
	var state updateCheckState
	err = json.Unmarshal(data, &state)
	return state, err
}

func saveUpdateCheckState(path string, state updateCheckState) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o600)
}

func newerReleaseVersion(candidate, current string) bool {
	candidateParts, candidateOK := releaseVersionParts(candidate)
	currentParts, currentOK := releaseVersionParts(current)
	if !candidateOK || !currentOK {
		return false
	}
	for i := range candidateParts {
		if candidateParts[i] != currentParts[i] {
			return candidateParts[i] > currentParts[i]
		}
	}
	return false
}

func validReleaseVersion(version string) bool {
	_, ok := releaseVersionParts(version)
	return ok
}

func releaseVersionParts(version string) ([3]int, bool) {
	var parsed [3]int
	parts := strings.Split(strings.TrimPrefix(version, "v"), ".")
	if len(parts) != len(parsed) {
		return parsed, false
	}
	for i, part := range parts {
		if part == "" || (len(part) > 1 && part[0] == '0') {
			return parsed, false
		}
		value, err := strconv.Atoi(part)
		if err != nil || value < 0 {
			return parsed, false
		}
		parsed[i] = value
	}
	return parsed, true
}
