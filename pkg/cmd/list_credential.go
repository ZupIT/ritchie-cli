package cmd

import (
	"fmt"

	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/credential"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

type listCredentialCmd struct {
	credential.Settings
}

func NewListCredentialCmd(
	ss credential.Settings) *cobra.Command {
	l := &listCredentialCmd{ss}

	cmd := &cobra.Command{
		Use:     "credential",
		Short:   "List all credential names and fields.",
		Example: "rit list credential",
		RunE:    l.run(),
	}

	return cmd
}

func hideCredential(credential string) string {
	noHiddenChars := len(credential) / 3
	var hiddenCredential []rune
	for i, r := range credential {
		if i < len(credential)-noHiddenChars {
			r = '*'
		}
		hiddenCredential = append(hiddenCredential, r)
	}
	return string(hiddenCredential)
}

func printCredentialsTable(fields credential.ListCredDatas) {
	table := uitable.New()
	table.MaxColWidth = 100
	table.Wrap = true
	table.AddRow(prompt.Bold("NAME"), prompt.Bold("PROVIDER"), prompt.Bold("CONTEXT"), prompt.Bold("VALUE"))
	switchColor := true
	for _, c := range fields {
		if switchColor {
			table.AddRow(c.Name, c.Provider, c.Context, hideCredential(c.Value))
			switchColor = false
		} else {
			table.AddRow(prompt.Cyan(c.Name), prompt.Cyan(c.Provider), prompt.Cyan(c.Context), prompt.Cyan(hideCredential(c.Value)))
			switchColor = true
		}

	}
	fmt.Println(table)
}

func (l listCredentialCmd) run() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {

		// TODO ler as pastas dentro de credentials para pegar os contextos

		// TODO separar os valores das credentials
		data := l.Settings.ReadCredentialsValue(credential.CredentialsPath())
		printCredentialsTable(data)
		return nil
	}
}
