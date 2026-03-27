# RepoHealth

[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go](https://img.shields.io/badge/go-1.22+-00ADD8.svg)](https://go.dev)
[![CI](https://github.com/spbuilds/repohealth/actions/workflows/ci.yml/badge.svg)](https://github.com/spbuilds/repohealth/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/spbuilds/repohealth)](https://goreportcard.com/report/github.com/spbuilds/repohealth)

**A Lighthouse-style health check for Git repositories.**

RepoHealth analyzes a repository's documentation, tests, CI/CD configuration, and maintenance activity to produce a composite health score (0-100) with actionable improvement recommendations. One command, one report, one score.

- **Deterministic** — same repo always produces the same score. No AI, no randomness.
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

```
$ repohealth .

  RepoHealth v0.5.0

  Repository: /home/dev/my-project
  Languages:  Go (74%), Shell (18%), Dockerfile (8%)
  Analyzed:   247 files in 45ms

  ──────────────────────────────────────────────

  Overall Score:  78 / 100    Grade: B+

  ──────────────────────────────────────────────

  Documentation                           13 / 15
    README exists                          ✓  README.md found
    README has content                     ✓  README has substantive content with sections
    LICENSE exists                         ✓  LICENSE found
    CONTRIBUTING exists                    ✗
    CODE_OF_CONDUCT exists                 ✗
    SECURITY.md exists                     ✓  SECURITY.md found
    CHANGELOG exists                       ✓  CHANGELOG.md found

  Testing                                 17 / 20
    Test files detected                    ✓  42 test files
    Test directory                         ✓  tests/
    Test framework configured              ✓  go test (built-in)
    Coverage config exists                 ✗
    Test-to-source ratio                   ✓  42 test files / 89 source files (47%)

  CI/CD                                   15 / 15
    CI configuration exists                ✓  GitHub Actions
    CI runs tests                          ✓  Test command found in CI
    CI runs linter                         ✓  Linter found in CI
    CI runs build                          ✓  Build command found in CI

  Dependencies                             9 / 9
    Lockfile exists                        ✓
    Package manager detected               ✓
    Lockfile freshness                     ✓
    Dependency count                       ✓

  Security                                 8 / 10
    No secrets in repo                     ✓
    .gitignore covers secrets              ✓
    Dependency pinning                     ✓
    Branch protection indicators           ✗

  Code Statistics                          4 / 5
    Source files exist                     ✓
    Language diversity                     ✓
    Comment ratio                          ◐
    No vendor bloat                        ✓

  Activity                                10 / 15
    Recent commit                          ✓  today
    Commit frequency                       ✓
    Contributors                           ◐  3
    Release exists                         ✓
    Bus factor                             ◐

  TODO / Technical Debt                    5 / 7
    TODO/FIXME count                       ◐
    TODO density per KLOC                  ✓
    No critical TODO markers               ✓

  ──────────────────────────────────────────────

  Suggestions (sorted by impact)
    +3 pts  Add coverage configuration
    +2 pts  Add a CODEOWNERS file or branch protection configuration
    +2 pts  Add CONTRIBUTING.md
    +1 pt   Add CODE_OF_CONDUCT.md

  ──────────────────────────────────────────────

  Improvement Plan
    78 → 81  Add coverage configuration
    81 → 83  Add CODEOWNERS file
    83 → 85  Add CONTRIBUTING.md
    85 → 86  Add CODE_OF_CONDUCT.md
```

## Real Repository Scores

Tested on well-known open-source projects (v0.5.0, 33 checks):

| Repository | Language | Score | Grade | Files | Time |
|------------|----------|-------|-------|-------|------|
| [gin](https://github.com/gin-gonic/gin) | Go | 77 | B | 130 | 75ms |
| [cobra](https://github.com/spf13/cobra) | Go | 64 | C | 66 | 55ms |
| [next.js](https://github.com/vercel/next.js) | JavaScript | 64 | C | 26,716 | 1.9s |
| [django](https://github.com/django/django) | Python | 55 | C- | 6,942 | 580ms |
| [rust](https://github.com/rust-lang/rust) | Rust | 56 | C- | 57,236 | 5.6s |
| [kubernetes](https://github.com/kubernetes/kubernetes) | Go | 53 | D | 23,674 | 2.2s |

Scores reflect all 8 categories: documentation, tests, CI/CD, dependencies, security, code statistics, activity, and TODO debt. Large monorepos score lower due to non-standard CI, missing community files, and high TODO counts.

## What It Checks

RepoHealth runs 33 checks across 8 categories:

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
# Output: 78/100 (B+)

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
repohealth . --score-only          # 78/100 (B+)
repohealth . --format json         # full JSON for dashboards
repohealth . --format markdown     # Markdown for PR comments
repohealth . --format html > report.html  # standalone HTML report
```

## Roadmap

| Version | What's Included | Status |
|---------|----------------|--------|
| **v0.1** | 14 checks, terminal + JSON output, scoring engine | Released |
| **v0.2** | 33 checks across 8 categories, Markdown output, CI mode | Released |
| **v0.3** | HTML report, config file, CI auto-detection, improvement plan | Released |
| **v0.4** | Accuracy improvements, secret patterns, CI parsing, GitHub Action | Released |
| **v0.5** | 51 reliability fixes, 99 tests, deterministic output, performance optimization | Released |
| **v0.6** | Homebrew tap, custom scoring weights | Planned |
| **v1.0** | Stable release, plugin system | Planned |

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for how to get started.

## License

MIT License. See [LICENSE](LICENSE) for details.
