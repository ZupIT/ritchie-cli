package cmd

import (
	"fmt"

	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/rtutorial"
)

const (
	totalReposMsg   = "There are %v repos"
	totalOneRepoMsg = "There is 1 repo"
)

type listRepoCmd struct {
	formula.RepositoryLister
	rt rtutorial.Finder
}

func NewListRepoCmd(rl formula.RepositoryLister, rtf rtutorial.Finder) *cobra.Command {
	lr := listRepoCmd{rl, rtf}
	cmd := &cobra.Command{
		Use:     "repo",
		Short:   "Show a list with all your available repositories",
		Example: "rit list repo",
		RunE:    lr.runFunc(),
	}
	return cmd
}

func (lr listRepoCmd) runFunc() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		repos, err := lr.List()
		if err != nil {
			return err
		}

		printRepos(repos)

		if len(repos) != 1 {
			prompt.Info(fmt.Sprintf(totalReposMsg, len(repos)))
		} else {
			prompt.Info(totalOneRepoMsg)
		}

		tutorialHolder, err := lr.rt.Find()
		if err != nil {
			return err
		}
		tutorialListRepo(tutorialHolder.Current)
		return nil
	}
}

func printRepos(repos formula.Repos) {
	table := uitable.New()
	table.AddRow("NAME", "VERSION", "PRIORITY")
	for _, repo := range repos {
		table.AddRow(repo.Name, repo.Version, repo.Priority)
	}
	raw := table.Bytes()
	raw = append(raw, []byte("\n")...)
	fmt.Println(string(raw))

}

func tutorialListRepo(tutorialStatus string) {
	const tagTutorial = "\n[TUTORIAL]"
	const MessageTitle = "To update all repositories or delete an repository:"
	const MessageBody = ` ∙ Run "rit update repo" to update your repositories
 ∙ Run "rit delete repo" to remove a repository` + "\n"

	if tutorialStatus == tutorialStatusEnabled {
		prompt.Info(tagTutorial)
		prompt.Info(MessageTitle)
		fmt.Println(MessageBody)
	}
}
