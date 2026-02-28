package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server         string `yaml:"server"`
	ClientID       string `yaml:"client_id,omitempty"`
	DefaultScooter int    `yaml:"default_scooter,omitempty"`
	Output         string `yaml:"output,omitempty"` // "text" or "json"
}

func DefaultConfig() *Config {
	return &Config{
		Server:   "https://sunshine.rescoot.org",
		ClientID: "sunshine-cli",
		Output:   "text",
	}
}

func Dir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "sunshine")
}

func Path() string {
	return filepath.Join(Dir(), "config.yaml")
}

func Load() (*Config, error) {
	cfg := DefaultConfig()

	data, err := os.ReadFile(Path())
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, fmt.Errorf("reading config: %w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	return cfg, nil
}

func (c *Config) Save() error {
	if err := os.MkdirAll(Dir(), 0o700); err != nil {
		return fmt.Errorf("creating config dir: %w", err)
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}

	return os.WriteFile(Path(), data, 0o600)
}
