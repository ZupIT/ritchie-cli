/*
 * Copyright 2020 ZUP IT SERVICOS EM TECNOLOGIA E INOVACAO SA
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cmd

import (
	"fmt"

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

		err := stdin.ReadJson(cmd.InOrStdin(), &obj)
		if err != nil {
			return err
		}

		_, err = o.Set(obj.Tutorial)
		if err != nil {
			return err
		}
		prompt.Success("Set tutorial successful!")

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

		_, err = o.Set(response)
		if err != nil {
			return err
		}
		prompt.Success("Set tutorial successful!")
		return nil
	}
}
