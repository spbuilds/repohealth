package scoring

import (
	"sort"

	"github.com/spbuilds/repohealth/internal/model"
)

// Recommendations generates improvement suggestions sorted by impact.
func Recommendations(results []model.CheckResult) []model.Suggestion {
	var suggestions []model.Suggestion

	for _, r := range results {
		if r.Status == model.StatusFull || r.Status == model.StatusSkipped {
			continue
		}
		if r.Suggestion == "" {
			continue
		}

		impact := r.MaxPoints - r.Points
		if impact <= 0 {
			continue
		}

		suggestions = append(suggestions, model.Suggestion{
			CheckID: r.ID,
			Impact:  impact,
			Message: r.Suggestion,
		})
	}

	sort.Slice(suggestions, func(i, j int) bool {
		return suggestions[i].Impact > suggestions[j].Impact
	})

	return suggestions
}
