package cmd

import (
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/spf13/cobra"
)

// addRepoCmd type for add repo command
type addRepoCmd struct {
	formula.RepoAdder
}

// NewRepoAddCmd creates a new cmd instance
func NewAddRepoCmd(ad formula.RepoAdder) *cobra.Command {
	a := &addRepoCmd{ad}

	cmd := &cobra.Command{
		Use:     "repo",
		Short:   "Add a repository.",
		Example: "rit add repo ",
		RunE:    a.RunFunc(),
	}

	return cmd
}

func (a addRepoCmd) RunFunc() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		rn, err := prompt.String("Name of the repository: ", true)
		if err != nil {
			return err
		}

		ur, err := prompt.URL("URL of the tree [http(s)://host:port/tree.json]: ", "")
		if err != nil {
			return err
		}

		pr, err := prompt.Integer("Priority [ps.: 0 is higher priority, the lower higher the priority] :")
		if err != nil {
			return err
		}

		r := formula.Repository{
			Priority: int(pr),
			Name:     rn,
			TreePath: ur,
		}

		if err = a.Add(r); err != nil {
			return err
		}

		return err
	}

}
