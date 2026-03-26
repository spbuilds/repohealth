package scanner

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
)

const maxLinesPerFile = 10000
const maxFileSize = 100 * 1024 // 100KB

// ReadFileLines reads up to maxLinesPerFile lines from a file.
// Returns nil if the file is binary (contains null bytes in first 512 bytes)
// or exceeds maxFileSize.
func ReadFileLines(repoPath, relPath string) ([]string, error) {
	fullPath := filepath.Join(repoPath, relPath)

	info, err := os.Stat(fullPath)
	if err != nil {
		return nil, err
	}
	if info.Size() > maxFileSize {
		return nil, nil // skip large files
	}

	f, err := os.Open(fullPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Check for binary file (null bytes in first 512 bytes)
	header := make([]byte, 512)
	n, err := f.Read(header)
	if n == 0 {
		if err == io.EOF || err == nil {
			return []string{}, nil // empty file, not binary
		}
		return nil, err
	}
	for i := 0; i < n; i++ {
		if header[i] == 0 {
			return nil, nil // binary file
		}
	}

	// Reset to beginning
	if _, err := f.Seek(0, 0); err != nil {
		return nil, err
	}

	var lines []string
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 256*1024), 256*1024) // 256KB max line length
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
		if len(lines) >= maxLinesPerFile {
			break
		}
	}

	return lines, scanner.Err()
}
