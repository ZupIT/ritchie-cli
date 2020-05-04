package cmd

import (
	"fmt"
	"github.com/ZupIT/ritchie-cli/pkg/formula"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/security"
)

// loginCmd type for init command
type loginCmd struct {
	security.LoginManager
	formula.Loader
	prompt.InputText
}

// NewLoginCmd creates new cmd instance
func NewLoginCmd(
	lm security.LoginManager,
	rm formula.Loader,
	it prompt.InputText) *cobra.Command {
	l := loginCmd{lm, rm, it}
	return &cobra.Command{
		Use:   "login",
		Short: "User login",
		Long:  "Authenticates and creates a session for the user of the organization",
		RunE:  l.runFunc(),
	}
}

func (l loginCmd) runFunc() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		org, err := l.Text("Enter your organization: ", true)
		if err != nil {
			return err
		}

		secret := security.Passcode(org)
		if err := l.Login(secret); err != nil {
			return err
		}

		if err := l.Load(); err != nil {
			return err
		}

		fmt.Println("Session created successfully!")
		return nil
	}
}
