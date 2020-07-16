package cmd

import (
	"os"
	"testing"
)

func TestNewTutorialCmd(t *testing.T) {
	cmd := NewTutorialCmd(os.TempDir(), inputFalseMock{})
	cmd.PersistentFlags().Bool("stdin", false, "input by stdin")

	if cmd == nil {
		t.Errorf("NewTutorialCmd got %v", cmd)
	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}
