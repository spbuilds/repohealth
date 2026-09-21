# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.5.3] - 2026-09-21

### Fixed
- Source-language coverage: every language the scanner recognises is now included in the secret scan (SEC-02), source-file statistics (STAT-01), test-to-source ratio (TST-05) and TODO scan (TODO-01 to TODO-03). Previously `.tsx`, `.jsx`, `.cs`, `.dart`, `.ex`, `.exs`, `.lua`, `.r`, `.sql`, `.html`, `.css`, `.scss`, `.bash` and `.zsh` files were skipped by these checks
- TODO detection: `TODO`, `FIXME`, `HACK` and `XXX` are counted only as whole words following a comment opener valid for the file's language; openers inside single-line string literals are ignored. Marker text in code, single-line string literals, URLs and identifiers no longer counts
- Test-file recognition now covers C# (`*Test.cs`, `*Tests.cs`), Dart (`*_test.dart`), Elixir (`*_test.exs`), Lua (`*_test.lua`, `*_spec.lua`) and R (`test-*.R`, `test_*.R`, also with a lower-case `.r` extension); these files count for TST-01 and TST-05 and are excluded from the secret scan and source-file counts like other languages' test files
- Scanner: a repository path that is a symlink is resolved before scanning. Previously it scanned as an empty repository
- TST-03: when a repository contains configuration files for several test frameworks, the framework named in the check details is chosen in a fixed order (Jest, Vitest, pytest, Mocha, PHPUnit). Previously the name could differ between runs on the same repository
- ACT-05: a repository whose commits are spread so widely that no single author holds more than 10% of them is now scored Full, with details saying so. Previously it was skipped with the details "No commits found", which was inaccurate

### Changed
- STAT-03 (comment ratio) and TST-05 (test-to-source ratio) measure programming-language files only; markup and stylesheets (`.html`, `.css`, `.scss`) still count as source files for STAT-01 and are still secret- and TODO-scanned
- The repository path shown in reports (`repo_path` in JSON output) is the symlink-resolved path
- Documentation synchronized with current behaviour: `docs/checks.md`, `docs/scoring.md` and `examples/report.json` regenerated for the 36-check, 8-category model; stale check-count claims corrected in the README and changelog; the `weights` configuration key is documented as reserved and not currently applied

## [0.5.2] - 2026-03-29

### Fixed
- CI-02: `npm run test`, `pnpm run test`, `yarn run test`, `bun run test` now detected as test commands
- SEC-02: Environment-specific files (`.env.production`, `.env.local`, `.env.development`, `.env.staging`, `.env.test`) excluded from secret detection
- SEC-02: Secret pattern scanning skips test files (`*_test.go`, `*.spec.ts`, etc.) and `testdata/`/`fixtures/` directories
- STAT-04: Vendor bloat check reads `.gitignore` before penalizing `node_modules`/`vendor` on disk

### Added
- Demo GIF in README

## [0.5.0] - 2026-03-27

### Fixed
- Scoring engine: category score capped at max after rounding (prevented >100% per category)
- Scoring engine: improvement plan normalizes impact by RawMax (was mixing raw points with percentages)
- Bus factor: uses float64 threshold (was inflated for small repos via integer division)
- Suggestions: stable sort with CheckID tiebreaker for deterministic ordering
- Comment ratio: files sorted by path before sampling (was OS-dependent walk order)
- Language formatting: alphabetical tiebreaker for equal file counts
- Secret scan: string prefix pre-filter skips regex when prefix absent (~99% fewer evaluations)
- Git shortlog: cached across ContributorCount + BusFactor (single subprocess call)
- Private key regex: matches PKCS#8 format (`BEGIN PRIVATE KEY`)
- CI content parsing: reads Taskfile.yml, CircleCI, Buildkite, Azure, Bitbucket configs
- CI workflow files: `.yaml` extension now detected (was only `.yml`)
- `.env.local`/`.env.production`/`.env.staging` detected as committed secrets
- `pyproject.toml` only credited as pytest when `[tool.pytest]` section present
- `isTestFile` uses precise suffix matching (no false positives like `latest_data.py`)
- Go mod parser: excludes indirect dependencies from count
- Lockfile freshness: uses `git log` date instead of unreliable `mtime`
- `--config` flag: reads the specified file directly (was using `filepath.Dir`)
- STAT-04 vendor check: uses `os.Stat` (was checking `ctx.Files` which never contains vendor/)
- HTML output: `html.EscapeString` applied to all user-derived values
- Scanner: `Truncated` flag set at 100K files with stderr warning
- Empty files: return empty slice instead of `io.EOF` error
- `bufio.Scanner` buffer: increased to 256KB (was 64KB, silently truncating minified files)
- Path-based excludes: scanner supports full path exclusions, not just directory basename

### Added
- 99 unit tests across 7 packages (up from 24)
- Test files: `activity_test.go`, `cicd_test.go`, `stats_test.go`, `loader_test.go`, `reader_test.go`
- `RawMax` field in JSON output for improvement plan normalization
- `Truncated` field in `ScanContext` for large repo awareness
- `schema_version: "1.0"` in JSON output

### Removed
- `internal/config/config.go` (dead code — unused Config struct)
- `internal/scoring/improvement.go` (dead code — inlined in renderers)
- Duplicate `containsAny()` function in cicd.go

## [0.4.0] - 2026-03-26

### Added
- GitHub Action composite (`action.yml`) for easy CI integration
- Improved CI detection: vitest, pnpm, yarn, turbo, nx, bun, biome, and many more
- Improved secret scanning: AWS keys, GitHub tokens, Slack tokens, Stripe keys, Google API keys, private key blocks, generic password/secret patterns
- README quality scoring: checks for section headings, not just file size
- Report output tests for Markdown, HTML, and JSON
- Version fallback via `debug.ReadBuildInfo` for `go install` users
- Build metadata: commit hash and date embedded in release binaries
- Performance guardrail: scanner stops at 100,000 files

### Fixed
- `repohealth --version` now shows correct version when installed via `go install`
- CI parsing no longer misses vitest, pnpm test, turbo test, and other modern tooling
- Deterministic output ordering for checks and suggestions

## [0.3.0] - 2026-03-26

### Added
- `--format html` flag: standalone HTML report with inline CSS
- `.repohealthrc.yaml` config file support (weights, disable, exclude, threshold)
- `--config` flag to specify custom config path
- Improvement plan section in terminal, markdown, and HTML reports
- Auto-detect CI environments (`CI=true`, `GITHUB_ACTIONS=true`) for no-color
- `schema_version` field in JSON output

### Fixed
- Performance warning on stderr when analysis exceeds 3 seconds

## [0.2.0] - 2026-03-26

### Added
- 23 new checks and 4 new categories (36 total, 8 categories)
- **Dependencies:** lockfile, package manager, freshness, dependency count
- **Security:** secret scanning, .gitignore coverage, dependency pinning, branch protection
- **Code Statistics:** source files, language diversity, comment ratio, vendor bloat
- **TODO / Debt:** TODO/FIXME count, density per KLOC, critical markers
- Extended testing: coverage config, test-to-source ratio
- Extended CI/CD: CI runs tests, linter, build (content parsing)
- Extended activity: commit frequency, releases, bus factor
- `--format markdown` flag
- `--ci` and `--threshold` flags for CI quality gates
- `ReadFileLines` helper with binary detection and size limits
- Git helpers: `CommitCountSince`, `TagCount`, `BusFactor`

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
- 13 checks across 4 categories (documentation, testing, CI/CD, activity)
- Composite scoring engine (0-100) with A+ through F grading
- Colored terminal report with actionable recommendations
