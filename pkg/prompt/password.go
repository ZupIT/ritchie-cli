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

const PasswordType = "password"

type SurveyPassword struct{}

func NewSurveyPassword() SurveyPassword {
	return SurveyPassword{}
}

func (SurveyPassword) Password(label string, helper ...string) (string, error) {
	password := ""
	prompt := &survey.Password{
		Message: label,
	}

	if len(helper) > 0 {
		prompt.Help = helper[0]
	}

	return password, survey.AskOne(prompt, &password)
}
