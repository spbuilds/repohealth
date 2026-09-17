package checks

import (
	"reflect"
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

func TestTestFrameworkCheck_PriorityOrder(t *testing.T) {
	seen := map[string]bool{}
	var order []string
	for _, fw := range testFrameworkConfigs {
		if !seen[fw.name] {
			seen[fw.name] = true
			order = append(order, fw.name)
		}
	}
	want := []string{"Jest", "Vitest", "pytest", "Mocha", "PHPUnit"}
	if !reflect.DeepEqual(order, want) {
		t.Errorf("framework priority order = %v, want %v", order, want)
	}
}

func TestTestFrameworkCheck_SeveralFrameworksPresent(t *testing.T) {
	tests := []struct {
		name  string
		files []string
		want  string
	}{
		{
			name:  "reverse priority order",
			files: []string{".mocharc.yml", "vitest.config.ts", "jest.config.js"},
			want:  "Jest",
		},
		{
			name:  "vitest and mocha",
			files: []string{".mocharc.yml", "vitest.config.ts"},
			want:  "Vitest",
		},
		{
			name:  "phpunit mocha and pytest",
			files: []string{"phpunit.xml", ".mocharc.yml", "pytest.ini"},
			want:  "pytest",
		},
		{
			name:  "phpunit dist and mocha",
			files: []string{"phpunit.xml.dist", ".mocharc.js"},
			want:  "Mocha",
		},
		{
			name:  "subdirectory config does not outrank root priority",
			files: []string{"packages/app/vitest.config.ts", "jest.config.js"},
			want:  "Jest",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &model.ScanContext{Files: fileInfos(tt.files...)}
			result := (&TestFrameworkCheck{}).Run(ctx)
			if result.Details != tt.want {
				t.Errorf("Details = %q, want %q", result.Details, tt.want)
			}
			if result.Status != model.StatusFull {
				t.Errorf("Status = %v, want Full", result.Status)
			}
			if result.Points != 4 {
				t.Errorf("Points = %d, want 4", result.Points)
			}
		})
	}
}

func TestTestFrameworkCheck_RepeatedRunsIdentical(t *testing.T) {
	ctx := &model.ScanContext{
		Files: fileInfos("jest.config.js", "vitest.config.ts", ".mocharc.yml", "src/a_test.go"),
	}
	check := &TestFrameworkCheck{}
	first := check.Run(ctx)
	if first.Details != "Jest" {
		t.Fatalf("Details = %q, want Jest", first.Details)
	}
	for i := 0; i < 50; i++ {
		result := check.Run(ctx)
		if !reflect.DeepEqual(result, first) {
			t.Fatalf("run %d = %+v, want %+v", i, result, first)
		}
	}
}

func TestTestFrameworkCheck_JestOnly(t *testing.T) {
	ctx := &model.ScanContext{Files: fileInfos("jest.config.js")}
	result := (&TestFrameworkCheck{}).Run(ctx)
	if result.Details != "Jest" {
		t.Errorf("Details = %q, want Jest", result.Details)
	}
	if result.Status != model.StatusFull || result.Points != 4 || result.MaxPoints != 4 {
		t.Errorf("result = %+v, want Full 4/4", result)
	}
}

func TestTestFrameworkCheck_VitestOnly(t *testing.T) {
	ctx := &model.ScanContext{Files: fileInfos("vitest.config.ts")}
	result := (&TestFrameworkCheck{}).Run(ctx)
	if result.Details != "Vitest" {
		t.Errorf("Details = %q, want Vitest", result.Details)
	}
}

func TestTestFrameworkCheck_PytestIni(t *testing.T) {
	ctx := &model.ScanContext{Files: fileInfos("pytest.ini")}
	result := (&TestFrameworkCheck{}).Run(ctx)
	if result.Details != "pytest" {
		t.Errorf("Details = %q, want pytest", result.Details)
	}
}

func TestTestFrameworkCheck_GoBuiltIn(t *testing.T) {
	ctx := &model.ScanContext{Files: fileInfos("pkg/x_test.go", "pkg/x.go")}
	result := (&TestFrameworkCheck{}).Run(ctx)
	if result.Details != "go test (built-in)" {
		t.Errorf("Details = %q, want go test (built-in)", result.Details)
	}
	if result.Status != model.StatusFull {
		t.Errorf("Status = %v, want Full", result.Status)
	}
}

func TestTestFrameworkCheck_NoneFound(t *testing.T) {
	ctx := &model.ScanContext{Files: fileInfos("main.go")}
	result := (&TestFrameworkCheck{}).Run(ctx)
	if result.Status != model.StatusNone {
		t.Errorf("Status = %v, want None", result.Status)
	}
	if result.Points != 0 {
		t.Errorf("Points = %d, want 0", result.Points)
	}
	if result.Details != "No test framework configuration found" {
		t.Errorf("Details = %q, want %q", result.Details, "No test framework configuration found")
	}
}
