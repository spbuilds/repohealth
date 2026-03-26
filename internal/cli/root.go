package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/spbuilds/repohealth/internal/checks"
	"github.com/spbuilds/repohealth/internal/report"
	"github.com/spbuilds/repohealth/internal/scanner"
	"github.com/spbuilds/repohealth/internal/scoring"
	"github.com/spf13/cobra"
)

// Version is set at build time via ldflags.
var Version = "dev"

var noColor bool
var format string
var scoreOnly bool

var rootCmd = &cobra.Command{
	Use:   "repohealth [path]",
	Short: "Analyze repository health and produce a unified score",
	Long:  "RepoHealth is a deterministic CLI tool that analyzes repository health across documentation, tests, CI/CD, and activity — producing a unified score (0-100) with actionable recommendations.",
	Args:  cobra.MaximumNArgs(1),
	RunE:  run,
}

func init() {
	rootCmd.Flags().BoolVar(&noColor, "no-color", false, "Disable colored output")
	rootCmd.Flags().StringVarP(&format, "format", "f", "terminal", "Output format: terminal, json")
	rootCmd.Flags().BoolVarP(&scoreOnly, "score-only", "s", false, "Output only the score")
	rootCmd.Version = Version
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	startTime := time.Now()

	if noColor || os.Getenv("NO_COLOR") != "" {
		color.NoColor = true
	}

	// Determine repo path
	repoPath := "."
	if len(args) > 0 {
		repoPath = args[0]
	}

	// Validate path exists
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		return fmt.Errorf("path does not exist: %s", repoPath)
	}

	// Scan repository
	ctx, err := scanner.Scan(repoPath, nil)
	if err != nil {
		return fmt.Errorf("failed to scan repository: %w", err)
	}

	// Create check registry and get all checks
	registry := checks.NewRegistry()
	allChecks := registry.All()

	// Run checks
	results := checks.Run(allChecks, ctx)

	// Score
	r := scoring.Score(results, ctx.RepoPath, ctx.Languages, len(ctx.Files), startTime)
	r.Version = Version

	if scoreOnly {
		fmt.Fprintf(os.Stdout, "%d/%d (%s)\n", r.Score, r.MaxScore, r.Grade)
		return nil
	}

	switch format {
	case "json":
		return report.JSON(os.Stdout, r)
	case "terminal":
		report.Terminal(os.Stdout, r, Version)
	default:
		return fmt.Errorf("unknown format: %s", format)
	}

	return nil
}
