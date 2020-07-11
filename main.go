package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/unvar/markfuence/files"
)

func main() {
	// get the current directory
	currentDir, _ := os.Getwd()

	// setup the argument parsing
	rootDirAbs := flag.String("root", currentDir, "absolute path to root directory")
	docsDirPattern := flag.String("docs", "**/docs/", "glob pattern for doc directories")
	changedOnly := flag.Bool("changed", false, "perform git check for changed files")
	gitDepth := flag.Int("depth", 1, "git commit depth when looking for changes")

	// parse the arguments
	flag.Parse()

	// find files that have changed since the last commit
	mdFilePattern := filepath.Join(*docsDirPattern, "**", "*.md")
	mdFiles := files.FindFilesInGit(*rootDirAbs, mdFilePattern, *changedOnly, *gitDepth)

	for _, file := range mdFiles {
		println(file)
	}
}
