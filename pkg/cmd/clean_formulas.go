package cmd

import (
	"fmt"
	"os"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/stdin"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

// cleanFormulasCmd type for clean formulas command
type cleanFormulasCmd struct {
	prompt.InputBool
}

// cleanFormulas type for stdin json decoder
type cleanFormulasStdin struct {
	Confirm bool `json:"confirm"`
}

// NewCleanFormulasCmd cleans formulas in .rit/formulas file
func NewCleanFormulasCmd() *cobra.Command {
	cleanCmd := &cleanFormulasCmd{
		prompt.NewSurveyBool(),
	}

	cmd := &cobra.Command{
		Use:     "formulas",
		Short:   "Cleans all downloaded formulas",
		Example: "rit clean formulas",
		RunE:    RunFuncE(cleanCmd.runStdin(), cleanCmd.runPrompt()),
	}

	cmd.LocalFlags()

	return cmd
}

func (cleanFormula cleanFormulasCmd) runPrompt() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {

		choice, _ := cleanFormula.Bool(
			fmt.Sprint("Are you sure you want to clean all your downloaded formulas?"),
			[]string{"yes", "no"},
		)
		if !choice {
			fmt.Println("Operation cancelled")
			return nil
		}

		formulasDirectory := fmt.Sprintf("%s/formulas", api.RitchieHomeDir())
		if err := os.RemoveAll(formulasDirectory); err != nil {
			return err
		}

		prompt.Success(fmt.Sprint("Your formulas folder has been cleaned!\n"))
		return nil
	}
}

func (cleanFormula cleanFormulasCmd) runStdin() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {

		stdinArgs := cleanFormulasStdin{}

		err := stdin.ReadJson(os.Stdin, &stdinArgs)
		if err != nil {
			prompt.Error(stdin.MsgInvalidInput)
			return err
		}

		if !stdinArgs.Confirm {
			fmt.Println("Operation cancelled")
			return nil
		}

		formulasDirectory := fmt.Sprintf("%s/formulas", api.RitchieHomeDir())
		if err := os.RemoveAll(formulasDirectory); err != nil {
			return err
		}

		prompt.Success(fmt.Sprint("Your formulas folder has been cleaned!\n"))
		return nil
	}
}
