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
	if len(conflictingCommands) >= 1 {
		cm.printConflictingCommandsWarning(conflictingCommands)
	}
}

func (cm CheckerManager) filterCommands() []string {
	var allCommands []string
	tree, _ := cm.tree.Tree()
	for _, t := range tree {
		for _, c := range t.Commands {
			if c.Formula {
				allCommands = append(allCommands, c.Id)
			}
		}
	}
	return allCommands
}

func (cm CheckerManager) conflictingCommands(commands []string) []string {
	duplicateFrequency := make(map[string]int)
	var duplicatedCommands []string
	for _, command := range commands {
		_, exist := duplicateFrequency[command]
		if exist {
			duplicateFrequency[command] += 1
			duplicatedCommands = append(duplicatedCommands, command)
		} else {
			duplicateFrequency[command] = 1
		}
	}
	fmt.Println(duplicatedCommands)
	return duplicatedCommands
}

func (cm CheckerManager) printConflictingCommandsWarning(conflictingCommands []string) {
	lastCommandIndex := len(conflictingCommands) - 1
	lastCommand := conflictingCommands[lastCommandIndex]
	lastCommand = strings.Replace(lastCommand, "root", "rit", 1)
	lastCommand = strings.ReplaceAll(lastCommand, "_", " ")
	msg := fmt.Sprintf("There's a total of %d formula conflicting commands, like:\n %s", len(conflictingCommands), lastCommand)
	msg = prompt.Yellow(msg)
	fmt.Println(msg)
}
