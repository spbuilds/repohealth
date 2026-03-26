package report

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spbuilds/repohealth/internal/model"
)

func sampleReport() *model.Report {
	return &model.Report{
		SchemaVersion: "1.0",
		Version:       "test",
		RepoPath:      "/tmp/test",
		Timestamp:     "2026-03-26T12:00:00Z",
		Score:         75,
		MaxScore:      100,
		Grade:         "B",
		Categories: []model.CategoryResult{
			{Name: "docs", Label: "Documentation", Score: 15, MaxScore: 15},
		},
		Checks: []model.CheckResult{
			{ID: "DOC-01", Category: "docs", Name: "README exists", Status: model.StatusFull, Points: 4, MaxPoints: 4, Details: "found"},
		},
		Suggestions: []model.Suggestion{
			{CheckID: "TST-01", Impact: 3, Message: "Add tests"},
		},
	}
}

func TestMarkdownOutput(t *testing.T) {
	var buf bytes.Buffer
	Markdown(&buf, sampleReport())
	out := buf.String()
	if !strings.Contains(out, "# Repository Health Report") {
		t.Error("missing header")
	}
	if !strings.Contains(out, "75 / 100") {
		t.Error("missing score")
	}
	if !strings.Contains(out, "Documentation") {
		t.Error("missing category")
	}
}

func TestHTMLOutput(t *testing.T) {
	var buf bytes.Buffer
	HTML(&buf, sampleReport())
	out := buf.String()
	if !strings.Contains(out, "<!DOCTYPE html>") {
		t.Error("missing doctype")
	}
	if !strings.Contains(out, "75 / 100") {
		t.Error("missing score")
	}
}

func TestJSONOutput(t *testing.T) {
	var buf bytes.Buffer
	err := JSON(&buf, sampleReport())
	if err != nil {
		t.Fatalf("JSON error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"schema_version"`) {
		t.Error("missing schema_version")
	}
	if !strings.Contains(out, `"score": 75`) {
		t.Error("missing score")
	}
}

func TestEmptyReport(t *testing.T) {
	r := &model.Report{Score: 0, Grade: "F", Timestamp: "2026-01-01T00:00:00Z"}
	var buf bytes.Buffer
	Markdown(&buf, r)
	if buf.Len() == 0 {
		t.Error("empty markdown output")
	}
	buf.Reset()
	HTML(&buf, r)
	if buf.Len() == 0 {
		t.Error("empty html output")
	}
}
