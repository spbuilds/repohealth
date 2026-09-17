# Check Reference

RepoHealth runs 36 checks across 8 categories. With every check applicable the checks total 96 raw points, which are normalised to a 0–100 score; see [scoring.md](scoring.md) for how skipped checks and categories are handled.

Each check returns one of four statuses: **Full** (all of its points), **Partial** (a check-specific share, shown in parentheses below), **None** (0) or **Skipped** (the check could not run). Check IDs have been unchanged since v0.2.0; `DEP-04` and `SEC-01` are not assigned.

## Documentation (15 points)

| ID | Check | Points | Detection |
|----|-------|-------:|-----------|
| DOC-01 | README exists | 4 | Full if `README.md`, `README`, `README.rst`, `readme.md` or `Readme.md` exists at the repository root; otherwise None. |
| DOC-02 | README has content | 2 | Full when the root README is larger than 100 bytes and has at least 3 lines starting with `#` (Markdown headings) whose text contains a recognised section keyword (install, usage, example, getting started, contribut, license, config, api, feature, setup, quickstart). Partial (1) when the README exists but is 100 bytes or smaller, or has fewer than 3 such lines. None when there is no README. |
| DOC-03 | LICENSE exists | 3 | Full if `LICENSE`, `LICENSE.md`, `LICENCE`, `COPYING` or `LICENSE.txt` exists at the root; otherwise None. |
| DOC-04 | CONTRIBUTING exists | 2 | Full if `CONTRIBUTING.md` or `CONTRIBUTING` exists at the root, or `.github/CONTRIBUTING.md` exists; otherwise None. |
| DOC-05 | CODE_OF_CONDUCT exists | 1 | Full if `CODE_OF_CONDUCT.md` exists at the root or `.github/CODE_OF_CONDUCT.md` exists; otherwise None. |
| DOC-06 | SECURITY.md exists | 2 | Full if `SECURITY.md` exists at the root or `.github/SECURITY.md` exists; otherwise None. |
| DOC-07 | CHANGELOG exists | 1 | Full if `CHANGELOG.md`, `CHANGELOG`, `HISTORY.md` or `CHANGES.md` exists at the root; otherwise None. |

## Testing (20 points)

| ID | Check | Points | Detection |
|----|-------|-------:|-----------|
| TST-01 | Test files detected | 8 | Full if at least one file matches a recognised test-file convention (see [Test file conventions](#test-file-conventions)); otherwise None. |
| TST-02 | Test directory exists | 2 | Full if a directory named `test`, `tests`, `__tests__` or `spec` exists at any depth; otherwise None. |
| TST-03 | Test framework configured | 4 | Full if a framework configuration file exists anywhere in the tree (`jest.config.js/.ts/.mjs/.cjs`, `vitest.config.ts/.js/.mts`, `pytest.ini`, `.mocharc.yml/.yaml/.js`, `phpunit.xml`, `phpunit.xml.dist`), or the first `pyproject.toml` found mentions `pytest`, or any `*_test.go` file exists (Go's built-in test runner); otherwise None. When several of these are present, the details name a single framework chosen by a fixed priority: Jest, Vitest, pytest (`pytest.ini`), Mocha, PHPUnit, then a `pyproject.toml` that mentions pytest, then Go's built-in runner. |
| TST-04 | Coverage config exists | 3 | Full if `.nycrc`, `.nycrc.json`, `.coveragerc`, `coverage.xml`, `codecov.yml` or `.coveralls.yml` exists anywhere in the tree, or the first `package.json` found contains `"coverage"` (for example a coverage script); otherwise None. |
| TST-05 | Test-to-source ratio | 3 | Ratio of test files (same conventions as TST-01) to non-test programming-language source files (markup and stylesheets excluded). Above 30% → Full; 10–30% → Partial (1); below 10% → None. Skipped when there are no source files. |

## CI/CD (15 points)

| ID | Check | Points | Detection |
|----|-------|-------:|-----------|
| CI-01 | CI configuration exists | 6 | Full if a `.github/workflows/`, `.circleci/` or `.buildkite/` directory exists at the root, or a `.gitlab-ci.yml`, `Jenkinsfile`, `.travis.yml`, `bitbucket-pipelines.yml`, `azure-pipelines.yml` or `Taskfile.yml` file exists; otherwise None. |
| CI-02 | CI runs tests | 4 | Reads the CI configuration at the repository root (YAML files under `.github/workflows/`, `.circleci/` and `.buildkite/`, plus root-level `.gitlab-ci.yml`, `.travis.yml`, `Jenkinsfile`, `azure-pipelines.yml`, `bitbucket-pipelines.yml` and `Taskfile.yml`) and looks, case-insensitively, for a test command: `npm test`, `npm run test`, `pnpm test`, `pnpm run test`, `yarn test`, `yarn run test`, `bun test`, `bun run test`, `vitest`, `jest`, `mocha`, `ava`, `pytest`, `python -m pytest`, `tox`, `go test`, `gotestsum`, `cargo test`, `mvn test`, `gradle test`, `make test`, `turbo test`, `nx test`, `rake test`, `rspec`, `bundle exec rspec` or `phpunit`. Found → Full; not found → None. Skipped when no CI configuration exists. |
| CI-03 | CI runs linter | 3 | Same CI files as CI-02; looks for a linter: `eslint`, `biome`, `oxlint`, `prettier --check`, `ruff`, `flake8`, `pylint`, `mypy`, `black --check`, `golangci-lint`, `go vet`, `staticcheck`, `clippy`, `cargo clippy`, `checkstyle`, `spotbugs`, `rubocop`, `standardrb`, `make lint`, `pnpm lint`, `yarn lint` or `npm run lint`. Found → Full; not found → None. Skipped when no CI configuration exists. |
| CI-04 | CI runs build | 2 | Same CI files as CI-02; looks for a build command: `npm run build`, `pnpm build`, `yarn build`, `bun build`, `vite build`, `next build`, `nuxt build`, `go build`, `goreleaser`, `cargo build`, `mvn package`, `gradle build`, `gradle assemble`, `make build`, `turbo build`, `nx build` or `docker build`. Found → Full; not found → None. Skipped when no CI configuration exists. |

## Dependencies (9 points)

| ID | Check | Points | Detection |
|----|-------|-------:|-----------|
| DEP-01 | Lockfile exists | 4 | Full if `package-lock.json`, `yarn.lock`, `pnpm-lock.yaml`, `Pipfile.lock`, `poetry.lock`, `Cargo.lock`, `go.sum`, `Gemfile.lock` or `composer.lock` exists at the root; otherwise None. |
| DEP-02 | Package manager detected | 2 | Full if `package.json`, `go.mod`, `Cargo.toml`, `pyproject.toml`, `requirements.txt`, `Gemfile`, `pom.xml` or `build.gradle` exists at the root; otherwise None. |
| DEP-03 | Lockfile freshness | 2 | Age of the root lockfile, taken from its last commit (`git log -1 -- <lockfile>`) or from the file's modification time when no git history is available. Up to 90 days → Full; 91–180 days → Partial (1); older → None. Skipped when there is no lockfile. |
| DEP-05 | Dependency count | 1 | Counts declared dependencies in a root `package.json` (`dependencies` and `devDependencies` entries) or, if there is none, a root `go.mod` (`require` entries, excluding `// indirect` lines inside `require (...)` blocks). Fewer than 50 → Full; 50–100 → Partial (0); more than 100 → None. Skipped when neither manifest is present. |

## Security (10 points)

| ID | Check | Points | Detection |
|----|-------|-------:|-----------|
| SEC-02 | No secrets in repo | 4 | None if a `.env` file, or a `.env.*` file other than `.env.example`, `.env.sample`, `.env.template`, `.env.production`, `.env.development`, `.env.staging`, `.env.local` or `.env.test`, is present at the repository root. Otherwise scans source-code files (all languages the scanner recognises, excluding test files and files under `testdata/` or `fixtures/` directories) line by line for known credential formats — AWS access keys, GitHub tokens (`ghp_`, `gho_`), Slack tokens, Stripe live keys, Google API keys, private key blocks — and for generic `password`/`secret`/`api_key`/`apikey`/`token` assignments to a quoted value of 8 or more characters. Any match → None; otherwise Full. Configuration files (`.yml`, `.yaml`, `.json`, `.toml`) and documentation (`.md`) are not scanned. |
| SEC-03 | .gitignore covers secrets | 2 | Reads the root `.gitignore` and looks for the entries `.env`, `*.pem`, `*.key` and `credentials`. All four → Full; some → Partial (1); none, or no `.gitignore` → None. |
| SEC-04 | Dependency pinning | 2 | Full if a lockfile (same list as DEP-01) exists at the root; otherwise None. |
| SEC-05 | Branch protection indicators | 2 | Full if `CODEOWNERS` exists at the root, or `.github/CODEOWNERS` or `.github/branch-protection.yml` exists; otherwise None. |

## Code Statistics (5 points)

| ID | Check | Points | Detection |
|----|-------|-------:|-----------|
| STAT-01 | Source files exist | 1 | Full if at least one non-test file has a recognised source-code extension (see [Source files](#source-files)); otherwise None. |
| STAT-02 | Language diversity | 1 | Full if at least one source-code language is detected from file extensions (YAML, JSON, TOML and Markdown do not count); None when no source-code language is detected. |
| STAT-03 | Comment ratio | 2 | Samples up to 50 non-test programming-language files (sorted by path; markup and stylesheets excluded) and counts lines that start with `//`, `#`, `--`, `/*`, `* ` or consist of `*` as comment lines. Above 10% of lines → Full; 5–10% → Partial (1); below 5% → None. Skipped when there are no such files or no readable lines. |
| STAT-04 | No vendor bloat | 1 | None if a `vendor/` or `node_modules/` directory exists at the root and is not listed in `.gitignore`; otherwise Full. |

## Activity (15 points)

| ID | Check | Points | Detection |
|----|-------|-------:|-----------|
| ACT-01 | Recent commit | 5 | Days since the last commit (`git log -1`). Up to 30 → Full; 31–90 → Partial (2); over 90 → None. Skipped when git is unavailable. |
| ACT-02 | Commit frequency | 3 | Number of commits in the last 6 months (`git log --since="6 months ago"`). More than 50 → Full; 10–50 → Partial (1); fewer than 10 → None. Skipped when git is unavailable. |
| ACT-03 | Contributors | 3 | Number of unique authors (`git shortlog -sn --no-merges HEAD`). More than 5 → Full; 2–5 → Partial (1); 1 → None. Skipped when git is unavailable or there are no commits. |
| ACT-04 | Release exists | 2 | Full if the repository has at least one git tag (`git tag -l`); otherwise None. Skipped when git is unavailable. |
| ACT-05 | Bus factor | 2 | Number of authors with more than 10% of non-merge commits (`git shortlog -sn --no-merges HEAD`). 3 or more → Full; 2 → Partial (1); 1 → None. Also Full when commits exist but no single author exceeds 10%, since contributions are spread across many authors. Skipped when git is unavailable or there are no commits. |

Activity checks require a `.git` entry (directory or file) at the root. Without one every Activity check is Skipped and the category is excluded from the score. Git commands run with a 5-second timeout.

## TODO / Technical Debt (7 points)

| ID | Check | Points | Detection |
|----|-------|-------:|-----------|
| TODO-01 | TODO/FIXME count | 3 | Counts `TODO`, `FIXME`, `HACK` and `XXX` (case-insensitive) as whole words appearing after a comment opener valid for the file's language, ignoring openers inside single-line string literals (see [TODO scanning](#todo-scanning)). Multi-line string contents are not tracked. 0 → Full; 1–20 → Partial (1); more than 20 → None. |
| TODO-02 | TODO density per KLOC | 2 | Markers counted by TODO-01 per 1,000 lines of scanned source. Below 2.0 → Full; 2.0–5.0 → Partial (1); above 5.0 → None. Skipped when no source lines were read. |
| TODO-03 | No critical TODO markers | 2 | None if any comment counted by TODO-01 also contains `SECURITY`, `VULNERABILITY` or `UNSAFE` (case-insensitive) as a whole word; otherwise Full. |

The TODO scan covers every recognised source-code file, including test files.

## Source files

The scanner recognises files by extension. Source-code extensions are:

`.go` `.py` `.js` `.jsx` `.ts` `.tsx` `.rs` `.java` `.rb` `.php` `.c` `.cpp` `.h` `.cs` `.swift` `.kt` `.sh` `.bash` `.zsh` `.sql` `.r` `.lua` `.dart` `.ex` `.exs` `.html` `.css` `.scss`

`.md`, `.yml`, `.yaml`, `.json` and `.toml` files are recognised for language detection but are treated as documentation or configuration: they are not counted as source files and are not secret- or TODO-scanned.

`.html`, `.css` and `.scss` count as source files (STAT-01) and are secret- and TODO-scanned, but are left out of the two checks that measure programming-language files: STAT-03 (comment ratio) and TST-05 (test-to-source ratio).

Directories named `node_modules`, `vendor`, `.venv`, `dist`, `build`, `target`, `out`, `.next`, `.nuxt`, `.output`, `.svelte-kit`, `__pycache__`, `.tox`, `.mypy_cache`, `.pytest_cache`, `.cache`, `.turbo`, `coverage`, `.nyc_output`, `.idea`, `.vscode` and `.git` are never scanned, and dotfiles are only read at the repository root. Scanning stops after 100,000 files. Checks that read file contents (DOC-02, TST-03, TST-04, CI-02 to CI-04, DEP-05, SEC-02, SEC-03, STAT-03, STAT-04, TODO-01 to TODO-03) skip files larger than 100 KB and binary files, and read at most the first 10,000 lines of a file.

## Test file conventions

TST-01 and TST-05 count files whose names match one of these patterns:

| Language | Patterns |
|----------|----------|
| Go | `*_test.go` |
| Python | `*_test.py`, `test_*.py` |
| JavaScript / TypeScript | `*.test.js`, `*.test.jsx`, `*.test.ts`, `*.test.tsx`, `*.spec.js`, `*.spec.jsx`, `*.spec.ts`, `*.spec.tsx` |
| Rust | `*_test.rs` |
| Ruby | `*_test.rb` |
| Java | `*Test.java`, `*Tests.java` |
| C# | `*Test.cs`, `*Tests.cs` |
| Dart | `*_test.dart` |
| Elixir | `*_test.exs` |
| Lua | `*_test.lua`, `*_spec.lua` |
| R | `test-*.R`, `test_*.R`, `test-*.r`, `test_*.r` |

Files matching these conventions — except `*_test.py` — are also left out of the secret scan (SEC-02) and of the source-file counts used by STAT-01, STAT-03 and TST-05.

## TODO scanning

TODO-01 to TODO-03 share one scan over every recognised source-code file. A line is only counted when a marker appears after a comment opener valid for the file's language:

| Extensions | Comment openers |
|------------|-----------------|
| `.go` `.js` `.jsx` `.ts` `.tsx` `.java` `.c` `.cpp` `.h` `.cs` `.swift` `.kt` `.rs` `.dart` `.scss` | `//`, `/*` |
| `.php` | `//`, `/*`, `#` |
| `.py` `.rb` `.sh` `.bash` `.zsh` `.r` `.ex` `.exs` | `#` |
| `.sql` | `--`, `/*` |
| `.lua` | `--` |
| `.html` | `<!--` |
| `.css` | `/*` |

Markers are matched case-insensitively and as whole words (`TODO`, `FIXME`, `HACK`, `XXX`), so `// todo` counts but identifiers such as `todoList` do not. An opener that appears inside a single-line string literal (`"…"`, `'…'`, `` `…` ``, backslash escapes honoured) is ignored, and `://` is never treated as a comment opener, so URLs do not count. In languages with `/* */` comments, a line whose first non-blank character is `*` is treated as a block-comment continuation.

Known limitations: the scanner works line by line, so string literals spanning multiple lines are not tracked (a line inside a multi-line string that looks like a comment is counted), and an unmatched quote character earlier on the line — for example a Rust lifetime such as `'a` — hides a comment that follows it on the same line.

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
