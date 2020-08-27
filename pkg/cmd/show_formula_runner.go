package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

type showFormulaRunnerCmd struct {
	formula.ConfigRunner
}

func NewShowFormulaRunnerCmd(c formula.ConfigRunner) *cobra.Command {
	s := showFormulaRunnerCmd{c}

	return &cobra.Command{
		Use:     "formula-runner",
		Short:   "Show the default formula runner",
		Example: "rit show formula-runner",
		RunE:    s.runFunc(),
	}
}

func (s showFormulaRunnerCmd) runFunc() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		runType, err := s.Find()
		if err != nil {
			return err
		}

		prompt.Info(fmt.Sprintf("Your default formula runner is: %q \n", runType.String()))
		return nil
	}
}
