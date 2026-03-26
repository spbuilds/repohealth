# RepoHealth

[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go](https://img.shields.io/badge/go-1.22+-00ADD8.svg)](https://go.dev)
[![CI](https://github.com/spbuilds/repohealth/actions/workflows/ci.yml/badge.svg)](https://github.com/spbuilds/repohealth/actions/workflows/ci.yml)

**A Lighthouse-style health check for Git repositories.**

RepoHealth analyzes a repository's documentation, tests, CI/CD configuration, and maintenance activity to produce a composite health score (0-100) with actionable improvement recommendations. One command, one report, one score.

- **Deterministic** — same repo always produces the same score. No AI, no randomness.
- **Zero-config** — works out of the box on any Git repository.
- **Fast** — analyzes 58,000 files in under 300ms.
- **Offline** — no network access, no API keys, no accounts.

**Use cases:** Pre-publish repo audit &middot; CI quality gates &middot; OSS evaluation before contributing

## Example Output

```
$ repohealth .

  RepoHealth v0.1.1

  Repository: /home/dev/my-project
  Languages:  Go (74%), Shell (18%), Dockerfile (8%)
  Analyzed:   247 files in 15ms

  ──────────────────────────────────────────────

  Overall Score:  78 / 100    Grade: B+

  ──────────────────────────────────────────────

  Documentation                           13 / 15
    README.md                              ✓
    LICENSE                                ✓
    CONTRIBUTING.md                        ✗
    CODE_OF_CONDUCT.md                     ✗
    SECURITY.md                            ✓
    CHANGELOG.md                           ✓

  Testing                                 14 / 14
    Test files detected                    ✓  (42 files)
    Test directory                         ✓  tests/
    Test framework configured              ✓  go test

  CI/CD                                    6 / 6
    CI configuration                       ✓  GitHub Actions

  Activity                                 6 / 8
    Last commit                            ✓  2 days ago
    Contributors                           ◐  3

  ──────────────────────────────────────────────

  Suggestions (sorted by impact)
    +3 pts  Add CONTRIBUTING.md
    +2 pts  Improve contributor count (bus factor)
    +1 pt   Add CODE_OF_CONDUCT.md
```

## Real Repository Scores

Tested on well-known open-source projects:

| Repository | Language | Score | Grade | Time |
|------------|----------|-------|-------|------|
| [axios](https://github.com/axios/axios) | JavaScript | 93 | A | 10ms |
| [requests](https://github.com/psf/requests) | Python | 93 | A | 9ms |
| [fastapi](https://github.com/tiangolo/fastapi) | Python | 88 | A- | 25ms |
| [hono](https://github.com/honojs/hono) | TypeScript | 81 | B+ | 11ms |
| [vscode](https://github.com/microsoft/vscode) | TypeScript | 79 | B | 91ms |
| [rust](https://github.com/rust-lang/rust) | Rust | 79 | B | 255ms |
| [kubernetes](https://github.com/kubernetes/kubernetes) | Go | 77 | B | 159ms |
| [redis](https://github.com/redis/redis) | C | 65 | C+ | 15ms |
| [express](https://github.com/expressjs/express) | JavaScript | 51 | D | 10ms |

Well-maintained libraries with complete docs and tests score highest. Large monorepos and older projects score lower due to non-standard CI and missing community files.

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
| **v0.2** | 33 checks across 8 categories, TODO scanning, deps, security, CI parsing | Released |
| **v0.3** | HTML report, config file, CI auto-detection, improvement plan, `--threshold` | Released |
| **v0.4** | Homebrew tap | Planned |
| **v0.5** | GitHub Action | Planned |
| **v1.0** | Stable release, plugin system | Planned |

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for how to get started.

## License

MIT License. See [LICENSE](LICENSE) for details.
