package cmd

import (
	"fmt"

	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
)

type ListRepoCmd struct {
	formula.RepositoryLister
}

func NewListRepoCmd(rl formula.RepositoryLister) *cobra.Command {
	lr := ListRepoCmd{rl}
	cmd := &cobra.Command{
		Use:     "repo",
		Short:   "Show a list with all your available repositories",
		Example: "rit list repo",
		RunE:    lr.runFunc(),
	}
	return cmd
}

func (lr ListRepoCmd) runFunc() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		repos, err := lr.List()
		if err != nil {
			return err
		}

		printRepos(repos)

		return nil
	}
}

func printRepos(repos formula.Repos) {
	table := uitable.New()
	table.AddRow("PRIORITY", "NAME", "VERSION")
	for _, repo := range repos {
		table.AddRow(repo.Priority, repo.Name, repo.Version)
	}
	raw := table.Bytes()
	raw = append(raw, []byte("\n")...)
	fmt.Println(string(raw))

}
