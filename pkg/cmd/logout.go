package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/security"
)

type logoutCmd struct {
	security.LogoutManager
}

// NewLogoutCmd creates new cmd instance of logout command
func NewLogoutCmd(lm security.LogoutManager) *cobra.Command {
	l := logoutCmd{lm}
	return &cobra.Command{
		Use:     "logout",
		Short:   "User logout",
		Long:    "Destroy the user session of the organization",
		Example: "rit logout",
		RunE:    l.runFunc(),
	}
}

func (l logoutCmd) runFunc() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		if err := l.Logout(); err != nil {
			return err
		}

		prompt.Success("Logout successful!")
		return nil
	}
}
