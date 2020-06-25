package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
)

// updateRepoCmd type for update command
type updateRepoCmd struct {
	formula.Updater
}

// NewUpdateRepoCmd creates a new cmd instance
func NewUpdateRepoCmd(up formula.Updater) *cobra.Command {
	u := &updateRepoCmd{up}

	cmd := &cobra.Command{
		Use:     "repo",
		Short:   "Update all repositories",
		Example: "rit update repo",
		RunE:    u.runFunc(),
	}

	return cmd
}

func (u updateRepoCmd) runFunc() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		if err := u.Update(); err != nil {
			return err
		}

		return nil
	}
}
