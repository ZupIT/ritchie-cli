package cmd

import (
	"fmt"
	"os"

	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/rtutorial"
	"github.com/ZupIT/ritchie-cli/pkg/stdin"
	"github.com/spf13/cobra"
)

type tutorialCmd struct {
	homePath string
	prompt.InputList
	rtutorial.FindSetter
}

const (
	tutorialStatusEnabled  = "enabled"
	tutorialStatusDisabled = "disabled"
)

// NewTutorialCmd creates tutorial command
func NewTutorialCmd(homePath string, il prompt.InputList, fs rtutorial.FindSetter) *cobra.Command {
	o := tutorialCmd{homePath, il, fs}

	cmd := &cobra.Command{
		Use:   "tutorial",
		Short: "Enable or disable the tutorial",
		Long:  "Enable or disable the tutorial",
		RunE:  RunFuncE(o.runStdin(), o.runPrompt()),
	}

	cmd.LocalFlags()

	return cmd
}

func (o tutorialCmd) runStdin() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {

		obj := struct {
			Tutorial string `json:"tutorial"`
		}{}

		err := stdin.ReadJson(os.Stdin, &obj)
		if err != nil {
			return err
		}

		fmt.Println(obj)

		return nil
	}
}

func (o tutorialCmd) runPrompt() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		msg := "Status tutorial?"
		var statusTypes = []string{tutorialStatusEnabled, tutorialStatusDisabled}

		tutorialHolder, err := o.Find()
		if err != nil {
			return err
		}

		tutorialStatusCurrent := tutorialHolder.Current
		fmt.Println("Current tutorial status: ", tutorialStatusCurrent)

		response, err := o.List(msg, statusTypes)
		if err != nil {
			return err
		}

		o.Set(response)

		prompt.Success("Set tutorial successful!")
		return nil
	}
}
