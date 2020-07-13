package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

const (
	deleteSuccessMsg = "Repository %q was deleted with success"
)

type DeleteRepoCmd struct {
	formula.RepositoryLister
	prompt.InputList
	formula.RepositoryDeleter
}

func NewDeleteRepoCmd(rl formula.RepositoryLister, il prompt.InputList, rd formula.RepositoryDeleter) *cobra.Command {
	dr := DeleteRepoCmd{rl, il, rd}
	cmd := &cobra.Command{
		Use:     "repo",
		Short:   "Delete a repository",
		Example: "rit delete repo",
		RunE:    dr.runFunc(),
	}
	return cmd
}

func (dr DeleteRepoCmd) runFunc() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		repos, err := dr.RepositoryLister.List()
		if err != nil {
			return err
		}

		if len(repos) == 0 {
			prompt.Warning("You don't have any repositories")
			return nil
		}

		var reposNames []string
		for _, r := range repos {
			reposNames = append(reposNames, r.Name.String())
		}

		repo, err := dr.InputList.List("Repository:", reposNames)
		if err != nil {
			return err
		}

		if err = dr.Delete(formula.RepoName(repo)); err != nil {
			return err
		}

		prompt.Success(fmt.Sprintf(deleteSuccessMsg, repo))
		return nil
	}
}
