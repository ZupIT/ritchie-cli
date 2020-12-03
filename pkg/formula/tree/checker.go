package tree

import (
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
func (cm CheckerManager) Check() []api.CommandID {
	hashTable := make(map[api.CommandID]bool)
	var conflicts []api.CommandID
	tree, _ := cm.tree.Tree()
	for k, t := range tree {
		if k == core {
			continue
		}

		for id, c := range t.Commands {
			if c.Formula {
				if added, exist := hashTable[id]; exist && !added {
					conflicts = append(conflicts, id)
					hashTable[id] = true
					continue
				}
				hashTable[id] = false
			}
		}
	}

	return conflicts
}
