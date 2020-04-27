package cmd

import (
	"testing"
)

func TestNewCreateCmd(t *testing.T) {
	cmd := NewCreateCmd()
	if cmd == nil {
		t.Errorf("NewCreateCmd got %v", cmd)

	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}
