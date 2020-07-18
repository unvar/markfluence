package files

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar"

	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

// FindFilesInGit finds all files in git in the given root
// matching the include pattern. optinally you can also find
// files changed in last commit
func FindFilesInGit(rootDirAbs, includePattern string, changedOnly bool, gitDepth int) []string {
	// prepare the git command
	var cmd *exec.Cmd
	if changedOnly {
		cmd = exec.Command("git", "--no-pager", "diff", "--name-only", fmt.Sprintf("HEAD..HEAD~%d", gitDepth))
	} else {
		cmd = exec.Command("git", "ls-tree", "-r", "HEAD", "--name-only")
	}
	cmd.Dir = rootDirAbs

	// run the command
	output, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	// process the output
	files := make([]string, 0)
	for _, file := range strings.Split(string(output), "\n") {
		if ok, _ := doublestar.Match(includePattern, file); ok {
			files = append(files, filepath.Join(rootDirAbs, file))
		}
	}

	return files
}

// ProcessFile processes a markdown file and upserts the
// content to confluence
func ProcessFile(filepath string) {
	convertMarkdownToHTML(filepath)
}

func convertMarkdownToHTML(filepath string) {
	// read the file
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}

	// create the md parser
	markdown := goldmark.New(
		goldmark.WithExtensions(
			meta.Meta,
			extension.GFM,
		),
	)

	// convert the file to html
	var buf bytes.Buffer
	context := parser.NewContext()
	if err := markdown.Convert(data, &buf, parser.WithContext(context)); err != nil {
		panic(err)
	}

	metaData := meta.Get(context)
	fmt.Println(metaData["title"], "|", metaData["space"], "|", metaData["parent"])
}
