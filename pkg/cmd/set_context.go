package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/rcontext"
	"github.com/ZupIT/ritchie-cli/pkg/stdin"
)

const newCtx = "Type new context?"

type setContextCmd struct {
	rcontext.FindSetter
	prompt.InputText
	prompt.InputList
}

func NewSetContextCmd(
	fs rcontext.FindSetter,
	it prompt.InputText,
	il prompt.InputList) *cobra.Command {
	s := setContextCmd{fs, it, il}

	cmd := &cobra.Command{
		Use:     "context",
		Short:   "Set context",
		Example: "rit set context",
		RunE: RunFuncE(s.runStdin(), s.runPrompt()),
	}

	cmd.LocalFlags()

	return cmd
}

func (s setContextCmd) runPrompt() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		ctxHolder, err := s.Find()
		if err != nil {
			return err
		}

		ctxHolder.All = append(ctxHolder.All, rcontext.DefaultCtx)
		ctxHolder.All = append(ctxHolder.All, newCtx)
		ctx, err := s.List("All:", ctxHolder.All)
		if err != nil {
			return err
		}

		if ctx == newCtx {
			ctx, err = s.Text("New context: ", true)
			if err != nil {
				return err
			}
		}

		if _, err := s.Set(ctx); err != nil {
			return err
		}

		fmt.Println("Set context successful!")
		return nil
	}

}

func (s setContextCmd) runStdin() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		data, err := stdin.Parse()
		if err != nil {
			return err
		}

		if _, err := s.Set(data[context]); err != nil {
			return err
		}

		fmt.Println("Set context successful!")
		return nil
	}
}
