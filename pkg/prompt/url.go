package prompt

import (
	"github.com/AlecAivazis/survey/v2"

	"github.com/ZupIT/ritchie-cli/pkg/validator"
)

type InputURL interface {
	URL(name, defaultValue string) (string, error)
}

// type inputURL struct{}
//
// func NewInputURL() inputURL {
// 	return inputURL{}
// }
//
// // URL show a prompt and parse to string.
// func (inputURL) URL(name, defaultValue string) (string, error) {
// 	prompt := promptui.Prompt{
// 		Label:     name,
// 		Pointer: promptui.PipeCursor,
// 		Default:   defaultValue,
// 		Validate:  validator.IsValidURL,
// 		Templates: defaultTemplate(),
// 	}
//
// 	return prompt.Run()
// }

type surveyURL struct{}

func NewSurveyURL() surveyURL {
	return surveyURL{}
}

// URL show a prompt and parse to string.
func (surveyURL) URL(name, defaultValue string) (string, error) {
	var value string

	validationQs := []*survey.Question{
		{
			Name:     "name",
			Prompt:   &survey.Input{
				Message:  name,
				Default:  defaultValue,
			},
			Validate: validator.IsValidURL,
		},
	}

	return value, survey.Ask(validationQs, &name)
}
