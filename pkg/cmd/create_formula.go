package cmd

import (
	"fmt"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/spf13/cobra"
)

// createFormulaCmd type for add formula command
type createFormulaCmd struct {
	formula.Creator
	prompt.InputText
}

// CreateFormulaCmd creates a new cmd instance
func NewCreateFormulaCmd(cf formula.Creator, it prompt.InputText) *cobra.Command {
	c := createFormulaCmd{cf, it}
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

		return c.Create(fCmd)
	}
}
