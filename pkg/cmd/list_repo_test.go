package cmd

import (
	"testing"
)

func TestNewListRepoCmd(t *testing.T) {
	cmd := NewListRepoCmd(repoListerMock{})
	if cmd == nil {
		t.Errorf("NewListRepoCmd got %v", cmd)

	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}
