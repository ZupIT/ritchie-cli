package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

type testFormulaCmd struct {
}

func NewTestFormulaCmd() *cobra.Command {
	s := testFormulaCmd{}

	cmd := &cobra.Command{
		Use:   "formula",
		Short: "Test your formulas locally. Use --watch flag and get real-time updates.",
		Long:  `Use this command to test your formulas locally. To make formulas development easier, you can run 
the command with the --watch flag and get real-time updates.`,
		RunE:  s.runFunc(),
	}
	cmd.Flags().BoolP("watch", "w", false, "Use this flag to watch your developing formulas")

	return cmd
}

func (s testFormulaCmd) runFunc() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		fmt.Println("Run test formula")

		return nil
	}
}
