package scanner

import (
	"context"
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
