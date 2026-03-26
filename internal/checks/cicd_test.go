package checks

import (
	"testing"

	"github.com/spbuilds/repohealth/internal/model"
)

func TestCIConfigExistsCheck_GitHubActions(t *testing.T) {
	ctx := &model.ScanContext{
		Dirs: []string{".github/workflows"},
	}
	check := &CIConfigExistsCheck{}
	result := check.Run(ctx)
	if result.Status != model.StatusFull {
		t.Errorf("Status = %v, want Full for .github/workflows", result.Status)
	}
	if result.Details != "GitHub Actions" {
		t.Errorf("Details = %q, want 'GitHub Actions'", result.Details)
	}
}

func TestCIConfigExistsCheck_GitLabCI(t *testing.T) {
	ctx := &model.ScanContext{
		Files: []model.FileInfo{{Path: ".gitlab-ci.yml", Name: ".gitlab-ci.yml", Size: 100}},
	}
	check := &CIConfigExistsCheck{}
	result := check.Run(ctx)
	if result.Status != model.StatusFull {
		t.Errorf("Status = %v, want Full for .gitlab-ci.yml", result.Status)
	}
}

func TestCIConfigExistsCheck_None(t *testing.T) {
	ctx := &model.ScanContext{}
	check := &CIConfigExistsCheck{}
	result := check.Run(ctx)
	if result.Status != model.StatusNone {
		t.Errorf("Status = %v, want None for empty repo", result.Status)
	}
}

func TestCIContainsAny(t *testing.T) {
	lines := []string{
		"name: CI",
		"    - run: npm test",
		"    - run: eslint .",
	}
	if !ciContainsAny(lines, []string{"npm test"}) {
		t.Error("expected to find 'npm test'")
	}
	if !ciContainsAny(lines, []string{"eslint"}) {
		t.Error("expected to find 'eslint'")
	}
	if ciContainsAny(lines, []string{"cargo test"}) {
		t.Error("should not find 'cargo test'")
	}
}

func TestCIContainsAny_CaseInsensitive(t *testing.T) {
	lines := []string{"- run: PYTEST"}
	if !ciContainsAny(lines, []string{"pytest"}) {
		t.Error("ciContainsAny should be case-insensitive")
	}
}
