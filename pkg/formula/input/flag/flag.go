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
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/pflag"

	"github.com/ZupIT/ritchie-cli/pkg/credential"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/formula/input"
)

const (
	errInvalidInputItemsMsg = "only these input items [%s] are accepted in the %q flag"
)

type InputManager struct {
	cred credential.Resolver
}

func NewInputManager(cred credential.Resolver) formula.InputRunner {
	return InputManager{cred: cred}
}

func (in InputManager) Inputs(cmd *exec.Cmd, setup formula.Setup, flags *pflag.FlagSet) error {
	var emptyInputs formula.Inputs
	inputs := setup.Config.Inputs
	for _, i := range inputs {
		var inputVal string
		var err error

		conditionPass, err := input.VerifyConditional(cmd, i)
		if err != nil {
			return err
		}
		if !conditionPass {
			continue
		}

		switch i.Type {
		case input.TextType, input.PassType, input.DynamicType:
			inputVal, err = flags.GetString(i.Name)
			if err := validateItem(i, inputVal); err != nil {
				return err
			}

			if err := matchWithRegex(i, inputVal); err != nil {
				return err
			}

		case input.BoolType:
			var inBool bool
			inBool, err = flags.GetBool(i.Name)
			inputVal = strconv.FormatBool(inBool)
		default:
			inputVal, err = in.cred.Resolve(i.Type)
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

	var emptyFlags []string
	for _, e := range emptyInputs {
		if input.IsRequired(e) {
			emptyFlags = append(emptyFlags, fmt.Sprintf("--%s", e.Name))
		}
	}

	if len(emptyInputs) > 0 {
		return fmt.Errorf("these flags cannot be empty [%s]", strings.Join(emptyFlags, ", "))
	}

	return nil
}

func matchWithRegex(i formula.Input, inputVal string) error {
	match, err := regexp.MatchString(i.Pattern.Regex, inputVal)
	if err != nil {
		return err
	}

	if !match {
		return errors.New(i.Pattern.MismatchText)
	}

	return nil
}

func validateItem(i formula.Input, inputVal string) error {
	if len(i.Items) > 0 && !i.Items.Contains(inputVal) {
		items := strings.Join(i.Items, ", ")
		formattedName := fmt.Sprintf("--%s", i.Name)
		return fmt.Errorf(errInvalidInputItemsMsg, items, formattedName)
	}

	return nil
}
