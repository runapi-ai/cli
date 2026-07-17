package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestProjectConfigLoadsFromGitRoot(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	nested := filepath.Join(root, "services", "worker")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, ".runapi.toml")
	if err := os.WriteFile(path, []byte("callback_api_key_id = \"token_project\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, gotPath, err := loadProjectConfig(nested)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.CallbackAPIKeyID != "token_project" {
		t.Fatalf("expected callback key id, got %q", cfg.CallbackAPIKeyID)
	}
	if gotPath != path {
		t.Fatalf("expected path %q, got %q", path, gotPath)
	}
}

func TestProjectConfigUsesCurrentDirectoryOutsideGitRepository(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, projectConfigName)
	if err := os.WriteFile(path, []byte("callback_api_key_id = \"token_local\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, gotPath, err := loadProjectConfig(dir)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.CallbackAPIKeyID != "token_local" || gotPath != path {
		t.Fatalf("expected local config %q, got %#v at %q", path, cfg, gotPath)
	}
}

func TestProjectConfigRejectsUnknownKeys(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(
		filepath.Join(root, ".runapi.toml"),
		[]byte("callback_api_key_id = \"token_project\"\napi_key = \"secret\"\n"),
		0o644,
	); err != nil {
		t.Fatal(err)
	}

	_, _, err := loadProjectConfig(root)
	if err == nil || !strings.Contains(err.Error(), "api_key") {
		t.Fatalf("expected unknown-key error naming api_key, got %v", err)
	}
}

func TestProjectConfigWritesOnlyCallbackAPIKeyID(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, ".runapi.toml")

	if err := saveProjectConfig(path, projectConfig{CallbackAPIKeyID: "token_project"}); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "callback_api_key_id = \"token_project\"\n" {
		t.Fatalf("unexpected project config: %q", data)
	}
	if matches, err := filepath.Glob(filepath.Join(root, ".runapi.toml.tmp-*")); err != nil || len(matches) != 0 {
		t.Fatalf("temporary config files remain: %v (err=%v)", matches, err)
	}
}
