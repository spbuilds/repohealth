package checks

import (
	"fmt"
	"time"

	"github.com/spbuilds/repohealth/internal/model"
	"github.com/spbuilds/repohealth/internal/scanner"
)

// ACT-01: Recent commit
type RecentCommitCheck struct{}

func (c *RecentCommitCheck) ID() string       { return "ACT-01" }
func (c *RecentCommitCheck) Category() string { return "activity" }
func (c *RecentCommitCheck) Name() string     { return "Recent commit" }
func (c *RecentCommitCheck) MaxPoints() int   { return 5 }

func (c *RecentCommitCheck) Run(ctx *model.ScanContext) model.CheckResult {
	if !ctx.GitAvailable {
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusSkipped, Points: 0, MaxPoints: c.MaxPoints(),
			Details: "Git not available",
		}
	}

	lastCommit, err := scanner.LastCommitDate(ctx.RepoPath)
	if err != nil {
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusSkipped, Points: 0, MaxPoints: c.MaxPoints(),
			Details: "Could not read git log",
		}
	}

	daysSince := int(time.Since(lastCommit).Hours() / 24)

	if daysSince <= 30 {
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusFull, Points: c.MaxPoints(), MaxPoints: c.MaxPoints(),
			Details: daysAgoText(daysSince),
		}
	}

	if daysSince <= 90 {
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusPartial, Points: c.MaxPoints() / 2, MaxPoints: c.MaxPoints(),
			Details:    daysAgoText(daysSince),
			Suggestion: "Repository has not been updated recently",
		}
	}

	return model.CheckResult{
		ID: c.ID(), Category: c.Category(), Name: c.Name(),
		Status: model.StatusNone, Points: 0, MaxPoints: c.MaxPoints(),
		Details:    daysAgoText(daysSince),
		Suggestion: "Repository appears inactive (last commit > 90 days ago)",
	}
}

// ACT-03: Contributor count
type ContributorCountCheck struct{}

func (c *ContributorCountCheck) ID() string       { return "ACT-03" }
func (c *ContributorCountCheck) Category() string { return "activity" }
func (c *ContributorCountCheck) Name() string     { return "Contributors" }
func (c *ContributorCountCheck) MaxPoints() int   { return 3 }

func (c *ContributorCountCheck) Run(ctx *model.ScanContext) model.CheckResult {
	if !ctx.GitAvailable {
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusSkipped, Points: 0, MaxPoints: c.MaxPoints(),
			Details: "Git not available",
		}
	}

	count, err := scanner.ContributorCount(ctx.RepoPath)
	if err != nil {
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusSkipped, Points: 0, MaxPoints: c.MaxPoints(),
			Details: "Could not count contributors",
		}
	}

	if count == 0 {
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusSkipped, Points: 0, MaxPoints: c.MaxPoints(),
			Details: "No commits found",
		}
	}

	if count > 5 {
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusFull, Points: c.MaxPoints(), MaxPoints: c.MaxPoints(),
			Details: fmt.Sprintf("%d contributors", count),
		}
	}

	if count >= 2 {
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusPartial, Points: c.MaxPoints() / 2, MaxPoints: c.MaxPoints(),
			Details:    fmt.Sprintf("%d contributors", count),
			Suggestion: "Increase contributor count to improve bus factor",
		}
	}

	return model.CheckResult{
		ID: c.ID(), Category: c.Category(), Name: c.Name(),
		Status: model.StatusNone, Points: 0, MaxPoints: c.MaxPoints(),
		Details:    "1 contributor",
		Suggestion: "Single contributor — consider inviting collaborators",
	}
}

// ACT-02: Commit frequency
type CommitFrequencyCheck struct{}

func (c *CommitFrequencyCheck) ID() string       { return "ACT-02" }
func (c *CommitFrequencyCheck) Category() string { return "activity" }
func (c *CommitFrequencyCheck) Name() string     { return "Commit frequency" }
func (c *CommitFrequencyCheck) MaxPoints() int   { return 3 }

func (c *CommitFrequencyCheck) Run(ctx *model.ScanContext) model.CheckResult {
	if !ctx.GitAvailable {
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusSkipped, Points: 0, MaxPoints: c.MaxPoints(),
			Details: "Git not available",
		}
	}

	count, err := scanner.CommitCountSince(ctx.RepoPath, 6)
	if err != nil {
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusSkipped, Points: 0, MaxPoints: c.MaxPoints(),
			Details: "Could not count commits",
		}
	}

	if count > 50 {
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusFull, Points: c.MaxPoints(), MaxPoints: c.MaxPoints(),
			Details: fmt.Sprintf("%d commits in last 6 months", count),
		}
	}
	if count >= 10 {
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusPartial, Points: 1, MaxPoints: c.MaxPoints(),
			Details:    fmt.Sprintf("%d commits in last 6 months", count),
			Suggestion: "Increase commit frequency to show active development",
		}
	}
	return model.CheckResult{
		ID: c.ID(), Category: c.Category(), Name: c.Name(),
		Status: model.StatusNone, Points: 0, MaxPoints: c.MaxPoints(),
		Details:    fmt.Sprintf("%d commits in last 6 months", count),
		Suggestion: "Repository has very few recent commits",
	}
}

// ACT-04: Release exists
type ReleaseExistsCheck struct{}

func (c *ReleaseExistsCheck) ID() string       { return "ACT-04" }
func (c *ReleaseExistsCheck) Category() string { return "activity" }
func (c *ReleaseExistsCheck) Name() string     { return "Release exists" }
func (c *ReleaseExistsCheck) MaxPoints() int   { return 2 }

func (c *ReleaseExistsCheck) Run(ctx *model.ScanContext) model.CheckResult {
	if !ctx.GitAvailable {
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusSkipped, Points: 0, MaxPoints: c.MaxPoints(),
			Details: "Git not available",
		}
	}

	count, err := scanner.TagCount(ctx.RepoPath)
	if err != nil {
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusSkipped, Points: 0, MaxPoints: c.MaxPoints(),
			Details: "Could not read git tags",
		}
	}

	if count > 0 {
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusFull, Points: c.MaxPoints(), MaxPoints: c.MaxPoints(),
			Details: fmt.Sprintf("%d release tags", count),
		}
	}
	return model.CheckResult{
		ID: c.ID(), Category: c.Category(), Name: c.Name(),
		Status: model.StatusNone, Points: 0, MaxPoints: c.MaxPoints(),
		Details:    "No release tags found",
		Suggestion: "Tag releases to make versioning visible (e.g. v1.0.0)",
	}
}

// ACT-05: Bus factor
type BusFactorCheck struct{}

func (c *BusFactorCheck) ID() string       { return "ACT-05" }
func (c *BusFactorCheck) Category() string { return "activity" }
func (c *BusFactorCheck) Name() string     { return "Bus factor" }
func (c *BusFactorCheck) MaxPoints() int   { return 2 }

func (c *BusFactorCheck) Run(ctx *model.ScanContext) model.CheckResult {
	if !ctx.GitAvailable {
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusSkipped, Points: 0, MaxPoints: c.MaxPoints(),
			Details: "Git not available",
		}
	}

	factor, err := scanner.BusFactor(ctx.RepoPath)
	if err != nil {
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusSkipped, Points: 0, MaxPoints: c.MaxPoints(),
			Details: "Could not compute bus factor",
		}
	}

	if factor == 0 {
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusSkipped, Points: 0, MaxPoints: c.MaxPoints(),
			Details: "No commits found",
		}
	}

	if factor >= 3 {
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusFull, Points: c.MaxPoints(), MaxPoints: c.MaxPoints(),
			Details: fmt.Sprintf("bus factor %d", factor),
		}
	}
	if factor == 2 {
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusPartial, Points: 1, MaxPoints: c.MaxPoints(),
			Details:    fmt.Sprintf("bus factor %d", factor),
			Suggestion: "Spread knowledge across more contributors",
		}
	}
	return model.CheckResult{
		ID: c.ID(), Category: c.Category(), Name: c.Name(),
		Status: model.StatusNone, Points: 0, MaxPoints: c.MaxPoints(),
		Details:    "bus factor 1",
		Suggestion: "Only 1 contributor holds >10% of commits — high bus factor risk",
	}
}

func daysAgoText(days int) string {
	if days == 0 {
		return "today"
	}
	if days == 1 {
		return "1 day ago"
	}
	return fmt.Sprintf("%d days ago", days)
}
