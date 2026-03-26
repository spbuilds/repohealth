# Check Reference

RepoHealth v0.1 includes 14 checks across 4 categories. Each check has an ID, a detection method, and a point value.

## Documentation (15 points)

| ID | Check | Detection | Points |
|----|-------|-----------|--------|
| DOC-01 | README exists | `README.md`, `README`, `README.rst`, `readme.md` | 4 |
| DOC-02 | README has content | README file size > 100 bytes | 2 |
| DOC-03 | LICENSE exists | `LICENSE`, `LICENSE.md`, `LICENCE`, `COPYING`, `LICENSE.txt` | 3 |
| DOC-04 | CONTRIBUTING exists | `CONTRIBUTING.md`, `.github/CONTRIBUTING.md` | 2 |
| DOC-05 | CODE_OF_CONDUCT exists | `CODE_OF_CONDUCT.md`, `.github/CODE_OF_CONDUCT.md` | 1 |
| DOC-06 | SECURITY.md exists | `SECURITY.md`, `.github/SECURITY.md` | 2 |
| DOC-07 | CHANGELOG exists | `CHANGELOG.md`, `CHANGELOG`, `HISTORY.md`, `CHANGES.md` | 1 |

All documentation checks are file-existence based. DOC-02 additionally checks file size to ensure the README has substantive content.

## Testing (14 points)

| ID | Check | Detection | Points |
|----|-------|-----------|--------|
| TST-01 | Test files detected | Files matching: `*_test.go`, `*_test.py`, `test_*.py`, `*.test.ts`, `*.test.js`, `*.spec.ts`, `*.spec.js`, `*_test.rs`, `*_test.rb`, `*Test.java`, `*Tests.java` | 8 |
| TST-02 | Test directory exists | Directories named: `test/`, `tests/`, `__tests__/`, `spec/` | 2 |
| TST-03 | Test framework configured | Config files: `jest.config.*`, `vitest.config.*`, `pytest.ini`, `pyproject.toml`, `.mocharc.*`, `phpunit.xml`; or `*_test.go` files (Go built-in) | 4 |

TST-01 uses glob-style pattern matching on filenames. It detects test files for Go, Python, JavaScript, TypeScript, Rust, Ruby, and Java.

## CI/CD (6 points)

| ID | Check | Detection | Points |
|----|-------|-----------|--------|
| CI-01 | CI configuration exists | Directories: `.github/workflows/`, `.circleci/`, `.buildkite/`; Files: `.gitlab-ci.yml`, `Jenkinsfile`, `.travis.yml`, `bitbucket-pipelines.yml`, `azure-pipelines.yml`, `Taskfile.yml` | 6 |

CI-01 detects the presence of any major CI/CD system. It does not parse the CI configuration content (that's planned for v0.2 checks CI-02 through CI-04).

## Activity (8 points)

| ID | Check | Detection | Points |
|----|-------|-----------|--------|
| ACT-01 | Recent commit | `git log -1 --format=%cI` — within 30 days = Full (5), within 90 days = Partial (2), older = None (0) | 5 |
| ACT-03 | Contributors | `git shortlog -sn --no-merges` — more than 5 = Full (3), 2-5 = Partial (1), 1 = None (0), 0 = Skipped | 3 |

Activity checks require a git repository. If `.git/` is not present, both checks return Skipped and the category is excluded from scoring.

Git commands use a 5-second timeout to prevent hangs on corrupted or network-mounted repositories.

## Adding New Checks

To add a check, implement the `Check` interface in `internal/checks/`:

```go
type Check interface {
    ID() string           // e.g., "DOC-01"
    Category() string     // e.g., "docs"
    Name() string         // e.g., "README exists"
    MaxPoints() int       // e.g., 4
    Run(ctx *model.ScanContext) model.CheckResult
}
```

Then register it in `NewRegistry()` in `internal/checks/check.go`. See [CONTRIBUTING.md](../CONTRIBUTING.md) for the full workflow.

## Planned Checks (v0.2)

| Category | Checks | Description |
|----------|--------|-------------|
| Testing | TST-04, TST-05 | Coverage config, test-to-source ratio |
| CI/CD | CI-02 to CI-04 | CI runs tests, linter, build |
| Dependencies | DEP-01 to DEP-05 | Lockfile, manager, freshness, outdated, count |
| Security | SEC-02 to SEC-05 | Secrets, gitignore, pinning, branch protection |
| Code Statistics | STAT-01 to STAT-04 | File count, languages, comment ratio, vendor bloat |
| Activity | ACT-02, ACT-04, ACT-05 | Commit frequency, releases, bus factor |
| TODO/Debt | TODO-01 to TODO-03 | Count, density, critical markers |
