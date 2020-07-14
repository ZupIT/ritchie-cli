package prompt

import (
	"github.com/AlecAivazis/survey/v2"

	"github.com/ZupIT/ritchie-cli/pkg/validator"
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
			Validate: validator.IsValidSurveyEmail,
		},
	}

	return value, survey.Ask(validationQs, &value)
}
