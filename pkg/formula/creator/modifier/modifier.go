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

package modifier

import "github.com/ZupIT/ritchie-cli/pkg/formula"

type Modifier interface {
	modify(b []byte) []byte
}

func NewModifiers(create formula.Create) []Modifier {
	return []Modifier{
		FormulaCmd{cf: create},
		FormulaTags{cf: create},
		TemplateRelease{},
	}
}

func Modify(b []byte, modifiers []Modifier) []byte {
	result := b
	for _, m := range modifiers {
		result = m.modify(result)
	}
	return result
}
