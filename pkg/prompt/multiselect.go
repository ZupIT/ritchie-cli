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
	"github.com/ZupIT/ritchie-cli/pkg/formula/input"
)

type SurveyMultiselect struct{}

func NewSurveyMultiselect() SurveyMultiselect {
	return SurveyMultiselect{}
}

func (SurveyMultiselect) Multiselect(in formula.Input) ([]string, error) {
	value := []string{}
	multiselect := &survey.MultiSelect{
		Message: in.Label,
		Options: in.Items,
		Help:    in.Tutorial,
	}
	multiQs := []*survey.Question{
		{
			Prompt: multiselect,
		},
	}

	if input.IsRequired(in) {
		multiQs[0].Validate = survey.Required
	}

	if err := survey.Ask(multiQs, &value); err != nil {
		return value, err
	}

	return value, nil
}
