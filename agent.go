package main

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	runapi "github.com/runapi-ai/cli/internal/runapi"
	"github.com/runapi-ai/core-sdk/go/core"
	"github.com/spf13/cobra"
)

const (
	defaultSkillSourceRepo = "runapi-ai/cli-skill"
	skillName              = "cli"
	skillSubdir            = "skills/cli/"
	skillDownloadTimeout   = 60 * time.Second
)

// Hard ceiling on the total decompressed bytes we will write while
// extracting the skill tarball. The source repo is small (well under
// 1 MB at the time of writing); 50 MB is a comfortable upper bound
// that still trips long before a real zip bomb causes harm. Declared
// as var so tests can lower the cap.
var skillMaxExtractBytes int64 = 50 * 1024 * 1024

func setSkillMaxExtractBytesForTest(n int64) { skillMaxExtractBytes = n }

// skillTargets maps a built-in target name to a path under the user's HOME
// where the skill directory should be created. The skill directory itself
// is always named "cli" inside that path (e.g. ~/.claude/skills/cli).
var skillTargets = map[string]string{
	"claude":   ".claude/skills",
	"codex":    ".agents/skills",
	"gemini":   ".gemini/skills",
	"openclaw": ".openclaw/skills",
	"hermes":   ".hermes/skills",
}

func supportedTargetNames() []string {
	return []string{"claude", "codex", "gemini", "openclaw", "hermes"}
}

func (c *cli) agentCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agent",
		Short: "Manage RunAPI skill installation for agent runtimes",
	}
	cmd.AddCommand(c.agentInstallSkillCommand())
	cmd.AddCommand(c.agentUninstallSkillCommand())
	cmd.AddCommand(c.agentListTargetsCommand())
	return cmd
}

type installSkillOpts struct {
	target    string
	targetDir string
	version   string
	source    string
	force     bool
}

func (c *cli) agentInstallSkillCommand() *cobra.Command {
	o := installSkillOpts{source: defaultSkillSourceRepo}
	cmd := &cobra.Command{
		Use:   "install-skill",
		Short: "Install the RunAPI CLI skill into a supported agent runtime",
		Long: "install-skill downloads the canonical skills/cli/ directory from the runapi-ai/cli-skill\n" +
			"release archive and writes it to the runtime's skill directory.\n\n" +
			"Targets: claude, codex, gemini, openclaw, hermes. Use --target-dir for a custom path.",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.runInstallSkill(cmd.Context(), o)
		},
	}
	cmd.Flags().StringVar(&o.target, "target", "", "Built-in target: "+strings.Join(supportedTargetNames(), ", "))
	cmd.Flags().StringVar(&o.targetDir, "target-dir", "", "Custom install directory (overrides --target). The skill is placed in <dir>/cli.")
	cmd.Flags().StringVar(&o.version, "version", "", "Skill tag to install (e.g. v0.1.0). Default: matches CLI version.")
	cmd.Flags().StringVar(&o.source, "source", defaultSkillSourceRepo, "Source GitHub repo")
	cmd.Flags().BoolVar(&o.force, "force", false, "Overwrite an existing skill directory")
	return cmd
}

func (c *cli) agentUninstallSkillCommand() *cobra.Command {
	var target, targetDir string
	cmd := &cobra.Command{
		Use:   "uninstall-skill",
		Short: "Remove the RunAPI CLI skill from an agent runtime",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			installDir, err := resolveSkillTargetDir(target, targetDir)
			if err != nil {
				return err
			}
			destDir := filepath.Join(installDir, skillName)
			if _, err := os.Stat(destDir); os.IsNotExist(err) {
				c.logf("nothing to remove at %s", destDir)
				return c.writeJSON(map[string]any{"removed": false, "target_dir": destDir})
			}
			if err := os.RemoveAll(destDir); err != nil {
				return err
			}
			c.logf("removed %s", destDir)
			return c.writeJSON(map[string]any{"removed": true, "target_dir": destDir})
		},
	}
	cmd.Flags().StringVar(&target, "target", "", "Built-in target: "+strings.Join(supportedTargetNames(), ", "))
	cmd.Flags().StringVar(&targetDir, "target-dir", "", "Custom install directory")
	return cmd
}

func (c *cli) agentListTargetsCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list-targets",
		Short: "Print supported agent runtime targets as JSON",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			home, _ := os.UserHomeDir()
			targets := make(map[string]string, len(skillTargets))
			for k, rel := range skillTargets {
				targets[k] = filepath.Join(home, rel, skillName)
			}
			return c.writeJSON(map[string]any{
				"skill":   skillName,
				"targets": targets,
			})
		},
	}
}

func (c *cli) runInstallSkill(ctx context.Context, o installSkillOpts) error {
	installDir, err := resolveSkillTargetDir(o.target, o.targetDir)
	if err != nil {
		return err
	}

	tag, err := resolveSkillTag(o.version)
	if err != nil {
		return err
	}

	source := strings.TrimSpace(o.source)
	if source == "" {
		source = defaultSkillSourceRepo
	}
	if !isSafeRepoSlug(source) {
		return core.NewError(core.ErrValidation, fmt.Sprintf("invalid source repo: %q", source), 422, "", nil, nil)
	}

	destDir := filepath.Join(installDir, skillName)
	if !o.force {
		if _, err := os.Stat(destDir); err == nil {
			return core.NewError(core.ErrConflict,
				fmt.Sprintf("skill already installed at %s (pass --force to overwrite)", destDir),
				409, "", nil, nil)
		}
	} else if err := os.RemoveAll(destDir); err != nil {
		return err
	}
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return err
	}

	archiveURL := fmt.Sprintf("https://github.com/%s/archive/refs/tags/%s.tar.gz", source, tag)
	if c.archiveBaseURL != "" {
		archiveURL = fmt.Sprintf("%s/%s/archive/refs/tags/%s.tar.gz", strings.TrimRight(c.archiveBaseURL, "/"), source, tag)
	}
	if err := assertSecureArchiveURL(archiveURL); err != nil {
		return err
	}

	c.logf("downloading %s", archiveURL)
	body, err := httpDownload(ctx, c.httpClient, archiveURL)
	if err != nil {
		return err
	}
	defer body.Close()

	files, err := extractSkillFromTarball(body, destDir)
	if err != nil {
		return err
	}

	c.logf("✓ installed skill %q (%d files) to %s", skillName, files, destDir)
	return c.writeJSON(map[string]any{
		"installed":  true,
		"skill":      skillName,
		"target":     o.target,
		"target_dir": destDir,
		"version":    tag,
		"source":     source,
		"files":      files,
	})
}

func resolveSkillTargetDir(target, targetDir string) (string, error) {
	if strings.TrimSpace(targetDir) != "" {
		return expandUserPath(targetDir)
	}
	if strings.TrimSpace(target) == "" {
		return "", core.NewError(core.ErrValidation,
			"--target or --target-dir is required (use `agent list-targets` to see options)",
			422, "", nil, nil)
	}
	rel, ok := skillTargets[target]
	if !ok {
		return "", core.NewError(core.ErrValidation,
			fmt.Sprintf("unknown target %q (valid: %s)", target, strings.Join(supportedTargetNames(), ", ")),
			422, "", nil, nil)
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, rel), nil
}

func resolveSkillTag(version string) (string, error) {
	v := strings.TrimSpace(version)
	if v != "" {
		if !strings.HasPrefix(v, "v") {
			v = "v" + v
		}
		if !isSafeTag(v) {
			return "", core.NewError(core.ErrValidation, fmt.Sprintf("invalid version %q", version), 422, "", nil, nil)
		}
		return v, nil
	}
	if runapi.Version == "" || strings.HasSuffix(runapi.Version, "-dev") {
		return "", core.NewError(core.ErrValidation,
			fmt.Sprintf("CLI is a dev build (%q); pass --version <tag> to pin the skill", runapi.Version),
			422, "", nil, nil)
	}
	return "v" + runapi.Version, nil
}

func expandUserPath(p string) (string, error) {
	if strings.HasPrefix(p, "~/") || p == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		if p == "~" {
			return home, nil
		}
		return filepath.Join(home, strings.TrimPrefix(p, "~/")), nil
	}
	abs, err := filepath.Abs(p)
	if err != nil {
		return "", err
	}
	return abs, nil
}

// isSafeRepoSlug accepts "<owner>/<name>" with conservative chars.
func isSafeRepoSlug(s string) bool {
	parts := strings.Split(s, "/")
	if len(parts) != 2 {
		return false
	}
	for _, part := range parts {
		if part == "" {
			return false
		}
		for _, r := range part {
			ok := (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' || r == '.'
			if !ok {
				return false
			}
		}
	}
	return true
}

// isSafeTag matches v<x>.<y>.<z>, optionally with a -<pre> suffix made of
// dot/dash/alphanumeric characters.
func isSafeTag(t string) bool {
	if !strings.HasPrefix(t, "v") {
		return false
	}
	rest := t[1:]
	if rest == "" {
		return false
	}
	for _, r := range rest {
		ok := (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '.' || r == '-' || r == '_'
		if !ok {
			return false
		}
	}
	_, err := url.Parse("https://example.com/" + t)
	return err == nil
}

// assertSecureArchiveURL enforces https:// for the skill archive download
// path. Loopback hosts (127.0.0.1, ::1, localhost) are permitted because
// tests inject an httptest.Server via cli.archiveBaseURL.
func assertSecureArchiveURL(rawURL string) error {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return core.NewError(core.ErrValidation, fmt.Sprintf("invalid archive URL: %v", err), 422, "", nil, err)
	}
	if parsed.Scheme == "https" {
		return nil
	}
	host := parsed.Hostname()
	if parsed.Scheme == "http" && (host == "127.0.0.1" || host == "::1" || host == "localhost") {
		return nil
	}
	return core.NewError(core.ErrValidation,
		fmt.Sprintf("refusing to fetch skill archive over %q (https required outside loopback)", parsed.Scheme),
		422, "", nil, nil)
}

func httpDownload(ctx context.Context, client *http.Client, rawURL string) (io.ReadCloser, error) {
	if client == nil {
		client = &http.Client{Timeout: skillDownloadTimeout}
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/octet-stream")
	resp, err := client.Do(req)
	if err != nil {
		return nil, core.NewError(core.ErrNetwork, fmt.Sprintf("download failed: %v", err), 0, "", nil, err)
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		resp.Body.Close()
		if resp.StatusCode == http.StatusNotFound {
			return nil, core.NewError(core.ErrNotFound, fmt.Sprintf("skill release not found: %s", rawURL), 404, "", nil, nil)
		}
		return nil, core.NewError(core.ErrServer, fmt.Sprintf("download failed: HTTP %d %s", resp.StatusCode, string(body)), resp.StatusCode, "", nil, nil)
	}
	return resp.Body, nil
}

// extractSkillFromTarball pulls the `skills/cli/` subtree of a GitHub archive
// into destDir (so destDir/SKILL.md, destDir/references/..., etc.) and returns
// the number of regular files written.
func extractSkillFromTarball(r io.Reader, destDir string) (int, error) {
	gz, err := gzip.NewReader(r)
	if err != nil {
		return 0, fmt.Errorf("gzip: %w", err)
	}
	defer gz.Close()
	tr := tar.NewReader(gz)

	absDest, err := filepath.Abs(destDir)
	if err != nil {
		return 0, err
	}

	count := 0
	budget := skillMaxExtractBytes
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, fmt.Errorf("tar: %w", err)
		}
		// Strip the top-level archive directory (e.g. "cli-skill-0.1.0/").
		parts := strings.SplitN(hdr.Name, "/", 2)
		if len(parts) != 2 {
			continue
		}
		rel := parts[1]
		if !strings.HasPrefix(rel, skillSubdir) {
			continue
		}
		inner := strings.TrimPrefix(rel, skillSubdir)
		if inner == "" {
			continue
		}
		if strings.Contains(inner, "..") || strings.HasPrefix(inner, "/") {
			return 0, fmt.Errorf("unsafe path in archive: %s", hdr.Name)
		}
		targetPath := filepath.Join(absDest, inner)
		if !strings.HasPrefix(targetPath, absDest+string(os.PathSeparator)) && targetPath != absDest {
			return 0, fmt.Errorf("path escapes destination: %s", hdr.Name)
		}
		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(targetPath, 0o755); err != nil {
				return 0, err
			}
		case tar.TypeReg, tar.TypeRegA:
			if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
				return 0, err
			}
			f, err := os.OpenFile(targetPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
			if err != nil {
				return 0, err
			}
			// Cap the per-file copy by the remaining budget so a
			// pathological tar entry cannot exhaust disk before we
			// notice. CopyN returns ErrUnexpectedEOF when the source
			// has fewer bytes than the limit, which is the normal
			// case; we differentiate that from "hit the cap".
			n, copyErr := io.CopyN(f, tr, budget+1)
			f.Close()
			if copyErr != nil && copyErr != io.EOF {
				return 0, copyErr
			}
			if n > budget {
				return 0, fmt.Errorf("archive exceeds %d-byte extract budget", skillMaxExtractBytes)
			}
			budget -= n
			count++
		}
	}
	if count == 0 {
		return 0, fmt.Errorf("archive contains no files under %s", skillSubdir)
	}
	return count, nil
}
