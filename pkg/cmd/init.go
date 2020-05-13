package cmd

import (
	"errors"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/security"
	"github.com/spf13/cobra"
)

type initCmd struct {
	edition api.Edition
	prompt.InputText
	prompt.InputPassword
	prompt.InputURL
	loginManager security.LoginManager
	repoLoader   formula.Loader
}

func NewSingleInitCmd(
	ip prompt.InputPassword,
	lm security.LoginManager,
	rl formula.Loader) *cobra.Command {

	o := initCmd{api.Single, nil, ip, nil, lm, rl}

	return newInitCobra(o)
}

func NewTeamInitCmd(
	it prompt.InputText,
	iu prompt.InputURL,
	lm security.LoginManager,
	rl formula.Loader) *cobra.Command {

	o := initCmd{api.Team, it, nil, iu, lm, rl}

	return newInitCobra(o)
}

func newInitCobra(o initCmd) *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Init rit",
		RunE:  o.runFunc(),
	}
}

func (o initCmd) runFunc() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		switch o.edition {
		case api.Single:
			return o.singlePrompt()
		case api.Team:
			return o.teamPrompt()
		default:
			return errors.New("invalid CLI build, no edition defined")
		}
	}
}

func (o initCmd) singlePrompt() error {
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
