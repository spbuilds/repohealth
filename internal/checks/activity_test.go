package checks

import (
	"testing"

	"github.com/spbuilds/repohealth/internal/model"
)

func TestRecentCommitCheck_NoGit(t *testing.T) {
	ctx := &model.ScanContext{GitAvailable: false}
	check := &RecentCommitCheck{}
	result := check.Run(ctx)
	if result.Status != model.StatusSkipped {
		t.Errorf("Status = %v, want Skipped when git unavailable", result.Status)
	}
}

func TestContributorCountCheck_NoGit(t *testing.T) {
	ctx := &model.ScanContext{GitAvailable: false}
	check := &ContributorCountCheck{}
	result := check.Run(ctx)
	if result.Status != model.StatusSkipped {
		t.Errorf("Status = %v, want Skipped when git unavailable", result.Status)
	}
}

func TestCommitFrequencyCheck_NoGit(t *testing.T) {
	ctx := &model.ScanContext{GitAvailable: false}
	check := &CommitFrequencyCheck{}
	result := check.Run(ctx)
	if result.Status != model.StatusSkipped {
		t.Errorf("Status = %v, want Skipped when git unavailable", result.Status)
	}
}

func TestReleaseExistsCheck_NoGit(t *testing.T) {
	ctx := &model.ScanContext{GitAvailable: false}
	check := &ReleaseExistsCheck{}
	result := check.Run(ctx)
	if result.Status != model.StatusSkipped {
		t.Errorf("Status = %v, want Skipped when git unavailable", result.Status)
	}
}

func TestBusFactorCheck_NoGit(t *testing.T) {
	ctx := &model.ScanContext{GitAvailable: false}
	check := &BusFactorCheck{}
	result := check.Run(ctx)
	if result.Status != model.StatusSkipped {
		t.Errorf("Status = %v, want Skipped when git unavailable", result.Status)
	}
}
