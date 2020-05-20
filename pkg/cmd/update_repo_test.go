package cmd

import (
	"testing"
)

func TestNewUpdateRepoCmd(t *testing.T) {
	cmd := NewUpdateRepoCmd(repoUpdaterMock{})
	cmd.PersistentFlags().Bool("stdin", false, "input by stdin")
	if cmd == nil {
		t.Errorf("NewUpdateRepoCmd got %v", cmd)
	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}
