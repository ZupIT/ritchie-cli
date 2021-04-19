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

package validator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/api"
)

var (
	ErrFormulaCmdNotBeEmpty        = errors.New("this input must not be empty")
	ErrFormulaCmdMustStartWithRit  = errors.New("rit formula's command needs to start with \"rit\" [ex.: rit group verb <noun>]")
	ErrInvalidFormulaCmdSize       = errors.New("rit formula's command needs at least 2 words following \"rit\" [ex.: rit group verb]")
	ErrInvalidCharactersFormulaCmd = errors.New(`these characters are not allowed in the formula command [\ /,> <@ -]`)
)

type ValidatorManager struct {
	formulaCMD string
}

func NewValidator() ValidatorManager {
	return ValidatorManager{}
}

func (v *ValidatorManager) FormulaCommmandValidator(formula string) error {
	v.formulaCMD = formula
	if len(strings.TrimSpace(v.formulaCMD)) < 1 {
		return ErrFormulaCmdNotBeEmpty
	}

	s := strings.Split(v.formulaCMD, " ")
	if s[0] != "rit" {
		return ErrFormulaCmdMustStartWithRit
	}

	if len(s) <= 2 {
		return ErrInvalidFormulaCmdSize
	}

	if err := v.characterValidator(); err != nil {
		return err
	}

	if err := v.coreCmdValidator(); err != nil {
		return err
	}

	return nil
}

func (v *ValidatorManager) characterValidator() error {
	if strings.ContainsAny(v.formulaCMD, `\/><,@`) {
		return ErrInvalidCharactersFormulaCmd
	}
	return nil
}

func (v *ValidatorManager) coreCmdValidator() error {
	wordAfterCore := strings.Split(v.formulaCMD, " ")[1]
	for i := range api.CoreCmds {
		if wordAfterCore == api.CoreCmds[i].Usage {
			errorString := fmt.Sprintf("core command verb %q after rit\n"+
				"Use your formula group before the verb\n"+
				"Example: rit aws list bucket\n",
				api.CoreCmds[i].Usage)

			return errors.New(errorString)
		}
	}
	return nil
}
