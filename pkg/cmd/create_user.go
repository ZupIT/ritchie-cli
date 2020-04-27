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
	prompt.InputText
	prompt.InputEmail
	prompt.InputPassword
}

// NewCreateUserCmd creates a new cmd instance
func NewCreateUserCmd(
	um security.UserManager,
	it prompt.InputText,
	ie prompt.InputEmail,
	ip prompt.InputPassword) *cobra.Command {
	c := &createUserCmd{um, it, ie, ip}

	return &cobra.Command{
		Use:   "user",
		Short: "Create user",
		Long:  `Create user of the organization`,
		RunE:  c.runFunc(),
	}
}

func (c createUserCmd) runFunc() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		org, err := c.Text("Organization: ", true)
		if err != nil {
			return err
		}
		fn, err := c.Text("First name: ", true)
		if err != nil {
			return err
		}
		ln, err := c.Text("Last name: ", true)
		if err != nil {
			return err
		}
		e, err := c.Email("Email: ")
		if err != nil {
			return err
		}
		un, err := c.Text("Username: ", true)
		if err != nil {
			return err
		}
		p, err := c.Password("Password: ")
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
