package checks

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

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

func TestBusFactorCheck_Evaluate(t *testing.T) {
	tests := []struct {
		commits, factor int
		wantStatus      model.Status
		wantPoints      int
		wantDetails     string
	}{
		{0, 0, model.StatusSkipped, 0, "No commits found"},
		{11, 0, model.StatusFull, 2, "no single author exceeds 10% of non-merge commits"},
		{10, 0, model.StatusFull, 2, "no single author exceeds 10% of non-merge commits"},
		{5, 1, model.StatusNone, 0, "bus factor 1"},
		{40, 2, model.StatusPartial, 1, "bus factor 2"},
		{30, 3, model.StatusFull, 2, "bus factor 3"},
		{100, 5, model.StatusFull, 2, "bus factor 5"},
	}

	check := &BusFactorCheck{}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("commits=%d,factor=%d", tt.commits, tt.factor), func(t *testing.T) {
			result := check.evaluate(tt.commits, tt.factor)
			if result.Status != tt.wantStatus {
				t.Errorf("Status = %v, want %v", result.Status, tt.wantStatus)
			}
			if result.Points != tt.wantPoints {
				t.Errorf("Points = %d, want %d", result.Points, tt.wantPoints)
			}
			if result.Details != tt.wantDetails {
				t.Errorf("Details = %q, want %q", result.Details, tt.wantDetails)
			}
			if result.MaxPoints != 2 {
				t.Errorf("MaxPoints = %d, want 2", result.MaxPoints)
			}
		})
	}
}

// gitRepoWithAuthors creates a fresh git repository in a temp directory and
// records commitsPerAuthor[i] empty commits authored by "Author i".
func gitRepoWithAuthors(t *testing.T, commitsPerAuthor ...int) string {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	dir := t.TempDir()
	baseEnv := append(os.Environ(), "GIT_TERMINAL_PROMPT=0")

	run := func(env []string, args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		cmd.Env = env
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			t.Fatalf("git %v: %v: %s", args, err, stderr.String())
		}
	}

	run(baseEnv, "init", "-q")

	for i, n := range commitsPerAuthor {
		name := fmt.Sprintf("Author %d", i)
		email := fmt.Sprintf("author%d@example.com", i)
		env := append(append([]string{}, baseEnv...),
			"GIT_AUTHOR_NAME="+name,
			"GIT_AUTHOR_EMAIL="+email,
			"GIT_COMMITTER_NAME="+name,
			"GIT_COMMITTER_EMAIL="+email,
		)
		for j := 0; j < n; j++ {
			run(env,
				"-c", "user.name="+name,
				"-c", "user.email="+email,
				"-c", "commit.gpgsign=false",
				"commit", "-q", "--allow-empty", "-m", "commit",
			)
		}
	}

	return dir
}

// git shortlog exits non-zero on a repository with no commits (unborn
// HEAD), so an empty repository surfaces as the failure path rather than
// as an empty result.
func TestBusFactorCheck_EmptyRepositoryIsSkipped(t *testing.T) {
	dir := gitRepoWithAuthors(t)
	ctx := &model.ScanContext{GitAvailable: true, RepoPath: dir}
	result := (&BusFactorCheck{}).Run(ctx)
	if result.Status != model.StatusSkipped {
		t.Errorf("Status = %v, want Skipped", result.Status)
	}
	if result.Points != 0 {
		t.Errorf("Points = %d, want 0", result.Points)
	}
	if result.Status == model.StatusFull {
		t.Error("Status = Full, want not Full")
	}
	if result.Details != "Could not compute bus factor" {
		t.Errorf("Details = %q, want %q", result.Details, "Could not compute bus factor")
	}

	plainDir := t.TempDir()
	ctx2 := &model.ScanContext{GitAvailable: true, RepoPath: plainDir}
	result2 := (&BusFactorCheck{}).Run(ctx2)
	if result2.Status != model.StatusSkipped {
		t.Errorf("Status = %v, want Skipped", result2.Status)
	}
	if result2.Details != "Could not compute bus factor" {
		t.Errorf("Details = %q, want %q", result2.Details, "Could not compute bus factor")
	}
}

func TestBusFactorCheck_GitTimeoutIsSkipped(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping slow git-timeout test in short mode")
	}

	t.Run("timeout", func(t *testing.T) {
		fakeDir := t.TempDir()
		if err := os.WriteFile(fakeDir+"/git", []byte("#!/bin/sh\nexec /bin/sleep 30\n"), 0755); err != nil {
			t.Fatalf("write fake git: %v", err)
		}
		t.Setenv("PATH", fakeDir)

		ctx := &model.ScanContext{GitAvailable: true, RepoPath: t.TempDir()}
		start := time.Now()
		result := (&BusFactorCheck{}).Run(ctx)
		elapsed := time.Since(start)
		if elapsed < 5*time.Second {
			t.Errorf("git returned after %v, want it killed by the scanner's 5 s timeout", elapsed)
		}
		if result.Status != model.StatusSkipped {
			t.Errorf("Status = %v, want Skipped", result.Status)
		}
		if result.Points != 0 {
			t.Errorf("Points = %d, want 0", result.Points)
		}
		if result.Details != "Could not compute bus factor" {
			t.Errorf("Details = %q, want %q", result.Details, "Could not compute bus factor")
		}
	})

	t.Run("nonzero exit", func(t *testing.T) {
		fakeDir := t.TempDir()
		if err := os.WriteFile(fakeDir+"/git", []byte("#!/bin/sh\nexit 128\n"), 0755); err != nil {
			t.Fatalf("write fake git: %v", err)
		}
		t.Setenv("PATH", fakeDir)

		ctx := &model.ScanContext{GitAvailable: true, RepoPath: t.TempDir()}
		result := (&BusFactorCheck{}).Run(ctx)
		if result.Status != model.StatusSkipped {
			t.Errorf("Status = %v, want Skipped", result.Status)
		}
		if result.Points != 0 {
			t.Errorf("Points = %d, want 0", result.Points)
		}
		if result.Details != "Could not compute bus factor" {
			t.Errorf("Details = %q, want %q", result.Details, "Could not compute bus factor")
		}
	})
}

func TestBusFactorCheck_DistributedAuthorsIsFull(t *testing.T) {
	counts := make([]int, 11)
	for i := range counts {
		counts[i] = 1
	}
	dir := gitRepoWithAuthors(t, counts...)
	ctx := &model.ScanContext{GitAvailable: true, RepoPath: dir}
	result := (&BusFactorCheck{}).Run(ctx)
	if result.Status != model.StatusFull {
		t.Errorf("Status = %v, want Full", result.Status)
	}
	if result.Points != 2 {
		t.Errorf("Points = %d, want 2", result.Points)
	}
	if result.Details != "no single author exceeds 10% of non-merge commits" {
		t.Errorf("Details = %q, want %q", result.Details, "no single author exceeds 10% of non-merge commits")
	}
}

func TestBusFactorCheck_SingleAuthorIsNone(t *testing.T) {
	dir := gitRepoWithAuthors(t, 3)
	ctx := &model.ScanContext{GitAvailable: true, RepoPath: dir}
	result := (&BusFactorCheck{}).Run(ctx)
	if result.Status != model.StatusNone {
		t.Errorf("Status = %v, want None", result.Status)
	}
	if result.Points != 0 {
		t.Errorf("Points = %d, want 0", result.Points)
	}
	if result.Details != "bus factor 1" {
		t.Errorf("Details = %q, want %q", result.Details, "bus factor 1")
	}
}

func TestBusFactorCheck_TwoAuthorsIsPartial(t *testing.T) {
	dir := gitRepoWithAuthors(t, 2, 2)
	ctx := &model.ScanContext{GitAvailable: true, RepoPath: dir}
	result := (&BusFactorCheck{}).Run(ctx)
	if result.Status != model.StatusPartial {
		t.Errorf("Status = %v, want Partial", result.Status)
	}
	if result.Points != 1 {
		t.Errorf("Points = %d, want 1", result.Points)
	}
	if result.Details != "bus factor 2" {
		t.Errorf("Details = %q, want %q", result.Details, "bus factor 2")
	}
}

func TestBusFactorCheck_ThreeAuthorsIsFull(t *testing.T) {
	dir := gitRepoWithAuthors(t, 2, 2, 2)
	ctx := &model.ScanContext{GitAvailable: true, RepoPath: dir}
	result := (&BusFactorCheck{}).Run(ctx)
	if result.Status != model.StatusFull {
		t.Errorf("Status = %v, want Full", result.Status)
	}
	if result.Points != 2 {
		t.Errorf("Points = %d, want 2", result.Points)
	}
	if result.Details != "bus factor 3" {
		t.Errorf("Details = %q, want %q", result.Details, "bus factor 3")
	}
}

func TestContributorCountCheck_UnaffectedByBusFactorChange(t *testing.T) {
	counts := make([]int, 11)
	for i := range counts {
		counts[i] = 1
	}
	dir := gitRepoWithAuthors(t, counts...)
	ctx := &model.ScanContext{GitAvailable: true, RepoPath: dir}

	// Run BusFactorCheck first, then ContributorCountCheck, to show the
	// shared shortlog cache still serves both checks correctly.
	busResult := (&BusFactorCheck{}).Run(ctx)
	if busResult.Status != model.StatusFull {
		t.Fatalf("BusFactorCheck Status = %v, want Full", busResult.Status)
	}

	contribResult := (&ContributorCountCheck{}).Run(ctx)
	if contribResult.Status != model.StatusFull {
		t.Errorf("Status = %v, want Full", contribResult.Status)
	}
	if contribResult.Details != "11 contributors" {
		t.Errorf("Details = %q, want %q", contribResult.Details, "11 contributors")
	}

	dir2 := gitRepoWithAuthors(t, 3, 3)
	ctx2 := &model.ScanContext{GitAvailable: true, RepoPath: dir2}
	contribResult2 := (&ContributorCountCheck{}).Run(ctx2)
	if contribResult2.Status != model.StatusPartial {
		t.Errorf("Status = %v, want Partial", contribResult2.Status)
	}
	if contribResult2.Details != "2 contributors" {
		t.Errorf("Details = %q, want %q", contribResult2.Details, "2 contributors")
	}
}
