package model

// Status represents the assessment level of a check.
type Status string

const (
	StatusFull    Status = "full"
	StatusPartial Status = "partial"
	StatusNone    Status = "none"
	StatusSkipped Status = "skipped"
)

// CheckResult is the output of a single check.
type CheckResult struct {
	ID         string `json:"id"`
	Category   string `json:"category"`
	Name       string `json:"name"`
	Status     Status `json:"status"`
	Points     int    `json:"points"`
	MaxPoints  int    `json:"max_points"`
	Details    string `json:"details"`
	Suggestion string `json:"suggestion"`
}

// CategoryResult aggregates check results for a category.
type CategoryResult struct {
	Name     string `json:"name"`
	Label    string `json:"label"`
	Score    int    `json:"score"`
	MaxScore int    `json:"max_score"`
}

// Suggestion is an actionable recommendation.
type Suggestion struct {
	CheckID string `json:"check_id"`
	Impact  int    `json:"impact"`
	Message string `json:"message"`
}

// Report is the complete output structure.
type Report struct {
	Version       string           `json:"version"`
	RepoPath      string           `json:"repo_path"`
	Timestamp     string           `json:"timestamp"`
	DurationMs    int64            `json:"duration_ms"`
	Score         int              `json:"score"`
	MaxScore      int              `json:"max_score"`
	Grade         string           `json:"grade"`
	Categories    []CategoryResult `json:"categories"`
	Checks        []CheckResult    `json:"checks"`
	Suggestions   []Suggestion     `json:"suggestions"`
	Languages     map[string]int   `json:"languages"`
	FilesAnalyzed int              `json:"files_analyzed"`
}

// FileInfo holds metadata about a file in the repo.
type FileInfo struct {
	Path  string
	Name  string
	Size  int64
	IsDir bool
}

// ScanContext provides shared repo metadata to all checks.
type ScanContext struct {
	RepoPath     string
	GitAvailable bool
	Files        []FileInfo
	Dirs         []string
	Languages    map[string]int
}

// HasFile checks if a file with any of the given names exists in the scan context.
func (ctx *ScanContext) HasFile(names ...string) (string, bool) {
	for _, f := range ctx.Files {
		for _, name := range names {
			if f.Name == name || f.Path == name {
				return f.Path, true
			}
		}
	}
	return "", false
}

// HasRootFile checks if a file exists at the repo root or in .github/.
// Use this for community files (README, LICENSE, etc.) that should only
// count when they are at the top level, not inside build output or dependencies.
func (ctx *ScanContext) HasRootFile(names ...string) (string, bool) {
	for _, f := range ctx.Files {
		for _, name := range names {
			// Match exact root path (e.g., "README.md") or .github/ path
			if f.Path == name {
				return f.Path, true
			}
		}
	}
	return "", false
}

// RootFileSize returns the size of a file at the repo root, or -1 if not found.
func (ctx *ScanContext) RootFileSize(names ...string) int64 {
	for _, f := range ctx.Files {
		for _, name := range names {
			if f.Path == name {
				return f.Size
			}
		}
	}
	return -1
}

// HasDir checks if a directory with any of the given names exists.
func (ctx *ScanContext) HasDir(names ...string) (string, bool) {
	for _, d := range ctx.Dirs {
		for _, name := range names {
			if d == name {
				return d, true
			}
		}
	}
	return "", false
}

// FileSize returns the size of a file by name, or -1 if not found.
func (ctx *ScanContext) FileSize(names ...string) int64 {
	for _, f := range ctx.Files {
		for _, name := range names {
			if f.Name == name || f.Path == name {
				return f.Size
			}
		}
	}
	return -1
}

// CountFilesMatching counts files whose names match any of the given suffixes.
func (ctx *ScanContext) CountFilesMatching(suffixes ...string) int {
	count := 0
	for _, f := range ctx.Files {
		if f.IsDir {
			continue
		}
		for _, suffix := range suffixes {
			if matchSuffix(f.Name, suffix) {
				count++
				break
			}
		}
	}
	return count
}

// matchSuffix checks if a filename matches a glob-like pattern.
// Supports: "*_test.go" (suffix), "test_*.py" (prefix+suffix), "*.spec.js" (suffix).
func matchSuffix(name, pattern string) bool {
	if len(pattern) == 0 {
		return false
	}

	starIdx := -1
	for i, c := range pattern {
		if c == '*' {
			starIdx = i
			break
		}
	}

	// No wildcard — exact match
	if starIdx < 0 {
		return name == pattern
	}

	prefix := pattern[:starIdx]
	suffix := pattern[starIdx+1:]

	if len(name) < len(prefix)+len(suffix) {
		return false
	}

	if prefix != "" && name[:len(prefix)] != prefix {
		return false
	}

	if suffix != "" && name[len(name)-len(suffix):] != suffix {
		return false
	}

	return true
}
