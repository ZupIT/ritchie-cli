package cmd

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/credential"
	"github.com/ZupIT/ritchie-cli/pkg/env"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

type inputDeleteCredential struct {
	provider string
}

// deleteCredentialCmd type for set credential command
type deleteCredentialCmd struct {
	credential.CredDelete
	credential.ReaderPather
	env.Finder
	prompt.InputBool
	prompt.InputList
}

var deleteCredentialFlags = flags{
	{
		name:        providerFlagName,
		kind:        reflect.String,
		defValue:    "",
		description: "Provider name to delete",
	},
}

// NewDeleteCredentialCmd creates a new cmd instance
func NewDeleteCredentialCmd(
	credDelete credential.CredDelete,
	credReader credential.ReaderPather,
	env env.Finder,
	inBool prompt.InputBool,
	inList prompt.InputList,
) *cobra.Command {
	s := &deleteCredentialCmd{
		CredDelete:   credDelete,
		ReaderPather: credReader,
		Finder:       env,
		InputBool:    inBool,
		InputList:    inList,
	}

	cmd := &cobra.Command{
		Use:       "credential",
		Short:     "Delete credential",
		Long:      `Delete credential from current env`,
		RunE:      s.runFormula(),
		ValidArgs: []string{""},
		Args:      cobra.OnlyValidArgs,
	}

	addReservedFlags(cmd.Flags(), deleteCredentialFlags)

	return cmd
}

func (d deleteCredentialCmd) runFormula() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		curEnv, err := d.currentEnv()
		if err != nil {
			return err
		}
		prompt.Info(fmt.Sprintf("Current env: %s", curEnv))

		inputParams, err := d.resolveInput(cmd, curEnv)
		if err != nil {
			return err
		} else if inputParams.provider == "" {
			return nil
		}

		if err := d.Delete(inputParams.provider); err != nil {
			return err
		}

		successMessage()
		return nil
	}
}

func (d *deleteCredentialCmd) resolveInput(cmd *cobra.Command, context string) (inputDeleteCredential, error) {
	if IsFlagInput(cmd) {
		return d.resolveFlags(cmd)
	}
	return d.resolvePrompt(context)
}

func (d *deleteCredentialCmd) resolvePrompt(context string) (inputDeleteCredential, error) {
	data, err := d.ReadCredentialsValueInEnv(d.CredentialsPath(), context)
	if err != nil {
		return inputDeleteCredential{}, err
	}

	if len(data) == 0 {
		return inputDeleteCredential{}, errors.New("you have no defined credentials in this env")
	}

	providers := make([]string, 0, len(data))
	for _, c := range data {
		providers = append(providers, c.Provider)
	}

	provider, err := d.List("Credentials: ", providers)
	if err != nil {
		return inputDeleteCredential{}, err
	}

	if b, err := d.Bool("Are you sure want to delete this credential?", []string{"yes", "no"}); err != nil {
		return inputDeleteCredential{}, err
	} else if !b {
		return inputDeleteCredential{}, nil
	}
	return inputDeleteCredential{provider}, nil
}

func (d *deleteCredentialCmd) resolveFlags(cmd *cobra.Command) (inputDeleteCredential, error) {
	provider, err := cmd.Flags().GetString(providerFlagName)
	if err != nil {
		return inputDeleteCredential{}, err
	} else if provider == "" {
		return inputDeleteCredential{}, errors.New("please provide a value for 'provider'")
	}

	return inputDeleteCredential{provider}, nil
}

func (d deleteCredentialCmd) currentEnv() (string, error) {
	envHolder, err := d.Find()
	if err != nil {
		return "", err
	}

	if envHolder.Current == "" {
		envHolder.Current = env.Default
	}

	return envHolder.Current, nil
}

func successMessage() {
	prompt.Success("Delete credential successful!")
	prompt.Info("Check your credentials using rit list credential")
}
