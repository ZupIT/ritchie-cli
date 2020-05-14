package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stdin"
)

// createFormulaCmd type for add formula command
type createFormulaCmd struct {
	formula.Creator
	prompt.InputText
	prompt.InputList
}

// createFormulaJsonDecoder type for stdin json decoder
type createFormulaJsonDecoder struct {
	formulaCmd string
	lang string
}

// CreateFormulaCmd creates a new cmd instance
func NewCreateFormulaCmd(cf formula.Creator, it prompt.InputText, il prompt.InputList) *cobra.Command {
	c := createFormulaCmd{
		cf,
		it,
		il,
	}

	cmd := &cobra.Command{
		Use:     "formula",
		Short:   "Create a new formula",
		Example: "rit create formula",
		RunE: RunFuncE(c.runStdin(), c.runPrompt()),
	}

	cmd.LocalFlags()

	return cmd
}

func (c createFormulaCmd) runPrompt() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		fmt.Println("Creating Formula ...")

		fCmd, err := c.Text("New formula's command ? [ex.: rit group verb <noun>]", true)
		if err != nil {
			return err
		}

		lang, err := c.List("Choose the language: ", []string{"Go", "Java", "Node", "Python", "Shell"})
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

func (c createFormulaCmd) runStdin() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		fmt.Println("Creating Formula ...")

		cf := createFormulaJsonDecoder{}

		err := stdin.ReadJson(&cf)
		if err != nil {
			fmt.Println("The STDIN inputs weren't informed correctly. Check the JSON used to execute the command.")
			return err
		}

		f, err := c.Create(
			cf.formulaCmd,
			cf.lang,
			)
		if err != nil {
			return err
		}

		log.Printf("Formula in %s successfully created!\n", cf.lang)
		log.Printf("Your formula is in %s", f.FormPath)

		return nil
	}
}
