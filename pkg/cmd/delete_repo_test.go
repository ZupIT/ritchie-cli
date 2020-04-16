package cmd

import (
	"testing"
)

func TestNewDeleteRepoCmd(t *testing.T) {
	cmd := NewDeleteRepoCmd(repoDeleterMock{}, inputTextMock{})
	if cmd == nil {
		t.Errorf("NewDeleteRepoCmd got %v", cmd)

	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}
