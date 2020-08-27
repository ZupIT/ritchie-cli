package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

type setFormulaRunnerCmd struct {
	config formula.ConfigRunner
	input  prompt.InputList
}

func NewSetFormulaRunnerCmd(c formula.ConfigRunner, i prompt.InputList) *cobra.Command {
	s := setFormulaRunnerCmd{c, i}

	return &cobra.Command{
		Use:     "formula-runner",
		Short:   "Set the default formula runner",
		Example: "rit set formula-runner",
		RunE:    s.runFunc(),
	}
}

func (c setFormulaRunnerCmd) runFunc() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		choose, err := c.input.List("Select a default formula run type", formula.RunnerTypes)
		if err != nil {
			return err
		}

		var runType formula.RunnerType
		for i := range formula.RunnerTypes {
			if formula.RunnerTypes[i] == choose {
				runType = formula.RunnerType(i)
				break
			}
		}

		if err := c.config.Create(runType); err != nil {
			return err
		}

		prompt.Success("The default formula runner has been successfully configured!")

		if runType == formula.Local {
			prompt.Warning(`
In order to run formulas locally, you must have the formula language installed on your machine,
if you don't want to install choose to run the formulas inside the docker.
`)
		}

		return nil
	}
}
