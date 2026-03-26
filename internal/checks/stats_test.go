package checks

import (
	"testing"

	"github.com/spbuilds/repohealth/internal/model"
)

func TestSourceFilesExistCheck_HasFiles(t *testing.T) {
	ctx := &model.ScanContext{
		Files: []model.FileInfo{
			{Path: "main.go", Name: "main.go", Size: 500},
			{Path: "app.py", Name: "app.py", Size: 300},
		},
	}
	check := &SourceFilesExistCheck{}
	result := check.Run(ctx)
	if result.Status != model.StatusFull {
		t.Errorf("Status = %v, want Full", result.Status)
	}
}

func TestSourceFilesExistCheck_NoSourceFiles(t *testing.T) {
	ctx := &model.ScanContext{
		Files: []model.FileInfo{
			{Path: "README.md", Name: "README.md", Size: 100},
		},
	}
	check := &SourceFilesExistCheck{}
	result := check.Run(ctx)
	if result.Status != model.StatusNone {
		t.Errorf("Status = %v, want None for no source files", result.Status)
	}
}

func TestLanguageDiversityCheck_Multiple(t *testing.T) {
	ctx := &model.ScanContext{
		Languages: map[string]int{"Go": 10, "Python": 5},
	}
	check := &LanguageDiversityCheck{}
	result := check.Run(ctx)
	if result.Status != model.StatusFull {
		t.Errorf("Status = %v, want Full for 2 languages", result.Status)
	}
}

func TestLanguageDiversityCheck_OnlyMarkdown(t *testing.T) {
	ctx := &model.ScanContext{
		Languages: map[string]int{"Markdown": 10},
	}
	check := &LanguageDiversityCheck{}
	result := check.Run(ctx)
	if result.Status != model.StatusNone {
		t.Errorf("Status = %v, want None for only non-code languages", result.Status)
	}
}

func TestNoVendorBloatCheck_NoVendor(t *testing.T) {
	ctx := &model.ScanContext{
		RepoPath: t.TempDir(),
	}
	check := &NoVendorBloatCheck{}
	result := check.Run(ctx)
	if result.Status != model.StatusFull {
		t.Errorf("Status = %v, want Full when no vendor dir", result.Status)
	}
}

func TestIsTestFile(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{"main_test.go", true},
		{"test_main.py", true},
		{"app.test.ts", true},
		{"app.spec.js", true},
		{"main.go", false},
		{"app.py", false},
		{"latest_data.py", false},
		{"contest.go", false},
	}
	for _, tt := range tests {
		got := isTestFile(tt.name)
		if got != tt.want {
			t.Errorf("isTestFile(%q) = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestIsSourceFile(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{"main.go", true},
		{"app.py", true},
		{"index.ts", true},
		{"README.md", false},
		{"config.yaml", false},
		{"data.json", false},
	}
	for _, tt := range tests {
		f := model.FileInfo{Name: tt.name}
		got := isSourceFile(f)
		if got != tt.want {
			t.Errorf("isSourceFile(%q) = %v, want %v", tt.name, got, tt.want)
		}
	}
}
