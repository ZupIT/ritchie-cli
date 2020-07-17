package cmd

import (
	"os"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/formula/creator/template"
)

func TestNewCreateFormulaCmd(t *testing.T) {

	tlpM := template.NewManager("../../testdata")
	cmd := NewCreateFormulaCmd(
		os.TempDir(),
		formCreator{},
		tlpM, workspaceForm{},
		inputTextMock{},
		inputTextValidatorMock{},
		inputListMock{},
	)
	cmd.PersistentFlags().Bool("stdin", false, "input by stdin")
	if cmd == nil {
		t.Errorf("NewCreateFormulaCmd got %v", cmd)

	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}
