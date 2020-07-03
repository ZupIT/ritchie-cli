package prompt

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/manifoldco/promptui"

	"github.com/ZupIT/ritchie-cli/pkg/validator"
)

type InputURL interface {
	URL(name, defaultValue string) (string, error)
}

type inputURL struct{}

type surveyURL struct{}

func NewInputURL() inputURL {
	return inputURL{}
}

func NewSurveyURL() surveyURL {
	return surveyURL{}
}

func (inputURL) URL(name, defaultValue string) (string, error) {
	prompt := promptui.Prompt{
		Label:     name,
		Pointer: promptui.PipeCursor,
		Default:   defaultValue,
		Validate:  validator.IsValidURL,
		Templates: defaultTemplate(),
	}

	return prompt.Run()
}

func (surveyURL) URL(name, defaultValue string) (string, error) {
	var value string

	validationQs := []*survey.Question{
		{
			Name:     "name",
			Prompt:   &survey.Input{
				Message:  name,
				Default:  defaultValue,
			},
			Validate: validator.IsValidSurveyURL,
		},
	}

	return value, survey.Ask(validationQs, &value)
}
