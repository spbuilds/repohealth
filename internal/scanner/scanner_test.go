package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanHealthyRepo(t *testing.T) {
	t.Helper()
	repoPath := filepath.Join("..", "..", "testdata", "healthy-repo")

	ctx, err := Scan(repoPath, nil)
	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}

	// Should have detected files
	if len(ctx.Files) == 0 {
		t.Error("expected files to be found in healthy-repo, got 0")
	}

	// Should have detected Python language (main.py, test_example.py)
	if ctx.Languages["Python"] == 0 {
		t.Errorf("expected Python language detected, got languages = %v", ctx.Languages)
	}

	// Should have detected YAML language (ci.yml)
	if ctx.Languages["YAML"] == 0 {
		t.Errorf("expected YAML language detected, got languages = %v", ctx.Languages)
	}

	// Should have detected Markdown language (README.md, CHANGELOG.md, etc.)
	if ctx.Languages["Markdown"] == 0 {
		t.Errorf("expected Markdown language detected, got languages = %v", ctx.Languages)
	}

	// Should have found the tests and .github/workflows dirs
	foundTests := false
	foundWorkflows := false
	for _, d := range ctx.Dirs {
		if d == "tests" {
			foundTests = true
		}
		if d == ".github/workflows" {
			foundWorkflows = true
		}
	}
	if !foundTests {
		t.Errorf("expected 'tests' dir in ctx.Dirs, got %v", ctx.Dirs)
	}
	if !foundWorkflows {
		t.Errorf("expected '.github/workflows' dir in ctx.Dirs, got %v", ctx.Dirs)
	}
}

func TestScanEmptyRepo(t *testing.T) {
	repoPath := filepath.Join("..", "..", "testdata", "empty-repo")

	ctx, err := Scan(repoPath, nil)
	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}

	// Only .gitkeep at root — no recognized language extensions
	if len(ctx.Languages) != 0 {
		t.Errorf("expected no languages detected, got %v", ctx.Languages)
	}

	// .gitkeep is a dotfile at the root level so it is included by the scanner
	if len(ctx.Files) != 1 {
		t.Errorf("expected 1 file (.gitkeep), got %d: %v", len(ctx.Files), ctx.Files)
	}
}

func TestScanNonexistentPath(t *testing.T) {
	_, err := Scan("/nonexistent/path/that/does/not/exist", nil)
	if err == nil {
		t.Error("expected error for nonexistent path, got nil")
	}
}

func TestScanFileNotDirectory(t *testing.T) {
	// Create a temp file and try to scan it as a directory
	f, err := os.CreateTemp("", "repohealth-test-*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	f.Close()
	defer os.Remove(f.Name())

	_, err = Scan(f.Name(), nil)
	if err == nil {
		t.Error("expected error when scanning a file (not a directory), got nil")
	}
}
