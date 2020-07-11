package files

import (
	"log"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar"
)

// FindFilesInGit finds all files in git in the given root
// matching the include pattern. optinally you can also find
// files changed in last commit
func FindFilesInGit(rootDirAbs, includePattern string, changedOnly bool) []string {
	// prepare the git command
	var cmd *exec.Cmd
	if changedOnly {
		cmd = exec.Command("git", "--no-pager", "diff", "--name-only", "HEAD..HEAD~10")
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
