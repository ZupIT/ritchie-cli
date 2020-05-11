package cmd

import (
	"fmt"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/stdin"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

// cleanRepoCmd type for clean repo command
type cleanRepoCmd struct {
	formula.Cleaner
	prompt.InputText
}

// NewCleanRepoCmd creates a new cmd instance
func NewCleanRepoCmd(cl formula.Cleaner, it prompt.InputText) *cobra.Command {
	c := &cleanRepoCmd{cl, it}

	cmd := &cobra.Command{
		Use:     "repo",
		Short:   "clean a repository.",
		Example: "rit clean repo ",
		RunE: RunFuncE(c.runStdin(), c.runPrompt()),
	}

	cmd.LocalFlags()

	return cmd
}

func (c cleanRepoCmd) runPrompt() CommandRunnerFunc {
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

func (c cleanRepoCmd) runStdin() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		data, err := stdin.Parse()
		if err != nil {
			return err
		}

		if err := c.Clean(data[name]); err != nil {
			return err
		}

		return nil
	}
}
