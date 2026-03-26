package checks

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spbuilds/repohealth/internal/model"
	"github.com/spbuilds/repohealth/internal/scanner"
)

var secretPatterns = []struct {
	name    string
	pattern *regexp.Regexp
}{
	{"AWS access key", regexp.MustCompile(`AKIA[0-9A-Z]{16}`)},
	{"GitHub token", regexp.MustCompile(`ghp_[A-Za-z0-9]{36}`)},
	{"GitHub OAuth", regexp.MustCompile(`gho_[A-Za-z0-9]{36}`)},
	{"Slack token", regexp.MustCompile(`xox[baprs]-[A-Za-z0-9-]{10,}`)},
	{"Stripe secret key", regexp.MustCompile(`sk_live_[A-Za-z0-9]{24,}`)},
	{"Google API key", regexp.MustCompile(`AIza[0-9A-Za-z_-]{35}`)},
	{"Private key block", regexp.MustCompile(`-----BEGIN\s+(RSA|EC|OPENSSH|DSA)?\s*PRIVATE KEY-----`)},
	{"Generic secret assignment", regexp.MustCompile(`(?i)(password|secret|api_key|apikey|token)\s*[:=]\s*['"][^'"]{8,}['"]`)},
}

// SEC-02: No secrets in repo
type NoSecretsCheck struct{}

func (c *NoSecretsCheck) ID() string       { return "SEC-02" }
func (c *NoSecretsCheck) Category() string { return "security" }
func (c *NoSecretsCheck) Name() string     { return "No secrets in repo" }
func (c *NoSecretsCheck) MaxPoints() int   { return 4 }

func (c *NoSecretsCheck) Run(ctx *model.ScanContext) model.CheckResult {
	// Check for .env files (not .env.example)
	for _, f := range ctx.Files {
		if f.IsDir {
			continue
		}
		base := filepath.Base(f.Path)
		if base == ".env" {
			return model.CheckResult{
				ID: c.ID(), Category: c.Category(), Name: c.Name(),
				Status: model.StatusNone, Points: 0, MaxPoints: c.MaxPoints(),
				Details:    ".env file found in repository",
				Suggestion: "Remove .env from the repository and add it to .gitignore",
			}
		}
	}

	// Scan source files for secret patterns
	for _, f := range ctx.Files {
		if f.IsDir {
			continue
		}
		ext := strings.ToLower(filepath.Ext(f.Name))
		if !sourceExtensions[ext] {
			continue
		}

		lines, err := scanner.ReadFileLines(ctx.RepoPath, f.Path)
		if err != nil || lines == nil {
			continue
		}

		for _, line := range lines {
			for _, sp := range secretPatterns {
				if sp.pattern.MatchString(line) {
					return model.CheckResult{
						ID: c.ID(), Category: c.Category(), Name: c.Name(),
						Status: model.StatusNone, Points: 0, MaxPoints: c.MaxPoints(),
						Details:    sp.name + " detected",
						Suggestion: "Remove secrets from source code and use environment variables",
					}
				}
			}
		}
	}

	return model.CheckResult{
		ID: c.ID(), Category: c.Category(), Name: c.Name(),
		Status: model.StatusFull, Points: c.MaxPoints(), MaxPoints: c.MaxPoints(),
		Details: "No secret patterns detected in source files",
	}
}

// SEC-03: .gitignore covers secrets
type GitignoreSecretsCheck struct{}

func (c *GitignoreSecretsCheck) ID() string       { return "SEC-03" }
func (c *GitignoreSecretsCheck) Category() string { return "security" }
func (c *GitignoreSecretsCheck) Name() string     { return ".gitignore covers secrets" }
func (c *GitignoreSecretsCheck) MaxPoints() int   { return 2 }

func (c *GitignoreSecretsCheck) Run(ctx *model.ScanContext) model.CheckResult {
	lines, err := scanner.ReadFileLines(ctx.RepoPath, ".gitignore")
	if err != nil || lines == nil {
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusNone, Points: 0, MaxPoints: c.MaxPoints(),
			Details:    "No .gitignore found",
			Suggestion: "Add a .gitignore with entries for .env, *.pem, *.key, and credentials",
		}
	}

	content := strings.Join(lines, "\n")
	required := []string{".env", "*.pem", "*.key", "credentials"}
	covered := 0
	for _, pattern := range required {
		if strings.Contains(content, pattern) {
			covered++
		}
	}

	switch {
	case covered == len(required):
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusFull, Points: c.MaxPoints(), MaxPoints: c.MaxPoints(),
			Details: ".gitignore covers common secret file patterns",
		}
	case covered > 0:
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusPartial, Points: 1, MaxPoints: c.MaxPoints(),
			Details:    ".gitignore covers some but not all secret patterns",
			Suggestion: "Add .env, *.pem, *.key, and credentials entries to .gitignore",
		}
	default:
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusNone, Points: 0, MaxPoints: c.MaxPoints(),
			Details:    ".gitignore does not cover secret file patterns",
			Suggestion: "Add .env, *.pem, *.key, and credentials entries to .gitignore",
		}
	}
}

// SEC-04: Dependency pinning
type DependencyPinningCheck struct{}

func (c *DependencyPinningCheck) ID() string       { return "SEC-04" }
func (c *DependencyPinningCheck) Category() string { return "security" }
func (c *DependencyPinningCheck) Name() string     { return "Dependency pinning" }
func (c *DependencyPinningCheck) MaxPoints() int   { return 2 }

func (c *DependencyPinningCheck) Run(ctx *model.ScanContext) model.CheckResult {
	found, ok := ctx.HasRootFile(lockfiles...)
	if ok {
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusFull, Points: c.MaxPoints(), MaxPoints: c.MaxPoints(),
			Details: found + " pins dependency versions",
		}
	}
	return model.CheckResult{
		ID: c.ID(), Category: c.Category(), Name: c.Name(),
		Status: model.StatusNone, Points: 0, MaxPoints: c.MaxPoints(),
		Details:    "No lockfile found — dependency versions are not pinned",
		Suggestion: "Commit a lockfile to pin dependency versions for reproducible builds",
	}
}

// SEC-05: Branch protection indicators
type BranchProtectionCheck struct{}

func (c *BranchProtectionCheck) ID() string       { return "SEC-05" }
func (c *BranchProtectionCheck) Category() string { return "security" }
func (c *BranchProtectionCheck) Name() string     { return "Branch protection indicators" }
func (c *BranchProtectionCheck) MaxPoints() int   { return 2 }

func (c *BranchProtectionCheck) Run(ctx *model.ScanContext) model.CheckResult {
	found, ok := ctx.HasRootFile(
		"CODEOWNERS",
		".github/CODEOWNERS",
		".github/branch-protection.yml",
	)
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
		Details:    "No branch protection indicators found",
		Suggestion: "Add a CODEOWNERS file or branch protection configuration",
	}
}
