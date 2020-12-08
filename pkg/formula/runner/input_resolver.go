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
	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

var ErrInputNotRecognized = prompt.NewError("terminal input not recognized")

var _ formula.InputResolver = Resolver{}

type Resolver struct {
	types formula.TermInputTypes
}

func NewInputResolver(types formula.TermInputTypes) Resolver {
	return Resolver{types: types}
}

func (r Resolver) Resolve(inType api.TermInputType) (formula.InputRunner, error) {
	inputRunner := r.types[inType]
	if inputRunner == nil {
		return nil, ErrInputNotRecognized
	}
	if inType == api.Stdin {
		prompt.Warning("stdin input is deprecated. Please use flags for programatic formula execution")
	}

	return inputRunner, nil
}
