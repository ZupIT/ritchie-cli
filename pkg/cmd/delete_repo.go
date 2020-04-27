package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/formula/repo"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

// deleteRepoCmd type for delete repo command
type deleteRepoCmd struct {
	repo.Deleter
	prompt.InputText
}

// NewDeleteRepoCmd delete repository instance
func NewDeleteRepoCmd(dl repo.Deleter, it prompt.InputText) *cobra.Command {
	d := &deleteRepoCmd{dl, it}

	return &cobra.Command{
		Use:     "repo [NAME_REPOSITORY]",
		Short:   "Delete a repository.",
		Example: "rit delete repo [NAME_REPOSITORY]",
		RunE:    d.runFunc(),
	}
}

func (d deleteRepoCmd) runFunc() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		n, err := d.Text("Name of the repository: ", true)
		if err != nil {
			return err
		}

		if err = d.Delete(n); err != nil {
			return err
		}

		fmt.Printf("%q has been removed from your repositories\n", n)

		return nil
	}
}
