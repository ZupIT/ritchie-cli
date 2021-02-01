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
	"errors"
	"fmt"
	"reflect"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/rtutorial"
	"github.com/ZupIT/ritchie-cli/pkg/stdin"
)

type tutorialCmd struct {
	homePath string
	prompt.InputList
	tutorial rtutorial.FindSetter
}

const (
	tutorialStatusEnabled  = "enabled"
	tutorialStatusDisabled = "disabled"
	enabledFlagName        = "enabled"
)

var tutorialFlags = flags{
	{
		name:        enabledFlagName,
		kind:        reflect.Bool,
		defValue:    false,
		description: "enable the tutorial",
	},
}

// NewTutorialCmd creates tutorial command.
func NewTutorialCmd(homePath string, il prompt.InputList, fs rtutorial.FindSetter) *cobra.Command {
	o := tutorialCmd{homePath, il, fs}

	cmd := &cobra.Command{
		Use:       "tutorial",
		Short:     "Enable or disable the tutorial",
		Long:      "Enable or disable the tutorial",
		RunE:      RunFuncE(o.runStdin(), o.runFormula()),
		ValidArgs: []string{""},
		Args:      cobra.OnlyValidArgs,
	}

	addReservedFlags(cmd.Flags(), tutorialFlags)

	return cmd
}

func (t *tutorialCmd) runFormula() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		enabled, err := t.resolveInput(cmd)
		if err != nil {
			return err
		}

		_, err = t.tutorial.Set(enabled)
		if err != nil {
			return err
		}
		prompt.Success("Tutorial " + enabled + "!")
		return nil
	}
}

func (t *tutorialCmd) resolveInput(cmd *cobra.Command) (string, error) {
	if IsFlagInput(cmd) {
		return t.resolveFlags(cmd)
	}
	return t.resolvePrompt()
}

func (t *tutorialCmd) resolveFlags(cmd *cobra.Command) (string, error) {
	enabled, err := cmd.Flags().GetBool(enabledFlagName)
	if err != nil {
		return "", errors.New(missingFlagText(enabledFlagName))
	} else if enabled {
		return tutorialStatusEnabled, nil
	} else {
		return tutorialStatusDisabled, nil
	}
}

func (t *tutorialCmd) resolvePrompt() (string, error) {
	msg := "Status tutorial?"
	var statusTypes = []string{tutorialStatusEnabled, tutorialStatusDisabled}

	tutorialHolder, err := t.tutorial.Find()
	if err != nil {
		return "", err
	}

	tutorialStatusCurrent := tutorialHolder.Current
	fmt.Println("Current tutorial status: ", tutorialStatusCurrent)

	response, err := t.List(msg, statusTypes)
	if err != nil {
		return "", err
	}
	return response, nil
}

func (t *tutorialCmd) runStdin() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		obj := struct {
			Tutorial string `json:"tutorial"`
		}{}

		if err := stdin.ReadJson(cmd.InOrStdin(), &obj); err != nil {
			return err
		}

		if _, err := t.tutorial.Set(obj.Tutorial); err != nil {
			return err
		}

		prompt.Success("Tutorial " + obj.Tutorial + "!")

		return nil
	}
}
