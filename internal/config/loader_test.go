package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig_NoFile(t *testing.T) {
	cfg, err := LoadConfig(t.TempDir(), "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg != nil {
		t.Error("expected nil config when no file exists")
	}
}

func TestLoadConfig_ValidFile(t *testing.T) {
	dir := t.TempDir()
	content := []byte("version: 1\nthreshold: 75\ndisable:\n  - STAT-03\nexclude:\n  - vendor/\n")
	os.WriteFile(filepath.Join(dir, ".repohealthrc.yaml"), content, 0644)

	cfg, err := LoadConfig(dir, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected config, got nil")
	}
	if cfg.Threshold != 75 {
		t.Errorf("Threshold = %d, want 75", cfg.Threshold)
	}
	if len(cfg.Disable) != 1 || cfg.Disable[0] != "STAT-03" {
		t.Errorf("Disable = %v, want [STAT-03]", cfg.Disable)
	}
	if len(cfg.Exclude) != 1 || cfg.Exclude[0] != "vendor/" {
		t.Errorf("Exclude = %v, want [vendor/]", cfg.Exclude)
	}
}

func TestLoadConfig_InvalidYAML(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, ".repohealthrc.yaml"), []byte("{{invalid"), 0644)

	_, err := LoadConfig(dir, "")
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func TestLoadConfig_ExplicitPath(t *testing.T) {
	dir := t.TempDir()
	customPath := filepath.Join(dir, "custom.yaml")
	os.WriteFile(customPath, []byte("threshold: 80\n"), 0644)

	cfg, err := LoadConfig(dir, customPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Threshold != 80 {
		t.Errorf("Threshold = %d, want 80", cfg.Threshold)
	}
}

func TestLoadConfig_ExplicitPathNotFound(t *testing.T) {
	_, err := LoadConfig(t.TempDir(), "/nonexistent/config.yaml")
	if err == nil {
		t.Error("expected error for nonexistent explicit config path")
	}
}

func TestLoadConfig_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, ".repohealthrc.yaml"), []byte(""), 0644)

	cfg, err := LoadConfig(dir, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Empty YAML unmarshals to zero-value struct, not nil
	if cfg == nil {
		t.Error("expected non-nil config for empty file")
	}
}
