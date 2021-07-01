package tree

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
)

func TestChecker(t *testing.T) {
	tests := []struct {
		name string
		in   map[formula.RepoName]formula.Tree
		want []api.CommandID
	}{
		{
			name: "should return conflicting commands",
			in:   conflictedTrees,
			want: []api.CommandID{
				"root_aws_create_bucket",
			},
		},
		{
			name: "should return empty conflict commands",
			in:   nonConflictTrees,
			want: []api.CommandID{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checker := NewChecker(treeMock{tree: tt.in})
			got := checker.Check()
			assert.Equal(t, tt.want, got)
		})
	}
}

func BenchmarkCheck(b *testing.B) {
	tree := NewChecker(treeMock{tree: conflictedTrees})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree.Check()
	}
}

var (
	createBucket = api.Command{
		Parent:  "root",
		Usage:   "bucket",
		Help:    "create bucket for aws",
		Formula: true,
	}

	createVpc = api.Command{
		Parent:  "root",
		Usage:   "vpc",
		Help:    "create vpc for aws",
		Formula: true,
	}

	conflictedTrees = map[formula.RepoName]formula.Tree{
		core: {
			Commands: coreCmds,
		},
		"repo1": {
			Commands: api.Commands{
				"root_aws_create_bucket": createBucket,
			},
		},
		"repo2": {
			Commands: api.Commands{
				"root_aws_create_vpc":    createVpc,
				"root_aws_create_bucket": createBucket,
			},
		},
		"repo3": {
			Commands: api.Commands{
				"root_aws_create_bucket": createBucket,
			},
		},
	}

	nonConflictTrees = map[formula.RepoName]formula.Tree{
		"repo1": {
			Commands: api.Commands{
				"root_aws_create_bucket": createBucket,
			},
		},
		"repo2": {
			Commands: api.Commands{
				"root_aws_create_vpc": createVpc,
			},
		},
	}
)

type treeMock struct {
	tree     map[formula.RepoName]formula.Tree
	repo     formula.RepoName
	repoTree formula.Tree
	error    error
}

func (t treeMock) Tree() (map[formula.RepoName]formula.Tree, error) {
	return t.tree, t.error
}

func (t treeMock) MergedTree(bool) formula.Tree {
	return t.tree[t.repo]
}

func (t treeMock) TreeByRepo(formula.RepoName) (formula.Tree, error) {
	return t.repoTree, t.error
}
