package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stdin"
)

var (
	msgNotAllowedCharacter = fmt.Sprintf(prompt.Error, `not allowed character on formula name \/,><@`)
	ErrNotAllowedCharacter = errors.New(msgNotAllowedCharacter)
)

const notAllowedChars = `\/><,@`

// createFormulaCmd type for add formula command
type createFormulaCmd struct {
	formula.Creator
	prompt.InputText
	prompt.InputList
	prompt.InputBool
}

// CreateFormulaCmd creates a new cmd instance
func NewCreateFormulaCmd(cf formula.Creator, it prompt.InputText, il prompt.InputList, ib prompt.InputBool) *cobra.Command {
	c := createFormulaCmd{
		cf,
		it,
		il,
		ib,
	}

	cmd := &cobra.Command{
		Use:     "formula",
		Short:   "Create a new formula",
		Example: "rit create formula",
		RunE:    RunFuncE(c.runStdin(), c.runPrompt()),
	}

	cmd.LocalFlags()

	return cmd
}

func (c createFormulaCmd) runPrompt() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		fCmd, err := c.Text("Enter the new formula command [ex.: rit group verb noun]", true)
		if err != nil {
			return err
		}

		if strings.ContainsAny(fCmd, notAllowedChars) {
			return ErrNotAllowedCharacter
		}

		lang, err := c.List("Choose the language: ", []string{"Go", "Java", "Node", "Python", "Shell"})
		if err != nil {
			return err
		}
		homeDir, _ := os.UserHomeDir()
		ritFormulasPath := fmt.Sprintf("%s/my-ritchie-formulas", homeDir)
		repoQuestion := fmt.Sprintf("Use default repo (%s)?", ritFormulasPath)
		var localRepoDir string
		choice, _ := c.Bool(repoQuestion, []string{"yes", "no"})
		if !choice {
			pathQuestion := fmt.Sprintf("Enter your path [ex.:%s]", ritFormulasPath)
			localRepoDir, err = c.Text(pathQuestion, true)
			if err != nil {
				return err
			}

		}

		cf := formula.Create{
			FormulaCmd:   fCmd,
			Lang:         lang,
			LocalRepoDir: localRepoDir,
		}

		f, err := c.Create(cf)
		if err != nil {
			return err
		}

		msg := fmt.Sprintf("%s formula successfully created!\n", lang)
		fmt.Printf(prompt.Success, msg)
		msg = fmt.Sprintf("Formula path is %s \n", f.FormPath)
		fmt.Printf(prompt.Info, msg)

		return nil
	}
}

func (c createFormulaCmd) runStdin() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {

		var cf formula.Create

		if err := stdin.ReadJson(os.Stdin, &cf); err != nil {
			fmt.Println("The STDIN inputs weren't informed correctly. Check the JSON used to execute the command.")
			return err
		}

		if strings.ContainsAny(cf.FormulaCmd, notAllowedChars) {
			return ErrNotAllowedCharacter
		}

		f, err := c.Create(cf)
		if err != nil {
			return err
		}

		msg := fmt.Sprintf("%s formula successfully created!\n", cf.Lang)
		fmt.Printf(prompt.Success, msg)
		msg = fmt.Sprintf("Formula path is %s \n", f.FormPath)
		fmt.Printf(prompt.Info, msg)

		return nil
	}
}
