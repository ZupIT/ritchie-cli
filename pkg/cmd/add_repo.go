package cmd

import (
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/spf13/cobra"
)

// addRepoCmd type for add repo command
type addRepoCmd struct {
	formula.RepoAdder
	prompt.InputText
	prompt.InputURL
	prompt.InputInt
}

// NewRepoAddCmd creates a new cmd instance
func NewAddRepoCmd(
	ad formula.RepoAdder,
	it prompt.InputText,
	iu prompt.InputURL,
	ii prompt.InputInt) *cobra.Command {
	a := &addRepoCmd{
		ad,
		it,
		iu,
		ii,
	}

	cmd := &cobra.Command{
		Use:     "repo",
		Short:   "Add a repository.",
		Example: "rit add repo ",
		RunE:    a.runFunc(),
	}

	return cmd
}

func (a addRepoCmd) runFunc() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		rn, err := a.Text("Name of the repository: ", true)
		if err != nil {
			return err
		}

		ur, err := a.URL("URL of the tree [http(s)://host:port/tree.json]: ", "")
		if err != nil {
			return err
		}

		pr, err := a.Int("Priority [ps.: 0 is higher priority, the lower higher the priority] :")
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
