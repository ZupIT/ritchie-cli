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
	MsgErrFormulaCanBeCoreCommand  = "core command verb %q after rit\n" + "Use your formula group before the verb\n" + "Example: rit aws list bucket\n"
)

type Manager struct{}

func New() Manager {
	return Manager{}
}

func (v *Manager) FormulaCommmandValidator(formula string) error {
	if len(strings.TrimSpace(formula)) < 1 {
		return ErrFormulaCmdNotBeEmpty
	}

	s := strings.Split(formula, " ")
	if s[0] != api.RootName {
		return ErrFormulaCmdMustStartWithRit
	}

	if len(s) <= 2 {
		return ErrInvalidFormulaCmdSize
	}

	if err := v.characterValidator(formula); err != nil {
		return err
	}

	if err := v.coreCmdValidator(formula); err != nil {
		return err
	}

	return nil
}

func (v *Manager) characterValidator(formula string) error {
	if strings.ContainsAny(formula, `\/><,@`) {
		return ErrInvalidCharactersFormulaCmd
	}
	return nil
}

func (v *Manager) coreCmdValidator(formula string) error {
	wordAfterCore := strings.Split(formula, " ")[1]
	for i := range api.CoreCmds {
		if wordAfterCore == api.CoreCmds[i].Usage {
			return fmt.Errorf(MsgErrFormulaCanBeCoreCommand, api.CoreCmds[i].Usage)
		}
	}
	return nil
}
