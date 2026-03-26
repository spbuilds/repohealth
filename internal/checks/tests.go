package checks

import (
	"fmt"
	"strings"

	"github.com/spbuilds/repohealth/internal/model"
)

// TST-01: Test files exist
type TestFilesExistCheck struct{}

func (c *TestFilesExistCheck) ID() string       { return "TST-01" }
func (c *TestFilesExistCheck) Category() string { return "tests" }
func (c *TestFilesExistCheck) Name() string     { return "Test files detected" }
func (c *TestFilesExistCheck) MaxPoints() int   { return 8 }

func (c *TestFilesExistCheck) Run(ctx *model.ScanContext) model.CheckResult {
	count := ctx.CountFilesMatching(
		"*_test.go",
		"*_test.py", "test_*.py",
		"*.test.ts", "*.test.js", "*.test.tsx", "*.test.jsx",
		"*.spec.ts", "*.spec.js", "*.spec.tsx", "*.spec.jsx",
		"*_test.rs",
		"*_test.rb",
		"*Test.java", "*Tests.java",
	)

	if count > 0 {
		detail := fmt.Sprintf("%d test files", count)
		if count == 1 {
			detail = "1 test file"
		}
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusFull, Points: c.MaxPoints(), MaxPoints: c.MaxPoints(),
			Details: detail,
		}
	}
	return model.CheckResult{
		ID: c.ID(), Category: c.Category(), Name: c.Name(),
		Status: model.StatusNone, Points: 0, MaxPoints: c.MaxPoints(),
		Details:    "No test files found",
		Suggestion: "Add test files for your code",
	}
}

// TST-02: Test directory exists
type TestDirExistsCheck struct{}

func (c *TestDirExistsCheck) ID() string       { return "TST-02" }
func (c *TestDirExistsCheck) Category() string { return "tests" }
func (c *TestDirExistsCheck) Name() string     { return "Test directory exists" }
func (c *TestDirExistsCheck) MaxPoints() int   { return 2 }

func (c *TestDirExistsCheck) Run(ctx *model.ScanContext) model.CheckResult {
	testDirs := []string{"test", "tests", "__tests__", "spec"}
	for _, d := range ctx.Dirs {
		base := d
		if idx := strings.LastIndex(d, "/"); idx >= 0 {
			base = d[idx+1:]
		}
		for _, td := range testDirs {
			if base == td {
				return model.CheckResult{
					ID: c.ID(), Category: c.Category(), Name: c.Name(),
					Status: model.StatusFull, Points: c.MaxPoints(), MaxPoints: c.MaxPoints(),
					Details: d + "/",
				}
			}
		}
	}
	return model.CheckResult{
		ID: c.ID(), Category: c.Category(), Name: c.Name(),
		Status: model.StatusNone, Points: 0, MaxPoints: c.MaxPoints(),
		Details:    "No test directory found",
		Suggestion: "Create a test/ or tests/ directory",
	}
}

// TST-03: Test framework configured
type TestFrameworkCheck struct{}

func (c *TestFrameworkCheck) ID() string       { return "TST-03" }
func (c *TestFrameworkCheck) Category() string { return "tests" }
func (c *TestFrameworkCheck) Name() string     { return "Test framework configured" }
func (c *TestFrameworkCheck) MaxPoints() int   { return 4 }

func (c *TestFrameworkCheck) Run(ctx *model.ScanContext) model.CheckResult {
	frameworks := map[string]string{
		"jest.config.js":    "Jest",
		"jest.config.ts":    "Jest",
		"jest.config.mjs":   "Jest",
		"jest.config.cjs":   "Jest",
		"vitest.config.ts":  "Vitest",
		"vitest.config.js":  "Vitest",
		"vitest.config.mts": "Vitest",
		"pytest.ini":        "pytest",
		".mocharc.yml":      "Mocha",
		".mocharc.yaml":     "Mocha",
		".mocharc.js":       "Mocha",
		"phpunit.xml":       "PHPUnit",
		"phpunit.xml.dist":  "PHPUnit",
	}

	for file, name := range frameworks {
		if _, ok := ctx.HasFile(file); ok {
			return model.CheckResult{
				ID: c.ID(), Category: c.Category(), Name: c.Name(),
				Status: model.StatusFull, Points: c.MaxPoints(), MaxPoints: c.MaxPoints(),
				Details: name,
			}
		}
	}

	if _, ok := ctx.HasFile("pyproject.toml"); ok {
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusFull, Points: c.MaxPoints(), MaxPoints: c.MaxPoints(),
			Details: "pyproject.toml (pytest)",
		}
	}

	count := ctx.CountFilesMatching("*_test.go")
	if count > 0 {
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusFull, Points: c.MaxPoints(), MaxPoints: c.MaxPoints(),
			Details: "go test (built-in)",
		}
	}

	return model.CheckResult{
		ID: c.ID(), Category: c.Category(), Name: c.Name(),
		Status: model.StatusNone, Points: 0, MaxPoints: c.MaxPoints(),
		Details:    "No test framework configuration found",
		Suggestion: "Configure a test framework for your project",
	}
}
