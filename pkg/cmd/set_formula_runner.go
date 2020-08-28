package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stdin"
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
		RunE:    RunFuncE(s.runStdin(), s.runPrompt()),
	}
}

func (c setFormulaRunnerCmd) runPrompt() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		choose, err := c.input.List("Select a default formula run type", formula.RunnerTypes)
		if err != nil {
			return err
		}

		runType := formula.RunnerType(-1)
		for i := range formula.RunnerTypes {
			if formula.RunnerTypes[i] == choose {
				runType = formula.RunnerType(i)
				break
			}
		}

		if runType == -1 {
			return ErrInvalidRunType
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

func (c setFormulaRunnerCmd) runStdin() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		stdinData := struct {
			RunType string `json:"runType"`
		}{}

		if err := stdin.ReadJson(cmd.InOrStdin(), &stdinData); err != nil {
			return err
		}

		runType := formula.RunnerType(-1)
		for i := range formula.RunnerTypes {
			if formula.RunnerTypes[i] == stdinData.RunType {
				runType = formula.RunnerType(i)
				break
			}
		}

		if runType == -1 {
			return ErrInvalidRunType
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
