package checks

import (
	"testing"

	"github.com/spbuilds/repohealth/internal/model"
)

func TestLockfileExistsCheck_Full(t *testing.T) {
	ctx := &model.ScanContext{
		Files: []model.FileInfo{
			{Path: "go.sum", Name: "go.sum"},
		},
	}

	check := &LockfileExistsCheck{}
	result := check.Run(ctx)

	if result.Status != model.StatusFull {
		t.Errorf("Status = %v, want Full", result.Status)
	}
	if result.Points != check.MaxPoints() {
		t.Errorf("Points = %d, want %d", result.Points, check.MaxPoints())
	}
}

func TestLockfileExistsCheck_None(t *testing.T) {
	ctx := &model.ScanContext{
		Files: []model.FileInfo{},
	}

	check := &LockfileExistsCheck{}
	result := check.Run(ctx)

	if result.Status != model.StatusNone {
		t.Errorf("Status = %v, want None", result.Status)
	}
	if result.Points != 0 {
		t.Errorf("Points = %d, want 0", result.Points)
	}
	if result.Suggestion == "" {
		t.Error("expected a suggestion for missing lockfile")
	}
}

func TestPackageManagerCheck_Full(t *testing.T) {
	ctx := &model.ScanContext{
		Files: []model.FileInfo{
			{Path: "go.mod", Name: "go.mod"},
		},
	}

	check := &PackageManagerCheck{}
	result := check.Run(ctx)

	if result.Status != model.StatusFull {
		t.Errorf("Status = %v, want Full", result.Status)
	}
}

func TestPackageManagerCheck_None(t *testing.T) {
	ctx := &model.ScanContext{
		Files: []model.FileInfo{},
	}

	check := &PackageManagerCheck{}
	result := check.Run(ctx)

	if result.Status != model.StatusNone {
		t.Errorf("Status = %v, want None", result.Status)
	}
}

func TestLockfileFreshnessCheck_Skipped_NoLockfile(t *testing.T) {
	ctx := &model.ScanContext{
		RepoPath: t.TempDir(),
		Files:    []model.FileInfo{},
	}

	check := &LockfileFreshnessCheck{}
	result := check.Run(ctx)

	if result.Status != model.StatusSkipped {
		t.Errorf("Status = %v, want Skipped", result.Status)
	}
}

func TestLockfileFreshnessCheck_Full(t *testing.T) {
	dir := makeTempRepo(t, map[string]string{
		"go.sum": "github.com/example/pkg v1.0.0 h1:abc123\n",
	})

	ctx := &model.ScanContext{
		RepoPath: dir,
		Files: []model.FileInfo{
			{Path: "go.sum", Name: "go.sum"},
		},
	}

	check := &LockfileFreshnessCheck{}
	result := check.Run(ctx)

	// Freshly created file should be within 90 days
	if result.Status != model.StatusFull {
		t.Errorf("Status = %v, want Full (newly created lockfile)", result.Status)
	}
}

func TestDependencyCountCheck_Skipped_NoManifest(t *testing.T) {
	ctx := &model.ScanContext{
		RepoPath: t.TempDir(),
		Files:    []model.FileInfo{},
	}

	check := &DependencyCountCheck{}
	result := check.Run(ctx)

	if result.Status != model.StatusSkipped {
		t.Errorf("Status = %v, want Skipped", result.Status)
	}
}

func TestDependencyCountCheck_Full_GoMod(t *testing.T) {
	gomod := `module example.com/app

go 1.22.0

require (
	github.com/foo/bar v1.0.0
	github.com/baz/qux v2.1.0
)
`
	dir := makeTempRepo(t, map[string]string{"go.mod": gomod})

	ctx := &model.ScanContext{
		RepoPath: dir,
		Files:    []model.FileInfo{{Path: "go.mod", Name: "go.mod"}},
	}

	check := &DependencyCountCheck{}
	result := check.Run(ctx)

	if result.Status != model.StatusFull {
		t.Errorf("Status = %v, want Full (2 deps < 50)", result.Status)
	}
}

func TestDependencyCountCheck_Full_PackageJSON(t *testing.T) {
	pkgjson := `{
  "name": "my-app",
  "dependencies": {
    "react": "^18.0.0",
    "lodash": "^4.17.21"
  },
  "devDependencies": {
    "typescript": "^5.0.0"
  }
}
`
	dir := makeTempRepo(t, map[string]string{"package.json": pkgjson})

	ctx := &model.ScanContext{
		RepoPath: dir,
		Files:    []model.FileInfo{{Path: "package.json", Name: "package.json"}},
	}

	check := &DependencyCountCheck{}
	result := check.Run(ctx)

	if result.Status != model.StatusFull {
		t.Errorf("Status = %v, want Full (3 deps < 50)", result.Status)
	}
}

func TestLockfileExistsCheck_GoSum(t *testing.T) {
	ctx := &model.ScanContext{
		Files: []model.FileInfo{{Path: "go.sum", Name: "go.sum", Size: 1000}},
	}
	check := &LockfileExistsCheck{}
	result := check.Run(ctx)
	if result.Status != model.StatusFull {
		t.Errorf("Status = %v, want Full", result.Status)
	}
}

func TestLockfileExistsCheck_NoFiles(t *testing.T) {
	ctx := &model.ScanContext{}
	check := &LockfileExistsCheck{}
	result := check.Run(ctx)
	if result.Status != model.StatusNone {
		t.Errorf("Status = %v, want None", result.Status)
	}
}

func TestPackageManagerCheck_GoMod(t *testing.T) {
	ctx := &model.ScanContext{
		Files: []model.FileInfo{{Path: "go.mod", Name: "go.mod", Size: 200}},
	}
	check := &PackageManagerCheck{}
	result := check.Run(ctx)
	if result.Status != model.StatusFull {
		t.Errorf("Status = %v, want Full for go.mod", result.Status)
	}
}

func TestPackageManagerCheck_PackageJSON(t *testing.T) {
	ctx := &model.ScanContext{
		Files: []model.FileInfo{{Path: "package.json", Name: "package.json", Size: 500}},
	}
	check := &PackageManagerCheck{}
	result := check.Run(ctx)
	if result.Status != model.StatusFull {
		t.Errorf("Status = %v, want Full for package.json", result.Status)
	}
}
