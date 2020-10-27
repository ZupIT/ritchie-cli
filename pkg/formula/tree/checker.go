package tree

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

// TODO rit home dir must be a parameter
type CheckerManager struct {
	dir  stream.DirLister
	file stream.FileReader
}

func NewChecker(dir stream.DirLister, file stream.FileReader) CheckerManager {
	return CheckerManager{dir: dir, file: file}
}

// CheckCommands is used to warn the user about conflicting
// formula commands on different repos. This function don't
// return err because print an error because of a unsuccessful
// warning attempt can be confusing to the user.
func (cm CheckerManager) CheckCommands() {
	trees := cm.readCommands()
	commands := filterCommands(trees)
	conflictingCommands := conflictingCommands(commands)
	printConflictingCommandsWarning(conflictingCommands)

}

func (cm CheckerManager) readCommands() []formula.Tree {
	repoDir := filepath.Join(api.RitchieHomeDir(), "repos")
	repos, _ := cm.dir.List(repoDir, false)

	tree := formula.Tree{}
	treeArr := make([]formula.Tree, len(repos))
	for _, r := range repos {
		path := fmt.Sprintf(treeRepoCmdPattern, api.RitchieHomeDir(), r)
		bytes, _ := cm.file.Read(path)
		_ = json.Unmarshal(bytes, &tree)
		treeArr = append(treeArr, tree)
	}
	return treeArr
}

func conflictingCommands(cmds []string) []string {
	moreThanOne := 0
	duplicatedCommands := []string{""}
	for _, c := range cmds {
		for _, c2 := range cmds {
			if c == c2 {
				moreThanOne++
				if moreThanOne == 2 {
					duplicatedCommands = append(duplicatedCommands, c)
				}
			}
		}
	moreThanOne = 0
	}

	return duplicatedCommands
}

func filterCommands(tree []formula.Tree) []string {
	allCommands := []string{""}
	for _, t := range tree {
		for _, c := range t.Commands {
			allCommands = append(allCommands, c.Id)
		}
	}
	return allCommands
}

func printConflictingCommandsWarning(conflictingCommands []string) {
	// TODO improve the warn information
	prompt.Warning("There are some conflicting commands on formulas" +
		"it can cause unexpected behaviors")

	fmt.Println(conflictingCommands)
}
