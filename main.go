package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"

	"github.com/box-bit/shoebox/internal/tree"
)

var errUsage = errors.New("usage")

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "shoebox:", err)
		if errors.Is(err, errUsage) {
			flag.Usage()
			os.Exit(2)
		}
		os.Exit(1)
	}
}

func checkDir(path string) error {

	info, err := os.Stat(path)
	if errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("%q: %w", path, fs.ErrNotExist)
	}
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return fmt.Errorf("%q is not a directory", path)
	}

	return nil
}

func elaborateInput(diff bool) (string, string, error) {
	if diff {
		if flag.NArg() != 2 {
			return "", "", fmt.Errorf("%w: expected 2 directories, got %d", errUsage, flag.NArg())
		}

		backup, target := flag.Arg(0), flag.Arg(1)

		if err := checkDir(backup); err != nil {
			return "", "", fmt.Errorf("backup directory: %w", err)
		}

		if err := checkDir(target); err != nil {
			return "", "", fmt.Errorf("target directory: %w", err)
		}
		return backup, target, nil
	}
	return "", "", fmt.Errorf("%w: no command given", errUsage)
}

func run() error {
	flag.Usage = func() {
		fmt.Fprintln(flag.CommandLine.Output(), "usage: shoebox [flags] <backup-dir> <root-dir>")
		flag.PrintDefaults()
	}
	diffFlag := flag.Bool("diff", false, "shows the diff between the backup and the target folder")

	flag.Parse()
	backup, root, err := elaborateInput(*diffFlag)
	if err != nil {
		return err
	}

	if err := tree.CompareTrees(backup, root); err != nil {
		return fmt.Errorf("compare trees: %w", err)
	}
	return nil
}
