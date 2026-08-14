package main

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/runapi-ai/core-sdk/go/core"
	"github.com/spf13/cobra"
)

const (
	defaultSkillSourceRepo = "runapi-ai/cli-skill"
	skillName              = "runapi-cli"
	skillSubdir            = "skills/runapi-cli/"
	skillDownloadTimeout   = 60 * time.Second
)

// Hard ceiling on the total decompressed bytes we will write while
// extracting the skill tarball. The source repo is small (well under
// 1 MB at the time of writing); 50 MB is a comfortable upper bound
// that still trips long before a real zip bomb causes harm. Declared
// as var so tests can lower the cap.
var skillMaxExtractBytes int64 = 50 * 1024 * 1024

func setSkillMaxExtractBytesForTest(n int64) { skillMaxExtractBytes = n }

type skillArchive struct {
	tag  string
	body io.ReadCloser
}

// skillTargets maps a built-in target name to a path under the user's HOME
// where the skill directory should be created. The skill directory itself
// is always named "runapi-cli" inside that path (e.g. ~/.claude/skills/runapi-cli).
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
		Long: "install-skill downloads the canonical skills/runapi-cli/ directory from the runapi-ai/cli-skill\n" +
			"release archive and writes it to the runtime's skill directory.\n\n" +
			"Targets: claude, codex, gemini, openclaw, hermes. Use --target-dir for a custom path.",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.runInstallSkill(cmd.Context(), o)
		},
	}
	cmd.Flags().StringVar(&o.target, "target", "", "Built-in target: "+strings.Join(supportedTargetNames(), ", "))
	cmd.Flags().StringVar(&o.targetDir, "target-dir", "", "Custom install directory (overrides --target). The skill is placed in <dir>/cli.")
	cmd.Flags().StringVar(&o.version, "version", "", "Skill tag to install. Default: latest stable skill tag.")
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
	}

	tag, err := c.resolveSkillTag(ctx, source, o.version)
	if err != nil {
		return err
	}

	archive, err := c.resolveSkillArchive(ctx, source, tag)
	if err != nil {
		return err
	}
	defer archive.body.Close()

	if o.force {
		if err := os.RemoveAll(destDir); err != nil {
			return err
		}
	}
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return err
	}

	files, err := extractSkillFromTarball(archive.body, destDir)
	if err != nil {
		return err
	}

	c.logf("✓ installed skill %q (%d files) to %s", skillName, files, destDir)
	return c.writeJSON(map[string]any{
		"installed":  true,
		"skill":      skillName,
		"target":     o.target,
		"target_dir": destDir,
		"version":    archive.tag,
		"source":     source,
		"files":      files,
	})
}

func (c *cli) resolveSkillTag(ctx context.Context, source, version string) (string, error) {
	v := strings.TrimSpace(version)
	if v != "" {
		return normalizeSkillTag(v)
	}
	return c.latestSkillTag(ctx, source)
}

func (c *cli) resolveSkillArchive(ctx context.Context, source, tag string) (*skillArchive, error) {
	body, err := c.downloadSkillArchive(ctx, source, tag)
	if err != nil {
		return nil, err
	}
	return &skillArchive{tag: tag, body: body}, nil
}

func (c *cli) downloadSkillArchive(ctx context.Context, source, tag string) (io.ReadCloser, error) {
	archiveURL := c.skillArchiveURL(source, tag)
	if err := assertSecureArchiveURL(archiveURL); err != nil {
		return nil, err
	}

	c.logf("downloading %s", archiveURL)
	return httpDownload(ctx, c.httpClient, archiveURL, "application/octet-stream")
}

func (c *cli) skillArchiveURL(source, tag string) string {
	if c.archiveBaseURL != "" {
		return fmt.Sprintf("%s/%s/archive/refs/tags/%s.tar.gz", strings.TrimRight(c.archiveBaseURL, "/"), source, tag)
	}
	return fmt.Sprintf("https://github.com/%s/archive/refs/tags/%s.tar.gz", source, tag)
}

func (c *cli) latestSkillTag(ctx context.Context, source string) (string, error) {
	tags, err := c.availableSkillTags(ctx, source)
	if err != nil {
		return "", err
	}
	return tags[0], nil
}

func (c *cli) availableSkillTags(ctx context.Context, source string) ([]string, error) {
	tagsURL := c.skillTagsURL(source)
	if err := assertSecureArchiveURL(tagsURL); err != nil {
		return nil, err
	}
	body, err := httpDownload(ctx, c.httpClient, tagsURL, "application/vnd.github+json")
	if err != nil {
		return nil, err
	}
	defer body.Close()

	tags, err := decodeSkillTags(body)
	if err != nil {
		return nil, err
	}
	if len(tags) == 0 {
		return nil, core.NewError(core.ErrNotFound, fmt.Sprintf("no release tags found for %s", source), 404, "", nil, nil)
	}
	sort.Slice(tags, func(i, j int) bool {
		return compareSkillTags(tags[i], tags[j]) > 0
	})
	return tags, nil
}

func (c *cli) skillTagsURL(source string) string {
	if c.archiveBaseURL != "" {
		return fmt.Sprintf("%s/repos/%s/tags", strings.TrimRight(c.archiveBaseURL, "/"), source)
	}
	return fmt.Sprintf("https://api.github.com/repos/%s/tags?per_page=100", source)
}

func decodeSkillTags(r io.Reader) ([]string, error) {
	var payload []struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r).Decode(&payload); err != nil {
		return nil, err
	}

	seen := map[string]bool{}
	tags := make([]string, 0, len(payload))
	for _, item := range payload {
		tag := strings.TrimSpace(item.Name)
		version, ok := parseSkillTagVersion(tag)
		if !ok || version.prerelease != "" {
			continue
		}
		if seen[tag] {
			continue
		}
		seen[tag] = true
		tags = append(tags, tag)
	}
	return tags, nil
}

type skillTagVersion struct {
	major      int
	minor      int
	patch      int
	prerelease string
}

func parseSkillTagVersion(tag string) (skillTagVersion, bool) {
	if !strings.HasPrefix(tag, "v") {
		return skillTagVersion{}, false
	}
	coreVersion, prerelease, hasPrerelease := strings.Cut(strings.TrimPrefix(tag, "v"), "-")
	parts := strings.Split(coreVersion, ".")
	if len(parts) != 3 {
		return skillTagVersion{}, false
	}

	major, ok := parseSemverNumber(parts[0])
	if !ok {
		return skillTagVersion{}, false
	}
	minor, ok := parseSemverNumber(parts[1])
	if !ok {
		return skillTagVersion{}, false
	}
	patch, ok := parseSemverNumber(parts[2])
	if !ok {
		return skillTagVersion{}, false
	}
	if hasPrerelease && !isSafePrerelease(prerelease) {
		return skillTagVersion{}, false
	}
	return skillTagVersion{major: major, minor: minor, patch: patch, prerelease: prerelease}, true
}

func parseSemverNumber(s string) (int, bool) {
	if s == "" {
		return 0, false
	}
	if len(s) > 1 && s[0] == '0' {
		return 0, false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return 0, false
		}
	}
	n, err := strconv.Atoi(s)
	return n, err == nil
}

func isSafePrerelease(s string) bool {
	if s == "" {
		return false
	}
	for _, part := range strings.Split(s, ".") {
		if part == "" {
			return false
		}
		for _, r := range part {
			ok := (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-'
			if !ok {
				return false
			}
		}
	}
	return true
}

func compareSkillTags(a, b string) int {
	av, aOK := parseSkillTagVersion(a)
	bv, bOK := parseSkillTagVersion(b)
	switch {
	case aOK && !bOK:
		return 1
	case !aOK && bOK:
		return -1
	case !aOK && !bOK:
		return strings.Compare(a, b)
	}

	for _, pair := range [][2]int{{av.major, bv.major}, {av.minor, bv.minor}, {av.patch, bv.patch}} {
		if pair[0] > pair[1] {
			return 1
		}
		if pair[0] < pair[1] {
			return -1
		}
	}
	return comparePrerelease(av.prerelease, bv.prerelease)
}

func comparePrerelease(a, b string) int {
	if a == b {
		return 0
	}
	if a == "" {
		return 1
	}
	if b == "" {
		return -1
	}

	aParts := strings.Split(a, ".")
	bParts := strings.Split(b, ".")
	for i := 0; i < len(aParts) && i < len(bParts); i++ {
		aNum, aIsNum := parseSemverNumber(aParts[i])
		bNum, bIsNum := parseSemverNumber(bParts[i])
		switch {
		case aIsNum && bIsNum:
			if aNum > bNum {
				return 1
			}
			if aNum < bNum {
				return -1
			}
		case aIsNum:
			return -1
		case bIsNum:
			return 1
		default:
			if cmp := strings.Compare(aParts[i], bParts[i]); cmp != 0 {
				return cmp
			}
		}
	}
	if len(aParts) > len(bParts) {
		return 1
	}
	if len(aParts) < len(bParts) {
		return -1
	}
	return 0
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

func normalizeSkillTag(version string) (string, error) {
	v := strings.TrimSpace(version)
	if !strings.HasPrefix(v, "v") {
		v = "v" + v
	}
	if !isSafeTag(v) {
		return "", core.NewError(core.ErrValidation, fmt.Sprintf("invalid version %q", version), 422, "", nil, nil)
	}
	return v, nil
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
	_, err := url.Parse("https://runapi.ai/" + t)
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

func httpDownload(ctx context.Context, client *http.Client, rawURL, accept string) (io.ReadCloser, error) {
	if client == nil {
		client = &http.Client{Timeout: skillDownloadTimeout}
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(accept) != "" {
		req.Header.Set("Accept", accept)
	}
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

// extractSkillFromTarball pulls the `skills/runapi-cli/` subtree of a GitHub archive
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
