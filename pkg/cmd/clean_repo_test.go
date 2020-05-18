package cmd

import (
	"testing"
)

func TestNewCleanRepoCmd(t *testing.T) {
	cmd := NewCleanRepoCmd(repoCleaner{}, inputTextMock{})
	cmd.PersistentFlags().Bool("stdin", false, "input by stdin")
	if cmd == nil {
		t.Errorf("NewCleanRepoCmd got %v", cmd)

	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}
