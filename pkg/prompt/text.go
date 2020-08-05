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
)

type SurveyText struct{}

func NewSurveyText() SurveyText {
	return SurveyText{}
}

func (SurveyText) Text(name string, required bool, helper ...string) (string, error) {

	var value string

	validationQs := []*survey.Question{
		{
			Name: "name",
		},
	}

	if required {
		validationQs[0].Validate = survey.Required
	}

	if len(helper) > 0 {
		validationQs[0].Prompt = &survey.Input{Message: name, Help: helper[0]}
	} else {
		validationQs[0].Prompt = &survey.Input{Message: name}
	}

	return value, survey.Ask(validationQs, &value)
}

