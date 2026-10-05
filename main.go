package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"

	"github.com/box-bit/shoebox/internal/tree"
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

func checkDir(path string) error {

	inf, err := os.Stat(path)
	if errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("path \"%s\" doesn't exists", path)
	}
	if err != nil {
		return fmt.Errorf("open directory: %w", err)
	}

	if !inf.IsDir() {
		return fmt.Errorf("path \"%s\" is not a directory", path)
	}

	return nil
}

func handleDiff(backup string, target string) error {
	if err := checkDir(backup); err != nil {
		return fmt.Errorf("backup folder: %w", err)
	}

	if err := checkDir(target); err != nil {
		return fmt.Errorf("target folder: %w", err)
	}

	tree.CompareTrees(backup, target)
	return nil

}

func run() error {
	flag.Usage = func() {
		fmt.Fprintln(flag.CommandLine.Output(), "usage: shoebox [flags] <backup-dir> <root-dir>")
		flag.PrintDefaults()
	}
	diffFlag := flag.Bool("diff", false, "shows the diff between the backup and the target folder")
	var backupRoot string
	var targetRoot string
	flag.StringVar(&backupRoot, "backup", "", "root directory of the photo database")
	flag.StringVar(&targetRoot, "target", "", "root directory of the photo you want to push with the backup")
	flag.Parse()

	if *diffFlag {
		if backupRoot != "" && targetRoot != "" {
			if flag.NArg() > 0 {
				return fmt.Errorf("%w: unexpected extra argument", errUsage)
			}
		} else if backupRoot == "" && targetRoot != "" {
			if flag.NArg() > 1 {
				return fmt.Errorf("%w: unexpected extra argument", errUsage)
			}
			backupRoot = flag.Arg(0)
		} else if targetRoot == "" && backupRoot != "" {
			if flag.NArg() > 1 {
				return fmt.Errorf("%w: unexpected extra argument", errUsage)
			}
			targetRoot = flag.Arg(0)
		} else {
			if flag.NArg() > 2 {
				return fmt.Errorf("%w: unexpected extra argument", errUsage)
			}
			if flag.NArg() < 2 {
				return fmt.Errorf("%w: too few arguments", errUsage)
			}
			backupRoot, targetRoot = flag.Arg(0), flag.Arg(1)

		}

		if err := handleDiff(backupRoot, targetRoot); err != nil {
			return err
		}
	}
	return nil
}
