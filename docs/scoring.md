# Scoring Model

RepoHealth produces a composite health score from 0 to 100, graded A+ through F.

## How the Score Is Calculated

1. Each check runs and returns one of four statuses:
   - **Full** — check passes completely (100% of max points)
   - **Partial** — check partially passes (50% of max points)
   - **None** — check fails (0 points)
   - **Skipped** — check could not run (e.g., no git history available)

2. Points are summed per category.

3. If a check is skipped, its points are **redistributed** proportionally within the same category. This ensures the score is always 0-100 regardless of which checks run. For example, if a repo has no `.git` directory, the activity checks are skipped and the score is based only on docs, tests, and CI/CD.

4. If an entire category is skipped (all checks in the category return Skipped), the category is **excluded** from the total — it does not penalize the score.

5. The composite score is normalized to a 0-100 scale.

## Categories and Weights

| Category | Max Points | Checks | Weight |
|----------|-----------|--------|--------|
| Documentation | 15 | 7 | 35% |
| Testing | 14 | 3 | 33% |
| CI/CD | 6 | 1 | 14% |
| Activity | 8 | 2 | 19% |
| **Total** | **43** | **13** | **100%** |

> v0.2 will add 4 more categories: Dependencies, Security, Code Statistics, and TODO/Technical Debt — bringing the total to 33 checks across 8 categories.

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

Every check that does not score Full generates a suggestion. Suggestions are sorted by **impact** (highest potential point gain first), so you always know the most valuable improvement to make next.

Example:
```
Suggestions (sorted by impact)
  +8 pts  Add test files for your code
  +6 pts  Add CI/CD configuration
  +3 pts  Add CONTRIBUTING.md
```

## Redistribution Example

If a repo has no `.git` directory:
- ACT-01 (Recent commit, 5 pts) → Skipped
- ACT-03 (Contributors, 3 pts) → Skipped
- Activity category is fully skipped → **excluded from scoring**
- Score is computed from Documentation + Testing + CI/CD only
- Result is still 0-100

This prevents non-git directories from being unfairly penalized on activity metrics they cannot produce.
