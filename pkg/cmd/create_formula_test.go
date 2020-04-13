package cmd

import (
	"testing"
)

type formCreator struct{}

func (formCreator) Create(formulaCmd string) error {
	return nil
}

func TestNewCreateFormulaCmd(t *testing.T) {
	cmd := NewCreateFormulaCmd(formCreator{}, inputTextMock{})
	if cmd == nil {
		t.Errorf("NewCreateFormulaCmd got %v", cmd)

	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}
