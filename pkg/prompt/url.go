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
