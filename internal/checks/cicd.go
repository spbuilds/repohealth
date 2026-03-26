package checks

import (
	"strings"

	"github.com/spbuilds/repohealth/internal/model"
	"github.com/spbuilds/repohealth/internal/scanner"
)

// CI-01: CI configuration exists
type CIConfigExistsCheck struct{}

func (c *CIConfigExistsCheck) ID() string       { return "CI-01" }
func (c *CIConfigExistsCheck) Category() string { return "cicd" }
func (c *CIConfigExistsCheck) Name() string     { return "CI configuration exists" }
func (c *CIConfigExistsCheck) MaxPoints() int   { return 6 }

func (c *CIConfigExistsCheck) Run(ctx *model.ScanContext) model.CheckResult {
	// Check for CI directories
	ciDirs := map[string]string{
		".github/workflows": "GitHub Actions",
		".circleci":         "CircleCI",
		".buildkite":        "Buildkite",
	}
	for dir, name := range ciDirs {
		if _, ok := ctx.HasDir(dir); ok {
			return model.CheckResult{
				ID: c.ID(), Category: c.Category(), Name: c.Name(),
				Status: model.StatusFull, Points: c.MaxPoints(), MaxPoints: c.MaxPoints(),
				Details: name,
			}
		}
	}

	// Check for CI files
	ciFiles := map[string]string{
		".gitlab-ci.yml":          "GitLab CI",
		"Jenkinsfile":             "Jenkins",
		".travis.yml":             "Travis CI",
		"bitbucket-pipelines.yml": "Bitbucket Pipelines",
		"azure-pipelines.yml":     "Azure Pipelines",
		"Taskfile.yml":            "Taskfile",
	}
	for file, name := range ciFiles {
		if _, ok := ctx.HasFile(file); ok {
			return model.CheckResult{
				ID: c.ID(), Category: c.Category(), Name: c.Name(),
				Status: model.StatusFull, Points: c.MaxPoints(), MaxPoints: c.MaxPoints(),
				Details: name,
			}
		}
	}

	return model.CheckResult{
		ID: c.ID(), Category: c.Category(), Name: c.Name(),
		Status: model.StatusNone, Points: 0, MaxPoints: c.MaxPoints(),
		Details:    "No CI/CD configuration found",
		Suggestion: "Add CI/CD configuration (GitHub Actions recommended)",
	}
}

// readCIConfigs collects all lines from all CI config files found in ctx.
func readCIConfigs(ctx *model.ScanContext) []string {
	var allLines []string
	for _, f := range ctx.Files {
		if strings.HasPrefix(f.Path, ".github/workflows/") && strings.HasSuffix(f.Name, ".yml") {
			lines, _ := scanner.ReadFileLines(ctx.RepoPath, f.Path)
			allLines = append(allLines, lines...)
		}
	}
	for _, ciFile := range []string{".gitlab-ci.yml", ".travis.yml", "Jenkinsfile"} {
		if _, ok := ctx.HasFile(ciFile); ok {
			lines, _ := scanner.ReadFileLines(ctx.RepoPath, ciFile)
			allLines = append(allLines, lines...)
		}
	}
	return allLines
}

// hasCIConfig returns true if the context contains any CI configuration.
func hasCIConfig(ctx *model.ScanContext) bool {
	ciDirs := []string{".github/workflows", ".circleci", ".buildkite"}
	for _, d := range ciDirs {
		if _, ok := ctx.HasDir(d); ok {
			return true
		}
	}
	ciFiles := []string{".gitlab-ci.yml", "Jenkinsfile", ".travis.yml",
		"bitbucket-pipelines.yml", "azure-pipelines.yml", "Taskfile.yml"}
	for _, f := range ciFiles {
		if _, ok := ctx.HasFile(f); ok {
			return true
		}
	}
	return false
}

// containsAny returns true if any line in lines contains any of the given substrings.
func containsAny(lines []string, substrings []string) bool {
	for _, line := range lines {
		lower := strings.ToLower(line)
		for _, s := range substrings {
			if strings.Contains(lower, s) {
				return true
			}
		}
	}
	return false
}

// ciContainsAny returns true if any CI config line contains any of the given patterns (case-insensitive).
func ciContainsAny(lines []string, patterns []string) bool {
	for _, line := range lines {
		lower := strings.ToLower(line)
		for _, p := range patterns {
			if strings.Contains(lower, p) {
				return true
			}
		}
	}
	return false
}

// CI-02: CI runs tests
type CIRunsTestsCheck struct{}

func (c *CIRunsTestsCheck) ID() string       { return "CI-02" }
func (c *CIRunsTestsCheck) Category() string { return "cicd" }
func (c *CIRunsTestsCheck) Name() string     { return "CI runs tests" }
func (c *CIRunsTestsCheck) MaxPoints() int   { return 4 }

func (c *CIRunsTestsCheck) Run(ctx *model.ScanContext) model.CheckResult {
	if !hasCIConfig(ctx) {
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusSkipped, Points: 0, MaxPoints: c.MaxPoints(),
			Details: "No CI configuration found",
		}
	}

	lines := readCIConfigs(ctx)
	testPatterns := []string{
		"npm test", "pnpm test", "yarn test", "bun test",
		"vitest", "jest", "mocha", "ava",
		"pytest", "python -m pytest", "tox",
		"go test", "gotestsum",
		"cargo test",
		"mvn test", "gradle test",
		"make test", "turbo test", "nx test",
		"rake test", "rspec", "bundle exec rspec",
		"phpunit",
	}

	if ciContainsAny(lines, testPatterns) {
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusFull, Points: c.MaxPoints(), MaxPoints: c.MaxPoints(),
			Details: "Test command found in CI",
		}
	}
	return model.CheckResult{
		ID: c.ID(), Category: c.Category(), Name: c.Name(),
		Status: model.StatusNone, Points: 0, MaxPoints: c.MaxPoints(),
		Details:    "No test command found in CI configuration",
		Suggestion: "Add a test step to your CI pipeline",
	}
}

// CI-03: CI runs linter
type CIRunsLinterCheck struct{}

func (c *CIRunsLinterCheck) ID() string       { return "CI-03" }
func (c *CIRunsLinterCheck) Category() string { return "cicd" }
func (c *CIRunsLinterCheck) Name() string     { return "CI runs linter" }
func (c *CIRunsLinterCheck) MaxPoints() int   { return 3 }

func (c *CIRunsLinterCheck) Run(ctx *model.ScanContext) model.CheckResult {
	if !hasCIConfig(ctx) {
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusSkipped, Points: 0, MaxPoints: c.MaxPoints(),
			Details: "No CI configuration found",
		}
	}

	lines := readCIConfigs(ctx)
	lintPatterns := []string{
		"eslint", "biome", "oxlint", "prettier --check",
		"ruff", "flake8", "pylint", "mypy", "black --check",
		"golangci-lint", "go vet", "staticcheck",
		"clippy", "cargo clippy",
		"checkstyle", "spotbugs",
		"rubocop", "standardrb",
		"make lint", "pnpm lint", "yarn lint", "npm run lint",
	}

	if ciContainsAny(lines, lintPatterns) {
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusFull, Points: c.MaxPoints(), MaxPoints: c.MaxPoints(),
			Details: "Linter found in CI",
		}
	}
	return model.CheckResult{
		ID: c.ID(), Category: c.Category(), Name: c.Name(),
		Status: model.StatusNone, Points: 0, MaxPoints: c.MaxPoints(),
		Details:    "No linter found in CI configuration",
		Suggestion: "Add a lint step to your CI pipeline",
	}
}

// CI-04: CI runs build
type CIRunsBuildCheck struct{}

func (c *CIRunsBuildCheck) ID() string       { return "CI-04" }
func (c *CIRunsBuildCheck) Category() string { return "cicd" }
func (c *CIRunsBuildCheck) Name() string     { return "CI runs build" }
func (c *CIRunsBuildCheck) MaxPoints() int   { return 2 }

func (c *CIRunsBuildCheck) Run(ctx *model.ScanContext) model.CheckResult {
	if !hasCIConfig(ctx) {
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusSkipped, Points: 0, MaxPoints: c.MaxPoints(),
			Details: "No CI configuration found",
		}
	}

	lines := readCIConfigs(ctx)
	buildPatterns := []string{
		"npm run build", "pnpm build", "yarn build", "bun build",
		"vite build", "next build", "nuxt build",
		"go build", "goreleaser",
		"cargo build",
		"mvn package", "gradle build", "gradle assemble",
		"make build", "turbo build", "nx build",
		"docker build",
	}

	if ciContainsAny(lines, buildPatterns) {
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusFull, Points: c.MaxPoints(), MaxPoints: c.MaxPoints(),
			Details: "Build command found in CI",
		}
	}
	return model.CheckResult{
		ID: c.ID(), Category: c.Category(), Name: c.Name(),
		Status: model.StatusNone, Points: 0, MaxPoints: c.MaxPoints(),
		Details:    "No build command found in CI configuration",
		Suggestion: "Add a build step to your CI pipeline",
	}
}
