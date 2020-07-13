package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/security"
	"github.com/ZupIT/ritchie-cli/pkg/stdin"
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
}


// NewSingleInitCmd creates init command for single edition
func NewSingleInitCmd(ip prompt.InputPassword, pm security.PassphraseManager) *cobra.Command {
	o := initSingleCmd{ip, pm}

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

		return nil
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

		return nil
	}
}
