package cmd

import (
	"testing"
)

func TestNewSetCmd(t *testing.T) {
	cmd := NewSetCmd()
	if cmd == nil {
		t.Errorf("NewSetCmd got %v", cmd)

	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}
