package testfs

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// directory paths must end with /
func Tree(t *testing.T, structure []string) string {
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

// returns filepath of new created file
func CreateFile(t *testing.T, path string, content string) string {
	t.Helper()
	root := Tree(t, []string{path})
	joined := filepath.Join(root, path)
	if err := os.WriteFile(joined, []byte(content), 0644); err != nil {
		t.Fatalf("can't write into file %q: %v", joined, err)
	}
	return joined
}
