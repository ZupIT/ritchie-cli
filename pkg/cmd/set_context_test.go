package cmd

import (
	"testing"
)

func TestNewSetContextCmd(t *testing.T) {
	cmd := NewSetContextCmd(ctxFindSetterMock{}, inputTextMock{}, inputListMock{})
	if cmd == nil {
		t.Errorf("NewSetContextCmd got %v", cmd)
	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}
