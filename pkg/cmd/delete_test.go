package cmd

import (
	"testing"
)

func TestNewDeleteCmd(t *testing.T) {
	cmd := NewDeleteCmd()
	if cmd == nil {
		t.Errorf("NewDeleteCmd got %v", cmd)

	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}
