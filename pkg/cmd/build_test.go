package cmd

import (
"testing"
)

func TestNewBuildCmd(t *testing.T) {
	cmd := NewBuildCmd()
	if cmd == nil {
		t.Errorf("NewBuildCmd got %v", cmd)

	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}
