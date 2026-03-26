# Contributing to RepoHealth

Thanks for your interest in contributing to RepoHealth!

## Getting Started

```bash
# Clone the repo
git clone https://github.com/spbuilds/repohealth.git
cd repohealth

# Build
make build

# Run tests
make test

# Run linter
make lint
```

## Adding a New Check

1. Create a new file in `internal/checks/` (e.g., `mycheck.go`)
2. Implement the `Check` interface:

```go
type MyCheck struct{}

func (c *MyCheck) ID() string       { return "MY-01" }
func (c *MyCheck) Category() string  { return "mycat" }
func (c *MyCheck) Name() string      { return "My check description" }
func (c *MyCheck) MaxPoints() int    { return 5 }
func (c *MyCheck) Run(ctx *model.ScanContext) model.CheckResult { ... }
```

3. Register the check in `NewRegistry()` in `internal/checks/check.go`
4. Add tests in `internal/checks/mycheck_test.go`

## Pull Request Process

1. Fork the repo
2. Create a feature branch (`git checkout -b feat/my-check`)
3. Write tests for your changes
4. Ensure `make test` and `make lint` pass
5. Submit a PR with a clear description

## Code Style

- Follow standard Go conventions (`gofmt`)
- Keep functions focused and small
- Add tests for new checks
- Use the existing check implementations as reference
