package checks

import (
	"os"
	"path/filepath"
	"strings"
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

func TestCommentBody(t *testing.T) {
	cases := []struct {
		line    string
		openers []string
		want    string
	}{
		{"x := 1 // TODO here", cStyleOpeners, " TODO HERE"},
		{"/* HACK */", cStyleOpeners, " HACK */"},
		{" * TODO continuation", cStyleOpeners, " TODO CONTINUATION"},
		{"u := \"https://x.io/todo\"", cStyleOpeners, ""},
		{"u := \"https://x.io\" // TODO", cStyleOpeners, " TODO"},
		{"strings.Contains(s, \"TODO\")", cStyleOpeners, ""},
		{"msg := \"// TODO not really\"", cStyleOpeners, ""},
		{"r := 'a' // TODO after rune", cStyleOpeners, " TODO AFTER RUNE"},
		{"s := \"a\\\"b\" // TODO after escape", cStyleOpeners, " TODO AFTER ESCAPE"},
		{"#include \"todo.h\"", cStyleOpeners, ""},
		{"# FIXME", hashOpeners, " FIXME"},
		{"s = \"# FIXME inside\"", hashOpeners, ""},
		{"s = 'it\\'s' # TODO", hashOpeners, " TODO"},
		{"SELECT 1 -- XXX", []string{"--", "/*"}, " XXX"},
		{"'don''t' -- TODO", []string{"--", "/*"}, " TODO"},
		{"<!-- TODO -->", []string{"<!--"}, " TODO -->"},
		{"#todo { color: red }", []string{"/*"}, ""},
		{" * TODO", hashOpeners, ""},
	}
	for _, c := range cases {
		if got := commentBody(strings.ToUpper(c.line), c.openers); got != c.want {
			t.Errorf("commentBody(%q) = %q, want %q", c.line, got, c.want)
		}
	}
}

func TestTodoScan_CountsCommentStyles(t *testing.T) {
	dir := makeTempRepo(t, map[string]string{
		"a.go":   "x := 1 // TODO: tune\n",
		"b.py":   "# FIXME broken\n",
		"c.sql":  "-- HACK: temp index\n",
		"d.html": "<!-- XXX remove before launch -->\n",
		"e.go":   "/*\n * TODO: multi-line block\n */\n",
		"f.css":  "a { color: red } /* TODO contrast */\n",
	})
	todoCache.key = ""
	ctx := &model.ScanContext{RepoPath: dir, Files: fileInfos("a.go", "b.py", "c.sql", "d.html", "e.go", "f.css")}
	result := (&TodoCountCheck{}).Run(ctx)
	if result.Details != "6 TODO/FIXME markers found" {
		t.Errorf("Details = %q, want 6 markers (one per comment style)", result.Details)
	}
}

func TestTodoScan_MarkersOutsideComments(t *testing.T) {
	dir := makeTempRepo(t, map[string]string{
		"scan.go": "return strings.Contains(s, \"TODO\") || strings.Contains(s, \"FIXME\")\n",
		"url.go":  "var url = \"https://example.com/todo\"\n",
		"name.go": "var hackathon = 1 // Hackathon signup\n",
		"msg.go":  "msg := \"// TODO: not really\"\n",
		"re.go":   "re := regexp.MustCompile(\"//.*TODO\")\n",
		"s.py":    "s = \"# FIXME inside string\"\n",
		"list.go": "x := y // see TODO_LIST\n",
	})
	todoCache.key = ""
	ctx := &model.ScanContext{RepoPath: dir, Files: fileInfos("scan.go", "url.go", "name.go", "msg.go", "re.go", "s.py", "list.go")}
	result := (&TodoCountCheck{}).Run(ctx)
	if result.Status != model.StatusFull {
		t.Errorf("Status = %v, want Full — none of these lines carries a TODO comment; details: %s", result.Status, result.Details)
	}
}

func TestTodoScan_ForeignOpenersAreNotComments(t *testing.T) {
	dir := makeTempRepo(t, map[string]string{
		"style.css": "#todo { color: red }\n",
		"app.js":    "document.querySelector(\"#todo\");\n",
		"main.c":    "#include \"todo.h\"\n",
		"count.go":  "n--; s := 1 -- TODO --\n",
		"page.html": "<a href=\"#todo\">x</a>\n",
	})
	todoCache.key = ""
	ctx := &model.ScanContext{RepoPath: dir, Files: fileInfos("style.css", "app.js", "main.c", "count.go", "page.html")}
	result := (&TodoCountCheck{}).Run(ctx)
	if result.Status != model.StatusFull {
		t.Errorf("Status = %v, want Full — # and -- are not comment openers in these languages; details: %s", result.Status, result.Details)
	}
}

// Pins the documented limitations so any change to them is deliberate and
// docs/checks.md stays accurate.
func TestTodoScan_DocumentedLimitations(t *testing.T) {
	dir := makeTempRepo(t, map[string]string{
		"tmpl.go": "tmpl := `\n// TODO inside a raw string is still counted\n`\n",
		"life.rs": "fn f<'a>() {} // TODO after a single lifetime tick is missed\n",
	})
	todoCache.key = ""
	ctx := &model.ScanContext{RepoPath: dir, Files: fileInfos("tmpl.go", "life.rs")}
	result := (&TodoCountCheck{}).Run(ctx)
	if result.Details != "1 TODO/FIXME markers found" {
		t.Errorf("Details = %q, want exactly 1 (raw-string line counted, lifetime line missed)", result.Details)
	}
}

func TestTodoCriticalCheck_IgnoresCriticalWordsOutsideComments(t *testing.T) {
	dir := makeTempRepo(t, map[string]string{
		"main.go": "const marker = \"SECURITY\" // TODO: rename\n",
	})
	todoCache.key = ""
	ctx := &model.ScanContext{RepoPath: dir, Files: fileInfos("main.go")}
	result := (&TodoCriticalCheck{}).Run(ctx)
	if result.Status != model.StatusFull {
		t.Errorf("Status = %v, want Full — SECURITY is in a string literal, not in the comment", result.Status)
	}
}
