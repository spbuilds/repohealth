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

	// Testing checks (TST-01 to TST-03)
	r.checks = append(r.checks,
		&TestFilesExistCheck{},
		&TestDirExistsCheck{},
		&TestFrameworkCheck{},
	)

	// CI/CD checks (CI-01)
	r.checks = append(r.checks,
		&CIConfigExistsCheck{},
	)

	// Activity checks (ACT-01, ACT-03)
	r.checks = append(r.checks,
		&RecentCommitCheck{},
		&ContributorCountCheck{},
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
