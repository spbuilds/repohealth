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
		{"LoginTests.cs", true},
		{"LoginTest.cs", true},
		{"latest.cs", false},
		{"Login.cs", false},
		{"login_test.dart", true},
		{"login.dart", false},
		{"login_test.exs", true},
		{"login_spec.lua", true},
		{"test-login.R", true},
		{"test_login.r", true},
		{"helpers.R", false},
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

func TestIsSourceFile_Extensions(t *testing.T) {
	for _, name := range []string{"App.tsx", "index.jsx", "Program.cs", "main.dart", "query.sql"} {
		if !isSourceFile(model.FileInfo{Name: name, Path: name}) {
			t.Errorf("isSourceFile(%q) = false, want true", name)
		}
	}
	for _, name := range []string{"README.md", "config.yaml", "package.json", "App.test.tsx", "Makefile"} {
		if isSourceFile(model.FileInfo{Name: name, Path: name}) {
			t.Errorf("isSourceFile(%q) = true, want false", name)
		}
	}
}

func TestIsProgramFile(t *testing.T) {
	for _, name := range []string{"main.go", "App.tsx", "Program.cs", "query.sql"} {
		if !isProgramFile(model.FileInfo{Name: name, Path: name}) {
			t.Errorf("isProgramFile(%q) = false, want true", name)
		}
	}
	for _, name := range []string{"index.html", "style.css", "theme.scss", "main_test.go", "README.md"} {
		if isProgramFile(model.FileInfo{Name: name, Path: name}) {
			t.Errorf("isProgramFile(%q) = true, want false", name)
		}
	}
}

func TestCommentRatioCheck_IgnoresMarkup(t *testing.T) {
	goSource := "// Package main does things.\n// It has comments.\n// Three of them.\npackage main\n\nfunc main() {}\n"
	page := "<html>\n<body>\n<p>one</p>\n<p>two</p>\n<p>three</p>\n<p>four</p>\n<p>five</p>\n<p>six</p>\n<p>seven</p>\n</body>\n</html>\n"
	dir := makeTempRepo(t, map[string]string{
		"docs/a.html": page,
		"docs/b.html": page,
		"docs/c.html": page,
		"main.go":     goSource,
	})
	ctx := &model.ScanContext{
		RepoPath: dir,
		Files:    fileInfos("docs/a.html", "docs/b.html", "docs/c.html", "main.go"),
	}
	result := (&CommentRatioCheck{}).Run(ctx)
	if result.Status != model.StatusFull {
		t.Errorf("Status = %v, want Full — markup pages must not dilute the comment ratio; details: %s", result.Status, result.Details)
	}
}
