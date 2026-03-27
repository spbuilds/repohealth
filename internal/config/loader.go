// Package config handles loading and parsing of RepoHealth configuration files.
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// FileConfig represents the .repohealthrc.yaml configuration.
type FileConfig struct {
	Version   int `yaml:"version"`
	Threshold int `yaml:"threshold"`
	// Weights allows overriding category scoring weights (not yet implemented).
	Weights map[string]int `yaml:"weights,omitempty"`
	Disable []string       `yaml:"disable"`
	Exclude []string       `yaml:"exclude"`
}

// LoadConfig loads configuration. If configFile is non-empty, it reads that
// exact file. Otherwise it looks for .repohealthrc.yaml in repoPath.
func LoadConfig(repoPath, configFile string) (*FileConfig, error) {
	configPath := filepath.Join(repoPath, ".repohealthrc.yaml")
	if configFile != "" {
		configPath = configFile
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) && configFile == "" {
			return nil, nil // auto-discovery: no config is fine
		}
		return nil, err // explicit --config path must exist
	}

	var cfg FileConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("invalid config file %s: %w", configPath, err)
	}

	return &cfg, nil
}
