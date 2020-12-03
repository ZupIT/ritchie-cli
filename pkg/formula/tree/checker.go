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

// Check returns an api.CommandID slice if it finds
// any commands that are in conflict in different
// repositories or an api.CommandID empty slice
// if no conflicts are found.
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
