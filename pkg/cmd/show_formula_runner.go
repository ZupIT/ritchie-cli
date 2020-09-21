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

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

type showFormulaRunnerCmd struct {
	formula.ConfigRunner
}

func NewShowFormulaRunnerCmd(c formula.ConfigRunner) *cobra.Command {
	s := showFormulaRunnerCmd{c}

	return &cobra.Command{
		Use:       "formula-runner",
		Short:     "Show the default formula runner",
		Example:   "rit show formula-runner",
		RunE:      s.runFunc(),
		ValidArgs: []string{""},
		Args:      cobra.OnlyValidArgs,
	}
}

func (s showFormulaRunnerCmd) runFunc() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		runType, err := s.Find()
		if err != nil {
			return err
		}

		prompt.Info(fmt.Sprintf("Your default formula runner is: %q \n", runType.String()))
		return nil
	}
}
