package cmd

import (
	"testing"
)

func TestNewSetServerCmd(t *testing.T) {
	cmd := NewSetServerCmd(setServerMock{}, inputURLMock{})
	if cmd == nil {
		t.Errorf("NewSetServerCmd got %v", cmd)
	}
	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}

type setServerMock struct {}

func (setServerMock) Set(url string) error {
	return nil
}