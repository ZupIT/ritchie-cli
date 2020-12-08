package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/credential"
	"github.com/ZupIT/ritchie-cli/pkg/env"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stdin"
)

// deleteCredentialCmd type for set credential command
type deleteCredentialCmd struct {
	credential.CredDelete
	credential.ReaderPather
	env.Finder
	prompt.InputBool
	prompt.InputList
}

// deleteCredential type for stdin json decoder
type deleteCredential struct {
	Provider string `json:"provider"`
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
		RunE:      RunFuncE(s.runStdin(), s.runPrompt()),
		ValidArgs: []string{""},
		Args:      cobra.OnlyValidArgs,
	}
	cmd.LocalFlags()
	return cmd
}

func (d deleteCredentialCmd) runPrompt() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		env, err := d.currentEnv()
		if err != nil {
			return err
		}
		prompt.Info(fmt.Sprintf("Current env: %s", env))

		data, err := d.ReadCredentialsValueInEnv(d.CredentialsPath(), env)
		if err != nil {
			return err
		}

		if len(data) <= 0 {
			prompt.Error("You have no defined credentials in this env")
			return nil
		}

		var providers []string
		for _, c := range data {
			providers = append(providers, c.Provider)
		}

		cred, err := d.List("Credentials: ", providers)
		if err != nil {
			return err
		}

		if b, err := d.Bool("Are you sure want to delete this credential?", []string{"yes", "no"}); err != nil {
			return err
		} else if !b {
			return nil
		}

		if err := d.Delete(cred); err != nil {
			return err
		}

		successMessage()
		return nil
	}
}

func (d deleteCredentialCmd) runStdin() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		dc, err := d.stdinResolver(cmd.InOrStdin())
		if err != nil {
			return err
		}

		env, err := d.currentEnv()
		if err != nil {
			return err
		}

		data, err := d.ReadCredentialsValueInEnv(d.CredentialsPath(), env)
		if err != nil {
			return err
		}

		mustDelete := false
		for _, c := range data {
			if c.Provider == dc.Provider {
				mustDelete = true
			}
		}

		if !mustDelete {
			prompt.Error("You do not have credentials defined for this provider!")
			return nil
		}

		if err := d.Delete(dc.Provider); err != nil {
			return err
		}

		successMessage()
		return nil
	}
}

func (d deleteCredentialCmd) stdinResolver(reader io.Reader) (deleteCredential, error) {
	dc := deleteCredential{}

	if err := stdin.ReadJson(reader, &dc); err != nil {
		return dc, err
	}
	return dc, nil
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
