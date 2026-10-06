package tree

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/fs"
	"maps"
	"path/filepath"
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

// check which files are aligned with the backup and which not
// returns 3 slices
// 1) both present
// 2) only in backup (good)
// 3) only in target (needs to be sync)
func CompareTrees(backup string, target string) error {
	backupScan, err := scanTree(backup)
	if err != nil {
		return err
	}
	targetScan, err := scanTree(target)
	if err != nil {
		return err
	}

	var onlyInBackup = map[string]bool{}
	var onlyInTarget = map[string]bool{}
	// compare by size first, hash only when sizes match
	var bothPresent = map[string]bool{}

	for k := range maps.Keys(targetScan) {
		if _, ok := backupScan[k]; !ok {
			onlyInTarget[k] = true
		} else {
			if backupScan[k].size == targetScan[k].size {

			}
			bothPresent[k] = true
		}
	}
	for k := range maps.Keys(backupScan) {
		if _, ok := targetScan[k]; ok != true {
			onlyInBackup[k] = true
		}
	}
	fmt.Println("only in backup")
	for k := range maps.Keys(onlyInBackup) {
		fmt.Println(k)
	}
	fmt.Println("\nonly in target")
	for k := range maps.Keys(onlyInTarget) {
		fmt.Println(k)
	}
	fmt.Println("\nboth present")
	for k := range maps.Keys(bothPresent) {
		fmt.Println(k)
	}
	return nil
}
