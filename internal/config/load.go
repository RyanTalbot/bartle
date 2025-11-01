package config

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type AI struct {
	Enabled     bool    `yaml:"enabled"`
	Provider    string  `yaml:"provider"`
	Model       string  `yaml:"model"`
	APIKey      string  `yaml:"api_key"`
	Temperature float64 `yaml:"temperature"`
}

type Rules struct {
	ScopeRequired  bool     `yaml:"scope_required"`
	MaxLineLength  int      `yaml:"max_line_length"`
	LowercaseStart bool     `yaml:"lowercase_start"`
	Types          []string `yaml:"types"`
}

type Hook struct {
	AutoApply   bool `yaml:"auto_apply"`
	BlockOnFail bool `yaml:"block_on_fail"`
}

type Config struct {
	Style string `yaml:"style"`
	AI    AI     `yaml:"ai"`
	Rules Rules  `yaml:"rules"`
	Hook  Hook   `yaml:"hook"`
}

var (
	ErrNotInGitRepo    = errors.New("not inside a git repository")
	ErrConfigNotFound  = errors.New("config file not found")
	ErrConfigMalformed = errors.New("config file is malformed")
)

// Default returns a Config struct populated with sensible defaults.
// These must always mirror what `bartle init` generates.
func Default() Config {
	return Config{
		Style: "conventional",
		AI: AI{
			Enabled:     false,
			Provider:    "openai",
			Model:       "gpt-5",
			Temperature: 0.2,
			APIKey:      "env:OPENAI_API_KEY",
		},
		Rules: Rules{
			ScopeRequired:  true,
			MaxLineLength:  72,
			LowercaseStart: false,
			Types:          []string{"feat", "fix", "docs", "refactor", "test", "chore"},
		},
		Hook: Hook{
			AutoApply:   false,
			BlockOnFail: true,
		},
	}
}

// findRepoConfigPath walks upward from the starting directory
// until it finds a `.git` folder, then returns the path to `.bartle.yaml`
// within that repo root.
func findRepoConfigPath(startDir string) (string, error) {
	currentDir := startDir
	for {
		if _, err := os.Stat(filepath.Join(currentDir, ".git")); err == nil {
			return filepath.Join(currentDir, ".bartle.yaml"), nil
		}
		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			return "", ErrNotInGitRepo
		}
		currentDir = parentDir
	}
}

// Load merges defaults with YAML config from the repo root.
// Returns the loaded config, the path to the config file, and any error.
func Load() (Config, string, error) {
	defaultConfig := Default()

	workingDir, err := os.Getwd()
	if err != nil {
		return defaultConfig, "", fmt.Errorf("get working directory: %w", err)
	}

	configPath, err := findRepoConfigPath(workingDir)
	if err != nil {
		return defaultConfig, "", err
	}

	rawBytes, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return defaultConfig, configPath, ErrConfigNotFound
		}
		return defaultConfig, configPath, fmt.Errorf("read config: %w", err)
	}

	// Normalize content (remove BOM and CRLF)
	rawBytes = bytes.TrimPrefix(rawBytes, []byte{0xEF, 0xBB, 0xBF}) // UTF-8 BOM
	rawBytes = bytes.ReplaceAll(rawBytes, []byte("\r\n"), []byte("\n"))

	decoder := yaml.NewDecoder(bytes.NewReader(rawBytes))
	decoder.KnownFields(true) // catch typos and unknown fields

	if err := decoder.Decode(&defaultConfig); err != nil {
		return defaultConfig, configPath, fmt.Errorf("%w: %v", ErrConfigMalformed, err)
	}

	return defaultConfig, configPath, nil
}
