package scoring

import "github.com/spbuilds/repohealth/internal/model"

// ImprovementStep represents a single improvement action with projected impact.
type ImprovementStep struct {
	Action         string `json:"action"`
	Impact         int    `json:"impact"`
	ProjectedScore int    `json:"projected_score"`
}

// ImprovementPlan calculates projected score improvements.
func ImprovementPlan(r *model.Report, maxSteps int) []ImprovementStep {
	if maxSteps <= 0 {
		maxSteps = 5
	}

	var steps []ImprovementStep
	projected := r.Score

	for i, s := range r.Suggestions {
		if i >= maxSteps {
			break
		}
		projected += s.Impact
		if projected > 100 {
			projected = 100
		}
		steps = append(steps, ImprovementStep{
			Action:         s.Message,
			Impact:         s.Impact,
			ProjectedScore: projected,
		})
	}

	return steps
}
