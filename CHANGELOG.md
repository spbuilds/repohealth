# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
