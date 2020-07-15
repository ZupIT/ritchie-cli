package cmd

import (
	"fmt"

	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/credential"
	"github.com/ZupIT/ritchie-cli/pkg/credential/set"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

type listCredentialCmd struct {
	set.SingleSettings
}

func NewListCredentialCmd(
	ss set.SingleSettings) *cobra.Command {
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
			r ='*'
		}
		hiddenCredential = append(hiddenCredential, r)
	}
	return string(hiddenCredential)
}

func printCredentialsTable(fields credential.Fields) {
	table := uitable.New()
	table.AddRow(prompt.Bold("CONTEXT"), prompt.Bold("PROVIDER"), prompt.Bold("NAME"), prompt.Bold("CREDENTIAL"))
	switchColor := true
	for c := range fields {
		provider := fields[c]
		for _, p := range provider {
			if switchColor {
				table.AddRow(c, p.Name)
				switchColor = false
			} else {
				table.AddRow(prompt.Cyan(c), prompt.Cyan(p.Name))
				switchColor = true
			}

		}

	}
	fmt.Println(table)
}



func (l listCredentialCmd) run() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		data, err := l.ReadCredentials(set.ProviderPath())
		if err != nil {
			return err
		}

		printCredentialsTable(data)
		return nil
	}
}