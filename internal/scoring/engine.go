package scoring

import (
	"math"
	"time"

	"github.com/spbuilds/repohealth/internal/model"
)

// categoryLabels maps category IDs to display labels.
var categoryLabels = map[string]string{
	"docs":     "Documentation",
	"tests":    "Testing",
	"cicd":     "CI/CD",
	"deps":     "Dependencies",
	"security": "Security",
	"stats":    "Code Statistics",
	"activity": "Activity",
	"todo":     "TODO / Technical Debt",
}

// Score aggregates check results into a complete report.
func Score(results []model.CheckResult, repoPath string, languages map[string]int, filesAnalyzed int, startTime time.Time) *model.Report {
	report := &model.Report{
		RepoPath:      repoPath,
		Timestamp:     time.Now().UTC().Format(time.RFC3339),
		DurationMs:    time.Since(startTime).Milliseconds(),
		MaxScore:      100,
		Checks:        results,
		Languages:     languages,
		FilesAnalyzed: filesAnalyzed,
	}

	// Group results by category
	catResults := make(map[string][]model.CheckResult)
	for _, r := range results {
		catResults[r.Category] = append(catResults[r.Category], r)
	}

	// Compute per-category scores
	totalScore := 0
	totalMax := 0

	// Process categories in display order
	catOrder := []string{"docs", "tests", "cicd", "deps", "security", "stats", "activity", "todo"}

	for _, cat := range catOrder {
		checks, ok := catResults[cat]
		if !ok {
			continue
		}

		catScore := 0
		catMax := 0
		skippedMax := 0

		for _, cr := range checks {
			if cr.Status == model.StatusSkipped {
				skippedMax += cr.MaxPoints
				continue
			}
			catScore += cr.Points
			catMax += cr.MaxPoints
		}

		// BUG FIX: If all checks in a category are skipped, exclude the
		// category entirely from scoring (don't penalize for unavailable data)
		if catMax == 0 {
			continue
		}

		// Redistribute skipped check points proportionally
		effectiveMax := catMax + skippedMax
		if skippedMax > 0 {
			factor := float64(effectiveMax) / float64(catMax)
			catScore = int(math.Round(float64(catScore) * factor))
		}

		label := categoryLabels[cat]
		if label == "" {
			label = cat
		}

		report.Categories = append(report.Categories, model.CategoryResult{
			Name:     cat,
			Label:    label,
			Score:    catScore,
			MaxScore: effectiveMax,
		})

		totalScore += catScore
		totalMax += effectiveMax
	}

	// Normalize to 0-100 scale
	if totalMax > 0 {
		report.Score = int(math.Round(float64(totalScore) / float64(totalMax) * 100))
	}

	report.Grade = Grade(report.Score)
	report.Suggestions = Recommendations(results)

	return report
}
