# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.4.0] - 2026-03-26

### Added
- GitHub Action composite (`action.yml`) for easy CI integration
- Improved CI detection: now detects vitest, pnpm, yarn, turbo, nx, bun, and many more
- Improved secret scanning: AWS keys, GitHub tokens, Slack tokens, Stripe keys, Google API keys, private key blocks, generic password/secret patterns
- README quality scoring: checks for section headings (install, usage, examples) not just file size
- Report output tests for Markdown, HTML, and JSON
- Version fallback via `debug.ReadBuildInfo` for `go install` users
- Build metadata: commit hash and date embedded in release binaries
- Performance guardrail: scanner stops at 100,000 files

### Fixed
- `repohealth --version` now shows correct version when installed via `go install`
- CI parsing no longer misses vitest, pnpm test, turbo test, and other modern tooling
- Deterministic output ordering for checks and suggestions

### Changed
- GoReleaser config now includes commit hash and build date in ldflags

## [0.1.1] - 2026-03-26

### Added
- `--format json` flag for machine-readable output
- `--score-only` flag for minimal output (e.g., `78/100 (B+)`)
- 24 unit tests across 4 packages (model, scanner, scoring, checks)
- Integration tests with testdata fixtures
- Example JSON report in `examples/`

### Fixed
- Scoring engine: all-skipped categories no longer penalize the score
- Scoring engine: use proper rounding instead of truncation
- Scanner: broken symlinks handled gracefully
- Scanner: descriptive error when path is not a directory
- Git commands: 5-second timeout prevents hangs on broken repos
- Activity check: zero contributors returns "skipped" instead of wrong message
- Cross-platform: paths normalized with forward slashes

### Changed
- Respects `NO_COLOR` environment variable

## [0.1.0] - 2026-03-26

### Added
- Initial release
- 14 checks across 4 categories (documentation, testing, CI/CD, activity)
- Composite scoring engine (0-100) with A+ through F grading
- Colored terminal report with actionable recommendations
