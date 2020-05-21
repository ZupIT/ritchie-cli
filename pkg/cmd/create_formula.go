package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

// createFormulaCmd type for add formula command
type createFormulaCmd struct {
	formula.Creator
	prompt.InputText
	prompt.InputList
	prompt.InputBool
}

// CreateFormulaCmd creates a new cmd instance
func NewCreateFormulaCmd(cf formula.Creator, it prompt.InputText, il prompt.InputList, ib prompt.InputBool) *cobra.Command {
	c := createFormulaCmd{cf, it, il, ib}
	return &cobra.Command{
		Use:     "formula",
		Short:   "Create a new formula",
		Example: "rit create formula",
		RunE:    c.runFunc(),
	}
}

func (c createFormulaCmd) runFunc() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		var localRepoDir = ""

		fCmd, err := c.Text("Enter the new formula command [ex.: rit group verb noun]", true)
		if err != nil {
			return err
		}
		lang, err := c.List("Choose the language: ", []string{"Go", "Java", "Node", "Python", "Shell"})
		if err != nil {
			return err
		}
		choice, err := c.Bool("Use default repo?", []string{"yes", "no"})

		if !choice {
			localRepoDir, err  = c.Text("Enter your path [ex.:/home/user/my-ritchie-formulas ]", true)
			fmt.Println("Make sure you have Makefile and tree.json")
			if err != nil {
				return err
			}

		}
		f, err := c.Create(fCmd, lang,localRepoDir)
		if err != nil {
			return err
		}

		log.Printf("Formula in %s successfully created!\n", lang)
		log.Printf("Your formula is in %s", f.FormPath)
		return nil
	}
}
