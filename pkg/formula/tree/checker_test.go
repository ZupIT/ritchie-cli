package tree

import (
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
)



func TestChecker(t *testing.T) {
	treeMock := treeMock{
		tree: formula.Tree{
			Commands: api.Commands{
				{
					Id:     "root_mock",
					Parent: "root",
					Usage:  "mock",
					Help:   "mock for add",
				},
				{
					Id:      "root_mock",
					Parent:  "root_mock",
					Usage:   "test",
					Help:    "test for add",
					Formula: true,
				},
			},
		},
	}

	tests := []struct{
		name string
	}{
		{
			name: "Should success run",

		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			treeChecker := NewChecker(treeMock)
			treeChecker.Check()
		})
	}

}

type treeMock struct {
	tree  formula.Tree
	error error
	value string
}

func (t treeMock) Tree() (map[string]formula.Tree, error) {
	if t.value != "" {
		return map[string]formula.Tree{t.value: t.tree}, t.error
	}
	return map[string]formula.Tree{"test": t.tree}, t.error
}

func (t treeMock) MergedTree(bool) formula.Tree {
	return t.tree
}