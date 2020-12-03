package tree

import (
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
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

// Check is used to warn the user about conflicting
// formula commands on different repos. This function doesn't
// return an error because printing an error from a unsuccessful
// warning attempt can be confusing to the user.
func (cm CheckerManager) Check() map[api.CommandID]string {
	conflictedCmds := make(map[api.CommandID]string)
	tree, _ := cm.tree.Tree()
	for k, t := range tree {
		if k == core {
			continue
		}

		for id, c := range t.Commands {
			if _, exist := conflictedCmds[id]; !exist && c.Formula {
				lastCommand := strings.Replace(id.String(), "root", "rit", 1)
				lastCommand = strings.ReplaceAll(lastCommand, "_", " ")
				conflictedCmds[id] = lastCommand
			}
		}
	}

	return conflictedCmds
}
