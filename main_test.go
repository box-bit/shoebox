package main

import (
	"errors"
	"io/fs"
	"path/filepath"
	"testing"

	"github.com/box-bit/shoebox/internal/testfs"
)

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
			rootDir := testfs.Tree(t, d.tree)
			err := checkDir(filepath.Join(rootDir, d.path))
			if !errors.Is(err, d.expected) {
				t.Errorf("Expected %v, got %v", d.expected, err)
			}
		})
	}
}
