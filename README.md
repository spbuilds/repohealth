# RepoHealth

[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go](https://img.shields.io/badge/go-1.22+-00ADD8.svg)](https://go.dev)
[![CI](https://github.com/spbuilds/repohealth/actions/workflows/ci.yml/badge.svg)](https://github.com/spbuilds/repohealth/actions/workflows/ci.yml)

Deterministic CLI tool that analyzes repository health and produces a unified score.

**Lighthouse for Git repositories.**

```
$ repohealth .

  RepoHealth v0.1.0

  Repository: /home/dev/my-project
  Languages:  Go (74%), Shell (18%), Dockerfile (8%)
  Analyzed:   247 files in 1.2s

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

## Features

- **Deterministic** — same repo always produces the same score. No AI, no randomness.
- **Zero-config** — works out of the box on any git repository.
- **Fast** — full analysis in under 3 seconds.
- **Offline** — core checks require no network access.
- **Actionable** — every finding includes a specific recommendation.

## Installation

```bash
go install github.com/spbuilds/repohealth/cmd/repohealth@latest
```

Or build from source:

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

# Disable colored output
repohealth . --no-color
```

## What It Checks

| Category | Checks | Max Points |
|----------|--------|------------|
| Documentation | README, LICENSE, CONTRIBUTING, CODE_OF_CONDUCT, SECURITY, CHANGELOG | 15 |
| Testing | Test files, test directories, framework config | 14 |
| CI/CD | GitHub Actions, GitLab CI, Jenkins, CircleCI | 6 |
| Activity | Last commit recency, contributor count | 8 |
| **Total** | **14 checks** | **43** |

Score is normalized to 0-100 and graded A+ through F.

> More check categories (dependencies, security, code stats, TODO scanning) coming in v0.2.

## Scoring

| Grade | Score | Meaning |
|-------|-------|---------|
| A+ | 95-100 | Exceptional |
| A / A- | 85-94 | Excellent |
| B+ / B / B- | 70-84 | Good |
| C+ / C / C- | 55-69 | Needs improvement |
| D | 40-54 | Failing |
| F | 0-39 | Critical |

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for how to get started.

## License

MIT License. See [LICENSE](LICENSE) for details.
