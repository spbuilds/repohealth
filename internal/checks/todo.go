package checks

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/spbuilds/repohealth/internal/model"
	"github.com/spbuilds/repohealth/internal/scanner"
)

type todoStats struct {
	count    int
	totalLOC int
	critical bool
}

var todoCache struct {
	sync.Mutex
	key   string
	stats todoStats
}

func scanTodos(ctx *model.ScanContext) todoStats {
	todoCache.Lock()
	defer todoCache.Unlock()

	if todoCache.key == ctx.RepoPath {
		return todoCache.stats
	}

	var stats todoStats
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

		stats.totalLOC += len(lines)
		for _, line := range lines {
			upper := strings.ToUpper(line)
			if strings.Contains(upper, "TODO") ||
				strings.Contains(upper, "FIXME") ||
				strings.Contains(upper, "HACK") ||
				strings.Contains(upper, "XXX") {
				stats.count++
				if strings.Contains(upper, "SECURITY") ||
					strings.Contains(upper, "VULNERABILITY") ||
					strings.Contains(upper, "UNSAFE") {
					stats.critical = true
				}
			}
		}
	}

	todoCache.key = ctx.RepoPath
	todoCache.stats = stats
	return stats
}

// TODO-01: TODO/FIXME count
type TodoCountCheck struct{}

func (c *TodoCountCheck) ID() string       { return "TODO-01" }
func (c *TodoCountCheck) Category() string { return "todo" }
func (c *TodoCountCheck) Name() string     { return "TODO/FIXME count" }
func (c *TodoCountCheck) MaxPoints() int   { return 3 }

func (c *TodoCountCheck) Run(ctx *model.ScanContext) model.CheckResult {
	stats := scanTodos(ctx)
	n := stats.count

	details := fmt.Sprintf("%d TODO/FIXME markers found", n)

	switch {
	case n == 0:
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusFull, Points: c.MaxPoints(), MaxPoints: c.MaxPoints(),
			Details: details,
		}
	case n <= 20:
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusPartial, Points: 1, MaxPoints: c.MaxPoints(),
			Details:    details,
			Suggestion: "Resolve or track TODO/FIXME markers as issues",
		}
	default:
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusNone, Points: 0, MaxPoints: c.MaxPoints(),
			Details:    details,
			Suggestion: "More than 20 TODO/FIXME markers found — consider tracking them as issues",
		}
	}
}

// TODO-02: Density per KLOC
type TodoDensityCheck struct{}

func (c *TodoDensityCheck) ID() string       { return "TODO-02" }
func (c *TodoDensityCheck) Category() string { return "todo" }
func (c *TodoDensityCheck) Name() string     { return "TODO density per KLOC" }
func (c *TodoDensityCheck) MaxPoints() int   { return 2 }

func (c *TodoDensityCheck) Run(ctx *model.ScanContext) model.CheckResult {
	stats := scanTodos(ctx)

	if stats.totalLOC == 0 {
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusSkipped, Points: 0, MaxPoints: c.MaxPoints(),
			Details: "No source lines found",
		}
	}

	density := float64(stats.count) / (float64(stats.totalLOC) / 1000.0)

	switch {
	case density < 2.0:
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusFull, Points: c.MaxPoints(), MaxPoints: c.MaxPoints(),
			Details: fmt.Sprintf("%.1f TODO/FIXME markers per KLOC", density),
		}
	case density <= 5.0:
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusPartial, Points: 1, MaxPoints: c.MaxPoints(),
			Details:    fmt.Sprintf("%.1f TODO/FIXME markers per KLOC", density),
			Suggestion: "TODO density is moderate — aim for under 2 per KLOC",
		}
	default:
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusNone, Points: 0, MaxPoints: c.MaxPoints(),
			Details:    fmt.Sprintf("%.1f TODO/FIXME markers per KLOC", density),
			Suggestion: "High TODO density — resolve markers or track them as issues",
		}
	}
}

// TODO-03: No critical markers
type TodoCriticalCheck struct{}

func (c *TodoCriticalCheck) ID() string       { return "TODO-03" }
func (c *TodoCriticalCheck) Category() string { return "todo" }
func (c *TodoCriticalCheck) Name() string     { return "No critical TODO markers" }
func (c *TodoCriticalCheck) MaxPoints() int   { return 2 }

func (c *TodoCriticalCheck) Run(ctx *model.ScanContext) model.CheckResult {
	stats := scanTodos(ctx)

	if !stats.critical {
		return model.CheckResult{
			ID: c.ID(), Category: c.Category(), Name: c.Name(),
			Status: model.StatusFull, Points: c.MaxPoints(), MaxPoints: c.MaxPoints(),
			Details: "No critical TODO markers found",
		}
	}
	return model.CheckResult{
		ID: c.ID(), Category: c.Category(), Name: c.Name(),
		Status: model.StatusNone, Points: 0, MaxPoints: c.MaxPoints(),
		Details:    "TODO/FIXME comments containing SECURITY, VULNERABILITY, or UNSAFE found",
		Suggestion: "Address security-related TODO markers immediately",
	}
}
