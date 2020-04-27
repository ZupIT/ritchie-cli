package cmd

import (
	"testing"
)

func TestNewCreateUserCmd(t *testing.T) {
	cmd := NewCreateUserCmd(userManagerMock{}, inputTextMock{}, inputEmailMock{}, inputPasswordMock{})
	if cmd == nil {
		t.Errorf("NewCreateUserCmd got %v", cmd)

	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}
