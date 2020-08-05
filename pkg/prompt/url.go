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
	"errors"
	"net/url"

	"github.com/AlecAivazis/survey/v2"
)

type SurveyURL struct{}

func NewSurveyURL() SurveyURL {
	return SurveyURL{}
}

func (SurveyURL) URL(name, defaultValue string) (string, error) {
	var value string

	validationQs := []*survey.Question{
		{
			Name: "name",
			Prompt: &survey.Input{
				Message: name,
				Default: defaultValue,
			},
			Validate: isValidSurveyURL,
		},
	}

	return value, survey.Ask(validationQs, &value)
}

func isValidSurveyURL(value interface{}) error {
	_, err := url.ParseRequestURI(value.(string))
	if err != nil {
		return errors.New("invalid URL")
	}
	return nil
}
