package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

const (
	newRepositoryPriority = "Now '%s' has priority %s"
)

type SetPriorityCmd struct {
	prompt.InputList
	prompt.InputInt
	formula.RepositoryLister
	formula.RepositoryPrioritySetter
}

func NewSetPriorityCmd(il prompt.InputList, ii prompt.InputInt, rl formula.RepositoryLister, rs formula.RepositoryPrioritySetter) *cobra.Command {
	s := SetPriorityCmd{il, ii, rl, rs}
	cmd := &cobra.Command{
		Use:     "priority",
		Short:   "Set a repository priority",
		Example: "rit set priority",
		RunE:    s.runFunc(),
	}
	return cmd
}

func (s SetPriorityCmd) runFunc() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		repositories, err := s.RepositoryLister.List()
		if err != nil {
			return err
		}

		var reposNames []string
		for _, r := range repositories {
			reposNames = append(reposNames, r.Name)
		}

		repoName, err := s.InputList.List("Repository list:", reposNames)
		if err != nil {
			return err
		}

		priority, err := s.InputInt.Int("New priority:")
		if err != nil {
			return err
		}

		var repo formula.Repo
		for _, r := range repositories {
			if r.Name == repoName {
				repo = r
				break
			}
		}

		err = s.SetPriority(repo, int(priority))
		if err != nil {
			return err
		}

		priorityString := strconv.Itoa(int(priority))
		successMsg := fmt.Sprintf(newRepositoryPriority, repoName, priorityString)
		prompt.Success(successMsg)
		return nil
	}
}
