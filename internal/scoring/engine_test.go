package scoring

import (
	"testing"
	"time"

	"github.com/spbuilds/repohealth/internal/model"
)

func TestGrade(t *testing.T) {
	tests := []struct {
		score int
		want  string
	}{
		{100, "A+"},
		{95, "A+"},
		{94, "A"},
		{90, "A"},
		{89, "A-"},
		{85, "A-"},
		{84, "B+"},
		{80, "B+"},
		{79, "B"},
		{75, "B"},
		{74, "B-"},
		{70, "B-"},
		{69, "C+"},
		{65, "C+"},
		{64, "C"},
		{60, "C"},
		{59, "C-"},
		{55, "C-"},
		{54, "D"},
		{40, "D"},
		{39, "F"},
		{0, "F"},
	}

	for _, tt := range tests {
		got := Grade(tt.score)
		if got != tt.want {
			t.Errorf("Grade(%d) = %q, want %q", tt.score, got, tt.want)
		}
	}
}

func TestScoreAllFull(t *testing.T) {
	results := []model.CheckResult{
		{ID: "DOC-01", Category: "docs", Status: model.StatusFull, Points: 4, MaxPoints: 4},
		{ID: "DOC-03", Category: "docs", Status: model.StatusFull, Points: 3, MaxPoints: 3},
		{ID: "TST-01", Category: "tests", Status: model.StatusFull, Points: 8, MaxPoints: 8},
	}

	r := Score(results, "/tmp/test", nil, 10, time.Now())

	if r.Score != 100 {
		t.Errorf("Score = %d, want 100 (all full)", r.Score)
	}
	if r.Grade != "A+" {
		t.Errorf("Grade = %q, want A+", r.Grade)
	}
}

func TestScoreAllNone(t *testing.T) {
	results := []model.CheckResult{
		{ID: "DOC-01", Category: "docs", Status: model.StatusNone, Points: 0, MaxPoints: 4},
		{ID: "TST-01", Category: "tests", Status: model.StatusNone, Points: 0, MaxPoints: 8},
	}

	r := Score(results, "/tmp/test", nil, 10, time.Now())

	if r.Score != 0 {
		t.Errorf("Score = %d, want 0 (all none)", r.Score)
	}
	if r.Grade != "F" {
		t.Errorf("Grade = %q, want F", r.Grade)
	}
}

func TestScoreMixed(t *testing.T) {
	results := []model.CheckResult{
		{ID: "DOC-01", Category: "docs", Status: model.StatusFull, Points: 4, MaxPoints: 4},
		{ID: "DOC-03", Category: "docs", Status: model.StatusNone, Points: 0, MaxPoints: 3},
		{ID: "TST-01", Category: "tests", Status: model.StatusFull, Points: 8, MaxPoints: 8},
		{ID: "TST-02", Category: "tests", Status: model.StatusNone, Points: 0, MaxPoints: 2},
	}

	r := Score(results, "/tmp/test", nil, 10, time.Now())

	// docs: 4/7, tests: 8/10 → total 12/17 → ~70%
	if r.Score < 60 || r.Score > 80 {
		t.Errorf("Score = %d, expected ~70 for mixed results", r.Score)
	}
}

func TestScoreWithSkipped(t *testing.T) {
	results := []model.CheckResult{
		{ID: "ACT-01", Category: "activity", Status: model.StatusFull, Points: 5, MaxPoints: 5},
		{ID: "ACT-03", Category: "activity", Status: model.StatusSkipped, Points: 0, MaxPoints: 3},
	}

	r := Score(results, "/tmp/test", nil, 10, time.Now())

	// Skipped check should be redistributed: 5/5 with redistribution → 100%
	if r.Score != 100 {
		t.Errorf("Score = %d, want 100 (full with skipped redistributed)", r.Score)
	}
}

func TestScoreAllSkippedCategory(t *testing.T) {
	// When all checks in a category are skipped (e.g., no git),
	// that category should be excluded entirely — not penalize the score.
	results := []model.CheckResult{
		{ID: "DOC-01", Category: "docs", Status: model.StatusFull, Points: 4, MaxPoints: 4},
		{ID: "DOC-03", Category: "docs", Status: model.StatusFull, Points: 3, MaxPoints: 3},
		{ID: "ACT-01", Category: "activity", Status: model.StatusSkipped, Points: 0, MaxPoints: 5},
		{ID: "ACT-03", Category: "activity", Status: model.StatusSkipped, Points: 0, MaxPoints: 3},
	}

	r := Score(results, "/tmp/test", nil, 10, time.Now())

	if r.Score != 100 {
		t.Errorf("Score = %d, want 100 (fully-skipped category should be excluded)", r.Score)
	}
}

func TestScoreEmptyResults(t *testing.T) {
	r := Score(nil, "/tmp/test", nil, 0, time.Now())

	if r.Score != 0 {
		t.Errorf("Score = %d, want 0 for empty results", r.Score)
	}
	if r.Grade != "F" {
		t.Errorf("Grade = %q, want F for empty results", r.Grade)
	}
	if len(r.Categories) != 0 {
		t.Errorf("Categories = %v, want empty for empty results", r.Categories)
	}
}

func TestScoreRoundingAtBoundary(t *testing.T) {
	// 7/10 total = exactly 70 → B-
	results := []model.CheckResult{
		{ID: "DOC-01", Category: "docs", Status: model.StatusFull, Points: 7, MaxPoints: 10},
	}

	r := Score(results, "/tmp/test", nil, 10, time.Now())

	if r.Score != 70 {
		t.Errorf("Score = %d, want 70 (exact boundary)", r.Score)
	}
	if r.Grade != "B-" {
		t.Errorf("Grade = %q, want B- for score 70", r.Grade)
	}
}

func TestScorePartialChecks(t *testing.T) {
	results := []model.CheckResult{
		// docs: full + partial = 4+1 = 5 out of 4+2 = 6
		{ID: "DOC-01", Category: "docs", Status: model.StatusFull, Points: 4, MaxPoints: 4},
		{ID: "DOC-02", Category: "docs", Status: model.StatusPartial, Points: 1, MaxPoints: 2},
		// tests: none only = 0 out of 8
		{ID: "TST-01", Category: "tests", Status: model.StatusNone, Points: 0, MaxPoints: 8},
		// cicd: full = 6 out of 6
		{ID: "CI-01", Category: "cicd", Status: model.StatusFull, Points: 6, MaxPoints: 6},
	}

	r := Score(results, "/tmp/test", nil, 10, time.Now())

	// total: 5+0+6 = 11 out of 6+8+6 = 20 → 55%
	if r.Score < 50 || r.Score > 60 {
		t.Errorf("Score = %d, expected ~55 for partial mix", r.Score)
	}

	// Must have exactly 3 categories
	if len(r.Categories) != 3 {
		t.Errorf("Categories count = %d, want 3", len(r.Categories))
	}

	// docs category score should be 5, max 6
	var docsCat *model.CategoryResult
	for i := range r.Categories {
		if r.Categories[i].Name == "docs" {
			docsCat = &r.Categories[i]
		}
	}
	if docsCat == nil {
		t.Fatal("docs category not found in report")
	}
	if docsCat.Score != 5 {
		t.Errorf("docs Score = %d, want 5", docsCat.Score)
	}
}

func TestRecommendations(t *testing.T) {
	results := []model.CheckResult{
		{ID: "DOC-01", Category: "docs", Status: model.StatusFull, Points: 4, MaxPoints: 4},
		{ID: "DOC-04", Category: "docs", Status: model.StatusNone, Points: 0, MaxPoints: 2, Suggestion: "Add CONTRIBUTING.md"},
		{ID: "TST-01", Category: "tests", Status: model.StatusNone, Points: 0, MaxPoints: 8, Suggestion: "Add test files"},
	}

	suggestions := Recommendations(results)

	if len(suggestions) != 2 {
		t.Fatalf("got %d suggestions, want 2", len(suggestions))
	}

	// Should be sorted by impact (8 before 2)
	if suggestions[0].Impact != 8 {
		t.Errorf("first suggestion impact = %d, want 8", suggestions[0].Impact)
	}
	if suggestions[1].Impact != 2 {
		t.Errorf("second suggestion impact = %d, want 2", suggestions[1].Impact)
	}
}

func TestScoreCategoryNeverExceedsMax(t *testing.T) {
	// Edge case: after redistribution, catScore should never exceed effectiveMax
	results := []model.CheckResult{
		{ID: "DOC-01", Category: "docs", Status: model.StatusFull, Points: 3, MaxPoints: 3},
		{ID: "DOC-02", Category: "docs", Status: model.StatusSkipped, Points: 0, MaxPoints: 1},
	}
	r := Score(results, "/tmp/test", nil, 10, time.Now())
	for _, cat := range r.Categories {
		if cat.Score > cat.MaxScore {
			t.Errorf("Category %s: score %d exceeds max %d", cat.Name, cat.Score, cat.MaxScore)
		}
	}
}

func TestScoreFullMarks(t *testing.T) {
	results := []model.CheckResult{
		{ID: "DOC-01", Category: "docs", Status: model.StatusFull, Points: 15, MaxPoints: 15},
		{ID: "TST-01", Category: "tests", Status: model.StatusFull, Points: 20, MaxPoints: 20},
	}
	r := Score(results, "/tmp/test", nil, 10, time.Now())
	if r.Score != 100 {
		t.Errorf("Score = %d, want 100 for all full marks", r.Score)
	}
	if r.Grade != "A+" {
		t.Errorf("Grade = %q, want A+", r.Grade)
	}
}

func TestScoreRawMaxSet(t *testing.T) {
	results := []model.CheckResult{
		{ID: "DOC-01", Category: "docs", Status: model.StatusFull, Points: 4, MaxPoints: 4},
		{ID: "TST-01", Category: "tests", Status: model.StatusNone, Points: 0, MaxPoints: 8},
	}
	r := Score(results, "/tmp/test", nil, 10, time.Now())
	if r.RawMax != 12 {
		t.Errorf("RawMax = %d, want 12", r.RawMax)
	}
}

func TestScoreSchemaVersion(t *testing.T) {
	r := Score(nil, "/tmp/test", nil, 0, time.Now())
	if r.SchemaVersion != "1.0" {
		t.Errorf("SchemaVersion = %q, want 1.0", r.SchemaVersion)
	}
}

func TestGradeBoundaries(t *testing.T) {
	tests := []struct {
		score int
		want  string
	}{
		{100, "A+"}, {95, "A+"}, {94, "A"}, {90, "A"},
		{89, "A-"}, {85, "A-"}, {84, "B+"}, {80, "B+"},
		{79, "B"}, {75, "B"}, {74, "B-"}, {70, "B-"},
		{69, "C+"}, {65, "C+"}, {64, "C"}, {60, "C"},
		{59, "C-"}, {55, "C-"}, {54, "D"}, {40, "D"},
		{39, "F"}, {0, "F"},
	}
	for _, tt := range tests {
		got := Grade(tt.score)
		if got != tt.want {
			t.Errorf("Grade(%d) = %q, want %q", tt.score, got, tt.want)
		}
	}
}
