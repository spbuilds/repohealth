package checks

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spbuilds/repohealth/internal/model"
	"github.com/spbuilds/repohealth/internal/scanner"
)

var sourceExtensions = map[string]bool{
	".go": true, ".py": true, ".js": true, ".ts": true,
	".java": true, ".rs": true, ".rb": true, ".c": true,
	".cpp": true, ".h": true, ".sh": true, ".php": true,
	".swift": true, ".kt": true,
}

var nonCodeLanguages = map[string]bool{
	"Markdown": true, "YAML": true, "JSON": true, "TOML": true,
}

// isTestFile returns true if the file name indicates a test file.
func isTestFile(name string) bool {
	lower := strings.ToLower(name)
	// Go: *_test.go, Rust: *_test.rs, Ruby: *_test.rb
	if strings.HasSuffix(lower, "_test.go") || strings.HasSuffix(lower, "_test.rs") || strings.HasSuffix(lower, "_test.rb") {
		return true
	}
	// Python: test_*.py (prefix only)
	if strings.HasPrefix(lower, "test_") && strings.HasSuffix(lower, ".py") {
		return true
	}
	// JS/TS: *.test.* or *.spec.*
	if strings.Contains(lower, ".test.") || strings.Contains(lower, ".spec.") {
		return true
	}
	// Java: *Test.java or *Tests.java
	if strings.HasSuffix(lower, "test.java") || strings.HasSuffix(lower, "tests.java") {
		return true
	}
	return false
}

// isSourceFile returns true if the file is a non-test source code file.
func isSourceFile(f model.FileInfo) bool {
	if f.IsDir {
		return false
	}
	dot := strings.LastIndex(f.Name, ".")
	if dot < 0 {
		return false
	}
	ext := strings.ToLower(f.Name[dot:])
	return sourceExtensions[ext] && !isTestFile(f.Name)
}

// STAT-01: Source files exist
type SourceFilesExistCheck struct{}

func (c *SourceFilesExistCheck) ID() string       { return "STAT-01" }
func (c *SourceFilesExistCheck) Category() string { return "stats" }
func (c *SourceFilesExistCheck) Name() string     { return "Source files exist" }
func (c *SourceFilesExistCheck) MaxPoints() int   { return 1 }

func (c *SourceFilesExistCheck) Run(ctx *model.ScanContext) model.CheckResult {
	count := 0
	for _, f := range ctx.Files {
		if isSourceFile(f) {
			count++
		}
	}

	if count > 0 {
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusFull, Points: c.MaxPoints(), MaxPoints: c.MaxPoints(),
			Details: fmt.Sprintf("%d source files", count),
		}
	}
	return model.CheckResult{
		ID: c.ID(), Category: c.Category(), Name: c.Name(),
		Status: model.StatusNone, Points: 0, MaxPoints: c.MaxPoints(),
		Details:    "No source files found",
		Suggestion: "Add source code files to the repository",
	}
}

// STAT-02: Language diversity
type LanguageDiversityCheck struct{}

func (c *LanguageDiversityCheck) ID() string       { return "STAT-02" }
func (c *LanguageDiversityCheck) Category() string { return "stats" }
func (c *LanguageDiversityCheck) Name() string     { return "Language diversity" }
func (c *LanguageDiversityCheck) MaxPoints() int   { return 1 }

func (c *LanguageDiversityCheck) Run(ctx *model.ScanContext) model.CheckResult {
	count := 0
	for lang := range ctx.Languages {
		if !nonCodeLanguages[lang] {
			count++
		}
	}

	if count >= 2 {
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusFull, Points: c.MaxPoints(), MaxPoints: c.MaxPoints(),
			Details: fmt.Sprintf("%d languages", count),
		}
	}
	if count == 1 {
		// Finding the primary language is a pass — polyglot is not required
		var primaryLang string
		for lang := range ctx.Languages {
			if !nonCodeLanguages[lang] {
				primaryLang = lang
				break
			}
		}
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusFull, Points: c.MaxPoints(), MaxPoints: c.MaxPoints(),
			Details: fmt.Sprintf("1 language (%s)", primaryLang),
		}
	}
	return model.CheckResult{
		ID: c.ID(), Category: c.Category(), Name: c.Name(),
		Status: model.StatusNone, Points: 0, MaxPoints: c.MaxPoints(),
		Details:    "No languages detected",
		Suggestion: "Add source code files to the repository",
	}
}

// STAT-03: Comment ratio
type CommentRatioCheck struct{}

func (c *CommentRatioCheck) ID() string       { return "STAT-03" }
func (c *CommentRatioCheck) Category() string { return "stats" }
func (c *CommentRatioCheck) Name() string     { return "Comment ratio" }
func (c *CommentRatioCheck) MaxPoints() int   { return 2 }

func (c *CommentRatioCheck) Run(ctx *model.ScanContext) model.CheckResult {
	var sourceFiles []model.FileInfo
	for _, f := range ctx.Files {
		if isSourceFile(f) {
			sourceFiles = append(sourceFiles, f)
		}
	}
	// Sort source files by path for deterministic sampling across platforms
	sort.Slice(sourceFiles, func(i, j int) bool {
		return sourceFiles[i].Path < sourceFiles[j].Path
	})

	maxSample := 50
	if len(sourceFiles) < maxSample {
		maxSample = len(sourceFiles)
	}
	sampled := sourceFiles[:maxSample]

	if len(sampled) == 0 {
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusSkipped, Points: 0, MaxPoints: c.MaxPoints(),
			Details: "No source files to sample",
		}
	}

	totalLines := 0
	commentLines := 0
	for _, f := range sampled {
		lines, err := scanner.ReadFileLines(ctx.RepoPath, f.Path)
		if err != nil || lines == nil {
			continue
		}
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			totalLines++
			if strings.HasPrefix(trimmed, "//") ||
				strings.HasPrefix(trimmed, "#") ||
				strings.HasPrefix(trimmed, "--") ||
				strings.HasPrefix(trimmed, "/*") ||
				strings.HasPrefix(trimmed, "* ") ||
				trimmed == "*" {
				commentLines++
			}
		}
	}

	if totalLines == 0 {
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusSkipped, Points: 0, MaxPoints: c.MaxPoints(),
			Details: "No readable lines found",
		}
	}

	ratio := float64(commentLines) / float64(totalLines)
	pct := int(ratio * 100)

	if ratio > 0.10 {
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusFull, Points: c.MaxPoints(), MaxPoints: c.MaxPoints(),
			Details: fmt.Sprintf("%d%% comment ratio", pct),
		}
	}
	if ratio >= 0.05 {
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusPartial, Points: 1, MaxPoints: c.MaxPoints(),
			Details:    fmt.Sprintf("%d%% comment ratio", pct),
			Suggestion: "Increase inline documentation (target >10% comment ratio)",
		}
	}
	return model.CheckResult{
		ID: c.ID(), Category: c.Category(), Name: c.Name(),
		Status: model.StatusNone, Points: 0, MaxPoints: c.MaxPoints(),
		Details:    fmt.Sprintf("%d%% comment ratio", pct),
		Suggestion: "Add comments to explain non-obvious code (currently <5%)",
	}
}

// STAT-04: No vendor bloat
type NoVendorBloatCheck struct{}

func (c *NoVendorBloatCheck) ID() string       { return "STAT-04" }
func (c *NoVendorBloatCheck) Category() string { return "stats" }
func (c *NoVendorBloatCheck) Name() string     { return "No vendor bloat" }
func (c *NoVendorBloatCheck) MaxPoints() int   { return 1 }

func (c *NoVendorBloatCheck) Run(ctx *model.ScanContext) model.CheckResult {
	// Check for vendor/node_modules directories via os.Stat since scanner
	// skips these dirs (they're in skipDirs) so ctx.Files won't contain them.
	// Only penalize if the directory exists AND is not already in .gitignore.
	gitignorePatterns := readGitignorePatterns(ctx)
	vendorDirs := []struct{ dir, name, suggestion string }{
		{"vendor", "vendor/ directory found", "Consider removing the vendor directory and using a package manager"},
		{"node_modules", "node_modules/ found in repository", "Add node_modules/ to .gitignore"},
	}
	for _, v := range vendorDirs {
		dirPath := filepath.Join(ctx.RepoPath, v.dir)
		if info, err := os.Stat(dirPath); err == nil && info.IsDir() {
			if isInGitignore(gitignorePatterns, v.dir) {
				continue
			}
			return model.CheckResult{
				ID: c.ID(), Category: c.Category(), Name: c.Name(),
				Status: model.StatusNone, Points: 0, MaxPoints: c.MaxPoints(),
				Details:    v.name,
				Suggestion: v.suggestion,
			}
		}
	}
	return model.CheckResult{
		ID: c.ID(), Category: c.Category(), Name: c.Name(),
		Status: model.StatusFull, Points: c.MaxPoints(), MaxPoints: c.MaxPoints(),
		Details: "No vendor bloat detected",
	}
}

// readGitignorePatterns reads .gitignore and returns all non-comment, non-empty lines.
func readGitignorePatterns(ctx *model.ScanContext) []string {
	lines, err := scanner.ReadFileLines(ctx.RepoPath, ".gitignore")
	if err != nil {
		return nil
	}
	var patterns []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			patterns = append(patterns, line)
		}
	}
	return patterns
}

// isInGitignore checks if a directory name is covered by any .gitignore pattern.
func isInGitignore(patterns []string, dir string) bool {
	for _, p := range patterns {
		p = strings.TrimSuffix(p, "/")
		p = strings.TrimPrefix(p, "/")
		if p == dir {
			return true
		}
	}
	return false
}
