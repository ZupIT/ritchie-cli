package cmd

import (
	"fmt"
	"os"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/security"
	"github.com/ZupIT/ritchie-cli/pkg/server"
	"github.com/ZupIT/ritchie-cli/pkg/stdin"
	"github.com/ZupIT/ritchie-cli/pkg/validator"
	"github.com/spf13/cobra"
)

const (
	MsgPhrase                    = "Define a passphrase for your machine: "
	MsgOrganization              = "Enter your organization: "
	msgOrganizationAlreadyExists = "The organization (%s) already exists. Do you like to override?"
	MsgServerURL                 = "URL of the server [http(s)://host]: "
	msgServerURLAlreadyExists    = "The server URL(%s) already exists. Do you like to override?"
	MsgLogin                     = "You can perform login to your organization now, or later using [rit login] command. Perform now?"
)

type initSingleCmd struct {
	prompt.InputPassword
	security.PassphraseManager
	formula.Loader
}

type initTeamCmd struct {
	prompt.InputText
	prompt.InputPassword
	prompt.InputURL
	prompt.InputBool
	server.FindSetter
	security.LoginManager
	formula.Loader
}

// NewSingleInitCmd creates init command for single edition
func NewSingleInitCmd(
	ip prompt.InputPassword,
	pm security.PassphraseManager,
	rl formula.Loader) *cobra.Command {

	o := initSingleCmd{ip, pm, rl}

	return newInitCmd(o.runStdin(), o.runPrompt())
}

// NewTeamInitCmd creates init command for team edition
func NewTeamInitCmd(
	it prompt.InputText,
	ip prompt.InputPassword,
	iu prompt.InputURL,
	ib prompt.InputBool,
	fs server.FindSetter,
	lm security.LoginManager,
	rl formula.Loader) *cobra.Command {

	o := initTeamCmd{it, ip,iu, ib, fs, lm, rl}

	return newInitCmd(o.runStdin(), o.runPrompt())
}

func newInitCmd(stdinFunc, promptFunc CommandRunnerFunc) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Init rit",
		Long:  "Initialize rit configuration",
		RunE:  RunFuncE(stdinFunc, promptFunc),
	}
	cmd.LocalFlags()
	return cmd
}

func (o initSingleCmd) runPrompt() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		pass, err := o.Password(MsgPhrase)
		if err != nil {
			return err
		}

		p := security.Passphrase(pass)
		if err := o.Save(p); err != nil {
			return err
		}

		return o.Load()
	}
}

func (o initSingleCmd) runStdin() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {

		obj := struct {
			Passphrase string `json:"passphrase"`
		}{}

		err := stdin.ReadJson(os.Stdin, &obj)
		if err != nil {
			fmt.Println(stdin.MsgInvalidInput)
			return err
		}

		p := security.Passphrase(obj.Passphrase)
		if err := o.Save(p); err != nil {
			return err
		}

		return o.Load()
	}
}

func (o initTeamCmd) runPrompt() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		cfg, err := o.Find()
		if err != nil {
			return err
		}

		if cfg.Organization != "" && len(cfg.Organization) > 0 {
			m := fmt.Sprintf(msgOrganizationAlreadyExists, cfg.Organization)
			y, err := o.Bool(m, []string{"no", "yes"})
			if err != nil {
				return err
			}
			if y {
				org, err := o.Text(MsgOrganization, true)
				if err != nil {
					return err
				}
				cfg.Organization = org
			}
		} else {
			org, err := o.Text(MsgOrganization, true)
			if err != nil {
				return err
			}
			cfg.Organization = org
		}

		if err := validator.IsValidURL(cfg.URL); err != nil {
			u, err := o.URL(MsgServerURL, "")
			if err != nil {
				return err
			}
			cfg.URL = u
		} else {
			m := fmt.Sprintf(msgServerURLAlreadyExists, cfg.URL)
			y, err := o.Bool(m, []string{"no", "yes"})
			if err != nil {
				return err
			}
			if y {
				u, err := o.URL(MsgServerURL, "")
				if err != nil {
					return err
				}
				cfg.URL = u
			}
		}

		if err := o.Set(cfg); err != nil {
			return err
		}

		y, err := o.Bool(MsgLogin, []string{"no", "yes"})
		if err != nil {
			return err
		}
		if y {
			u, err := o.Text(MsgUsername, true)
			if err != nil {
				return err
			}
			p, err := o.Password(MsgPassword)
			if err != nil {
				return err
			}
			us := security.User{
				Username: u,
				Password: p,
			}
			if err := o.Login(us); err != nil {
				return err
			}
			if err := o.Load(); err != nil {
				return err
			}
			fmt.Println("Login successfully!")
		}

		return nil
	}
}

func (o initTeamCmd) runStdin() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		cfg := server.Config{}

		err := stdin.ReadJson(os.Stdin, &cfg)
		if err != nil {
			prompt.Error(stdin.MsgInvalidInput)
			return err
		}

		if err := o.Set(cfg); err != nil {
			return err
		}

		return nil
	}
}