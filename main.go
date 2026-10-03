package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

var errUsage = errors.New("Usage error")

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		if errors.Is(err, errUsage) {
			flag.Usage()
			os.Exit(2)
		}
		os.Exit(1)
	}
}

func run() error {
	flag.Usage = func() {
		fmt.Fprintln(flag.CommandLine.Output(), "usage: shoebox [flags] <dir>")
		flag.PrintDefaults()
	}
	var rootDir string
	flag.StringVar(&rootDir, "dir", "", "root directory of the photo database")
	flag.Parse()

	if len(rootDir) < 1 {
		if flag.NArg() < 1 {
			return fmt.Errorf("%w: no folder given as input", errUsage)
		}
		if flag.NArg() > 1 {
			return fmt.Errorf("%w: expected 1 argument, %d given", errUsage, flag.NArg())
		}
		rootDir = flag.Arg(0)
	}

	fileInfo, err := os.Stat(rootDir)
	if os.IsNotExist(err) {
		return fmt.Errorf("path \"%s\" doesn't exists", rootDir)
	}
	if err != nil {
		return fmt.Errorf("open directory: %w", err)
	}

	if !fileInfo.IsDir() {
		return fmt.Errorf("path \"%s\" is not a directory", rootDir)
	}

	treeMap, err := exploreTree(rootDir)
	for v, _ := range treeMap {
		fmt.Println(v)
		fmt.Println("\t", treeMap[v].modTime)
	}
	if err != nil {
		return err
	}

	return nil

}

type fileInfo struct {
	name      string
	extension string
	modTime   time.Time
}

func exploreTree(rootFolder string) (map[string]fileInfo, error) {

	m := map[string]fileInfo{}

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
			m[relPath] = fileInfo{info.Name(), filepath.Ext(relPath), info.ModTime()}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return m, nil
}
