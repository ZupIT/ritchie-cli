package cmd

import (
	"fmt"

	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/credential"
	"github.com/ZupIT/ritchie-cli/pkg/credential/credsingle"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

type listCredentialCmd struct {
	credential.SingleSettings
}

func NewListCredentialCmd(
	ss credential.SingleSettings) *cobra.Command {
	l := &listCredentialCmd{ss}

	cmd := &cobra.Command{
		Use:     "credential",
		Short:   "List all credential names and fields.",
		Example: "rit list credential",
		RunE:    l.run(),
	}

	return cmd
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
		data, err := l.ReadCredentials(credsingle.ProviderPath())
		if err != nil {
			return err
		}

		printCredentialsTable(data)
		return nil
	}
}