package cmd

import (
	"github.com/ZupIT/ritchie-cli/pkg/formula/repo"

	"github.com/spf13/cobra"
)

// updateRepoCmd type for update command
type updateRepoCmd struct {
	repo.Updater
}

// NewUpdateRepoCmd creates a new cmd instance
func NewUpdateRepoCmd(up repo.Updater) *cobra.Command {
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
