package main

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func createFakeTree(t *testing.T, structure []string) string {
	t.Helper()
	root := t.TempDir()
	for _, path := range structure {
		if strings.HasSuffix(path, "/") {
			if err := os.MkdirAll(filepath.Join(root, path), 0755); err != nil {
				t.Fatalf("can't create dir %q: %v", path, err)
			}
		} else {
			if err := os.MkdirAll(filepath.Join(root, filepath.Dir(path)), 0755); err != nil {
				t.Fatalf("can't create directory %q for file %q: %v", filepath.Dir(path), path, err)
			}
			if err := os.WriteFile(filepath.Join(root, path), []byte{}, 0644); err != nil {
				t.Fatalf("can't create file %q: %v", path, err)
			}
		}
	}
	return root
}
func TestCheckDir(t *testing.T) {
	data := []struct {
		name     string
		tree     []string
		path     string
		expected error
	}{
		{"nested dir", []string{"A/B/C/"}, "A/B/C/", nil},
		{"sibling nested dir", []string{"A/B/D/"}, "A/B/D/", nil},
		{"top-level dir", []string{"E/"}, "E/", nil},
		{"file with extension", []string{"A/B/f.png"}, "A/B/f.png", errNotDir},
		{"file without extension", []string{"A/Z"}, "A/Z", errNotDir},
		{"hidden file in hidden dir", []string{".G/.f"}, ".G/.f", errNotDir},
		{"missing nested path", nil, "A/M", fs.ErrNotExist},
		{"missing dir", nil, "M/", fs.ErrNotExist},
		{"missing hidden dir", nil, ".M/", fs.ErrNotExist},
		{"missing child of existing dir", []string{"A/"}, "A/M", fs.ErrNotExist},
		{"path through a file", []string{"A/Z"}, "A/Z/f", fs.ErrNotExist},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			rootDir := createFakeTree(t, d.tree)
			err := checkDir(filepath.Join(rootDir, d.path))
			if !errors.Is(err, d.expected) {
				t.Errorf("Expected %v, got %v", d.expected, err)
			}
		})
	}
}
