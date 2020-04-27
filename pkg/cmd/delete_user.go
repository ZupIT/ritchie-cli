package cmd

import (
	"fmt"

	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/security"
	"github.com/spf13/cobra"
)

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

	return &cobra.Command{
		Use:   "user",
		Short: "Delete user",
		Long:  `Delete user of the organization`,
		RunE:  d.runFunc(),
	}
}

func (d deleteUserCmd) runFunc() CommandRunnerFunc {
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
