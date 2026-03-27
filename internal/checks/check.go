// Package checks implements the health check modules for RepoHealth.
package checks

import "github.com/spbuilds/repohealth/internal/model"

// Check is the interface every analysis module implements.
type Check interface {
	ID() string
	Category() string
	Name() string
	MaxPoints() int
	Run(ctx *model.ScanContext) model.CheckResult
}

// Registry holds all registered checks.
type Registry struct {
	checks []Check
}

// NewRegistry creates a registry with all built-in checks.
func NewRegistry() *Registry {
	r := &Registry{}

	// Documentation checks (DOC-01 to DOC-07)
	r.checks = append(r.checks,
		&ReadmeExistsCheck{},
		&ReadmeContentCheck{},
		&LicenseExistsCheck{},
		&ContributingExistsCheck{},
		&CodeOfConductExistsCheck{},
		&SecurityExistsCheck{},
		&ChangelogExistsCheck{},
	)

	// Testing checks (TST-01 to TST-05)
	r.checks = append(r.checks,
		&TestFilesExistCheck{},
		&TestDirExistsCheck{},
		&TestFrameworkCheck{},
		&CoverageConfigCheck{},
		&TestToSourceRatioCheck{},
	)

	// CI/CD checks (CI-01 to CI-04)
	r.checks = append(r.checks,
		&CIConfigExistsCheck{},
		&CIRunsTestsCheck{},
		&CIRunsLinterCheck{},
		&CIRunsBuildCheck{},
	)

	// Dependency checks (DEP-01, DEP-02, DEP-03, DEP-05; DEP-04 deferred)
	r.checks = append(r.checks,
		&LockfileExistsCheck{},
		&PackageManagerCheck{},
		&LockfileFreshnessCheck{},
		&DependencyCountCheck{},
	)

	// Security checks (SEC-02 to SEC-05)
	r.checks = append(r.checks,
		&NoSecretsCheck{},
		&GitignoreSecretsCheck{},
		&DependencyPinningCheck{},
		&BranchProtectionCheck{},
	)

	// Code statistics (STAT-01 to STAT-04)
	r.checks = append(r.checks,
		&SourceFilesExistCheck{},
		&LanguageDiversityCheck{},
		&CommentRatioCheck{},
		&NoVendorBloatCheck{},
	)

	// Activity checks (ACT-01 to ACT-05)
	r.checks = append(r.checks,
		&RecentCommitCheck{},
		&CommitFrequencyCheck{},
		&ContributorCountCheck{},
		&ReleaseExistsCheck{},
		&BusFactorCheck{},
	)

	// TODO / Technical Debt (TODO-01 to TODO-03)
	r.checks = append(r.checks,
		&TodoCountCheck{},
		&TodoDensityCheck{},
		&TodoCriticalCheck{},
	)

	return r
}

// All returns all registered checks.
func (r *Registry) All() []Check {
	return r.checks
}

// Filter returns checks matching the given categories, excluding disabled IDs.
func (r *Registry) Filter(categories []string, disabled []string) []Check {
	disabledSet := make(map[string]bool)
	for _, id := range disabled {
		disabledSet[id] = true
	}

	catSet := make(map[string]bool)
	for _, c := range categories {
		catSet[c] = true
	}

	var result []Check
	for _, c := range r.checks {
		if disabledSet[c.ID()] {
			continue
		}
		if len(categories) > 0 && !catSet[c.Category()] {
			continue
		}
		result = append(result, c)
	}
	return result
}

// Run executes all given checks and returns results.
func Run(checks []Check, ctx *model.ScanContext) []model.CheckResult {
	results := make([]model.CheckResult, 0, len(checks))
	for _, c := range checks {
		results = append(results, c.Run(ctx))
	}
	return results
}
