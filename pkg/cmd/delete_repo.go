package cmd

import (
	"fmt"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/stdin"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

// deleteRepoCmd type for delete repo command
type deleteRepoCmd struct {
	formula.Deleter
	prompt.InputText
}

// deleteRepoJsonDecoder type for stdin json decoder
type deleteRepoJsonDecoder struct {
	name string
}

// NewDeleteRepoCmd delete repository instance
func NewDeleteRepoCmd(dl formula.Deleter, it prompt.InputText) *cobra.Command {
	d := &deleteRepoCmd{dl, it}

	cmd := &cobra.Command{
		Use:     "repo [NAME_REPOSITORY]",
		Short:   "Delete a repository.",
		Example: "rit delete repo [NAME_REPOSITORY]",
		RunE: RunFuncE(d.runStdin(), d.runPrompt()),
	}

	cmd.LocalFlags()

	return cmd
}

func (d deleteRepoCmd) runPrompt() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		n, err := d.Text("Name of the repository: ", true)
		if err != nil {
			return err
		}

		if err = d.Delete(n); err != nil {
			return err
		}

		fmt.Printf("%q has been removed from your repositories\n", n)

		return nil
	}
}

func (d deleteRepoCmd) runStdin() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {

		dr := deleteRepoJsonDecoder{}

		err := stdin.ReadJson(&dr)
		if err != nil {
			fmt.Println("The STDIN inputs weren't informed correctly. Check the JSON used to execute the command.")
			return err
		}

		if err = d.Delete(dr.name); err != nil {
			return err
		}

		fmt.Printf("%q has been removed from your repositories\n", dr.name)

		return nil
	}
}