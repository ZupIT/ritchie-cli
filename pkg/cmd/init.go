package cmd

import (
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/security"
	"github.com/ZupIT/ritchie-cli/pkg/server"
	"github.com/spf13/cobra"
)

type initSingleCmd struct {
	prompt.InputPassword
	security.PassphraseManager
	formula.Loader
}

type initTeamCmd struct {
	prompt.InputText
	prompt.InputURL
	server.Setter
	security.LoginManager
	formula.Loader
}

// NewSingleInitCmd creates init command for single edition
func NewSingleInitCmd(
	ip prompt.InputPassword,
	pm security.PassphraseManager,
	rl formula.Loader) *cobra.Command {

	o := initSingleCmd{ip, pm, rl}

	return newInitCmd(o.runFunc())
}

// NewTeamInitCmd creates init command for team edition
func NewTeamInitCmd(
	it prompt.InputText,
	iu prompt.InputURL,
	st server.Setter,
	lm security.LoginManager,
	rl formula.Loader) *cobra.Command {

	o := initTeamCmd{it, iu, st, lm, rl}

	return newInitCmd(o.runFunc())
}

func newInitCmd(fn CommandRunnerFunc) *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Init rit",
		Long:  "Long desc for Single and Team",
		RunE:  fn,
	}
}

func (o initSingleCmd) runFunc() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		pass, err := o.Password("Define a passphrase for your machine: ")
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

func (o initTeamCmd) runFunc() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		org, err := o.Text("Enter your organization: ", true)
		if err != nil {
			return err
		}

		u, err := o.URL("URL of the server [http(s)://host]", "")
		if err != nil {
			return err
		}

		cfg := server.Config{
			Organization: org,
			URL:          u,
		}
		if err := o.Set(cfg); err != nil {
			return err
		}

		if err := o.Login(); err != nil {
			return err
		}

		if err := o.Load(); err != nil {
			return err
		}

		return nil
	}
}
