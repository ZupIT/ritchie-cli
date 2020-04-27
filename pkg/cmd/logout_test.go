package cmd

import (
	"testing"
)

func TestNewLogoutCmd(t *testing.T) {
	cmd := NewLogoutCmd(logoutManagerMock{})
	if cmd == nil {
		t.Errorf("NewLogoutCmd got %v", cmd)

	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}
