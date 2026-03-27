package scanner

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadFileLines_Normal(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")
	if err := os.WriteFile(path, []byte("line1\nline2\nline3\n"), 0644); err != nil {
		t.Fatal(err)
	}

	lines, err := ReadFileLines(dir, "test.txt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 3 {
		t.Errorf("got %d lines, want 3", len(lines))
	}
}

func TestReadFileLines_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.txt")
	if err := os.WriteFile(path, []byte(""), 0644); err != nil {
		t.Fatal(err)
	}

	lines, err := ReadFileLines(dir, "empty.txt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lines == nil {
		t.Error("expected empty slice, got nil")
	}
	if len(lines) != 0 {
		t.Errorf("got %d lines, want 0", len(lines))
	}
}

func TestReadFileLines_BinaryFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "binary.bin")
	data := make([]byte, 100)
	data[50] = 0 // null byte = binary
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}

	lines, err := ReadFileLines(dir, "binary.bin")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lines != nil {
		t.Error("expected nil for binary file, got lines")
	}
}

func TestReadFileLines_LargeFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "large.txt")
	// Create file > 100KB
	data := strings.Repeat("x", 200*1024)
	if err := os.WriteFile(path, []byte(data), 0644); err != nil {
		t.Fatal(err)
	}

	lines, err := ReadFileLines(dir, "large.txt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lines != nil {
		t.Error("expected nil for large file, got lines")
	}
}

func TestReadFileLines_NonexistentFile(t *testing.T) {
	_, err := ReadFileLines(t.TempDir(), "nonexistent.txt")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}
