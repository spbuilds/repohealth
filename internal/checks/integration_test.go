package checks

import (
	"path/filepath"
	"testing"

	"github.com/spbuilds/repohealth/internal/model"
	"github.com/spbuilds/repohealth/internal/scanner"
)

func TestIntegration_HealthyRepo(t *testing.T) {
	repoPath := filepath.Join("..", "..", "testdata", "healthy-repo")

	ctx, err := scanner.Scan(repoPath, nil)
	if err != nil {
		t.Fatalf("scanner.Scan() error = %v", err)
	}

	registry := NewRegistry()
	results := Run(registry.All(), ctx)

	// Collect docs results and sum points
	docsScore := 0
	docsMax := 0
	for _, r := range results {
		if r.Category == "docs" {
			docsScore += r.Points
			docsMax += r.MaxPoints
		}
	}

	if docsScore != 15 {
		t.Errorf("docs score = %d, want 15 (all docs checks should pass)", docsScore)
	}
	if docsMax != 15 {
		t.Errorf("docs max = %d, want 15", docsMax)
	}

	// Verify test files are detected (TST-01)
	var tst01 *model.CheckResult
	for i := range results {
		if results[i].ID == "TST-01" {
			tst01 = &results[i]
		}
	}
	if tst01 == nil {
		t.Fatal("TST-01 check result not found")
	}
	if tst01.Status != model.StatusFull {
		t.Errorf("TST-01 status = %v, want Full (test_example.py should be detected)", tst01.Status)
	}

	// Verify CI is detected (CI-01)
	var ci01 *model.CheckResult
	for i := range results {
		if results[i].ID == "CI-01" {
			ci01 = &results[i]
		}
	}
	if ci01 == nil {
		t.Fatal("CI-01 check result not found")
	}
	if ci01.Status != model.StatusFull {
		t.Errorf("CI-01 status = %v, want Full (.github/workflows should be detected)", ci01.Status)
	}
}

func TestIntegration_EmptyRepo(t *testing.T) {
	repoPath := filepath.Join("..", "..", "testdata", "empty-repo")

	ctx, err := scanner.Scan(repoPath, nil)
	if err != nil {
		t.Fatalf("scanner.Scan() error = %v", err)
	}

	registry := NewRegistry()
	results := Run(registry.All(), ctx)

	// docs score should be 0 — no community files
	docsScore := 0
	for _, r := range results {
		if r.Category == "docs" {
			docsScore += r.Points
		}
	}
	if docsScore != 0 {
		t.Errorf("docs score = %d, want 0 for empty-repo", docsScore)
	}

	// tests score should be 0 — no test files or dirs
	testsScore := 0
	for _, r := range results {
		if r.Category == "tests" {
			testsScore += r.Points
		}
	}
	if testsScore != 0 {
		t.Errorf("tests score = %d, want 0 for empty-repo", testsScore)
	}

	// cicd score should be 0 — no CI config
	cicdScore := 0
	for _, r := range results {
		if r.Category == "cicd" {
			cicdScore += r.Points
		}
	}
	if cicdScore != 0 {
		t.Errorf("cicd score = %d, want 0 for empty-repo", cicdScore)
	}
}
