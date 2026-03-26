package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// FileConfig represents the .repohealthrc.yaml configuration.
type FileConfig struct {
	Version   int            `yaml:"version"`
	Threshold int            `yaml:"threshold"`
	Weights   map[string]int `yaml:"weights"`
	Disable   []string       `yaml:"disable"`
	Exclude   []string       `yaml:"exclude"`
}

// LoadConfig loads .repohealthrc.yaml from the given directory.
// Returns nil if no config file found (not an error).
func LoadConfig(repoPath string) (*FileConfig, error) {
	configPath := filepath.Join(repoPath, ".repohealthrc.yaml")

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // no config file is fine
		}
		return nil, err
	}

	var cfg FileConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("invalid .repohealthrc.yaml: %w", err)
	}

	return &cfg, nil
}
