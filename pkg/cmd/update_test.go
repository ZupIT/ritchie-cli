package cmd

import (
	"testing"
)

func TestNewUpdateCmd(t *testing.T) {
	cmd := NewUpdateCmd()
	if cmd == nil {
		t.Errorf("NewUpdateCmd got %v", cmd)

	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}
