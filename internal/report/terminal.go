package report

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/spbuilds/repohealth/internal/model"
)

// Terminal renders a colored terminal report.
func Terminal(w io.Writer, r *model.Report, version string) {
	bold := color.New(color.Bold)
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)
	red := color.New(color.FgRed)
	dim := color.New(color.Faint)

	fmt.Fprintln(w)
	bold.Fprintf(w, "  RepoHealth %s\n", version)
	fmt.Fprintln(w)
	fmt.Fprintf(w, "  Repository: %s\n", r.RepoPath)

	// Languages
	if len(r.Languages) > 0 {
		fmt.Fprintf(w, "  Languages:  %s\n", formatLanguages(r.Languages))
	}

	fmt.Fprintf(w, "  Analyzed:   %d files in %s\n", r.FilesAnalyzed, formatDuration(r.DurationMs))
	fmt.Fprintln(w)

	// Separator
	dim.Fprintln(w, "  "+strings.Repeat("\u2500", 46))
	fmt.Fprintln(w)

	// Overall score
	gradeColor := green
	if r.Score < 70 {
		gradeColor = yellow
	}
	if r.Score < 55 {
		gradeColor = red
	}

	fmt.Fprint(w, "  Overall Score:  ")
	gradeColor.Fprintf(w, "%d / %d", r.Score, r.MaxScore)
	fmt.Fprint(w, "    Grade: ")
	gradeColor.Fprintln(w, r.Grade)

	fmt.Fprintln(w)
	dim.Fprintln(w, "  "+strings.Repeat("\u2500", 46))
	fmt.Fprintln(w)

	// Categories
	for _, cat := range r.Categories {
		catColor := green
		pct := 0
		if cat.MaxScore > 0 {
			pct = cat.Score * 100 / cat.MaxScore
		}
		if pct < 70 {
			catColor = yellow
		}
		if pct < 50 {
			catColor = red
		}

		label := fmt.Sprintf("%-36s", cat.Label)
		score := fmt.Sprintf("%d / %d", cat.Score, cat.MaxScore)
		fmt.Fprintf(w, "  %s", label)
		catColor.Fprintln(w, score)

		// Show individual checks for this category
		for _, check := range r.Checks {
			if check.Category != cat.Name {
				continue
			}

			status := statusIcon(check.Status)
			detail := ""
			if check.Details != "" {
				detail = "  " + check.Details
			}

			fmt.Fprintf(w, "    %-34s %s%s\n", check.Name, status, dim.Sprint(detail))
		}
		fmt.Fprintln(w)
	}

	// Suggestions
	if len(r.Suggestions) > 0 {
		dim.Fprintln(w, "  "+strings.Repeat("\u2500", 46))
		fmt.Fprintln(w)
		bold.Fprintln(w, "  Suggestions (sorted by impact)")

		for _, s := range r.Suggestions {
			pts := fmt.Sprintf("+%d pts", s.Impact)
			if s.Impact == 1 {
				pts = "+1 pt "
			}
			yellow.Fprintf(w, "    %s", pts)
			fmt.Fprintf(w, "  %s\n", s.Message)
		}
		fmt.Fprintln(w)
	}

	// Improvement plan
	if len(r.Suggestions) > 0 {
		dim.Fprintln(w, "  "+strings.Repeat("\u2500", 46))
		fmt.Fprintln(w)
		bold.Fprintln(w, "  Improvement Plan")
		projected := r.Score
		for i, s := range r.Suggestions {
			if i >= 5 {
				break
			}
			prev := projected
			if r.RawMax > 0 {
				projected += s.Impact * 100 / r.RawMax
			} else {
				projected += s.Impact
			}
			if projected > 100 {
				projected = 100
			}
			green.Fprintf(w, "    %d \u2192 %d", prev, projected)
			fmt.Fprintf(w, "  %s\n", s.Message)
		}
		fmt.Fprintln(w)
	}
}

func statusIcon(s model.Status) string {
	switch s {
	case model.StatusFull:
		return color.GreenString("\u2713")
	case model.StatusPartial:
		return color.YellowString("\u25D0")
	case model.StatusNone:
		return color.RedString("\u2717")
	case model.StatusSkipped:
		return color.HiBlackString("-")
	default:
		return "?"
	}
}

func formatLanguages(langs map[string]int) string {
	type langCount struct {
		name  string
		count int
	}

	total := 0
	var sorted []langCount
	for name, count := range langs {
		// Skip non-code languages
		if name == "Markdown" || name == "YAML" || name == "JSON" || name == "TOML" {
			continue
		}
		sorted = append(sorted, langCount{name, count})
		total += count
	}

	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].count != sorted[j].count {
			return sorted[i].count > sorted[j].count
		}
		return sorted[i].name < sorted[j].name
	})

	if len(sorted) == 0 {
		return "(no source files)"
	}

	var parts []string
	shown := 0
	for _, lc := range sorted {
		if shown >= 4 {
			break
		}
		pct := lc.count * 100 / total
		if pct < 1 {
			pct = 1
		}
		parts = append(parts, fmt.Sprintf("%s (%d%%)", lc.name, pct))
		shown++
	}

	return strings.Join(parts, ", ")
}

func formatDuration(ms int64) string {
	if ms < 1000 {
		return fmt.Sprintf("%dms", ms)
	}
	return fmt.Sprintf("%.1fs", float64(ms)/1000.0)
}
