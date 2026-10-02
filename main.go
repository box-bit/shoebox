package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

var errUsage = errors.New("Usage error")

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		if errors.Is(err, errUsage) {
			// called wrong
			fmt.Fprintln(os.Stderr, "Usage:")
			flag.PrintDefaults()
			os.Exit(2)
		}
		os.Exit(1)
	}
}

func run() error {
	flag.Usage = func() {
		fmt.Fprintln(flag.CommandLine.Output(), "usage: photo-tree [flags] <dir>")
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

	// check if directory
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

	err = exploreTree(rootDir)
	if err != nil {
		return err
	}
	return nil

}

func exploreTree(rootFolder string) error {
	filepath.WalkDir(rootFolder, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			fmt.Printf("%s of type directory\n", path)
		} else {
			fmt.Printf("%s with ext %s\n", path, filepath.Ext(path))
			fileInfo, err := os.Stat(path)
			if err != nil {
				return err
			}
			fmt.Printf("Modification time: %v\n\n", fileInfo.ModTime())
		}
		return nil
	})
	return nil
}
