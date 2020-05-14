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

// cleanRepoJsonDecoder type for stdin json decoder
type cleanRepoJsonDecoder struct {
	name string
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

		f := cleanRepoJsonDecoder{}

		err := stdin.ReadJson(&f)
		if err != nil {
			fmt.Println("The STDIN inputs weren't informed correctly. Check the JSON used to execute the command.")
			return err
		}

		if err := c.Clean(f.name); err != nil {
			return err
		}

		fmt.Printf("%q has been cleaned successfully\n", f.name)

		return nil
	}
}
