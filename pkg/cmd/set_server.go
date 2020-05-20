package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/server"
	"github.com/ZupIT/ritchie-cli/pkg/stdin"
)

// setServerCmd type for add repo command
type setServerCmd struct {
	server.Setter
	prompt.InputURL
}

// setServer type for stdin json decoder
type setServer struct {
	Url string `json:"url"`
}

func NewSetServerCmd(
	st server.Setter,
	iu prompt.InputURL,
) *cobra.Command {

	o := setServerCmd{
		Setter:   st,
		InputURL: iu,
	}

	cmd := &cobra.Command{
		Use:   "server",
		Short: "Set server",
		Long:  `Set organization Server url `,
		RunE: RunFuncE(o.runStdin(), o.runPrompt()),
	}

	cmd.LocalFlags()

	return cmd
}

func (s setServerCmd) runPrompt() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		u, err := s.URL("URL of the server [http(s)://host]", "")
		if err != nil {
			return err
		}
		if err := s.Set(u); err != nil {
			return err
		}
		fmt.Sprintln("Organization server URL saved!")
		return nil
	}
}

func (s setServerCmd) runStdin() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {

		ss := setServer{}

		err := stdin.ReadJson(os.Stdin, &ss)
		if err != nil {
			fmt.Println("The STDIN inputs weren't informed correctly. Check the JSON used to execute the command.")
			return err
		}

		if err := s.Set(ss.Url); err != nil {
			return err
		}

		fmt.Sprintln("Organization server url saved!")
		return nil
	}
}
