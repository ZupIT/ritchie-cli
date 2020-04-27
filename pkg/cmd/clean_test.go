package cmd

import (
	"testing"
)

func TestNewCleanCmd(t *testing.T) {
	cmd := NewCleanCmd()
	if cmd == nil {
		t.Errorf("NewCleanCmd got %v", cmd)

	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}
