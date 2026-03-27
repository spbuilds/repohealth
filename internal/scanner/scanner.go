// Package scanner walks repository directories, detects languages, and collects file metadata.
package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spbuilds/repohealth/internal/model"
)

const maxFilesToScan = 100000

// skipDirs are directories that should never be scanned.
// These are build outputs, caches, dependency dirs, and IDE configs
// that do not represent the repository's own source or documentation.
var skipDirs = map[string]bool{
	// Version control
	".git": true,
	// Dependencies
	"node_modules": true,
	"vendor":       true,
	".venv":        true,
	// Build output
	"dist":        true,
	"build":       true,
	"target":      true,
	"out":         true,
	".next":       true,
	".nuxt":       true,
	".output":     true,
	".svelte-kit": true,
	// Caches
	"__pycache__":   true,
	".tox":          true,
	".mypy_cache":   true,
	".pytest_cache": true,
	".cache":        true,
	".turbo":        true,
	// Coverage output
	"coverage":    true,
	".nyc_output": true,
	// IDE
	".idea":   true,
	".vscode": true,
}

// languageExtensions maps file extensions to language names.
var languageExtensions = map[string]string{
	".go":    "Go",
	".py":    "Python",
	".js":    "JavaScript",
	".ts":    "TypeScript",
	".tsx":   "TypeScript",
	".jsx":   "JavaScript",
	".rs":    "Rust",
	".java":  "Java",
	".rb":    "Ruby",
	".php":   "PHP",
	".c":     "C",
	".cpp":   "C++",
	".h":     "C",
	".cs":    "C#",
	".swift": "Swift",
	".kt":    "Kotlin",
	".sh":    "Shell",
	".bash":  "Shell",
	".zsh":   "Shell",
	".yml":   "YAML",
	".yaml":  "YAML",
	".json":  "JSON",
	".toml":  "TOML",
	".md":    "Markdown",
	".html":  "HTML",
	".css":   "CSS",
	".scss":  "SCSS",
	".sql":   "SQL",
	".r":     "R",
	".lua":   "Lua",
	".dart":  "Dart",
	".ex":    "Elixir",
	".exs":   "Elixir",
}

// Scan walks the repository directory and collects metadata.
func Scan(repoPath string, excludes []string) (*model.ScanContext, error) {
	absPath, err := filepath.Abs(repoPath)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(absPath)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("path is not a directory: %s", absPath)
	}

	ctx := &model.ScanContext{
		RepoPath:  absPath,
		Files:     make([]model.FileInfo, 0, 256),
		Dirs:      make([]string, 0, 32),
		Languages: make(map[string]int),
	}

	// Check for .git directory
	if _, err := os.Stat(filepath.Join(absPath, ".git")); err == nil {
		ctx.GitAvailable = true
	}

	excludeSet := make(map[string]bool)
	for _, e := range excludes {
		excludeSet[e] = true
	}

	err = filepath.Walk(absPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // skip unreadable files/dirs gracefully
		}
		if info == nil {
			return nil // skip broken symlinks
		}

		relPath, relErr := filepath.Rel(absPath, path)
		if relErr != nil {
			return nil
		}
		if relPath == "." {
			return nil
		}

		// Normalize to forward slashes for cross-platform consistency
		relPath = filepath.ToSlash(relPath)
		name := info.Name()

		if info.IsDir() {
			if skipDirs[name] || excludeSet[name] || excludeSet[relPath] {
				return filepath.SkipDir
			}
			ctx.Dirs = append(ctx.Dirs, relPath)
			return nil
		}

		// Skip hidden files in subdirectories (allow dotfiles at root like .gitignore)
		if strings.HasPrefix(name, ".") {
			dir := filepath.ToSlash(filepath.Dir(relPath))
			if dir != "." {
				return nil
			}
		}

		fi := model.FileInfo{
			Path:  relPath,
			Name:  name,
			Size:  info.Size(),
			IsDir: false,
		}
		ctx.Files = append(ctx.Files, fi)
		if len(ctx.Files) >= maxFilesToScan {
			ctx.Truncated = true
			return filepath.SkipAll
		}

		// Detect language
		ext := strings.ToLower(filepath.Ext(name))
		if lang, ok := languageExtensions[ext]; ok {
			ctx.Languages[lang]++
		}

		return nil
	})

	return ctx, err
}
