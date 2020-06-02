package cmd

import "testing"

func TestNewTestFormulaCmd(t *testing.T) {
	cmd := NewTestFormulaCmd()
	if cmd == nil {
		t.Errorf("NewTestFormulaCmd got %v", cmd)
		return
	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}
