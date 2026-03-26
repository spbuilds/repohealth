package config

// Config holds RepoHealth configuration.
type Config struct {
	Weights  map[string]int
	Disabled []string
	Excludes []string
}

// DefaultConfig returns the default configuration with PRD-defined weights.
func DefaultConfig() *Config {
	return &Config{
		Weights: map[string]int{
			"docs":     15,
			"tests":    20,
			"cicd":     15,
			"deps":     13,
			"security": 10,
			"stats":    5,
			"activity": 15,
			"todo":     7,
		},
		Disabled: nil,
		Excludes: nil,
	}
}
