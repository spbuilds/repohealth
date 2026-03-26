package checks

import (
	"testing"

	"github.com/spbuilds/repohealth/internal/model"
)

func TestReadmeExistsCheck(t *testing.T) {
	ctx := &model.ScanContext{
		Files: []model.FileInfo{
			{Path: "README.md", Name: "README.md", Size: 500},
		},
	}

	check := &ReadmeExistsCheck{}
	result := check.Run(ctx)

	if result.Status != model.StatusFull {
		t.Errorf("Status = %v, want Full", result.Status)
	}
	if result.Points != 4 {
		t.Errorf("Points = %d, want 4", result.Points)
	}
}

func TestReadmeExistsCheck_Missing(t *testing.T) {
	ctx := &model.ScanContext{
		Files: []model.FileInfo{},
	}

	check := &ReadmeExistsCheck{}
	result := check.Run(ctx)

	if result.Status != model.StatusNone {
		t.Errorf("Status = %v, want None", result.Status)
	}
	if result.Points != 0 {
		t.Errorf("Points = %d, want 0", result.Points)
	}
	if result.Suggestion == "" {
		t.Error("expected a suggestion for missing README")
	}
}

func TestReadmeContentCheck(t *testing.T) {
	ctx := &model.ScanContext{
		Files: []model.FileInfo{
			{Path: "README.md", Name: "README.md", Size: 500},
		},
	}

	check := &ReadmeContentCheck{}
	result := check.Run(ctx)

	if result.Status != model.StatusFull {
		t.Errorf("Status = %v, want Full (size > 100)", result.Status)
	}
}

func TestReadmeContentCheck_Small(t *testing.T) {
	ctx := &model.ScanContext{
		Files: []model.FileInfo{
			{Path: "README.md", Name: "README.md", Size: 50},
		},
	}

	check := &ReadmeContentCheck{}
	result := check.Run(ctx)

	if result.Status != model.StatusPartial {
		t.Errorf("Status = %v, want Partial (size < 100)", result.Status)
	}
}

func TestLicenseExistsCheck(t *testing.T) {
	ctx := &model.ScanContext{
		Files: []model.FileInfo{
			{Path: "LICENSE", Name: "LICENSE", Size: 1000},
		},
	}

	check := &LicenseExistsCheck{}
	result := check.Run(ctx)

	if result.Status != model.StatusFull {
		t.Errorf("Status = %v, want Full", result.Status)
	}
}

func TestAllDocsChecks_FullRepo(t *testing.T) {
	ctx := &model.ScanContext{
		Files: []model.FileInfo{
			{Path: "README.md", Name: "README.md", Size: 500},
			{Path: "LICENSE", Name: "LICENSE", Size: 1000},
			{Path: "CONTRIBUTING.md", Name: "CONTRIBUTING.md", Size: 200},
			{Path: "CODE_OF_CONDUCT.md", Name: "CODE_OF_CONDUCT.md", Size: 300},
			{Path: "SECURITY.md", Name: "SECURITY.md", Size: 150},
			{Path: "CHANGELOG.md", Name: "CHANGELOG.md", Size: 100},
		},
	}

	allChecks := []Check{
		&ReadmeExistsCheck{},
		&ReadmeContentCheck{},
		&LicenseExistsCheck{},
		&ContributingExistsCheck{},
		&CodeOfConductExistsCheck{},
		&SecurityExistsCheck{},
		&ChangelogExistsCheck{},
	}

	totalPoints := 0
	for _, c := range allChecks {
		result := c.Run(ctx)
		totalPoints += result.Points
	}

	if totalPoints != 15 {
		t.Errorf("Total docs points = %d, want 15 (all full)", totalPoints)
	}
}
