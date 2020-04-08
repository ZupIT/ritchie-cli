package cmd

import (
	"fmt"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/security"
	"github.com/spf13/cobra"
)

type deleteUserCmd struct {
	userManager security.UserManager
}

func NewDeleteUserCmd(userManager security.UserManager) *cobra.Command {
	d := &deleteUserCmd{userManager}

	return &cobra.Command{
		Use:   "user",
		Short: "Delete user",
		Long:  `Delete user of the organization`,
		RunE:  d.RunFunc(),
	}
}

func (d deleteUserCmd) RunFunc() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		un, err := prompt.String("Username: ", true)
		if err != nil {
			return err
		}

		e, err := prompt.String("Email: ", true)
		if err != nil {
			return err
		}

		if d, err := prompt.ListBool("Are you sure want to delete this user?", []string{"yes", "no"}); err != nil {
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
