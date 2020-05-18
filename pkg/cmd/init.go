package cmd

import (
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/security"
	"github.com/spf13/cobra"
)

type initSingleCmd struct {
	prompt.InputPassword
	security.LoginManager
	formula.Loader
}

type initTeamCmd struct {
	prompt.InputText
	prompt.InputURL
	security.LoginManager
	formula.Loader
}

func NewSingleInitCmd(
	ip prompt.InputPassword,
	lm security.LoginManager,
	rl formula.Loader) *cobra.Command {

	o := initSingleCmd{ip, lm, rl}

	return newInitCobra(o)
}

func NewTeamInitCmd(
	it prompt.InputText,
	iu prompt.InputURL,
	lm security.LoginManager,
	rl formula.Loader) *cobra.Command {

	o := initTeamCmd{it, iu, lm, rl}

	return newInitCobra(o)
}

func newInitCobra(o initSingleCmd) *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Init rit",
		RunE:  o.runFunc(),
	}
}

func (o initSingleCmd) runFunc() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		pass, err := o.Password("Define a passphrase for your machine: ")
		if err != nil {
			return err
		}

		p := security.Passcode(pass)
		if err := o.loginManager.Login(p); err != nil {
			return err
		}

		return o.repoLoader.Load()
	}
}

func (o initCmd) singlePrompt() error {

}

func (o initCmd) teamPrompt() error {
	// set org in home
	pass, err := o.Text("Enter your organization: ", true)
	if err != nil {
		return err
	}
	p := security.Passcode(pass)

	// set serverURL
	u, err := o.URL("URL of the server [http(s)://host]", "")
	if err != nil {
		return err
	}

	return nil
}
