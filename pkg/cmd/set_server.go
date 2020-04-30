package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/server"
)

type setServerCmd struct {
	server.Setter
	prompt.InputURL
}

func NewSetServerCmd(
	st server.Setter,
	iu prompt.InputURL,
) *cobra.Command {

	o := setServerCmd{
		Setter:   st,
		InputURL: iu,
	}

	return &cobra.Command{
		Use:   "server",
		Short: "Set server",
		Long:  `Set organization Server url `,
		RunE:  o.runFunc(),
	}
}

func (s setServerCmd) runFunc() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		u, err := s.URL("URL of the server [http(s)://host]", "")
		if err != nil {
			return err
		}
		if err := s.Set(u); err != nil {
			return err
		}
		fmt.Sprintln("Organization server url saved!")
		return nil
	}
}
