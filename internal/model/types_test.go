package model

import "testing"

func TestMatchSuffix(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		want    bool
	}{
		{"foo_test.go", "*_test.go", true},
		{"test_main.py", "test_*.py", true},
		{"app.spec.js", "*.spec.js", true},
		{"main.go", "*_test.go", false},
		{"", "*_test.go", false},
		{"test.go", "", false},
		{"exact.txt", "exact.txt", true},
	}

	for _, tt := range tests {
		got := matchSuffix(tt.name, tt.pattern)
		if got != tt.want {
			t.Errorf("matchSuffix(%q, %q) = %v, want %v", tt.name, tt.pattern, got, tt.want)
		}
	}
}

func TestHasFile(t *testing.T) {
	ctx := &ScanContext{
		Files: []FileInfo{
			{Path: "README.md", Name: "README.md", Size: 200},
			{Path: "LICENSE", Name: "LICENSE", Size: 1000},
			{Path: "src/main.go", Name: "main.go", Size: 500},
		},
	}

	// Found by name
	path, ok := ctx.HasFile("README.md")
	if !ok {
		t.Error("HasFile(README.md) = false, want true")
	}
	if path != "README.md" {
		t.Errorf("HasFile(README.md) path = %q, want %q", path, "README.md")
	}

	// Found by path
	path, ok = ctx.HasFile("src/main.go")
	if !ok {
		t.Error("HasFile(src/main.go) = false, want true")
	}
	if path != "src/main.go" {
		t.Errorf("HasFile(src/main.go) path = %q, want %q", path, "src/main.go")
	}

	// Not found
	_, ok = ctx.HasFile("CONTRIBUTING.md")
	if ok {
		t.Error("HasFile(CONTRIBUTING.md) = true, want false")
	}

	// Multiple names — first match wins
	_, ok = ctx.HasFile("NONEXISTENT.md", "LICENSE")
	if !ok {
		t.Error("HasFile(NONEXISTENT.md, LICENSE) = false, want true (LICENSE exists)")
	}
}

func TestHasDir(t *testing.T) {
	ctx := &ScanContext{
		Dirs: []string{"src", "tests", ".github/workflows"},
	}

	// Found
	path, ok := ctx.HasDir("tests")
	if !ok {
		t.Error("HasDir(tests) = false, want true")
	}
	if path != "tests" {
		t.Errorf("HasDir(tests) path = %q, want %q", path, "tests")
	}

	// Found nested dir
	_, ok = ctx.HasDir(".github/workflows")
	if !ok {
		t.Error("HasDir(.github/workflows) = false, want true")
	}

	// Not found
	_, ok = ctx.HasDir("nonexistent")
	if ok {
		t.Error("HasDir(nonexistent) = true, want false")
	}
}

func TestCountFilesMatching(t *testing.T) {
	ctx := &ScanContext{
		Files: []FileInfo{
			{Name: "main_test.go", IsDir: false},
			{Name: "util_test.go", IsDir: false},
			{Name: "main.go", IsDir: false},
			{Name: "test_helper.py", IsDir: false},
			{Name: "test_main.py", IsDir: false},
			{Name: "app.spec.js", IsDir: false},
			{Name: "src", IsDir: true},
		},
	}

	// Go test files
	count := ctx.CountFilesMatching("*_test.go")
	if count != 2 {
		t.Errorf("CountFilesMatching(*_test.go) = %d, want 2", count)
	}

	// Python test files
	count = ctx.CountFilesMatching("test_*.py")
	if count != 2 {
		t.Errorf("CountFilesMatching(test_*.py) = %d, want 2", count)
	}

	// JS spec files
	count = ctx.CountFilesMatching("*.spec.js")
	if count != 1 {
		t.Errorf("CountFilesMatching(*.spec.js) = %d, want 1", count)
	}

	// Dirs should not be counted
	count = ctx.CountFilesMatching("*")
	if count != 6 {
		t.Errorf("CountFilesMatching(*) = %d, want 6 (dirs excluded)", count)
	}

	// Multiple patterns, no double-counting
	count = ctx.CountFilesMatching("*_test.go", "test_*.py")
	if count != 4 {
		t.Errorf("CountFilesMatching(*_test.go, test_*.py) = %d, want 4", count)
	}
}
