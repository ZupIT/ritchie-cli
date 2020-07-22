package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/rcontext"
)

type showContextCmd struct {
	rcontext.CtxFinder
}

func NewShowContextCmd(f rcontext.CtxFinder) *cobra.Command {
	s := showContextCmd{f}

	return &cobra.Command{
		Use:     "context",
		Short:   "Show current context",
		Example: "rit show context",
		RunE:    s.runFunc(),
	}
}

func (s showContextCmd) runFunc() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		ctx, err := s.Find()
		if err != nil {
			return err
		}

		if ctx.Current == "" {
			ctx.Current = rcontext.DefaultCtx
		}

		prompt.Info(fmt.Sprintf("Current context: %s \n", ctx.Current))
		return nil
	}
}
