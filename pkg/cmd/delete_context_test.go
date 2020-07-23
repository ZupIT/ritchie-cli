package cmd

import (
	"testing"
)

func TestNewDeleteContextCmd(t *testing.T) {
	cmd := NewDeleteContextCmd(ctxFindRemoverMock{}, inputTrueMock{}, inputListMock{})
	cmd.PersistentFlags().Bool("stdin", false, "input by stdin")
	if cmd == nil {
		t.Errorf("NewDeleteContextCmd got %v", cmd)

	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}
