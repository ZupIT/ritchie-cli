package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/rcontext"
)

type deleteContextCmd struct {
	rcontext.FindRemover
}

func NewDeleteContextCmd(fr rcontext.FindRemover) *cobra.Command {
	c := deleteContextCmd{fr}

	return &cobra.Command{
		Use:     "context",
		Short:   "Delete context for Ritchie-cli",
		Example: "rit delete context",
		RunE:    c.RunFunc(),
	}
}

func (c deleteContextCmd) RunFunc() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		ctxHolder, err := c.Find()
		if err != nil {
			return err
		}

		if len(ctxHolder.All) <= 0 {
			fmt.Println("You have no defined contexts")
			return nil
		}

		for i := range ctxHolder.All {
			if ctxHolder.All[i] == ctxHolder.Current {
				ctxHolder.All[i] = fmt.Sprintf("%s%s", rcontext.CurrentCtx, ctxHolder.Current)
			}
		}

		ctx, err := prompt.List("All:", ctxHolder.All)
		if err != nil {
			return err
		}

		if b, err := prompt.ListBool("Are you sure want to delete this context?", []string{"yes", "no"}); err != nil {
			return err
		} else if !b {
			return nil
		}

		if _, err := c.Remove(ctx); err != nil {
			return err
		}

		fmt.Println("Delete context successful!")
		return nil
	}
}
