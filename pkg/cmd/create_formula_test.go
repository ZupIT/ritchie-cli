package cmd

import (
	"os"
	"testing"
)

func TestNewCreateFormulaCmd(t *testing.T) {
	cmd := NewCreateFormulaCmd(os.TempDir(), formCreator{}, workspaceForm{}, inputTextMock{}, inputListMock{})
	cmd.PersistentFlags().Bool("stdin", false, "input by stdin")
	if cmd == nil {
		t.Errorf("NewCreateFormulaCmd got %v", cmd)

	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}
