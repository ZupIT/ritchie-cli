package cmd

import (
	"fmt"
	"os"

	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stdin"
	"github.com/spf13/cobra"
)

type tutorialSingleCmd struct {
	homePath string
	prompt.InputBool
}

const (
	tutorialStatusOn    = "on"
	tutorialStatusOff   = "off"
	TutorialFilePattern = "%s/tutorial"
)

// NewTutorialCmd creates tutorial command
func NewTutorialCmd(homePath string, ib prompt.InputBool) *cobra.Command {
	o := tutorialSingleCmd{homePath, ib}

	cmd := &cobra.Command{
		Use:   "tutorial",
		Short: "Turns the tutorial on or off",
		Long:  "Turns the tutorial on or off",
		RunE:  RunFuncE(o.runStdin(), o.runPrompt()),
	}

	cmd.LocalFlags()

	return cmd
}

func (o tutorialSingleCmd) runStdin() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {

		obj := struct {
			Passphrase string `json:"passphrase"`
		}{}

		err := stdin.ReadJson(os.Stdin, &obj)
		if err != nil {
			fmt.Println(stdin.MsgInvalidInput)
			return err
		}

		fmt.Println(obj)

		return nil
	}
}

func (o tutorialSingleCmd) runPrompt() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		pathTutorial := fmt.Sprintf(TutorialFilePattern, o.homePath)
		msg := "Enable tutorial?"
		tutorialStatus, _ := currentTutorial(pathTutorial)
		tutorialEnabled := tutorialStatus == tutorialStatusOn
		fmt.Println("STATUS TUTORIAL: ", tutorialStatus)

		if tutorialEnabled {
			msg = "Disable tutorial?"
		}

		y, err := o.Bool(msg, []string{"yes", "no"})
		if err != nil {
			return err
		}

		if y {
			invertsTutorialStatus(pathTutorial, tutorialStatus)
		}

		fmt.Println("TUDO OK! SUA RESPOSTA: ", y)
		tutorialStatus, _ = currentTutorial(pathTutorial)
		fmt.Println("STATUS TUTORIAL: ", tutorialStatus)
		return nil
	}
}

func currentTutorial(path string) (string, error) {
	currentStatus := tutorialStatusOn

	if fileutil.Exists(path) {
		status, err := fileutil.ReadFile(path)
		if err != nil {
			return tutorialStatusOn, err
		}
		currentStatus = string(status)
	} else {
		err := createTutorial(path)
		if err != nil {
			return tutorialStatusOn, err
		}
	}

	return string(currentStatus), nil
}

func createTutorial(path string) error {
	if err := fileutil.WriteFile(path, []byte(tutorialStatusOn)); err != nil {
		return err
	}

	return nil
}

func invertsTutorialStatus(path string, currentStatus string) error {
	nextStatus := tutorialStatusOn

	if currentStatus == tutorialStatusOn {
		nextStatus = tutorialStatusOff
	}

	if err := fileutil.WriteFile(path, []byte(nextStatus)); err != nil {
		return err
	}

	return nil
}
