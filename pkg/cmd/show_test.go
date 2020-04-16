package cmd

import (
	"testing"
)

func TestNewShowCmd(t *testing.T) {
	cmd := NewShowCmd()
	if cmd == nil {
		t.Errorf("NewShowCmd got %v", cmd)

	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}
