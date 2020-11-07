package cmd

import (
	"io"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/credential"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/rcontext"
	"github.com/ZupIT/ritchie-cli/pkg/stdin"
)

// deleteCredentialCmd type for set credential command
type deleteCredentialCmd struct {
	credential.CredDelete
	credential.ReaderPather
	rcontext.Finder
	prompt.InputBool
	prompt.InputList
}

// deleteCredential type for stdin json decoder
type deleteCredential struct {
	Service string `json:"service"`
}

// NewDeleteCredentialCmd creates a new cmd instance
func NewDeleteCredentialCmd(
	credDelete credential.CredDelete,
	credReader credential.ReaderPather,
	ctxFinder rcontext.Finder,
	inBool prompt.InputBool,
	inList prompt.InputList,
) *cobra.Command {
	s := &deleteCredentialCmd{
		CredDelete:   credDelete,
		ReaderPather: credReader,
		Finder:       ctxFinder,
		InputBool:    inBool,
		InputList:    inList,
	}

	cmd := &cobra.Command{
		Use:       "credential",
		Short:     "Delete credential",
		Long:      `Delete credential from current context`,
		RunE:      RunFuncE(s.runStdin(), s.runPrompt()),
		ValidArgs: []string{""},
		Args:      cobra.OnlyValidArgs,
	}
	cmd.LocalFlags()
	return cmd
}

func (d deleteCredentialCmd) runPrompt() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		context, err := d.getCurrentContext()
		if err != nil {
			return err
		}

		data, err := d.ReadCredentialsValueInContext(d.CredentialsPath(), context)
		if err != nil {
			return err
		}

		if len(data) <= 0 {
			prompt.Error("You have no defined credentials in this context")
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

		context, err := d.getCurrentContext()
		if err != nil {
			return err
		}

		data, err := d.ReadCredentialsValueInContext(d.CredentialsPath(), context)
		if err != nil {
			return err
		}

		mustDelete := false
		for _, c := range data {
			if c.Provider == dc.Service {
				mustDelete = true
			}
		}

		if !mustDelete {
			prompt.Error("You do not have credentials defined for this provider!")
			return nil
		}

		if err := d.Delete(dc.Service); err != nil {
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

func (d deleteCredentialCmd) getCurrentContext() (string, error) {
	ctxHolder, err := d.Find()
	if err != nil {
		return "", err
	}

	if ctxHolder.Current == "" {
		ctxHolder.Current = rcontext.DefaultCtx
	}

	return ctxHolder.Current, nil
}

func successMessage() {
	prompt.Success("Delete credential successful!")
	prompt.Info("Check your credentials using rit list credential")
}
