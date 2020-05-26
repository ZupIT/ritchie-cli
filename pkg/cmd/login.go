package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/security"
)

// loginCmd type for init command
type loginCmd struct {
	security.LoginManager
	formula.Loader
}

// NewLoginCmd creates new cmd instance
func NewLoginCmd(
	lm security.LoginManager,
	rm formula.Loader) *cobra.Command {
	l := loginCmd{lm, rm}
	return &cobra.Command{
		Use:   "login",
		Short: "User login",
		Long:  "Authenticates and creates a session for the user of the organization",
		RunE:  l.runFunc(),
	}
}

func (l loginCmd) runFunc() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		if err := l.Login(); err != nil {
			return err
		}

		if err := l.Load(); err != nil {
			return err
		}

		fmt.Println("Session created successfully!")
		return nil
	}
}
