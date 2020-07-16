package cmd

import (
	"os"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/stdin"
)

func TestNewCleanFormulasCmd(t *testing.T) {
	cmd := NewCleanFormulasCmd()
	cmd.PersistentFlags().Bool("stdin", false, "input by stdin")
	if cmd == nil {
		t.Errorf("NewCleanFormulasCmd got %v", cmd)
	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}

func TestNewCleanFormulasCmdStdin(t *testing.T) {
	t.Run("confirm", func(t *testing.T) {
		runStdinCommand(t, `{"confirm": true}`)
	})

	t.Run("not confirm", func(t *testing.T) {
		runStdinCommand(t, `{"confirm": false}`)
	})
}

func runStdinCommand(t *testing.T, content string) {
	tmpfile, oldStdin, err := stdin.WriteToStdin(content)
	defer os.Remove(tmpfile.Name())
	defer func() { os.Stdin = oldStdin }()
	if err != nil {
		t.Errorf("TestNewCleanFormulasCmdStdin got error %v", err)
	}

	cmd := NewCleanFormulasCmd()
	cmd.PersistentFlags().Bool("stdin", true, "input by stdin")
	if cmd == nil {
		t.Errorf("NewCleanFormulasCmd got %v", cmd)
	}

	if err := cmd.Execute(); err != nil {
		t.Errorf("%s = %v, want %v", cmd.Use, err, nil)
	}
}

