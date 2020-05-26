package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stdin"
)

var ErrNotAllowedCharacter = errors.New(`not allowed character on formula name \/,><@`)


// createFormulaCmd type for add formula command
type createFormulaCmd struct {
	formula.Creator
	prompt.InputText
	prompt.InputList
	prompt.InputBool
}

// createFormula type for stdin json decoder
type createFormula struct {
	FormulaCmd   string `json:"formulaCmd"`
	Lang         string `json:"lang"`
	LocalRepoDir string `json:"localRepoDir"`
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

		var localRepoDir string

		fCmd, err := c.Text("Enter the new formula command [ex.: rit group verb noun]", true)
		notAllowed := `\/><,@`
		if strings.ContainsAny(fCmd, notAllowed){
			return ErrNotAllowedCharacter
		}

		fmt.Println("Creating Formula ...")
		if err != nil {
			return err
		}

		lang, err := c.List("Choose the language: ", []string{"Go", "Java", "Node", "Python", "Shell"})
		if err != nil {
			return err
		}
		homeDir, _ := os.UserHomeDir()
		ritFormulasPath := fmt.Sprintf("%s/my-ritchie-formulas", homeDir)
		repoQuestion := fmt.Sprintf("Use default repo (%s)?", ritFormulasPath)
		choice, _ := c.Bool(repoQuestion, []string{"yes", "no"})
		if !choice {
			pathQuestion := fmt.Sprintf("Enter your path [ex.:%s]",ritFormulasPath)
			localRepoDir, err = c.Text(pathQuestion, true)
			if err != nil {
				return err
			}

		}

		f, err := c.Create(fCmd, lang, localRepoDir)
		if err != nil {
			return err
		}

		log.Printf("Formula in %s successfully created!\n", lang)
		log.Printf("Your formula is in %s", f.FormPath)

		return nil
	}
}

func (c createFormulaCmd) runStdin() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		fmt.Println("Creating Formula ...")

		cf := createFormula{}

		err := stdin.ReadJson(os.Stdin, &cf)
		if err != nil {
			fmt.Println("The STDIN inputs weren't informed correctly. Check the JSON used to execute the command.")
			return err
		}

		f, err := c.Create(
			cf.FormulaCmd,
			cf.Lang,
			cf.LocalRepoDir,
		)
		if err != nil {
			return err
		}

		log.Printf("Formula in %s successfully created!\n", cf.Lang)
		log.Printf("Your formula is in %s", f.FormPath)

		return nil
	}
}
