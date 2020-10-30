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
	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stdin"
)

type setFormulaRunnerCmd struct {
	config formula.ConfigRunner
	input  prompt.InputList
}

func NewSetFormulaRunnerCmd(c formula.ConfigRunner, i prompt.InputList) *cobra.Command {
	s := setFormulaRunnerCmd{c, i}

	return &cobra.Command{
		Use:       "formula-runner",
		Short:     "Set the default formula runner",
		Example:   "rit set formula-runner",
		RunE:      RunFuncE(s.runStdin(), s.runPrompt()),
		ValidArgs: []string{""},
		Args:      cobra.OnlyValidArgs,
	}
}

func (c setFormulaRunnerCmd) runPrompt() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		choose, err := c.input.List("Select a default formula run type", formula.RunnerTypes)
		if err != nil {
			return err
		}

		runType := formula.DefaultRun
		for i := range formula.RunnerTypes {
			if formula.RunnerTypes[i] == choose {
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

func (c setFormulaRunnerCmd) runStdin() CommandRunnerFunc {
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
