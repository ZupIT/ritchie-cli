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

