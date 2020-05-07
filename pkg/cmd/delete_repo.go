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
	prompt.InputText
	prompt.InputList
}

// NewDeleteRepoCmd delete repository instance
func NewDeleteRepoCmd(dl formula.DelLister, it prompt.InputText, il prompt.InputList) *cobra.Command {
	d := &deleteRepoCmd{
		dl,
		it,
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

		rn, err := d.ListI("Choose a repository to delete:", options)
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
