package cmd

import (
	"fmt"
	"os"

	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/security"
	"github.com/ZupIT/ritchie-cli/pkg/stdin"

	"github.com/spf13/cobra"
)

// deleteUserCmd type for clean repo command
type deleteUserCmd struct {
	userManager security.UserManager
	prompt.InputBool
	prompt.InputText
}

func NewDeleteUserCmd(
	um security.UserManager,
	ib prompt.InputBool,
	it prompt.InputText) *cobra.Command {
	d := &deleteUserCmd{um, ib, it}

	cmd := &cobra.Command{
		Use:   "user",
		Short: "Delete user",
		Long:  `Delete user of the organization`,
		RunE: RunFuncE(d.runStdin(), d.runPrompt()),
	}

	cmd.LocalFlags()

	return cmd
}

func (d deleteUserCmd) runPrompt() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		un, err := d.Text("Username: ", true)
		if err != nil {
			return err
		}

		e, err := d.Text("Email: ", true)
		if err != nil {
			return err
		}

		if d, err := d.Bool("Are you sure want to delete this user?", []string{"yes", "no"}); err != nil {
			return err
		} else if !d {
			return nil
		}

		fmt.Println("Deleting user...")

		u := security.User{
			Email:    e,
			Username: un,
		}
		if err = d.userManager.Delete(u); err != nil {
			return err
		}

		fmt.Printf("User %s deleted!", u.Username)

		return nil
	}
}

func (d deleteUserCmd) runStdin() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {

		u := security.User{}

		err := stdin.ReadJson(os.Stdin, &u)
		if err != nil {
			fmt.Println(prompt.Error(stdin.MsgInvalidInput))
			return err
		}

		if err = d.userManager.Delete(u); err != nil {
			return err
		}

		fmt.Printf("User %s deleted!", u.Username)

		return nil
	}
}