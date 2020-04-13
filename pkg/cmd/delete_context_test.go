package cmd

import (
	"testing"
)

func TestNewDeleteContextCmd(t *testing.T) {
	cmd := NewDeleteContextCmd(findRemoverMock{}, inputBoolMock{}, inputListMock{})
	if cmd == nil {
		t.Errorf("NewDeleteContextCmd got %v", cmd)

	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}
