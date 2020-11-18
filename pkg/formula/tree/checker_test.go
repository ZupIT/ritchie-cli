package tree

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
)

func TestChecker(t *testing.T) {
	treeMock := treeMock{
		tree: formula.Tree{
			Commands: api.Commands{
				"root_mock": {
					Parent:  "root",
					Usage:   "mock",
					Help:    "mock for add",
					Formula: true,
				},
			},
		},
	}

	tests := []struct {
		name string
	}{
		{
			name: "Should warn conflicting commands",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := captureCheckerStdout(treeMock)
			if !strings.Contains(out, "rit mock") {
				t.Error("Wrong output on tree checker function")
			}
		})
	}
}

func captureCheckerStdout(tree formula.TreeManager) string {
	treeChecker := NewChecker(tree)

	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	treeChecker.Check()

	_ = w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout
	return string(out)
}

type treeMock struct {
	tree  formula.Tree
	error error
	value string
}

func (t treeMock) Tree() (map[formula.RepoName]formula.Tree, error) {
	if t.value != "" {
		v := formula.RepoName(t.value)
		return map[formula.RepoName]formula.Tree{v: t.tree}, t.error
	}
	return map[formula.RepoName]formula.Tree{"test": t.tree}, t.error
}

func (t treeMock) MergedTree(bool) formula.Tree {
	return t.tree
}
