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
	"fmt"
	"regexp"

	"github.com/AlecAivazis/survey/v2"
)

type SurveyEmail struct{}

func NewSurveyEmail() SurveyEmail {
	return SurveyEmail{}
}

func (SurveyEmail) Email(name string) (string, error) {

	var value string

	validationQs := []*survey.Question{
		{
			Name: "name",
			Prompt: &survey.Input{
				Message: name,
			},
			Validate: isValidSurveyEmail,
		},
	}

	return value, survey.Ask(validationQs, &value)
}

func isValidSurveyEmail(email interface{}) error {
	rgx := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !rgx.MatchString(email.(string)) {
		return fmt.Errorf("%s is not a valid email", email)
	}
	return nil
}
