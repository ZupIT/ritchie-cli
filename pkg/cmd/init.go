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
	msgPassphrase             = "Define a passphrase for your machine: "
	msgOrganization           = "Enter your organization: "
	msgServerURL              = "URL of the server [http(s)://host]: "
	msgServerURLAlreadyExists = "The server URL(%s) already exists. Do you like to override?"
	msgLogin                  = "You can perform login to your organization now, or later using [rit login] command. Perform now?"
)

type initSingleCmd struct {
	prompt.InputPassword
	security.PassphraseManager
	formula.Loader
}

type initTeamCmd struct {
	prompt.InputText
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
	iu prompt.InputURL,
	ib prompt.InputBool,
	fs server.FindSetter,
	lm security.LoginManager,
	rl formula.Loader) *cobra.Command {

	o := initTeamCmd{it, iu, ib, fs, lm, rl}

	return newInitCmd(o.runStdin(), o.runPrompt())
}

func newInitCmd(stdinFunc, promptFunc CommandRunnerFunc) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Init rit",
		Long:  "Initiliaze rit configuration",
		RunE:  RunFuncE(stdinFunc, promptFunc),
	}
	cmd.LocalFlags()
	return cmd
}

func (o initSingleCmd) runPrompt() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		pass, err := o.Password(msgPassphrase)
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

		org, err := o.Text(msgOrganization, true)
		if err != nil {
			return err
		}
		cfg.Organization = org

		if err := validator.IsValidURL(cfg.URL); err != nil {
			u, err := o.URL(msgServerURL, "")
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
				u, err := o.URL(msgServerURL, "")
				if err != nil {
					return err
				}
				cfg.URL = u
			}
		}

		if err := o.Set(cfg); err != nil {
			return err
		}

		y, err := o.Bool(msgLogin, []string{"no", "yes"})
		if err != nil {
			return err
		} else if y {
			if err := o.Login(); err != nil {
				return err
			}
			if err := o.Load(); err != nil {
				return err
			}
		}

		return nil
	}
}

func (o initTeamCmd) runStdin() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		cfg := server.Config{}

		err := stdin.ReadJson(os.Stdin, &cfg)
		if err != nil {
			fmt.Println(stdin.MsgInvalidInput)
			return err
		}

		if err := o.Set(cfg); err != nil {
			return err
		}

		return nil
	}
}
