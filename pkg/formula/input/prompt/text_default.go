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

package prompt

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	finput "github.com/ZupIT/ritchie-cli/pkg/formula/input"
)

type SurveyDefault struct{}

func NewSurveyDefault() SurveyDefault {
	return SurveyDefault{}
}

// As validações devem ficar aqui
func (SurveyDefault) Text(i formula.Input) (string, error) {
	var value string

	input := &survey.Input{Message: i.Label}
	validationQs := []*survey.Question{
		{
			Name:   "name",
			Prompt: input,
		},
	}

	if len(i.Tutorial) > 0 {
		input.Help = i.Tutorial
	}

	if len(i.Default) > 0 {
		input.Default = i.Default
	}

	if finput.IsRequired(i) {
		validationQs[0].Validate = survey.Required
	}

	return value, survey.Ask(validationQs, &value)
}
