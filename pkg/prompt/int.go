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
	"strconv"

	"github.com/AlecAivazis/survey/v2"
)

var ErrInvalidNumber = errors.New("invalid number")

type SurveyInt struct{}

func NewSurveyInt() SurveyInt {
	return SurveyInt{}
}

func (SurveyInt) Int(name string, defaultValue string, helper ...string) (int64, error) {
	var value string
	var input = &survey.Input{Message: name, Default: defaultValue}

	validationQs := []*survey.Question{
		{
			Name:     "name",
			Validate: validateSurveyIntIn,
			Prompt:   input,
		},
	}

	if len(helper) > 0 {
		input.Help = helper[0]
	}

	if err := survey.Ask(validationQs, &value); err != nil {
		return 0, err
	}

	parseInt, err := strconv.ParseInt(value, 0, 64)
	if err != nil {
		return 0, err
	}
	return parseInt, nil
}

func validateSurveyIntIn(input interface{}) error {
	if _, err := strconv.ParseInt(input.(string), 0, 64); err != nil {
		return ErrInvalidNumber
	}
	return nil
}
