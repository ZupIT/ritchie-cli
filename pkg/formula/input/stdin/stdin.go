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

package stdin

import (
	"fmt"
	"os/exec"

	"github.com/spf13/pflag"

	"github.com/ZupIT/ritchie-cli/pkg/env"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/input"
	"github.com/ZupIT/ritchie-cli/pkg/stdin"
)

type InputManager struct {
	envResolvers env.Resolvers
}

func NewInputManager(env env.Resolvers) formula.InputRunner {
	return InputManager{envResolvers: env}
}

func (in InputManager) Inputs(cmd *exec.Cmd, setup formula.Setup, _ *pflag.FlagSet) error {
	data := make(map[string]interface{})
	if err := stdin.ReadJson(cmd.Stdin, &data); err != nil {
		return err
	}

	inputs := setup.Config.Inputs
	for _, i := range inputs {
		var inputVal string
		var err error
		switch iType := i.Type; iType {
		case input.TextType, input.BoolType, input.PassType:
			inputVal = fmt.Sprintf("%v", data[i.Name])
		default:
			inputVal, err = input.ResolveIfReserved(in.envResolvers, i)
			if err != nil {
				return err
			}
		}

		if len(inputVal) != 0 {
			input.AddEnv(cmd, i.Name, inputVal)
		}
	}
	return nil
}
