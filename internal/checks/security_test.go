package checks

import (
	"testing"

	"github.com/spbuilds/repohealth/internal/model"
)

func TestNoSecretsCheck_Full(t *testing.T) {
	dir := makeTempRepo(t, map[string]string{
		"main.go": "package main\n\nfunc main() {}\n",
	})

	ctx := &model.ScanContext{
		RepoPath: dir,
		Files:    fileInfos("main.go"),
	}

	check := &NoSecretsCheck{}
	result := check.Run(ctx)

	if result.Status != model.StatusFull {
		t.Errorf("Status = %v, want Full (no secrets)", result.Status)
	}
	if result.Points != check.MaxPoints() {
		t.Errorf("Points = %d, want %d", result.Points, check.MaxPoints())
	}
}

func TestNoSecretsCheck_None_DotEnv(t *testing.T) {
	ctx := &model.ScanContext{
		RepoPath: t.TempDir(),
		Files: []model.FileInfo{
			{Path: ".env", Name: ".env"},
		},
	}

	check := &NoSecretsCheck{}
	result := check.Run(ctx)

	if result.Status != model.StatusNone {
		t.Errorf("Status = %v, want None (.env present)", result.Status)
	}
	if result.Points != 0 {
		t.Errorf("Points = %d, want 0", result.Points)
	}
}

func TestNoSecretsCheck_None_AWSKey(t *testing.T) {
	dir := makeTempRepo(t, map[string]string{
		"config.go": `package config

const key = "AKIAIOSFODNN7EXAMPLE"
`,
	})

	ctx := &model.ScanContext{
		RepoPath: dir,
		Files:    fileInfos("config.go"),
	}

	check := &NoSecretsCheck{}
	result := check.Run(ctx)

	if result.Status != model.StatusNone {
		t.Errorf("Status = %v, want None (AWS key pattern present)", result.Status)
	}
}

func TestNoSecretsCheck_Full_Empty(t *testing.T) {
	ctx := &model.ScanContext{
		RepoPath: t.TempDir(),
		Files:    []model.FileInfo{},
	}

	check := &NoSecretsCheck{}
	result := check.Run(ctx)

	if result.Status != model.StatusFull {
		t.Errorf("Status = %v, want Full (no files)", result.Status)
	}
}

func TestGitignoreSecretsCheck_Full(t *testing.T) {
	dir := makeTempRepo(t, map[string]string{
		".gitignore": ".env\n*.pem\n*.key\ncredentials\nnode_modules/\n",
	})

	ctx := &model.ScanContext{
		RepoPath: dir,
		Files:    fileInfos(".gitignore"),
	}

	check := &GitignoreSecretsCheck{}
	result := check.Run(ctx)

	if result.Status != model.StatusFull {
		t.Errorf("Status = %v, want Full (all patterns covered)", result.Status)
	}
}

func TestGitignoreSecretsCheck_None_NoFile(t *testing.T) {
	ctx := &model.ScanContext{
		RepoPath: t.TempDir(),
		Files:    []model.FileInfo{},
	}

	check := &GitignoreSecretsCheck{}
	result := check.Run(ctx)

	if result.Status != model.StatusNone {
		t.Errorf("Status = %v, want None (no .gitignore)", result.Status)
	}
}

func TestGitignoreSecretsCheck_Partial(t *testing.T) {
	dir := makeTempRepo(t, map[string]string{
		".gitignore": ".env\nnode_modules/\n",
	})

	ctx := &model.ScanContext{
		RepoPath: dir,
		Files:    fileInfos(".gitignore"),
	}

	check := &GitignoreSecretsCheck{}
	result := check.Run(ctx)

	if result.Status != model.StatusPartial {
		t.Errorf("Status = %v, want Partial (some patterns covered)", result.Status)
	}
}

func TestDependencyPinningCheck_Full(t *testing.T) {
	ctx := &model.ScanContext{
		Files: []model.FileInfo{
			{Path: "package-lock.json", Name: "package-lock.json"},
		},
	}

	check := &DependencyPinningCheck{}
	result := check.Run(ctx)

	if result.Status != model.StatusFull {
		t.Errorf("Status = %v, want Full", result.Status)
	}
}

func TestDependencyPinningCheck_None(t *testing.T) {
	ctx := &model.ScanContext{
		Files: []model.FileInfo{},
	}

	check := &DependencyPinningCheck{}
	result := check.Run(ctx)

	if result.Status != model.StatusNone {
		t.Errorf("Status = %v, want None", result.Status)
	}
}

func TestBranchProtectionCheck_Full(t *testing.T) {
	ctx := &model.ScanContext{
		Files: []model.FileInfo{
			{Path: ".github/CODEOWNERS", Name: "CODEOWNERS"},
		},
	}

	check := &BranchProtectionCheck{}
	result := check.Run(ctx)

	if result.Status != model.StatusFull {
		t.Errorf("Status = %v, want Full", result.Status)
	}
}

func TestBranchProtectionCheck_None(t *testing.T) {
	ctx := &model.ScanContext{
		Files: []model.FileInfo{},
	}

	check := &BranchProtectionCheck{}
	result := check.Run(ctx)

	if result.Status != model.StatusNone {
		t.Errorf("Status = %v, want None", result.Status)
	}
	if result.Suggestion == "" {
		t.Error("expected a suggestion for missing branch protection")
	}
}
