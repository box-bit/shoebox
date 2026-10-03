package tree

import (
	"io/fs"
	"path/filepath"
	"time"
)

type Entry struct {
	Name      string
	Extension string
	ModTime   time.Time
}

func ExploreTree(rootFolder string) (map[string]Entry, error) {

	m := map[string]Entry{}

	err := filepath.WalkDir(rootFolder, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && d.Type().IsRegular() {
			info, err := d.Info()
			if err != nil {
				return err
			}
			relPath, err := filepath.Rel(rootFolder, path)
			if err != nil {
				return err
			}
			m[relPath] = Entry{info.Name(), filepath.Ext(relPath), info.ModTime()}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return m, nil
}
