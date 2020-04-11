package cmd

import (
	"testing"
)

func TestNewAddCmd(t *testing.T) {
	cmd := NewAddCmd()
	if cmd == nil {
		t.Errorf("NewAddCmd got %v", cmd)

	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}
