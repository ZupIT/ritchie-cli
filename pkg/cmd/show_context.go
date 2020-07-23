package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/rcontext"
	"github.com/ZupIT/ritchie-cli/pkg/rtutorial"
)

type showContextCmd struct {
	rcontext.Finder
	rt rtutorial.Finder
}

func NewShowContextCmd(f rcontext.Finder, tf rtutorial.Finder) *cobra.Command {
	s := showContextCmd{f, tf}

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

		tutorialHolder, err := s.rt.Find()
		if err != nil {
			return err
		}

		tutorialShowCtx(tutorialHolder.Current)
		return nil
	}
}

func tutorialShowCtx(tutorialStatus string) {
	if tutorialStatus == tutorialStatusOn {
		prompt.Info("\n[TUTORIAL] The next step is \"rit delete context\"")
	}
}
