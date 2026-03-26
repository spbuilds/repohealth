package scanner

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"sync"
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

// shortlogCache avoids running git shortlog twice (ContributorCount + BusFactor).
var shortlogCache struct {
	mu           sync.Mutex
	repoPath     string
	contributors int
	busFactor    int
	err          error
	done         bool
}

func shortlogStats(repoPath string) (int, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), gitTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "git", "shortlog", "-sn", "--no-merges", "HEAD")
	cmd.Dir = repoPath
	out, err := cmd.Output()
	if err != nil {
		return 0, 0, err
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return 0, 0, nil
	}

	contributors := len(lines)
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

	busFactor := 0
	if totalCommits > 0 {
		thresholdF := float64(totalCommits) * 0.10
		for _, c := range counts {
			if float64(c) > thresholdF {
				busFactor++
			}
		}
	}

	return contributors, busFactor, nil
}

func cachedShortlog(repoPath string) (int, int, error) {
	shortlogCache.mu.Lock()
	if shortlogCache.done && shortlogCache.repoPath == repoPath {
		c, b, e := shortlogCache.contributors, shortlogCache.busFactor, shortlogCache.err
		shortlogCache.mu.Unlock()
		return c, b, e
	}
	shortlogCache.mu.Unlock()

	c, b, err := shortlogStats(repoPath)

	shortlogCache.mu.Lock()
	shortlogCache.repoPath = repoPath
	shortlogCache.contributors = c
	shortlogCache.busFactor = b
	shortlogCache.err = err
	shortlogCache.done = true
	shortlogCache.mu.Unlock()

	return c, b, err
}

// ContributorCount returns the number of unique commit authors.
func ContributorCount(repoPath string) (int, error) {
	c, _, err := cachedShortlog(repoPath)
	return c, err
}

// CommitCountSince returns the number of commits in the last N months.
func CommitCountSince(repoPath string, months int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), gitTimeout)
	defer cancel()

	since := fmt.Sprintf("%d months ago", months)
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

// FileLastCommitDate returns the date of the most recent commit that touched filePath.
func FileLastCommitDate(repoPath, filePath string) (time.Time, error) {
	ctx, cancel := context.WithTimeout(context.Background(), gitTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "git", "log", "-1", "--format=%cI", "--", filePath)
	cmd.Dir = repoPath
	out, err := cmd.Output()
	if err != nil {
		return time.Time{}, err
	}

	return time.Parse(time.RFC3339, strings.TrimSpace(string(out)))
}

// BusFactor returns the number of contributors with > 10% of total commits.
func BusFactor(repoPath string) (int, error) {
	_, b, err := cachedShortlog(repoPath)
	return b, err
}
