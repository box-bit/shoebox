package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"syscall"

	"github.com/box-bit/shoebox/internal/tree"
)

var errUsage = errors.New("usage")
var errNotDir = errors.New("not a directory")

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
	// syscall.ENOTDIR is raised when OS try to resolve a path that is actually a file (ex /A/B/C where B is a file)
	if errors.Is(err, fs.ErrNotExist) || errors.Is(err, syscall.ENOTDIR) {
		return fmt.Errorf("%q: %w", path, fs.ErrNotExist)
	}
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return fmt.Errorf("%q: %w", path, errNotDir)
	}

	return nil
}
func elaborateDiff() error {

	if flag.NArg() != 2 {
		return fmt.Errorf("%w: expected 2 directories, got %d", errUsage, flag.NArg())
	}

	backup, target := flag.Arg(0), flag.Arg(1)

	if err := checkDir(backup); err != nil {
		return fmt.Errorf("backup directory: %w", err)
	}

	if err := checkDir(target); err != nil {
		return fmt.Errorf("target directory: %w", err)
	}

	diff, err := tree.CompareTrees(backup, target)
	if err != nil {
		return fmt.Errorf("compare trees: %w", err)
	}

	fmt.Println("Differences found")
	fmt.Println("\nChanged files:")
	for _, p := range diff.Changed {
		fmt.Printf("\t %q\n", p)
	}
	fmt.Println("\nNew files to sync:")
	for _, p := range diff.OnlyInTarget {
		fmt.Printf("\t %q\n", p)
	}
	return nil
}

func elaborateDupl() error {

	if flag.NArg() != 1 {
		return fmt.Errorf("%w: expected 1 directory, got %d", errUsage, flag.NArg())
	}

	target := flag.Arg(0)

	if err := checkDir(target); err != nil {
		return fmt.Errorf("target directory: %w", err)
	}

	dupl, err := tree.FindDuplicates(target)
	if err != nil {
		return fmt.Errorf("find duplicates: %w", err)
	}

	fmt.Println("duplicates found")
	fmt.Println(dupl)
	return nil
}

func run() error {
	flag.Usage = func() {
		fmt.Fprintln(flag.CommandLine.Output(), "usage: shoebox [flags] <backup-dir> <target-dir>")
		flag.PrintDefaults()
	}
	diffFlag := flag.Bool("diff", false, "shows the diff between the backup and the target folder")
	duplFlag := flag.Bool("dupl", false, "shows the duplicates found in the given directory")

	flag.Parse()
	switch {
	case *diffFlag:
		if err := elaborateDiff(); err != nil {
			return err
		}
	case *duplFlag:
		if err := elaborateDupl(); err != nil {
			return err
		}
	default:
		return fmt.Errorf("%w: no command given", errUsage)
	}

	return nil
}
