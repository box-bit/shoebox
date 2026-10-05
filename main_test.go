package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func createFakeTree(t *testing.T, structure []string) {
	// Later, your diff will compare modification times, and files you create all get "now". Look up os.Chtimes so your helper can set times you control.
	t.Helper()
	root := t.TempDir()
	for _, path := range structure {
		if strings.HasSuffix(path, "/") {
			if err := os.MkdirAll(filepath.Join(root, path), 0755); err != nil {
				t.Fatalf("can't create dir %q", path)
			}
		} else {
			if err := os.MkdirAll(filepath.Join(root, filepath.Dir(path)), 0755); err != nil {
				t.Fatalf("can't create directory %q for file %q", filepath.Dir(path), path)
			}
			if err := os.WriteFile(filepath.Join(root, path), []byte{}, 0755); err != nil {
				t.Fatalf("can't create file %q", path)
			}
		}
	}
}
func Test_checkDir(t *testing.T) {

}
