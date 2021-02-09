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
	"strings"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stdin"
)

const runnerFlagName = "runner"

type setFormulaRunnerCmd struct {
	config formula.ConfigRunner
	input  prompt.InputList
}

var setFormulaRunnerFlags = flags{
	{
		name:        runnerFlagName,
		kind:        reflect.String,
		defValue:    "",
		description: fmt.Sprintf("runner name (%s)", strings.Join(formula.RunnerTypes, "|")),
	},
}

func NewSetFormulaRunnerCmd(c formula.ConfigRunner, i prompt.InputList) *cobra.Command {
	s := setFormulaRunnerCmd{c, i}

	cmd := &cobra.Command{
		Use:       "formula-runner",
		Short:     "Set the default formula runner",
		Example:   "rit set formula-runner",
		RunE:      RunFuncE(s.runStdin(), s.runFormula()),
		ValidArgs: []string{""},
		Args:      cobra.OnlyValidArgs,
	}

	addReservedFlags(cmd.Flags(), setFormulaRunnerFlags)

	return cmd
}

func (c *setFormulaRunnerCmd) runFormula() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		runType, err := c.resolveInput(cmd)
		if err != nil {
			return err
		}

		if err := c.config.Create(runType); err != nil {
			return err
		}

		prompt.Success("The default formula runner has been successfully configured!")

		if runType == formula.LocalRun {
			prompt.Warning(FormulaLocalRunWarning)
		}

		return nil
	}
}

func (c *setFormulaRunnerCmd) resolveInput(cmd *cobra.Command) (formula.RunnerType, error) {
	if IsFlagInput(cmd) {
		return c.resolveFlags(cmd)
	}
	return c.resolvePrompt()
}

func (c *setFormulaRunnerCmd) resolveFlags(cmd *cobra.Command) (formula.RunnerType, error) {
	runner, err := cmd.Flags().GetString(runnerFlagName)
	if err != nil || runner == "" {
		return formula.DefaultRun, errors.New(missingFlagText(runnerFlagName))
	}

	return c.resolveRunner(runner)
}

func (c *setFormulaRunnerCmd) resolvePrompt() (formula.RunnerType, error) {
	choose, err := c.input.List("Select a default formula run type", formula.RunnerTypes)
	if err != nil {
		return formula.DefaultRun, err
	}

	return c.resolveRunner(choose)
}

func (c *setFormulaRunnerCmd) resolveRunner(runner string) (formula.RunnerType, error) {
	for i := range formula.RunnerTypes {
		if formula.RunnerTypes[i] == runner {
			return formula.RunnerType(i), nil
		}
	}

	return formula.DefaultRun, ErrInvalidRunType
}

// TODO: Remove stdin after deprecation
func (c *setFormulaRunnerCmd) runStdin() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		stdinData := struct {
			RunType string `json:"runType"`
		}{}

		if err := stdin.ReadJson(cmd.InOrStdin(), &stdinData); err != nil {
			return err
		}

		runType := formula.DefaultRun
		for i := range formula.RunnerTypes {
			if formula.RunnerTypes[i] == stdinData.RunType {
				runType = formula.RunnerType(i)
				break
			}
		}

		if runType == formula.DefaultRun {
			return ErrInvalidRunType
		}

		if err := c.config.Create(runType); err != nil {
			return err
		}

		prompt.Success("The default formula runner has been successfully configured!")

		if runType == formula.LocalRun {
			prompt.Warning(FormulaLocalRunWarning)
		}

		return nil
	}
}
