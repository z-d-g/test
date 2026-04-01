package cursor

import (
	"os"
	"path/filepath"
	"testing"
)

func testStore(t *testing.T) *PositionStore {
	t.Helper()
	dir := t.TempDir()
	return &PositionStore{cacheDir: dir}
}

func TestSetAndGetPosition(t *testing.T) {
	ps := testStore(t)

	cfg := FileConfig{CursorLine: 5, CursorCol: 10}
	if err := ps.SetPosition("/path/to/file.md", cfg); err != nil {
		t.Fatalf("SetPosition() error = %v", err)
	}

	got, ok := ps.GetPosition("/path/to/file.md")
	if !ok {
		t.Fatal("GetPosition() returned false")
	}
	if got.CursorLine != 5 || got.CursorCol != 10 {
		t.Errorf("got %+v, want {CursorLine:5 CursorCol:10}", got)
	}
}

func TestGetPositionNotFound(t *testing.T) {
	ps := testStore(t)
	_, ok := ps.GetPosition("/nonexistent.md")
	if ok {
		t.Error("GetPosition() should return false for nonexistent file")
	}
}

func TestRemovePosition(t *testing.T) {
	ps := testStore(t)

	cfg := FileConfig{CursorLine: 1, CursorCol: 0}
	if err := ps.SetPosition("/file.md", cfg); err != nil {
		t.Fatal(err)
	}

	if err := ps.RemovePosition("/file.md"); err != nil {
		t.Fatalf("RemovePosition() error = %v", err)
	}

	_, ok := ps.GetPosition("/file.md")
	if ok {
		t.Error("GetPosition() should return false after remove")
	}
}

func TestConfigPath(t *testing.T) {
	ps := &PositionStore{cacheDir: "/tmp/test"}
	got := ps.configPath("/path/to/file.md")
	want := filepath.Join("/tmp/test", "_path_to_file.md.json")
	if got != want {
		t.Errorf("configPath() = %q, want %q", got, want)
	}
}

func TestSanitize(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"/path/to/file.md", "_path_to_file.md"},
		{"C:\\Users\\file.md", "C__Users_file.md"},
		{".hidden", "_hidden"},
		{"normal.md", "normal.md"},
	}
	for _, tt := range tests {
		if got := sanitize(tt.path); got != tt.want {
			t.Errorf("sanitize(%q) = %q, want %q", tt.path, got, tt.want)
		}
	}
}

func TestGetPositionCorruptJSON(t *testing.T) {
	ps := testStore(t)

	// Write invalid JSON
	path := ps.configPath("/bad.json")
	os.WriteFile(path, []byte("not json"), 0644)

	_, ok := ps.GetPosition("/bad.json")
	if ok {
		t.Error("GetPosition() should return false for corrupt JSON")
	}
}

func TestNewPositionStore(t *testing.T) {
	ps, err := NewPositionStore()
	if err != nil {
		t.Fatalf("NewPositionStore() error = %v", err)
	}
	if ps == nil {
		t.Fatal("NewPositionStore() returned nil")
	}
}
