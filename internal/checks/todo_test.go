package checks

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spbuilds/repohealth/internal/model"
)

func makeTempRepo(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	for name, content := range files {
		full := filepath.Join(dir, name)
		if err := os.MkdirAll(filepath.Dir(full), 0755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}
		if err := os.WriteFile(full, []byte(content), 0644); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
	}
	return dir
}

func fileInfos(paths ...string) []model.FileInfo {
	infos := make([]model.FileInfo, 0, len(paths))
	for _, p := range paths {
		infos = append(infos, model.FileInfo{Path: p, Name: filepath.Base(p)})
	}
	return infos
}

func TestTodoCountCheck_Full(t *testing.T) {
	dir := makeTempRepo(t, map[string]string{
		"main.go": "package main\n\nfunc main() {}\n",
	})
	// Reset cache so this temp dir is scanned fresh
	todoCache.key = ""

	ctx := &model.ScanContext{
		RepoPath: dir,
		Files:    fileInfos("main.go"),
	}

	check := &TodoCountCheck{}
	result := check.Run(ctx)

	if result.Status != model.StatusFull {
		t.Errorf("Status = %v, want Full (no TODOs)", result.Status)
	}
	if result.Points != check.MaxPoints() {
		t.Errorf("Points = %d, want %d", result.Points, check.MaxPoints())
	}
}

func TestTodoCountCheck_Partial(t *testing.T) {
	dir := makeTempRepo(t, map[string]string{
		"main.go": "// TODO: fix this\n// FIXME: broken\n// HACK: workaround\n",
	})
	todoCache.key = ""

	ctx := &model.ScanContext{
		RepoPath: dir,
		Files:    fileInfos("main.go"),
	}

	check := &TodoCountCheck{}
	result := check.Run(ctx)

	if result.Status != model.StatusPartial {
		t.Errorf("Status = %v, want Partial (1-20 TODOs)", result.Status)
	}
}

func TestTodoCountCheck_None_Empty(t *testing.T) {
	todoCache.key = ""

	ctx := &model.ScanContext{
		RepoPath: t.TempDir(),
		Files:    []model.FileInfo{},
	}

	check := &TodoCountCheck{}
	result := check.Run(ctx)

	if result.Status != model.StatusFull {
		t.Errorf("Status = %v, want Full (no source files = 0 TODOs)", result.Status)
	}
}

func TestTodoDensityCheck_Full(t *testing.T) {
	// 1000 lines, 1 TODO → density < 2/KLOC
	content := ""
	for i := 0; i < 999; i++ {
		content += "x := 1\n"
	}
	content += "// TODO: minor\n"

	dir := makeTempRepo(t, map[string]string{"app.go": content})
	todoCache.key = ""

	ctx := &model.ScanContext{
		RepoPath: dir,
		Files:    fileInfos("app.go"),
	}

	check := &TodoDensityCheck{}
	result := check.Run(ctx)

	if result.Status != model.StatusFull {
		t.Errorf("Status = %v, want Full (density < 2/KLOC)", result.Status)
	}
}

func TestTodoDensityCheck_Skipped_NoFiles(t *testing.T) {
	todoCache.key = ""

	ctx := &model.ScanContext{
		RepoPath: t.TempDir(),
		Files:    []model.FileInfo{},
	}

	check := &TodoDensityCheck{}
	result := check.Run(ctx)

	if result.Status != model.StatusSkipped {
		t.Errorf("Status = %v, want Skipped (no source lines)", result.Status)
	}
}

func TestTodoCriticalCheck_Full(t *testing.T) {
	dir := makeTempRepo(t, map[string]string{
		"main.go": "// TODO: improve performance\n",
	})
	todoCache.key = ""

	ctx := &model.ScanContext{
		RepoPath: dir,
		Files:    fileInfos("main.go"),
	}

	check := &TodoCriticalCheck{}
	result := check.Run(ctx)

	if result.Status != model.StatusFull {
		t.Errorf("Status = %v, want Full (no critical markers)", result.Status)
	}
	if result.Points != check.MaxPoints() {
		t.Errorf("Points = %d, want %d", result.Points, check.MaxPoints())
	}
}

func TestTodoCriticalCheck_None(t *testing.T) {
	dir := makeTempRepo(t, map[string]string{
		"auth.go": "// FIXME: SECURITY - validate input\n",
	})
	todoCache.key = ""

	ctx := &model.ScanContext{
		RepoPath: dir,
		Files:    fileInfos("auth.go"),
	}

	check := &TodoCriticalCheck{}
	result := check.Run(ctx)

	if result.Status != model.StatusNone {
		t.Errorf("Status = %v, want None (critical marker present)", result.Status)
	}
	if result.Points != 0 {
		t.Errorf("Points = %d, want 0", result.Points)
	}
}
