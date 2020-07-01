package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/http/headers"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/security"
	"github.com/ZupIT/ritchie-cli/pkg/server"
	"github.com/ZupIT/ritchie-cli/pkg/stdin"
)

// loginCmd type for init command
type loginCmd struct {
	security.LoginManager
	formula.RepoLoader
	prompt.InputText
	prompt.InputPassword
	server.Finder
	hc *http.Client
}

const (
	MsgUsername = "Enter your username: "
	MsgPassword = "Enter your password: "
	MsgOtp      = "Enter your two factor authentication code: "
	OtpURL      = "%s/otp"
)

// NewLoginCmd creates new cmd instance
func NewLoginCmd(
	t prompt.InputText,
	p prompt.InputPassword,
	lm security.LoginManager,
	fm formula.RepoLoader,
	sf server.Finder,
	hc *http.Client) *cobra.Command {
	l := loginCmd{
		LoginManager:  lm,
		RepoLoader:    fm,
		InputText:     t,
		InputPassword: p,
		Finder:        sf,
		hc:            hc,
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

		otpFlag, err := l.requestOtpFlag(cfg.URL, cfg.Organization)
		if err != nil {
			return err
		}

		if otpFlag {
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

func (l loginCmd) requestOtpFlag(serverUrl string, org string) (bool, error) {
	url := fmt.Sprintf(OtpURL, serverUrl)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return false, err
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set(headers.XOrg, org)

	resp, err := l.hc.Do(req)
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	bodyJson := struct {
		Otp bool `json:"otp"`
	}{}
	err = json.Unmarshal(b, &bodyJson)
	if err != nil {
		return false, err
	}
	return bodyJson.Otp, nil
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
