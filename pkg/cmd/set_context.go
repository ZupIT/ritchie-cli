package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/rcontext"
)

const newCtx = "Type new context?"

type setContextCmd struct {
	rcontext.FindSetter
}

func NewSetContextCmd(fs rcontext.FindSetter) *cobra.Command {
	s := setContextCmd{fs}

	return &cobra.Command{
		Use:     "context",
		Short:   "Set context",
		Example: "rit set context",
		RunE:    s.RunFunc(),
	}
}

func (s setContextCmd) RunFunc() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		ctxHolder, err := s.Find()
		if err != nil {
			return err
		}

		ctxHolder.All = append(ctxHolder.All, rcontext.DefaultCtx)
		ctxHolder.All = append(ctxHolder.All, newCtx)
		ctx, err := prompt.List("All:", ctxHolder.All)
		if err != nil {
			return err
		}

		if ctx == newCtx {
			ctx, err = prompt.String("New context: ", true)
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
