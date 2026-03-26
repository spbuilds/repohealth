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

func daysAgoText(days int) string {
	if days == 0 {
		return "today"
	}
	if days == 1 {
		return "1 day ago"
	}
	return fmt.Sprintf("%d days ago", days)
}
