// Package cli implements the command-line interface for RepoHealth.
package cli

import (
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"github.com/fatih/color"
	"github.com/spbuilds/repohealth/internal/checks"
	"github.com/spbuilds/repohealth/internal/config"
	"github.com/spbuilds/repohealth/internal/report"
	"github.com/spbuilds/repohealth/internal/scanner"
	"github.com/spbuilds/repohealth/internal/scoring"
	"github.com/spf13/cobra"
)

// Build metadata — set at build time via ldflags.
var Version = "dev"
var Commit = "none"
var Date = "unknown"

var noColor bool
var format string
var scoreOnly bool
var ciMode bool
var threshold int
var configPath string

var rootCmd = &cobra.Command{
	Use:   "repohealth [path]",
	Short: "Analyze repository health and produce a unified score",
	Long:  "RepoHealth is a deterministic CLI tool that analyzes repository health across documentation, tests, CI/CD, dependencies, security, activity, and code statistics — producing a unified score (0-100) with actionable recommendations.",
	Args:  cobra.MaximumNArgs(1),
	RunE:  run,
}

func init() {
	// Version fallback for go install (no ldflags)
	if Version == "dev" {
		if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
			Version = info.Main.Version
		}
	}

	rootCmd.Flags().BoolVar(&noColor, "no-color", false, "Disable colored output")
	rootCmd.Flags().StringVarP(&format, "format", "f", "terminal", "Output format: terminal, json, markdown, html")
	rootCmd.Flags().BoolVarP(&scoreOnly, "score-only", "s", false, "Output only the score")
	rootCmd.Flags().BoolVar(&ciMode, "ci", false, "CI mode: exit with code 2 if score below threshold")
	rootCmd.Flags().IntVarP(&threshold, "threshold", "t", 70, "Minimum passing score (used with --ci)")
	rootCmd.Flags().StringVar(&configPath, "config", "", "Path to config file")
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

	if noColor || os.Getenv("NO_COLOR") != "" || os.Getenv("CI") == "true" || os.Getenv("GITHUB_ACTIONS") == "true" {
		color.NoColor = true
	}

	repoPath := "."
	if len(args) > 0 {
		repoPath = args[0]
	}

	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		return fmt.Errorf("path does not exist: %s", repoPath)
	}

	// Load config
	fileCfg, err := config.LoadConfig(repoPath, configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Determine excludes from config
	var excludes []string
	if fileCfg != nil && len(fileCfg.Exclude) > 0 {
		excludes = fileCfg.Exclude
	}

	ctx, scanErr := scanner.Scan(repoPath, excludes)
	if scanErr != nil {
		return fmt.Errorf("failed to scan repository: %w", scanErr)
	}
	if ctx.Truncated {
		fmt.Fprintf(os.Stderr, "Warning: scan truncated at %d files. Results may be incomplete.\n", len(ctx.Files))
	}

	registry := checks.NewRegistry()

	// Filter checks if config disables any
	var activeChecks []checks.Check
	if fileCfg != nil && len(fileCfg.Disable) > 0 {
		activeChecks = registry.Filter(nil, fileCfg.Disable)
	} else {
		activeChecks = registry.All()
	}

	results := checks.Run(activeChecks, ctx)

	r := scoring.Score(results, ctx.RepoPath, ctx.Languages, len(ctx.Files), startTime)
	r.Version = Version

	if elapsed := time.Since(startTime); elapsed > 3*time.Second {
		fmt.Fprintf(os.Stderr, "Warning: analysis took %.1fs (target: <3s). Consider adding excludes to .repohealthrc.yaml\n", elapsed.Seconds())
	}

	// Use config threshold if --threshold not explicitly set
	if fileCfg != nil && fileCfg.Threshold > 0 && !cmd.Flags().Changed("threshold") {
		threshold = fileCfg.Threshold
	}

	if scoreOnly {
		fmt.Fprintf(os.Stdout, "%d/%d (%s)\n", r.Score, r.MaxScore, r.Grade)
		if ciMode && r.Score < threshold {
			os.Exit(2)
		}
		return nil
	}

	switch format {
	case "json":
		if err := report.JSON(os.Stdout, r); err != nil {
			return err
		}
	case "markdown":
		report.Markdown(os.Stdout, r)
	case "html":
		report.HTML(os.Stdout, r)
	case "terminal":
		report.Terminal(os.Stdout, r, Version)
	default:
		return fmt.Errorf("unknown format: %s", format)
	}

	if ciMode && r.Score < threshold {
		os.Exit(2)
	}

	return nil
}
