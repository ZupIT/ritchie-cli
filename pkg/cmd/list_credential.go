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

func NewListCredentialCmd(ss credential.Settings) *cobra.Command {
	l := &listCredentialCmd{ss}

	cmd := &cobra.Command{
		Use:     "credential",
		Short:   "List credentials, fields and part of values",
		Example: "rit list credential",
		RunE:    l.run(),
	}

	return cmd
}

func printCredentialsTable(fields credential.ListCredDatas) {
	table := uitable.New()
	table.Wrap = true
	table.AddRow(
		prompt.Bold("CREDENTIAL"),
		prompt.Bold("PROVIDER"),
		prompt.Bold("CONTEXT"),
	)

	for _, c := range fields {
		table.AddRow(c.Credential, c.Provider, c.Context)
	}
	if len(table.Rows) < 2 {
		setCmd := prompt.Cyan("rit set credential")
		fmt.Printf("You dont have any credential, use %s\n", setCmd)
	} else {
		fmt.Println(table)
	}
}

func (l listCredentialCmd) run() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		data, err := l.Settings.ReadCredentialsValue(l.CredentialsPath())
		if err != nil {
			return err
		}
		printCredentialsTable(data)
		return nil
	}
}
