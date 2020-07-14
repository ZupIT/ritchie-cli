package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

const (
	newRepositoryPriority = "Now %q repository has priority %v"
)

type setPriorityCmd struct {
	prompt.InputList
	prompt.InputInt
	formula.RepositoryLister
	formula.RepositoryPrioritySetter
}

func NewSetPriorityCmd(
	inList prompt.InputList,
	inInt prompt.InputInt,
	repoLister formula.RepositoryLister,
	repoPriority formula.RepositoryPrioritySetter,
) *cobra.Command {
	s := setPriorityCmd{
		InputList:                inList,
		InputInt:                 inInt,
		RepositoryLister:         repoLister,
		RepositoryPrioritySetter: repoPriority,
	}

	cmd := &cobra.Command{
		Use:     "repo-priority",
		Short:   "Set a repository priority",
		Example: "rit set repo-priority",
		RunE:    s.runFunc(),
	}

	return cmd
}

func (s setPriorityCmd) runFunc() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		repositories, err := s.RepositoryLister.List()
		if err != nil {
			return err
		}

		if len(repositories) == 0 {
			prompt.Warning("You should add a repository first")
			prompt.Info("Command: rit add repo")
			return nil
		}

		var reposNames []string
		for _, r := range repositories {
			reposNames = append(reposNames, r.Name.String())
		}

		repoName, err := s.InputList.List("Repository:", reposNames)
		if err != nil {
			return err
		}

		priority, err := s.InputInt.Int("New priority:")
		if err != nil {
			return err
		}

		var repo formula.Repo
		for _, r := range repositories {
			if r.Name == formula.RepoName(repoName) {
				repo = r
				break
			}
		}

		err = s.SetPriority(repo.Name, int(priority))
		if err != nil {
			return err
		}

		successMsg := fmt.Sprintf(newRepositoryPriority, repoName, priority)
		prompt.Success(successMsg)
		return nil
	}
}
