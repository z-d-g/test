package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsMarkdownFile(t *testing.T) {
	tests := []struct {
		name string
		file string
		want bool
	}{
		{"md file", "readme.md", true},
		{"txt file", "readme.txt", false},
		{"no ext", "readme", false},
		{"path", "/path/to/file.md", true},
		{"empty", "", false},
	}
	for _, tt := range tests {
		if got := IsMarkdownFile(tt.file); got != tt.want {
			t.Errorf("IsMarkdownFile(%q) = %v, want %v", tt.file, got, tt.want)
		}
	}
}

func TestFilterMarkdownFiles(t *testing.T) {
	input := []string{"a.md", "b.txt", "c.md", "d"}
	got := FilterMarkdownFiles(input)
	if len(got) != 2 {
		t.Fatalf("expected 2 files, got %d: %v", len(got), got)
	}
	if got[0] != "a.md" || got[1] != "c.md" {
		t.Errorf("got %v, want [a.md c.md]", got)
	}
}

func TestFilterMarkdownFilesEmpty(t *testing.T) {
	got := FilterMarkdownFiles(nil)
	if len(got) != 0 {
		t.Errorf("expected empty, got %v", got)
	}
}

func TestReadFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.md")
	content := []byte("# Hello\n")
	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatal(err)
	}

	got, err := ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if string(got) != string(content) {
		t.Errorf("ReadFile() = %q, want %q", got, content)
	}
}

func TestReadFileNotFound(t *testing.T) {
	_, err := ReadFile("/nonexistent/file.md")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestWriteFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.md")
	content := []byte("# Written\n")

	if err := WriteFile(path, content); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(content) {
		t.Errorf("file content = %q, want %q", got, content)
	}
}
