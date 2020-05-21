package cmd

import (
	"testing"
)

func TestNewCreateFormulaCmd(t *testing.T) {
	cmd := NewCreateFormulaCmd(formCreator{}, inputTextMock{}, inputListMock{}, inputTrueMock{})
	if cmd == nil {
		t.Errorf("NewCreateFormulaCmd got %v", cmd)

	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}
