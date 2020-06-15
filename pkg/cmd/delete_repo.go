package cmd

import (
	"fmt"
	"os"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/stdin"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

// deleteRepoCmd type for delete repo command
type deleteRepoCmd struct {
	repo formula.DelLister
	prompt.InputList
	prompt.InputBool
}

// deleteRepo type for stdin json decoder
type deleteRepo struct {
	Name string `json:"name"`
}

// NewDeleteRepoCmd delete repository instance
func NewDeleteRepoCmd(dl formula.DelLister, il prompt.InputList, ib prompt.InputBool) *cobra.Command {
	d := &deleteRepoCmd{
		dl,
		il,
		ib,
	}

	cmd := &cobra.Command{
		Use:     "repo [NAME_REPOSITORY]",
		Short:   "Delete a repository.",
		Example: "rit delete repo [NAME_REPOSITORY]",
		RunE: RunFuncE(d.runStdin(), d.runPrompt()),
	}

	cmd.LocalFlags()

	return cmd
}
func rNameList(r []formula.Repository) []string {
	var names []string

	for _, repo := range r {
		names = append(names, repo.Name)
	}

	return names
}

func (d deleteRepoCmd) runPrompt() CommandRunnerFunc {
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

func (d deleteRepoCmd) runStdin() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {

		dr := deleteRepo{}

		err := stdin.ReadJson(os.Stdin, &dr)
		if err != nil {
			prompt.Error(stdin.MsgInvalidInput)
			return err
		}

		if err = d.repo.Delete(dr.Name); err != nil {
			return err
		}

		fmt.Printf("%q has been removed from your repositories\n", dr.Name)

		return nil
	}
}