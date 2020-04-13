package cmd

import (
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
)

type repoAdder struct{}

func (repoAdder) Add(d formula.Repository) error {
	return nil
}

func TestNewAddRepoCmd(t *testing.T) {
	cmd := NewAddRepoCmd(repoAdder{}, inputTextMock{}, inputURLMock{}, inputIntMock{})
	if cmd == nil {
		t.Errorf("NewAddRepoCmd got %v", cmd)

	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}
