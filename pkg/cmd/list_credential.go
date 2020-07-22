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
	mustHideIndex := len(credential) / 3
	var hiddenCredential []rune
	for i, r := range credential {
		if i > mustHideIndex {
			r = '*'
		}
		hiddenCredential = append(hiddenCredential, r)
	}
	return string(hiddenCredential)
}

func printCredentialsTable(fields credential.ListCredDatas) {
	table := uitable.New()
	table.MaxColWidth = 50
	table.Wrap = true

	table.AddRow(
		prompt.Bold("NAME"),
		prompt.Bold("VALUE"),
		prompt.Bold("PROVIDER"),
		prompt.Bold("CONTEXT"),
	)

	for _, c := range fields {
		table.AddRow(c.Name, hideCredential(c.Value), c.Provider, c.Context)
	}
	fmt.Println(table)
}

func (l listCredentialCmd) run() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		data, _ := l.Settings.ReadCredentialsValue()
		printCredentialsTable(data)
		return nil
	}
}
