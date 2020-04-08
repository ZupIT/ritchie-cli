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
}

// CreateFormulaCmd creates a new cmd instance
func NewCreateFormulaCmd(cf formula.Creator) *cobra.Command {
	c := createFormulaCmd{cf}
	return &cobra.Command{
		Use:     "formula",
		Short:   "Create a new formula",
		Example: "rit create formula",
		RunE:    c.RunFunc(),
	}
}

func (c createFormulaCmd) RunFunc() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		fmt.Println("Creating Formula ...")
		fCmd, err := prompt.String("New formula's command ? [ex.: rit group verb <noun>]", true)
		if err != nil {
			return err
		}

		return c.Create(fCmd)
	}
}
