package cmd

import (
	"fmt"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/security"
	"github.com/spf13/cobra"
)

// loginCmd type for init command
type loginCmd struct {
	security.LoginManager
	formula.RepoLoader
}

// NewLoginCmd creates new cmd instance
func NewLoginCmd(lm security.LoginManager, rm formula.RepoLoader) *cobra.Command {
	l := loginCmd{lm, rm}
	return &cobra.Command{
		Use:   "login",
		Short: "User login",
		Long:  "Authenticates and creates a session for the user of the organization",
		RunE:  l.RunFunc(),
	}
}

func (l loginCmd) RunFunc() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		s, err := prompt.String("Enter your organization: ", true)
		if err != nil {
			return err
		}

		fmt.Println("Starting login...")
		secret := security.Passcode(s)
		if err := l.Login(secret); err != nil {
			return err
		}
		fmt.Println("Login successful!")

		fmt.Println("Loading repositories...")
		if err := l.Load(); err != nil {
			return err
		}
		fmt.Println("Done.")

		return nil
	}
}
