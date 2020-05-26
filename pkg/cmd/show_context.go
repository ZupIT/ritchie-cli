package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/rcontext"
)

type showContextCmd struct {
	rcontext.Finder
}

func NewShowContextCmd(f rcontext.Finder) *cobra.Command {
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

		fmt.Println(fmt.Sprintf("Current context: %s", ctx.Current))

		return nil
	}
}
