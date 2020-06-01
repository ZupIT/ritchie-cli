package cmd

import (
	"testing"
)

func TestNewLoginCmd(t *testing.T) {
	cmd := NewLoginCmd(loginManagerMock{}, repoLoaderMock{})
	if cmd == nil {
		t.Errorf("NewLoginCmd got %v", cmd)

	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}
