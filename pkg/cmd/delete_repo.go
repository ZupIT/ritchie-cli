package cmd

import (
	"fmt"

	"github.com/ZupIT/ritchie-cli/pkg/formula"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

// deleteRepoCmd type for delete repo command
type deleteRepoCmd struct {
	formula.DelLister
	input prompt.InputList
}

// NewDeleteRepoCmd delete repository instance
func NewDeleteRepoCmd(dl formula.DelLister, il prompt.InputList) *cobra.Command {
	d := &deleteRepoCmd{
		dl,
		il,
	}

	return &cobra.Command{
		Use:     "repo [NAME_REPOSITORY]",
		Short:   "Delete a repository.",
		Example: "rit delete repo [NAME_REPOSITORY]",
		RunE:    d.runFunc(),
	}
}
func rNameList(r []formula.Repository) []string {
	var names []string

	for _, repo := range r {
		names = append(names, repo.Name)
	}

	return names
}

func (d deleteRepoCmd) runFunc() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {

		repos, err := d.List()
		if err != nil {
			return err
		}

		if len(repos) <= 0 {
			fmt.Println("You dont have any repository to delete")
			return nil
		}

		options := rNameList(repos)

		rn, err := d.input.List("Choose a repository to delete:", options)
		if err != nil {
			return err
		}

		if err = d.Delete(rn); err != nil {
			return err
		}

		fmt.Printf("%q has been removed from your repositories\n", rn)

		return nil
	}
}
