package checks

import (
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/spbuilds/repohealth/internal/model"
	"github.com/spbuilds/repohealth/internal/scanner"
)

var lockfiles = []string{
	"package-lock.json",
	"yarn.lock",
	"pnpm-lock.yaml",
	"Pipfile.lock",
	"poetry.lock",
	"Cargo.lock",
	"go.sum",
	"Gemfile.lock",
	"composer.lock",
}

var manifests = []string{
	"package.json",
	"go.mod",
	"Cargo.toml",
	"pyproject.toml",
	"requirements.txt",
	"Gemfile",
	"pom.xml",
	"build.gradle",
}

// DEP-01: Lockfile exists
type LockfileExistsCheck struct{}

func (c *LockfileExistsCheck) ID() string       { return "DEP-01" }
func (c *LockfileExistsCheck) Category() string { return "deps" }
func (c *LockfileExistsCheck) Name() string     { return "Lockfile exists" }
func (c *LockfileExistsCheck) MaxPoints() int   { return 4 }

func (c *LockfileExistsCheck) Run(ctx *model.ScanContext) model.CheckResult {
	found, ok := ctx.HasRootFile(lockfiles...)
	if ok {
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusFull, Points: c.MaxPoints(), MaxPoints: c.MaxPoints(),
			Details: found + " found",
		}
	}
	return model.CheckResult{
		ID: c.ID(), Category: c.Category(), Name: c.Name(),
		Status: model.StatusNone, Points: 0, MaxPoints: c.MaxPoints(),
		Details:    "No lockfile found",
		Suggestion: "Commit a lockfile to ensure reproducible dependency installs",
	}
}

// DEP-02: Package manager detected
type PackageManagerCheck struct{}

func (c *PackageManagerCheck) ID() string       { return "DEP-02" }
func (c *PackageManagerCheck) Category() string { return "deps" }
func (c *PackageManagerCheck) Name() string     { return "Package manager detected" }
func (c *PackageManagerCheck) MaxPoints() int   { return 2 }

func (c *PackageManagerCheck) Run(ctx *model.ScanContext) model.CheckResult {
	found, ok := ctx.HasRootFile(manifests...)
	if ok {
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusFull, Points: c.MaxPoints(), MaxPoints: c.MaxPoints(),
			Details: found + " found",
		}
	}
	return model.CheckResult{
		ID: c.ID(), Category: c.Category(), Name: c.Name(),
		Status: model.StatusNone, Points: 0, MaxPoints: c.MaxPoints(),
		Details:    "No package manager manifest found",
		Suggestion: "Add a manifest file (e.g. package.json, go.mod, Cargo.toml)",
	}
}

// DEP-03: Lockfile freshness
type LockfileFreshnessCheck struct{}

func (c *LockfileFreshnessCheck) ID() string       { return "DEP-03" }
func (c *LockfileFreshnessCheck) Category() string { return "deps" }
func (c *LockfileFreshnessCheck) Name() string     { return "Lockfile freshness" }
func (c *LockfileFreshnessCheck) MaxPoints() int   { return 2 }

func (c *LockfileFreshnessCheck) Run(ctx *model.ScanContext) model.CheckResult {
	lockPath, ok := ctx.HasRootFile(lockfiles...)
	if !ok {
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusSkipped, Points: 0, MaxPoints: c.MaxPoints(),
			Details: "No lockfile found",
		}
	}

	fullPath := filepath.Join(ctx.RepoPath, lockPath)
	info, err := os.Stat(fullPath)
	if err != nil {
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusSkipped, Points: 0, MaxPoints: c.MaxPoints(),
			Details: "Could not stat lockfile",
		}
	}

	age := time.Since(info.ModTime())
	days := int(age.Hours() / 24)

	switch {
	case days <= 90:
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusFull, Points: c.MaxPoints(), MaxPoints: c.MaxPoints(),
			Details: lockPath + " updated within 90 days",
		}
	case days <= 180:
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusPartial, Points: 1, MaxPoints: c.MaxPoints(),
			Details:    lockPath + " not updated in over 90 days",
			Suggestion: "Update your dependencies to keep the lockfile current",
		}
	default:
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusNone, Points: 0, MaxPoints: c.MaxPoints(),
			Details:    lockPath + " not updated in over 180 days",
			Suggestion: "Lockfile is stale — run your package manager update command",
		}
	}
}

// DEP-05: Dependency count
type DependencyCountCheck struct{}

func (c *DependencyCountCheck) ID() string       { return "DEP-05" }
func (c *DependencyCountCheck) Category() string { return "deps" }
func (c *DependencyCountCheck) Name() string     { return "Dependency count" }
func (c *DependencyCountCheck) MaxPoints() int   { return 1 }

func (c *DependencyCountCheck) Run(ctx *model.ScanContext) model.CheckResult {
	count, ok := parseDependencyCount(ctx)
	if !ok {
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusSkipped, Points: 0, MaxPoints: c.MaxPoints(),
			Details: "Could not parse dependency manifest",
		}
	}

	switch {
	case count < 50:
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusFull, Points: c.MaxPoints(), MaxPoints: c.MaxPoints(),
			Details: formatDepCount(count),
		}
	case count <= 100:
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusPartial, Points: 0, MaxPoints: c.MaxPoints(),
			Details:    formatDepCount(count),
			Suggestion: "Consider auditing dependencies for unused packages",
		}
	default:
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusNone, Points: 0, MaxPoints: c.MaxPoints(),
			Details:    formatDepCount(count),
			Suggestion: "High dependency count — audit for unused or redundant packages",
		}
	}
}

func formatDepCount(n int) string {
	if n == 1 {
		return "1 dependency declared"
	}
	s := ""
	for _, c := range []rune{'0' + rune(n/100%10), '0' + rune(n/10%10), '0' + rune(n%10)} {
		if c != '0' || s != "" {
			s += string(c)
		}
	}
	if s == "" {
		s = "0"
	}
	return s + " dependencies declared"
}

func parseDependencyCount(ctx *model.ScanContext) (int, bool) {
	// Try package.json first
	if _, ok := ctx.HasRootFile("package.json"); ok {
		return parsePackageJSON(ctx)
	}
	// Try go.mod
	if _, ok := ctx.HasRootFile("go.mod"); ok {
		return parseGoMod(ctx)
	}
	return 0, false
}

func parsePackageJSON(ctx *model.ScanContext) (int, bool) {
	lines, err := scanner.ReadFileLines(ctx.RepoPath, "package.json")
	if err != nil || lines == nil {
		return 0, false
	}

	content := strings.Join(lines, "\n")
	count := 0
	for _, section := range []string{"\"dependencies\"", "\"devDependencies\""} {
		idx := strings.Index(content, section)
		if idx < 0 {
			continue
		}
		// Find the opening brace of this section
		start := strings.Index(content[idx:], "{")
		if start < 0 {
			continue
		}
		start += idx + 1
		// Find the matching closing brace
		depth := 1
		i := start
		for i < len(content) && depth > 0 {
			switch content[i] {
			case '{':
				depth++
			case '}':
				depth--
			}
			i++
		}
		block := content[start : i-1]
		// Count occurrences of `":` which marks each key-value pair
		count += strings.Count(block, "\":")
	}
	return count, true
}

func parseGoMod(ctx *model.ScanContext) (int, bool) {
	lines, err := scanner.ReadFileLines(ctx.RepoPath, "go.mod")
	if err != nil || lines == nil {
		return 0, false
	}

	count := 0
	inRequire := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "require (" {
			inRequire = true
			continue
		}
		if inRequire {
			if trimmed == ")" {
				inRequire = false
				continue
			}
			// Single-line require outside block: "require module v1.0.0"
		}
		// Count lines in require block that start with whitespace (indented deps)
		if inRequire && len(line) > 0 && unicode.IsSpace(rune(line[0])) && trimmed != "" {
			count++
		}
		// Also count single-line require statements outside blocks
		if !inRequire && strings.HasPrefix(trimmed, "require ") && !strings.Contains(trimmed, "(") {
			count++
		}
	}
	return count, true
}
