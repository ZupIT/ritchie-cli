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

package flag

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/spf13/pflag"

	"github.com/ZupIT/ritchie-cli/pkg/env"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/input"
)

type InputManager struct {
	envResolvers env.Resolvers
}

func NewInputManager(env env.Resolvers) formula.InputRunner {
	return InputManager{envResolvers: env}
}

func (in InputManager) Inputs(cmd *exec.Cmd, setup formula.Setup, flags *pflag.FlagSet) error {
	var emptyInputs formula.Inputs
	inputs := setup.Config.Inputs
	for _, i := range inputs {
		var inputVal string
		var err error
		switch i.Type {
		case input.TextType, input.PassType:
			inputVal, err = flags.GetString(i.Name)
		case input.BoolType:
			var inBool bool
			inBool, err = flags.GetBool(i.Name)
			inputVal = strconv.FormatBool(inBool)
		default:
			inputVal, err = input.ResolveIfReserved(in.envResolvers, i)
		}

		if err != nil {
			return err
		}

		if len(inputVal) != 0 {
			input.AddEnv(cmd, i.Name, inputVal)
		} else {
			emptyInputs = append(emptyInputs, i)
		}
	}

	if len(emptyInputs) > 0 {
		var emptyFlags []string
		for i := 0; i < len(emptyInputs); i++ {
			if input.IsRequired(emptyInputs[i]) {
				emptyFlags = append(emptyFlags, fmt.Sprintf("--%s", emptyInputs[i].Name))
			}
		}

		return fmt.Errorf("this flags cannot be empty [%s]", strings.Join(emptyFlags, ", "))
	}

	return nil
}
