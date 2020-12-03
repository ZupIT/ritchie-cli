package tree

import (
	"reflect"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
)

func TestChecker(t *testing.T) {
	tests := []struct {
		name string
		want map[api.CommandID]string
	}{
		{
			name: "should return conflicting commands",
			want: map[api.CommandID]string{
				"root_aws_create_bucket": "rit aws create bucket",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checker := NewChecker(treeMock{})
			got := checker.Check()

			if !reflect.DeepEqual(tt.want, got) {
				t.Fatalf("Check(%s) got = %v, but want = %v", tt.name, got, tt.want)
			}
		})
	}
}

type treeMock struct {
	tree  formula.Tree
	error error
}

func (t treeMock) Tree() (map[formula.RepoName]formula.Tree, error) {
	m := map[formula.RepoName]formula.Tree{
		"repo1": {
			Commands: api.Commands{
				"root_aws_create_bucket": {
					Parent:  "root",
					Usage:   "bucket",
					Help:    "create bucket for aws",
					Formula: true,
				},
			},
		},
		"repo2": {
			Commands: api.Commands{
				"root_aws_create_bucket": {
					Parent:  "root",
					Usage:   "bucket",
					Help:    "create bucket for aws",
					Formula: true,
				},
			},
		},
	}

	return m, t.error
}

func (t treeMock) MergedTree(bool) formula.Tree {
	return t.tree
}
