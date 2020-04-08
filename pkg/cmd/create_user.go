package cmd

import (
	"fmt"
	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/security"
)

// createUserCmd type for create user command
type createUserCmd struct {
	security.UserManager
}

// NewCreateUserCmd creates a new cmd instance
func NewCreateUserCmd(userManager security.UserManager) *cobra.Command {
	c := &createUserCmd{userManager}

	return &cobra.Command{
		Use:   "user",
		Short: "Create user",
		Long:  `Create user of the organization`,
		RunE:  c.RunFunc(),
	}
}

func (c createUserCmd) RunFunc() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		org, err := prompt.String("Organization: ", true)
		if err != nil {
			return err
		}
		fn, err := prompt.String("First name: ", true)
		if err != nil {
			return err
		}
		ln, err := prompt.String("Last name: ", true)
		if err != nil {
			return err
		}
		e, err := prompt.Email("Email: ")
		if err != nil {
			return err
		}
		un, err := prompt.String("Username: ", true)
		if err != nil {
			return err
		}
		p, err := prompt.Password("Password: ")
		if err != nil {
			return err
		}

		u := security.User{
			Organization: org,
			FirstName:    fn,
			LastName:     ln,
			Email:        e,
			Username:     un,
			Password:     p,
		}
		if err = c.Create(u); err != nil {
			return err
		}

		fmt.Println("User created!")

		return err
	}
}
