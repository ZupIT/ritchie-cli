package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/rcontext"
	"github.com/ZupIT/ritchie-cli/pkg/stdin"
)

const (
	context = "ctx"
)

type deleteContextCmd struct {
	rcontext.FindRemover
	prompt.InputBool
	prompt.InputList
}

func NewDeleteContextCmd(
	fr rcontext.FindRemover,
	ib prompt.InputBool,
	il prompt.InputList) *cobra.Command {
	d := deleteContextCmd{fr, ib, il}

	cmd := &cobra.Command{
		Use:     "context",
		Short:   "Delete context for Ritchie-cli",
		Example: "rit delete context",
		RunE: RunFuncE(d.runStdin(), d.runPrompt()),
	}

	cmd.LocalFlags()

	return cmd
}

func (d deleteContextCmd) runPrompt() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		ctxHolder, err := d.Find()
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

		ctx, err := d.List("All:", ctxHolder.All)
		if err != nil {
			return err
		}

		if b, err := d.Bool("Are you sure want to delete this context?", []string{"yes", "no"}); err != nil {
			return err
		} else if !b {
			return nil
		}

		if _, err := d.Remove(ctx); err != nil {
			return err
		}

		fmt.Println("Delete context successful!")
		return nil
	}
}

func (d deleteContextCmd) runStdin() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		ctxHolder, err := d.Find()
		if err != nil {
			return err
		}

		if len(ctxHolder.All) <= 0 {
			fmt.Println("You have no defined contexts")
			return nil
		}

		data, err := stdin.Parse()
		if err != nil {
			return err
		}

		if _, err := d.Remove(data[context]); err != nil {
			return err
		}

		fmt.Println("Delete context successful!")
		return nil
	}
}
