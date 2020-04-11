package cmd

import (
	"fmt"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/spf13/cobra"
)

// cleanRepoCmd type for clean repo command
type cleanRepoCmd struct {
	formula.RepoCleaner
	prompt.InputText
}

// NewCleanRepoCmd creates a new cmd instance
func NewCleanRepoCmd(cl formula.RepoCleaner, it prompt.InputText) *cobra.Command {
	c := &cleanRepoCmd{cl, it}

	cmd := &cobra.Command{
		Use:     "repo",
		Short:   "clean a repository.",
		Example: "rit clean repo ",
		RunE:    c.RunFunc(),
	}

	return cmd
}

func (c cleanRepoCmd) RunFunc() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		n, err := c.Text("Name of the repository: ", true)
		if err != nil {
			return err
		}

		if err = c.Clean(n); err != nil {
			return err
		}

		fmt.Printf("%q has been cleaned successfully\n", n)

		return nil
	}
}
