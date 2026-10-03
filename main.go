package main

import (
	"errors"
	"flag"
	"fmt"
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

	treeMap, err := tree.ExploreTree(rootDir)
	for v, _ := range treeMap {
		fmt.Println(v)
		fmt.Println("\t", treeMap[v].ModTime)
	}
	if err != nil {
		return err
	}

	return nil

}
