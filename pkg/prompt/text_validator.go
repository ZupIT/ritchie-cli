package prompt

import (
	"github.com/AlecAivazis/survey/v2"
)

type SurveyTextValidator struct{}

func NewSurveyTextValidator() SurveyTextValidator {
	return SurveyTextValidator{}
}

func (SurveyTextValidator) Text(name string, validate func(interface{}) error, helper ...string) (string, error) {
	var value string
	validationQs := []*survey.Question{
		{
			Name:     "name",
			Validate: validate,
		},
	}

	if len(helper) > 0 {
		validationQs[0].Prompt = &survey.Input{Message: name, Help: helper[0]}
	} else {
		validationQs[0].Prompt = &survey.Input{Message: name}
	}

	return value, survey.Ask(validationQs, &value)
}
