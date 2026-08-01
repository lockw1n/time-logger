package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

const (
	dirEnv     = "TL_CONFIG_DIR"
	configFile = "config.json"
)

type Config struct {
	BaseURL          string `json:"base_url,omitempty"`
	DefaultCompanyID uint64 `json:"default_company_id,omitempty"`
	DefaultActivity  string `json:"default_activity,omitempty"`
	Email            string `json:"email,omitempty"`
	// CACertPath trusts a self-signed backend cert (PEM) for HTTPS base URLs.
	CACertPath string `json:"ca_cert,omitempty"`
	// Insecure skips TLS verification entirely — convenient for localhost dev.
	Insecure bool `json:"insecure,omitempty"`
}

func Dir() (string, error) {
	if dir := os.Getenv(dirEnv); dir != "" {
		return dir, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolving home directory: %w", err)
	}

	return filepath.Join(home, ".config", "tl"), nil
}

func Load() (Config, error) {
	dir, err := Dir()
	if err != nil {
		return Config{}, err
	}

	data, err := os.ReadFile(filepath.Join(dir, configFile))
	if errors.Is(err, fs.ErrNotExist) {
		return Config{}, nil
	}
	if err != nil {
		return Config{}, fmt.Errorf("reading %s: %w", configFile, err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parsing %s: %w", configFile, err)
	}

	return cfg, nil
}

func Save(cfg Config) error {
	dir, err := Dir()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding %s: %w", configFile, err)
	}

	return WriteFileAtomic(filepath.Join(dir, configFile), data, 0o600)
}
