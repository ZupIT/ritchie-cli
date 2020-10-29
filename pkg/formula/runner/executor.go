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

package runner

import (
	"github.com/ZupIT/ritchie-cli/pkg/formula"
)

var _ formula.Executor = ExecutorManager{}

type ExecutorManager struct {
	runners       formula.Runners
	preRunBuilder formula.PreRunBuilder
	config        formula.ConfigRunner
}

func NewExecutor(runners formula.Runners, preRunBuilder formula.PreRunBuilder, config formula.ConfigRunner) ExecutorManager {
	return ExecutorManager{
		runners:       runners,
		preRunBuilder: preRunBuilder,
		config:        config,
	}
}

func (ex ExecutorManager) Execute(exe formula.ExecuteData) error {
	runType := exe.RunType
	runner := ex.runners[runType]

	if runner == nil {
		configType, err := ex.config.Find()
		if err != nil {
			return err
		}

		runner = ex.runners[configType]
	}

	if exe.Def.RepoName == "local" {
		ex.preRunBuilder.Build(exe.Def.Path)
	}

	if err := runner.Run(exe.Def, exe.InType, exe.Verbose, exe.Flags); err != nil {
		return err
	}

	return nil
}
