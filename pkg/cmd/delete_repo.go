package cmd

import (
	"fmt"

	"github.com/ZupIT/ritchie-cli/pkg/formula"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

// deleteRepoCmd type for delete repo command
type deleteRepoCmd struct {
	repo formula.DelLister
	prompt.InputList
	prompt.InputBool
}

// NewDeleteRepoCmd delete repository instance
func NewDeleteRepoCmd(dl formula.DelLister, il prompt.InputList, ib prompt.InputBool) *cobra.Command {
	d := &deleteRepoCmd{
		dl,
		il,
		ib,
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

		repos, err := d.repo.List()
		if err != nil {
			return err
		}

		if len(repos) == 0 {
			fmt.Println("You dont have any repository to delete")
			return nil
		}

		options := rNameList(repos)

		rn, err := d.List("Choose a repository to delete:", options)
		if err != nil {
			return err
		}

		choice, _ := d.Bool(fmt.Sprintf("Want to delete %s?", rn), []string{"yes", "no"})
		if !choice {
			fmt.Println("Operation cancelled")
			return nil
		}

		if err = d.repo.Delete(rn); err != nil {
			return err
		}

		fmt.Printf("%q has been removed from your repositories\n", rn)

		return nil
	}
}
