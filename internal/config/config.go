package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// Config stores CLI configuration values.
type Config struct {
	Token    string `json:"token,omitempty"`
	APIBase  string `json:"api_base,omitempty"`
	Project  string `json:"default_project,omitempty"`
	Labels   string `json:"default_labels,omitempty"`
	LabelCLI bool   `json:"label_cli,omitempty"`
	Filename string `json:"-"`
}

// DefaultPath returns the default config file path.
func DefaultPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("config dir: %w", err)
	}
	return filepath.Join(dir, "todi", "config.json"), nil
}

// Load reads config values from the provided path.
func Load(path string) (*Config, error) {
	if path == "" {
		return &Config{}, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &Config{Filename: path}, nil
		}
		return nil, fmt.Errorf("read config: %w", err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	cfg.Filename = path
	return &cfg, nil
}

// Save writes config values to the provided path.
func (c *Config) Save(path string) error {
	if path == "" {
		return fmt.Errorf("config path empty")
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("mkdir config dir: %w", err)
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("encode config: %w", err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("write config: %w", err)
	}
	return nil
}
