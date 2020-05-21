package cmd

import (
	"testing"
)

func TestNewDeleteRepoCmd(t *testing.T) {
	cmd := NewDeleteRepoCmd(repoDeleterMock{}, inputListMock{}, inputTrueMock{})
	cmd.PersistentFlags().Bool("stdin", false, "input by stdin")
	if cmd == nil {
		t.Errorf("NewDeleteRepoCmd got %v", cmd)

	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}
