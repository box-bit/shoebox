package tree

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/box-bit/shoebox/internal/testfs"
)

func TestEncodeFile(t *testing.T) {
	data := []struct {
		name    string
		file    string
		content string
		want    string
	}{
		{"empty file", "empty.txt", "", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
		{"file with content", "file.txt", "hello", "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			f := testfs.CreateFile(t, d.file, d.content)
			file, err := os.Open(f)
			if err != nil {
				t.Fatalf("can't open file %q: %v", f, err)
			}
			defer file.Close()
			encoded, err := encodeFile(file)
			if err != nil {
				t.Fatalf("error occurred: %v", err)
			}

			if fmt.Sprintf("%x", encoded) != d.want {
				t.Errorf("encodeFile(%q) = %x, wants %x", f, d.want, fmt.Sprintf("%x", encoded))
			}
		})
	}
}

func TestCompareTrees(t *testing.T) {
	data := []struct {
		name       string
		backupTree map[string]string
		targetTree map[string]string
		want       Diff
	}{
		{
			name:       "both empty",
			backupTree: map[string]string{},
			targetTree: map[string]string{},
			want:       Diff{},
		},
		{
			name:       "identical trees",
			backupTree: map[string]string{"a.jpg": "aaa", "2024/b.jpg": "bbb"},
			targetTree: map[string]string{"a.jpg": "aaa", "2024/b.jpg": "bbb"},
			want:       Diff{Same: []string{"2024/b.jpg", "a.jpg"}},
		},
		{
			name:       "only in backup",
			backupTree: map[string]string{"a.jpg": "aaa", "old.jpg": "old"},
			targetTree: map[string]string{"a.jpg": "aaa"},
			want:       Diff{Same: []string{"a.jpg"}, OnlyInBackup: []string{"old.jpg"}},
		},
		{
			name:       "only in target",
			backupTree: map[string]string{"a.jpg": "aaa"},
			targetTree: map[string]string{"a.jpg": "aaa", "new.jpg": "new"},
			want:       Diff{Same: []string{"a.jpg"}, OnlyInTarget: []string{"new.jpg"}},
		},
		{
			name:       "changed with different size",
			backupTree: map[string]string{"a.jpg": "short"},
			targetTree: map[string]string{"a.jpg": "much longer"},
			want:       Diff{Changed: []string{"a.jpg"}},
		},
		{
			name:       "changed with same size",
			backupTree: map[string]string{"a.jpg": "abc"},
			targetTree: map[string]string{"a.jpg": "xyz"},
			want:       Diff{Changed: []string{"a.jpg"}},
		},
		{
			name:       "empty directories are ignored",
			backupTree: map[string]string{"empty/": ""},
			targetTree: map[string]string{"other/": ""},
			want:       Diff{},
		},
		{
			name: "mixed",
			backupTree: map[string]string{
				"same.jpg":     "same",
				"changed.jpg":  "v1",
				"deleted.jpg":  "gone",
				"sub/keep.jpg": "keep",
			},
			targetTree: map[string]string{
				"same.jpg":     "same",
				"changed.jpg":  "v2",
				"added.jpg":    "new",
				"sub/keep.jpg": "keep",
			},
			want: Diff{
				Same:         []string{"same.jpg", "sub/keep.jpg"},
				Changed:      []string{"changed.jpg"},
				OnlyInBackup: []string{"deleted.jpg"},
				OnlyInTarget: []string{"added.jpg"},
			},
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			backup := testfs.TreeWithContent(t, d.backupTree)
			target := testfs.TreeWithContent(t, d.targetTree)
			got, err := CompareTrees(backup, target)
			if err != nil {
				t.Fatalf("CompareTrees() error: %v", err)
			}
			if !reflect.DeepEqual(got, d.want) {
				t.Errorf("CompareTrees() = %+v, want %+v", got, d.want)
			}
		})
	}
}

func TestCompareTreesMissingDir(t *testing.T) {
	existing := t.TempDir()
	missing := filepath.Join(t.TempDir(), "does-not-exist")

	if _, err := CompareTrees(missing, existing); err == nil {
		t.Error("CompareTrees() with missing backup dir: expected error, got nil")
	}
	if _, err := CompareTrees(existing, missing); err == nil {
		t.Error("CompareTrees() with missing target dir: expected error, got nil")
	}
}
