package scanner

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

const gitTimeout = 5 * time.Second

// LastCommitDate returns the date of the most recent commit.
func LastCommitDate(repoPath string) (time.Time, error) {
	ctx, cancel := context.WithTimeout(context.Background(), gitTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "git", "log", "-1", "--format=%cI")
	cmd.Dir = repoPath
	out, err := cmd.Output()
	if err != nil {
		return time.Time{}, err
	}

	return time.Parse(time.RFC3339, strings.TrimSpace(string(out)))
}

// ContributorCount returns the number of unique commit authors.
func ContributorCount(repoPath string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), gitTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "git", "shortlog", "-sn", "--no-merges", "HEAD")
	cmd.Dir = repoPath
	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return 0, nil
	}

	return len(lines), nil
}

// CommitCountSince returns the number of commits in the last N months.
func CommitCountSince(repoPath string, months int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), gitTimeout)
	defer cancel()

	since := strings.TrimSpace(strings.Replace("N months ago", "N", fmt.Sprintf("%d", months), 1))
	cmd := exec.CommandContext(ctx, "git", "log", "--oneline", "--since="+since)
	cmd.Dir = repoPath
	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	text := strings.TrimSpace(string(out))
	if text == "" {
		return 0, nil
	}

	return len(strings.Split(text, "\n")), nil
}

// TagCount returns the number of git tags.
func TagCount(repoPath string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), gitTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "git", "tag", "-l")
	cmd.Dir = repoPath
	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	text := strings.TrimSpace(string(out))
	if text == "" {
		return 0, nil
	}

	return len(strings.Split(text, "\n")), nil
}

// BusFactor returns the number of contributors with > 10% of total commits.
func BusFactor(repoPath string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), gitTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "git", "shortlog", "-sn", "--no-merges", "HEAD")
	cmd.Dir = repoPath
	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) == 0 || (len(lines) == 1 && lines[0] == "") {
		return 0, nil
	}

	// Parse commit counts
	totalCommits := 0
	var counts []int
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "\t", 2)
		if len(parts) < 1 {
			continue
		}
		n := 0
		for _, c := range strings.TrimSpace(parts[0]) {
			if c >= '0' && c <= '9' {
				n = n*10 + int(c-'0')
			}
		}
		counts = append(counts, n)
		totalCommits += n
	}

	if totalCommits == 0 {
		return 0, nil
	}

	threshold := totalCommits / 10 // 10%
	busFactor := 0
	for _, c := range counts {
		if c > threshold {
			busFactor++
		}
	}

	return busFactor, nil
}
