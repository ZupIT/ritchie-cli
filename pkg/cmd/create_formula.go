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
}

// CreateFormulaCmd creates a new cmd instance
func NewCreateFormulaCmd(cf formula.Creator, it prompt.InputText, il prompt.InputList) *cobra.Command {
	c := createFormulaCmd{cf, it, il}
	return &cobra.Command{
		Use:     "formula",
		Short:   "Create a new formula",
		Example: "rit create formula",
		RunE:    c.runFunc(),
	}
}

func (c createFormulaCmd) runFunc() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		fmt.Println("Creating Formula ...")
		fCmd, err := c.Text("New formula's command ? [ex.: rit group verb <noun>]", true)
		if err != nil {
			return err
		}
		lang, err := c.ListI("Choose the language: ", []string{"Go", "Java", "Node", "Python", "Shell"})
		if err != nil {
			return err
		}
		f, err := c.Create(fCmd, lang)
		if err != nil {
			return err
		}
		log.Printf("Formula in %s successfully created!\n", lang)
		log.Printf("Your formula is in %s", f.FormPath)
		return nil
	}
}
