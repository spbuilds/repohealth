# RepoHealth

[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go](https://img.shields.io/badge/go-1.22+-00ADD8.svg)](https://go.dev)
[![CI](https://github.com/spbuilds/repohealth/actions/workflows/ci.yml/badge.svg)](https://github.com/spbuilds/repohealth/actions/workflows/ci.yml)

**A Lighthouse-style health check for Git repositories.**

RepoHealth analyzes a repository's documentation, tests, CI/CD configuration, and maintenance activity to produce a composite health score (0-100) with actionable improvement recommendations. One command, one report, one score.

![RepoHealth Demo](demo-full.gif)

- **Deterministic** — reproducible scoring for the same repository state and evaluation date. No AI, no randomness.
- **Zero-config** — works out of the box on any Git repository.
- **Fast** — analyzes most repositories in under 3 seconds.
- **Offline** — no network access, no API keys, no accounts.

**What RepoHealth is NOT:**
- Not a code quality analyzer (use SonarQube)
- Not a security vulnerability scanner (use Trivy, Snyk)
- Not a test runner or coverage tool (use your framework's CLI)
- Not a dependency update tool (use Dependabot, Renovate)

RepoHealth measures *repository maturity and project hygiene*, not code quality.

**Use cases:** Pre-publish repo audit &middot; CI quality gates &middot; OSS evaluation before contributing

## Example Output

Sample output from `repohealth .` on a checkout of this repository on 2026-09-15. Timings vary by machine, and the recent-commit, commit-frequency and lockfile-freshness checks compare git dates with the current date, so a later run can change those lines and the overall score.

```
$ repohealth .

  RepoHealth v0.5.3

  Repository: /home/dev/repohealth
  Languages:  Go (95%), Python (5%)
  Analyzed:   67 files in 45ms

  ──────────────────────────────────────────────

  Overall Score:  81 / 100    Grade: B+

  ──────────────────────────────────────────────

  Documentation                       15 / 15
    README exists                      ✓  README.md found
    README has content                 ✓  README has substantive content with sections
    LICENSE exists                     ✓  LICENSE found
    CONTRIBUTING exists                ✓  CONTRIBUTING.md found
    CODE_OF_CONDUCT exists             ✓  CODE_OF_CONDUCT.md found
    SECURITY.md exists                 ✓  SECURITY.md found
    CHANGELOG exists                   ✓  CHANGELOG.md found

  Testing                             17 / 20
    Test files detected                ✓  16 test files
    Test directory exists              ✓  testdata/healthy-repo/tests/
    Test framework configured          ✓  go test (built-in)
    Coverage config exists             ✗  No coverage configuration found
    Test-to-source ratio               ✓  16 test files / 24 source files (67%)

  CI/CD                               15 / 15
    CI configuration exists            ✓  GitHub Actions
    CI runs tests                      ✓  Test command found in CI
    CI runs linter                     ✓  Linter found in CI
    CI runs build                      ✓  Build command found in CI

  Dependencies                        8 / 9
    Lockfile exists                    ✓  go.sum found
    Package manager detected           ✓  go.mod found
    Lockfile freshness                 ◐  go.sum not updated in over 90 days
    Dependency count                   ✓  2 dependencies declared

  Security                            7 / 10
    No secrets in repo                 ✓  No secret patterns detected in source files
    .gitignore covers secrets          ◐  .gitignore covers some but not all secret patterns
    Dependency pinning                 ✓  go.sum pins dependency versions
    Branch protection indicators       ✗  No branch protection indicators found

  Code Statistics                     3 / 5
    Source files exist                 ✓  24 source files
    Language diversity                 ✓  2 languages
    Comment ratio                      ✗  4% comment ratio
    No vendor bloat                    ✓  No vendor bloat detected

  Activity                            8 / 15
    Recent commit                      ✓  today
    Commit frequency                   ◐  24 commits in last 6 months
    Contributors                       ✗  1 contributor
    Release exists                     ✓  10 release tags
    Bus factor                         ✗  bus factor 1

  TODO / Technical Debt               5 / 7
    TODO/FIXME count                   ◐  6 TODO/FIXME markers found
    TODO density per KLOC              ✓  1.0 TODO/FIXME markers per KLOC
    No critical TODO markers           ✓  No critical TODO markers found

  ──────────────────────────────────────────────

  Suggestions (sorted by impact)
    +3 pts  Single contributor — consider inviting collaborators
    +3 pts  Add coverage configuration (e.g. codecov.yml or .nycrc)
    +2 pts  Increase commit frequency to show active development
    +2 pts  Only 1 contributor holds >10% of commits — high bus factor risk
    +2 pts  Add a CODEOWNERS file or branch protection configuration
    +2 pts  Add comments to explain non-obvious code (currently <5%)
    +2 pts  Resolve or track TODO/FIXME markers as issues
    +1 pt   Update your dependencies to keep the lockfile current
    +1 pt   Add .env, *.pem, *.key, and credentials entries to .gitignore

  ──────────────────────────────────────────────

  Improvement Plan
    81 → 84  Single contributor — consider inviting collaborators
    84 → 87  Add coverage configuration (e.g. codecov.yml or .nycrc)
    87 → 89  Increase commit frequency to show active development
    89 → 91  Only 1 contributor holds >10% of commits — high bus factor risk
    91 → 93  Add a CODEOWNERS file or branch protection configuration
```

## Real Repository Scores

RepoHealth measures repository hygiene and project-maintenance signals, not source-code quality, security posture, or overall project quality. The table reports RepoHealth results for the pinned revisions listed in the manifest; it is not a judgment of the projects.

v0.5.3 benchmark: 36 checks, 19 repositories, evaluated on 2026-09-17 (UTC) against pinned upstream commits with full git history. Each repository was scanned twice and produced identical results apart from timestamp and duration. Activity and dependency-freshness checks depend on the evaluation date, so the scores are a snapshot rather than a permanent rating. Pinned commits, per-category scores and skipped checks are in [benchmarks/v0.5.3.json](benchmarks/v0.5.3.json). Earlier benchmark runs used older RepoHealth versions, other repository revisions and in some cases shallow history, so they are not directly comparable.

| Repository | Score | Grade | Files |
|------------|-------|-------|-------|
| [prometheus/prometheus](https://github.com/prometheus/prometheus) | 86 | A- | 1,668 |
| [gin-gonic/gin](https://github.com/gin-gonic/gin) | 85 | A- | 130 |
| [hashicorp/terraform](https://github.com/hashicorp/terraform) | 82 | B+ | 5,452 |
| [docker/compose](https://github.com/docker/compose) | 82 | B+ | 835 |
| [facebook/react](https://github.com/facebook/react) | 82 | B+ | 7,195 |
| [grafana/grafana](https://github.com/grafana/grafana) | 81 | B+ | 22,738 |
| [vercel/next.js](https://github.com/vercel/next.js) | 72 | B- | 30,709 |
| [spf13/cobra](https://github.com/spf13/cobra) | 70 | B- | 66 |
| [rails/rails](https://github.com/rails/rails) | 68 | C+ | 4,960 |
| [python/cpython](https://github.com/python/cpython) | 68 | C+ | 6,268 |
| [rust-lang/rust](https://github.com/rust-lang/rust) | 68 | C+ | 61,340 |
| [pallets/flask](https://github.com/pallets/flask) | 67 | C+ | 232 |
| [fastapi/fastapi](https://github.com/fastapi/fastapi) | 66 | C+ | 3,139 |
| [django/django](https://github.com/django/django) | 66 | C+ | 7,006 |
| [kubernetes/kubernetes](https://github.com/kubernetes/kubernetes) | 65 | C+ | 25,788 |
| [vuejs/vue](https://github.com/vuejs/vue) | 65 | C+ | 505 |
| [expressjs/express](https://github.com/expressjs/express) | 55 | C- | 212 |
| [laravel/laravel](https://github.com/laravel/laravel) | 53 | D | 52 |
| [golang/go](https://github.com/golang/go) | 51 | D | 14,698 |

## What It Checks

RepoHealth runs 36 checks across 8 categories:

| Category | What It Measures | Checks |
|----------|-----------------|--------|
| **Documentation** | README, LICENSE, CONTRIBUTING, CODE_OF_CONDUCT, SECURITY, CHANGELOG | 7 |
| **Testing** | Test files, directories, framework config, coverage config, test-to-source ratio | 5 |
| **CI/CD** | CI presence, runs tests, runs linter, runs build | 4 |
| **Dependencies** | Lockfile, package manager, freshness, dependency count | 4 |
| **Security** | Secret scanning, .gitignore coverage, dependency pinning, branch protection | 4 |
| **Code Statistics** | Source files, language diversity, comment ratio, vendor bloat | 4 |
| **Activity** | Last commit, commit frequency, contributors, releases, bus factor | 5 |
| **TODO / Debt** | TODO count, density per KLOC, critical markers | 3 |

Each check contributes points. The total is normalized to 0-100 and graded A+ through F.

## How Scoring Works

| Grade | Score | Meaning |
|-------|-------|---------|
| A+ | 95-100 | Exceptional — production-grade, well-governed |
| A / A- | 85-94 | Excellent — strong across all dimensions |
| B+ / B / B- | 70-84 | Good — solid fundamentals, clear improvement areas |
| C+ / C / C- | 55-69 | Needs improvement — notable gaps |
| D | 40-54 | Failing — major investment needed |
| F | 0-39 | Critical — fundamental project hygiene missing |

Every check that scores below full generates a specific, actionable suggestion sorted by potential point impact.

## Installation

**Go install** (requires Go 1.22+):

```bash
go install github.com/spbuilds/repohealth/cmd/repohealth@latest
```

**Download binary** from [GitHub Releases](https://github.com/spbuilds/repohealth/releases):

```bash
# macOS / Linux
curl -sSL https://github.com/spbuilds/repohealth/releases/latest/download/repohealth_$(uname -s)_$(uname -m).tar.gz | tar xz
sudo mv repohealth /usr/local/bin/
```

**Build from source:**

```bash
git clone https://github.com/spbuilds/repohealth.git
cd repohealth
make build
```

## Usage

```bash
# Analyze current directory
repohealth .

# Analyze a specific repo
repohealth /path/to/repo

# JSON output for CI pipelines, scripts, and dashboards
repohealth . --format json

# Score only — single line for badges and automation
repohealth . --score-only
# Output: 81/100 (B+)

# CI quality gate — fail if below threshold
repohealth . --ci --threshold 70

# HTML report
repohealth . --format html > report.html

# Markdown report
repohealth . --format markdown

# Use custom config
repohealth . --config .repohealthrc.yaml

# Disable colored output
repohealth . --no-color
```

## CI Integration

Add RepoHealth to your GitHub Actions workflow:

```yaml
- name: Install RepoHealth
  run: go install github.com/spbuilds/repohealth/cmd/repohealth@latest

- name: Check repo health
  run: repohealth . --ci --threshold 70
```

RepoHealth auto-detects CI environments (`CI=true`, `GITHUB_ACTIONS=true`) and disables colors. Exit code 2 means the score is below threshold.

**Output formats in CI:**

```bash
repohealth . --score-only          # 81/100 (B+)
repohealth . --format json         # full JSON for dashboards
repohealth . --format markdown     # Markdown for PR comments
repohealth . --format html > report.html  # standalone HTML report
```

## Release History

| Version | What's Included |
|---------|----------------|
| **v0.1** | 13 checks, terminal + JSON output, scoring engine |
| **v0.2** | 36 checks across 8 categories, Markdown output, CI mode |
| **v0.3** | HTML report, config file, CI auto-detection, improvement plan |
| **v0.4** | Accuracy improvements, secret patterns, CI parsing, GitHub Action |
| **v0.5** | Reliability fixes, 99 tests, deterministic output, performance optimization |

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for how to get started.

## License

MIT License. See [LICENSE](LICENSE) for details.
