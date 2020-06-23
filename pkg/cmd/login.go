package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/security"
	"github.com/ZupIT/ritchie-cli/pkg/server"
	"github.com/ZupIT/ritchie-cli/pkg/stdin"
)

// loginCmd type for init command
type loginCmd struct {
	security.LoginManager
	formula.Loader
	prompt.InputText
	prompt.InputPassword
	server.Finder
}

const (
	MsgUsername = "Enter your username: "
	MsgPassword = "Enter your password: "
	MsgOtp      = "Enter your two factor authentication code: "
)

// NewLoginCmd creates new cmd instance
func NewLoginCmd(
	t prompt.InputText,
	p prompt.InputPassword,
	lm security.LoginManager,
	fm formula.Loader,
	sf server.Finder) *cobra.Command {
	l := loginCmd{
		LoginManager:  lm,
		Loader:        fm,
		InputText:     t,
		InputPassword: p,
		Finder:        sf,
	}
	return &cobra.Command{
		Use:   "login",
		Short: "User login",
		Long:  "Authenticates and creates a session for the user of the organization",
		RunE:  RunFuncE(l.runStdin(), l.runPrompt()),
	}
}

func (l loginCmd) runPrompt() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		cfg, err := l.Find()
		if err != nil {
			return err
		}
		u, err := l.Text(MsgUsername, true)
		if err != nil {
			return err
		}
		p, err := l.Password(MsgPassword)
		if err != nil {
			return err
		}
		var totp string
		if cfg.Otp {
			totp, err = l.Text(MsgOtp, true)
			if err != nil {
				return err
			}
		}
		us := security.User{
			Username: u,
			Password: p,
			Totp:     totp,
		}
		if err = l.Login(us); err != nil {
			return err
		}
		if err := l.Load(); err != nil {
			return err
		}
		prompt.Success("Login successfully!")
		return err
	}
}

func (l loginCmd) runStdin() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {

		u := security.User{}

		err := stdin.ReadJson(os.Stdin, &u)
		if err != nil {
			prompt.Error(stdin.MsgInvalidInput)
			return err
		}

		if err = l.Login(u); err != nil {
			return err
		}
		if err := l.Load(); err != nil {
			return err
		}
		prompt.Success("Session created successfully!")
		return err

	}
}
