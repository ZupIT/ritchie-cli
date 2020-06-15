package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/security"
	"github.com/ZupIT/ritchie-cli/pkg/stdin"
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

	cmd := &cobra.Command{
		Use:   "user",
		Short: "Create user",
		Long:  `Create user of the organization`,
		RunE:  RunFuncE(c.runStdin(), c.runPrompt()),
	}

	cmd.LocalFlags()

	return cmd
}

func (c createUserCmd) runPrompt() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
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
			FirstName: fn,
			LastName:  ln,
			Email:     e,
			Username:  un,
			Password:  p,
		}
		if err = c.Create(u); err != nil {
			return err
		}

		fmt.Println("User created!")

		return err
	}
}

func (c createUserCmd) runStdin() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {

		u := security.User{}

		err := stdin.ReadJson(os.Stdin, &u)
		if err != nil {
			fmt.Println(prompt.Error(stdin.MsgInvalidInput))
			return err
		}

		if err := c.Create(u); err != nil {
			return err
		}

		fmt.Println("User created!")

		return nil
	}
}
