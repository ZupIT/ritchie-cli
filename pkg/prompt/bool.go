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

var (
	boolOpts = map[string]bool{"yes": true, "no": false, "true": true, "false": false}
)

type SurveyBool struct{}

func NewSurveyBool() SurveyBool {
	return SurveyBool{}
}

func (SurveyBool) Bool(name string, items []string) (bool, error) {
	choice := ""
	prompt := &survey.Select{
		Message: name,
		Options: items,
	}
	if err := survey.AskOne(prompt, &choice); err != nil {
		return false, err
	}

	return boolOpts[choice], nil
}
