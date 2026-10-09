package tree

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"time"
)

type Entry struct {
	size    int64
	ModTime time.Time
}

func encodeFile(r io.Reader) ([sha256.Size]byte, error) {
	h := sha256.New()
	_, err := io.Copy(h, r)
	if err != nil {
		return [sha256.Size]byte{}, err
	}
	return [sha256.Size]byte(h.Sum(nil)), nil
}

func isSameContent(f1 string, f2 string) (bool, error) {

	of1, err := os.Open(f1)
	if err != nil {
		return false, fmt.Errorf("can't open file %q: %w", f1, err)
	}
	defer of1.Close()
	of2, err := os.Open(f2)
	if err != nil {
		return false, fmt.Errorf("can't open file %q: %w", f2, err)
	}
	defer of2.Close()
	encoded1, err := encodeFile(of1)
	if err != nil {
		return false, fmt.Errorf("encoding file %q: %w", f1, err)
	}
	encoded2, err := encodeFile(of2)
	if err != nil {
		return false, fmt.Errorf("encoding file %q: %w", f2, err)
	}
	if encoded1 == encoded2 {
		return true, nil
	}
	return false, nil
}

func scanTree(root string) (map[string]Entry, error) {

	m := map[string]Entry{}

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && d.Type().IsRegular() {
			info, err := d.Info()
			if err != nil {
				return err
			}
			relPath, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			m[relPath] = Entry{info.Size(), info.ModTime()}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return m, nil
}

type Diff struct {
	Same         []string
	Changed      []string
	OnlyInBackup []string
	OnlyInTarget []string
}

// check which files are aligned with the backup and which not
func CompareTrees(backup string, target string) (Diff, error) {
	backupScan, err := scanTree(backup)
	if err != nil {
		return Diff{}, err
	}
	targetScan, err := scanTree(target)
	if err != nil {
		return Diff{}, err
	}

	var d = Diff{}
	// compare by size first, hash only when sizes match
	// keys are sorted so the result is deterministic

	for _, p := range slices.Sorted(maps.Keys(targetScan)) {
		if _, ok := backupScan[p]; !ok {
			d.OnlyInTarget = append(d.OnlyInTarget, p)
		} else {
			if backupScan[p].size == targetScan[p].size {
				same, err := isSameContent(filepath.Join(backup, p), filepath.Join(target, p))
				if err != nil {
					return Diff{}, err
				}
				if same {
					d.Same = append(d.Same, p)
				} else {
					d.Changed = append(d.Changed, p)
				}
			} else {
				d.Changed = append(d.Changed, p)
			}
		}
	}
	for _, p := range slices.Sorted(maps.Keys(backupScan)) {
		if _, ok := targetScan[p]; ok != true {
			d.OnlyInBackup = append(d.OnlyInBackup, p)
		}
	}
	return d, nil
}
