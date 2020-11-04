package tree

import (
	"fmt"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

type CheckerManager struct {
	tree formula.TreeManager
}

func NewChecker(
	tree formula.TreeManager,
) CheckerManager {
	return CheckerManager{
		tree: tree,
	}
}

// CheckCommands is used to warn the user about conflicting
// formula commands on different repos. This function doesn't
// return an error because printing an error from a unsuccessful
// warning attempt can be confusing to the user.
func (cm CheckerManager) Check() {
	commands := cm.filterCommands()
	conflictingCommands := cm.conflictingCommands(commands)
	if len(conflictingCommands) > 1 {
		cm.printConflictingCommandsWarning(conflictingCommands)
	}
}

func (cm CheckerManager) filterCommands() []string {
	allCommands := []string{""}
	tree, _ := cm.tree.Tree()

	for _, t := range tree {
		for _, c := range t.Commands {
			allCommands = append(allCommands, c.Id)
		}
	}
	return allCommands
}

func (cm CheckerManager) conflictingCommands(commands []string) []string {
	duplicateFrequency := make(map[string]int)
	duplicatedCommands := []string{""}
	for _, item := range commands {
		_, exist := duplicateFrequency[item]
		if exist {
			duplicateFrequency[item] += 1 // increase counter by 1 if already in the map
			duplicatedCommands = append(duplicatedCommands, item)
		} else {
			duplicateFrequency[item] = 1 // else start counting from 1
		}
	}
	return duplicatedCommands
}

func (cm CheckerManager) printConflictingCommandsWarning(conflictingCommands []string) {
	lastCommandIndex := len(conflictingCommands) - 1
	lastCommand := conflictingCommands[lastCommandIndex]
	lastCommand = strings.Replace(lastCommand, "root", "rit", 1)
	lastCommand = strings.ReplaceAll(lastCommand, "_", " ")
	msg := fmt.Sprintf("The following formula command are conflicting: %s", lastCommand)
	msg = prompt.Yellow(msg)
	fmt.Println(msg)
}
