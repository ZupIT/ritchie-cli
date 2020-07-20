package cmd

import (
	"testing"
)

func TestNewListCmd(t *testing.T) {
	cmd := NewListCmd()
	if cmd == nil {
		t.Errorf("NewListCmd got %v", cmd)

	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}
