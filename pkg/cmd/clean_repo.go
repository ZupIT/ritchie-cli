package cmd

import (
	"fmt"
	"os"

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

// cleanRepo type for stdin json decoder
type cleanRepo struct {
	Name string `json:"name"`
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

		prompt.Success(fmt.Sprintf("%q has been cleaned successfully", n))
		return nil
	}
}

func (c cleanRepoCmd) runStdin() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {

		f := cleanRepo{}

		err := stdin.ReadJson(os.Stdin, &f)
		if err != nil {
			prompt.Error(stdin.MsgInvalidInput)
			return err
		}

		if err := c.Clean(f.Name); err != nil {
			return err
		}

		prompt.Success(fmt.Sprintf("%q has been cleaned successfully\n", f.Name))

		return nil
	}
}
