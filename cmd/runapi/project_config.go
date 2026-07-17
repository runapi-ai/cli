package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

const projectConfigName = ".runapi.toml"

type projectConfig struct {
	CallbackAPIKeyID string `toml:"callback_api_key_id"`
}

func loadProjectConfig(startDir string) (projectConfig, string, error) {
	path, err := projectConfigPath(startDir)
	if err != nil {
		return projectConfig{}, "", err
	}

	cfg, err := loadProjectConfigFile(path)
	return cfg, path, err
}

func loadProjectConfigFile(path string) (projectConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return projectConfig{}, nil
		}
		return projectConfig{}, err
	}

	var cfg projectConfig
	metadata, err := toml.Decode(string(data), &cfg)
	if err != nil {
		return projectConfig{}, fmt.Errorf("invalid %s: %w", path, err)
	}
	if undecoded := metadata.Undecoded(); len(undecoded) > 0 {
		keys := make([]string, 0, len(undecoded))
		for _, key := range undecoded {
			keys = append(keys, key.String())
		}
		return projectConfig{}, fmt.Errorf("invalid %s: unknown key(s): %s", path, strings.Join(keys, ", "))
	}
	if strings.TrimSpace(cfg.CallbackAPIKeyID) == "" {
		return projectConfig{}, fmt.Errorf("invalid %s: callback_api_key_id cannot be blank", path)
	}
	return cfg, nil
}

func projectConfigPath(startDir string) (string, error) {
	root, err := findGitRoot(startDir)
	if err != nil {
		return "", err
	}
	if root != "" {
		return filepath.Join(root, projectConfigName), nil
	}

	dir := strings.TrimSpace(startDir)
	if dir == "" {
		dir, err = os.Getwd()
		if err != nil {
			return "", err
		}
	}
	dir, err = filepath.Abs(dir)
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, projectConfigName), nil
}

func saveProjectConfig(path string, cfg projectConfig) error {
	if strings.TrimSpace(cfg.CallbackAPIKeyID) == "" {
		return fmt.Errorf("callback_api_key_id cannot be blank")
	}
	var data bytes.Buffer
	if err := toml.NewEncoder(&data).Encode(cfg); err != nil {
		return err
	}

	temporary, err := os.CreateTemp(filepath.Dir(path), projectConfigName+".tmp-*")
	if err != nil {
		return err
	}
	temporaryPath := temporary.Name()
	defer os.Remove(temporaryPath)

	if err := temporary.Chmod(0o644); err != nil {
		temporary.Close()
		return err
	}
	if _, err := temporary.Write(data.Bytes()); err != nil {
		temporary.Close()
		return err
	}
	if err := temporary.Sync(); err != nil {
		temporary.Close()
		return err
	}
	if err := temporary.Close(); err != nil {
		return err
	}
	return os.Rename(temporaryPath, path)
}

func findGitRoot(startDir string) (string, error) {
	start := strings.TrimSpace(startDir)
	if start == "" {
		var err error
		start, err = os.Getwd()
		if err != nil {
			return "", err
		}
	}
	dir, err := filepath.Abs(start)
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir, nil
		} else if !os.IsNotExist(err) {
			return "", err
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", nil
		}
		dir = parent
	}
}
