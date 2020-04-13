package cmd

import (
	"testing"
)

type repoCleaner struct{}

func (repoCleaner) Clean(name string) error {
	return nil
}

func TestNewCleanRepoCmd(t *testing.T) {
	cmd := NewCleanRepoCmd(repoCleaner{}, inputTextMock{})
	if cmd == nil {
		t.Errorf("NewCleanRepoCmd got %v", cmd)

	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}
