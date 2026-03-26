package checks

import "github.com/spbuilds/repohealth/internal/model"

// DOC-01: README exists
type ReadmeExistsCheck struct{}

func (c *ReadmeExistsCheck) ID() string       { return "DOC-01" }
func (c *ReadmeExistsCheck) Category() string { return "docs" }
func (c *ReadmeExistsCheck) Name() string     { return "README exists" }
func (c *ReadmeExistsCheck) MaxPoints() int   { return 4 }

func (c *ReadmeExistsCheck) Run(ctx *model.ScanContext) model.CheckResult {
	found, ok := ctx.HasFile("README.md", "README", "README.rst", "readme.md", "Readme.md")
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
		Details:    "No README found",
		Suggestion: "Add a README.md file describing your project",
	}
}

// DOC-02: README has content
type ReadmeContentCheck struct{}

func (c *ReadmeContentCheck) ID() string       { return "DOC-02" }
func (c *ReadmeContentCheck) Category() string { return "docs" }
func (c *ReadmeContentCheck) Name() string     { return "README has content" }
func (c *ReadmeContentCheck) MaxPoints() int   { return 2 }

func (c *ReadmeContentCheck) Run(ctx *model.ScanContext) model.CheckResult {
	size := ctx.FileSize("README.md", "README", "README.rst", "readme.md", "Readme.md")
	if size > 100 {
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusFull, Points: c.MaxPoints(), MaxPoints: c.MaxPoints(),
			Details: "README has substantive content",
		}
	}
	if size >= 0 {
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusPartial, Points: c.MaxPoints() / 2, MaxPoints: c.MaxPoints(),
			Details:    "README exists but has minimal content",
			Suggestion: "Expand README with project description, usage, and installation instructions",
		}
	}
	return model.CheckResult{
		ID: c.ID(), Category: c.Category(), Name: c.Name(),
		Status: model.StatusNone, Points: 0, MaxPoints: c.MaxPoints(),
		Details:    "No README found",
		Suggestion: "Add a README.md with project description, usage, and installation instructions",
	}
}

// DOC-03: LICENSE exists
type LicenseExistsCheck struct{}

func (c *LicenseExistsCheck) ID() string       { return "DOC-03" }
func (c *LicenseExistsCheck) Category() string { return "docs" }
func (c *LicenseExistsCheck) Name() string     { return "LICENSE exists" }
func (c *LicenseExistsCheck) MaxPoints() int   { return 3 }

func (c *LicenseExistsCheck) Run(ctx *model.ScanContext) model.CheckResult {
	found, ok := ctx.HasFile("LICENSE", "LICENSE.md", "LICENCE", "COPYING", "LICENSE.txt")
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
		Details:    "No LICENSE file found",
		Suggestion: "Add a LICENSE file (MIT, Apache 2.0, or BSD recommended)",
	}
}

// DOC-04: CONTRIBUTING exists
type ContributingExistsCheck struct{}

func (c *ContributingExistsCheck) ID() string       { return "DOC-04" }
func (c *ContributingExistsCheck) Category() string { return "docs" }
func (c *ContributingExistsCheck) Name() string     { return "CONTRIBUTING exists" }
func (c *ContributingExistsCheck) MaxPoints() int   { return 2 }

func (c *ContributingExistsCheck) Run(ctx *model.ScanContext) model.CheckResult {
	found, ok := ctx.HasFile("CONTRIBUTING.md", "CONTRIBUTING", ".github/CONTRIBUTING.md")
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
		Details:    "No CONTRIBUTING file found",
		Suggestion: "Add a CONTRIBUTING.md with contribution guidelines",
	}
}

// DOC-05: CODE_OF_CONDUCT exists
type CodeOfConductExistsCheck struct{}

func (c *CodeOfConductExistsCheck) ID() string       { return "DOC-05" }
func (c *CodeOfConductExistsCheck) Category() string { return "docs" }
func (c *CodeOfConductExistsCheck) Name() string     { return "CODE_OF_CONDUCT exists" }
func (c *CodeOfConductExistsCheck) MaxPoints() int   { return 1 }

func (c *CodeOfConductExistsCheck) Run(ctx *model.ScanContext) model.CheckResult {
	found, ok := ctx.HasFile("CODE_OF_CONDUCT.md", ".github/CODE_OF_CONDUCT.md")
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
		Details:    "No CODE_OF_CONDUCT file found",
		Suggestion: "Add a CODE_OF_CONDUCT.md (Contributor Covenant recommended)",
	}
}

// DOC-06: SECURITY.md exists
type SecurityExistsCheck struct{}

func (c *SecurityExistsCheck) ID() string       { return "DOC-06" }
func (c *SecurityExistsCheck) Category() string { return "docs" }
func (c *SecurityExistsCheck) Name() string     { return "SECURITY.md exists" }
func (c *SecurityExistsCheck) MaxPoints() int   { return 2 }

func (c *SecurityExistsCheck) Run(ctx *model.ScanContext) model.CheckResult {
	found, ok := ctx.HasFile("SECURITY.md", ".github/SECURITY.md")
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
		Details:    "No SECURITY.md found",
		Suggestion: "Add a SECURITY.md with vulnerability reporting instructions",
	}
}

// DOC-07: CHANGELOG exists
type ChangelogExistsCheck struct{}

func (c *ChangelogExistsCheck) ID() string       { return "DOC-07" }
func (c *ChangelogExistsCheck) Category() string { return "docs" }
func (c *ChangelogExistsCheck) Name() string     { return "CHANGELOG exists" }
func (c *ChangelogExistsCheck) MaxPoints() int   { return 1 }

func (c *ChangelogExistsCheck) Run(ctx *model.ScanContext) model.CheckResult {
	found, ok := ctx.HasFile("CHANGELOG.md", "CHANGELOG", "HISTORY.md", "CHANGES.md")
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
		Details:    "No CHANGELOG found",
		Suggestion: "Add a CHANGELOG.md to track releases",
	}
}
