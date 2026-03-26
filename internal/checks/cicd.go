package checks

import "github.com/spbuilds/repohealth/internal/model"

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
