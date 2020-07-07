package prompt

import (
	"github.com/AlecAivazis/survey/v2"

	"github.com/ZupIT/ritchie-cli/pkg/validator"
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
			Validate: validator.IsValidSurveyURL,
		},
	}

	return value, survey.Ask(validationQs, &value)
}
