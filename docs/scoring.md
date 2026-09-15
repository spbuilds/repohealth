# Scoring Model

RepoHealth produces a composite health score from 0 to 100, graded A+ through F.

## How the Score Is Calculated

1. Each check runs and returns one of four statuses:
   - **Full** — the check passes completely (all of its points)
   - **Partial** — the check partially passes (a check-specific share of its points; see [checks.md](checks.md))
   - **None** — the check fails (0 points)
   - **Skipped** — the check could not run (for example, no git history is available)

2. Points are summed per category.

3. If a check is skipped, its points are **redistributed** proportionally within the same category: the category's earned points are scaled by `(category max) / (max of the checks that ran)`, rounded to the nearest integer and capped at the category max. A category is never penalised for a check that could not run.

4. If an entire category is skipped (every check in the category returns Skipped), the category is **excluded** from the total — it does not appear in the report and does not affect the score.

5. The composite score is `round(total points / total max × 100)`, where the total max is the sum of the included categories' maximums. This is the `raw_max` value in JSON output: 96 when every category is included, less when a category is excluded.

## Categories and Weights

With every check applicable there are 36 checks and 96 raw points. A category's weight is its share of those points:

| Category | Max Points | Checks | Weight |
|----------|-----------:|-------:|-------:|
| Documentation | 15 | 7 | 15.6% |
| Testing | 20 | 5 | 20.8% |
| CI/CD | 15 | 4 | 15.6% |
| Dependencies | 9 | 4 | 9.4% |
| Security | 10 | 4 | 10.4% |
| Code Statistics | 5 | 4 | 5.2% |
| Activity | 15 | 5 | 15.6% |
| TODO / Technical Debt | 7 | 3 | 7.3% |
| **Total** | **96** | **36** | **100%** |

Weights are fixed by the point values of the checks. The `weights` key accepted in `.repohealthrc.yaml` is reserved and is not currently applied to scoring.

## Grade Scale

| Grade | Score Range | Meaning |
|-------|------------|---------|
| A+ | 95 - 100 | Exceptional — production-grade, well-governed |
| A | 90 - 94 | Excellent — strong across all dimensions |
| A- | 85 - 89 | Very good — minor gaps |
| B+ | 80 - 84 | Good — solid fundamentals, some areas to improve |
| B | 75 - 79 | Above average — clear improvement areas |
| B- | 70 - 74 | Acceptable — meets minimum quality bar |
| C+ | 65 - 69 | Below average — notable gaps |
| C | 60 - 64 | Needs improvement — significant gaps |
| C- | 55 - 59 | Poor — multiple areas failing |
| D | 40 - 54 | Failing — major investment needed |
| F | 0 - 39 | Critical — fundamental project hygiene missing |

## Recommendations

Every check that returns Partial or None generates a suggestion. Suggestions are sorted by **impact** (the points the check is missing, highest first; ties broken by check ID), so the most valuable improvement is always listed first.

Example:
```
Suggestions (sorted by impact)
  +8 pts  Add test files for your code
  +6 pts  Add CI/CD configuration (GitHub Actions recommended)
  +2 pts  Add a CONTRIBUTING.md with contribution guidelines
```

The improvement plan shown in reports projects the score after each of the top five suggestions, adding `impact × 100 / raw max` (rounded down) at each step and capping at 100.

## Redistribution Example

If a repository has no `.git` directory:
- ACT-01 to ACT-05 (Activity, 15 points) all return Skipped
- The Activity category is fully skipped → **excluded from scoring**
- `raw_max` becomes 81 and the score is computed from the remaining seven categories
- The result is still 0-100

This prevents non-git directories from being unfairly penalised on activity metrics they cannot produce.

If only some checks in a category are skipped — for example DEP-03 (Lockfile freshness) and DEP-05 (Dependency count) when a repository has no lockfile or manifest — the points of the checks that did run are scaled up to the category's full maximum, so the category still counts for its usual weight.

## Disabling Checks

Checks can be disabled per repository by listing their IDs under `disable` in `.repohealthrc.yaml`:

```yaml
disable:
  - DOC-05
  - ACT-05
```

A disabled check does not run and does not appear in the report. Its points are removed from its category's maximum (and from `raw_max`) rather than redistributed, so the score is computed over the remaining checks only. Disabling every check in a category removes the category from the report.

## Determinism

Given the same working tree and git history, RepoHealth produces the same score: checks run in a fixed order, suggestions use a stable sort, comment-ratio sampling is path-ordered, and no network input is used. Checks that compare dates against the current date — ACT-01 (Recent commit), ACT-02 (Commit frequency) and DEP-03 (Lockfile freshness) — can change as time passes.
