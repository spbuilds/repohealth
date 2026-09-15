package checks

import (
	"testing"

	"github.com/spbuilds/repohealth/internal/model"
)

func TestTestFilesExistCheck_CSharpAndDart(t *testing.T) {
	ctx := &model.ScanContext{
		Files: fileInfos("src/Login.cs", "tests/LoginTests.cs", "lib/auth.dart", "test/auth_test.dart"),
	}
	result := (&TestFilesExistCheck{}).Run(ctx)
	if result.Details != "2 test files" {
		t.Errorf("Details = %q, want LoginTests.cs and auth_test.dart counted as test files", result.Details)
	}
}

func TestTestToSourceRatioCheck_CSharp(t *testing.T) {
	ctx := &model.ScanContext{
		Files: fileInfos("src/Login.cs", "src/User.cs", "tests/LoginTests.cs"),
	}
	result := (&TestToSourceRatioCheck{}).Run(ctx)
	if result.Details != "1 test files / 2 source files (50%)" {
		t.Errorf("Details = %q, want the test file recognised and excluded from the source count", result.Details)
	}
	if result.Status != model.StatusFull {
		t.Errorf("Status = %v, want Full", result.Status)
	}
}

func TestTestToSourceRatioCheck_IgnoresMarkup(t *testing.T) {
	ctx := &model.ScanContext{
		Files: fileInfos("app.js", "app.test.js", "index.html", "style.css", "theme.scss"),
	}
	result := (&TestToSourceRatioCheck{}).Run(ctx)
	if result.Details != "1 test files / 1 source files (100%)" {
		t.Errorf("Details = %q, want markup and stylesheets excluded from the source count", result.Details)
	}
}
