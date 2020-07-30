package cmd

import (
	"testing"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
)

func TestFormulaCommand_Add(t *testing.T) {
	treeMock := treeMock{
		tree: formula.Tree{
			Commands: api.Commands{
				{
					Id:     "root_mock",
					Parent: "root",
					Usage:  "mock",
					Help:   "mock for add",
				},
				{
					Id:      "root_mock_test",
					Parent:  "root_mock",
					Usage:   "test",
					Help:    "test for add",
					Formula: true,
				},
			},
		},
	}
	formulaCmd := NewFormulaCommand(api.CoreCmds, treeMock, runnerMock{})
	rootCmd := &cobra.Command{
		Use: "rit",
	}
	rootCmd.PersistentFlags().Bool("stdin", false, "input by stdin")
	got := formulaCmd.Add(rootCmd)
	if got != nil {
		t.Errorf("Add got %v, want nil", got)
	}

	tests := []struct {
		name string
		args []string
	}{
		{
			name: "success default",
			args: []string{"mock", "test"},
		},
		{
			name: "success docker",
			args: []string{"mock", "test", "--docker"},
		},
		{
			name: "success stdin",
			args: []string{"mock", "test", "--stdin"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rootCmd.SetArgs(tt.args)

			if err := rootCmd.Execute(); err != nil {
				t.Errorf("%s = %v, want %v", rootCmd.Use, err, nil)
			}
		})
	}
}
